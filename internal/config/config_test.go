package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigCompleteness(t *testing.T) {
	cfg := DefaultConfig()

	// Verify exact Cue defaults
	if cfg.General.Filter != "modifiers" {
		t.Errorf("expected filter modifiers, got %s", cfg.General.Filter)
	}
	if cfg.General.ShowEventHistory != false {
		t.Errorf("expected showEventHistory false, got %v", cfg.General.ShowEventHistory)
	}
	if cfg.General.MaxHistory != 5 {
		t.Errorf("expected maxHistory 5, got %d", cfg.General.MaxHistory)
	}
	if cfg.General.LingerDurationMs != 5000 {
		t.Errorf("expected lingerDurationMs 5000, got %d", cfg.General.LingerDurationMs)
	}
	if cfg.Appearance.Alignment != AlignBottomCenter {
		t.Errorf("expected alignment bottom-center, got %s", cfg.Appearance.Alignment)
	}
	if cfg.Appearance.MarginX != 100 || cfg.Appearance.MarginY != 100 {
		t.Errorf("expected margins 100, got (%d, %d)", cfg.Appearance.MarginX, cfg.Appearance.MarginY)
	}
	if cfg.Appearance.Animation != "fade" {
		t.Errorf("expected animation fade, got %s", cfg.Appearance.Animation)
	}
	if cfg.Appearance.Style != "lowprofile" {
		t.Errorf("expected style lowprofile, got %s", cfg.Appearance.Style)
	}
	if cfg.Color.Color != "#ffffff" || cfg.Color.SecondaryColor != "#1a1a1a" {
		t.Errorf("expected color #ffffff/#1a1a1a, got %s/%s", cfg.Color.Color, cfg.Color.SecondaryColor)
	}
	if cfg.Mouse.Size != 150 || cfg.Mouse.Color != "#009dff" {
		t.Errorf("expected mouse size 150, color #009dff, got %d/%s", cfg.Mouse.Size, cfg.Mouse.Color)
	}
}

func TestStorePersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test_config.json")
	store, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Mutate config
	err = store.Update(func(c *AppConfig) {
		c.General.MaxHistory = 8
		c.Appearance.Alignment = AlignTopRight
		c.Color.Color = "#00B1D2"
	})
	if err != nil {
		t.Fatalf("failed to update config: %v", err)
	}

	// Reload from new store instance
	reloaded, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("failed to reload store: %v", err)
	}
	got := reloaded.Get()
	if got.General.MaxHistory != 8 || got.Appearance.Alignment != AlignTopRight || got.Color.Color != "#00B1D2" {
		t.Errorf("reloaded config mismatch: %+v", got)
	}
}
