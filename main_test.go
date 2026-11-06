package main

import (
	"testing"

	"github.com/egoist/mygo"
	"github.com/yviscool/cue/internal/config"
)

func TestHistoryChordLimit(t *testing.T) {
	cfg := config.DefaultConfig()

	// Default: history off means one cue at a time.
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
