package theme

import (
	"testing"
)

func TestAllThemesCompleteness(t *testing.T) {
	expectedThemes := []string{
		"Liquid Glass",
		"MacBook Scissor",
		"Dieter Rams 1960",
		"Cyberpunk Neon",
		"Classic Mechanical",
		"IBM Model M 1984",
		"Minimal Chalk",
		"Obsidian OLED",
		"Tokyo Blossom",
		"Emerald Matrix",
		"Sunset Horizon",
		"Royal Velvet",
	}

	if len(AllThemes) < len(expectedThemes) {
		t.Fatalf("expected at least %d themes, got %d", len(expectedThemes), len(AllThemes))
	}

	themeMap := make(map[string]ThemeConfig)
	for _, thm := range AllThemes {
		themeMap[thm.Name] = thm
	}

	for _, name := range expectedThemes {
		thm, exists := themeMap[name]
		if !exists {
			t.Errorf("missing expected theme: %s", name)
			continue
		}

		if thm.KeyCapRadius != 12.5 {
			t.Errorf("theme %s expected KeyCapRadius 12.5, got %f", name, thm.KeyCapRadius)
		}

		if thm.DisplayName == "" {
			t.Errorf("theme %s has empty DisplayName", name)
		}

		if thm.Description == "" {
			t.Errorf("theme %s has empty Description", name)
		}

		// Ensure primary colors have non-zero alpha
		if thm.KeyCapTextColor.A == 0 {
			t.Errorf("theme %s has transparent KeyCapTextColor", name)
		}
		if thm.KeyCapBorder.A == 0 {
			t.Errorf("theme %s has transparent KeyCapBorder", name)
		}
		if thm.KeyCapHeldBackground.A == 0 {
			t.Errorf("theme %s has transparent KeyCapHeldBackground", name)
		}
	}
}

func TestThemeByNameLookup(t *testing.T) {
	tests := []struct {
		query    string
		expected string
	}{
		{"Liquid Glass", "Liquid Glass"},
		{"liquid-glass", "Liquid Glass"},
		{"LIQUID GLASS", "Liquid Glass"},
		{"MacBook Scissor", "MacBook Scissor"},
		{"macbook-scissor", "MacBook Scissor"},
		{"Dieter Rams 1960", "Dieter Rams 1960"},
		{"dieter-rams-1960", "Dieter Rams 1960"},
		{"Cyberpunk Neon", "Cyberpunk Neon"},
		{"cyberpunk-neon", "Cyberpunk Neon"},
		{"Classic Mechanical", "Classic Mechanical"},
		{"classic-mechanical", "Classic Mechanical"},
		{"IBM Model M 1984", "IBM Model M 1984"},
		{"ibm-model-m-1984", "IBM Model M 1984"},
		{"Minimal Chalk", "Minimal Chalk"},
		{"minimal-chalk", "Minimal Chalk"},
		{"Obsidian OLED", "Obsidian OLED"},
		{"obsidian-oled", "Obsidian OLED"},
		{"Tokyo Blossom", "Tokyo Blossom"},
		{"tokyo-blossom", "Tokyo Blossom"},
		{"Emerald Matrix", "Emerald Matrix"},
		{"emerald-matrix", "Emerald Matrix"},
		{"Sunset Horizon", "Sunset Horizon"},
		{"sunset-horizon", "Sunset Horizon"},
		{"Royal Velvet", "Royal Velvet"},
		{"royal-velvet", "Royal Velvet"},
		// Legacy aliases
		{"cupertino-pro", "MacBook Scissor"},
		{"braun-tactile", "Dieter Rams 1960"},
		// Fallbacks
		{"NonExistentTheme", "Liquid Glass"},
		{"", "Liquid Glass"},
		{"   ", "Liquid Glass"},
	}

	for _, tt := range tests {
		got := ThemeByName(tt.query)
		if got.Name != tt.expected {
			t.Errorf("ThemeByName(%q) = %q; expected %q", tt.query, got.Name, tt.expected)
		}

		// Also check GetTheme alias
		aliasGot := GetTheme(tt.query)
		if aliasGot.Name != tt.expected {
			t.Errorf("GetTheme(%q) = %q; expected %q", tt.query, aliasGot.Name, tt.expected)
		}
	}
}

func TestThemeByDisplayName(t *testing.T) {
	for _, thm := range AllThemes {
		got := ThemeByName(thm.DisplayName)
		if got.Name != thm.Name {
			t.Errorf("ThemeByName(DisplayName=%q) = %q; expected %q", thm.DisplayName, got.Name, thm.Name)
		}
	}
}

func TestDefaultStyleForTheme(t *testing.T) {
	if s := DefaultStyleForTheme("Cyberpunk Neon"); s != "neon" {
		t.Errorf("expected neon, got %s", s)
	}
	if s := DefaultStyleForTheme("Liquid Glass"); s != "glass" {
		t.Errorf("expected glass, got %s", s)
	}
	if s := DefaultStyleForTheme("MacBook Scissor"); s != "laptop" {
		t.Errorf("expected laptop, got %s", s)
	}
	if s := DefaultStyleForTheme("Classic Mechanical"); s != "pbt" {
		t.Errorf("expected pbt, got %s", s)
	}
}

func TestColorToHex(t *testing.T) {
	opaque := rgba(255, 0, 128, 255)
	if h := ColorToHex(opaque); h != "#ff0080" {
		t.Errorf("expected #ff0080, got %s", h)
	}
	trans := rgba(255, 255, 255, 42)
	if h := ColorToHex(trans); h != "#ffffff2a" {
		t.Errorf("expected #ffffff2a, got %s", h)
	}
}

