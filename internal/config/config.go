package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Alignment defines the 9-grid screen anchoring position.
type Alignment string

const (
	AlignTopLeft      Alignment = "top-left"
	AlignTopCenter    Alignment = "top-center"
	AlignTopRight     Alignment = "top-right"
	AlignCenterLeft   Alignment = "center-left"
	AlignCenter       Alignment = "center"
	AlignCenterRight  Alignment = "center-right"
	AlignBottomLeft   Alignment = "bottom-left"
	AlignBottomCenter Alignment = "bottom-center"
	AlignBottomRight  Alignment = "bottom-right"
)

// AppearanceSettings controls HUD window placement, animations, and monitor targets.
type AppearanceSettings struct {
	Monitor               string    `json:"monitor"`               // "" / "primary" = primary display, otherwise display ID or label
	FlexDirection         string    `json:"flexDirection"`         // "row" | "column"
	Alignment             Alignment `json:"alignment"`             // default "bottom-right"
	MarginX               int       `json:"marginX"`               // default 100
	MarginY               int       `json:"marginY"`               // default 100
	Animation             string    `json:"animation"`             // "spring" (default) | "none" | "fade" | "zoom" | "float" | "slide"
	AnimationDuration     float64   `json:"animationDuration"`     // default 0.25 (seconds)
	DisablePressAnimation bool      `json:"disablePressAnimation"` // disables keycap sinking & press highlight FX
	Style                 string    `json:"style"`                 // "lowprofile" | "minimal" | "laptop" | "pbt"
	Opacity               float64   `json:"opacity"`               // overall keycap overlay opacity, 1.0 = fully opaque
}

// Opacity bounds. The floor keeps the overlay visible enough to be dragged back
// up: a slider that can reach zero would look like the overlay had vanished.
const (
	MinOpacity     = 0.10
	MaxOpacity     = 1.00
	DefaultOpacity = 1.00
)

// EffectiveOpacity returns the overlay opacity clamped into the usable band, so
// a hand-edited config.json cannot make the HUD invisible or fully opaque when
// the user asked for something in between.
func (a AppearanceSettings) EffectiveOpacity() float64 {
	switch {
	case a.Opacity <= 0:
		return DefaultOpacity
	case a.Opacity < MinOpacity:
		return MinOpacity
	case a.Opacity > MaxOpacity:
		return MaxOpacity
	default:
		return a.Opacity
	}
}

// LayoutSettings defines internal keycap element toggles.
type LayoutSettings struct {
	ShowIcon       bool   `json:"showIcon"`       // default true
	ShowSymbol     bool   `json:"showSymbol"`     // default true
	ShowPressCount bool   `json:"showPressCount"` // default true
	IconAlignment  string `json:"iconAlignment"`  // "flex-start" | "center" | "flex-end"
}

// AppBehaviorSettings groups operating-system behaviour: launch, tray, and the
// privacy boundary around foreground-application inspection.
type AppBehaviorSettings struct {
	StartAtLogin    bool     `json:"startAtLogin"`           // register Cue in the user's login items
	StartHidden     bool     `json:"startHidden"`            // start with the overlay hidden
	TrayClickAction string   `json:"trayClickAction"`        // "settings" | "toggle"
	ExcludedApps    []string `json:"excludedApps,omitempty"` // process names that never contribute context
}

// ColorSettings defines keycap primary/secondary gradient fills.
type ColorSettings struct {
	Color          string `json:"color"`          // default "#ffffff"
	SecondaryColor string `json:"secondaryColor"` // default "#1a1a1a"
	UseGradient    bool   `json:"useGradient"`    // default true
}

// ModifierSettings defines modifier key accenting.
type ModifierSettings struct {
	Highlight      bool   `json:"highlight"`      // default false
	Color          string `json:"color"`          // default "#3a86ff"
	SecondaryColor string `json:"secondaryColor"` // default "#000000"
	TextColor      string `json:"textColor"`      // default "#000000"
	BorderColor    string `json:"borderColor"`    // default "#000000"
}

// TextSettings defines typography within the keycaps.
type TextSettings struct {
	Size      int       `json:"size"`      // default 32
	Color     string    `json:"color"`     // default "#000000"
	Caps      string    `json:"caps"`      // "uppercase" | "capitalize" | "lowercase"
	Variant   string    `json:"variant"`   // "icon" | "text" | "text-short"
	Alignment Alignment `json:"alignment"` // default "center"
}

// BorderSettings defines keycap outline borders.
type BorderSettings struct {
	Enabled bool    `json:"enabled"` // default true
	Color   string  `json:"color"`   // default "#1a1a1a"
	Width   int     `json:"width"`   // default 2
	Radius  float64 `json:"radius"`  // default 0.5 (or pixels)
}

// BackgroundSettings defines overall HUD background.
type BackgroundSettings struct {
	Enabled bool   `json:"enabled"` // default true
	Color   string `json:"color"`   // default "#ffffff99"
}

// MouseSettings mirrors all Cue mouse configuration properties.
type MouseSettings struct {
	ShowClicks       bool   `json:"showClicks"`       // default false
	Size             int    `json:"size"`             // default 150
	Color            string `json:"color"`            // default "#00FF6A"
	ThemePreset      string `json:"themePreset"`      // default "neon_green" ("neon_green", "cyber_cyan", "minimal_white", "obsidian_crimson", "deep_blue", "aurora_bicolor")
	KeepHighlight    bool   `json:"keepHighlight"`    // default false
	ShowIndicator    bool   `json:"showIndicator"`    // default true
	KeepIndicator    bool   `json:"keepIndicator"`    // default true
	IndicatorSize    int    `json:"indicatorSize"`    // default 135
	IndicatorOffsetX int    `json:"indicatorOffsetX"` // default 50
	IndicatorOffsetY int    `json:"indicatorOffsetY"` // default 50
	ClickLingerMs    int    `json:"clickLingerMs"`    // default 180 (ms)
	DragThreshold    int    `json:"dragThreshold"`    // default 50 (px)

	// Opacity is the mouse overlay's own alpha, independent of the keycap
	// overlay: the cursor ring and button indicator sit under the pointer and
	// often need to be fainter than the keycaps. Shares the Appearance opacity
	// bounds so both sliders behave identically.
	Opacity float64 `json:"opacity"` // 1.0 = fully opaque
}

// EffectiveOpacity returns the mouse overlay opacity clamped into the same
// usable band as the keycap overlay, so a hand-edited config.json cannot make
// the cursor ring invisible.
func (m MouseSettings) EffectiveOpacity() float64 {
	switch {
	case m.Opacity <= 0:
		return DefaultOpacity
	case m.Opacity < MinOpacity:
		return MinOpacity
	case m.Opacity > MaxOpacity:
		return MaxOpacity
	default:
		return m.Opacity
	}
}

// CustomShortcut is a user-defined interpretation for a shortcut chord. App may
// be an executable name ("devcpp.exe"), a profile id ("vscode"), a category
// ("editor"), or empty / "*" for a global mapping.
type CustomShortcut struct {
	App         string `json:"app,omitempty"`
	Chord       string `json:"chord"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// GeneralSettings mirrors Cue General tab settings.
type GeneralSettings struct {
	Filter           string   `json:"filter"`                  // "none" | "modifiers" | "custom" | "silent"
	AllowedKeys      []string `json:"allowedKeys"`             // custom allowed raw keys
	ShowEventHistory bool     `json:"showEventHistory"`        // default true
	MaxHistory       int      `json:"maxHistory"`              // default 5 (2 to 12)
	LingerDurationMs int      `json:"lingerDurationMs"`        // default 4500 (ms)
	ToggleShortcut   []string `json:"toggleShortcut"`          // default ["ShiftLeft", "F10"]
	DragThreshold    int      `json:"dragThreshold,omitempty"` // deprecated: migrated to mouse.dragThreshold
	WordGrouping     bool     `json:"wordGrouping"`            // continuous alphanumeric word grouping
	PrivacyMode      bool     `json:"privacyMode"`             // mask keystrokes typed into password, login and banking windows
	ShowInjectedKeys bool     `json:"showInjectedKeys"`        // also show keys synthesized by AutoHotkey / RDP / key remappers
}

// AppConfig is the complete, 100% faithful configuration model matching Cue specifications.
type AppConfig struct {
	// Version is the schema version used for one-time configuration migrations.
	Version int `json:"version"`

	General    GeneralSettings     `json:"general"`
	Appearance AppearanceSettings  `json:"appearance"`
	Layout     LayoutSettings      `json:"layout"`
	Color      ColorSettings       `json:"color"`
	Modifier   ModifierSettings    `json:"modifier"`
	Text       TextSettings        `json:"text"`
	Border     BorderSettings      `json:"border"`
	Background BackgroundSettings  `json:"background"`
	Mouse      MouseSettings       `json:"mouse"`
	App        AppBehaviorSettings `json:"app"`

	// Jobsian Enhanced Extensions (Optional enhancements)
	TactileAudio    bool             `json:"tactileAudio"`    // Micro-acoustics ON/OFF
	AudioProfile    string           `json:"audioProfile"`    // "auto", "thock", "click", "scissor", "glass", "typewriter", "hinoki", "pebble", "paper", "felt", "droplet", "crystal", "chip", "shutter"
	AudioVolume     float64          `json:"audioVolume"`     // master keystroke volume, 1.0 = unity gain
	SemanticSub     bool             `json:"semanticSub"`     // Semantic Action subtitles ON/OFF
	SemanticContext bool             `json:"semanticContext"` // Context-aware shortcut interpretation
	ShowAppBadge    bool             `json:"showAppBadge"`    // Show the recognized application badge in subtitles
	CustomShortcuts []CustomShortcut `json:"customShortcuts,omitempty"`
	ThemePreset     string           `json:"themePreset"` // e.g. "Liquid Glass", "MacBook Scissor", "Dieter Rams 1960"
	Language        string           `json:"language"`    // "zh-CN", "en-US", "zh-TW", "ja-JP", or "" (system default)
}

// Audio volume bounds. 1.0 is unity gain, matching the loudness of the
// synthesized profiles before the slider existed, so upgrading does not change
// how loud the app sounds.
const (
	MinAudioVolume     = 0.0
	MaxAudioVolume     = 1.0
	DefaultAudioVolume = 1.0
)

// EffectiveVolume returns the master keystroke volume clamped to [0, 1]. Zero is
// a legitimate value (muted), so the default for configs written before the
// field existed is supplied by migrateLegacyConfig rather than guessed here.
func (c AppConfig) EffectiveVolume() float64 {
	if c.AudioVolume < MinAudioVolume {
		return MinAudioVolume
	}
	if c.AudioVolume > MaxAudioVolume {
		return MaxAudioVolume
	}
	return c.AudioVolume
}

// configVersion history:
//
//	2 - theme-derived values inherit from the selected preset
//	3 - context intelligence (semanticContext / showAppBadge / custom shortcuts)
//	4 - presentation controls (keycap overlay opacity, keystroke volume)
//	5 - mouse overlay opacity
//	6 - settings IA cleanup: drag threshold moves to the mouse section,
//	    previously inert presentation fields get safe defaults, and app
//	    behaviour (launch, tray, privacy) is introduced.
const configVersion = 6

// DefaultConfig returns the official Cue defaults. Theme-derived values are left
// empty/zero on purpose: the selected theme is the single source of truth for
// keycap color, typography, corner radius, and physical style. Users can still
// override any of them explicitly in the settings studio.
func DefaultConfig() AppConfig {
	return AppConfig{
		Version: configVersion,
		General: GeneralSettings{
			Filter:           "none",
			AllowedKeys:      []string{"ControlLeft", "MetaLeft", "Alt"},
			ShowEventHistory: true,
			MaxHistory:       5,
			LingerDurationMs: 4500,
			ToggleShortcut:   []string{"ShiftLeft", "F10"},
			WordGrouping:     true,
			PrivacyMode:      false,
			ShowInjectedKeys: true,
		},
		Appearance: AppearanceSettings{
			Monitor:               "",
			FlexDirection:         "column",
			Alignment:             AlignBottomRight,
			MarginX:               100,
			MarginY:               100,
			Animation:             "spring",
			AnimationDuration:     0.25,
			DisablePressAnimation: true,
			Style:                 "",
			Opacity:               DefaultOpacity,
		},
		Layout: LayoutSettings{
			ShowIcon:       true,
			ShowSymbol:     true,
			ShowPressCount: true,
			IconAlignment:  "center",
		},
		Color: ColorSettings{
			Color:          "",
			SecondaryColor: "",
			UseGradient:    false,
		},
		Modifier: ModifierSettings{
			Highlight:      false,
			Color:          "#3a86ff",
			SecondaryColor: "#000000",
			TextColor:      "#000000",
			BorderColor:    "#000000",
		},
		Text: TextSettings{
			Size:      0,
			Color:     "",
			Caps:      "capitalize",
			Variant:   "text-short",
			Alignment: AlignCenter,
		},
		Border: BorderSettings{
			Enabled: true,
			Width:   2,
			Color:   "",
			Radius:  0,
		},
		Background: BackgroundSettings{
			Enabled: false,
			Color:   "",
		},
		Mouse: MouseSettings{
			ShowClicks:       false,
			Size:             40,
			Color:            "",
			ThemePreset:      "neon_green",
			KeepHighlight:    false,
			ShowIndicator:    false,
			KeepIndicator:    false,
			IndicatorSize:    135,
			IndicatorOffsetX: 35,
			IndicatorOffsetY: 35,
			ClickLingerMs:    60,
			Opacity:          DefaultOpacity,
			DragThreshold:    50,
		},
		App: AppBehaviorSettings{
			StartAtLogin:    false,
			StartHidden:     false,
			TrayClickAction: "settings",
			ExcludedApps:    nil,
		},
		TactileAudio:    true,
		AudioProfile:    "chip",
		AudioVolume:     DefaultAudioVolume,
		SemanticSub:     true,
		SemanticContext: true,
		ShowAppBadge:    true,
		ThemePreset:     "Swiss Typography 1957",
		Language:        "",
	}
}

// Decode parses a pasted configuration, starting from the factory defaults so a
// partial file keeps every setting it does not mention. It applies the same
// migrations as a file load, so a shared v5 document is safe to import.
// Without the defaults, importing a three-key JSON snippet would silently zero
// sizes, filters, and colors.
func Decode(data []byte) (AppConfig, error) {
	cfg := DefaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{}, fmt.Errorf("parse config: %w", err)
	}
	migrateLegacyConfig(&cfg)
	return cfg, nil
}

// Encode returns the indented JSON representation used on disk and by the
// settings console's clipboard export.
func Encode(cfg AppConfig) ([]byte, error) {
	return json.MarshalIndent(cfg, "", "  ")
}

// Store manages thread-safe configuration file persistence and change watching.
type Store struct {
	mu             sync.RWMutex
	filePath       string
	cfg            AppConfig
	listeners      []func(AppConfig)
	errorListeners []func(error)
	lastHash       [sha256.Size]byte
	hasHash        bool
	badHash        [sha256.Size]byte
	hasBadHash     bool
	lastError      string
	watchStop      chan struct{}
	watchDone      chan struct{}
	watchOnce      sync.Once
	stopOnce       sync.Once
}

// OnChange registers a callback that is invoked whenever configuration updates.
func (s *Store) OnChange(fn func(AppConfig)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listeners = append(s.listeners, fn)
}

// OnError registers a callback for external edits that cannot be loaded. The
// active configuration is kept intact when this fires.
func (s *Store) OnError(fn func(error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errorListeners = append(s.errorListeners, fn)
}

func (s *Store) notifyListenersLocked() {
	c := s.cfg
	for _, fn := range s.listeners {
		cb := fn
		go cb(c)
	}
}

func (s *Store) notifyErrorLocked(err error) {
	for _, fn := range s.errorListeners {
		cb := fn
		go cb(err)
	}
}

// Path returns the configuration file path.
func (s *Store) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.filePath
}

// LastError returns the most recent external reload error, or "" when the file
// is healthy. The error clears automatically after a successful reload.
func (s *Store) LastError() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastError
}

// NewStore initializes or loads config from the specified path.
func NewStore(filePath string) (*Store, error) {
	if filePath == "" {
		home, _ := os.UserHomeDir()
		dir := filepath.Join(home, ".cue")
		_ = os.MkdirAll(dir, 0755)
		filePath = filepath.Join(dir, "config.json")
	}

	s := &Store{
		filePath: filePath,
		cfg:      DefaultConfig(),
	}

	if err := s.load(); err != nil {
		if os.IsNotExist(err) {
			_ = s.Save() // persist defaults only when no file exists yet
		} else {
			// Never overwrite a corrupt user file on startup. Keep defaults in
			// memory and surface the error so the watcher can recover when the
			// user fixes the JSON.
			s.reportError(err, nil)
		}
	}
	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	// Start from the defaults so a key the file does not mention keeps its
	// documented default instead of the Go zero value. Without this, every
	// option added after a user's config was written silently reads as false/0 —
	// which is how "showInjectedKeys" would have arrived switched off for
	// everybody who already had a config file.
	loaded := DefaultConfig()
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	migrated := migrateLegacyConfig(&loaded)
	sum := sha256.Sum256(data)

	s.mu.Lock()
	s.cfg = loaded
	s.lastHash = sum
	s.hasHash = true
	s.lastError = ""
	s.hasBadHash = false
	s.mu.Unlock()

	if migrated {
		if out, err := json.MarshalIndent(loaded, "", "  "); err == nil {
			if err := writeFileAtomic(s.filePath, out); err == nil {
				s.noteSaved(out)
			}
		}
	}
	return nil
}

// Reload re-reads the configuration file. It returns nil without notifying when
// the file content has not changed. Invalid JSON leaves the active config
// untouched and reports through OnError.
func (s *Store) Reload() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		s.reportError(err, nil)
		return err
	}
	sum := sha256.Sum256(data)

	s.mu.RLock()
	unchanged := s.hasHash && sum == s.lastHash
	s.mu.RUnlock()
	if unchanged {
		return nil
	}

	var loaded AppConfig
	if err := json.Unmarshal(data, &loaded); err != nil {
		parseErr := fmt.Errorf("parse config: %w", err)
		s.reportError(parseErr, &sum)
		return parseErr
	}

	migrated := migrateLegacyConfig(&loaded)
	s.mu.Lock()
	s.cfg = loaded
	s.lastHash = sum
	s.hasHash = true
	s.lastError = ""
	s.hasBadHash = false
	s.notifyListenersLocked()
	s.mu.Unlock()

	if migrated {
		if out, err := json.MarshalIndent(loaded, "", "  "); err == nil {
			if err := writeFileAtomic(s.filePath, out); err == nil {
				s.noteSaved(out)
			}
		}
	}
	return nil
}

// Watch polls the configuration file so manual edits are applied without a
// restart. The interval is clamped to a value that keeps the worst-case
// detection latency comfortably below one second.
func (s *Store) Watch(interval time.Duration) {
	if interval <= 0 || interval > 500*time.Millisecond {
		interval = 400 * time.Millisecond
	}
	s.watchOnce.Do(func() {
		stop := make(chan struct{})
		done := make(chan struct{})
		s.mu.Lock()
		s.watchStop = stop
		s.watchDone = done
		s.mu.Unlock()

		go func() {
			defer close(done)
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				select {
				case <-stop:
					return
				case <-ticker.C:
					_ = s.Reload()
				}
			}
		}()
	})
}

// StopWatching stops the polling goroutine.
func (s *Store) StopWatching() {
	s.mu.RLock()
	stop := s.watchStop
	done := s.watchDone
	s.mu.RUnlock()
	if stop == nil || done == nil {
		return
	}
	s.stopOnce.Do(func() { close(stop) })
	<-done
}

func (s *Store) reportError(err error, badHash *[sha256.Size]byte) {
	s.mu.Lock()
	if badHash != nil {
		if s.hasBadHash && *badHash == s.badHash {
			s.mu.Unlock()
			return
		}
		s.badHash = *badHash
		s.hasBadHash = true
	}
	s.lastError = err.Error()
	s.notifyErrorLocked(err)
	s.mu.Unlock()
}

func (s *Store) noteSaved(data []byte) {
	sum := sha256.Sum256(data)
	s.mu.Lock()
	s.lastHash = sum
	s.hasHash = true
	s.lastError = ""
	s.hasBadHash = false
	s.mu.Unlock()
}

// migrateLegacyConfig upgrades configuration written before version 2. Those
// files always stored the old hard-coded Keyviz palette (white keycaps, black
// text, 32 px type, low-profile shape) even when the user had picked a preset
// theme. Treat values that still equal those old defaults as "inherit from the
// theme" so an existing installation gets the same first impression as a new
// one, without discarding genuine customizations.
func migrateLegacyConfig(c *AppConfig) bool {
	if c.Version >= configVersion {
		return false
	}

	migrated := false
	if c.Version < 2 {
		// Version 2: theme-derived values inherit from the selected preset.
		if c.Appearance.Style == "lowprofile" {
			c.Appearance.Style = ""
		}
		if c.Text.Size == 32 {
			c.Text.Size = 0
		}
		if strings.EqualFold(c.Text.Color, "#000000") {
			c.Text.Color = ""
		}
		if strings.EqualFold(c.Color.Color, "#ffffff") {
			c.Color.Color = ""
		}
		if strings.EqualFold(c.Color.SecondaryColor, "#1a1a1a") {
			c.Color.SecondaryColor = ""
		}
		if c.Border.Radius == 0.5 {
			c.Border.Radius = 0
		}
		if strings.EqualFold(c.Border.Color, "#1a1a1a") {
			c.Border.Color = ""
		}

		// Canonicalize the legacy slug names so the settings picker can show a
		// checkmark on the active preset.
		switch strings.ToLower(strings.TrimSpace(c.ThemePreset)) {
		case "", "liquid-glass":
			c.ThemePreset = "Liquid Glass"
		case "cupertino-pro":
			c.ThemePreset = "MacBook Scissor"
		case "braun-tactile":
			c.ThemePreset = "Dieter Rams 1960"
		}
		c.Version = 2
		migrated = true
	}

	if c.Version < 3 {
		// Version 3: context intelligence defaults on for existing installs.
		c.SemanticContext = true
		c.ShowAppBadge = true
		c.Version = 3
		migrated = true
	}

	if c.Version < 4 {
		// Version 4: the opacity and volume sliders did not exist before, so a
		// zero value here means "unset". Both default to their neutral position:
		// fully opaque and unity gain, which preserves the previous look and
		// loudness exactly.
		if c.Appearance.Opacity <= 0 {
			c.Appearance.Opacity = DefaultOpacity
		}
		if c.AudioVolume <= 0 {
			c.AudioVolume = DefaultAudioVolume
		}
		c.Version = 4
		migrated = true
	}

	if c.Version < 5 {
		// Version 5: the mouse overlay gained its own opacity control. The zero
		// value means "unset" (the slider bottoms out at MinOpacity), so existing
		// installs start fully opaque and keep their previous appearance.
		if c.Mouse.Opacity <= 0 {
			c.Mouse.Opacity = DefaultOpacity
		}
		c.Version = 5
		migrated = true
	}

	if c.Version < 6 {
		// Version 6: fields that existed in the schema but were never rendered
		// get the values that preserve the pre-v6 look. Drag threshold moves to
		// the mouse section, where the gesture actually lives.
		// The deprecated key wins over the prefilled default: it is the value the
		// user actually chose in v5, while mouse.dragThreshold is only there
		// because the struct started from DefaultConfig.
		if c.General.DragThreshold > 0 {
			c.Mouse.DragThreshold = c.General.DragThreshold
		}
		if c.Mouse.DragThreshold <= 0 {
			c.Mouse.DragThreshold = 50
		}
		c.General.DragThreshold = 0
		c.Color.UseGradient = false
		c.Layout.IconAlignment = "center"
		c.Text.Alignment = AlignCenter
		c.Background.Enabled = false
		if c.Mouse.Color == "#00FF6A" {
			// Before v6 every theme selection copied the neon-green seed into
			// mouse.color, which silently disabled the theme's own colors. Empty
			// now means "follow the selected mouse theme".
			c.Mouse.Color = ""
		}
		if c.App.TrayClickAction == "" {
			c.App.TrayClickAction = "settings"
		}
		c.Version = 6
		migrated = true
	}

	c.Version = configVersion
	return migrated
}

// Get returns an immutable copy of the active configuration.
func (s *Store) Get() AppConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

// Update modifies configuration with a callback and saves to disk. The
// in-memory configuration is only replaced after the write succeeds, so a
// failed save cannot leave the UI and disk out of sync.
func (s *Store) Update(fn func(*AppConfig)) error {
	s.mu.Lock()
	next := s.cfg
	fn(&next)
	data, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		s.mu.Unlock()
		return err
	}
	path := s.filePath
	if err := writeFileAtomic(path, data); err != nil {
		s.mu.Unlock()
		return err
	}
	s.cfg = next
	s.lastHash = sha256.Sum256(data)
	s.hasHash = true
	s.lastError = ""
	s.hasBadHash = false
	s.notifyListenersLocked()
	s.mu.Unlock()
	return nil
}

// Save writes current config to disk.
func (s *Store) Save() error {
	s.mu.Lock()
	data, err := json.MarshalIndent(s.cfg, "", "  ")
	if err != nil {
		s.mu.Unlock()
		return err
	}
	path := s.filePath
	if err := writeFileAtomic(path, data); err != nil {
		s.mu.Unlock()
		return err
	}
	s.lastHash = sha256.Sum256(data)
	s.hasHash = true
	s.mu.Unlock()
	return nil
}

// writeFileAtomic writes through a temporary file so the watcher can never read
// a half-written configuration.
func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".cue-config-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	_ = os.Chmod(tmpName, 0644)
	return os.Rename(tmpName, path)
}
