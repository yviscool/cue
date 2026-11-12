package hud

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/theme"
)

// TestShortcutTokens guards the global-hotkey recorder. The token names must
// match what main.isToggleShortcutMatch parses, otherwise recording a hotkey
// would silently save something that never fires.
func TestShortcutTokens(t *testing.T) {
	cases := []struct {
		name string
		mods ui.Modifiers
		key  ui.Key
		want []string
	}{
		{"default", ui.Shift, ui.KeyF10, []string{"ShiftLeft", "F10"}},
		{"ctrl letter", ui.Ctrl, ui.KeyK, []string{"ControlLeft", "K"}},
		{"all modifiers", ui.Ctrl | ui.Shift | ui.Alt | ui.Super, ui.KeySpace, []string{"ControlLeft", "ShiftLeft", "AltLeft", "MetaLeft", "Space"}},
		{"digit", 0, ui.Key4, []string{"4"}},
		{"empty for modifier only", ui.Ctrl, ui.KeyUnknown, nil},
	}
	for _, tc := range cases {
		got := shortcutTokens(tc.mods, tc.key)
		if len(got) != len(tc.want) {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
				break
			}
		}
	}
}

func TestKeyTokenSpecials(t *testing.T) {
	cases := map[ui.Key]string{
		ui.KeyEnter:   "Enter",
		ui.KeyEscape:  "Esc",
		ui.KeySpace:   "Space",
		ui.KeyPageUp:  "PageUp",
		ui.KeyF12:     "F12",
		ui.KeyA:       "A",
		ui.Key0:       "0",
		ui.KeyUnknown: "",
	}
	for key, want := range cases {
		if got := keyToken(key); got != want {
			t.Errorf("keyToken(%v) = %q, want %q", key, got, want)
		}
	}
}

func TestColorHex(t *testing.T) {
	if got := colorHex(ui.RGB(255, 0, 16)); got != "#FF0010" {
		t.Errorf("opaque color = %q", got)
	}
	if got := colorHex(ui.RGBA(0, 0, 0, 0.5)); got != "#0000007F" && got != "#00000080" {
		t.Errorf("translucent color = %q", got)
	}
}

func TestSplitList(t *testing.T) {
	got := splitList(" Code.exe, slack.exe;\ncode.exe\n\nfigma ")
	want := []string{"Code.exe", "slack.exe", "figma"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v", got, want)
			break
		}
	}
}

// TestThemeMenuSelectsTheme proves the menu is not decoration: opening it and
// choosing an entry applies the theme and releases the previous overrides,
// which is the contract the old shelf carried.
func TestThemeMenuSelectsTheme(t *testing.T) {
	dir := t.TempDir()
	store, err := config.NewStore(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Update(func(ac *config.AppConfig) {
		ac.Appearance.Style = "pbt"
		ac.Color.Color = "#123456"
		ac.Text.Color = "#abcdef"
	}); err != nil {
		t.Fatal(err)
	}

	sc := NewSettingsController(store, nil, nil)
	tt := ui.NewTester(sc.render, 960, 900)
	tt.SetDark(true)
	tt.Frame()

	// The trigger wears the current theme, so its label is the handle to open.
	if err := tt.Click(themeShortName(store.Get().ThemePreset)); err != nil {
		t.Fatalf("open theme menu: %v", err)
	}
	tt.Frame()

	// The menu opens scrolled to the selection and is shorter than the list, so
	// only its neighbours are on screen. Pick the entry directly above it rather
	// than a hard-coded name that may sit outside the viewport.
	target := theme.AllThemes[themeIndex(store.Get().ThemePreset)-1].Name
	if err := tt.Click(themeShortName(target)); err != nil {
		t.Fatalf("choose theme: %v", err)
	}
	tt.Frame()

	cfg := store.Get()
	if cfg.ThemePreset != target {
		t.Errorf("theme preset = %q, want %q", cfg.ThemePreset, target)
	}
	if cfg.Appearance.Style != "" {
		t.Errorf("style = %q, want auto", cfg.Appearance.Style)
	}
	if cfg.Color.Color != "" || cfg.Text.Color != "" {
		t.Errorf("color overrides survived: %q / %q", cfg.Color.Color, cfg.Text.Color)
	}
}

// TestModelMenuSelectsStyle covers the physical-model menu end to end: Auto is
// stored as the empty string, a build as its own id.
func TestModelMenuSelectsStyle(t *testing.T) {
	dir := t.TempDir()
	store, err := config.NewStore(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	sc := NewSettingsController(store, nil, nil)
	tt := ui.NewTester(sc.render, 960, 900)
	tt.SetDark(true)
	tt.Frame()

	if got := store.Get().Appearance.Style; got != "" {
		t.Fatalf("default style = %q, want auto", got)
	}

	if err := tt.Click(physicalStyleShortName("")); err != nil {
		t.Fatalf("open model menu: %v", err)
	}
	tt.Frame()
	if err := tt.Click(physicalStyleShortName("pbt")); err != nil {
		t.Fatalf("choose pbt: %v", err)
	}
	tt.Frame()
	if got := store.Get().Appearance.Style; got != "pbt" {
		t.Errorf("style = %q, want pbt", got)
	}

	// And back to Auto, which must clear the override rather than store "auto".
	if err := tt.Click(physicalStyleShortName("pbt")); err != nil {
		t.Fatalf("reopen model menu: %v", err)
	}
	tt.Frame()
	if err := tt.Click(physicalStyleShortName("")); err != nil {
		t.Fatalf("choose auto: %v", err)
	}
	tt.Frame()
	if got := store.Get().Appearance.Style; got != "" {
		t.Errorf("style = %q, want auto after choosing Auto", got)
	}
}

func TestModelStyleID(t *testing.T) {
	if got := modelStyleID("auto"); got != "" {
		t.Errorf("modelStyleID(auto) = %q, want empty", got)
	}
	for _, id := range modelChoices[1:] {
		if got := modelStyleID(id); got != id {
			t.Errorf("modelStyleID(%q) = %q", id, got)
		}
	}
	if modelChoices[0] != "auto" {
		t.Errorf("first entry = %q, want auto", modelChoices[0])
	}
}

// TestAuditionButtonFeedback covers the whole audition contract: the click
// plays the chosen profile and lights the pill, and the pill settles back on
// its own, driven by the animation rather than by the next mouse event.
func TestAuditionButtonFeedback(t *testing.T) {
	dir := t.TempDir()
	store, err := config.NewStore(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	sc := NewSettingsController(store, nil, nil)
	sc.activeTab = "Sound"

	played := ""
	sc.SetAudioPreview(func(profile string) { played = profile })

	tt := ui.NewTester(sc.render, 960, 900)
	tt.SetDark(true)
	tt.Frame()

	if err := tt.Click(i18n.T("ui_btn_preview")); err != nil {
		t.Fatalf("click audition: %v", err)
	}
	tt.Frame()

	if played == "" {
		t.Error("audition did not reach the player")
	}
	if sc.audioPulse != 1 {
		t.Errorf("pulse target = %v, want 1", sc.audioPulse)
	}
	if sc.audioPulseAt.IsZero() {
		t.Error("pulse clock was not set, the meter would freeze")
	}

	tt.Advance(700 * time.Millisecond)
	tt.Frame()
	if sc.audioPulse != 0 {
		t.Errorf("pulse did not settle: %v", sc.audioPulse)
	}
}

// TestPaletteSelection covers both halves of the palette strip: a preset writes
// its three colors, and the custom entry materialises the theme colors instead
// of leaving the wells empty.
func TestPaletteSelection(t *testing.T) {
	dir := t.TempDir()
	store, err := config.NewStore(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	sc := NewSettingsController(store, nil, nil)

	blue := MasterColorSchemes[5]
	sc.selectPalette(blue, false)
	cfg := store.Get()
	if cfg.Color.Color != blue.Primary || cfg.Color.SecondaryColor != blue.Secondary || cfg.Text.Color != blue.Text {
		t.Errorf("blue palette = %q / %q / %q", cfg.Color.Color, cfg.Color.SecondaryColor, cfg.Text.Color)
	}

	sc.materializePalette()
	cfg = store.Get()
	if !paletteIsCustom(cfg) {
		t.Errorf("custom palette not active: %q", cfg.Color.Color)
	}
	if cfg.Color.Color == "" || cfg.Color.SecondaryColor == "" || cfg.Text.Color == "" {
		t.Errorf("custom palette left wells empty: %q / %q / %q", cfg.Color.Color, cfg.Color.SecondaryColor, cfg.Text.Color)
	}

	sc.selectPalette(ColorSchemePreset{}, true)
	cfg = store.Get()
	if cfg.Color.Color != "" || cfg.Text.Color != "" || cfg.Border.Color != "" {
		t.Errorf("inherit did not clear overrides: %q / %q / %q", cfg.Color.Color, cfg.Text.Color, cfg.Border.Color)
	}
}

func TestPaletteIsCustom(t *testing.T) {
	cases := []struct {
		name string
		cfg  config.AppConfig
		want bool
	}{
		{"inherit", config.AppConfig{}, false},
		{"preset", config.AppConfig{Color: config.ColorSettings{Color: "#2196f3"}}, false},
		{"custom", config.AppConfig{Color: config.ColorSettings{Color: "#123456"}}, true},
	}
	for _, tc := range cases {
		if got := paletteIsCustom(tc.cfg); got != tc.want {
			t.Errorf("%s: paletteIsCustom = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestShortcutDraftIssues covers the save-time validation behind the recorder:
// an unreadable chord or a duplicated chord in the same app blocks saving,
// while overriding a built-in shortcut does not.
func TestShortcutDraftIssues(t *testing.T) {
	cases := []struct {
		name   string
		drafts []config.CustomShortcut
		want   int
	}{
		{
			name:   "incomplete rows are not errors",
			drafts: []config.CustomShortcut{{Label: "Copy"}, {Chord: "Ctrl+K"}},
			want:   0,
		},
		{
			name:   "same chord in different apps is fine",
			drafts: []config.CustomShortcut{{App: "code.exe", Chord: "Ctrl+K", Label: "a"}, {App: "figma.exe", Chord: "Ctrl+K", Label: "b"}},
			want:   0,
		},
		{
			name:   "duplicate chord in one app",
			drafts: []config.CustomShortcut{{App: "code.exe", Chord: "Ctrl+K", Label: "a"}, {App: "Code.exe", Chord: "ctrl + k", Label: "b"}},
			want:   1,
		},
		{
			name:   "unreadable chord",
			drafts: []config.CustomShortcut{{App: "code.exe", Chord: "Ctrl+Shift", Label: "a"}},
			want:   1,
		},
	}
	for _, tc := range cases {
		if got := len(shortcutDraftIssues(tc.drafts)); got != tc.want {
			t.Errorf("%s: %d issues, want %d (%v)", tc.name, got, tc.want, shortcutDraftIssues(tc.drafts))
		}
	}
}

// TestRecordShortcutChord pins the recorder behaviour: a real key press becomes
// a chord the matcher can read back, and Backspace clears instead of recording.
func TestRecordShortcutChord(t *testing.T) {
	draft := config.CustomShortcut{}
	if !recordShortcutChord(&draft, ui.TextKey{Mods: ui.Ctrl | ui.Shift, Key: ui.KeyK}) {
		t.Fatal("Ctrl+Shift+K was not accepted")
	}
	if draft.Chord != "Ctrl + Shift + K" {
		t.Errorf("chord = %q, want Ctrl + Shift + K", draft.Chord)
	}
	if !recordShortcutChord(&draft, ui.TextKey{Key: ui.KeyBackspace}) {
		t.Fatal("Backspace was not consumed")
	}
	if draft.Chord != "" {
		t.Errorf("chord after Backspace = %q, want empty", draft.Chord)
	}
	if recordShortcutChord(&draft, ui.TextKey{Mods: ui.Ctrl}) {
		t.Error("a modifier-only press should be ignored")
	}
}

// themeIndex returns the position of name in theme.AllThemes, or the index of the
// default theme when the name is unknown.
func themeIndex(name string) int {
	for i, th := range theme.AllThemes {
		if th.Name == name {
			return i
		}
	}
	return 1
}
