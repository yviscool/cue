package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	Monitor           string    `json:"monitor"`
	FlexDirection     string    `json:"flexDirection"`     // "row" | "column"
	Alignment         Alignment `json:"alignment"`         // default "bottom-center"
	MarginX           int       `json:"marginX"`           // default 100
	MarginY           int       `json:"marginY"`           // default 100
	Animation         string    `json:"animation"`         // "none" | "fade" | "zoom" | "float" | "slide"
	AnimationDuration float64   `json:"animationDuration"` // default 0.25 (seconds)
	Style             string    `json:"style"`             // "lowprofile" | "minimal" | "laptop" | "pbt"
}

// LayoutSettings defines internal keycap element toggles.
type LayoutSettings struct {
	ShowIcon       bool   `json:"showIcon"`       // default true
	ShowSymbol     bool   `json:"showSymbol"`     // default true
	ShowPressCount bool   `json:"showPressCount"` // default true
	IconAlignment  string `json:"iconAlignment"`  // "flex-start" | "center" | "flex-end"
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
	Color            string `json:"color"`            // default "#009dff"
	KeepHighlight    bool   `json:"keepHighlight"`    // default false
	ShowIndicator    bool   `json:"showIndicator"`    // default true
	KeepIndicator    bool   `json:"keepIndicator"`    // default true
	IndicatorSize    int    `json:"indicatorSize"`    // default 50
	IndicatorOffsetX int    `json:"indicatorOffsetX"` // default 50
	IndicatorOffsetY int    `json:"indicatorOffsetY"` // default 50
}

// GeneralSettings mirrors Cue General tab settings.
type GeneralSettings struct {
	Filter           string   `json:"filter"`           // "none" | "modifiers" | "custom"
	AllowedKeys      []string `json:"allowedKeys"`      // custom allowed raw keys
	ShowEventHistory bool     `json:"showEventHistory"` // default false
	MaxHistory       int      `json:"maxHistory"`       // default 5 (2 to 12)
	LingerDurationMs int      `json:"lingerDurationMs"` // default 5000 (ms)
	ToggleShortcut   []string `json:"toggleShortcut"`   // default ["ShiftLeft", "F10"]
	DragThreshold    int      `json:"dragThreshold"`    // default 50
	WordGrouping     bool     `json:"wordGrouping"`     // continuous alphanumeric word grouping
}

// AppConfig is the complete, 100% faithful configuration model matching Cue specifications.
type AppConfig struct {
	// Version is the schema version used for one-time configuration migrations.
	Version int `json:"version"`

	General    GeneralSettings    `json:"general"`
	Appearance AppearanceSettings `json:"appearance"`
	Layout     LayoutSettings     `json:"layout"`
	Color      ColorSettings      `json:"color"`
	Modifier   ModifierSettings   `json:"modifier"`
	Text       TextSettings       `json:"text"`
	Border     BorderSettings     `json:"border"`
	Background BackgroundSettings `json:"background"`
	Mouse      MouseSettings      `json:"mouse"`

	// Jobsian Enhanced Extensions (Optional enhancements)
	TactileAudio bool   `json:"tactileAudio"` // Micro-acoustics ON/OFF
	SemanticSub  bool   `json:"semanticSub"`  // Semantic Action subtitles ON/OFF
	ThemePreset  string `json:"themePreset"`  // e.g. "Liquid Glass", "MacBook Scissor", "Dieter Rams 1960"
	Language     string `json:"language"`     // "zh-CN", "en-US", "zh-TW", "ja-JP", or "" (system default)
}

// configVersion 2 made theme-derived values inherit from the selected preset
// instead of hard-coding the old Keyviz palette into every new installation.
const configVersion = 2

// DefaultConfig returns the official Cue defaults. Theme-derived values are left
// empty/zero on purpose: the selected theme is the single source of truth for
// keycap color, typography, corner radius, and physical style. Users can still
// override any of them explicitly in the settings studio.
func DefaultConfig() AppConfig {
	return AppConfig{
		Version: configVersion,
		General: GeneralSettings{
			Filter:           "modifiers",
			AllowedKeys:      []string{"ControlLeft", "MetaLeft", "Alt"},
			ShowEventHistory: false,
			MaxHistory:       5,
			LingerDurationMs: 5000,
			ToggleShortcut:   []string{"ShiftLeft", "F10"},
			DragThreshold:    50,
			WordGrouping:     true,
		},
		Appearance: AppearanceSettings{
			Monitor:           "",
			FlexDirection:     "column",
			Alignment:         AlignBottomCenter,
			MarginX:           100,
			MarginY:           100,
			Animation:         "fade",
			AnimationDuration: 0.25,
			Style:             "",
		},
		Layout: LayoutSettings{
			ShowIcon:       true,
			ShowSymbol:     true,
			ShowPressCount: true,
			IconAlignment:  "flex-end",
		},
		Color: ColorSettings{
			Color:          "",
			SecondaryColor: "",
			UseGradient:    true,
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
			Enabled: true,
			Color:   "#ffffff99",
		},
		Mouse: MouseSettings{
			ShowClicks:       false,
			Size:             150,
			Color:            "#009dff",
			KeepHighlight:    false,
			ShowIndicator:    false,
			KeepIndicator:    false,
			IndicatorSize:    50,
			IndicatorOffsetX: 50,
			IndicatorOffsetY: 50,
		},
		TactileAudio: true,
		SemanticSub:  true,
		ThemePreset:  "Liquid Glass",
		Language:     "",
	}
}

// Store manages thread-safe configuration file persistence.
type Store struct {
	mu        sync.RWMutex
	filePath  string
	cfg       AppConfig
	listeners []func(AppConfig)
}

// OnChange registers a callback that is invoked whenever configuration updates.
func (s *Store) OnChange(fn func(AppConfig)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listeners = append(s.listeners, fn)
}

func (s *Store) notifyListenersLocked() {
	c := s.cfg
	for _, fn := range s.listeners {
		cb := fn
		go cb(c)
	}
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
		_ = s.Save() // persist defaults if file didn't exist
	}
	return s, nil
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	var loaded AppConfig
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}
	if migrateLegacyConfig(&loaded) {
		s.cfg = loaded
		if migrated, err := json.MarshalIndent(s.cfg, "", "  "); err == nil {
			_ = os.WriteFile(s.filePath, migrated, 0644)
		}
		return nil
	}
	s.cfg = loaded
	return nil
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

	c.Version = configVersion
	return true
}

// Get returns an immutable copy of the active configuration.
func (s *Store) Get() AppConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

// Update modifies configuration with a callback and saves to disk.
func (s *Store) Update(fn func(*AppConfig)) error {
	s.mu.Lock()
	fn(&s.cfg)
	data, err := json.MarshalIndent(s.cfg, "", "  ")
	if err != nil {
		s.mu.Unlock()
		return err
	}
	path := s.filePath
	s.notifyListenersLocked()
	s.mu.Unlock()
	return os.WriteFile(path, data, 0644)
}

// Save writes current config to disk.
func (s *Store) Save() error {
	s.mu.RLock()
	data, err := json.MarshalIndent(s.cfg, "", "  ")
	path := s.filePath
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
