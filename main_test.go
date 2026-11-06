package main

import (
	"testing"

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
