package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// Note-taking and knowledge tools.

func (r *Registry) initNotes() {
	r.initObsidian()
	r.initNotion()
}

func (r *Registry) initObsidian() {
	obsidian := newProfile("obsidian", "Obsidian", context.CategoryNotes,
		"obsidian.exe",
	).source("Obsidian default hotkeys", "Obsidian 1.6 (Windows)")

	// Navigation & search
	obsidian.k(hook.ModCtrl, "O", "Quick Open").
		k(hook.ModCtrl, "P", "Command Palette").
		k(hook.ModCtrl, "F", "Find in Note").
		k(hook.ModCtrl, "H", "Replace in Note").
		k(hook.ModCtrl|hook.ModShift, "F", "Search in Vault").
		k(hook.ModCtrl, "G", "Open Graph View").
		k(hook.ModCtrl|hook.ModShift, "G", "Open Local Graph").
		k(hook.ModCtrl|hook.ModAlt, "LEFT", "Navigate Back").
		k(hook.ModCtrl|hook.ModAlt, "RIGHT", "Navigate Forward")

	// Notes & tabs
	obsidian.k(hook.ModCtrl, "N", "New Note").
		k(hook.ModCtrl, "W", "Close Tab").
		k(hook.ModCtrl|hook.ModShift, "W", "Close Window").
		k(hook.ModCtrl|hook.ModShift, "T", "Reopen Closed Tab").
		k(hook.ModNone, "F2", "Rename File").
		k(hook.ModCtrl, ",", "Open Settings")

	// Formatting
	obsidian.k(hook.ModCtrl, "E", "Toggle Reading View").
		k(hook.ModCtrl, "B", "Toggle Bold").
		k(hook.ModCtrl, "I", "Toggle Italic").
		k(hook.ModCtrl, "K", "Insert Link").
		k(hook.ModCtrl|hook.ModShift, "V", "Paste and Match Style").
		k(hook.ModCtrl, "ENTER", "Toggle Checkbox").
		k(hook.ModAlt, "ENTER", "Follow Link")

	// Indentation
	obsidian.k(hook.ModNone, "TAB", "Indent").
		k(hook.ModShift, "TAB", "Outdent")

	r.addProfile(obsidian)
}

func (r *Registry) initNotion() {
	notion := newProfile("notion", "Notion", context.CategoryNotes,
		"notion.exe",
	).source("Notion keyboard shortcuts", "Notion desktop (Windows)")

	// Navigation
	notion.k(hook.ModCtrl, "P", "Quick Find").
		k(hook.ModCtrl, "N", "New Page").
		k(hook.ModCtrl|hook.ModShift, "N", "New Window").
		k(hook.ModCtrl, "[", "Navigate Back").
		k(hook.ModCtrl, "]", "Navigate Forward").
		k(hook.ModCtrl, "\\", "Toggle Sidebar").
		k(hook.ModCtrl|hook.ModShift, "L", "Toggle Dark Mode")

	// Blocks
	notion.k(hook.ModCtrl, "D", "Duplicate Block").
		k(hook.ModNone, "ENTER", "New Block").
		k(hook.ModShift, "ENTER", "New Line").
		k(hook.ModNone, "TAB", "Indent Block").
		k(hook.ModShift, "TAB", "Outdent Block")

	// Formatting
	notion.k(hook.ModCtrl, "B", "Bold").
		k(hook.ModCtrl, "I", "Italic").
		k(hook.ModCtrl, "U", "Underline").
		k(hook.ModCtrl|hook.ModShift, "S", "Strikethrough").
		k(hook.ModCtrl, "E", "Inline Code").
		k(hook.ModCtrl, "K", "Insert Link").
		k(hook.ModCtrl|hook.ModShift, "M", "Add Comment").
		k(hook.ModCtrl, "A", "Select All")

	r.addProfile(notion)
}
