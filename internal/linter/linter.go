package linter

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gobig/internal/parser"
	"gopkg.in/yaml.v3"
)

// Issue represents a linting issue found in the markdown
type Issue struct {
	Line     int    // Line number where the issue was found (0 for general issues)
	Severity string // "error", "warning", "info"
	Message  string // Description of the issue
}

// Result contains the results of linting a markdown file
type Result struct {
	Issues   []Issue
	FilePath string
	HasError bool // True if any error-level issues were found
}

// Linter performs validation on presentation markdown files
type Linter struct {
	baseDir         string
	validLayouts    map[string]bool
	imageRegex      *regexp.Regexp
	yamlSlideRegex  *regexp.Regexp
	yamlPresentRegex *regexp.Regexp
}

// NewLinter creates a new linter instance
func NewLinter(baseDir string) *Linter {
	return &Linter{
		baseDir: baseDir,
		validLayouts: map[string]bool{
			"50-50":       true,
			"75-25":       true,
			"25-75":       true,
			"50-50-rows":  true,
			"75-25-rows":  true,
			"25-75-rows":  true,
			"grid-3x2":    true,
			"grid-2x3":    true,
		},
		imageRegex:       regexp.MustCompile(`!\[.*?\]\((.*?)\)`),
		yamlSlideRegex:   regexp.MustCompile(`(?s)<!--\s*slide\s+(.*?)\s*-->`),
		yamlPresentRegex: regexp.MustCompile(`(?s)^\s*<!--\s*presentation\s+(.*?)\s*-->`),
	}
}

// LintFile validates a markdown file and returns any issues found
func (l *Linter) LintFile(filename string) (*Result, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return l.LintString(string(content), filename)
}

// LintString validates markdown content and returns any issues found
func (l *Linter) LintString(content, filename string) (*Result, error) {
	result := &Result{
		Issues:   make([]Issue, 0),
		FilePath: filename,
	}

	// Create a parser to check basic structure
	p := parser.NewParser()
	parseErr := p.ParseString(content)

	// Check for parsing errors
	if parseErr != nil {
		result.Issues = append(result.Issues, Issue{
			Line:     0,
			Severity: "error",
			Message:  fmt.Sprintf("Failed to parse markdown: %v", parseErr),
		})
		result.HasError = true
		return result, nil
	}

	// Now perform detailed line-by-line validation
	lines := strings.Split(content, "\n")

	// Check presentation metadata YAML syntax
	l.checkPresentationMetadata(content, lines, result)

	// Check slide metadata YAML syntax
	l.checkSlideMetadata(content, lines, result)

	// Check each slide
	l.checkSlides(p.GetSlides(), lines, result)

	// Check for broken image references
	l.checkImageReferences(content, lines, result)

	// Check for malformed horizontal rules
	l.checkHorizontalRules(lines, result)

	return result, nil
}

// checkPresentationMetadata validates presentation-level metadata
func (l *Linter) checkPresentationMetadata(content string, lines []string, result *Result) {
	matches := l.yamlPresentRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		yamlContent := matches[1]
		var metadata parser.PresentationMetadata
		err := yaml.Unmarshal([]byte(yamlContent), &metadata)
		if err != nil {
			// Find the line number of the presentation metadata
			lineNum := l.findLineNumber(lines, "<!-- presentation")
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Severity: "error",
				Message:  fmt.Sprintf("Invalid presentation metadata YAML: %v", err),
			})
			result.HasError = true
		}
	}
}

// checkSlideMetadata validates slide-level metadata YAML syntax
func (l *Linter) checkSlideMetadata(content string, lines []string, result *Result) {
	matches := l.yamlSlideRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			yamlContent := match[1]
			var metadata parser.SlideMetadata
			err := yaml.Unmarshal([]byte(yamlContent), &metadata)
			if err != nil {
				// Find the line number of this slide metadata block
				lineNum := l.findLineNumber(lines, match[0])
				result.Issues = append(result.Issues, Issue{
					Line:     lineNum,
					Severity: "error",
					Message:  fmt.Sprintf("Invalid slide metadata YAML: %v", err),
				})
				result.HasError = true
			}
		}
	}
}

// checkSlides validates individual slides
func (l *Linter) checkSlides(slides []*parser.Slide, lines []string, result *Result) {
	if len(slides) == 0 {
		result.Issues = append(result.Issues, Issue{
			Line:     0,
			Severity: "error",
			Message:  "No slides found in presentation",
		})
		result.HasError = true
		return
	}

	slideNum := 0
	inSlide := true // Start in the first slide
	currentSlideStart := 1

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Detect slide boundaries
		if l.isHorizontalRule(trimmed) {
			if inSlide && slideNum < len(slides) {
				l.checkSlideContent(slides[slideNum], currentSlideStart, lineNum-1, lines, result)
				slideNum++
			}
			inSlide = true
			currentSlideStart = lineNum + 1
		}
	}

	// Check the last slide
	if slideNum < len(slides) {
		l.checkSlideContent(slides[slideNum], currentSlideStart, len(lines), lines, result)
	}
}

// checkSlideContent validates a single slide's content
func (l *Linter) checkSlideContent(slide *parser.Slide, startLine, endLine int, lines []string, result *Result) {
	// Check layout validity
	layout := slide.Metadata.Layout
	if layout != "" && !l.validLayouts[layout] && !l.isCSSLayout(layout) {
		// Find the line with the layout metadata
		lineNum := l.findMetadataLine(lines, startLine, endLine, "layout:")
		result.Issues = append(result.Issues, Issue{
			Line:     lineNum,
			Severity: "warning",
			Message:  fmt.Sprintf("Unknown layout '%s'. Valid layouts: %s", layout, l.getValidLayoutsString()),
		})
	}

	// Check for negative auto-advance times
	if slide.Metadata.TimeToNext < 0 {
		lineNum := l.findMetadataLine(lines, startLine, endLine, "time-to-next:")
		result.Issues = append(result.Issues, Issue{
			Line:     lineNum,
			Severity: "error",
			Message:  "time-to-next cannot be negative",
		})
		result.HasError = true
	}
}

// isCSSLayout returns true if the layout string looks like a CSS declaration
// (e.g., "grid-template-columns: 1fr 2fr;"). Such layouts are allowed by the
// generator and should not be treated as unknown named layouts.
func (l *Linter) isCSSLayout(layout string) bool {
	return strings.Contains(layout, ":") || strings.Contains(layout, ";")
}

// checkImageReferences validates that local image files exist
func (l *Linter) checkImageReferences(content string, lines []string, result *Result) {
	// Iterate line by line so we can accurately report the line number
	for i, line := range lines {
		matches := l.imageRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				imagePath := match[1]

				// Skip external URLs
				if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") || strings.HasPrefix(imagePath, "data:") {
					continue
				}

				// Check if local file exists
				fullPath := filepath.Join(l.baseDir, imagePath)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					lineNum := i + 1
					result.Issues = append(result.Issues, Issue{
						Line:     lineNum,
						Severity: "error",
						Message:  fmt.Sprintf("Image file not found: %s", imagePath),
					})
					result.HasError = true
				}
			}
		}
	}
}

// checkHorizontalRules validates horizontal rule syntax
func (l *Linter) checkHorizontalRules(lines []string, result *Result) {
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip HTML comment closers (-->)
		if trimmed == "-->" {
			continue
		}

		// Check for potential malformed horizontal rules
		if strings.HasPrefix(trimmed, "--") && !l.isHorizontalRule(trimmed) {
			lineNum := i + 1
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Severity: "warning",
				Message:  "Potential malformed horizontal rule. Use at least 3 dashes (---) on a line by themselves",
			})
		}
	}
}

// isHorizontalRule checks if a line is a valid horizontal rule
func (l *Linter) isHorizontalRule(line string) bool {
	if len(line) < 3 {
		return false
	}

	for _, ch := range line {
		if ch != '-' && ch != ' ' {
			return false
		}
	}

	dashCount := strings.Count(line, "-")
	return dashCount >= 3
}

// findLineNumber finds the line number containing a specific string
func (l *Linter) findLineNumber(lines []string, searchStr string) int {
	for i, line := range lines {
		if strings.Contains(line, searchStr) {
			return i + 1
		}
	}
	return 0
}

// findMetadataLine finds the line number of a specific metadata key within a range
func (l *Linter) findMetadataLine(lines []string, startLine, endLine int, key string) int {
	for i := startLine - 1; i < endLine && i < len(lines); i++ {
		if strings.Contains(lines[i], key) {
			return i + 1
		}
	}
	return startLine
}

// getValidLayoutsString returns a comma-separated list of valid layouts
func (l *Linter) getValidLayoutsString() string {
	layouts := make([]string, 0, len(l.validLayouts))
	for layout := range l.validLayouts {
		layouts = append(layouts, layout)
	}
	return strings.Join(layouts, ", ")
}

// FormatIssues formats the linting issues for display
func (r *Result) FormatIssues() string {
	if len(r.Issues) == 0 {
		return "No issues found"
	}

	var output strings.Builder
	for _, issue := range r.Issues {
		if issue.Line > 0 {
			fmt.Fprintf(&output, "%s:%d: [%s] %s\n", r.FilePath, issue.Line, issue.Severity, issue.Message)
		} else {
			fmt.Fprintf(&output, "%s: [%s] %s\n", r.FilePath, issue.Severity, issue.Message)
		}
	}
	return output.String()
}

// GetSummary returns a summary of linting results
func (r *Result) GetSummary() string {
	if len(r.Issues) == 0 {
		return "✓ No issues found"
	}

	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, issue := range r.Issues {
		switch issue.Severity {
		case "error":
			errorCount++
		case "warning":
			warningCount++
		case "info":
			infoCount++
		}
	}

	return fmt.Sprintf("Found %d issue(s): %d error(s), %d warning(s), %d info(s)",
		len(r.Issues), errorCount, warningCount, infoCount)
}
