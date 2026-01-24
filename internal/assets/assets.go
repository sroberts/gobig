package assets

import (
	"embed"
	"fmt"

	"gobig/internal/theme"
)

//go:embed embed/big.js embed/big.css embed/themes/*.css embed/themes/*.yaml
var files embed.FS

// GetBigJS returns the big.js JavaScript content
func GetBigJS() (string, error) {
	content, err := files.ReadFile("embed/big.js")
	if err != nil {
		return "", fmt.Errorf("failed to read big.js: %w", err)
	}
	return string(content), nil
}

// GetBigCSS returns the big.css stylesheet content
func GetBigCSS() (string, error) {
	content, err := files.ReadFile("embed/big.css")
	if err != nil {
		return "", fmt.Errorf("failed to read big.css: %w", err)
	}
	return string(content), nil
}

// GetTheme returns the theme CSS content for the specified theme
// Valid themes: "dark", "light", "white"
// Tries to load YAML theme first, falls back to CSS if not found
func GetTheme(themeName string) (string, error) {
	// Try YAML theme first
	yamlPath := fmt.Sprintf("embed/themes/%s.yaml", themeName)
	yamlContent, err := files.ReadFile(yamlPath)
	if err == nil {
		// Parse YAML and generate CSS
		t, err := theme.LoadFromYAML(yamlContent)
		if err != nil {
			return "", fmt.Errorf("failed to parse theme YAML %s: %w", themeName, err)
		}
		css, err := t.ToCSS()
		if err != nil {
			return "", fmt.Errorf("failed to generate CSS for theme %s: %w", themeName, err)
		}
		return css, nil
	}

	// Fall back to CSS theme (for backwards compatibility)
	cssPath := fmt.Sprintf("embed/themes/%s.css", themeName)
	cssContent, err := files.ReadFile(cssPath)
	if err != nil {
		return "", fmt.Errorf("failed to read theme %s: %w", themeName, err)
	}
	return string(cssContent), nil
}

// ValidateTheme checks if a theme name is valid by checking if the theme file exists
func ValidateTheme(themeName string) bool {
	// Try YAML theme first
	yamlPath := fmt.Sprintf("embed/themes/%s.yaml", themeName)
	if _, err := files.ReadFile(yamlPath); err == nil {
		return true
	}

	// Fall back to CSS theme
	cssPath := fmt.Sprintf("embed/themes/%s.css", themeName)
	if _, err := files.ReadFile(cssPath); err == nil {
		return true
	}

	return false
}

// ListAvailableThemes returns a list of all available theme names
func ListAvailableThemes() ([]string, error) {
	entries, err := files.ReadDir("embed/themes")
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	themes := make([]string, 0, len(entries))
	seen := make(map[string]bool)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Remove .yaml or .css extension
		if len(name) > 5 && name[len(name)-5:] == ".yaml" {
			themeName := name[:len(name)-5]
			if !seen[themeName] {
				themes = append(themes, themeName)
				seen[themeName] = true
			}
		} else if len(name) > 4 && name[len(name)-4:] == ".css" {
			themeName := name[:len(name)-4]
			if !seen[themeName] {
				themes = append(themes, themeName)
				seen[themeName] = true
			}
		}
	}

	return themes, nil
}
