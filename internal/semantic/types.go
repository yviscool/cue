package semantic

import "github.com/yviscool/cue/internal/hook"

// ActionInfo describes the human-perceivable meaning of a shortcut.
type ActionInfo struct {
	Title       string // e.g., "Copy", "Command Palette", "Undo"
	Category    string // e.g., "System", "Editor", "Navigation", "Window"
	Description string // Optional longer explanation
}

// Matcher analyzes key combinations and provides human-readable action subtitles.
type Matcher interface {
	// MatchAction returns an ActionInfo if the combo corresponds to a known shortcut.
	MatchAction(appContext string, mods hook.Modifiers, keyName string) (ActionInfo, bool)

	// IsSensitive returns true if current keystrokes should be masked (e.g. password fields).
	IsSensitive(appContext string, keyName string) bool
}
