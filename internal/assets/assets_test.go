package assets

import (
	"strings"
	"testing"
)

func TestGetBigJS(t *testing.T) {
	content, err := GetBigJS()
	if err != nil {
		t.Fatalf("GetBigJS() failed: %v", err)
	}

	if len(content) == 0 {
		t.Error("GetBigJS() returned empty content")
	}

	// Check for some expected content in big.js
	if !strings.Contains(content, "ASPECT_RATIO") {
		t.Error("big.js does not contain expected ASPECT_RATIO variable")
	}

	if !strings.Contains(content, "function") {
		t.Error("big.js does not contain any functions")
	}
}

func TestGetBigCSS(t *testing.T) {
	content, err := GetBigCSS()
	if err != nil {
		t.Fatalf("GetBigCSS() failed: %v", err)
	}

	if len(content) == 0 {
		t.Error("GetBigCSS() returned empty content")
	}

	// Check for expected CSS content
	if !strings.Contains(content, "body") {
		t.Error("big.css does not contain body styles")
	}
}

func TestGetTheme(t *testing.T) {
	tests := []struct {
		name      string
		theme     string
		wantError bool
	}{
		{
			name:      "dark theme",
			theme:     "dark",
			wantError: false,
		},
		{
			name:      "light theme",
			theme:     "light",
			wantError: false,
		},
		{
			name:      "white theme",
			theme:     "white",
			wantError: false,
		},
		{
			name:      "invalid theme",
			theme:     "invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := GetTheme(tt.theme)

			if tt.wantError {
				if err == nil {
					t.Error("GetTheme() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetTheme() failed: %v", err)
			}

			if len(content) == 0 {
				t.Error("GetTheme() returned empty content")
			}

			// Check for expected CSS content
			if !strings.Contains(content, "body") && !strings.Contains(content, "a") {
				t.Error("theme CSS does not contain expected styles")
			}
		})
	}
}

func TestValidateTheme(t *testing.T) {
	tests := []struct {
		theme string
		want  bool
	}{
		{"dark", true},
		{"light", true},
		{"white", true},
		{"invalid", false},
		{"", false},
		{"Dark", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.theme, func(t *testing.T) {
			got := ValidateTheme(tt.theme)
			if got != tt.want {
				t.Errorf("ValidateTheme(%q) = %v, want %v", tt.theme, got, tt.want)
			}
		})
	}
}
