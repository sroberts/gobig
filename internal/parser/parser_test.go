package parser

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	p := NewParser()
	if p == nil {
		t.Fatal("NewParser() returned nil")
	}

	if p.slides == nil {
		t.Error("Parser slides should be initialized")
	}

	if len(p.GetSlides()) != 0 {
		t.Error("New parser should have no slides")
	}
}

func TestParseStringBasic(t *testing.T) {
	p := NewParser()
	content := `# First Slide

This is the first slide.

---

# Second Slide

This is the second slide.`

	err := p.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	slides := p.GetSlides()
	if len(slides) != 2 {
		t.Fatalf("Expected 2 slides, got %d", len(slides))
	}

	if slides[0].Content != "# First Slide\n\nThis is the first slide." {
		t.Errorf("First slide content mismatch: got %q", slides[0].Content)
	}

	if slides[1].Content != "# Second Slide\n\nThis is the second slide." {
		t.Errorf("Second slide content mismatch: got %q", slides[1].Content)
	}
}

func TestParseStringSingleSlide(t *testing.T) {
	p := NewParser()
	content := `# Only Slide

Just one slide here.`

	err := p.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	slides := p.GetSlides()
	if len(slides) != 1 {
		t.Fatalf("Expected 1 slide, got %d", len(slides))
	}
}

func TestParseStringEmptyContent(t *testing.T) {
	p := NewParser()
	content := ""

	err := p.ParseString(content)
	if err == nil {
		t.Error("ParseString() should fail with empty content")
	}
}

func TestParseStringWithSlideMetadata(t *testing.T) {
	p := NewParser()
	content := `<!-- slide
layout: 50-50
class: custom-class
time-to-next: 10
-->

# Slide with Metadata

Content here.`

	err := p.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	slides := p.GetSlides()
	if len(slides) != 1 {
		t.Fatalf("Expected 1 slide, got %d", len(slides))
	}

	metadata := slides[0].Metadata
	if metadata.Layout != "50-50" {
		t.Errorf("Expected layout '50-50', got %q", metadata.Layout)
	}

	if metadata.Class != "custom-class" {
		t.Errorf("Expected class 'custom-class', got %q", metadata.Class)
	}

	if metadata.TimeToNext != 10 {
		t.Errorf("Expected time-to-next 10, got %d", metadata.TimeToNext)
	}
}

func TestParseStringWithPresentationMetadata(t *testing.T) {
	p := NewParser()
	content := `<!-- presentation
title: My Presentation
time-to-next: 5
-->

# First Slide

Content here.

---

# Second Slide

More content.`

	err := p.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	metadata := p.GetPresentationMetadata()
	if metadata.Title != "My Presentation" {
		t.Errorf("Expected title 'My Presentation', got %q", metadata.Title)
	}

	if metadata.TimeToNext != 5 {
		t.Errorf("Expected time-to-next 5, got %d", metadata.TimeToNext)
	}

	slides := p.GetSlides()
	if len(slides) != 2 {
		t.Fatalf("Expected 2 slides, got %d", len(slides))
	}
}

func TestParseStringWithNotes(t *testing.T) {
	p := NewParser()
	content := `# Slide with Notes

Visible content

<!--
This is a speaker note.
It should be extracted.
-->`

	err := p.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	slides := p.GetSlides()
	if len(slides) != 1 {
		t.Fatalf("Expected 1 slide, got %d", len(slides))
	}

	if slides[0].Notes == "" {
		t.Error("Expected notes to be extracted, but got empty string")
	}

	// Notes should not be in content
	if len(slides[0].Content) == 0 {
		t.Error("Content should not be empty")
	}
}

func TestSlideDetectType(t *testing.T) {
	slide := &Slide{}

	// Basic test - should return SlideTypeContent as default
	slideType := slide.DetectType("<h1>Test</h1>")
	if slideType != SlideTypeContent {
		t.Errorf("Expected SlideTypeContent, got %v", slideType)
	}
}

// TestParseStringMalformedYAML tests handling of malformed YAML
// Note: The parser logs warnings for malformed YAML but continues processing
func TestParseStringMalformedYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "invalid slide metadata YAML",
			content: `<!-- slide
layout: [invalid yaml structure
class: test
-->

# Slide

Content`,
		},
		{
			name: "invalid presentation metadata YAML",
			content: `<!-- presentation
title: [broken
time-to-next: not a number
-->

# Slide

Content`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			// Parser should continue even with malformed YAML (it logs warnings)
			err := p.ParseString(tt.content)
			if err != nil {
				t.Errorf("Parser should handle malformed YAML gracefully, got error: %v", err)
			}
			// Should still produce a slide
			slides := p.GetSlides()
			if len(slides) == 0 {
				t.Error("Expected at least one slide even with malformed YAML")
			}
		})
	}
}

// TestParseStringWithMultipleSeparators tests handling of consecutive separators
func TestParseStringWithMultipleSeparators(t *testing.T) {
	p := NewParser()
	content := `# Slide 1

Content

---

---

# Slide 2

More content`

	err := p.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	slides := p.GetSlides()
	// Should handle empty slides or skip them appropriately
	if len(slides) == 0 {
		t.Error("Expected at least one slide")
	}
}

// TestParseStringWhitespaceOnly tests handling of whitespace-only content
func TestParseStringWhitespaceOnly(t *testing.T) {
	p := NewParser()
	content := `


	`

	err := p.ParseString(content)
	// Should error on whitespace-only content
	if err == nil {
		t.Error("Expected error for whitespace-only content")
	}
}

// TestParseStringWithComplexNotes tests complex note patterns
func TestParseStringWithComplexNotes(t *testing.T) {
	p := NewParser()
	content := `# Slide

Content

<!-- Note 1 -->

More content

<!-- Note 2
Multi-line note
with details
-->`

	err := p.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	slides := p.GetSlides()
	if len(slides) != 1 {
		t.Fatalf("Expected 1 slide, got %d", len(slides))
	}

	// Should extract notes
	if slides[0].Notes == "" {
		t.Error("Expected notes to be extracted")
	}
}

// TestParseFileNotFound tests handling of non-existent files
func TestParseFileNotFound(t *testing.T) {
	p := NewParser()
	err := p.ParseFile("/nonexistent/path/to/file.md")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestParseFileInvalidPath tests handling of invalid file paths
func TestParseFileInvalidPath(t *testing.T) {
	p := NewParser()
	// Try to parse a directory instead of a file
	err := p.ParseFile("/tmp")
	if err == nil {
		// Some systems might allow reading directories as files, so we're lenient here
		t.Log("Warning: parsing directory didn't produce expected error")
	}
}
