package i18n

import (
	"strings"
	"testing"

	"github.com/yviscool/cue/internal/semantic"
)

// TestActionTranslations guarantees that every action title in the semantic
// registry is localized in all four supported languages. This is what keeps the
// action vocabulary complete as new application profiles are added.
func TestActionTranslations(t *testing.T) {
	matcher := semantic.NewMatcher()
	titles := matcher.AllActionTitles()
	if len(titles) == 0 {
		t.Fatal("semantic registry returned no action titles")
	}

	for _, lang := range []Language{LangEnUS, LangZhCN, LangZhTW, LangJaJP} {
		dict := catalogue[lang]
		for _, title := range titles {
			text, ok := dict[title]
			if !ok {
				t.Errorf("%s: missing action translation for %q", lang, title)
				continue
			}
			if strings.TrimSpace(text) == "" {
				t.Errorf("%s: empty action translation for %q", lang, title)
			}
		}
	}
}

func TestContextSettingsTranslations(t *testing.T) {
	keys := []string{
		"ctx_group_title",
		"ctx_smart_title",
		"ctx_smart_desc",
		"ctx_sub_title",
		"ctx_badge_title",
		"ctx_current_app",
		"ctx_unknown_app",
		"ctx_custom_title",
	}
	for _, lang := range []Language{LangEnUS, LangZhCN, LangZhTW, LangJaJP} {
		dict := catalogue[lang]
		for _, key := range keys {
			if text, ok := dict[key]; !ok || strings.TrimSpace(text) == "" {
				t.Errorf("%s: missing UI translation for %q", lang, key)
			}
		}
	}
}
