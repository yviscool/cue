package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Language represents a supported localization language.
type Language string

const (
	LangEnUS Language = "en-US"
	LangZhCN Language = "zh-CN"
	LangZhTW Language = "zh-TW"
	LangJaJP Language = "ja-JP"
)

// SupportedLanguages lists the locales in the order they appear in the language
// pickers. The Name is deliberately written in its own script: language menus
// are conventionally presented in the target language, not the current one.
var SupportedLanguages = []struct {
	Code Language
	Name string
}{
	{LangZhCN, "简体中文 (Simplified Chinese)"},
	{LangEnUS, "English (US)"},
	{LangZhTW, "繁體中文 (Traditional Chinese)"},
	{LangJaJP, "日本語 (Japanese)"},
}

// LocaleDict is the flat message table for one locale. It mirrors the on-disk
// message file one-to-one and is kept in memory for diagnostics and tests.
type LocaleDict map[string]string

// localeFS holds the message files. They are the single source of truth for
// every user-visible string: the settings studio, the tray, the flyout, and the
// action subtitles all resolve through the bundle built from these files.
//
//go:embed locales/*.json
var localeFS embed.FS

const localeDir = "locales"

var (
	// bundle is the go-i18n message store. English is the default language, so
	// a key missing from the active locale falls back to English and only then
	// to the raw key.
	bundle = goi18n.NewBundle(language.English)
	// catalogue exposes the parsed tables per locale for tests and diagnostics.
	catalogue = make(map[Language]LocaleDict, len(SupportedLanguages))
	// loadErr records the first message-file problem. A broken translation file
	// must not take the whole application down, so the manager degrades to
	// returning raw keys and the error is surfaced by the test suite.
	loadErr error
)

func init() {
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, entry := range SupportedLanguages {
		path := localeDir + "/active." + string(entry.Code) + ".json"
		raw, err := localeFS.ReadFile(path)
		if err != nil {
			recordLoadError(fmt.Errorf("reading %s: %w", path, err))
			continue
		}
		dict := make(LocaleDict)
		if err := json.Unmarshal(raw, &dict); err != nil {
			recordLoadError(fmt.Errorf("decoding %s: %w", path, err))
			continue
		}
		if _, err := bundle.ParseMessageFileBytes(raw, path); err != nil {
			recordLoadError(fmt.Errorf("loading %s into bundle: %w", path, err))
			continue
		}
		catalogue[entry.Code] = dict
	}
}

func recordLoadError(err error) {
	if loadErr == nil {
		loadErr = err
	}
}

// LoadError reports the first problem encountered while loading the message
// files, or nil when every locale loaded cleanly.
func LoadError() error { return loadErr }

// I18nManager coordinates the active locale and resolves messages through
// go-i18n. The localizer is rebuilt whenever the language changes, which is
// cheap (it only holds a tag list) and keeps reads lock-free.
type I18nManager struct {
	mu        sync.RWMutex
	current   Language
	localizer *goi18n.Localizer
}

var (
	globalManager *I18nManager
	once          sync.Once
)

// Instance returns the singleton I18nManager.
func Instance() *I18nManager {
	once.Do(func() {
		current := detectDefaultLanguage()
		globalManager = &I18nManager{
			current:   current,
			localizer: goi18n.NewLocalizer(bundle, string(current)),
		}
	})
	return globalManager
}

// detectDefaultLanguage queries the Windows system locale and maps it onto a
// supported language, falling back to Simplified Chinese.
func detectDefaultLanguage() Language {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetUserDefaultLocaleName")
	if proc != nil {
		buf := make([]uint16, 85)
		r, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if r > 0 {
			loc := strings.ToLower(syscall.UTF16ToString(buf))
			switch {
			case strings.HasPrefix(loc, "zh-tw"), strings.HasPrefix(loc, "zh-hk"), strings.HasPrefix(loc, "zh-mo"):
				return LangZhTW
			case strings.HasPrefix(loc, "zh"):
				return LangZhCN
			case strings.HasPrefix(loc, "ja"):
				return LangJaJP
			}
		}
	}
	return LangZhCN
}

// SetLanguage updates the active locale and rebuilds the localizer.
func (m *I18nManager) SetLanguage(lang Language) {
	if _, ok := catalogue[lang]; !ok {
		return
	}
	loc := goi18n.NewLocalizer(bundle, string(lang))
	m.mu.Lock()
	defer m.mu.Unlock()
	m.current = lang
	m.localizer = loc
}

// GetLanguage returns the active locale.
func (m *I18nManager) GetLanguage() Language {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

// T translates a key into the active locale.
//
// Resolution order is handled by go-i18n: active locale, then the bundle's
// default language (English), then the default message, which is the key
// itself. A missing translation therefore degrades to the raw identifier rather
// than an empty label.
func (m *I18nManager) T(key string) string {
	m.mu.RLock()
	loc := m.localizer
	m.mu.RUnlock()
	if loc == nil {
		return key
	}
	// go-i18n reports MessageNotFoundErr together with a usable fallback string
	// whenever it had to leave the requested locale, so the returned text is
	// authoritative and only an empty result means "nothing to show".
	text, _ := loc.Localize(&goi18n.LocalizeConfig{
		DefaultMessage: &goi18n.Message{ID: key, Other: key},
	})
	if text == "" {
		return key
	}
	return text
}

// Has reports whether key exists in the given locale.
func (m *I18nManager) Has(lang Language, key string) bool {
	_, ok := catalogue[lang][key]
	return ok
}

// T is the global package helper.
func T(key string) string {
	return Instance().T(key)
}

// TF is the global package helper with printf-style formatting applied to the
// resolved message.
func TF(key string, args ...any) string {
	raw := Instance().T(key)
	if len(args) == 0 {
		return raw
	}
	return fmt.Sprintf(raw, args...)
}
