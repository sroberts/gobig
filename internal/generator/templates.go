package generator

import (
	"fmt"
)

// htmlTemplate is the base HTML structure for the presentation
const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0" />
  <title>%s</title>
  <style>
%s
  </style>
  <style>
%s
  </style>
  <style>
    /* Mermaid SVG diagram styling - responsive sizing for presentations */
    .mermaid {
      display: flex !important;
      justify-content: center !important;
      align-items: center !important;
      max-width: 100%% !important;
      max-height: 100%% !important;
    }
    .mermaid svg {
      max-width: 90%% !important;
      max-height: 90%% !important;
      width: 100%% !important;
      height: 100%% !important;
    }
  </style>
  <style>
    /* DeckSet [fit] header support - scales text to fill slide */
    .fit {
      display: block;
      font-size: 10vw;
      line-height: 1.2;
      font-weight: bold;
    }
    
    /* Background color support for slides */
    div[data-background-color] {
      background-color: var(--bg-color) !important;
    }
    
    /* Autoscale support */
    div[data-autoscale="true"] {
      font-size: 0.9em;
    }
  </style>
  <style>
    /* Force dark theme colors for Mermaid diagrams with proper contrast */
    .dark .mermaid rect,
    .light .mermaid rect:not([fill]),
    .white .mermaid rect:not([fill]) {
      fill: #1f2020 !important;
      stroke: #ccc !important;
    }
    .dark .mermaid polygon {
      fill: #1f2020 !important;
      stroke: #ccc !important;
    }
    /* Ensure text is light colored on dark backgrounds */
    .dark .mermaid text,
    .dark .mermaid tspan,
    .dark .mermaid .nodeLabel,
    .dark .mermaid .edgeLabel {
      fill: #f0f0f0 !important;
      color: #f0f0f0 !important;
    }
    /* Allow diagrams to scale properly without text cutoff */
    .dark .mermaid foreignObject {
      overflow: visible !important;
    }
  </style>
  %s
  <script>
%s
  </script>
</head>
<body class="%s">
%s
</body>
</html>`

// generateHTML generates the complete HTML document
func generateHTML(title, bigCSS, themeCSS, customCSS, bigJS, theme, slides string) string {
	return fmt.Sprintf(
		htmlTemplate,
		title,     // %s - title
		bigCSS,    // %s - big.css
		themeCSS,  // %s - theme CSS
		customCSS, // %s - aspect ratio script
		bigJS,     // %s - big.js
		theme,     // %s - body class (theme)
		slides,    // %s - slides HTML
	)
}

// aspectRatioScript generates the aspect ratio configuration script
func aspectRatioScript(ratio string) string {
	script := ""
	
	if ratio != "" && ratio != "1.6" {
		if ratio == "false" || ratio == "none" {
			script += "<script>BIG_ASPECT_RATIO = false;</script>"
		} else {
			script += fmt.Sprintf("<script>BIG_ASPECT_RATIO = %s;</script>", ratio)
		}
	}
	
	// Add DeckSet features script
	script += `
<script>
// Handle DeckSet background colors and other data attributes
document.addEventListener('DOMContentLoaded', function() {
  // Process background colors
  document.querySelectorAll('div[data-background-color]').forEach(function(slide) {
    var color = slide.getAttribute('data-background-color');
    if (color) {
      slide.style.backgroundColor = color;
    }
  });
});
</script>`
	
	return script
}
