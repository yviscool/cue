package engine

import (
	"time"

	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/semantic"
)

// FilterMode controls how aggressively key events are filtered.
type FilterMode int

const (
	// FilterSmart only visualizes shortcuts (keys with modifiers) and special keys (Enter, Esc, Tab, F-keys, etc.).
	FilterSmart FilterMode = iota
	// FilterAll visualizes every keystroke including plain alphanumeric typing.
	FilterAll
	// FilterSilent suppresses all visualizations (DND / Presentation Pause).
	FilterSilent
	// FilterCustom only visualizes keystrokes matching the configured AllowedKeys whitelist.
	FilterCustom
)

func (m FilterMode) String() string {
	switch m {
	case FilterSmart:
		return "Smart (Shortcuts & Specials)"
	case FilterAll:
		return "All Keys"
	case FilterSilent:
		return "Silent (Paused)"
	case FilterCustom:
		return "Custom (Allowed Keys)"
	default:
		return "Unknown"
	}
}

// KeyCap represents a single rendered key badge (e.g. "Ctrl", "Shift", "P").
type KeyCap struct {
	Label      string // e.g., "⌘", "Ctrl", "Alt", "⇧", "P", "↵"
	IsModifier bool   // true if this keycap is a modifier key
	IsHeld     bool   // true if this key is physically pressed down right now
}

// DisplayChord represents an active or lingering visualized key cluster.
type DisplayChord struct {
	ID         string              // Unique ID for animation tracking
	Keys       []KeyCap            // e.g. [Ctrl, Shift, P]
	Streak     int                 // Repeat count (e.g. ×3), 1 if single press
	Action     semantic.ActionInfo // Semantic subtitle if matched
	Timestamp  time.Time           // When this chord was first triggered
	ExpiresAt  time.Time           // When this chord will fade out
	IsActive   bool                // True while keys are still physically held down
	IsPrivate  bool                // True if masked by privacy shield
	Velocity   int                 // Cadence velocity: 0: Adagio (<35 WPM), 1: Andante (35-75 WPM), 2: Allegro (>75 WPM)
	IsFinisher bool                // True if this chord is a sequence terminator (Enter, Tab, etc.)
}

// MouseVisualState represents a mouse action visualization.
type MouseVisualState struct {
	Button    hook.MouseButton
	X         int
	Y         int
	Active    bool
	Dragging  bool
	Wheel     int
	Modifiers hook.Modifiers
	Timestamp time.Time
}

// UIState is the immutable snapshot consumed by the MyGo HUD / UI renderer.
type UIState struct {
	Chords       []DisplayChord   // Active chords to display on screen (ordered newest at end or start)
	MouseState   MouseVisualState // Recent mouse action to render
	FilterMode   FilterMode       // Current filter mode
	SoundActive  bool             // True if micro-acoustic feedback is enabled
	ThemeName    string           // Current active theme name
	WordGrouping bool             // True if word grouping is active
	CadenceWPM   int              // Current typing velocity in words per minute
}

// SoundType designates the audio profile to trigger.
type SoundType int

const (
	SoundClick SoundType = iota
	SoundThock
)

// SoundPlayer is an abstract sink for synthesized tactile feedback.
type SoundPlayer interface {
	PlaySound(s SoundType)
}

// SpatialSoundPlayer extends SoundPlayer with key-aware spatial acoustic stage placement.
type SpatialSoundPlayer interface {
	SoundPlayer
	PlaySpatialSound(s SoundType, key string)
}

// SoundPlayerFunc adapter allows simple closures to act as SoundPlayer.
type SoundPlayerFunc func(s SoundType)

func (f SoundPlayerFunc) PlaySound(s SoundType) {
	f(s)
}

// SpatialSoundPlayerFunc adapter allows simple closures to act as SpatialSoundPlayer.
type SpatialSoundPlayerFunc func(s SoundType, key string)

func (f SpatialSoundPlayerFunc) PlaySound(s SoundType) {
	f(s, "")
}

func (f SpatialSoundPlayerFunc) PlaySpatialSound(s SoundType, key string) {
	f(s, key)
}
