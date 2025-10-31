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
    /* Mermaid diagram styling */
    .mermaid {
      display: flex;
      justify-content: center;
      align-items: center;
    }
  </style>
  %s
  <script>
%s
  </script>
  <script>
%s
  </script>
  <script>
    // Initialize Mermaid with theme
    if (typeof mermaid !== 'undefined') {
      mermaid.initialize({
        startOnLoad: true,
        theme: '%s' === 'dark' ? 'dark' : 'default',
        securityLevel: 'loose'
      });
    }
  </script>
</head>
<body class="%s">
%s
</body>
</html>`

// generateHTML generates the complete HTML document
func generateHTML(title, bigCSS, themeCSS, customCSS, bigJS, mermaidJS, theme, themeClass, slides string) string {
	return fmt.Sprintf(
		htmlTemplate,
		title,      // %s - title
		bigCSS,     // %s - big.css
		themeCSS,   // %s - theme CSS
		customCSS,  // %s - aspect ratio script
		bigJS,      // %s - big.js
		mermaidJS,  // %s - mermaid.js
		theme,      // %s - theme for mermaid initialization
		themeClass, // %s - body class
		slides,     // %s - slides HTML
	)
}

// aspectRatioScript generates the aspect ratio configuration script
func aspectRatioScript(ratio string) string {
	if ratio == "" || ratio == "1.6" {
		return "" // Default is 1.6, no need to override
	}

	if ratio == "false" || ratio == "none" {
		return "<script>BIG_ASPECT_RATIO = false;</script>"
	}

	return fmt.Sprintf("<script>BIG_ASPECT_RATIO = %s;</script>", ratio)
}
