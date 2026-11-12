package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// File managers.

// initTotalCommander transcribes Total Commander's own Keyboard.txt, which ships
// in the installation directory next to the executable. Every chord below is
// taken from that file rather than from memory: in a keyboard-driven manager a
// wrong label is worse than no label, because F5/F6/F7/F8 mean copy, move, mkdir
// and delete there, not refresh, address bar and so on.
func (r *Registry) initTotalCommander() {
	tc := newProfile("totalcommander", "Total Commander", context.CategorySystem,
		"totalcmd.exe", "totalcmd64.exe",
	).source("Total Commander Keyboard.txt (installation directory)", "Total Commander 11")

	// Function keys are the core file operations.
	tc.k(hook.ModNone, "F1", "Help").
		k(hook.ModNone, "F2", "Reread Source").
		k(hook.ModNone, "F3", "View File").
		k(hook.ModNone, "F4", "Edit File").
		k(hook.ModNone, "F5", "Copy Files").
		k(hook.ModNone, "F6", "Rename / Move Files").
		k(hook.ModNone, "F7", "New Folder").
		k(hook.ModNone, "F8", "Delete Files").
		k(hook.ModNone, "DELETE", "Delete Files").
		k(hook.ModNone, "F9", "Activate Menu").
		k(hook.ModNone, "F10", "Activate Left Menu")

	// Panels and navigation.
	tc.k(hook.ModNone, "TAB", "Switch Panel").
		k(hook.ModShift, "TAB", "Switch to Tree").
		k(hook.ModNone, "BACKSPACE", "Up One Level").
		k(hook.ModCtrl, "PAGEUP", "Up One Level").
		k(hook.ModCtrl, "PAGEDOWN", "Open Directory / Archive").
		k(hook.ModCtrl, "\\", "Jump to Root").
		k(hook.ModCtrl, "U", "Swap Panels").
		k(hook.ModCtrl|hook.ModShift, "U", "Swap Panels and Tabs").
		k(hook.ModAlt, "LEFT", "Back").
		k(hook.ModAlt, "RIGHT", "Forward").
		k(hook.ModAlt, "DOWN", "Directory History").
		k(hook.ModCtrl, "LEFT", "Open in Target Panel").
		k(hook.ModCtrl, "RIGHT", "Open in Target Panel").
		k(hook.ModAlt, "F1", "Change Left Drive").
		k(hook.ModAlt, "F2", "Change Right Drive").
		k(hook.ModCtrl, "D", "Directory Hotlist").
		k(hook.ModCtrl, "I", "Switch to Target Directory").
		k(hook.ModCtrl, "P", "Copy Path to Command Line")

	// Files and selection.
	tc.k(hook.ModShift, "F4", "New Text File").
		k(hook.ModShift, "F5", "Copy in Place").
		k(hook.ModShift, "F6", "Rename in Place").
		k(hook.ModShift, "F7", "New Folder in Other Panel").
		k(hook.ModShift, "F10", "Context Menu").
		k(hook.ModAlt, "ENTER", "Properties").
		k(hook.ModAlt, "F5", "Pack Files").
		k(hook.ModAlt, "F6", "Unpack Files").
		k(hook.ModAlt, "F9", "Unpack Files").
		k(hook.ModAlt, "F7", "Find Files").
		k(hook.ModCtrl, "M", "Multi-Rename Tool").
		k(hook.ModCtrl, "L", "Calculate Size").
		k(hook.ModCtrl, "Z", "Edit Comment").
		k(hook.ModCtrl, "C", "Copy to Clipboard").
		k(hook.ModCtrl, "X", "Cut to Clipboard").
		k(hook.ModCtrl, "V", "Paste to Current Directory").
		k(hook.ModCtrl, "A", "Select All").
		k(hook.ModNone, "INSERT", "Select File").
		k(hook.ModNone, "SPACE", "Select / Count Size").
		k(hook.ModNone, "NUM+", "Select by Mask").
		k(hook.ModNone, "NUM-", "Deselect by Mask").
		k(hook.ModNone, "NUM*", "Invert Selection")

	// View modes and sorting.
	tc.k(hook.ModCtrl, "F1", "Brief View").
		k(hook.ModCtrl|hook.ModShift, "F1", "Thumbnails View").
		k(hook.ModCtrl, "F2", "Full View").
		k(hook.ModCtrl|hook.ModShift, "F2", "Comments View").
		k(hook.ModCtrl, "F3", "Sort by Name").
		k(hook.ModCtrl, "F4", "Sort by Extension").
		k(hook.ModCtrl, "F5", "Sort by Date").
		k(hook.ModCtrl, "F6", "Sort by Size").
		k(hook.ModCtrl, "F7", "Unsorted").
		k(hook.ModCtrl, "F8", "Directory Tree").
		k(hook.ModCtrl, "F10", "Show All Files").
		k(hook.ModCtrl, "F11", "Show Only Programs").
		k(hook.ModCtrl, "F12", "User-Defined Files").
		k(hook.ModCtrl, "Q", "Quick View")

	// Folder tabs.
	tc.k(hook.ModCtrl, "T", "New Tab").
		k(hook.ModCtrl|hook.ModShift, "T", "New Tab in Background").
		k(hook.ModCtrl, "W", "Close Tab").
		k(hook.ModCtrl|hook.ModShift, "W", "Close All Tabs").
		k(hook.ModCtrl, "TAB", "Next Tab").
		k(hook.ModCtrl|hook.ModShift, "TAB", "Previous Tab")

	// Filtering, branching, and remote connections.
	tc.k(hook.ModCtrl, "S", "Quick Filter").
		k(hook.ModCtrl|hook.ModShift, "S", "Quick Filter (Last Used)").
		k(hook.ModCtrl, "B", "Branch View").
		k(hook.ModCtrl|hook.ModShift, "B", "Branch View (Selected)").
		k(hook.ModCtrl, "F", "Connect to FTP").
		k(hook.ModCtrl|hook.ModShift, "F", "Disconnect from FTP")

	r.addProfile(tc)
}
