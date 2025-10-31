package generator

import (
	"fmt"
)

// generateHTML generates the complete HTML document
func generateHTML(title, bigCSS, themeCSS, codeBlocksCSS, customCSS, bigJS, shikiJS, themeClass, slides string) string {
	// Build shiki script tag if shikiJS is provided
	shikiScriptTag := ""
	if shikiJS != "" {
		shikiScriptTag = fmt.Sprintf(`  <script type="module">
%s
  </script>`, shikiJS)
	}

	htmlTemplate := `<!DOCTYPE html>
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
%s
  </style>
  %s
  <script>
%s
  </script>
%s
</head>
<body class="%s">
%s
</body>
</html>`

	return fmt.Sprintf(
		htmlTemplate,
		title,
		bigCSS,
		themeCSS,
		codeBlocksCSS,
		customCSS,
		bigJS,
		shikiScriptTag,
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
