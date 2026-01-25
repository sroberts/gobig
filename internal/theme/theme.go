package theme

import (
	"bytes"
	"fmt"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Theme represents a presentation theme loaded from YAML
type Theme struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Colors      Colors      `yaml:"colors"`
	Typography  Typography  `yaml:"typography"`
	Links       Links       `yaml:"links"`
	CustomCSS   string      `yaml:"custom_css"`
}

// Colors defines the color scheme for the theme
type Colors struct {
	Background          string `yaml:"background"`
	Foreground          string `yaml:"foreground"`
	Emphasis            string `yaml:"emphasis"`
	Link                string `yaml:"link"`
	LinkHoverBackground string `yaml:"link_hover_background"`
}

// Typography defines font settings for the theme
type Typography struct {
	FontFamily string `yaml:"font_family"`
	FontWeight int    `yaml:"font_weight"`
}

// Links defines link styling options
type Links struct {
	TextDecorationSkip string `yaml:"text_decoration_skip"`
}

// cssTemplate is the template for generating CSS from theme data
const cssTemplate = `body {
  background: {{.Colors.Background}};
  color: {{.Colors.Foreground}};
{{- if .Typography.FontFamily}}
  font-family: {{.Typography.FontFamily}};
{{- end}}
{{- if .Typography.FontWeight}}
  font-weight: {{.Typography.FontWeight}};
{{- end}}
}
{{if .Colors.Emphasis}}
em {
  color: {{.Colors.Emphasis}};
}
{{end}}
{{if .Colors.Link}}
a {
  color: {{.Colors.Link}};
{{- if .Links.TextDecorationSkip}}
  text-decoration-skip: {{.Links.TextDecorationSkip}};
{{- end}}
}
{{end}}
{{if .Colors.LinkHoverBackground}}
a:hover {
  background: {{.Colors.LinkHoverBackground}};
}
{{end}}
{{.CustomCSS}}`

// LoadFromYAML parses a YAML theme file and returns a Theme
func LoadFromYAML(data []byte) (*Theme, error) {
	var t Theme
	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("failed to parse theme YAML: %w", err)
	}
	return &t, nil
}

// ToCSS generates CSS from the theme
func (t *Theme) ToCSS() (string, error) {
	tmpl, err := template.New("theme").Parse(cssTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse CSS template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, t); err != nil {
		return "", fmt.Errorf("failed to generate CSS: %w", err)
	}

	return buf.String(), nil
}
