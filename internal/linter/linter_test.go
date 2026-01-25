package linter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLinter(t *testing.T) {
	linter := NewLinter("/tmp")
	if linter == nil {
		t.Fatal("Expected non-nil linter")
	}
	if linter.baseDir != "/tmp" {
		t.Errorf("Expected baseDir to be /tmp, got %s", linter.baseDir)
	}
}

func TestLintStringValid(t *testing.T) {
	linter := NewLinter(".")
	content := `# Hello World

This is a test slide.

---

# Second Slide

More content here.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.HasError {
		t.Errorf("Expected no errors, got: %v", result.Issues)
	}

	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues, got: %v", result.Issues)
	}
}

func TestLintStringNoSlides(t *testing.T) {
	linter := NewLinter(".")
	content := ``

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.HasError {
		t.Error("Expected error for file with no slides")
	}

	if len(result.Issues) == 0 {
		t.Error("Expected issues for file with no slides")
	}
}

func TestLintStringInvalidLayout(t *testing.T) {
	linter := NewLinter(".")
	content := `<!-- slide
layout: invalid-layout
-->

# Test Slide

Content here.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have a warning about unknown layout
	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Severity == "warning" && strings.Contains(issue.Message, "Unknown layout") {
			foundWarning = true
			break
		}
	}

	if !foundWarning {
		t.Errorf("Expected warning about invalid layout, got issues: %v", result.Issues)
	}
}

func TestLintStringFirstSlideInvalidLayout(t *testing.T) {
	linter := NewLinter(".")
	content := `<!-- slide
layout: invalid-layout
-->

# First Slide

Content here.

---

# Second Slide

More content.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have a warning about unknown layout in the first slide
	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Severity == "warning" && strings.Contains(issue.Message, "Unknown layout") {
			foundWarning = true
			break
		}
	}

	if !foundWarning {
		t.Errorf("Expected warning about invalid layout in first slide, got issues: %v", result.Issues)
	}
}

func TestLintStringFirstSlideNegativeTimeToNext(t *testing.T) {
	linter := NewLinter(".")
	content := `<!-- slide
time-to-next: -5
-->

# First Slide

Content here.

---

# Second Slide

More content.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have an error about negative time-to-next in the first slide
	foundError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && strings.Contains(issue.Message, "time-to-next cannot be negative") {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Errorf("Expected error about negative time-to-next in first slide, got issues: %v", result.Issues)
	}
}

func TestLintStringValidLayout(t *testing.T) {
	linter := NewLinter(".")
	content := `<!-- slide
layout: 50-50
-->

# Test Slide

Content here.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.HasError {
		t.Errorf("Expected no errors for valid layout, got: %v", result.Issues)
	}
}

func TestLintStringCSSCustomLayout(t *testing.T) {
	linter := NewLinter(".")
	content := `<!-- slide
layout: grid-template-columns: 1fr 2fr;
-->

# Test Slide

Content here.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should not have a warning about unknown layout for CSS-style layouts
	for _, issue := range result.Issues {
		if strings.Contains(issue.Message, "Unknown layout") {
			t.Errorf("Unexpected warning about CSS-style layout: %v", issue)
		}
	}
}

func TestLintStringInvalidSlideMetadata(t *testing.T) {
	linter := NewLinter(".")
	content := `<!-- slide
layout: "Test"
invalid-yaml: [unclosed array
-->

# Test Slide

Content here.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.HasError {
		t.Error("Expected error for invalid YAML in slide metadata")
	}

	foundError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && strings.Contains(issue.Message, "Invalid slide metadata YAML") {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Errorf("Expected error about invalid slide metadata YAML, got issues: %v", result.Issues)
	}
}

func TestLintStringNegativeTimeToNext(t *testing.T) {
	linter := NewLinter(".")
	content := `<!-- slide
time-to-next: -5
-->

# Test Slide

Content here.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.HasError {
		t.Error("Expected error for negative time-to-next")
	}

	foundError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && strings.Contains(issue.Message, "time-to-next cannot be negative") {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Errorf("Expected error about negative time-to-next, got issues: %v", result.Issues)
	}
}

func TestLintStringInvalidPresentationMetadata(t *testing.T) {
	linter := NewLinter(".")
	content := `<!-- presentation
title: "Test"
invalid-yaml: [unclosed array
-->

# Test Slide

Content here.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.HasError {
		t.Error("Expected error for invalid YAML in presentation metadata")
	}
}

func TestLintStringMissingImage(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	linter := NewLinter(tmpDir)
	content := `# Test Slide

![Test Image](nonexistent.png)
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.HasError {
		t.Error("Expected error for missing image file")
	}

	foundError := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" && strings.Contains(issue.Message, "Image file not found") {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Errorf("Expected error about missing image, got issues: %v", result.Issues)
	}
}

func TestLintStringExistingImage(t *testing.T) {
	// Create a temporary directory and image file
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.png")
	err := os.WriteFile(imagePath, []byte("fake image"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	linter := NewLinter(tmpDir)
	content := `# Test Slide

![Test Image](test.png)
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should not have errors about missing image
	for _, issue := range result.Issues {
		if strings.Contains(issue.Message, "Image file not found") {
			t.Errorf("Unexpected error about missing image: %v", issue)
		}
	}
}

func TestLintStringExternalImage(t *testing.T) {
	linter := NewLinter(".")
	content := `# Test Slide

![External Image](https://example.com/image.png)
![HTTP Image](http://example.com/image.png)
![Data URI](data:image/png;base64,abc123)
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should not have errors about missing external images
	for _, issue := range result.Issues {
		if strings.Contains(issue.Message, "Image file not found") {
			t.Errorf("Unexpected error about missing external image: %v", issue)
		}
	}
}

func TestLintStringMalformedHorizontalRule(t *testing.T) {
	linter := NewLinter(".")
	content := `# First Slide

Content here.

-- (only 2 dashes)

# Second Slide

More content.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Severity == "warning" && strings.Contains(issue.Message, "malformed horizontal rule") {
			foundWarning = true
			break
		}
	}

	if !foundWarning {
		t.Errorf("Expected warning about malformed horizontal rule, got issues: %v", result.Issues)
	}
}

func TestLintStringEmptySlide(t *testing.T) {
	// Note: The parser filters out completely empty slides, so they won't
	// appear in the parsed results. This test verifies that the linter
	// handles this gracefully and doesn't crash.
	linter := NewLinter(".")
	content := `# First Slide

Content here.

---

<!-- This slide has only a comment -->

---

# Third Slide

More content.
`

	result, err := linter.LintString(content, "test.md")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should not have any errors (empty slides are filtered by parser)
	if result.HasError {
		t.Errorf("Unexpected errors: %v", result.Issues)
	}
}

func TestFormatIssues(t *testing.T) {
	result := &Result{
		FilePath: "test.md",
		Issues: []Issue{
			{Line: 5, Severity: "error", Message: "Test error"},
			{Line: 10, Severity: "warning", Message: "Test warning"},
			{Line: 0, Severity: "info", Message: "General info"},
		},
	}

	formatted := result.FormatIssues()

	if !strings.Contains(formatted, "test.md:5: [error] Test error") {
		t.Error("Formatted output should contain line 5 error")
	}

	if !strings.Contains(formatted, "test.md:10: [warning] Test warning") {
		t.Error("Formatted output should contain line 10 warning")
	}

	if !strings.Contains(formatted, "test.md: [info] General info") {
		t.Error("Formatted output should contain general info")
	}
}

func TestGetSummary(t *testing.T) {
	result := &Result{
		Issues: []Issue{
			{Severity: "error", Message: "Error 1"},
			{Severity: "error", Message: "Error 2"},
			{Severity: "warning", Message: "Warning 1"},
			{Severity: "info", Message: "Info 1"},
		},
	}

	summary := result.GetSummary()

	if !strings.Contains(summary, "4 issue(s)") {
		t.Error("Summary should contain total issue count")
	}

	if !strings.Contains(summary, "2 error(s)") {
		t.Error("Summary should contain error count")
	}

	if !strings.Contains(summary, "1 warning(s)") {
		t.Error("Summary should contain warning count")
	}
}

func TestGetSummaryNoIssues(t *testing.T) {
	result := &Result{
		Issues: []Issue{},
	}

	summary := result.GetSummary()

	if !strings.Contains(summary, "No issues found") {
		t.Errorf("Expected 'No issues found', got: %s", summary)
	}
}
