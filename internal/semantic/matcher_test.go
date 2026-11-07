package semantic_test

import (
	"sync"
	"testing"

	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/semantic"
)

func app(exe, title string) context.AppContext {
	return context.AppContext{Exe: exe, Title: title}
}

func TestUniversalShortcuts(t *testing.T) {
	m := semantic.NewMatcher()

	tests := []struct {
		name          string
		mods          hook.Modifiers
		keyName       string
		expectedTitle string
	}{
		{"Ctrl+C", hook.ModCtrl, "C", "Copy"},
		{"Ctrl+c lowercase", hook.ModCtrl, "c", "Copy"},
		{"Ctrl+X", hook.ModCtrl, "X", "Cut"},
		{"Ctrl+V", hook.ModCtrl, "V", "Paste"},
		{"Ctrl+Z", hook.ModCtrl, "Z", "Undo"},
		{"Ctrl+Y", hook.ModCtrl, "Y", "Redo"},
		{"Ctrl+Shift+Z", hook.ModCtrl | hook.ModShift, "Z", "Redo"},
		{"Ctrl+A", hook.ModCtrl, "A", "Select All"},
		{"Ctrl+S", hook.ModCtrl, "S", "Save"},
		{"Ctrl+F", hook.ModCtrl, "F", "Find"},
		{"Ctrl+H", hook.ModCtrl, "H", "Replace"},
		{"Ctrl+P universal print", hook.ModCtrl, "P", "Print"},
		{"Alt+Tab", hook.ModAlt, "Tab", "Switch App"},
		{"Win+D", hook.ModSuper, "D", "Show Desktop"},
		{"Win+E", hook.ModSuper, "E", "File Explorer"},
		{"Win+L", hook.ModSuper, "L", "Lock PC"},
		{"Win+V", hook.ModSuper, "V", "Clipboard History"},
		{"Win+Shift+S", hook.ModSuper | hook.ModShift, "S", "Snipping Tool"},
		{"Alt+F4", hook.ModAlt, "F4", "Close Window"},
		{"Ctrl+Plus", hook.ModCtrl, "+", "Zoom In"},
		{"Ctrl+=", hook.ModCtrl, "=", "Zoom In"},
		{"Ctrl+Minus", hook.ModCtrl, "-", "Zoom Out"},
		{"Ctrl+0", hook.ModCtrl, "0", "Reset Zoom"},
		{"F1", hook.ModNone, "F1", "Help"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := m.MatchAction(context.AppContext{}, tt.mods, tt.keyName)
			if !res.Matched() {
				t.Fatalf("expected shortcut %s to match, but got no match", tt.name)
			}
			if res.Action.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, res.Action.Title)
			}
			if res.Confidence != semantic.ConfidenceUniversal {
				t.Errorf("expected universal confidence, got %s", res.Confidence)
			}
		})
	}
}

func TestVSCodeProfile(t *testing.T) {
	m := semantic.NewMatcher()
	ctx := app("Code.exe", "main.go - Visual Studio Code")

	tests := []struct {
		name          string
		mods          hook.Modifiers
		keyName       string
		expectedTitle string
	}{
		{"Ctrl+Shift+P", hook.ModCtrl | hook.ModShift, "P", "Command Palette"},
		{"Ctrl+P", hook.ModCtrl, "P", "Quick Open"},
		{"Ctrl+B", hook.ModCtrl, "B", "Toggle Sidebar"},
		{"Ctrl+Backtick", hook.ModCtrl, "`", "Toggle Terminal"},
		{"Ctrl+Grave", hook.ModCtrl, "Grave", "Toggle Terminal"},
		{"Ctrl+Slash", hook.ModCtrl, "/", "Toggle Line Comment"},
		{"Alt+Up", hook.ModAlt, "Up", "Move Line Up"},
		{"Alt+ArrowUp", hook.ModAlt, "ArrowUp", "Move Line Up"},
		{"Alt+Down", hook.ModAlt, "Down", "Move Line Down"},
		{"Shift+Alt+Down", hook.ModShift | hook.ModAlt, "Down", "Copy Line Down"},
		{"Ctrl+D", hook.ModCtrl, "D", "Add Selection to Next Find Match"},
		{"F12", hook.ModNone, "F12", "Go to Definition"},
		{"Shift+F12", hook.ModShift, "F12", "Find References"},
		{"Ctrl+Shift+O", hook.ModCtrl | hook.ModShift, "O", "Go to Symbol"},
		{"Ctrl+Shift+Pipe", hook.ModCtrl | hook.ModShift, "|", "Jump to Matching Bracket"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := m.MatchAction(ctx, tt.mods, tt.keyName)
			if !res.Matched() {
				t.Fatalf("expected shortcut %s to match, but got no match", tt.name)
			}
			if res.Action.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, res.Action.Title)
			}
			if res.Confidence != semantic.ConfidenceProfile {
				t.Errorf("expected profile confidence, got %s", res.Confidence)
			}
			if res.AppName != "VS Code" {
				t.Errorf("expected app name VS Code, got %q", res.AppName)
			}
		})
	}

	// VS Code now defines Ctrl+G itself, so it resolves at profile level.
	res := m.MatchAction(ctx, hook.ModCtrl, "G")
	if res.Action.Title != "Go to Line" || res.Confidence != semantic.ConfidenceProfile {
		t.Errorf("expected profile-level Go to Line, got %q (%s)", res.Action.Title, res.Confidence)
	}

	// An editor without a profile still falls back to the category dictionary.
	generic := context.AppContext{Exe: "unknown-editor.exe", Category: context.CategoryEditor}
	genericRes := m.MatchAction(generic, hook.ModCtrl, "G")
	if genericRes.Action.Title != "Go to Line" || genericRes.Confidence != semantic.ConfidenceCategory {
		t.Errorf("expected category-level Go to Line, got %q (%s)", genericRes.Action.Title, genericRes.Confidence)
	}
}

func TestDevCppProfile(t *testing.T) {
	m := semantic.NewMatcher()
	ctx := app("devcpp.exe", "main.cpp - Dev-C++ 6.3")

	tests := []struct {
		name          string
		mods          hook.Modifiers
		keyName       string
		expectedTitle string
	}{
		{"F9 compile", hook.ModNone, "F9", "Compile"},
		{"F10 run", hook.ModNone, "F10", "Run"},
		{"F11 compile and run", hook.ModNone, "F11", "Compile & Run"},
		{"F5 debug", hook.ModNone, "F5", "Debug"},
		{"F6 stop", hook.ModNone, "F6", "Stop Execution"},
		{"Ctrl+Shift+F find in files", hook.ModCtrl | hook.ModShift, "F", "Find in Files"},
		{"Ctrl+Slash comment", hook.ModCtrl, "/", "Toggle Comment"},
		{"Ctrl+Space completion", hook.ModCtrl, "SPACE", "Code Completion"},
		{"Ctrl+C still copies", hook.ModCtrl, "C", "Copy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := m.MatchAction(ctx, tt.mods, tt.keyName)
			if !res.Matched() {
				t.Fatalf("expected shortcut %s to match, but got no match", tt.name)
			}
			if res.Action.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, res.Action.Title)
			}
		})
	}
}

func TestIDLEProfile(t *testing.T) {
	m := semantic.NewMatcher()
	ctx := app("python.exe", "Untitled - IDLE (Python 3.12)")

	tests := []struct {
		name          string
		mods          hook.Modifiers
		keyName       string
		expectedTitle string
	}{
		{"F5 run module", hook.ModNone, "F5", "Run Module"},
		{"Ctrl+N new file", hook.ModCtrl, "N", "New File"},
		{"Ctrl+O open file", hook.ModCtrl, "O", "Open File"},
		{"Alt+G go to line", hook.ModAlt, "G", "Go to Line"},
		{"Alt+3 comment region", hook.ModAlt, "3", "Comment Region"},
		{"Alt+4 uncomment region", hook.ModAlt, "4", "Uncomment Region"},
		{"Ctrl+] indent", hook.ModCtrl, "]", "Indent Region"},
		{"Ctrl+[ dedent", hook.ModCtrl, "[", "Dedent Region"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := m.MatchAction(ctx, tt.mods, tt.keyName)
			if !res.Matched() {
				t.Fatalf("expected shortcut %s to match, but got no match", tt.name)
			}
			if res.Action.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, res.Action.Title)
			}
			if res.AppName != "IDLE" {
				t.Errorf("expected app name IDLE, got %q", res.AppName)
			}
		})
	}
}

func TestIDLEShellAndDirectExecutable(t *testing.T) {
	m := semantic.NewMatcher()

	// The IDLE Shell window is titled "Python 3.x Shell" rather than "IDLE".
	shell := app("pythonw.exe", "Python 3.12.0 Shell")
	if res := m.MatchAction(shell, hook.ModNone, "F5"); res.Action.Title != "Run Module" || res.AppName != "IDLE" {
		t.Errorf("expected IDLE Run Module for the IDLE shell, got %q (%q)", res.Action.Title, res.AppName)
	}

	// idle.exe itself is unambiguous and needs no title hint.
	direct := app("idle.exe", "Python Editor")
	if res := m.MatchAction(direct, hook.ModNone, "F5"); res.Action.Title != "Run Module" {
		t.Errorf("expected Run Module for idle.exe, got %q", res.Action.Title)
	}

	// A plain Python console must remain generic.
	plain := app("python.exe", "Python 3.12.0")
	if res := m.MatchAction(plain, hook.ModNone, "F5"); res.Matched() {
		t.Errorf("plain python.exe F5 should not resolve to IDLE, got %q", res.Action.Title)
	}
}

func TestPythonWithoutIDLETitleIsNotMislabeled(t *testing.T) {
	m := semantic.NewMatcher()
	ctx := app("python.exe", "train.py - Python")

	// Ctrl+N is universal "New"; it must not become IDLE's "New File".
	res := m.MatchAction(ctx, hook.ModCtrl, "N")
	if res.Action.Title != "New" {
		t.Fatalf("expected universal 'New' for plain python.exe, got %q", res.Action.Title)
	}
	if res.Confidence != semantic.ConfidenceUniversal {
		t.Fatalf("expected universal confidence for plain python.exe, got %s", res.Confidence)
	}
	if res.AppName == "IDLE" {
		t.Fatal("plain python.exe must never be labelled IDLE without a title hint")
	}
}

func TestBrowserAndTerminalCategories(t *testing.T) {
	m := semantic.NewMatcher()

	browser := app("chrome.exe", "Example - Google Chrome")
	if res := m.MatchAction(browser, hook.ModCtrl, "D"); res.Action.Title != "Bookmark Tab" {
		t.Errorf("expected 'Bookmark Tab' in browser, got %q", res.Action.Title)
	}
	if res := m.MatchAction(browser, hook.ModCtrl|hook.ModShift, "T"); res.Action.Title != "Reopen Closed Tab" {
		t.Errorf("expected 'Reopen Closed Tab' in browser, got %q", res.Action.Title)
	}

	terminal := app("WindowsTerminal.exe", "PowerShell")
	if res := m.MatchAction(terminal, hook.ModCtrl, "C"); res.Action.Title != "Interrupt Process" {
		t.Errorf("expected 'Interrupt Process' in terminal, got %q", res.Action.Title)
	}
	if res := m.MatchAction(terminal, hook.ModCtrl|hook.ModShift, "T"); res.Action.Title != "New Tab" {
		t.Errorf("expected 'New Tab' in terminal, got %q", res.Action.Title)
	}

	// A plain editor falls back to the universal meaning for Ctrl+C.
	editor := app("notepad.exe", "Untitled - Notepad")
	res := m.MatchAction(editor, hook.ModCtrl, "C")
	if res.Action.Title != "Copy" || res.Confidence != semantic.ConfidenceUniversal {
		t.Errorf("expected universal Copy in Notepad, got %q (%s)", res.Action.Title, res.Confidence)
	}
	// And no accidental app-specific F9 mapping.
	if res := m.MatchAction(editor, hook.ModNone, "F9"); res.Matched() {
		t.Errorf("Notepad F9 should have no semantic match, got %q", res.Action.Title)
	}
}

func TestOfficeProfiles(t *testing.T) {
	m := semantic.NewMatcher()

	excel := app("EXCEL.EXE", "Book1 - Excel")
	if res := m.MatchAction(excel, hook.ModCtrl, "1"); res.Action.Title != "Format Cells" {
		t.Errorf("expected Format Cells, got %q", res.Action.Title)
	}
	if res := m.MatchAction(excel, hook.ModCtrl|hook.ModShift, ":"); res.Action.Title != "Insert Time" {
		t.Errorf("expected Insert Time, got %q", res.Action.Title)
	}

	powerpoint := app("POWERPNT.EXE", "Deck - PowerPoint")
	if res := m.MatchAction(powerpoint, hook.ModNone, "F5"); res.Action.Title != "Start Slideshow" {
		t.Errorf("expected Start Slideshow, got %q", res.Action.Title)
	}
}

func TestSequences(t *testing.T) {
	m := semantic.NewMatcher()
	ctx := app("code.exe", "main.go - Visual Studio Code")

	res, ok := m.MatchSequence(ctx, hook.ModCtrl, "K", hook.ModCtrl, "S")
	if !ok {
		t.Fatal("expected Ctrl+K Ctrl+S to resolve")
	}
	// In VS Code Ctrl+K Ctrl+S opens the Keyboard Shortcuts editor; Save All
	// lives on Ctrl+K S.
	if res.Action.Title != "Keyboard Shortcuts" {
		t.Errorf("expected Keyboard Shortcuts, got %q", res.Action.Title)
	}
	if !res.IsSequence {
		t.Error("expected IsSequence=true")
	}
	if res.Confidence != semantic.ConfidenceProfile {
		t.Errorf("expected profile confidence, got %s", res.Confidence)
	}

	// Case-insensitive primary keys.
	if res, ok := m.MatchSequence(ctx, hook.ModCtrl, "k", hook.ModCtrl, "u"); !ok || res.Action.Title != "Remove Line Comment" {
		t.Errorf("expected Remove Line Comment, got %q (%t)", res.Action.Title, ok)
	}

	// Unknown pair.
	if _, ok := m.MatchSequence(ctx, hook.ModCtrl, "K", hook.ModCtrl, "Q"); ok {
		t.Error("Ctrl+K Ctrl+Q should not resolve")
	}

	// Visual Studio shares the same chord family but a different meaning.
	vs := app("devenv.exe", "Solution - Microsoft Visual Studio")
	if res, ok := m.MatchSequence(vs, hook.ModCtrl, "K", hook.ModCtrl, "D"); !ok || res.Action.Title != "Format Document" {
		t.Errorf("expected VS Code Format Document mapping in Visual Studio, got %q (%t)", res.Action.Title, ok)
	}
}

func TestUserOverrides(t *testing.T) {
	m := semantic.NewMatcher()
	m.SetOverrides([]semantic.Override{
		{
			App:    "devcpp.exe",
			Mods:   hook.ModCtrl | hook.ModAlt,
			Key:    "B",
			Action: semantic.ActionInfo{Title: "Build My Project"},
		},
		{
			App:    "",
			Mods:   hook.ModCtrl | hook.ModShift,
			Key:    "Q",
			Action: semantic.ActionInfo{Title: "Quick Note"},
		},
	})

	res := m.MatchAction(app("devcpp.exe", "main.cpp - Dev-C++"), hook.ModCtrl|hook.ModAlt, "B")
	if res.Action.Title != "Build My Project" {
		t.Fatalf("expected user override, got %q", res.Action.Title)
	}
	if res.Confidence != semantic.ConfidenceUser {
		t.Fatalf("expected user confidence, got %s", res.Confidence)
	}

	// Global override applies everywhere.
	if res := m.MatchAction(app("chrome.exe", "Example"), hook.ModCtrl|hook.ModShift, "Q"); res.Action.Title != "Quick Note" {
		t.Fatalf("expected global override, got %q", res.Action.Title)
	}

	// Existing built-in mappings still work outside the override scope.
	if res := m.MatchAction(app("devcpp.exe", "main.cpp - Dev-C++"), hook.ModCtrl, "C"); res.Action.Title != "Copy" {
		t.Fatalf("expected built-in Copy, got %q", res.Action.Title)
	}
}

func TestMatchUniversalIgnoresProfiles(t *testing.T) {
	m := semantic.NewMatcher()
	m.SetOverrides([]semantic.Override{
		{
			App:    "devcpp.exe",
			Mods:   hook.ModCtrl,
			Key:    "C",
			Action: semantic.ActionInfo{Title: "Compile Copy"},
		},
		{
			App:    "editor",
			Mods:   hook.ModCtrl | hook.ModAlt,
			Key:    "H",
			Action: semantic.ActionInfo{Title: "Editor Thing"},
		},
		{
			Mods:   hook.ModCtrl | hook.ModAlt,
			Key:    "G",
			Action: semantic.ActionInfo{Title: "Global Thing"},
		},
	})

	// App- and category-scoped overrides must not leak into universal mode.
	if res := m.MatchUniversal(hook.ModCtrl, "C"); res.Action.Title != "Copy" || res.Confidence != semantic.ConfidenceUniversal {
		t.Errorf("expected universal Copy, got %q (%s)", res.Action.Title, res.Confidence)
	}
	if res := m.MatchUniversal(hook.ModCtrl|hook.ModAlt, "H"); res.Matched() {
		t.Errorf("category override must not apply in universal mode, got %q", res.Action.Title)
	}

	// Global overrides do apply everywhere.
	res := m.MatchUniversal(hook.ModCtrl|hook.ModAlt, "G")
	if res.Action.Title != "Global Thing" || res.Confidence != semantic.ConfidenceUser {
		t.Errorf("expected global user override, got %q (%s)", res.Action.Title, res.Confidence)
	}
	if res.AppName != "" {
		t.Errorf("universal mode must not attach an application name, got %q", res.AppName)
	}

	// App-only shortcuts have no universal meaning.
	if res := m.MatchUniversal(hook.ModNone, "F9"); res.Matched() {
		t.Errorf("app-only F9 must not resolve without context, got %q", res.Action.Title)
	}
}

func TestParseOverride(t *testing.T) {
	o, err := semantic.ParseOverride("devcpp.exe", "Ctrl+Shift+F9", "Compile All", "")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if o.App != "devcpp.exe" || o.Key != "F9" {
		t.Fatalf("unexpected override: %+v", o)
	}
	if o.Mods != hook.ModCtrl|hook.ModShift {
		t.Fatalf("unexpected modifiers: %v", o.Mods)
	}
	if _, err := semantic.ParseOverride("", "", "x", ""); err == nil {
		t.Fatal("empty chord should fail")
	}
	if _, err := semantic.ParseOverride("", "Ctrl+", "", ""); err == nil {
		t.Fatal("empty label should fail")
	}
}

func TestIsSensitive(t *testing.T) {
	m := semantic.NewMatcher()

	sensitiveContexts := []context.AppContext{
		app("1Password.exe", "1Password - Password Manager"),
		app("chrome.exe", "Bitwarden Vault"),
		app("KeePass.exe", "KeePass Password Safe"),
		app("chrome.exe", "Login - Microsoft Account"),
		app("chrome.exe", "Sign in to your account - Google Chrome"),
		app("chrome.exe", "Bank of America - Sign In"),
		app("chrome.exe", "PayPal: Summary"),
		app("LogonUI.exe", "Windows Security - Enter your PIN"),
		app("UserAccountControl.exe", "User Account Control"),
		app("firefox.exe", "Private Browsing - Firefox"),
		app("chrome.exe", "Incognito - Chrome"),
	}

	for _, ctx := range sensitiveContexts {
		// Each context has a distinct Key, so the matcher recomputes the flag.
		if !m.IsSensitive(ctx, "a") {
			t.Errorf("expected context %q to be marked sensitive", ctx.Title)
		}
	}

	// Non-sensitive context.
	m.Reset()
	if m.IsSensitive(app("code.exe", "Visual Studio Code - main.go"), "a") {
		t.Errorf("normal VS Code keystroke should not be sensitive")
	}

	// Direct sensitive keyName.
	if !m.IsSensitive(app("notepad.exe", "Notepad"), "password") {
		t.Errorf("direct 'password' keyName should be marked sensitive")
	}

	// Consecutive pattern typing: p a s s w o r d
	m.Reset()
	for _, k := range []string{"p", "a", "s", "s", "w", "o", "r"} {
		m.IsSensitive(app("notepad.exe", "Notepad"), k)
	}
	if !m.IsSensitive(app("notepad.exe", "Notepad"), "d") {
		t.Errorf("completing trigger word 'password' should trigger sensitive mode")
	}
	if !m.IsSensitive(app("notepad.exe", "Notepad"), "x") {
		t.Errorf("keystroke following sensitive trigger should remain sensitive")
	}

	// Enter resets the sensitive streak.
	m.IsSensitive(app("notepad.exe", "Notepad"), "Enter")
	if m.IsSensitive(app("notepad.exe", "Notepad"), "z") {
		t.Errorf("keystroke after Enter should no longer be sensitive")
	}
}

func TestConcurrentMatcher(t *testing.T) {
	m := semantic.NewMatcher()
	appCtx := app("code.exe", "main.go - Visual Studio Code")
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = m.MatchAction(appCtx, hook.ModCtrl, "C")
				_, _ = m.MatchSequence(appCtx, hook.ModCtrl, "K", hook.ModCtrl, "S")
				_ = m.IsSensitive(app("notepad.exe", "Notepad"), "a")
			}
		}()
	}
	wg.Wait()
}
