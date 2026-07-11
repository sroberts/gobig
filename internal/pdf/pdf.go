// Package pdf renders a generated gobig presentation to a PDF file by driving
// a headless Chrome/Chromium instance. Each slide is captured at the
// presentation's aspect ratio and placed on its own page, so the PDF is
// suitable for upload to SpeakerDeck, Slideshare, and similar services.
package pdf

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-pdf/fpdf"
)

// Options configures a PDF render.
type Options struct {
	HTML       string // self-contained presentation HTML
	SlideCount int    // number of slides in the deck
	Width      int    // slide width in CSS pixels
	Height     int    // slide height in CSS pixels
	OutputPath string // destination PDF path
}

// perSlideTimeout bounds how long we wait for a single slide to render.
const perSlideTimeout = 20 * time.Second

// Render captures each slide of the presentation with headless Chrome and
// writes them to a single PDF, one slide per page.
func Render(opts Options) error {
	if opts.SlideCount < 1 {
		return fmt.Errorf("no slides to render")
	}

	// Write the self-contained HTML to a temp file so Chrome can load it via
	// a file:// URL (avoids data-URL length limits with base64 images).
	tmp, err := os.CreateTemp("", "gobig-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(opts.HTML); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	fileURL := "file://" + tmp.Name()

	// Capture a PNG of every slide.
	shots, err := captureSlides(fileURL, opts)
	if err != nil {
		return err
	}

	return writePDF(shots, opts.OutputPath)
}

// captureSlides launches headless Chrome, then navigates to each slide by hash
// and screenshots the viewport. Screenshots use the screen media type, so
// big.js stays in its single-slide "talk" mode (its print listener, which
// switches to a 2-up handout layout, is never triggered).
func captureSlides(fileURL string, opts Options) ([][]byte, error) {
	// Render at 2x for crisp output in the PDF.
	const scale = 2

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.WindowSize(opts.Width, opts.Height),
		)...,
	)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Fail fast with a helpful message if Chrome can't be started.
	startCtx, startCancel := context.WithTimeout(ctx, 30*time.Second)
	defer startCancel()
	if err := chromedp.Run(startCtx,
		chromedp.EmulateViewport(int64(opts.Width), int64(opts.Height), chromedp.EmulateScale(scale)),
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("div.slide-container", chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("failed to start headless Chrome (is Google Chrome or Chromium installed?): %w", err)
	}

	shots := make([][]byte, 0, opts.SlideCount)
	for i := 0; i < opts.SlideCount; i++ {
		var buf []byte
		slideCtx, slideCancel := context.WithTimeout(ctx, perSlideTimeout)
		err := chromedp.Run(slideCtx,
			// Select the slide by hash; big.js listens for hashchange and
			// scales the target slide to fill the viewport.
			chromedp.Evaluate(fmt.Sprintf("window.location.hash = '#%d'", i), nil),
			chromedp.Sleep(250*time.Millisecond),
			chromedp.CaptureScreenshot(&buf),
		)
		slideCancel()
		if err != nil {
			return nil, fmt.Errorf("failed to capture slide %d: %w", i+1, err)
		}
		shots = append(shots, buf)
	}

	return shots, nil
}

// writePDF places each PNG screenshot on its own page, sized to the image so
// the slide fills the page edge to edge.
func writePDF(shots [][]byte, outputPath string) error {
	if len(shots) == 0 {
		return fmt.Errorf("no slides captured")
	}

	// Derive the page size from the first screenshot's pixel dimensions,
	// treating each pixel as a PDF point (aspect ratio is what matters).
	cfg, err := png.DecodeConfig(bytes.NewReader(shots[0]))
	if err != nil {
		return fmt.Errorf("failed to read screenshot: %w", err)
	}
	wPt, hPt := float64(cfg.Width), float64(cfg.Height)

	// Always use "P" orientation and pass the exact pixel dimensions: gofpdf
	// swaps width/height for "L", which would flip a landscape slide to a
	// portrait page.
	size := fpdf.SizeType{Wd: wPt, Ht: hPt}
	pdf := fpdf.NewCustom(&fpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "pt",
		Size:           size,
	})
	pdf.SetAutoPageBreak(false, 0)

	for i, shot := range shots {
		pdf.AddPageFormat("P", size)
		name := fmt.Sprintf("slide-%d", i)
		pdf.RegisterImageOptionsReader(name, fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(shot))
		pdf.ImageOptions(name, 0, 0, wPt, hPt, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	}

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	return nil
}
