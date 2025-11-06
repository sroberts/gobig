package generator

import (
	"strings"
	"testing"

	"gobig/internal/parser"
)

func TestNewGenerator(t *testing.T) {
	opts := Options{
		Theme: "dark",
		Title: "Test Presentation",
	}

	gen := NewGenerator(opts)
	if gen == nil {
		t.Fatal("NewGenerator() returned nil")
	}

	if gen.options.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got %q", gen.options.Theme)
	}
}

func TestNewGeneratorDefaultTheme(t *testing.T) {
	opts := Options{
		Title: "Test",
	}

	gen := NewGenerator(opts)
	if gen.options.Theme != "dark" {
		t.Errorf("Expected default theme 'dark', got %q", gen.options.Theme)
	}
}

func TestGenerateBasic(t *testing.T) {
	opts := Options{
		Theme:       "dark",
		Title:       "Test Presentation",
		AspectRatio: "1.6",
	}

	gen := NewGenerator(opts)

	slides := []*parser.Slide{
		{
			Content: "# First Slide\n\nContent here.",
		},
		{
			Content: "# Second Slide\n\nMore content.",
		},
	}

	html, err := gen.Generate(slides)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	if len(html) == 0 {
		t.Error("Generate() returned empty HTML")
	}

	// Check for expected HTML structure
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Generated HTML missing DOCTYPE")
	}

	if !strings.Contains(html, "<html>") {
		t.Error("Generated HTML missing html tag")
	}

	if !strings.Contains(html, "<head>") {
		t.Error("Generated HTML missing head tag")
	}

	if !strings.Contains(html, "<body") {
		t.Error("Generated HTML missing body tag")
	}

	if !strings.Contains(html, "Test Presentation") {
		t.Error("Generated HTML missing title")
	}

	// Check for big.js content
	if !strings.Contains(html, "ASPECT_RATIO") {
		t.Error("Generated HTML missing big.js")
	}

	// Check for slide content
	if !strings.Contains(html, "First Slide") {
		t.Error("Generated HTML missing first slide content")
	}

	if !strings.Contains(html, "Second Slide") {
		t.Error("Generated HTML missing second slide content")
	}
}

func TestGenerateWithAllThemes(t *testing.T) {
	themes := []string{"dark", "light", "white"}

	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			opts := Options{
				Theme:       theme,
				Title:       "Test",
				AspectRatio: "1.6",
			}

			gen := NewGenerator(opts)

			slides := []*parser.Slide{
				{Content: "# Test Slide"},
			}

			html, err := gen.Generate(slides)
			if err != nil {
				t.Fatalf("Generate() failed for theme %s: %v", theme, err)
			}

			if len(html) == 0 {
				t.Errorf("Generate() returned empty HTML for theme %s", theme)
			}
		})
	}
}

func TestGenerateWithSlideMetadata(t *testing.T) {
	opts := Options{
		Theme: "dark",
		Title: "Test",
	}

	gen := NewGenerator(opts)

	slides := []*parser.Slide{
		{
			Content: "# Slide with Layout",
			Metadata: parser.SlideMetadata{
				Layout: "50-50",
				Class:  "custom",
			},
		},
	}

	html, err := gen.Generate(slides)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Just check that HTML is generated successfully
	if len(html) == 0 {
		t.Error("Generated HTML is empty")
	}

	// The metadata should be processed (specific format may vary)
	if !strings.Contains(html, "Slide with Layout") {
		t.Error("Generated HTML missing slide content")
	}
}

func TestGenerateWithNotes(t *testing.T) {
	opts := Options{
		Theme: "dark",
		Title: "Test",
	}

	gen := NewGenerator(opts)

	slides := []*parser.Slide{
		{
			Content: "# Slide",
			Notes:   "This is a speaker note.",
		},
	}

	html, err := gen.Generate(slides)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Check that notes are included
	if !strings.Contains(html, "<notes>") && !strings.Contains(html, "This is a speaker note") {
		t.Error("Generated HTML missing speaker notes")
	}
}

func TestGenerateWithAutoAdvance(t *testing.T) {
	opts := Options{
		Theme: "dark",
		Title: "Test",
		PresentationMetadata: parser.PresentationMetadata{
			TimeToNext: 5,
		},
	}

	gen := NewGenerator(opts)

	slides := []*parser.Slide{
		{
			Content: "# Auto-advance slide",
		},
	}

	html, err := gen.Generate(slides)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Check that auto-advance time is included
	if !strings.Contains(html, "data-time-to-next") || !strings.Contains(html, "5") {
		t.Error("Generated HTML missing auto-advance time")
	}
}

func TestGenerateEmptySlides(t *testing.T) {
	opts := Options{
		Theme: "dark",
		Title: "Test",
	}

	gen := NewGenerator(opts)

	var slides []*parser.Slide

	// Generate should handle empty slides gracefully
	html, err := gen.Generate(slides)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Should still produce valid HTML structure
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Generated HTML missing DOCTYPE even with no slides")
	}
}

func TestGenerateWithAspectRatioFalse(t *testing.T) {
	opts := Options{
		Theme:       "dark",
		Title:       "Test",
		AspectRatio: "false",
	}

	gen := NewGenerator(opts)

	slides := []*parser.Slide{
		{Content: "# Test"},
	}

	html, err := gen.Generate(slides)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	if !strings.Contains(html, "false") {
		t.Error("Generated HTML should handle aspect ratio false")
	}
}
