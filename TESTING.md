# Testing Guide

This document describes the testing infrastructure for the GoBig project.

## Overview

GoBig uses Go's built-in testing framework with comprehensive test coverage across all packages:

- **Unit tests**: Test individual functions and components
- **Integration tests**: Test the complete pipeline from parsing to generation
- **Edge case tests**: Test error handling and boundary conditions

## Current Test Coverage

```
Package                  Coverage
----------------------   ---------
internal/parser          97.3%
internal/assets          86.7%
internal/generator       63.2%
cmd/gobig                18.6%
----------------------   ---------
Total                    52.8%
```

## Running Tests

### Quick Test Run

Run all tests:
```bash
make test
```

Or using Go directly:
```bash
go test ./...
```

### Test with Race Detection

Detect race conditions:
```bash
make test-race
```

Or:
```bash
go test -race ./...
```

### Test with Coverage

Generate coverage report:
```bash
make test-coverage
```

View coverage summary:
```bash
go tool cover -func=coverage.txt | grep total
```

### HTML Coverage Report

Generate and view HTML coverage report:
```bash
make test-coverage-html
open coverage.html
```

### Verbose Testing

Run tests with verbose output per package:
```bash
make test-verbose
```

### Check Coverage Threshold

Verify coverage meets minimum threshold (50%):
```bash
make test-coverage-check
```

### Run Benchmarks

```bash
make test-bench
```

## Test Organization

### Unit Tests

Located alongside source files with `_test.go` suffix:

- `internal/parser/parser_test.go` - Parser tests
- `internal/generator/generator_test.go` - Generator tests
- `internal/assets/assets_test.go` - Asset loading tests
- `cmd/gobig/main_test.go` - CLI and helper function tests

### Integration Tests

Located in `test/integration/`:

- `integration_test.go` - End-to-end pipeline tests

## Test Categories

### Parser Tests

Test markdown parsing, metadata extraction, and slide separation:

```go
func TestParseStringBasic(t *testing.T) {
    p := parser.NewParser()
    err := p.ParseString(`# Slide 1\n\n---\n\n# Slide 2`)
    // ...
}
```

**What's tested:**
- Basic slide parsing
- Slide metadata (layout, class, timing)
- Presentation metadata (title, auto-advance)
- Speaker notes extraction
- Malformed YAML handling
- Edge cases (empty content, whitespace, multiple separators)

### Generator Tests

Test HTML generation from parsed slides:

```go
func TestGenerateBasic(t *testing.T) {
    gen := generator.NewGenerator(opts)
    html, err := gen.Generate(slides)
    // ...
}
```

**What's tested:**
- HTML structure generation
- Theme application (dark, light, white)
- Aspect ratio handling
- Slide layouts (50-50, 75-25, grid-3x2, etc.)
- Markdown rendering (code, lists, links, emphasis)
- Auto-advance timing
- Special character escaping

### Main Package Tests

Test CLI functionality and helper functions:

```go
func TestGenerateFromFile(t *testing.T) {
    html, err := generateFromFile(mdFile, false)
    // ...
}
```

**What's tested:**
- File generation workflow
- Image handling (base64 encoding, local paths)
- External image URL rewriting for proxy
- Error handling for invalid inputs

### Integration Tests

Test the complete pipeline:

```go
func TestParserGeneratorIntegration(t *testing.T) {
    p := parser.NewParser()
    p.ParseString(markdown)
    gen := generator.NewGenerator(opts)
    html, err := gen.Generate(p.GetSlides())
    // ...
}
```

**What's tested:**
- Full pipeline from markdown to HTML
- Multiple themes with same content
- Various layouts
- Complex markdown features
- Metadata propagation

## Writing New Tests

### Test Structure

Follow Go testing conventions:

```go
func TestFeatureName(t *testing.T) {
    // Arrange - setup test data
    input := "test data"

    // Act - execute the code under test
    result, err := FunctionUnderTest(input)

    // Assert - verify the results
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Table-Driven Tests

For testing multiple scenarios:

```go
func TestMultipleScenarios(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"scenario 1", "input1", "output1"},
        {"scenario 2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Process(tt.input)
            if result != tt.expected {
                t.Errorf("Expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

### Using Temporary Files

Use `t.TempDir()` for file-based tests:

```go
func TestWithFile(t *testing.T) {
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.md")
    os.WriteFile(testFile, []byte("content"), 0644)

    // Test with testFile
    // tmpDir is automatically cleaned up
}
```

## Continuous Integration

Tests run automatically on:
- Every push to main branch
- Every pull request

CI performs:
1. Test execution with race detection
2. Coverage report generation
3. Coverage threshold check (must be ≥50%)
4. Coverage breakdown by package

## Coverage Goals

Current thresholds:
- **Minimum overall coverage**: 50%
- **Individual package targets**:
  - `internal/parser`: ≥90% (currently 97.3%)
  - `internal/assets`: ≥80% (currently 86.7%)
  - `internal/generator`: ≥60% (currently 63.2%)

## Common Test Commands

```bash
# Run all tests
make test

# Run with race detection
make test-race

# Generate coverage report
make test-coverage

# Check coverage threshold
make test-coverage-check

# Generate HTML report
make test-coverage-html

# Run verbose tests
make test-verbose

# Run benchmarks
make test-bench

# Run specific package tests
go test ./internal/parser
go test ./internal/generator
go test ./cmd/gobig

# Run specific test
go test -run TestParseStringBasic ./internal/parser
```

## Debugging Failed Tests

### Enable Verbose Output

```bash
go test -v ./internal/parser
```

### Run Specific Test

```bash
go test -run TestGenerateBasic ./internal/generator
```

### Show Test Coverage for Specific Package

```bash
go test -cover ./internal/parser
```

### Debug with Print Statements

```go
func TestDebug(t *testing.T) {
    result := Calculate(5)
    t.Logf("Result: %v", result)  // Only shown with -v flag
    // ...
}
```

## Best Practices

1. **Write tests first**: Consider TDD for new features
2. **Test edge cases**: Empty inputs, nil values, malformed data
3. **Use table-driven tests**: For testing multiple scenarios
4. **Keep tests isolated**: Don't rely on test execution order
5. **Use t.TempDir()**: For file operations
6. **Check errors**: Always verify error conditions
7. **Meaningful test names**: Use descriptive names that explain what's being tested
8. **Don't test external dependencies**: Mock or skip external services

## Troubleshooting

### Tests Pass Locally But Fail in CI

- Check for race conditions: `go test -race ./...`
- Verify file paths are relative, not absolute
- Ensure tests don't depend on local environment

### Coverage Lower Than Expected

- Check which functions are untested: `go tool cover -func=coverage.txt`
- Add tests for uncovered code paths
- Consider if some code is untestable (e.g., server initialization)

### Flaky Tests

- Use `-race` flag to detect race conditions
- Avoid sleep-based timing in tests
- Use proper synchronization primitives

## Additional Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Go Code Coverage](https://go.dev/blog/cover)
