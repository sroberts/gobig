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
  </style>%s
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
func generateHTML(title, bigCSS, themeCSS, customCSS, aspectRatioScript, bigJS, themeClass, slides string) string {
	// Format custom CSS with style tags if present
	customCSSBlock := ""
	if customCSS != "" {
		customCSSBlock = fmt.Sprintf("\n  <style>\n%s\n  </style>", customCSS)
	}

	return fmt.Sprintf(
		htmlTemplate,
		title,
		bigCSS,
		themeCSS,
		customCSSBlock,
		aspectRatioScript,
		bigJS,
		themeClass,
		slides,
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
