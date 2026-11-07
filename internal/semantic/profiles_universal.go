package semantic

import "github.com/yviscool/cue/internal/hook"

// initUniversal registers shortcuts that mean the same thing in virtually every
// Windows application. App-specific interpretations live in profiles and take
// precedence, so e.g. Ctrl+P is Print here but Quick Open in editors.
//
// Rule of thumb for this table: if a shortcut would teach the user something
// wrong in a mainstream app, it does not belong here. PowerPoint's B (black
// screen) and VLC's F (fullscreen) stay in their profiles.
func (r *Registry) initUniversal() {
	// --- Clipboard & editing ---------------------------------------------
	r.addUc(hook.ModCtrl, "C", "Copy", "Edit")
	r.addUc(hook.ModCtrl, "X", "Cut", "Edit")
	r.addUc(hook.ModCtrl, "V", "Paste", "Edit")
	r.addUc(hook.ModCtrl, "Z", "Undo", "Edit")
	r.addUc(hook.ModCtrl, "Y", "Redo", "Edit")
	r.addUc(hook.ModCtrl|hook.ModShift, "Z", "Redo", "Edit")
	r.addUc(hook.ModCtrl, "A", "Select All", "Edit")
	// Legacy Windows editing chords still work in most text surfaces.
	r.addUc(hook.ModCtrl, "INSERT", "Copy", "Edit")
	r.addUc(hook.ModShift, "INSERT", "Paste", "Edit")
	r.addUc(hook.ModShift, "DELETE", "Cut", "Edit")

	// --- File & search ----------------------------------------------------
	r.addUc(hook.ModCtrl, "S", "Save", "File")
	r.addUc(hook.ModCtrl|hook.ModShift, "S", "Save As", "File")
	r.addUc(hook.ModCtrl, "O", "Open", "File")
	r.addUc(hook.ModCtrl, "N", "New", "File")
	r.addUc(hook.ModCtrl, "P", "Print", "File")
	r.addUc(hook.ModCtrl, "F", "Find", "Search")
	r.addUc(hook.ModCtrl, "H", "Replace", "Search")

	// --- Document & tab lifetime -----------------------------------------
	r.addUc(hook.ModCtrl, "W", "Close", "File")
	r.addUc(hook.ModCtrl, "F4", "Close Document", "File")
	r.addUc(hook.ModCtrl, "TAB", "Next Tab", "Navigation")
	r.addUc(hook.ModCtrl|hook.ModShift, "TAB", "Previous Tab", "Navigation")

	// --- Zoom (browsers, editors, Office, image viewers) ------------------
	r.addUc(hook.ModCtrl, "+", "Zoom In", "View")
	r.addUc(hook.ModCtrl, "=", "Zoom In", "View")
	r.addUc(hook.ModCtrl|hook.ModShift, "+", "Zoom In", "View")
	r.addUc(hook.ModCtrl, "SCROLL UP", "Zoom In", "View")
	r.addUc(hook.ModCtrl, "-", "Zoom Out", "View")
	r.addUc(hook.ModCtrl|hook.ModShift, "-", "Zoom Out", "View")
	r.addUc(hook.ModCtrl, "SCROLL DOWN", "Zoom Out", "View")
	r.addUc(hook.ModCtrl, "0", "Reset Zoom", "View")

	// --- Windows shell ----------------------------------------------------
	r.addUc(hook.ModSuper, "D", "Show Desktop", "System")
	r.addUc(hook.ModSuper, "E", "File Explorer", "System")
	r.addUc(hook.ModSuper, "L", "Lock PC", "System")
	r.addUc(hook.ModSuper, "V", "Clipboard History", "System")
	r.addUc(hook.ModSuper|hook.ModShift, "S", "Snipping Tool", "System")
	r.addUc(hook.ModSuper, "R", "Run Dialog", "System")
	r.addUc(hook.ModSuper, "I", "Settings", "System")
	r.addUc(hook.ModSuper, "A", "Action Center", "System")
	r.addUc(hook.ModSuper, "S", "Search", "System")
	r.addUc(hook.ModSuper, "TAB", "Task View", "System")
	r.addUc(hook.ModSuper, "P", "Project Display", "System")
	r.addUc(hook.ModSuper, "X", "Power User Menu", "System")
	r.addUc(hook.ModSuper, "G", "Game Bar", "System")
	r.addUc(hook.ModSuper, "H", "Voice Typing", "System")
	r.addUc(hook.ModSuper, "K", "Cast", "System")
	r.addUc(hook.ModSuper, "N", "Notification Center", "System")
	r.addUc(hook.ModSuper, "W", "Widgets", "System")
	r.addUc(hook.ModSuper, "Z", "Snap Layouts", "System")
	r.addUc(hook.ModSuper, ",", "Peek Desktop", "System")
	r.addUc(hook.ModSuper, ".", "Emoji Picker", "System")
	r.addUc(hook.ModSuper, ";", "Emoji Picker", "System")
	r.addUc(hook.ModSuper, "PAUSE", "System Properties", "System")
	r.addUc(hook.ModSuper, "U", "Accessibility Settings", "System")
	r.addUc(hook.ModSuper, "F", "Feedback Hub", "System")
	r.addUc(hook.ModSuper, "T", "Cycle Taskbar", "System")
	r.addUc(hook.ModSuper, "B", "Focus Notification Area", "System")
	r.addUc(hook.ModSuper, "SPACE", "Switch Input Language", "System")
	r.addUc(hook.ModSuper, "PRINTSCREEN", "Save Screenshot", "System")
	r.addUc(hook.ModSuper|hook.ModAlt, "PRINTSCREEN", "Game Bar Screenshot", "System")
	r.addUc(hook.ModCtrl, "ESC", "Start Menu", "System")

	// Taskbar slots: Win+1 … Win+9 launch or focus the nth pinned app.
	for _, key := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"} {
		r.addUc(hook.ModSuper, key, "Taskbar App", "System")
	}

	// --- Window management ------------------------------------------------
	r.addUc(hook.ModSuper, "UP", "Maximize Window", "Window")
	r.addUc(hook.ModSuper, "DOWN", "Restore / Minimize", "Window")
	r.addUc(hook.ModSuper, "LEFT", "Snap Left", "Window")
	r.addUc(hook.ModSuper, "RIGHT", "Snap Right", "Window")
	r.addUc(hook.ModSuper|hook.ModShift, "LEFT", "Move to Previous Monitor", "Window")
	r.addUc(hook.ModSuper|hook.ModShift, "RIGHT", "Move to Next Monitor", "Window")
	r.addUc(hook.ModSuper, "HOME", "Minimize All But Active", "Window")
	r.addUc(hook.ModAlt, "TAB", "Switch App", "Window")
	r.addUc(hook.ModAlt|hook.ModShift, "TAB", "Switch App", "Window")
	r.addUc(hook.ModCtrl|hook.ModAlt, "TAB", "App Switcher", "Window")
	r.addUc(hook.ModAlt, "F4", "Close Window", "Window")
	r.addUc(hook.ModAlt, "SPACE", "Window Menu", "Window")

	// --- Virtual desktops -------------------------------------------------
	r.addUc(hook.ModSuper|hook.ModCtrl, "D", "New Virtual Desktop", "System")
	r.addUc(hook.ModSuper|hook.ModCtrl, "LEFT", "Previous Virtual Desktop", "System")
	r.addUc(hook.ModSuper|hook.ModCtrl, "RIGHT", "Next Virtual Desktop", "System")
	r.addUc(hook.ModSuper|hook.ModCtrl, "F4", "Close Virtual Desktop", "System")

	// --- Accessibility & low-level ---------------------------------------
	r.addUc(hook.ModSuper, "+", "Magnifier Zoom In", "System")
	r.addUc(hook.ModSuper, "-", "Magnifier Zoom Out", "System")
	r.addUc(hook.ModSuper, "ESC", "Exit Magnifier", "System")
	r.addUc(hook.ModSuper|hook.ModCtrl, "ENTER", "Narrator", "System")
	r.addUc(hook.ModSuper|hook.ModCtrl|hook.ModShift, "B", "Reset Graphics Driver", "System")
	r.addUc(hook.ModSuper|hook.ModAlt, "R", "Start/Stop Recording", "System")

	// --- Context menu & help ---------------------------------------------
	r.addUc(hook.ModShift, "F10", "Context Menu", "Edit")
	r.addUc(hook.ModNone, "APPS", "Context Menu", "Edit")
	r.addUc(hook.ModNone, "F1", "Help", "System")
	r.addUc(hook.ModNone, "PRINTSCREEN", "Screenshot", "System")
	r.addUc(hook.ModAlt, "PRINTSCREEN", "Active Window Screenshot", "System")
	r.addUc(hook.ModCtrl|hook.ModShift, "ESC", "Task Manager", "System")

	// --- Hardware media keys ---------------------------------------------
	r.addUc(hook.ModNone, "VOLMUTE", "Mute Audio", "Media")
	r.addUc(hook.ModNone, "VOLUP", "Volume Up", "Media")
	r.addUc(hook.ModNone, "VOLDOWN", "Volume Down", "Media")
	r.addUc(hook.ModNone, "MEDIAPLAYPAUSE", "Play/Pause", "Media")
	r.addUc(hook.ModNone, "MEDIANEXT", "Next Track", "Media")
	r.addUc(hook.ModNone, "MEDIAPREV", "Previous Track", "Media")
	r.addUc(hook.ModNone, "MEDIASTOP", "Stop Playback", "Media")
}
