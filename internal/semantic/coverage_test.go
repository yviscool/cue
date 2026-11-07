package semantic_test

import (
	"testing"

	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/semantic"
)

// minimumMappings is the "richness" floor for each profile. A profile that
// drops below its floor fails the build instead of quietly getting thinner,
// which is what happened to Sublime Text before this file existed.
//
// Raise a number when you add coverage; never lower one to make a test pass —
// lower it only with a deliberate decision recorded in the commit message.
var minimumMappings = map[string]int{
	// Editors & IDEs
	"vscode":       55,
	"jetbrains":    55,
	"eclipse":      30,
	"visualstudio": 30,
	"sublime":      38,
	"devcpp":       25,
	"idle":         20,
	"notepadpp":    25,
	"zed":          20,

	// Office
	"word":       35,
	"excel":      40,
	"powerpoint": 25,
	"outlook":    20,

	// Design
	"figma":       30,
	"photoshop":   32,
	"illustrator": 22,

	// Media
	"premiere": 20,
	"vlc":      20,
	"obs":      8,

	// Chat & meetings
	"slack":   20,
	"teams":   20,
	"discord": 15,
	"wechat":  5,
	"zoom":    12,

	// Notes
	"obsidian": 20,
	"notion":   18,

	// Shell
	"explorer": 28,
}

func TestProfileCoverageFloors(t *testing.T) {
	r := semantic.NewRegistry()
	for id, floor := range minimumMappings {
		p := r.ProfileByID(id)
		if p == nil {
			t.Errorf("profile %q is missing from the registry", id)
			continue
		}
		if got := p.ShortcutCount(); got < floor {
			t.Errorf("profile %q has %d mappings, floor is %d", id, got, floor)
		}
	}
}

// TestRegistryCoverageTargets enforces the headline numbers. They are the
// "is this actually rich enough to be useful" contract for the whole dataset.
func TestRegistryCoverageTargets(t *testing.T) {
	s := semantic.NewRegistry().Stats()

	targets := []struct {
		name string
		got  int
		min  int
	}{
		{"universal", s.Universal, 90},
		{"categories", s.Categories, 120},
		{"profile mappings", s.ProfileMap, 650},
		{"total mappings", s.Total, 850},
		{"distinct action titles", s.Titles, 450},
	}
	for _, target := range targets {
		if target.got < target.min {
			t.Errorf("%s: %d, want at least %d", target.name, target.got, target.min)
		}
	}
	if s.Duplicates != 0 {
		t.Errorf("registry has %d duplicate chords, want 0", s.Duplicates)
	}
}

// TestProfilesHaveProvenance keeps the "no blog-post data" rule enforceable:
// every profile must name the keymap it came from and the build it was checked
// against.
func TestProfilesHaveProvenance(t *testing.T) {
	for _, p := range semantic.NewRegistry().Profiles() {
		if p.Source == "" {
			t.Errorf("profile %q has no Source", p.ID)
		}
		if p.VerifiedAt == "" {
			t.Errorf("profile %q has no VerifiedAt", p.ID)
		}
	}
}

func TestRegistryHasNoDuplicates(t *testing.T) {
	for _, d := range semantic.NewRegistry().Audit() {
		t.Errorf("duplicate chord: %s", d)
	}
}

// TestShiftedSymbolEquivalence pins the normalisation layer: the same physical
// key must resolve no matter whether the profile or the user wrote the shifted
// symbol or the base key plus Shift.
func TestShiftedSymbolEquivalence(t *testing.T) {
	pairs := [][2]string{
		{"Ctrl+!", "Ctrl+Shift+1"},
		{"Ctrl+@", "Ctrl+Shift+2"},
		{"Ctrl+#", "Ctrl+Shift+3"},
		{"Ctrl+$", "Ctrl+Shift+4"},
		{"Ctrl+%", "Ctrl+Shift+5"},
		{"Ctrl+|", "Ctrl+Shift+\\"},
		{"Ctrl+:", "Ctrl+Shift+;"},
		{"Ctrl+?", "Ctrl+Shift+/"},
		{"Ctrl+>", "Ctrl+Shift+."},
		{"Ctrl+<", "Ctrl+Shift+,"},
	}
	for _, pair := range pairs {
		modsA, keyA, errA := semantic.ParseChord(pair[0])
		modsB, keyB, errB := semantic.ParseChord(pair[1])
		if errA != nil || errB != nil {
			t.Fatalf("parse error for %q / %q: %v / %v", pair[0], pair[1], errA, errB)
		}
		if modsA != modsB || keyA != keyB {
			t.Errorf("%q parsed as (%v,%q), %q parsed as (%v,%q); want identical",
				pair[0], modsA, keyA, pair[1], modsB, keyB)
		}
	}
}

// TestShiftedSymbolMatching verifies the end-to-end path: a hook event that
// reports the shifted symbol must hit the mapping registered with the base key.
func TestShiftedSymbolMatching(t *testing.T) {
	m := semantic.NewMatcher()

	excel := app("EXCEL.EXE", "Book1 - Excel")
	cases := []struct {
		mods  hook.Modifiers
		key   string
		title string
	}{
		{hook.ModCtrl | hook.ModShift, "$", "Currency Format"},
		{hook.ModCtrl | hook.ModShift, "4", "Currency Format"},
		{hook.ModCtrl | hook.ModShift, "%", "Percent Format"},
		{hook.ModCtrl | hook.ModShift, "5", "Percent Format"},
		{hook.ModCtrl | hook.ModShift, ":", "Insert Time"},
		{hook.ModCtrl | hook.ModShift, ";", "Insert Time"},
	}
	for _, tc := range cases {
		res := m.MatchAction(excel, tc.mods, tc.key)
		if res.Action.Title != tc.title {
			t.Errorf("Excel %s+%s: got %q, want %q", tc.mods, tc.key, res.Action.Title, tc.title)
		}
	}

	// A shifted pipe must reach the VS Code bracket jump registered as "\".
	code := app("code.exe", "main.go - Visual Studio Code")
	if res := m.MatchAction(code, hook.ModCtrl|hook.ModShift, "|"); res.Action.Title != "Jump to Matching Bracket" {
		t.Errorf("VS Code Ctrl+Shift+|: got %q, want Jump to Matching Bracket", res.Action.Title)
	}
	if res := m.MatchAction(code, hook.ModCtrl|hook.ModShift, "\\"); res.Action.Title != "Jump to Matching Bracket" {
		t.Errorf("VS Code Ctrl+Shift+\\: got %q, want Jump to Matching Bracket", res.Action.Title)
	}
}

// TestUserOverrideShiftedSymbol checks that a customShortcuts entry written as
// "Ctrl+Shift+1" also fires when the hook reports "!".
func TestUserOverrideShiftedSymbol(t *testing.T) {
	m := semantic.NewMatcher()
	override, err := semantic.ParseOverride("", "Ctrl+Shift+1", "My Macro", "")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	m.SetOverrides([]semantic.Override{override})

	res := m.MatchAction(context.AppContext{}, hook.ModCtrl|hook.ModShift, "!")
	if res.Action.Title != "My Macro" {
		t.Fatalf("shifted symbol did not match override: got %q", res.Action.Title)
	}
	if res.Confidence != semantic.ConfidenceUser {
		t.Fatalf("expected user confidence, got %s", res.Confidence)
	}
}

// TestNoAppProfileLeaksIntoUniversal guards the "never teach the user something
// wrong" rule: app-only chords such as PowerPoint's B must stay unresolvable
// without context.
func TestNoAppProfileLeaksIntoUniversal(t *testing.T) {
	m := semantic.NewMatcher()
	for _, key := range []string{"B", "W", "F", "S", "V", "A", "I", "O", "X", "M"} {
		if res := m.MatchUniversal(hook.ModNone, key); res.Matched() {
			t.Errorf("bare %q resolved to %q without application context", key, res.Action.Title)
		}
	}
}
