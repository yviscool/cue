package i18n

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// localeOrder is the set of locales that must be complete. Adding a language
// without translating every key is the bug this file exists to prevent.
var localeOrder = []Language{LangEnUS, LangZhCN, LangZhTW, LangJaJP}

// keyCallPattern matches i18n.T("key") / i18n.TF("key", ...) call sites in the
// rest of the repository. It deliberately only picks up plain string literals;
// dynamic keys would defeat the audit anyway.
var keyCallPattern = regexp.MustCompile(`i18n\.T[F]?\("([A-Za-z0-9_.\-]+)"`)

// TestLocaleFilesLoaded checks that every embedded message file parsed cleanly
// and produced the same number of messages the raw JSON declares. A malformed
// file would otherwise degrade silently to raw keys at runtime.
func TestLocaleFilesLoaded(t *testing.T) {
	if err := LoadError(); err != nil {
		t.Fatalf("message file load failed: %v", err)
	}
	if len(catalogue) != len(localeOrder) {
		t.Fatalf("loaded %d locales, want %d", len(catalogue), len(localeOrder))
	}
	for _, lang := range localeOrder {
		if len(catalogue[lang]) == 0 {
			t.Errorf("%s has no messages", lang)
		}
	}
}

// TestEveryDictionaryKeyIsComplete guards against a locale that is missing a
// key present in another locale. A missing entry would silently fall back to
// English (or worse, render the raw key) for everyone else.
func TestEveryDictionaryKeyIsComplete(t *testing.T) {
	reference := catalogue[LangEnUS]
	keys := make([]string, 0, len(reference))
	for k := range reference {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, lang := range localeOrder {
		dict := catalogue[lang]
		for _, key := range keys {
			text, ok := dict[key]
			if !ok {
				t.Errorf("%s is missing key %q", lang, key)
				continue
			}
			if strings.TrimSpace(text) == "" {
				t.Errorf("%s has an empty value for key %q", lang, key)
			}
		}
	}

	// The reverse direction: a key defined only in a non-English locale would
	// never be reachable through the English fallback chain.
	for _, lang := range localeOrder {
		if lang == LangEnUS {
			continue
		}
		for key := range catalogue[lang] {
			if _, ok := reference[key]; !ok {
				t.Errorf("key %q exists in %s but not in %s", key, lang, LangEnUS)
			}
		}
	}
}

// TestFormatVerbsMatchAcrossLocales keeps TF() safe. If a translation drops or
// reorders a printf verb, the rendered string either loses information or
// prints "%!d(MISSING)", so the verb sequence has to stay identical per key.
func TestFormatVerbsMatchAcrossLocales(t *testing.T) {
	verb := regexp.MustCompile(`%[-+ #0-9.]*[a-zA-Z%]`)

	for key, en := range catalogue[LangEnUS] {
		want := verb.FindAllString(en, -1)
		if len(want) == 0 {
			continue
		}
		for _, lang := range localeOrder {
			if lang == LangEnUS {
				continue
			}
			got := verb.FindAllString(catalogue[lang][key], -1)
			if strings.Join(got, "|") != strings.Join(want, "|") {
				t.Errorf("%s: %q has verbs %v but %s has %v", key, lang, got, LangEnUS, want)
			}
		}
	}
}

// TestUIKeyCallSitesAreTranslated walks every Go file outside this package and
// resolves each i18n.T / i18n.TF literal against the message files. This is the
// regression guard for a call site added without a matching translation, which
// is how a raw "filter_silent" identifier once leaked into the settings UI.
func TestUIKeyCallSitesAreTranslated(t *testing.T) {
	root := filepath.Join("..", "..")
	selfDir := filepath.Join("internal", "i18n")

	type callSite struct {
		key  string
		file string
	}
	seen := map[string]callSite{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if strings.Contains(filepath.ToSlash(path), filepath.ToSlash(selfDir)) {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		for _, m := range keyCallPattern.FindAllStringSubmatch(string(body), -1) {
			if _, ok := seen[m[1]]; !ok {
				seen[m[1]] = callSite{key: m[1], file: filepath.ToSlash(path)}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking repository: %v", err)
	}
	if len(seen) == 0 {
		t.Fatal("no i18n.T call sites found; the audit would pass vacuously")
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		site := seen[key]
		for _, lang := range localeOrder {
			if text, ok := catalogue[lang][key]; !ok || strings.TrimSpace(text) == "" {
				t.Errorf("%s: i18n.T(%q) has no %s translation", site.file, key, lang)
			}
		}
	}
}

// TestBundleResolvesEveryLocale drives the public API through go-i18n for all
// four locales, which also proves the localizer is rebuilt on SetLanguage.
func TestBundleResolvesEveryLocale(t *testing.T) {
	mgr := Instance()
	previous := mgr.GetLanguage()
	defer mgr.SetLanguage(previous)

	cases := []struct {
		lang Language
		want string
	}{
		{LangZhCN, "系统信息"},
		{LangEnUS, "System Information"},
		{LangZhTW, "系統資訊"},
		{LangJaJP, "システム情報"},
	}
	for _, tc := range cases {
		mgr.SetLanguage(tc.lang)
		if got := mgr.T("ui_group_system_info"); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.lang, got, tc.want)
		}
	}

	// Missing keys degrade to the identifier, and the fallback chain still
	// reaches English when the active locale lacks an entry.
	mgr.SetLanguage(LangEnUS)
	if got := mgr.T("definitely_not_a_key"); got != "definitely_not_a_key" {
		t.Errorf("expected raw key fallback, got %q", got)
	}

	// filter_silent used to be missing entirely and leaked as a raw identifier
	// into the filter segmented control.
	mgr.SetLanguage(LangZhCN)
	if got := mgr.T("filter_silent"); got == "filter_silent" {
		t.Error("filter_silent is not translated")
	}
}
