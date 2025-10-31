package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gobig/internal/assets"
	"gobig/internal/generator"
	"gobig/internal/parser"
)

const version = "1.0.0"

var (
	outputFile        = flag.String("o", "", "Output HTML file (default: stdout)")
	theme             = flag.String("theme", "dark", "Theme: dark, light, or white")
	aspectRatio       = flag.String("aspect-ratio", "1.6", "Aspect ratio (e.g., 1.6, 2, false)")
	title             = flag.String("title", "", "Presentation title (default: from first slide)")
	syntaxHighlighting = flag.Bool("syntax-highlighting", false, "Enable language-specific syntax highlighting with Shiki")
	showVersion       = flag.Bool("version", false, "Show version information")
	showHelp          = flag.Bool("help", false, "Show help message")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("gobig version %s\n", version)
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		usage()
		os.Exit(0)
	}

	// Get input file
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: exactly one input file required")
		usage()
		os.Exit(1)
	}

	inputFile := args[0]

	// Validate theme
	if !assets.ValidateTheme(*theme) {
		fmt.Fprintf(os.Stderr, "Error: invalid theme '%s'. Valid themes: dark, light, white\n", *theme)
		os.Exit(1)
	}

	// Run the conversion
	if err := run(inputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(inputFile string) error {
	// Parse markdown file
	p := parser.NewParser()
	if err := p.ParseFile(inputFile); err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	slides := p.GetSlides()
	presentationMetadata := p.GetPresentationMetadata()

	// Get base path for resolving relative image paths
	basePath, err := filepath.Abs(filepath.Dir(inputFile))
	if err != nil {
		basePath = filepath.Dir(inputFile)
	}

	// Generate HTML
	opts := generator.Options{
		Theme:                *theme,
		Title:                *title,
		AspectRatio:          *aspectRatio,
		BasePath:             basePath,
		PresentationMetadata: presentationMetadata,
		SyntaxHighlighting:   *syntaxHighlighting,
	}

	gen := generator.NewGenerator(opts)
	html, err := gen.Generate(slides)
	if err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Output HTML
	if *outputFile != "" {
		// Write to file
		if err := os.WriteFile(*outputFile, []byte(html), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Presentation generated: %s\n", *outputFile)
	} else {
		// Write to stdout
		fmt.Print(html)
	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `gobig - Generate big.js presentations from Markdown

Usage:
  gobig [options] <input.md>

Options:
  -o <file>              Output HTML file (default: stdout)
  -theme <name>          Theme: dark, light, or white (default: dark)
  -aspect-ratio <ratio>  Aspect ratio: number or "false" to disable (default: 1.6)
  -title <title>         Presentation title (default: from first slide)
  -syntax-highlighting   Enable language-specific syntax highlighting with Shiki (default: false)
  -version               Show version information
  -help                  Show this help message

Examples:
  gobig -o index.html presentation.md
  gobig -theme light -o output.html slides.md
  gobig -aspect-ratio 2 -title "My Talk" -o slides.html talk.md
  gobig -syntax-highlighting -o slides.html code-heavy-talk.md

Markdown Syntax:
  Slides:      Separate with --- (horizontal rule)
  Notes:       Use HTML comments: <!-- speaker notes here -->
  Metadata:    Use YAML frontmatter in comments:
               <!-- slide
               layout: 50-50
               class: custom-class
               -->

Layouts:
  50-50         Two columns (50%% each)
  75-25         Two columns (75%%, 25%%)
  25-75         Two columns (25%%, 75%%)
  50-50-rows    Two rows (50%% each)
  75-25-rows    Two rows (75%%, 25%%)
  grid-3x2      3 columns, 2 rows
  Custom CSS    Use custom grid-template syntax

For more information: https://github.com/tmcw/big
`)
}
