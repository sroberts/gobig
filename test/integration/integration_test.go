package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gobig/internal/generator"
	"gobig/internal/parser"
)

// TestParserGeneratorIntegration tests the full pipeline from parsing to generation
func TestParserGeneratorIntegration(t *testing.T) {
	tests := []struct {
		name         string
		markdown     string
		theme        string
		wantSlides   int
		checkContent func(t *testing.T, html string)
	}{
		{
			name: "basic two slides",
			markdown: `# First Slide

Content 1

---

# Second Slide

Content 2`,
			theme:      "dark",
			wantSlides: 2,
			checkContent: func(t *testing.T, html string) {
				if !strings.Contains(html, "First Slide") {
					t.Error("Expected HTML to contain 'First Slide'")
				}
				if !strings.Contains(html, "Second Slide") {
					t.Error("Expected HTML to contain 'Second Slide'")
				}
				if !strings.Contains(html, "Content 1") {
					t.Error("Expected HTML to contain 'Content 1'")
				}
				if !strings.Contains(html, "Content 2") {
					t.Error("Expected HTML to contain 'Content 2'")
				}
			},
		},
		{
			name: "with metadata and notes",
			markdown: `<!-- presentation
title: Test Presentation
theme: light
-->

<!-- slide
layout: 50-50
-->

# Slide 1

Content here

<!-- This is a speaker note -->`,
			theme:      "light",
			wantSlides: 1,
			checkContent: func(t *testing.T, html string) {
				if !strings.Contains(html, "Slide 1") {
					t.Error("Expected HTML to contain 'Slide 1'")
				}
				if !strings.Contains(html, "Content here") {
					t.Error("Expected HTML to contain 'Content here'")
				}
				// Notes should be in the HTML but hidden
				if !strings.Contains(html, "This is a speaker note") {
					t.Error("Expected HTML to contain speaker notes")
				}
			},
		},
		{
			name: "complex formatting",
			markdown: `# Code Example

` + "```go\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```" + `

---

# Lists and Links

- Item 1
- Item 2
- [Link](https://example.com)

---

# Images

![Test](test.png)`,
			theme:      "white",
			wantSlides: 3,
			checkContent: func(t *testing.T, html string) {
				if !strings.Contains(html, "Code Example") {
					t.Error("Expected HTML to contain 'Code Example'")
				}
				if !strings.Contains(html, "func main()") {
					t.Error("Expected HTML to contain code")
				}
				if !strings.Contains(html, "Lists and Links") {
					t.Error("Expected HTML to contain 'Lists and Links'")
				}
				if !strings.Contains(html, "Item 1") {
					t.Error("Expected HTML to contain list items")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse
			p := parser.NewParser()
			if err := p.ParseString(tt.markdown); err != nil {
				t.Fatalf("ParseString() error = %v", err)
			}

			slides := p.GetSlides()
			if len(slides) != tt.wantSlides {
				t.Errorf("Got %d slides, want %d", len(slides), tt.wantSlides)
			}

			// Generate
			opts := generator.Options{
				Theme:                tt.theme,
				Title:                "",
				AspectRatio:          "1.6",
				BasePath:             "",
				SkipBase64Images:     false,
				PresentationMetadata: p.GetPresentationMetadata(),
			}

			gen := generator.NewGenerator(opts)
			html, err := gen.Generate(slides)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			// Check content
			if tt.checkContent != nil {
				tt.checkContent(t, html)
			}

			// Verify HTML structure
			if !strings.Contains(html, "<!DOCTYPE html>") {
				t.Error("Expected HTML to contain DOCTYPE")
			}
			if !strings.Contains(html, "<html") {
				t.Error("Expected HTML to contain <html tag")
			}
			if !strings.Contains(html, "</html>") {
				t.Error("Expected HTML to contain closing </html tag")
			}
		})
	}
}

// TestEndToEndFileWorkflow tests the complete workflow with real files
func TestEndToEndFileWorkflow(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	mdFile := filepath.Join(tmpDir, "presentation.md")
	imgFile := filepath.Join(tmpDir, "diagram.png")

	markdown := `# Welcome

![Diagram](diagram.png)

---

# Code

` + "```python\nprint('hello')\n```"

	// Write files
	if err := os.WriteFile(mdFile, []byte(markdown), 0644); err != nil {
		t.Fatalf("Failed to write markdown: %v", err)
	}
	if err := os.WriteFile(imgFile, []byte("fake image data"), 0644); err != nil {
		t.Fatalf("Failed to write image: %v", err)
	}

	// Parse file
	p := parser.NewParser()
	if err := p.ParseFile(mdFile); err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	slides := p.GetSlides()
	if len(slides) != 2 {
		t.Errorf("Got %d slides, want 2", len(slides))
	}

	// Generate HTML with base path
	basePath, _ := filepath.Abs(tmpDir)
	opts := generator.Options{
		Theme:                "dark",
		Title:                "Test",
		AspectRatio:          "1.6",
		BasePath:             basePath,
		SkipBase64Images:     false,
		PresentationMetadata: p.GetPresentationMetadata(),
	}

	gen := generator.NewGenerator(opts)
	html, err := gen.Generate(slides)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify content
	if !strings.Contains(html, "Welcome") {
		t.Error("Expected HTML to contain 'Welcome'")
	}
	if !strings.Contains(html, "print('hello')") {
		t.Error("Expected HTML to contain code")
	}
}

// TestMultipleThemes verifies that all themes work with the same content
func TestMultipleThemes(t *testing.T) {
	markdown := `# Test Slide

Test content.`

	p := parser.NewParser()
	if err := p.ParseString(markdown); err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	themes := []string{"dark", "light", "white"}

	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			opts := generator.Options{
				Theme:                theme,
				Title:                "",
				AspectRatio:          "1.6",
				BasePath:             "",
				SkipBase64Images:     false,
				PresentationMetadata: p.GetPresentationMetadata(),
			}

			gen := generator.NewGenerator(opts)
			html, err := gen.Generate(p.GetSlides())
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			if !strings.Contains(html, "Test Slide") {
				t.Error("Expected HTML to contain 'Test Slide'")
			}
			if !strings.Contains(html, "Test content") {
				t.Error("Expected HTML to contain 'Test content'")
			}
		})
	}
}

// TestLayouts tests various slide layouts
func TestLayouts(t *testing.T) {
	layouts := []struct {
		name   string
		layout string
	}{
		{"50-50 columns", "50-50"},
		{"75-25 columns", "75-25"},
		{"25-75 columns", "25-75"},
		{"50-50 rows", "50-50-rows"},
		{"75-25 rows", "75-25-rows"},
		{"grid 3x2", "grid-3x2"},
	}

	for _, tt := range layouts {
		t.Run(tt.name, func(t *testing.T) {
			markdown := `<!-- slide
layout: ` + tt.layout + `
-->

# ` + tt.name + `

Content here.`

			p := parser.NewParser()
			if err := p.ParseString(markdown); err != nil {
				t.Fatalf("ParseString() error = %v", err)
			}

			opts := generator.Options{
				Theme:                "dark",
				Title:                "",
				AspectRatio:          "1.6",
				BasePath:             "",
				SkipBase64Images:     false,
				PresentationMetadata: p.GetPresentationMetadata(),
			}

			gen := generator.NewGenerator(opts)
			html, err := gen.Generate(p.GetSlides())
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			if !strings.Contains(html, tt.name) {
				t.Errorf("Expected HTML to contain '%s'", tt.name)
			}
		})
	}
}

// TestPresentationMetadataOverrides tests that presentation metadata is respected
func TestPresentationMetadataOverrides(t *testing.T) {
	markdown := `<!-- presentation
title: Custom Title
auto-advance: 5000
-->

# Slide 1

Content`

	p := parser.NewParser()
	if err := p.ParseString(markdown); err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	metadata := p.GetPresentationMetadata()
	if metadata.Title != "Custom Title" {
		t.Errorf("Expected title 'Custom Title', got '%v'", metadata.Title)
	}

	opts := generator.Options{
		Theme:                "dark",
		Title:                "Override Title", // Should be ignored if metadata has title
		AspectRatio:          "1.6",
		BasePath:             "",
		SkipBase64Images:     false,
		PresentationMetadata: metadata,
	}

	gen := generator.NewGenerator(opts)
	html, err := gen.Generate(p.GetSlides())
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// The HTML should contain the content
	if !strings.Contains(html, "Slide 1") {
		t.Error("Expected HTML to contain 'Slide 1'")
	}
}
