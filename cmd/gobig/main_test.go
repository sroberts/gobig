package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateFromFile tests the generateFromFile helper function
func TestGenerateFromFile(t *testing.T) {
	// Create a temporary markdown file
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")

	content := `# Hello World

This is a test slide.

---

# Slide 2

Another slide.`

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name         string
		inputFile    string
		skipBase64   bool
		wantErr      bool
		checkContent func(t *testing.T, html string)
	}{
		{
			name:       "basic generation",
			inputFile:  mdFile,
			skipBase64: false,
			wantErr:    false,
			checkContent: func(t *testing.T, html string) {
				if !strings.Contains(html, "Hello World") {
					t.Error("Expected HTML to contain 'Hello World'")
				}
				if !strings.Contains(html, "Slide 2") {
					t.Error("Expected HTML to contain 'Slide 2'")
				}
				if !strings.Contains(html, "<html") {
					t.Error("Expected HTML to contain <html tag")
				}
			},
		},
		{
			name:       "skip base64 images",
			inputFile:  mdFile,
			skipBase64: true,
			wantErr:    false,
			checkContent: func(t *testing.T, html string) {
				if !strings.Contains(html, "Hello World") {
					t.Error("Expected HTML to contain 'Hello World'")
				}
			},
		},
		{
			name:       "nonexistent file",
			inputFile:  filepath.Join(tmpDir, "nonexistent.md"),
			skipBase64: false,
			wantErr:    true,
			checkContent: func(t *testing.T, html string) {
				// HTML should be empty on error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := generateFromFile(tt.inputFile, tt.skipBase64)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkContent != nil {
				tt.checkContent(t, html)
			}
		})
	}
}

// TestGenerateFromFileWithImages tests image handling
func TestGenerateFromFileWithImages(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test image file
	imgFile := filepath.Join(tmpDir, "test.png")
	if err := os.WriteFile(imgFile, []byte("fake png data"), 0644); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Create markdown with local image reference
	mdFile := filepath.Join(tmpDir, "test.md")
	content := `# Test Slide

![Test Image](test.png)`

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with base64 encoding (skipBase64 = false)
	html, err := generateFromFile(mdFile, false)
	if err != nil {
		t.Fatalf("generateFromFile() error = %v", err)
	}

	// In base64 mode, should contain data: URI
	if !strings.Contains(html, "Test Image") {
		t.Error("Expected HTML to contain 'Test Image'")
	}

	// Test without base64 encoding (skipBase64 = true)
	html, err = generateFromFile(mdFile, true)
	if err != nil {
		t.Fatalf("generateFromFile() error = %v", err)
	}

	if !strings.Contains(html, "Test Image") {
		t.Error("Expected HTML to contain 'Test Image'")
	}
}

// TestGenerateFromFileWithMetadata tests presentation-level metadata
func TestGenerateFromFileWithMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")

	content := `<!-- presentation
title: Test Presentation
theme: light
-->

# First Slide

Content here.`

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	html, err := generateFromFile(mdFile, false)
	if err != nil {
		t.Fatalf("generateFromFile() error = %v", err)
	}

	if !strings.Contains(html, "First Slide") {
		t.Error("Expected HTML to contain 'First Slide'")
	}
}

// TestRewriteExternalImageURLs tests URL rewriting for external images
func TestRewriteExternalImageURLs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "http URL",
			input:    `<img src="http://example.com/image.png" alt="Test">`,
			expected: `<img src="/_proxy/?url=http%3A%2F%2Fexample.com%2Fimage.png" alt="Test">`,
		},
		{
			name:     "https URL",
			input:    `<img src="https://example.com/image.png" alt="Test">`,
			expected: `<img src="/_proxy/?url=https%3A%2F%2Fexample.com%2Fimage.png" alt="Test">`,
		},
		{
			name:     "local relative path unchanged",
			input:    `<img src="local/image.png" alt="Test">`,
			expected: `<img src="local/image.png" alt="Test">`,
		},
		{
			name:     "data URI unchanged",
			input:    `<img src="data:image/png;base64,iVBORw0KGg" alt="Test">`,
			expected: `<img src="data:image/png;base64,iVBORw0KGg" alt="Test">`,
		},
		{
			name:     "multiple external images",
			input:    `<img src="https://example.com/1.png"><img src="http://test.com/2.jpg">`,
			expected: `<img src="/_proxy/?url=https%3A%2F%2Fexample.com%2F1.png"><img src="/_proxy/?url=http%3A%2F%2Ftest.com%2F2.jpg">`,
		},
		{
			name:     "mixed local and external",
			input:    `<img src="local.png"><img src="https://example.com/remote.png">`,
			expected: `<img src="local.png"><img src="/_proxy/?url=https%3A%2F%2Fexample.com%2Fremote.png">`,
		},
		{
			name:     "image with attributes before src",
			input:    `<img alt="Test" class="big" src="https://example.com/img.png">`,
			expected: `<img alt="Test" class="big" src="/_proxy/?url=https%3A%2F%2Fexample.com%2Fimg.png">`,
		},
		{
			name:     "image with attributes after src",
			input:    `<img src="https://example.com/img.png" alt="Test" class="big">`,
			expected: `<img src="/_proxy/?url=https%3A%2F%2Fexample.com%2Fimg.png" alt="Test" class="big">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rewriteExternalImageURLs(tt.input)
			if result != tt.expected {
				t.Errorf("rewriteExternalImageURLs() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestRewriteExternalImageURLsComplexHTML tests URL rewriting with real HTML
func TestRewriteExternalImageURLsComplexHTML(t *testing.T) {
	input := `<!DOCTYPE html>
<html>
<body>
	<div>
		<img src="https://example.com/logo.png" alt="Logo" class="header">
		<p>Some text</p>
		<img src="local/image.jpg" alt="Local">
		<img src="http://cdn.example.com/banner.png" alt="Banner">
	</div>
</body>
</html>`

	result := rewriteExternalImageURLs(input)

	// Check that external URLs were rewritten
	if !strings.Contains(result, `/_proxy/?url=https%3A%2F%2Fexample.com%2Flogo.png`) {
		t.Error("Expected https URL to be proxied")
	}
	if !strings.Contains(result, `/_proxy/?url=http%3A%2F%2Fcdn.example.com%2Fbanner.png`) {
		t.Error("Expected http URL to be proxied")
	}

	// Check that local URLs were NOT rewritten
	if !strings.Contains(result, `src="local/image.jpg"`) {
		t.Error("Expected local URL to remain unchanged")
	}

	// Check that structure is preserved
	if !strings.Contains(result, `alt="Logo"`) {
		t.Error("Expected attributes to be preserved")
	}
	if !strings.Contains(result, `class="header"`) {
		t.Error("Expected class attribute to be preserved")
	}
}

// TestRun tests the run function with file output
func TestRun(t *testing.T) {
	// Save original flags
	origOutputFile := *outputFile
	origTheme := *theme
	origAspectRatio := *aspectRatio
	origTitle := *title

	defer func() {
		*outputFile = origOutputFile
		*theme = origTheme
		*aspectRatio = origAspectRatio
		*title = origTitle
	}()

	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")
	outFile := filepath.Join(tmpDir, "output.html")

	content := `# Test Slide

Test content.`

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Set flags
	*outputFile = outFile
	*theme = "dark"
	*aspectRatio = "1.6"
	*title = ""

	// Run
	err := run(mdFile)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	// Check that output file was created
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Check content
	htmlContent, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	html := string(htmlContent)
	if !strings.Contains(html, "Test Slide") {
		t.Error("Output file does not contain expected content")
	}
	if !strings.Contains(html, "<html") {
		t.Error("Output file is not valid HTML")
	}
}

// TestRunWithInvalidInput tests error handling
func TestRunWithInvalidInput(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.md")

	err := run(nonExistentFile)
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

// TestRunWithEmptyContent tests handling of empty markdown
func TestRunWithEmptyContent(t *testing.T) {
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "empty.md")

	if err := os.WriteFile(mdFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err := run(mdFile)
	if err == nil {
		t.Error("Expected error for empty content")
	}
}
