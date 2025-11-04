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
	// Regex to match presentation frontmatter: <!-- presentation ... -->
	presentationFrontmatterRegex = regexp.MustCompile(`(?s)^\s*<!--\s*presentation\s+(.*?)\s*-->`)

	// Regex to match slide frontmatter: <!-- slide ... -->
	slideFrontmatterRegex = regexp.MustCompile(`(?s)<!--\s*slide\s+(.*?)\s*-->`)

	// Regex to match any HTML comment
	htmlCommentRegex = regexp.MustCompile(`(?s)<!--(.*?)-->`)
	
	// DeckSet per-slide directive regex: [.command: value]
	decksetDirectiveRegex = regexp.MustCompile(`(?m)^\s*\[\.([\w-]+):\s*([^\]]+)\]\s*$`)
	
	// DeckSet global config regex (key: value at start of file)
	decksetGlobalConfigRegex = regexp.MustCompile(`(?m)^([\w-]+):\s*(.+)\s*$`)
	
	// DeckSet speaker note regex: ^ Note text
	decksetNoteRegex = regexp.MustCompile(`(?m)^\^(.*)$`)
)

// Parser handles parsing markdown files into slides
type Parser struct {
	slides               []*Slide
	presentationMetadata PresentationMetadata
}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{
		slides:               make([]*Slide, 0),
		presentationMetadata: PresentationMetadata{},
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
	// Extract presentation-level frontmatter first
	content = p.extractPresentationFrontmatter(content)

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

// GetPresentationMetadata returns the presentation-level metadata
func (p *Parser) GetPresentationMetadata() PresentationMetadata {
	return p.presentationMetadata
}

// extractPresentationFrontmatter extracts and parses presentation-level YAML frontmatter
func (p *Parser) extractPresentationFrontmatter(content string) string {
	// First try HTML comment style (gobig native)
	matches := presentationFrontmatterRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		yamlContent := matches[1]

		// Parse YAML
		err := yaml.Unmarshal([]byte(yamlContent), &p.presentationMetadata)
		if err != nil {
			// If YAML parsing fails, just ignore the frontmatter
			fmt.Fprintf(os.Stderr, "Warning: failed to parse presentation metadata: %v\n", err)
		}

		// Remove frontmatter from content
		content = presentationFrontmatterRegex.ReplaceAllString(content, "")
		return content
	}
	
	// Try DeckSet global configuration format (key: value at top of file)
	content = p.extractDeckSetGlobalConfig(content)
	
	return content
}

// extractDeckSetGlobalConfig extracts DeckSet global configuration from the top of the file
func (p *Parser) extractDeckSetGlobalConfig(content string) string {
	lines := strings.Split(content, "\n")
	configLines := 0
	configMap := make(map[string]string)
	
	// Read configuration lines from the top of the file
	// They must be consecutive with no blank lines between them
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Stop at first blank line or non-config line
		if trimmed == "" {
			if configLines > 0 {
				break // Found config, now hit blank line
			}
			continue // Skip leading blank lines
		}
		
		// Check if it's a slide separator
		if isHorizontalRule(trimmed) {
			break
		}
		
		// Try to match config line
		if matches := decksetGlobalConfigRegex.FindStringSubmatch(line); len(matches) > 2 {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])
			configMap[key] = value
			configLines = i + 1
		} else if configLines > 0 {
			// Non-config line after we've started reading config - stop
			break
		} else {
			// First non-config, non-blank line - no global config
			break
		}
	}
	
	// If we found config, parse it and remove from content
	if configLines > 0 {
		// Build YAML from config map
		var yamlBuilder strings.Builder
		for key, value := range configMap {
			yamlBuilder.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
		
		// Parse YAML
		err := yaml.Unmarshal([]byte(yamlBuilder.String()), &p.presentationMetadata)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse DeckSet global config: %v\n", err)
		}
		
		// Remove config lines from content
		remainingLines := lines[configLines:]
		content = strings.Join(remainingLines, "\n")
	}
	
	return content
}

// parseSlide parses a single slide's content
func (p *Parser) parseSlide(content string) (*Slide, error) {
	slide := &Slide{
		Metadata: SlideMetadata{},
	}

	// Extract frontmatter (gobig style)
	content = p.extractFrontmatter(content, slide)

	// Extract DeckSet directives
	content = p.extractDeckSetDirectives(content, slide)
	
	// Process DeckSet image modifiers
	content = p.processDeckSetImages(content, slide)

	// Extract speaker notes (both styles)
	content = p.extractNotes(content, slide)
	content = p.extractDeckSetNotes(content, slide)

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
		if slide.Notes != "" {
			slide.Notes += "\n"
		}
		slide.Notes += strings.Join(notes, "\n")
	}

	// Remove all HTML comments from content
	content = htmlCommentRegex.ReplaceAllString(content, "")

	return content
}

// extractDeckSetDirectives extracts DeckSet per-slide directives [.command: value]
func (p *Parser) extractDeckSetDirectives(content string, slide *Slide) string {
	matches := decksetDirectiveRegex.FindAllStringSubmatch(content, -1)
	
	if len(matches) == 0 {
		return content
	}
	
	// Build YAML from directives
	var yamlBuilder strings.Builder
	for _, match := range matches {
		if len(match) > 2 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])
			yamlBuilder.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}
	
	// Parse YAML into slide metadata
	if yamlBuilder.Len() > 0 {
		err := yaml.Unmarshal([]byte(yamlBuilder.String()), &slide.Metadata)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse DeckSet directives: %v\n", err)
		}
	}
	
	// Remove directives from content
	content = decksetDirectiveRegex.ReplaceAllString(content, "")
	
	return content
}

// extractDeckSetNotes extracts DeckSet speaker notes (lines starting with ^)
func (p *Parser) extractDeckSetNotes(content string, slide *Slide) string {
	var notes []string
	
	matches := decksetNoteRegex.FindAllStringSubmatch(content, -1)
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
		if slide.Notes != "" {
			slide.Notes += "\n"
		}
		slide.Notes += strings.Join(notes, "\n")
	}
	
	// Remove DeckSet notes from content
	content = decksetNoteRegex.ReplaceAllString(content, "")
	
	return content
}

// processDeckSetImages processes DeckSet image modifiers and sets appropriate layout
func (p *Parser) processDeckSetImages(content string, slide *Slide) string {
	// Regex to match DeckSet image syntax: ![modifiers](path)
	decksetImageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	
	matches := decksetImageRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return content
	}
	
	// Check for left/right positioning which implies a layout
	hasLeft := false
	hasRight := false
	
	for _, match := range matches {
		if len(match) > 1 {
			modifiers := strings.ToLower(match[1])
			if strings.Contains(modifiers, "left") {
				hasLeft = true
			}
			if strings.Contains(modifiers, "right") {
				hasRight = true
			}
		}
	}
	
	// Auto-set layout based on image positions
	if slide.Metadata.Layout == "" {
		if hasLeft && hasRight {
			// Two images side by side
			slide.Metadata.Layout = "50-50"
		} else if hasLeft || hasRight {
			// One image on side, content on other
			if hasLeft {
				slide.Metadata.Layout = "50-50"
			} else {
				slide.Metadata.Layout = "50-50"
			}
		}
	}
	
	// Convert DeckSet image modifiers to HTML classes or remove them
	// For now, we'll strip the modifiers and keep the standard markdown syntax
	// The actual rendering will be handled by the layout
	content = decksetImageRegex.ReplaceAllStringFunc(content, func(match string) string {
		submatches := decksetImageRegex.FindStringSubmatch(match)
		if len(submatches) > 2 {
			// Keep the image but without special modifiers in alt text
			// modifiers := submatches[1]
			path := submatches[2]
			return fmt.Sprintf("![%s](%s)", "", path)
		}
		return match
	})
	
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
