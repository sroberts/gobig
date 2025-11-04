package parser

// SlideMetadata represents the YAML frontmatter for a slide
type SlideMetadata struct {
	Layout     string `yaml:"layout"`       // e.g., "50-50", "75-25-rows", "grid-3x2"
	Class      string `yaml:"class"`        // Custom CSS classes
	BodyStyle  string `yaml:"body-style"`   // Custom body styling for this slide
	BodyClass  string `yaml:"body-class"`   // Custom body class for this slide
	TimeToNext int    `yaml:"time-to-next"` // Auto-advance time in seconds
	
	// DeckSet per-slide directives
	Autoscale       *bool  `yaml:"autoscale"`        // Scale text to fit slide
	BuildLists      *bool  `yaml:"build-lists"`      // Reveal list items one by one
	SlideNumbers    *bool  `yaml:"slidenumbers"`     // Show/hide slide numbers
	Footer          string `yaml:"footer"`           // Custom footer for this slide
	HideFooter      bool   `yaml:"hide-footer"`      // Hide footer on this slide
	SlideTransition string `yaml:"slide-transition"` // Slide transition effect
	BackgroundColor string `yaml:"background-color"` // Background color for slide
}

// PresentationMetadata represents presentation-level metadata
type PresentationMetadata struct {
	Title      string `yaml:"title"`        // Presentation title
	TimeToNext int    `yaml:"time-to-next"` // Default auto-advance time for all slides
	
	// DeckSet global configuration
	Autoscale       bool   `yaml:"autoscale"`        // Scale text to fit slides
	BuildLists      bool   `yaml:"build-lists"`      // Reveal list items one by one
	SlideNumbers    bool   `yaml:"slidenumbers"`     // Add slide numbers
	Footer          string `yaml:"footer"`           // Persistent footer
	BackgroundImage string `yaml:"background-image"` // Background image for all slides
	FitHeader       bool   `yaml:"fit-header"`       // Fit all headers to full slide
	SlideDividers   string `yaml:"slide-dividers"`   // Auto-create slides at header levels
	SlideTransition string `yaml:"slide-transition"` // Default transition effect
	Theme           string `yaml:"theme"`            // Theme name and color variant
}

// Slide represents a single presentation slide
type Slide struct {
	Metadata SlideMetadata // Parsed frontmatter
	Content  string        // Raw markdown content (without frontmatter)
	Notes    string        // Speaker notes extracted from HTML comments
}

// SlideType represents the detected type of slide
type SlideType int

const (
	SlideTypeContent SlideType = iota
	SlideTypeTitle
	SlideTypeSection
	SlideTypeTable
)

// DetectType attempts to detect the slide type based on content
func (s *Slide) DetectType(htmlContent string) SlideType {
	// This will be implemented after we have HTML content
	// For now, return content type
	return SlideTypeContent
}
