package semantic_test

import (
	"testing"

	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/semantic"
)

// Table-driven coverage for the profiles that were expanded most aggressively.
// These assertions are the contract for the data itself: if a chord is retyped
// wrongly during a future edit, the build fails instead of teaching the user
// the wrong thing on stage.

type chordCase struct {
	mods hook.Modifiers
	key  string
	want string
}

type profileTable struct {
	profile string
	exe     string
	title   string
	cases   []chordCase
}

var profileTables = []profileTable{
	{
		profile: "sublime", exe: "sublime_text.exe", title: "main.py - Sublime Text",
		cases: []chordCase{
			{hook.ModCtrl, "ENTER", "Insert Line After"},
			{hook.ModCtrl | hook.ModShift, "ENTER", "Insert Line Before"},
			{hook.ModCtrl | hook.ModShift, "R", "Goto Symbol in Project"},
			{hook.ModAlt, "F3", "Select All Occurrences"},
			{hook.ModCtrl | hook.ModShift, "L", "Split Selection into Lines"},
			{hook.ModCtrl, "J", "Join Lines"},
			{hook.ModNone, "F9", "Sort Lines"},
			{hook.ModCtrl | hook.ModShift, "SPACE", "Expand Selection to Scope"},
			{hook.ModCtrl, "M", "Jump to Matching Bracket"},
			{hook.ModCtrl | hook.ModShift, "V", "Paste and Indent"},
			{hook.ModCtrl, "`", "Show Console"},
			{hook.ModShift, "F11", "Distraction Free Mode"},
		},
	},
	{
		profile: "devcpp", exe: "devcpp.exe", title: "main.cpp - Dev-C++ 6.3",
		cases: []chordCase{
			{hook.ModCtrl | hook.ModShift, "A", "Format Code"},
			{hook.ModNone, "F9", "Compile"},
			{hook.ModCtrl, "F9", "Compile"},
			{hook.ModCtrl, "F10", "Run"},
			{hook.ModCtrl, "I", "Indent"},
			{hook.ModCtrl | hook.ModShift, "I", "Unindent"},
			{hook.ModCtrl | hook.ModShift, "S", "Save All"},
			{hook.ModCtrl | hook.ModAlt, "F2", "Terminate Program"},
			{hook.ModCtrl, "BACKSPACE", "Delete Word Left"},
			{hook.ModCtrl, "G", "Go to Line"},
		},
	},
	{
		profile: "idle", exe: "python.exe", title: "Untitled - IDLE (Python 3.12)",
		cases: []chordCase{
			{hook.ModShift, "F5", "Run Customized"},
			{hook.ModAlt, "X", "Check Module"},
			{hook.ModCtrl | hook.ModShift, "S", "Save As"},
			{hook.ModCtrl, "Q", "Quit"},
			{hook.ModCtrl, "G", "Find Next"},
			{hook.ModCtrl | hook.ModShift, "G", "Find Previous"},
			{hook.ModAlt, "P", "Previous Command"},
			{hook.ModAlt, "N", "Next Command"},
			{hook.ModNone, "TAB", "Indent"},
			{hook.ModShift, "TAB", "Dedent"},
		},
	},
	{
		profile: "notepadpp", exe: "notepad++.exe", title: "config.ini - Notepad++",
		cases: []chordCase{
			{hook.ModCtrl, "D", "Duplicate Line"},
			{hook.ModCtrl, "L", "Delete Line"},
			{hook.ModCtrl, "T", "Transpose Line"},
			{hook.ModCtrl | hook.ModShift, "T", "Copy Current Line"},
			{hook.ModCtrl, "U", "Lower Case"},
			{hook.ModCtrl | hook.ModShift, "U", "Upper Case"},
			{hook.ModCtrl, "Q", "Toggle Block Comment"},
			{hook.ModAlt, "0", "Fold All"},
			{hook.ModCtrl | hook.ModAlt, "F", "Fold Current Level"},
		},
	},
	{
		profile: "excel", exe: "EXCEL.EXE", title: "Book1 - Excel",
		cases: []chordCase{
			{hook.ModCtrl | hook.ModShift, "4", "Currency Format"},
			{hook.ModCtrl | hook.ModShift, "5", "Percent Format"},
			{hook.ModCtrl | hook.ModShift, "3", "Date Format"},
			{hook.ModCtrl, "D", "Fill Down"},
			{hook.ModCtrl | hook.ModShift, "DOWN", "Extend Selection Down"},
			{hook.ModCtrl, "9", "Hide Row"},
			{hook.ModCtrl | hook.ModShift, "9", "Unhide Row"},
			{hook.ModCtrl, "`", "Show Formulas"},
			{hook.ModNone, "F4", "Repeat Last Action"},
			{hook.ModAlt, "=", "AutoSum"},
			{hook.ModCtrl, "UP", "Jump to Edge Up"},
			{hook.ModShift, "F2", "Cell Comment"},
		},
	},
	{
		profile: "powerpoint", exe: "POWERPNT.EXE", title: "Deck - PowerPoint",
		cases: []chordCase{
			{hook.ModAlt, "F5", "Presenter View"},
			{hook.ModCtrl | hook.ModShift, "D", "Duplicate Slide"},
			{hook.ModCtrl, "END", "Last Slide"},
			{hook.ModAlt | hook.ModShift, "UP", "Move Slide Up"},
			{hook.ModCtrl, "ENTER", "Next Placeholder"},
			{hook.ModCtrl, "T", "Font Dialog"},
		},
	},
	{
		profile: "outlook", exe: "OUTLOOK.EXE", title: "Inbox - Outlook",
		cases: []chordCase{
			{hook.ModCtrl, "2", "Go to Calendar"},
			{hook.ModCtrl | hook.ModShift, "V", "Move to Folder"},
			{hook.ModCtrl, "Q", "Mark as Read"},
			{hook.ModAlt, "S", "Send Message"},
			{hook.ModCtrl | hook.ModShift, "K", "New Task"},
			{hook.ModCtrl, "6", "Folder List"},
		},
	},
	{
		profile: "figma", exe: "figma.exe", title: "Untitled - Figma",
		cases: []chordCase{
			{hook.ModCtrl, "/", "Search Menu"},
			{hook.ModShift, "1", "Zoom to Fit"},
			{hook.ModCtrl | hook.ModAlt, "U", "Union Selection"},
			{hook.ModCtrl | hook.ModShift, "E", "Export"},
			{hook.ModNone, "TAB", "Hide UI"},
			{hook.ModShift, "A", "Auto Layout"},
			{hook.ModCtrl | hook.ModShift, "[", "Send to Back"},
		},
	},
	{
		profile: "photoshop", exe: "photoshop.exe", title: "Untitled-1 - Photoshop",
		cases: []chordCase{
			{hook.ModCtrl | hook.ModAlt | hook.ModShift, "E", "Stamp Visible"},
			{hook.ModCtrl, "L", "Levels"},
			{hook.ModAlt, "]", "Select Layer Above"},
			{hook.ModCtrl | hook.ModShift, "I", "Inverse Selection"},
			{hook.ModCtrl | hook.ModAlt, "G", "Create Clipping Mask"},
			{hook.ModCtrl | hook.ModAlt, "Z", "Step Backward"},
		},
	},
	{
		profile: "illustrator", exe: "illustrator.exe", title: "Untitled-1 - Illustrator",
		cases: []chordCase{
			{hook.ModCtrl, "7", "Clipping Mask"},
			{hook.ModCtrl | hook.ModAlt, "7", "Release Clipping Mask"},
			{hook.ModCtrl, "8", "Compound Path"},
			{hook.ModCtrl, "F", "Paste in Front"},
			{hook.ModCtrl, "B", "Paste in Back"},
			{hook.ModCtrl, "2", "Lock Selection"},
			{hook.ModCtrl, "Y", "Outline Mode"},
		},
	},
	{
		profile: "premiere", exe: "adobe premiere pro.exe", title: "Project - Premiere Pro",
		cases: []chordCase{
			{hook.ModCtrl, "K", "Add Edit"},
			{hook.ModCtrl, "D", "Apply Video Transition"},
			{hook.ModNone, "I", "Mark In"},
			{hook.ModNone, "M", "Add Marker"},
			{hook.ModCtrl | hook.ModShift, "E", "Export Frame"},
			{hook.ModNone, "L", "Play Forward"},
		},
	},
	{
		profile: "vlc", exe: "vlc.exe", title: "VLC media player",
		cases: []chordCase{
			{hook.ModNone, "V", "Cycle Subtitles"},
			{hook.ModNone, "B", "Cycle Audio Track"},
			{hook.ModNone, "M", "Mute"},
			{hook.ModCtrl, "T", "Go to Time"},
			{hook.ModNone, "F", "Full Screen"},
		},
	},
	{
		profile: "obs", exe: "obs64.exe", title: "OBS 30.0 - Profile: Untitled",
		cases: []chordCase{
			{hook.ModCtrl, "D", "Duplicate"},
			{hook.ModCtrl, "S", "Save Scene Collection"},
			{hook.ModNone, "F2", "Rename"},
			{hook.ModNone, "DELETE", "Remove"},
		},
	},
	{
		profile: "slack", exe: "slack.exe", title: "Slack - team",
		cases: []chordCase{
			{hook.ModCtrl | hook.ModShift, "A", "All Unreads"},
			{hook.ModCtrl | hook.ModShift, "L", "Browse Channels"},
			{hook.ModCtrl, "UP", "Edit Last Message"},
			{hook.ModCtrl | hook.ModShift, ">", "Block Quote"},
			{hook.ModShift, "ESC", "Mark All as Read"},
		},
	},
	{
		profile: "teams", exe: "ms-teams.exe", title: "Microsoft Teams",
		cases: []chordCase{
			{hook.ModCtrl | hook.ModShift, "M", "Toggle Mute"},
			{hook.ModCtrl | hook.ModShift, "E", "Start Screen Share"},
			{hook.ModCtrl, "3", "Go to Teams"},
			{hook.ModCtrl | hook.ModShift, "H", "Hang Up"},
		},
	},
	{
		profile: "discord", exe: "discord.exe", title: "Discord",
		cases: []chordCase{
			{hook.ModCtrl | hook.ModShift, "D", "Toggle Deafen"},
			{hook.ModAlt | hook.ModShift, "UP", "Previous Server"},
			{hook.ModShift, "ESC", "Mark Server Read"},
			{hook.ModCtrl, "/", "Keyboard Shortcuts"},
		},
	},
	{
		profile: "zoom", exe: "zoom.exe", title: "Zoom Meeting",
		cases: []chordCase{
			{hook.ModAlt, "A", "Toggle Mute"},
			{hook.ModAlt, "V", "Toggle Video"},
			{hook.ModAlt, "H", "Raise Hand"},
			{hook.ModAlt, "S", "Start Screen Share"},
			{hook.ModAlt, "Q", "End Meeting"},
		},
	},
	{
		profile: "obsidian", exe: "obsidian.exe", title: "Note - Obsidian",
		cases: []chordCase{
			{hook.ModCtrl, "G", "Open Graph View"},
			{hook.ModCtrl | hook.ModAlt, "LEFT", "Navigate Back"},
			{hook.ModAlt, "ENTER", "Follow Link"},
			{hook.ModCtrl, "ENTER", "Toggle Checkbox"},
			{hook.ModNone, "F2", "Rename File"},
		},
	},
	{
		profile: "notion", exe: "notion.exe", title: "Page - Notion",
		cases: []chordCase{
			{hook.ModCtrl, "\\", "Toggle Sidebar"},
			{hook.ModCtrl | hook.ModShift, "L", "Toggle Dark Mode"},
			{hook.ModCtrl, "E", "Inline Code"},
			{hook.ModCtrl | hook.ModShift, "M", "Add Comment"},
			{hook.ModCtrl, "D", "Duplicate Block"},
		},
	},
	{
		profile: "explorer", exe: "explorer.exe", title: "Downloads",
		cases: []chordCase{
			{hook.ModCtrl | hook.ModShift, "C", "Copy Path"},
			{hook.ModAlt, "P", "Preview Pane"},
			{hook.ModCtrl | hook.ModShift, "6", "Details View"},
			{hook.ModNone, "NUM*", "Expand All Folders"},
			{hook.ModAlt, "ENTER", "Properties"},
		},
	},
}

func TestProfileShortcutTables(t *testing.T) {
	m := semantic.NewMatcher()
	for _, table := range profileTables {
		ctx := app(table.exe, table.title)
		t.Run(table.profile, func(t *testing.T) {
			for _, tc := range table.cases {
				res := m.MatchAction(ctx, tc.mods, tc.key)
				if res.Action.Title != tc.want {
					t.Errorf("%s + %s: got %q, want %q", tc.mods, tc.key, res.Action.Title, tc.want)
				}
				if res.ProfileID != table.profile {
					t.Errorf("%s + %s resolved through profile %q, want %q", tc.mods, tc.key, res.ProfileID, table.profile)
				}
			}
		})
	}
}

// TestSequenceProfiles covers the two-chord families that were extended.
func TestSequenceProfiles(t *testing.T) {
	m := semantic.NewMatcher()

	sublime := app("sublime_text.exe", "main.py - Sublime Text")
	sequences := []struct {
		firstMods  hook.Modifiers
		firstKey   string
		secondMods hook.Modifiers
		secondKey  string
		want       string
	}{
		{hook.ModCtrl, "K", hook.ModCtrl, "B", "Toggle Sidebar"},
		{hook.ModCtrl, "K", hook.ModCtrl, "U", "Upper Case"},
		{hook.ModCtrl, "K", hook.ModCtrl, "L", "Lower Case"},
		{hook.ModCtrl, "K", hook.ModCtrl, "K", "Delete to End of Line"},
	}
	for _, tc := range sequences {
		res, ok := m.MatchSequence(sublime, tc.firstMods, tc.firstKey, tc.secondMods, tc.secondKey)
		if !ok || res.Action.Title != tc.want {
			t.Errorf("%s+%s then %s+%s: got %q (%t), want %q",
				tc.firstMods, tc.firstKey, tc.secondMods, tc.secondKey, res.Action.Title, ok, tc.want)
		}
	}

	// Visual Studio keeps its own Ctrl+K family.
	vs := app("devenv.exe", "Solution - Microsoft Visual Studio")
	if res, ok := m.MatchSequence(vs, hook.ModCtrl, "K", hook.ModCtrl, "C"); !ok || res.Action.Title != "Comment Selection" {
		t.Errorf("VS Ctrl+K Ctrl+C: got %q (%t)", res.Action.Title, ok)
	}

	// And the same chord must not leak into an unrelated app.
	notepad := app("notepad.exe", "Untitled - Notepad")
	if _, ok := m.MatchSequence(notepad, hook.ModCtrl, "K", hook.ModCtrl, "B"); ok {
		t.Error("Ctrl+K Ctrl+B must not resolve in Notepad")
	}
}
