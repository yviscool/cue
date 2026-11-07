package semantic

import (
	"time"

	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// DefaultSequenceWindow is how long a two-chord sequence stays valid.
const DefaultSequenceWindow = 1500 * time.Millisecond

// ActionInfo describes the human-perceivable meaning of a shortcut.
// Title doubles as the stable localization key; i18n.T(Title) resolves it.
type ActionInfo struct {
	Title       string // e.g., "Copy", "Command Palette", "Undo"
	Category    string // e.g., "System", "Editor", "Navigation", "Window"
	Description string // Optional longer explanation
}

// Confidence describes how specific the match is. It lets the HUD and future
// UI layers distinguish an app-exact interpretation from a generic fallback.
type Confidence uint8

const (
	ConfidenceNone Confidence = iota
	ConfidenceUniversal
	ConfidenceCategory
	ConfidenceProfile
	ConfidenceUser
)

func (c Confidence) String() string {
	switch c {
	case ConfidenceUniversal:
		return "universal"
	case ConfidenceCategory:
		return "category"
	case ConfidenceProfile:
		return "profile"
	case ConfidenceUser:
		return "user"
	default:
		return "none"
	}
}

// MatchResult is the semantic interpretation of one chord.
type MatchResult struct {
	Action      ActionInfo
	Confidence  Confidence
	ProfileID   string
	ProfileName string
	AppName     string
	IsSequence  bool
}

// Matched reports whether any action was recognized.
func (m MatchResult) Matched() bool {
	return m.Action.Title != ""
}

// Matcher analyzes key combinations and provides human-readable action subtitles.
type Matcher interface {
	// MatchAction returns the best interpretation for one chord.
	MatchAction(app context.AppContext, mods hook.Modifiers, keyName string) MatchResult

	// MatchUniversal resolves a chord without application context. It is used
	// when context intelligence is disabled: universal shortcuts and global
	// user overrides keep working, while app profiles and categories are
	// skipped. It never touches the context cache, so privacy detection stays
	// coherent with the real foreground application.
	MatchUniversal(mods hook.Modifiers, keyName string) MatchResult

	// MatchSequence returns an interpretation for a two-chord sequence such as
	// Ctrl+K Ctrl+S. firstKey/secondKey are the primary keys, not modifiers.
	MatchSequence(app context.AppContext, firstMods hook.Modifiers, firstKey string, secondMods hook.Modifiers, secondKey string) (MatchResult, bool)

	// IsSensitive returns true if current keystrokes should be masked (e.g. password fields).
	IsSensitive(app context.AppContext, keyName string) bool
}

// Override is a user-defined shortcut interpretation. App may be an executable
// name ("code.exe"), a profile id ("vscode"), a category ("editor"), or empty
// for a global mapping.
type Override struct {
	App    string
	Mods   hook.Modifiers
	Key    string
	Action ActionInfo
}
