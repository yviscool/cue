package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfigCompleteness(t *testing.T) {
	cfg := DefaultConfig()

	// Verify exact Cue defaults
	if cfg.General.Filter != "none" {
		t.Errorf("expected filter none, got %s", cfg.General.Filter)
	}
	if cfg.General.ShowEventHistory != true {
		t.Errorf("expected showEventHistory true, got %v", cfg.General.ShowEventHistory)
	}
	if cfg.General.MaxHistory != 5 {
		t.Errorf("expected maxHistory 5, got %d", cfg.General.MaxHistory)
	}
	if cfg.General.LingerDurationMs != 4500 {
		t.Errorf("expected lingerDurationMs 4500, got %d", cfg.General.LingerDurationMs)
	}
	if cfg.Appearance.Alignment != AlignBottomRight {
		t.Errorf("expected alignment bottom-right, got %s", cfg.Appearance.Alignment)
	}
	if cfg.Appearance.MarginX != 100 || cfg.Appearance.MarginY != 100 {
		t.Errorf("expected margins 100, got (%d, %d)", cfg.Appearance.MarginX, cfg.Appearance.MarginY)
	}
	if cfg.Appearance.Animation != "spring" {
		t.Errorf("expected animation spring, got %s", cfg.Appearance.Animation)
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
	if cfg.ThemePreset != "Swiss Typography 1957" {
		t.Errorf("expected Swiss Typography theme, got %s", cfg.ThemePreset)
	}
	if cfg.Mouse.Size != 40 || cfg.Mouse.Color != "" {
		t.Errorf("expected mouse size 40 and theme-following color, got %d/%s", cfg.Mouse.Size, cfg.Mouse.Color)
	}
	if cfg.Mouse.DragThreshold != 50 {
		t.Errorf("expected mouse drag threshold 50, got %d", cfg.Mouse.DragThreshold)
	}
	if cfg.General.DragThreshold != 0 {
		t.Errorf("expected deprecated general drag threshold to stay empty, got %d", cfg.General.DragThreshold)
	}
	if cfg.Background.Enabled || cfg.Color.UseGradient {
		t.Errorf("expected backdrop and gradient off by default, got %v/%v", cfg.Background.Enabled, cfg.Color.UseGradient)
	}
	if cfg.Text.Alignment != AlignCenter || cfg.Layout.IconAlignment != "center" {
		t.Errorf("expected centered text and icon defaults, got %s/%s", cfg.Text.Alignment, cfg.Layout.IconAlignment)
	}
	if cfg.App.TrayClickAction != "settings" {
		t.Errorf("expected tray click to open settings, got %q", cfg.App.TrayClickAction)
	}
	if !cfg.SemanticContext || !cfg.ShowAppBadge {
		t.Errorf("expected context intelligence defaults on, got context=%v badge=%v", cfg.SemanticContext, cfg.ShowAppBadge)
	}
	if !cfg.SemanticSub || !cfg.TactileAudio {
		t.Errorf("expected semantic subtitles and audio on by default")
	}
	if len(cfg.CustomShortcuts) != 0 {
		t.Errorf("expected no custom shortcuts by default, got %d", len(cfg.CustomShortcuts))
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
	if !got.SemanticContext || !got.ShowAppBadge {
		t.Errorf("v3 migration must enable context intelligence by default, got %v/%v", got.SemanticContext, got.ShowAppBadge)
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

func TestV3ExplicitContextOptOutIsPreserved(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_optout")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "v3.json")
	cfg := DefaultConfig()
	cfg.SemanticContext = false
	cfg.ShowAppBadge = false
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(filePath)
	if err != nil {
		t.Fatal(err)
	}
	got := store.Get()
	if got.SemanticContext || got.ShowAppBadge {
		t.Fatalf("explicit v3 opt-out was overwritten: context=%v badge=%v", got.SemanticContext, got.ShowAppBadge)
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

func TestStoreHotReload(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_reload")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "hot_reload.json")
	store, err := NewStore(filePath)
	if err != nil {
		t.Fatal(err)
	}

	notifyCh := make(chan struct{}, 1)
	store.OnChange(func(c AppConfig) {
		if c.General.Filter == "silent" {
			select {
			case notifyCh <- struct{}{}:
			default:
			}
		}
	})

	// External write simulating manual edit of config.json
	cfg := store.Get()
	cfg.General.Filter = "silent"
	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Fatal(err)
	}

	if err := store.Reload(); err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	select {
	case <-notifyCh:
	case <-time.After(500 * time.Millisecond):
		t.Errorf("timeout waiting for OnChange callback on hot reload")
	}
	if got := store.Get().General.Filter; got != "silent" {
		t.Errorf("expected filter silent, got %s", got)
	}
}

func TestStoreInvalidJsonPreservesConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_bad")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "corrupt.json")
	store, err := NewStore(filePath)
	if err != nil {
		t.Fatal(err)
	}

	orig := store.Get()

	errCh := make(chan error, 1)
	store.OnError(func(err error) {
		select {
		case errCh <- err:
		default:
		}
	})

	// External corrupt write (syntax error)
	if err := os.WriteFile(filePath, []byte(`{ "general": { "filter": `), 0644); err != nil {
		t.Fatal(err)
	}

	err = store.Reload()
	if err == nil {
		t.Errorf("expected error on corrupt JSON")
	}
	select {
	case <-errCh:
	case <-time.After(500 * time.Millisecond):
		t.Errorf("timeout waiting for OnError callback to fire")
	}
	// Active config in memory must remain completely intact
	if store.Get().General.Filter != orig.General.Filter {
		t.Errorf("corrupt reload corrupted in-memory config: %+v", store.Get())
	}
}

// TestV4PresentationDefaults checks the migration path for the opacity and
// volume sliders. A config written by a v3 build has neither key, so both must
// land on their neutral position instead of a zero that would render the overlay
// invisible or mute the app.
func TestV4PresentationDefaults(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_v4")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "v3.json")
	// Deliberately written by hand: no opacity, no audioVolume.
	raw := `{
  "version": 3,
  "appearance": {"alignment": "bottom-center", "marginX": 120},
  "tactileAudio": true,
  "audioProfile": "thock",
  "themePreset": "Liquid Glass"
}`
	if err := os.WriteFile(filePath, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("failed to load v3 store: %v", err)
	}
	got := store.Get()

	if got.Version != configVersion {
		t.Errorf("expected version %d, got %d", configVersion, got.Version)
	}
	if got.Appearance.Opacity != DefaultOpacity {
		t.Errorf("opacity should default to %v, got %v", DefaultOpacity, got.Appearance.Opacity)
	}
	if got.AudioVolume != DefaultAudioVolume {
		t.Errorf("volume should default to %v, got %v", DefaultAudioVolume, got.AudioVolume)
	}
	if got.Appearance.EffectiveOpacity() != 1.0 {
		t.Errorf("effective opacity should be 1.0, got %v", got.Appearance.EffectiveOpacity())
	}
	if got.EffectiveVolume() != 1.0 {
		t.Errorf("effective volume should be unity gain, got %v", got.EffectiveVolume())
	}
	if got.Mouse.Opacity != DefaultOpacity {
		t.Errorf("mouse opacity should default to %v, got %v", DefaultOpacity, got.Mouse.Opacity)
	}
	// Existing preferences survive the migration untouched.
	if got.Appearance.MarginX != 120 || got.AudioProfile != "thock" {
		t.Errorf("existing preferences were lost: %+v", got.Appearance)
	}
}

// TestV4ExplicitPresentationValuesArePreserved makes sure a muted overlay or a
// lowered volume that the user actually chose is never overwritten by the
// migration defaults.
func TestV4ExplicitPresentationValuesArePreserved(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_v4_explicit")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "v4.json")
	raw := `{
  "version": 4,
  "appearance": {"opacity": 0.45},
  "audioVolume": 0.3
}`
	if err := os.WriteFile(filePath, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("failed to load v4 store: %v", err)
	}
	got := store.Get()
	if got.Appearance.Opacity != 0.45 {
		t.Errorf("explicit opacity was overwritten: %v", got.Appearance.Opacity)
	}
	if got.AudioVolume != 0.3 {
		t.Errorf("explicit volume was overwritten: %v", got.AudioVolume)
	}
	// A v4 file predates the mouse overlay slider, so that key is absent.
	if got.Mouse.Opacity != DefaultOpacity {
		t.Errorf("mouse opacity should default to %v, got %v", DefaultOpacity, got.Mouse.Opacity)
	}
}

// TestV5MouseOpacityPreserved checks the migration path for the mouse overlay
// opacity slider: a config written by a v5 build keeps the user's chosen value
// instead of being reset to fully opaque.
func TestV5MouseOpacityPreserved(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_v5_explicit")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "v5.json")
	raw := `{
  "version": 5,
  "mouse": {"opacity": 0.35, "indicatorSize": 90}
}`
	if err := os.WriteFile(filePath, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("failed to load v5 store: %v", err)
	}
	got := store.Get()
	if got.Mouse.Opacity != 0.35 {
		t.Errorf("explicit mouse opacity was overwritten: %v", got.Mouse.Opacity)
	}
	if got.Mouse.IndicatorSize != 90 {
		t.Errorf("unrelated mouse settings were lost: %v", got.Mouse.IndicatorSize)
	}
}

// TestV6SettingsCleanup checks the schema cleanup: the drag threshold moves to
// the mouse section, and fields that were previously inert get defaults that
// preserve the v5 appearance instead of suddenly changing it.
func TestV6SettingsCleanup(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cue_cfg_v6_cleanup")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "v5.json")
	raw := `{
  "version": 5,
  "general": {"dragThreshold": 80},
  "appearance": {"monitor": "12345"},
  "background": {"enabled": true, "color": "#ffffff99"},
  "color": {"useGradient": true},
  "layout": {"iconAlignment": "flex-end"},
  "mouse": {"color": "#00FF6A", "indicatorSize": 70}
}`
	if err := os.WriteFile(filePath, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("failed to load v5 store: %v", err)
	}
	got := store.Get()
	if got.Version != configVersion {
		t.Errorf("expected version %d, got %d", configVersion, got.Version)
	}
	if got.Mouse.DragThreshold != 80 {
		t.Errorf("drag threshold was not migrated: %d", got.Mouse.DragThreshold)
	}
	if got.General.DragThreshold != 0 {
		t.Errorf("deprecated drag threshold should be cleared: %d", got.General.DragThreshold)
	}
	if got.Background.Enabled {
		t.Error("the previously inert backdrop must stay off for existing installs")
	}
	if got.Color.UseGradient {
		t.Error("the previously inert gradient must stay off for existing installs")
	}
	if got.Layout.IconAlignment != "center" || got.Text.Alignment != AlignCenter {
		t.Errorf("alignment should reset to the v5 visual result, got %s/%s", got.Layout.IconAlignment, got.Text.Alignment)
	}
	if got.Mouse.Color != "" {
		t.Errorf("old neon seed should clear to follow the theme, got %q", got.Mouse.Color)
	}
	if got.Mouse.IndicatorSize != 70 || got.Appearance.Monitor != "12345" {
		t.Error("unrelated settings were lost during migration")
	}
}

func TestEffectiveMouseOpacityClamps(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{0, DefaultOpacity},    // unset or hand-edited to zero: never invisible
		{-0.5, DefaultOpacity}, // nonsense
		{0.02, MinOpacity},     // below the usable band
		{0.4, 0.4},             // in range
		{1, MaxOpacity},        // fully opaque
		{3, MaxOpacity},        // nonsense
	}
	for _, tc := range cases {
		got := MouseSettings{Opacity: tc.in}.EffectiveOpacity()
		if got != tc.want {
			t.Errorf("MouseSettings.EffectiveOpacity(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestEffectiveVolumeClamps(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{-1, 0},
		{0, 0},
		{0.35, 0.35},
		{1, 1},
		{2.5, 1},
	}
	for _, tc := range cases {
		got := AppConfig{AudioVolume: tc.in}.EffectiveVolume()
		if got != tc.want {
			t.Errorf("EffectiveVolume(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestEffectiveOpacityClamps(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{0, DefaultOpacity},    // unset or hand-edited to zero: never invisible
		{-0.5, DefaultOpacity}, // nonsense
		{0.02, MinOpacity},     // below the usable band
		{0.4, 0.4},             // in range
		{1, MaxOpacity},        // fully opaque
		{3, MaxOpacity},        // nonsense
	}
	for _, tc := range cases {
		got := AppearanceSettings{Opacity: tc.in}.EffectiveOpacity()
		if got != tc.want {
			t.Errorf("EffectiveOpacity(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

// TestDefaultConfigIsFullyOpaqueAndAudible pins the shipped defaults so a future
// edit cannot silently dim the overlay or mute the keystroke sound.
func TestDefaultConfigIsFullyOpaqueAndAudible(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Appearance.EffectiveOpacity() != 1.0 {
		t.Errorf("default overlay opacity should be 1.0, got %v", cfg.Appearance.EffectiveOpacity())
	}
	if cfg.EffectiveVolume() != 1.0 {
		t.Errorf("default volume should be unity gain, got %v", cfg.EffectiveVolume())
	}
	if cfg.Mouse.EffectiveOpacity() != 1.0 {
		t.Errorf("default mouse opacity should be 1.0, got %v", cfg.Mouse.EffectiveOpacity())
	}
}

// TestDecodeKeepsDefaultsForOmittedFields covers clipboard import: a small
// snippet must not zero every setting it does not mention.
func TestDecodeKeepsDefaultsForOmittedFields(t *testing.T) {
	got, err := Decode([]byte(`{"themePreset":"Neon Night","general":{"showEventHistory":true}}`))
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	def := DefaultConfig()
	if got.ThemePreset != "Neon Night" || !got.General.ShowEventHistory {
		t.Errorf("provided fields were not applied: %+v", got)
	}
	if got.Mouse.Size != def.Mouse.Size || got.Mouse.DragThreshold != def.Mouse.DragThreshold {
		t.Errorf("omitted mouse fields should keep defaults, got %+v", got.Mouse)
	}
	if got.Version != def.Version {
		t.Errorf("expected current schema version %d, got %d", def.Version, got.Version)
	}
}

func TestDecodeRejectsInvalidJSON(t *testing.T) {
	if _, err := Decode([]byte(`{"themePreset":`)); err == nil {
		t.Fatal("expected malformed JSON to be rejected")
	}
}

// TestEncodeDecodeRoundTrip guards the export/import pair used by the settings
// window.
func TestEncodeDecodeRoundTrip(t *testing.T) {
	orig := DefaultConfig()
	orig.ThemePreset = "Dieter Rams 1960"
	orig.App.ExcludedApps = []string{"code.exe"}
	orig.CustomShortcuts = []CustomShortcut{{Chord: "Ctrl+K", Label: "Command Palette"}}

	data, err := Encode(orig)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if got.ThemePreset != orig.ThemePreset {
		t.Errorf("theme lost: %q", got.ThemePreset)
	}
	if len(got.App.ExcludedApps) != 1 || got.App.ExcludedApps[0] != "code.exe" {
		t.Errorf("exclusions lost: %v", got.App.ExcludedApps)
	}
	if len(got.CustomShortcuts) != 1 || got.CustomShortcuts[0].Chord != "Ctrl+K" {
		t.Errorf("shortcuts lost: %v", got.CustomShortcuts)
	}
}
