package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gobig/internal/assets"
	"gobig/internal/generator"
	"gobig/internal/parser"
)

const version = "1.0.0"

var (
	outputFile  = flag.String("o", "", "Output HTML file (default: stdout)")
	theme       = flag.String("theme", "dark", "Theme: dark, light, or white")
	aspectRatio = flag.String("aspect-ratio", "1.6", "Aspect ratio (e.g., 1.6, 2, false)")
	title       = flag.String("title", "", "Presentation title (default: from first slide)")
	serve       = flag.Bool("serve", false, "Run as web server instead of generating file")
	port        = flag.Int("port", 8080, "Port for web server (default: 8080)")
	watch       = flag.Bool("watch", false, "Watch markdown file for changes and regenerate (only with -serve)")
	showVersion = flag.Bool("version", false, "Show version information")
	showHelp    = flag.Bool("help", false, "Show help message")
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

	// Validate flags
	if *watch && !*serve {
		fmt.Fprintln(os.Stderr, "Error: -watch can only be used with -serve")
		os.Exit(1)
	}

	if *serve && *outputFile != "" {
		fmt.Fprintln(os.Stderr, "Error: -serve and -o cannot be used together")
		os.Exit(1)
	}

	// Run in serve mode or conversion mode
	if *serve {
		if err := runServer(inputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := run(inputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

// generateFromFile parses a markdown file and generates HTML
func generateFromFile(inputFile string) (string, error) {
	// Parse markdown file
	p := parser.NewParser()
	if err := p.ParseFile(inputFile); err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
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
	}

	gen := generator.NewGenerator(opts)
	html, err := gen.Generate(slides)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML: %w", err)
	}

	return html, nil
}

func run(inputFile string) error {
	html, err := generateFromFile(inputFile)
	if err != nil {
		return err
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

func runServer(inputFile string) error {
	// Keep track of the generated HTML
	var (
		currentHTML   string
		lastMod       time.Time
		mu            sync.RWMutex
		generateError error
	)

	// Function to generate HTML from markdown
	generateHTML := func() error {
		html, err := generateFromFile(inputFile)
		if err != nil {
			mu.Lock()
			generateError = err
			mu.Unlock()
			return err
		}

		mu.Lock()
		currentHTML = html
		generateError = nil
		mu.Unlock()

		return nil
	}

	// Initial generation
	if err := generateHTML(); err != nil {
		return err
	}

	// Get initial file info
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	
	mu.Lock()
	lastMod = fileInfo.ModTime()
	mu.Unlock()

	log.Printf("Serving presentation from %s", inputFile)

	// Start file watcher if enabled
	if *watch {
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				fileInfo, err := os.Stat(inputFile)
				if err != nil {
					log.Printf("Error checking file: %v", err)
					continue
				}

				mu.RLock()
				currentLastMod := lastMod
				mu.RUnlock()

				if fileInfo.ModTime().After(currentLastMod) {
					mu.Lock()
					lastMod = fileInfo.ModTime()
					mu.Unlock()
					
					log.Printf("File changed, regenerating...")
					if err := generateHTML(); err != nil {
						log.Printf("Error regenerating: %v", err)
					} else {
						log.Printf("Presentation regenerated successfully")
					}
				}
			}
		}()
		log.Printf("Watching %s for changes...", inputFile)
	}

	// HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		html := currentHTML
		genErr := generateError
		mu.RUnlock()

		if genErr != nil {
			http.Error(w, fmt.Sprintf("Error generating presentation: %v", genErr), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	})

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server listening on http://localhost%s", addr)
	log.Printf("Press Ctrl+C to stop")
	return http.ListenAndServe(addr, nil)
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
  -serve                 Run as web server instead of generating file
  -port <port>           Port for web server (default: 8080)
  -watch                 Watch markdown file for changes (only with -serve)
  -version               Show version information
  -help                  Show this help message

Examples:
  # Generate HTML file
  gobig -o index.html presentation.md
  gobig -theme light -o output.html slides.md
  gobig -aspect-ratio 2 -title "My Talk" -o slides.html talk.md

  # Run as web server
  gobig -serve presentation.md
  gobig -serve -port 3000 presentation.md
  gobig -serve -watch -theme light presentation.md

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
