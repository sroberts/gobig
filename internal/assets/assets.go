package assets

import (
	"embed"
	"fmt"
)

//go:embed embed/big.js embed/big.css embed/themes/*.css embed/shiki-init.js embed/code-blocks.css
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
func GetTheme(theme string) (string, error) {
	themePath := fmt.Sprintf("embed/themes/%s.css", theme)
	content, err := files.ReadFile(themePath)
	if err != nil {
		return "", fmt.Errorf("failed to read theme %s: %w", theme, err)
	}
	return string(content), nil
}

// ValidateTheme checks if a theme name is valid
func ValidateTheme(theme string) bool {
	validThemes := map[string]bool{
		"dark":  true,
		"light": true,
		"white": true,
	}
	return validThemes[theme]
}

// GetShikiJS returns the Shiki initialization JavaScript content
func GetShikiJS() (string, error) {
	content, err := files.ReadFile("embed/shiki-init.js")
	if err != nil {
		return "", fmt.Errorf("failed to read shiki-init.js: %w", err)
	}
	return string(content), nil
}

// GetCodeBlocksCSS returns the code blocks CSS content
func GetCodeBlocksCSS() (string, error) {
	content, err := files.ReadFile("embed/code-blocks.css")
	if err != nil {
		return "", fmt.Errorf("failed to read code-blocks.css: %w", err)
	}
	return string(content), nil
}
