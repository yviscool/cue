package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// Editors and IDEs. Every profile records the keymap it was transcribed from so
// a future maintainer can re-verify a chord instead of trusting a blog post.

func (r *Registry) initVSCodeFamily() {
	p := newProfile("vscode", "VS Code", context.CategoryEditor,
		"code.exe", "code - insiders.exe", "cursor.exe", "vscodium.exe",
		"windsurf.exe", "trae.exe",
	).source("VS Code Default Keybindings (Windows)", "VS Code 1.9x (Windows)")

	// Command palette & navigation
	p.k(hook.ModCtrl|hook.ModShift, "P", "Command Palette").
		k(hook.ModNone, "F1", "Command Palette").
		k(hook.ModCtrl, "P", "Quick Open").
		k(hook.ModCtrl, "T", "Show All Symbols").
		k(hook.ModCtrl, "G", "Go to Line").
		k(hook.ModCtrl|hook.ModShift, "O", "Go to Symbol").
		k(hook.ModNone, "F12", "Go to Definition").
		k(hook.ModAlt, "F12", "Peek Definition").
		k(hook.ModShift, "F12", "Find References").
		k(hook.ModNone, "F2", "Rename Symbol").
		k(hook.ModCtrl, ".", "Quick Fix").
		k(hook.ModCtrl, "SPACE", "Trigger Suggestion").
		k(hook.ModCtrl|hook.ModShift, "SPACE", "Trigger Parameter Hints")

	// Panels & views
	p.k(hook.ModCtrl, "B", "Toggle Sidebar").
		k(hook.ModCtrl, "J", "Toggle Panel").
		k(hook.ModCtrl, "`", "Toggle Terminal").
		k(hook.ModCtrl|hook.ModShift, "`", "New Terminal").
		k(hook.ModCtrl|hook.ModShift, "E", "Show Explorer").
		k(hook.ModCtrl|hook.ModShift, "G", "Show Source Control").
		k(hook.ModCtrl|hook.ModShift, "D", "Show Run and Debug").
		k(hook.ModCtrl|hook.ModShift, "X", "Show Extensions").
		k(hook.ModCtrl|hook.ModShift, "M", "Show Problems").
		k(hook.ModCtrl|hook.ModShift, "U", "Show Output").
		k(hook.ModCtrl, "\\", "Split Editor").
		k(hook.ModCtrl, ",", "Open Settings")

	// Editors & windows
	p.k(hook.ModCtrl, "W", "Close Editor").
		k(hook.ModCtrl|hook.ModShift, "W", "Close Window").
		k(hook.ModCtrl|hook.ModShift, "N", "New Window").
		k(hook.ModCtrl|hook.ModShift, "T", "Reopen Closed Editor").
		k(hook.ModCtrl, "TAB", "Next Editor").
		k(hook.ModCtrl|hook.ModShift, "TAB", "Previous Editor")

	// Editing
	p.k(hook.ModCtrl, "/", "Toggle Line Comment").
		k(hook.ModCtrl|hook.ModShift, "A", "Toggle Block Comment").
		k(hook.ModAlt, "UP", "Move Line Up").
		k(hook.ModAlt, "DOWN", "Move Line Down").
		k(hook.ModShift|hook.ModAlt, "UP", "Copy Line Up").
		k(hook.ModShift|hook.ModAlt, "DOWN", "Copy Line Down").
		k(hook.ModCtrl|hook.ModShift, "K", "Delete Line").
		k(hook.ModCtrl, "ENTER", "Insert Line Below").
		k(hook.ModCtrl|hook.ModShift, "ENTER", "Insert Line Above").
		k(hook.ModCtrl|hook.ModShift, "\\", "Jump to Matching Bracket").
		k(hook.ModCtrl, "]", "Indent Line").
		k(hook.ModCtrl, "[", "Outdent Line").
		k(hook.ModCtrl|hook.ModShift, "[", "Fold Region").
		k(hook.ModCtrl|hook.ModShift, "]", "Unfold Region")

	// Multi-cursor
	p.k(hook.ModCtrl, "D", "Add Selection to Next Find Match").
		k(hook.ModCtrl|hook.ModShift, "L", "Select All Occurrences").
		k(hook.ModCtrl, "F2", "Select All Occurrences of Current Word").
		k(hook.ModCtrl|hook.ModAlt, "UP", "Add Cursor Above").
		k(hook.ModCtrl|hook.ModAlt, "DOWN", "Add Cursor Below").
		k(hook.ModCtrl, "U", "Undo Last Cursor Operation").
		k(hook.ModShift|hook.ModAlt, "I", "Insert Cursor at End of Each Line")

	// Search & problems
	p.k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "H", "Replace").
		k(hook.ModCtrl|hook.ModShift, "F", "Find in Files").
		k(hook.ModCtrl|hook.ModShift, "H", "Replace in Files").
		k(hook.ModNone, "F8", "Next Problem").
		k(hook.ModShift, "F8", "Previous Problem").
		k(hook.ModShift|hook.ModAlt, "F", "Format Document")

	// Debug
	p.k(hook.ModNone, "F9", "Toggle Breakpoint").
		k(hook.ModNone, "F5", "Start Debugging").
		k(hook.ModCtrl, "F5", "Run Without Debugging").
		k(hook.ModShift, "F5", "Stop Debugging").
		k(hook.ModNone, "F10", "Step Over").
		k(hook.ModNone, "F11", "Step Into").
		k(hook.ModShift, "F11", "Step Out")

	// Two-chord sequences
	p.seq(hook.ModCtrl, "K", hook.ModCtrl, "S", "Keyboard Shortcuts").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "T", "Color Theme").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "W", "Close All Editors").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "C", "Add Line Comment").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "U", "Remove Line Comment").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "F", "Format Selection").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "I", "Show Hover").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "0", "Fold All").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "J", "Unfold All").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "V", "Markdown Preview to the Side").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "Z", "Zen Mode")

	r.addProfile(p)
}

func (r *Registry) initJetBrains() {
	p := newProfile("jetbrains", "JetBrains IDE", context.CategoryEditor,
		"idea64.exe", "pycharm64.exe", "webstorm64.exe", "clion64.exe",
		"goland64.exe", "rider64.exe", "datagrip64.exe", "phpstorm64.exe",
		"rubymine64.exe", "studio64.exe", "aqua64.exe", "writerside64.exe",
		"fleet.exe", "rustrover64.exe",
	).source("JetBrains Default Keymap (Windows)", "IntelliJ Platform 2024.x (Windows)")

	// Navigation
	p.k(hook.ModCtrl|hook.ModShift, "A", "Find Action").
		k(hook.ModCtrl, "N", "Go to Class").
		k(hook.ModCtrl|hook.ModShift, "N", "Go to File").
		k(hook.ModCtrl|hook.ModAlt|hook.ModShift, "N", "Go to Symbol").
		k(hook.ModCtrl, "E", "Recent Files").
		k(hook.ModCtrl|hook.ModShift, "E", "Recent Locations").
		k(hook.ModCtrl, "B", "Go to Declaration").
		k(hook.ModCtrl|hook.ModAlt, "B", "Go to Implementation").
		k(hook.ModCtrl|hook.ModShift, "B", "Go to Type Declaration").
		k(hook.ModCtrl, "U", "Go to Super Method").
		k(hook.ModCtrl, "F12", "File Structure").
		k(hook.ModCtrl, "H", "Type Hierarchy").
		k(hook.ModCtrl|hook.ModAlt, "H", "Call Hierarchy")

	// Search & usages
	p.k(hook.ModAlt, "F7", "Find Usages").
		k(hook.ModCtrl, "F7", "Find Usages in File").
		k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "R", "Replace").
		k(hook.ModCtrl|hook.ModShift, "F", "Find in Files").
		k(hook.ModCtrl|hook.ModShift, "R", "Replace in Files")

	// Editing
	p.k(hook.ModCtrl, "/", "Toggle Comment").
		k(hook.ModCtrl|hook.ModShift, "/", "Block Comment").
		k(hook.ModCtrl, "D", "Duplicate Line").
		k(hook.ModCtrl, "Y", "Delete Line").
		k(hook.ModCtrl|hook.ModShift, "UP", "Move Line Up").
		k(hook.ModCtrl|hook.ModShift, "DOWN", "Move Line Down").
		k(hook.ModCtrl|hook.ModAlt, "L", "Reformat Code").
		k(hook.ModCtrl|hook.ModAlt, "O", "Optimize Imports").
		k(hook.ModCtrl|hook.ModAlt, "I", "Auto-Indent Lines").
		k(hook.ModCtrl, "W", "Extend Selection").
		k(hook.ModCtrl|hook.ModShift, "W", "Shrink Selection").
		k(hook.ModCtrl, "J", "Insert Live Template")

	// Code intelligence
	p.k(hook.ModCtrl, "SPACE", "Code Completion").
		k(hook.ModCtrl|hook.ModShift, "SPACE", "Smart Completion").
		k(hook.ModCtrl|hook.ModShift, "ENTER", "Complete Statement").
		k(hook.ModCtrl, "P", "Parameter Info").
		k(hook.ModCtrl, "Q", "Quick Documentation").
		k(hook.ModAlt, "ENTER", "Show Context Actions").
		k(hook.ModNone, "F2", "Next Highlighted Error").
		k(hook.ModShift, "F2", "Previous Highlighted Error").
		k(hook.ModAlt, "INSERT", "Generate Code").
		k(hook.ModCtrl, "O", "Override Methods").
		k(hook.ModCtrl, "I", "Implement Methods")

	// Refactoring
	p.k(hook.ModShift, "F6", "Rename Symbol").
		k(hook.ModCtrl|hook.ModAlt|hook.ModShift, "T", "Refactor This").
		k(hook.ModCtrl|hook.ModAlt, "V", "Extract Variable").
		k(hook.ModCtrl|hook.ModAlt, "M", "Extract Method").
		k(hook.ModCtrl|hook.ModAlt, "P", "Extract Parameter").
		k(hook.ModCtrl|hook.ModAlt, "N", "Inline").
		k(hook.ModCtrl|hook.ModAlt, "T", "Surround With")

	// Build, run, debug
	p.k(hook.ModShift, "F10", "Run").
		k(hook.ModShift, "F9", "Debug").
		k(hook.ModCtrl, "F2", "Stop").
		k(hook.ModNone, "F8", "Step Over").
		k(hook.ModNone, "F7", "Step Into").
		k(hook.ModShift, "F8", "Step Out").
		k(hook.ModCtrl, "F8", "Toggle Line Breakpoint").
		k(hook.ModNone, "F11", "Toggle Bookmark").
		k(hook.ModShift, "F11", "Show Bookmarks")

	// Version control
	p.k(hook.ModCtrl, "K", "Commit Changes").
		k(hook.ModCtrl|hook.ModShift, "K", "Push Commits").
		k(hook.ModCtrl, "T", "Update Project")

	// Windows
	p.k(hook.ModCtrl, "F4", "Close Editor").
		k(hook.ModShift, "ESC", "Hide Active Tool Window").
		k(hook.ModCtrl|hook.ModShift, "F12", "Hide All Windows")

	r.addProfile(p)
}

func (r *Registry) initVisualStudio() {
	p := newProfile("visualstudio", "Visual Studio", context.CategoryEditor,
		"devenv.exe",
	).source("Visual Studio Default Keyboard Shortcuts", "Visual Studio 2022 (C# profile)")

	// Debug
	p.k(hook.ModNone, "F5", "Start Debugging").
		k(hook.ModCtrl, "F5", "Start Without Debugging").
		k(hook.ModShift, "F5", "Stop Debugging").
		k(hook.ModCtrl|hook.ModShift, "F5", "Restart Debugging").
		k(hook.ModNone, "F9", "Toggle Breakpoint").
		k(hook.ModNone, "F10", "Step Over").
		k(hook.ModNone, "F11", "Step Into").
		k(hook.ModShift, "F11", "Step Out").
		k(hook.ModCtrl|hook.ModAlt, "P", "Attach to Process").
		k(hook.ModCtrl|hook.ModAlt, "E", "Exception Settings")

	// Build & solution
	p.k(hook.ModCtrl|hook.ModShift, "B", "Build Solution").
		k(hook.ModCtrl|hook.ModShift, "N", "New Project").
		k(hook.ModCtrl|hook.ModShift, "A", "Add New Item").
		k(hook.ModCtrl|hook.ModAlt, "L", "Solution Explorer")

	// Navigation
	p.k(hook.ModNone, "F12", "Go to Definition").
		k(hook.ModCtrl, "F12", "Go to Declaration").
		k(hook.ModShift, "F12", "Find References").
		k(hook.ModAlt, "F12", "Peek Definition").
		k(hook.ModCtrl, "T", "Go to All").
		k(hook.ModCtrl, "G", "Go to Line")

	// Search
	p.k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "H", "Replace").
		k(hook.ModCtrl|hook.ModShift, "F", "Find in Files").
		k(hook.ModCtrl|hook.ModShift, "H", "Replace in Files")

	// Two-chord sequences
	p.seq(hook.ModCtrl, "K", hook.ModCtrl, "D", "Format Document").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "F", "Format Selection").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "C", "Comment Selection").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "U", "Uncomment Selection").
		seq(hook.ModCtrl, "M", hook.ModCtrl, "O", "Collapse to Definitions").
		seq(hook.ModCtrl, "M", hook.ModCtrl, "M", "Toggle Outlining").
		seq(hook.ModCtrl, "R", hook.ModCtrl, "R", "Rename Symbol").
		seq(hook.ModCtrl, "R", hook.ModCtrl, "M", "Extract Method").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "K", "Toggle Bookmark")

	r.addProfile(p)
}

func (r *Registry) initEclipse() {
	p := newProfile("eclipse", "Eclipse", context.CategoryEditor,
		"eclipse.exe",
	).source("Eclipse Default Key Bindings", "Eclipse 2024-09 (Windows)")

	p.k(hook.ModNone, "F3", "Open Declaration").
		k(hook.ModCtrl|hook.ModShift, "T", "Open Type").
		k(hook.ModCtrl|hook.ModShift, "R", "Open Resource").
		k(hook.ModCtrl|hook.ModShift, "O", "Organize Imports").
		k(hook.ModCtrl|hook.ModShift, "F", "Format Document").
		k(hook.ModCtrl|hook.ModShift, "G", "Find References").
		k(hook.ModCtrl, "D", "Delete Line").
		k(hook.ModCtrl, "/", "Toggle Comment").
		k(hook.ModCtrl|hook.ModShift, "/", "Block Comment").
		k(hook.ModCtrl|hook.ModShift, "\\", "Remove Block Comment").
		k(hook.ModAlt|hook.ModShift, "R", "Rename Symbol").
		k(hook.ModAlt|hook.ModShift, "M", "Extract Method").
		k(hook.ModAlt|hook.ModShift, "L", "Extract Local Variable").
		k(hook.ModAlt|hook.ModShift, "I", "Inline").
		k(hook.ModAlt|hook.ModShift, "V", "Move").
		k(hook.ModCtrl, "1", "Quick Fix").
		k(hook.ModCtrl, "SPACE", "Code Completion").
		k(hook.ModNone, "F4", "Open Type Hierarchy").
		k(hook.ModCtrl, "L", "Go to Line").
		k(hook.ModCtrl, "H", "Search").
		k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "J", "Incremental Find").
		k(hook.ModCtrl, "K", "Find Next").
		k(hook.ModCtrl|hook.ModShift, "K", "Find Previous").
		k(hook.ModCtrl, "O", "Quick Outline").
		k(hook.ModCtrl, "3", "Quick Access").
		k(hook.ModCtrl, "M", "Maximize Editor").
		k(hook.ModAlt, "UP", "Move Line Up").
		k(hook.ModAlt, "DOWN", "Move Line Down").
		k(hook.ModCtrl|hook.ModAlt, "DOWN", "Duplicate Line")

	// Debug
	p.k(hook.ModCtrl|hook.ModShift, "B", "Toggle Breakpoint").
		k(hook.ModNone, "F5", "Step Into").
		k(hook.ModNone, "F6", "Step Over").
		k(hook.ModNone, "F7", "Step Out").
		k(hook.ModNone, "F8", "Resume").
		k(hook.ModNone, "F11", "Debug").
		k(hook.ModCtrl, "F11", "Run")

	r.addProfile(p)
}

func (r *Registry) initSublime() {
	p := newProfile("sublime", "Sublime Text", context.CategoryEditor,
		"sublime_text.exe",
	).source("Default (Windows).sublime-keymap", "Sublime Text 4 (Build 4180, Windows)")

	// Command palette & navigation
	p.k(hook.ModCtrl|hook.ModShift, "P", "Command Palette").
		k(hook.ModCtrl, "P", "Goto Anything").
		k(hook.ModCtrl, "R", "Goto Symbol").
		k(hook.ModCtrl|hook.ModShift, "R", "Goto Symbol in Project").
		k(hook.ModCtrl, "G", "Goto Line")

	// Search
	p.k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "H", "Replace").
		k(hook.ModCtrl|hook.ModShift, "F", "Find in Files").
		k(hook.ModNone, "F4", "Next Result").
		k(hook.ModShift, "F4", "Previous Result")

	// Multi-cursor & selection
	p.k(hook.ModCtrl, "D", "Add Selection to Next Find Match").
		k(hook.ModAlt, "F3", "Select All Occurrences").
		k(hook.ModCtrl|hook.ModShift, "L", "Split Selection into Lines").
		k(hook.ModCtrl, "L", "Expand Selection to Line").
		k(hook.ModCtrl|hook.ModShift, "SPACE", "Expand Selection to Scope").
		k(hook.ModCtrl|hook.ModShift, "M", "Expand Selection to Brackets").
		k(hook.ModCtrl, "M", "Jump to Matching Bracket").
		k(hook.ModCtrl|hook.ModShift, "J", "Expand Selection to Indentation")

	// Line operations
	p.k(hook.ModCtrl, "ENTER", "Insert Line After").
		k(hook.ModCtrl|hook.ModShift, "ENTER", "Insert Line Before").
		k(hook.ModCtrl|hook.ModShift, "UP", "Move Line Up").
		k(hook.ModCtrl|hook.ModShift, "DOWN", "Move Line Down").
		k(hook.ModCtrl, "J", "Join Lines").
		k(hook.ModCtrl|hook.ModShift, "D", "Duplicate Line").
		k(hook.ModCtrl|hook.ModShift, "K", "Delete Line").
		k(hook.ModCtrl|hook.ModShift, "V", "Paste and Indent").
		k(hook.ModNone, "F9", "Sort Lines")

	// Comments
	p.k(hook.ModCtrl, "/", "Toggle Comment").
		k(hook.ModCtrl|hook.ModShift, "/", "Toggle Block Comment")

	// Files & windows
	p.k(hook.ModCtrl, "N", "New File").
		k(hook.ModCtrl, "W", "Close File").
		k(hook.ModCtrl|hook.ModShift, "S", "Save As").
		k(hook.ModCtrl|hook.ModShift, "N", "New Window").
		k(hook.ModCtrl|hook.ModShift, "W", "Close Window").
		k(hook.ModCtrl|hook.ModShift, "T", "Reopen Closed File").
		k(hook.ModCtrl, "`", "Show Console").
		k(hook.ModNone, "F11", "Full Screen").
		k(hook.ModShift, "F11", "Distraction Free Mode")

	// Two-chord sequences
	p.seq(hook.ModCtrl, "K", hook.ModCtrl, "B", "Toggle Sidebar").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "U", "Upper Case").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "L", "Lower Case").
		seq(hook.ModCtrl, "K", hook.ModCtrl, "K", "Delete to End of Line")

	r.addProfile(p)
}

func (r *Registry) initDevCpp() {
	p := newProfile("devcpp", "Dev-C++", context.CategoryEditor,
		"devcpp.exe", "devc++.exe", "dev-c++.exe",
	).source("Dev-C++ default shortcut list", "Dev-C++ 6.3 / 5.11 (Windows)")

	// Build & run
	p.k(hook.ModNone, "F9", "Compile").
		k(hook.ModNone, "F10", "Run").
		k(hook.ModNone, "F11", "Compile & Run").
		k(hook.ModCtrl, "F9", "Compile").
		k(hook.ModCtrl, "F10", "Run").
		k(hook.ModCtrl, "F11", "Compile & Run")

	// Debug
	p.k(hook.ModNone, "F5", "Debug").
		k(hook.ModNone, "F6", "Stop Execution").
		k(hook.ModCtrl|hook.ModAlt, "F2", "Terminate Program")

	// Source formatting & indentation
	p.k(hook.ModCtrl|hook.ModShift, "A", "Format Code").
		k(hook.ModCtrl, "I", "Indent").
		k(hook.ModCtrl|hook.ModShift, "I", "Unindent").
		k(hook.ModCtrl|hook.ModShift, "UP", "Move Line Up").
		k(hook.ModCtrl|hook.ModShift, "DOWN", "Move Line Down")

	// Files
	p.k(hook.ModCtrl, "N", "New File").
		k(hook.ModCtrl, "O", "Open File").
		k(hook.ModCtrl, "S", "Save").
		k(hook.ModCtrl|hook.ModShift, "S", "Save All").
		k(hook.ModCtrl, "F4", "Close File")

	// Search & navigation
	p.k(hook.ModCtrl, "G", "Go to Line").
		k(hook.ModCtrl, "F", "Find").
		k(hook.ModNone, "F3", "Find Next").
		k(hook.ModShift, "F3", "Find Previous").
		k(hook.ModCtrl, "R", "Replace").
		k(hook.ModCtrl|hook.ModShift, "F", "Find in Files")

	// Editing
	p.k(hook.ModCtrl, "/", "Toggle Comment").
		k(hook.ModCtrl, "SPACE", "Code Completion").
		k(hook.ModCtrl, "BACKSPACE", "Delete Word Left").
		k(hook.ModCtrl, "DELETE", "Delete Word Right")

	r.addProfile(p)
}

func (r *Registry) initIDLE() {
	// python.exe is ambiguous (scripts, tools, REPL). Title hints keep IDLE
	// detection exact without mislabeling every Python process. idle.exe is
	// unambiguous and is accepted without a hint.
	p := newProfile("idle", "IDLE", context.CategoryEditor,
		"python.exe", "pythonw.exe", "idle.exe", "idlew.exe",
	).titleHints("idle", " shell").ambiguous("python.exe", "pythonw.exe").
		source("IDLE keys.py keydefs", "Python 3.12 IDLE (Windows)")

	// Run
	p.k(hook.ModNone, "F5", "Run Module").
		k(hook.ModShift, "F5", "Run Customized").
		k(hook.ModAlt, "X", "Check Module").
		k(hook.ModNone, "F1", "Help")

	// Files
	p.k(hook.ModCtrl, "N", "New File").
		k(hook.ModCtrl, "O", "Open File").
		k(hook.ModCtrl, "S", "Save").
		k(hook.ModCtrl|hook.ModShift, "S", "Save As").
		k(hook.ModCtrl, "W", "Close").
		k(hook.ModCtrl, "Q", "Quit")

	// Search
	p.k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "G", "Find Next").
		k(hook.ModCtrl|hook.ModShift, "G", "Find Previous").
		k(hook.ModCtrl, "H", "Replace").
		k(hook.ModAlt, "G", "Go to Line")

	// Editing
	p.k(hook.ModAlt, "3", "Comment Region").
		k(hook.ModAlt, "4", "Uncomment Region").
		k(hook.ModCtrl, "]", "Indent Region").
		k(hook.ModCtrl, "[", "Dedent Region").
		k(hook.ModNone, "TAB", "Indent").
		k(hook.ModShift, "TAB", "Dedent").
		k(hook.ModCtrl, "Z", "Undo").
		k(hook.ModCtrl|hook.ModShift, "Z", "Redo").
		k(hook.ModCtrl, "A", "Select All")

	// Shell history
	p.k(hook.ModAlt, "P", "Previous Command").
		k(hook.ModAlt, "N", "Next Command")

	r.addProfile(p)
}

func (r *Registry) initNotepadPP() {
	p := newProfile("notepadpp", "Notepad++", context.CategoryEditor,
		"notepad++.exe",
	).source("Notepad++ Shortcut Mapper defaults", "Notepad++ 8.6 (Windows)")

	// Files
	p.k(hook.ModCtrl, "N", "New File").
		k(hook.ModCtrl, "O", "Open File").
		k(hook.ModCtrl, "S", "Save").
		k(hook.ModCtrl|hook.ModShift, "S", "Save As").
		k(hook.ModCtrl|hook.ModAlt, "S", "Save All").
		k(hook.ModCtrl, "W", "Close").
		k(hook.ModCtrl, "P", "Print")

	// Search
	p.k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "H", "Replace").
		k(hook.ModNone, "F3", "Find Next").
		k(hook.ModShift, "F3", "Find Previous").
		k(hook.ModCtrl|hook.ModShift, "F", "Find in Files").
		k(hook.ModCtrl, "G", "Go to Line")

	// Line operations
	p.k(hook.ModCtrl, "D", "Duplicate Line").
		k(hook.ModCtrl, "L", "Delete Line").
		k(hook.ModCtrl, "T", "Transpose Line").
		k(hook.ModCtrl|hook.ModShift, "T", "Copy Current Line").
		k(hook.ModCtrl|hook.ModShift, "UP", "Move Line Up").
		k(hook.ModCtrl|hook.ModShift, "DOWN", "Move Line Down")
	p.k(hook.ModCtrl, "U", "Lower Case").
		k(hook.ModCtrl|hook.ModShift, "U", "Upper Case")

	// Comments & folding
	p.k(hook.ModCtrl, "Q", "Toggle Block Comment").
		k(hook.ModCtrl|hook.ModShift, "Q", "Toggle Stream Comment").
		k(hook.ModAlt, "0", "Fold All").
		k(hook.ModAlt|hook.ModShift, "0", "Unfold All").
		k(hook.ModCtrl|hook.ModAlt, "F", "Fold Current Level").
		k(hook.ModCtrl|hook.ModAlt|hook.ModShift, "F", "Unfold Current Level")

	// Editing
	p.k(hook.ModCtrl, "BACKSPACE", "Delete Word Left").
		k(hook.ModCtrl, "DELETE", "Delete Word Right").
		k(hook.ModNone, "F11", "Full Screen")

	r.addProfile(p)
}

func (r *Registry) initZed() {
	p := newProfile("zed", "Zed", context.CategoryEditor,
		"zed.exe",
	).source("Zed default.json keymap (Windows)", "Zed 0.16x (Windows)")

	p.k(hook.ModCtrl|hook.ModShift, "P", "Command Palette").
		k(hook.ModCtrl, "P", "Quick Open").
		k(hook.ModCtrl|hook.ModShift, "F", "Project Search").
		k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "G", "Go to Line").
		k(hook.ModCtrl, "B", "Toggle Left Dock").
		k(hook.ModCtrl, "J", "Toggle Bottom Dock").
		k(hook.ModCtrl|hook.ModShift, "E", "Toggle Project Panel").
		k(hook.ModCtrl, "`", "Toggle Terminal").
		k(hook.ModCtrl, "D", "Add Selection to Next Find Match").
		k(hook.ModCtrl|hook.ModShift, "L", "Select All Occurrences").
		k(hook.ModAlt, "UP", "Move Line Up").
		k(hook.ModAlt, "DOWN", "Move Line Down").
		k(hook.ModCtrl|hook.ModShift, "K", "Delete Line").
		k(hook.ModCtrl, "/", "Toggle Comment").
		k(hook.ModNone, "F12", "Go to Definition").
		k(hook.ModCtrl, ",", "Open Settings").
		k(hook.ModCtrl|hook.ModShift, "T", "Reopen Closed Item").
		k(hook.ModCtrl, "W", "Close Tab").
		k(hook.ModCtrl|hook.ModShift, "W", "Close Window").
		k(hook.ModCtrl|hook.ModShift, "N", "New Window").
		k(hook.ModCtrl|hook.ModShift, "S", "Save As").
		k(hook.ModNone, "F11", "Full Screen").
		k(hook.ModCtrl, "Q", "Quit")

	r.addProfile(p)
}

func (r *Registry) initBrowsers() {
	r.addProfile(newProfile("browser", "Browser", context.CategoryBrowser,
		"chrome.exe", "msedge.exe", "firefox.exe", "brave.exe", "vivaldi.exe",
		"opera.exe", "arc.exe", "waterfox.exe", "librewolf.exe", "thorium.exe",
	).source("Chromium/Gecko default shortcuts", "Chrome 13x / Edge 13x (Windows)"))
}

func (r *Registry) initTerminals() {
	r.addProfile(newProfile("terminal", "Terminal", context.CategoryTerminal,
		"windowsterminal.exe", "wt.exe", "cmd.exe", "powershell.exe", "pwsh.exe",
		"conhost.exe", "bash.exe", "git-bash.exe", "mintty.exe", "wezterm-gui.exe",
		"wezterm.exe", "alacritty.exe", "kitty.exe", "tabby.exe", "hyper.exe",
		"conemu64.exe", "conemu.exe",
	).source("Windows Terminal default keybindings", "Windows Terminal 1.2x (Windows)"))
}
