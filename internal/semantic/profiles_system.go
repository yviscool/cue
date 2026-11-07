package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// Windows shell.

func (r *Registry) initExplorer() {
	explorer := newProfile("explorer", "File Explorer", context.CategorySystem,
		"explorer.exe",
	).source("Windows 11 File Explorer shortcuts", "Windows 11 23H2")

	// Navigation
	explorer.k(hook.ModCtrl, "L", "Focus Address Bar").
		k(hook.ModAlt, "D", "Focus Address Bar").
		k(hook.ModNone, "F4", "Address Bar History").
		k(hook.ModAlt, "UP", "Up One Level").
		k(hook.ModAlt, "LEFT", "Back").
		k(hook.ModAlt, "RIGHT", "Forward").
		k(hook.ModNone, "BACKSPACE", "Back").
		k(hook.ModCtrl, "N", "New Window")

	// Files
	explorer.k(hook.ModCtrl|hook.ModShift, "N", "New Folder").
		k(hook.ModNone, "F2", "Rename").
		k(hook.ModNone, "DELETE", "Delete").
		k(hook.ModShift, "DELETE", "Permanently Delete").
		k(hook.ModCtrl|hook.ModShift, "C", "Copy Path").
		k(hook.ModCtrl|hook.ModShift, "V", "Paste Shortcut").
		k(hook.ModNone, "F5", "Refresh").
		k(hook.ModAlt, "ENTER", "Properties").
		k(hook.ModCtrl, "Z", "Undo").
		k(hook.ModCtrl, "Y", "Redo").
		k(hook.ModCtrl, "A", "Select All")

	// Search & panes
	explorer.k(hook.ModCtrl, "F", "Search").
		k(hook.ModNone, "F3", "Search").
		k(hook.ModAlt, "P", "Preview Pane").
		k(hook.ModAlt|hook.ModShift, "P", "Details Pane").
		k(hook.ModNone, "F11", "Full Screen")

	// Tree expansion
	explorer.k(hook.ModNone, "NUM*", "Expand All Folders").
		k(hook.ModNone, "NUM+", "Expand Folder").
		k(hook.ModNone, "NUM-", "Collapse Folder")

	// View modes
	explorer.k(hook.ModCtrl|hook.ModShift, "1", "Extra Large Icons").
		k(hook.ModCtrl|hook.ModShift, "2", "Large Icons").
		k(hook.ModCtrl|hook.ModShift, "3", "Medium Icons").
		k(hook.ModCtrl|hook.ModShift, "4", "Small Icons").
		k(hook.ModCtrl|hook.ModShift, "5", "List View").
		k(hook.ModCtrl|hook.ModShift, "6", "Details View").
		k(hook.ModCtrl|hook.ModShift, "8", "Content View")

	r.addProfile(explorer)
}
