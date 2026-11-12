package theme

import (
	"regexp"
	"testing"
)

func TestAllThemesCompleteness(t *testing.T) {
	expectedThemes := []string{
		"Liquid Glass",
		"MacBook Scissor",
		"Dieter Rams 1960",
		"Cyberpunk Neon",
		"Azure Pulse",
		"Classic Mechanical",
		"IBM Model M 1984",
		"Minimal Chalk",
		"Obsidian OLED",
		"Tokyo Blossom",
		"Emerald Matrix",
		"Sunset Horizon",
		"Royal Velvet",
		"Sony Walkman 1979",
		"Teenage OP-1",
		"Amber Cathode 1977",
		"Swiss Typography 1957",
		"Kyoto Kintsugi",
		"Blueprint Cyanotype",
		"Leica M",
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

		// Signature Physical Archetype validations
		if thm.DefaultStyle == "" {
			t.Errorf("theme %s has empty DefaultStyle", name)
		}
		if thm.HeldPhysics == "" {
			t.Errorf("theme %s has empty HeldPhysics", name)
		}
		if thm.AcousticProfile == "" {
			t.Errorf("theme %s has empty AcousticProfile", name)
		}
		if thm.EffectiveSpringBounce() <= 0 {
			t.Errorf("theme %s has invalid SpringBounce %f", name, thm.EffectiveSpringBounce())
		}
		if thm.EffectiveDurationScale() <= 0 {
			t.Errorf("theme %s has invalid DurationScale %f", name, thm.EffectiveDurationScale())
		}
	}
}

func TestThemeArchetypePhysics(t *testing.T) {
	// Liquid Glass: optical pulse mode, 0 physical descent
	glass := ThemeLiquidGlass
	if glass.HeldPhysics != HeldPulse {
		t.Errorf("expected Liquid Glass HeldPulse, got %v", glass.HeldPhysics)
	}
	if glass.EffectiveTravelDepth() != 0.0 {
		t.Errorf("expected Liquid Glass 0 travel depth, got %f", glass.EffectiveTravelDepth())
	}

	// MacBook Scissor: shallow, snappy scissor travel (1.2px)
	mac := ThemeMacBookScissor
	if mac.HeldPhysics != HeldSnappy {
		t.Errorf("expected MacBook Scissor HeldSnappy, got %v", mac.HeldPhysics)
	}
	if mac.EffectiveTravelDepth() != 1.2 {
		t.Errorf("expected MacBook Scissor 1.2 travel depth, got %f", mac.EffectiveTravelDepth())
	}

	// IBM Model M: deep buckling spring travel (3.8px)
	modelM := ThemeIBMModelM
	if modelM.EffectiveTravelDepth() != 3.8 {
		t.Errorf("expected IBM Model M 3.8 travel depth, got %f", modelM.EffectiveTravelDepth())
	}
}

// TestExplicitModelOwnsItsPhysics pins the contract the model gallery previews:
// an explicitly chosen model supplies its own travel and held response, while
// Auto keeps the theme's hand-tuned signature.
func TestExplicitModelOwnsItsPhysics(t *testing.T) {
	glass := ThemeLiquidGlass // pulse theme with zero signature travel

	if got := glass.EffectiveTravelDepthFor(""); got != 0 {
		t.Errorf("Auto on Liquid Glass: got %f travel, want 0", got)
	}
	if got := glass.EffectiveTravelDepthFor("pbt"); got != 3.0 {
		t.Errorf("PBT override: got %f travel, want 3.0", got)
	}
	if got := glass.EffectiveHeldPhysicsFor("pbt"); got != HeldSink {
		t.Errorf("PBT override physics: got %v, want sink", got)
	}
	if got := glass.EffectiveHeldPhysicsFor("neon"); got != HeldPulse {
		t.Errorf("Neon override physics: got %v, want pulse", got)
	}

	modelM := ThemeIBMModelM
	if got := modelM.EffectiveTravelDepthFor("laptop"); got != 1.2 {
		t.Errorf("Laptop override: got %f travel, want 1.2", got)
	}
	if got := modelM.EffectiveTravelDepthFor(""); got != 3.8 {
		t.Errorf("Auto on Model M: got %f travel, want 3.8", got)
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

// TestMouseThemesAreComplete guards the cursor presets as a catalogue: the
// settings menu lists them in order, MouseThemeByID resolves every one of them,
// and every colour must be a hex value the overlay can parse. A typo in a
// colour used to show up as an invisible cursor rather than as an error.
func TestMouseThemesAreComplete(t *testing.T) {
	const wantAtLeast = 12
	if len(AllMouseThemes) < wantAtLeast {
		t.Fatalf("expected at least %d mouse themes, got %d", wantAtLeast, len(AllMouseThemes))
	}

	hex := regexp.MustCompile(`^#[0-9a-fA-F]{6}([0-9a-fA-F]{2})?$`)
	ids := make(map[string]bool, len(AllMouseThemes))
	for _, m := range AllMouseThemes {
		if m.ID == "" || m.Name == "" || m.DisplayName == "" {
			t.Errorf("mouse theme %q is missing an identifier or a name", m.ID)
		}
		if ids[m.ID] {
			t.Errorf("duplicate mouse theme id %q", m.ID)
		}
		ids[m.ID] = true

		for label, color := range map[string]string{
			"left":  m.LeftColor,
			"right": m.RightColor,
			"wheel": m.WheelColor,
			"ring":  m.RingColor,
			"badge": m.BadgeColor,
		} {
			if !hex.MatchString(color) {
				t.Errorf("%s: %s colour %q is not hex", m.ID, label, color)
			}
		}
		if got := MouseThemeByID(m.ID); got.ID != m.ID {
			t.Errorf("MouseThemeByID(%q) returned %q", m.ID, got.ID)
		}
	}

	if got := MouseThemeByID("no_such_theme"); got.ID != AllMouseThemes[0].ID {
		t.Errorf("unknown id fell back to %q, want %q", got.ID, AllMouseThemes[0].ID)
	}
}
