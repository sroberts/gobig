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
%s
  </style>
  %s
  <script>
%s
  </script>
</head>
<body class="%s">
%s
%s
</body>
</html>`

// fullscreenSnippet adds a full-screen toggle to the presentation: press "f"
// or click the on-screen button to enter/exit full screen (Esc also exits).
// It is self-contained so it survives updates to the vendored big.js, and
// binds only the "f" key, which big.js does not use.
const fullscreenSnippet = `<button id="big-fullscreen" type="button" aria-label="Toggle full screen" title="Full screen (f)">&#x26F6;</button>
<style>
#big-fullscreen {
  position: fixed;
  bottom: 1rem;
  right: 1rem;
  z-index: 1000;
  width: 2.5rem;
  height: 2.5rem;
  padding: 0;
  border: none;
  border-radius: 0.375rem;
  background: rgba(128, 128, 128, 0.25);
  color: inherit;
  font: inherit;
  font-size: 1.25rem;
  line-height: 2.5rem;
  text-align: center;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s ease;
}
body:hover #big-fullscreen { opacity: 0.4; }
#big-fullscreen:hover { opacity: 1; }
:fullscreen #big-fullscreen { display: none; }
@media print { #big-fullscreen { display: none; } }
</style>
<script>
(function () {
  var btn = document.getElementById("big-fullscreen");
  function fsElement() {
    return document.fullscreenElement || document.webkitFullscreenElement;
  }
  function toggle() {
    if (fsElement()) {
      (document.exitFullscreen || document.webkitExitFullscreen || function () {}).call(document);
    } else {
      var el = document.documentElement;
      (el.requestFullscreen || el.webkitRequestFullscreen || function () {}).call(el);
    }
  }
  if (btn) {
    btn.addEventListener("click", function (e) {
      // Stop the click from bubbling to big.js's document handler, which
      // would otherwise advance to the next slide.
      e.preventDefault();
      e.stopPropagation();
      toggle();
      btn.blur();
    });
  }
  document.addEventListener("keydown", function (e) {
    if (e.key !== "f" && e.key !== "F") return;
    if (e.ctrlKey || e.metaKey || e.altKey) return;
    var t = e.target;
    if (t && (t.tagName === "INPUT" || t.tagName === "TEXTAREA" || t.isContentEditable)) return;
    e.preventDefault();
    toggle();
  });
})();
</script>`

// codeCSS styles inline and fenced code blocks so they render in a familiar
// monospace font inside a distinguishable panel. Colors are neutral and
// semi-transparent so they adapt to every theme (dark, light, white) without
// per-theme rules, and it is injected by the generator so it survives updates
// to the vendored big.css.
const codeCSS = `pre, code, kbd, samp {
  font-family: ui-monospace, "SF Mono", SFMono-Regular, Menlo, Consolas, "Liberation Mono", "Courier New", monospace;
  font-weight: 500;
}

/* Inline code */
code {
  font-size: 0.9em;
  background: rgba(127, 127, 127, 0.22);
  padding: 0.08em 0.32em;
  border-radius: 0.2em;
}

/* Fenced / indented code blocks */
pre {
  max-width: 100%;
  margin: 0.4em 0;
  padding: 0.6em 0.8em;
  font-size: 0.6em;
  line-height: 1.45;
  text-align: left;
  white-space: pre;
  overflow-x: auto;
  background: rgba(127, 127, 127, 0.16);
  border: 1px solid rgba(127, 127, 127, 0.28);
  border-radius: 0.35em;
}

/* Code inside a block should not get the inline pill treatment */
pre code {
  font-size: 1em;
  background: none;
  padding: 0;
  border-radius: 0;
  white-space: inherit;
}`

// generateHTML generates the complete HTML document
func generateHTML(title, bigCSS, themeCSS, customCSS, bigJS, theme, slides string) string {
	return fmt.Sprintf(
		htmlTemplate,
		title,             // %s - title
		bigCSS,            // %s - big.css
		themeCSS,          // %s - theme CSS
		codeCSS,           // %s - code block CSS
		customCSS,         // %s - aspect ratio script
		bigJS,             // %s - big.js
		theme,             // %s - body class (theme)
		slides,            // %s - slides HTML
		fullscreenSnippet, // %s - full-screen toggle button + script
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
