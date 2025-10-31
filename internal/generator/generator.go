package generator

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"gobig/internal/assets"
	parserPkg "gobig/internal/parser"
)

// Options contains configuration for HTML generation
type Options struct {
	Theme                string                         // "dark", "light", or "white"
	Title                string                         // Presentation title
	AspectRatio          string                         // Aspect ratio (e.g., "1.6", "2", "false")
	BasePath             string                         // Base path for resolving relative image paths
	PresentationMetadata parserPkg.PresentationMetadata // Presentation-level metadata
}

// Generator handles HTML generation from parsed slides
type Generator struct {
	options Options
	md      goldmark.Markdown
}

// NewGenerator creates a new generator with the given options
func NewGenerator(opts Options) *Generator {
	// Set default theme if not specified
	if opts.Theme == "" {
		opts.Theme = "dark"
	}

	// Create goldmark markdown processor
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,   // GitHub Flavored Markdown
			extension.Table, // Tables
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // Allow raw HTML
		),
	)

	return &Generator{
		options: opts,
		md:      md,
	}
}

// Generate creates the final HTML output from slides
func (g *Generator) Generate(slides []*parserPkg.Slide) (string, error) {
	// Get embedded assets
	bigJS, err := assets.GetBigJS()
	if err != nil {
		return "", fmt.Errorf("failed to get big.js: %w", err)
	}

	bigCSS, err := assets.GetBigCSS()
	if err != nil {
		return "", fmt.Errorf("failed to get big.css: %w", err)
	}

	themeCSS, err := assets.GetTheme(g.options.Theme)
	if err != nil {
		return "", fmt.Errorf("failed to get theme: %w", err)
	}

	// Generate slides HTML
	slidesHTML := g.generateSlides(slides)

	// Determine title (use first slide's text if not specified)
	title := g.options.Title
	if title == "" && len(slides) > 0 {
		title = extractTitle(slides[0].Content)
	}
	if title == "" {
		title = "Presentation"
	}

	// Generate aspect ratio script
	aspectRatioScript := aspectRatioScript(g.options.AspectRatio)

	// Generate final HTML
	html := generateHTML(
		title,
		bigCSS,
		themeCSS,
		aspectRatioScript,
		bigJS,
		g.options.Theme,
		slidesHTML,
	)

	return html, nil
}

// generateSlides converts all slides to HTML
func (g *Generator) generateSlides(slides []*parserPkg.Slide) string {
	var sb strings.Builder

	for _, slide := range slides {
		slideHTML := g.generateSlide(slide)
		sb.WriteString(slideHTML)
		sb.WriteString("\n")
	}

	return sb.String()
}

// generateSlide converts a single slide to HTML
func (g *Generator) generateSlide(slide *parserPkg.Slide) string {
	var sb strings.Builder

	// Start slide div with optional attributes
	sb.WriteString("  <div")

	// Add data attributes
	// Determine time-to-next: slide-level overrides presentation-level
	timeToNext := slide.Metadata.TimeToNext
	if timeToNext == 0 && g.options.PresentationMetadata.TimeToNext > 0 {
		timeToNext = g.options.PresentationMetadata.TimeToNext
	}
	if timeToNext > 0 {
		sb.WriteString(fmt.Sprintf(` data-time-to-next="%d"`, timeToNext))
	}

	if slide.Metadata.BodyStyle != "" {
		sb.WriteString(fmt.Sprintf(` data-body-style="%s"`, escapeAttr(slide.Metadata.BodyStyle)))
	}
	if slide.Metadata.BodyClass != "" {
		sb.WriteString(fmt.Sprintf(` data-body-class="%s"`, escapeAttr(slide.Metadata.BodyClass)))
	}

	sb.WriteString(">")

	// Handle layouts
	if slide.Metadata.Layout != "" {
		sb.WriteString(g.generateLayoutSlide(slide))
	} else {
		// Regular slide - convert markdown to HTML
		html := g.markdownToHTML(slide.Content)
		sb.WriteString(html)
	}

	// Add speaker notes if present
	if slide.Notes != "" {
		sb.WriteString(fmt.Sprintf("\n    <notes>%s</notes>", escapeHTML(slide.Notes)))
	}

	sb.WriteString("\n  </div>")

	return sb.String()
}

// generateLayoutSlide generates a slide with CSS Grid layout
func (g *Generator) generateLayoutSlide(slide *parserPkg.Slide) string {
	gridStyle := layoutToGridStyle(slide.Metadata.Layout)

	// Split content by image/text blocks
	parts := splitContentForLayout(slide.Content)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n    <div class=\"layout\" style=\"%s\">", gridStyle))

	for _, part := range parts {
		html := g.markdownToHTML(part)
		sb.WriteString("\n      ")
		sb.WriteString(html)
	}

	sb.WriteString("\n    </div>")

	return sb.String()
}

// markdownToHTML converts markdown to HTML
func (g *Generator) markdownToHTML(markdown string) string {
	var buf bytes.Buffer
	if err := g.md.Convert([]byte(markdown), &buf); err != nil {
		return markdown // Fallback to raw content
	}

	html := buf.String()

	// Process images for base64 encoding (for single-file output)
	html = g.processImages(html)

	return strings.TrimSpace(html)
}

// processImages converts local image paths to base64 data URIs
func (g *Generator) processImages(html string) string {
	if g.options.BasePath == "" {
		return html
	}

	// Regex to find image tags
	imgRegex := regexp.MustCompile(`<img[^>]+src="([^"]+)"[^>]*>`)

	return imgRegex.ReplaceAllStringFunc(html, func(match string) string {
		// Extract src attribute
		srcRegex := regexp.MustCompile(`src="([^"]+)"`)
		srcMatches := srcRegex.FindStringSubmatch(match)

		if len(srcMatches) < 2 {
			return match
		}

		src := srcMatches[1]

		// Skip URLs (http://, https://, data:, etc.)
		if strings.HasPrefix(src, "http://") ||
			strings.HasPrefix(src, "https://") ||
			strings.HasPrefix(src, "data:") {
			return match
		}

		// Try to read and encode the image
		imagePath := filepath.Join(g.options.BasePath, src)
		data, err := os.ReadFile(imagePath)
		if err != nil {
			// If file doesn't exist, return original
			return match
		}

		// Detect content type
		contentType := detectContentType(imagePath)

		// Encode to base64
		encoded := base64.StdEncoding.EncodeToString(data)
		dataURI := fmt.Sprintf("data:%s;base64,%s", contentType, encoded)

		// Replace src
		return srcRegex.ReplaceAllString(match, fmt.Sprintf(`src="%s"`, dataURI))
	})
}

// detectContentType detects the MIME type from file extension
func detectContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// layoutToGridStyle converts layout name to CSS Grid style
func layoutToGridStyle(layout string) string {
	switch layout {
	case "50-50":
		return "grid-template-columns: 50% 50%;"
	case "75-25":
		return "grid-template-columns: 75% 25%;"
	case "25-75":
		return "grid-template-columns: 25% 75%;"
	case "75-25-rows":
		return "grid-template-rows: 75% 25%;"
	case "25-75-rows":
		return "grid-template-rows: 25% 75%;"
	case "50-50-rows":
		return "grid-template-rows: 50% 50%;"
	case "grid-3x2":
		return "grid-template-columns: repeat(3, 1fr); grid-template-rows: repeat(2, 1fr);"
	case "grid-2x3":
		return "grid-template-columns: repeat(2, 1fr); grid-template-rows: repeat(3, 1fr);"
	default:
		// Allow custom grid styles
		return layout
	}
}

// splitContentForLayout splits content into parts for layout
// Each paragraph or image becomes a grid item
func splitContentForLayout(content string) []string {
	var parts []string

	// Split on double newlines or image boundaries
	sections := strings.Split(content, "\n\n")

	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section != "" {
			parts = append(parts, section)
		}
	}

	return parts
}

// extractTitle extracts a title from markdown content
func extractTitle(markdown string) string {
	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

// escapeHTML escapes HTML special characters
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// escapeAttr escapes HTML attribute values
func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
