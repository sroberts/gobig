package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	// Regex to match slide frontmatter: <!-- slide ... -->
	slideFrontmatterRegex = regexp.MustCompile(`(?s)<!--\s*slide\s+(.*?)\s*-->`)

	// Regex to match any HTML comment
	htmlCommentRegex = regexp.MustCompile(`(?s)<!--(.*?)-->`)
)

// Parser handles parsing markdown files into slides
type Parser struct {
	slides []*Slide
}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{
		slides: make([]*Slide, 0),
	}
}

// ParseFile reads and parses a markdown file
func (p *Parser) ParseFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return p.ParseString(string(content))
}

// ParseString parses markdown content from a string
func (p *Parser) ParseString(content string) error {
	// Split on horizontal rules (---)
	// We need to be careful to only split on standalone ---
	slideContents := splitOnHorizontalRule(content)

	for _, slideContent := range slideContents {
		slide, err := p.parseSlide(slideContent)
		if err != nil {
			return fmt.Errorf("failed to parse slide: %w", err)
		}

		// Only add non-empty slides
		if strings.TrimSpace(slide.Content) != "" {
			p.slides = append(p.slides, slide)
		}
	}

	if len(p.slides) == 0 {
		return fmt.Errorf("no slides found in input")
	}

	return nil
}

// GetSlides returns all parsed slides
func (p *Parser) GetSlides() []*Slide {
	return p.slides
}

// parseSlide parses a single slide's content
func (p *Parser) parseSlide(content string) (*Slide, error) {
	slide := &Slide{
		Metadata: SlideMetadata{},
	}

	// Extract frontmatter
	content = p.extractFrontmatter(content, slide)

	// Extract speaker notes
	content = p.extractNotes(content, slide)

	// Remaining content is the slide content
	slide.Content = strings.TrimSpace(content)

	return slide, nil
}

// extractFrontmatter extracts and parses YAML frontmatter from slide content
func (p *Parser) extractFrontmatter(content string, slide *Slide) string {
	matches := slideFrontmatterRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		yamlContent := matches[1]

		// Parse YAML
		err := yaml.Unmarshal([]byte(yamlContent), &slide.Metadata)
		if err != nil {
			// If YAML parsing fails, just ignore the frontmatter
			fmt.Fprintf(os.Stderr, "Warning: failed to parse slide metadata: %v\n", err)
		}

		// Remove frontmatter from content
		content = slideFrontmatterRegex.ReplaceAllString(content, "")
	}

	return content
}

// extractNotes extracts speaker notes from HTML comments
// This is called AFTER extractFrontmatter, so all remaining comments are notes
func (p *Parser) extractNotes(content string, slide *Slide) string {
	var notes []string

	matches := htmlCommentRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			noteContent := strings.TrimSpace(match[1])
			// Only add non-empty notes
			if noteContent != "" {
				notes = append(notes, noteContent)
			}
		}
	}

	if len(notes) > 0 {
		slide.Notes = strings.Join(notes, "\n")
	}

	// Remove all HTML comments from content
	content = htmlCommentRegex.ReplaceAllString(content, "")

	return content
}

// splitOnHorizontalRule splits content on standalone horizontal rules (---)
func splitOnHorizontalRule(content string) []string {
	var slides []string
	var currentSlide strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()

		// Check if line is a horizontal rule (3 or more dashes, possibly with spaces)
		trimmed := strings.TrimSpace(line)
		if isHorizontalRule(trimmed) {
			// Save current slide if it has content
			slideContent := currentSlide.String()
			if strings.TrimSpace(slideContent) != "" {
				slides = append(slides, slideContent)
			}
			currentSlide.Reset()
		} else {
			currentSlide.WriteString(line)
			currentSlide.WriteString("\n")
		}
	}

	// Add the last slide
	slideContent := currentSlide.String()
	if strings.TrimSpace(slideContent) != "" {
		slides = append(slides, slideContent)
	}

	return slides
}

// isHorizontalRule checks if a line is a horizontal rule
func isHorizontalRule(line string) bool {
	// Must be at least 3 dashes, with optional spaces
	if len(line) < 3 {
		return false
	}

	// Check if it's only dashes and spaces
	for _, ch := range line {
		if ch != '-' && ch != ' ' {
			return false
		}
	}

	// Count dashes
	dashCount := strings.Count(line, "-")
	return dashCount >= 3
}
