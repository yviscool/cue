package main

import (
	"testing"

	"github.com/egoist/mygo"
	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

func TestHistoryChordLimit(t *testing.T) {
	cfg := config.DefaultConfig()

	// The shipping default keeps history on, so the stage holds MaxHistory cues.
	if got := historyChordLimit(cfg); got != cfg.General.MaxHistory {
		t.Errorf("history on should yield %d chords, got %d", cfg.General.MaxHistory, got)
	}

	// History off collapses the stage to a single cue.
	cfg.General.ShowEventHistory = false
	if got := historyChordLimit(cfg); got != 1 {
		t.Errorf("history off should yield 1 chord, got %d", got)
	}

	cfg.General.ShowEventHistory = true
	cfg.General.MaxHistory = 7
	if got := historyChordLimit(cfg); got != 7 {
		t.Errorf("history on should yield 7 chords, got %d", got)
	}

	cfg.General.MaxHistory = 0
	if got := historyChordLimit(cfg); got != 2 {
		t.Errorf("invalid low history should clamp to 2, got %d", got)
	}

	cfg.General.MaxHistory = 99
	if got := historyChordLimit(cfg); got != 12 {
		t.Errorf("invalid high history should clamp to 12, got %d", got)
	}
}

func TestIsAppExcluded(t *testing.T) {
	excluded := []string{"Code.exe", "slack", "C:\\Program Files\\Figma\\Figma.exe"}

	cases := []struct {
		exe  string
		want bool
	}{
		{"C:\\Apps\\Code.exe", true},
		{"code", true},
		{"SLACK.EXE", true},
		{"figma.exe", true},
		{"chrome.exe", false},
		{"", false},
	}
	for _, tc := range cases {
		got := isAppExcluded(context.AppContext{Exe: tc.exe}, excluded)
		if got != tc.want {
			t.Errorf("isAppExcluded(%q) = %v, want %v", tc.exe, got, tc.want)
		}
	}
	if isAppExcluded(context.AppContext{Exe: "code.exe"}, nil) {
		t.Error("an empty exclusion list must never exclude anything")
	}
}

func TestResolveDisplay(t *testing.T) {
	displays := []mygo.Display{
		{ID: 101, Label: "Dell 4K Primary", Primary: true, Bounds: mygo.Rectangle{Width: 3840, Height: 2160}},
		{ID: 102, Label: "LG Side Monitor", Primary: false, Bounds: mygo.Rectangle{X: -1920, Y: 0, Width: 1920, Height: 1080}},
	}

	// 1. Empty / primary selection picks primary display
	d, ok := resolveDisplay(displays, "")
	if !ok || d.ID != 101 {
		t.Errorf("expected primary display 101, got %v (%d)", ok, d.ID)
	}

	// 2. Select by ID
	d, ok = resolveDisplay(displays, "102")
	if !ok || d.ID != 102 || d.Bounds.X != -1920 {
		t.Errorf("expected side display 102, got %v (%d, X=%d)", ok, d.ID, d.Bounds.X)
	}

	// 3. Select by Label
	d, ok = resolveDisplay(displays, "LG Side Monitor")
	if !ok || d.ID != 102 {
		t.Errorf("expected match by label, got %d", d.ID)
	}

	// 4. Disconnected monitor gracefully falls back to primary
	d, ok = resolveDisplay(displays, "999")
	if !ok || d.ID != 101 {
		t.Errorf("expected disconnected monitor to fallback to primary 101, got %d", d.ID)
	}
}

// TestToggleShortcutMatchesRecordedTokens covers the path from the settings
// recorder to the global hook: the recorder saves plain names, while the hook
// reports glyph names for several keys.
func TestToggleShortcutMatchesRecordedTokens(t *testing.T) {
	cases := []struct {
		name    string
		tokens  []string
		keyName string
		mods    hook.Modifiers
		want    bool
	}{
		{"default F10", []string{"ShiftLeft", "F10"}, "F10", hook.ModShift, true},
		{"ctrl space", []string{"ControlLeft", "Space"}, "␣", hook.ModCtrl, true},
		{"ctrl enter", []string{"ControlLeft", "Enter"}, "↵", hook.ModCtrl, true},
		{"ctrl page up", []string{"ControlLeft", "PageUp"}, "PgUp", hook.ModCtrl, true},
		{"alt backspace", []string{"AltLeft", "Backspace"}, "⌫", hook.ModAlt, true},
		{"win left arrow", []string{"MetaLeft", "Left"}, "←", hook.ModSuper, true},
		{"wrong modifier", []string{"ControlLeft", "Space"}, "␣", hook.ModShift, false},
		{"wrong key", []string{"ControlLeft", "Space"}, "↵", hook.ModCtrl, false},
		{"key up ignored", []string{"ShiftLeft", "F10"}, "F10", hook.ModShift, false},
	}
	for _, tc := range cases {
		eventType := hook.EventKeyDown
		if tc.name == "key up ignored" {
			eventType = hook.EventKeyUp
		}
		ev := hook.InputEvent{Type: eventType, KeyName: tc.keyName, Modifiers: tc.mods}
		if got := isToggleShortcutMatch(ev, tc.tokens); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}
