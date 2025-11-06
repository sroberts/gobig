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
