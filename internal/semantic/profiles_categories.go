package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// initCategoryDictionaries registers the family fallbacks. A category entry must
// hold in most members of the family, because it also fires for applications
// that have no exact profile yet. Anything that differs between VS Code,
// JetBrains and Dev-C++ belongs in those profiles instead.
func (r *Registry) initCategoryDictionaries() {
	r.initEditorCategory()
	r.initBrowserCategory()
	r.initTerminalCategory()
}

// initEditorCategory holds the actions that are safe across almost all code
// editors.
func (r *Registry) initEditorCategory() {
	cat := context.CategoryEditor

	// Search & navigation
	r.addCat(cat, hook.ModNone, "F3", "Find Next")
	r.addCat(cat, hook.ModShift, "F3", "Find Previous")
	r.addCat(cat, hook.ModCtrl, "G", "Go to Line")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "F", "Find in Files")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "P", "Command Palette")
	r.addCat(cat, hook.ModCtrl, "HOME", "Start of File")
	r.addCat(cat, hook.ModCtrl, "END", "End of File")
	r.addCat(cat, hook.ModCtrl, "PAGEUP", "Previous Editor Tab")
	r.addCat(cat, hook.ModCtrl, "PAGEDOWN", "Next Editor Tab")

	// Symbols & completion
	r.addCat(cat, hook.ModNone, "F12", "Go to Definition")
	r.addCat(cat, hook.ModCtrl, "SPACE", "Code Completion")

	// Editing
	r.addCat(cat, hook.ModCtrl, "/", "Toggle Comment")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "/", "Block Comment")
	r.addCat(cat, hook.ModCtrl, "BACKSPACE", "Delete Word Left")
	r.addCat(cat, hook.ModCtrl, "DELETE", "Delete Word Right")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "UP", "Move Line Up")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "DOWN", "Move Line Down")

	// Windows
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "T", "Reopen Closed Editor")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "N", "New Window")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "W", "Close Window")

	// Debugging (shared by most IDE-grade editors)
	r.addCat(cat, hook.ModNone, "F9", "Toggle Breakpoint")
	r.addCat(cat, hook.ModNone, "F10", "Step Over")
}

// initBrowserCategory covers Chromium/Gecko-style navigation. Edge and Chrome
// agree on nearly everything here.
func (r *Registry) initBrowserCategory() {
	cat := context.CategoryBrowser

	// Tabs & windows
	r.addCat(cat, hook.ModCtrl, "T", "New Tab")
	r.addCat(cat, hook.ModCtrl, "W", "Close Tab")
	r.addCat(cat, hook.ModCtrl, "F4", "Close Tab")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "T", "Reopen Closed Tab")
	r.addCat(cat, hook.ModCtrl, "N", "New Window")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "N", "New Incognito Window")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "W", "Close Window")
	r.addCat(cat, hook.ModCtrl, "PAGEUP", "Previous Tab")
	r.addCat(cat, hook.ModCtrl, "PAGEDOWN", "Next Tab")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "PAGEUP", "Move Tab Left")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "PAGEDOWN", "Move Tab Right")
	for _, key := range []string{"1", "2", "3", "4", "5", "6", "7", "8"} {
		r.addCat(cat, hook.ModCtrl, key, "Switch to Tab")
	}
	r.addCat(cat, hook.ModCtrl, "9", "Last Tab")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "A", "Search Tabs")

	// Bookmarks & history
	r.addCat(cat, hook.ModCtrl, "D", "Bookmark Tab")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "B", "Toggle Bookmarks Bar")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "O", "Bookmark Manager")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "D", "Bookmark All Tabs")
	r.addCat(cat, hook.ModAlt|hook.ModShift, "B", "Focus Bookmarks Bar")
	r.addCat(cat, hook.ModCtrl, "H", "History")
	r.addCat(cat, hook.ModCtrl, "J", "Downloads")

	// Address bar
	r.addCat(cat, hook.ModCtrl, "L", "Focus Address Bar")
	r.addCat(cat, hook.ModAlt, "D", "Focus Address Bar")
	r.addCat(cat, hook.ModNone, "F6", "Focus Address Bar")
	r.addCat(cat, hook.ModCtrl, "E", "Search from Address Bar")
	r.addCat(cat, hook.ModCtrl, "K", "Search from Address Bar")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "V", "Paste and Go")
	r.addCat(cat, hook.ModAlt, "HOME", "Home Page")

	// Page
	r.addCat(cat, hook.ModCtrl, "R", "Refresh")
	r.addCat(cat, hook.ModNone, "F5", "Refresh")
	r.addCat(cat, hook.ModCtrl, "F5", "Hard Refresh")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "R", "Hard Refresh")
	r.addCat(cat, hook.ModAlt, "LEFT", "Back")
	r.addCat(cat, hook.ModAlt, "RIGHT", "Forward")
	r.addCat(cat, hook.ModCtrl, "G", "Find Next")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "G", "Find Previous")
	r.addCat(cat, hook.ModCtrl, "U", "View Source")
	r.addCat(cat, hook.ModCtrl, "P", "Print")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "DELETE", "Clear Browsing Data")
	r.addCat(cat, hook.ModNone, "F11", "Full Screen")
	r.addCat(cat, hook.ModNone, "F7", "Caret Browsing")

	// Developer tools
	r.addCat(cat, hook.ModNone, "F12", "Developer Tools")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "I", "Developer Tools")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "J", "Console")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "C", "Inspect Element")
	r.addCat(cat, hook.ModShift, "ESC", "Browser Task Manager")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "M", "Profile Menu")
}

// initTerminalCategory covers Windows Terminal, PowerShell, cmd and the common
// third-party terminals.
func (r *Registry) initTerminalCategory() {
	cat := context.CategoryTerminal

	// Process control
	r.addCat(cat, hook.ModCtrl, "C", "Interrupt Process")
	r.addCat(cat, hook.ModCtrl, "D", "Exit / EOF")
	r.addCat(cat, hook.ModCtrl, "L", "Clear Terminal")

	// Clipboard
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "C", "Copy Selection")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "V", "Paste")
	r.addCat(cat, hook.ModCtrl, "BACKSPACE", "Delete Word Left")
	r.addCat(cat, hook.ModCtrl, "DELETE", "Delete Word Right")

	// Tabs & panes
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "T", "New Tab")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "N", "New Window")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "W", "Close Tab/Pane")
	r.addCat(cat, hook.ModCtrl, "TAB", "Next Tab")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "TAB", "Previous Tab")
	r.addCat(cat, hook.ModAlt|hook.ModShift, "+", "Split Pane Right")
	r.addCat(cat, hook.ModAlt|hook.ModShift, "-", "Split Pane Down")
	r.addCat(cat, hook.ModAlt|hook.ModShift, "D", "Duplicate Pane")
	r.addCat(cat, hook.ModAlt, "LEFT", "Focus Pane Left")
	r.addCat(cat, hook.ModAlt, "RIGHT", "Focus Pane Right")
	r.addCat(cat, hook.ModAlt, "UP", "Focus Pane Up")
	r.addCat(cat, hook.ModAlt, "DOWN", "Focus Pane Down")
	r.addCat(cat, hook.ModAlt|hook.ModShift, "LEFT", "Resize Pane Left")
	r.addCat(cat, hook.ModAlt|hook.ModShift, "RIGHT", "Resize Pane Right")
	r.addCat(cat, hook.ModAlt|hook.ModShift, "UP", "Resize Pane Up")
	r.addCat(cat, hook.ModAlt|hook.ModShift, "DOWN", "Resize Pane Down")
	for _, key := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"} {
		r.addCat(cat, hook.ModCtrl|hook.ModShift, key, "New Tab by Profile")
	}

	// Search
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "F", "Find")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "P", "Command Palette")

	// Scrolling
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "UP", "Scroll Up")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "DOWN", "Scroll Down")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "PAGEUP", "Scroll Page Up")
	r.addCat(cat, hook.ModCtrl|hook.ModShift, "PAGEDOWN", "Scroll Page Down")
	r.addCat(cat, hook.ModShift, "PAGEUP", "Scroll Up")
	r.addCat(cat, hook.ModShift, "PAGEDOWN", "Scroll Page Down")
	r.addCat(cat, hook.ModCtrl, "HOME", "Scroll to Top")
	r.addCat(cat, hook.ModCtrl, "END", "Scroll to Bottom")

	// Display
	r.addCat(cat, hook.ModNone, "F11", "Full Screen")
	r.addCat(cat, hook.ModAlt, "ENTER", "Toggle Fullscreen")
}
