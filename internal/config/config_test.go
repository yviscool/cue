package config

import (
	"encoding/json"
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
	if cfg.Version != configVersion {
		t.Errorf("expected config version %d, got %d", configVersion, cfg.Version)
	}
	if cfg.Appearance.Style != "" {
		t.Errorf("expected style to inherit from theme, got %s", cfg.Appearance.Style)
	}
	if cfg.Color.Color != "" || cfg.Color.SecondaryColor != "" {
		t.Errorf("expected keycap colors to inherit, got %s/%s", cfg.Color.Color, cfg.Color.SecondaryColor)
	}
	if cfg.Text.Size != 0 || cfg.Text.Color != "" {
		t.Errorf("expected typography to inherit, got %d/%s", cfg.Text.Size, cfg.Text.Color)
	}
	if cfg.Border.Radius != 0 || cfg.Border.Color != "" {
		t.Errorf("expected border to inherit, got %.1f/%s", cfg.Border.Radius, cfg.Border.Color)
	}
	if cfg.ThemePreset != "Liquid Glass" {
		t.Errorf("expected Liquid Glass theme, got %s", cfg.ThemePreset)
	}
	if cfg.Mouse.Size != 150 || cfg.Mouse.Color != "#009dff" {
		t.Errorf("expected mouse size 150, color #009dff, got %d/%s", cfg.Mouse.Size, cfg.Mouse.Color)
	}
}

func TestLegacyConfigMigration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_migrate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "legacy.json")
	legacy := DefaultConfig()
	legacy.Version = 0
	legacy.Appearance.Style = "lowprofile"
	legacy.Color.Color = "#ffffff"
	legacy.Color.SecondaryColor = "#1a1a1a"
	legacy.Text.Size = 32
	legacy.Text.Color = "#000000"
	legacy.Border.Radius = 0.5
	legacy.Border.Color = "#1a1a1a"
	legacy.ThemePreset = "cupertino-pro"
	legacy.General.LingerDurationMs = 3200 // a genuine preference, kept

	data, err := json.MarshalIndent(legacy, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("failed to load legacy store: %v", err)
	}
	got := store.Get()
	if got.Version != configVersion {
		t.Errorf("expected migrated version %d, got %d", configVersion, got.Version)
	}
	if got.Appearance.Style != "" || got.Color.Color != "" || got.Color.SecondaryColor != "" ||
		got.Text.Size != 0 || got.Text.Color != "" || got.Border.Radius != 0 || got.Border.Color != "" {
		t.Errorf("legacy default overrides were not cleared: %+v", got)
	}
	if got.ThemePreset != "MacBook Scissor" {
		t.Errorf("expected legacy theme alias to canonicalize, got %s", got.ThemePreset)
	}
	if got.General.LingerDurationMs != 3200 {
		t.Errorf("genuine preference was discarded: %d", got.General.LingerDurationMs)
	}

	// Migration is persisted, not only applied in memory.
	persisted, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	var onDisk AppConfig
	if err := json.Unmarshal(persisted, &onDisk); err != nil {
		t.Fatal(err)
	}
	if onDisk.Version != configVersion {
		t.Errorf("expected migrated version persisted, got %d", onDisk.Version)
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
