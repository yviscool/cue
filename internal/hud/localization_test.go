package hud

import (
	"regexp"
	"strings"
	"testing"

	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/theme"
)

var cjkRun = regexp.MustCompile(`[\x{4e00}-\x{9fff}]`)

// TestEveryThemeHasLocalizedLabels guards the settings dropdowns and the quick
// flyout. Both used to render theme.DisplayName directly, which is Chinese-first
// by design; adding a theme without a matching entry in themeNameKeys would
// therefore leak Chinese into the English and Japanese interfaces.
func TestEveryThemeHasLocalizedLabels(t *testing.T) {
	mgr := i18n.Instance()
	previous := mgr.GetLanguage()
	defer mgr.SetLanguage(previous)

	for _, th := range theme.AllThemes {
		name := th.Name
		if _, ok := themeNameKeys[strings.ToLower(strings.TrimSpace(name))]; !ok {
			t.Errorf("theme %q has no localization key", name)
			continue
		}

		mgr.SetLanguage(i18n.LangEnUS)
		enName := themeDisplayName(name)
		enDesc := themeDescription(name)
		if cjkRun.MatchString(enName) {
			t.Errorf("theme %q still renders Chinese in en-US: %q", name, enName)
		}
		if cjkRun.MatchString(enDesc) {
			t.Errorf("theme %q description still renders Chinese in en-US: %q", name, enDesc)
		}

		// Japanese legitimately uses kanji, so "no CJK" only holds for English.
		// For Japanese compare against the Chinese value instead: an identical
		// string means the entry was copy-pasted rather than translated.
		mgr.SetLanguage(i18n.LangZhCN)
		zhName := themeDisplayName(name)
		mgr.SetLanguage(i18n.LangJaJP)
		if jaName := themeDisplayName(name); jaName == zhName {
			t.Errorf("theme %q has an untranslated ja-JP label: %q", name, jaName)
		}
	}
}

// TestEveryMouseThemeHasLocalizedLabels is the same guard for the mouse tab.
func TestEveryMouseThemeHasLocalizedLabels(t *testing.T) {
	mgr := i18n.Instance()
	previous := mgr.GetLanguage()
	defer mgr.SetLanguage(previous)

	for _, mt := range theme.AllMouseThemes {
		if _, ok := mouseThemeKeys[mt.ID]; !ok {
			t.Errorf("mouse theme %q has no localization key", mt.ID)
			continue
		}
		mgr.SetLanguage(i18n.LangEnUS)
		if got := mouseThemeDisplayName(mt); cjkRun.MatchString(got) {
			t.Errorf("mouse theme %q still renders Chinese in en-US: %q", mt.ID, got)
		}
	}
}

// TestPhysicalStyleNamesAreLocalized covers the keycap model names that appear
// in the "· archetype preset" summary line.
func TestPhysicalStyleNamesAreLocalized(t *testing.T) {
	mgr := i18n.Instance()
	previous := mgr.GetLanguage()
	defer mgr.SetLanguage(previous)

	ids := []string{"", "lowprofile", "minimal", "laptop", "pbt", "glass", "neon", "retro", "outline"}
	mgr.SetLanguage(i18n.LangEnUS)
	for _, id := range ids {
		if got := physicalStyleName(id); got == "" || cjkRun.MatchString(got) {
			t.Errorf("physical style %q renders %q in en-US", id, got)
		}
	}
}

// TestTravelDepthLabelIsLocalized checks both branches: the kinetic pulse label
// and the millimetre readout.
func TestTravelDepthLabelIsLocalized(t *testing.T) {
	mgr := i18n.Instance()
	previous := mgr.GetLanguage()
	defer mgr.SetLanguage(previous)

	mgr.SetLanguage(i18n.LangEnUS)

	pulse := theme.ThemeConfig{HeldPhysics: theme.HeldPulse}
	if got := travelDepthLabel(pulse); cjkRun.MatchString(got) || got == "" {
		t.Errorf("held-pulse label not localized: %q", got)
	}

	travel := theme.ThemeConfig{HeldPhysics: theme.HeldSink, TravelDepth: 2.0}
	got := travelDepthLabel(travel)
	if cjkRun.MatchString(got) {
		t.Errorf("travel depth label not localized: %q", got)
	}
	if !strings.Contains(got, "2.0") {
		t.Errorf("travel depth label lost its value: %q", got)
	}
}
