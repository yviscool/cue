package semantic_test

import (
	"sync"
	"testing"

	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/semantic"
)

func TestUniversalShortcuts(t *testing.T) {
	m := semantic.NewMatcher()

	tests := []struct {
		name          string
		appContext    string
		mods          hook.Modifiers
		keyName       string
		expectedTitle string
	}{
		{"Ctrl+C", "", hook.ModCtrl, "C", "Copy"},
		{"Ctrl+c lowercase", "", hook.ModCtrl, "c", "Copy"},
		{"Ctrl+X", "", hook.ModCtrl, "X", "Cut"},
		{"Ctrl+V", "", hook.ModCtrl, "V", "Paste"},
		{"Ctrl+Z", "", hook.ModCtrl, "Z", "Undo"},
		{"Ctrl+Y", "", hook.ModCtrl, "Y", "Redo"},
		{"Ctrl+Shift+Z", "", hook.ModCtrl | hook.ModShift, "Z", "Redo"},
		{"Ctrl+A", "", hook.ModCtrl, "A", "Select All"},
		{"Ctrl+S", "", hook.ModCtrl, "S", "Save"},
		{"Ctrl+F", "", hook.ModCtrl, "F", "Find"},
		{"Ctrl+H", "", hook.ModCtrl, "H", "Replace"},
		{"Ctrl+P", "", hook.ModCtrl, "P", "Quick Open"},
		{"Ctrl+Shift+P", "", hook.ModCtrl | hook.ModShift, "P", "Command Palette"},
		{"Alt+Tab", "", hook.ModAlt, "Tab", "Switch App"},
		{"Win+D", "", hook.ModSuper, "D", "Show Desktop"},
		{"Win+E", "", hook.ModSuper, "E", "File Explorer"},
		{"Win+L", "", hook.ModSuper, "L", "Lock PC"},
		{"Win+V", "", hook.ModSuper, "V", "Clipboard History"},
		{"Win+Shift+S", "", hook.ModSuper | hook.ModShift, "S", "Snipping Tool"},
		{"Alt+F4", "", hook.ModAlt, "F4", "Close Window"},
		{"Ctrl+W", "", hook.ModCtrl, "W", "Close Tab"},
		{"Ctrl+T", "", hook.ModCtrl, "T", "New Tab"},
		{"Ctrl+Shift+T", "", hook.ModCtrl | hook.ModShift, "T", "Reopen Closed Tab"},
		{"Ctrl+R", "", hook.ModCtrl, "R", "Refresh"},
		{"F5", "", hook.ModNone, "F5", "Refresh"},
		{"Ctrl+Plus", "", hook.ModCtrl, "+", "Zoom In"},
		{"Ctrl+=", "", hook.ModCtrl, "=", "Zoom In"},
		{"Ctrl+Shift+Plus", "", hook.ModCtrl | hook.ModShift, "+", "Zoom In"},
		{"Ctrl+Minus", "", hook.ModCtrl, "-", "Zoom Out"},
		{"Ctrl+0", "", hook.ModCtrl, "0", "Reset Zoom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, ok := m.MatchAction(tt.appContext, tt.mods, tt.keyName)
			if !ok {
				t.Fatalf("expected shortcut %s to match, but got no match", tt.name)
			}
			if action.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, action.Title)
			}
		})
	}
}

func TestCodeEditorShortcuts(t *testing.T) {
	m := semantic.NewMatcher()

	tests := []struct {
		name          string
		appContext    string
		mods          hook.Modifiers
		keyName       string
		expectedTitle string
	}{
		{"Ctrl+B", "Code.exe", hook.ModCtrl, "B", "Toggle Sidebar"},
		{"Ctrl+`", "Code.exe", hook.ModCtrl, "`", "Toggle Terminal"},
		{"Ctrl+Grave", "Code.exe", hook.ModCtrl, "Grave", "Toggle Terminal"},
		{"Ctrl+/", "Code.exe", hook.ModCtrl, "/", "Toggle Comment"},
		{"Alt+Up", "Code.exe", hook.ModAlt, "Up", "Move Line Up"},
		{"Alt+ArrowUp", "Code.exe", hook.ModAlt, "ArrowUp", "Move Line Up"},
		{"Alt+Down", "Code.exe", hook.ModAlt, "Down", "Move Line Down"},
		{"Alt+ArrowDown", "Code.exe", hook.ModAlt, "ArrowDown", "Move Line Down"},
		{"Shift+Alt+Down", "Code.exe", hook.ModShift | hook.ModAlt, "Down", "Duplicate Line Down"},
		{"Shift+Alt+Up", "Code.exe", hook.ModShift | hook.ModAlt, "Up", "Duplicate Line Up"},
		{"Ctrl+D", "Code.exe", hook.ModCtrl, "D", "Add Selection to Next Find Match"},
		{"F12", "Code.exe", hook.ModNone, "F12", "Go to Definition"},
		{"Shift+F12", "Code.exe", hook.ModShift, "F12", "Find References"},
		{"F2", "Code.exe", hook.ModNone, "F2", "Rename Symbol"},
		{"Ctrl+Shift+O", "Code.exe", hook.ModCtrl | hook.ModShift, "O", "Go to Symbol"},
		{"Ctrl+G", "Code.exe", hook.ModCtrl, "G", "Go to Line"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, ok := m.MatchAction(tt.appContext, tt.mods, tt.keyName)
			if !ok {
				t.Fatalf("expected shortcut %s to match, but got no match", tt.name)
			}
			if action.Title != tt.expectedTitle {
				t.Errorf("expected title %q, got %q", tt.expectedTitle, action.Title)
			}
		})
	}
}

func TestAppOverrides(t *testing.T) {
	m := semantic.NewMatcher()

	// In browser, Ctrl+D is "Bookmark Tab"
	actionBrowser, ok := m.MatchAction("Google Chrome", hook.ModCtrl, "D")
	if !ok || actionBrowser.Title != "Bookmark Tab" {
		t.Errorf("expected 'Bookmark Tab' in browser, got %v (%t)", actionBrowser.Title, ok)
	}

	// In terminal, Ctrl+C is "Interrupt Process"
	actionTerm, ok := m.MatchAction("WindowsTerminal.exe", hook.ModCtrl, "C")
	if !ok || actionTerm.Title != "Interrupt Process" {
		t.Errorf("expected 'Interrupt Process' in terminal, got %v (%t)", actionTerm.Title, ok)
	}

	// In editor or default, Ctrl+C is "Copy"
	actionGeneral, ok := m.MatchAction("Code.exe", hook.ModCtrl, "C")
	if !ok || actionGeneral.Title != "Copy" {
		t.Errorf("expected 'Copy' in editor, got %v (%t)", actionGeneral.Title, ok)
	}
}

func TestIsSensitive(t *testing.T) {
	m := semantic.NewMatcher()

	// Context indicating password / auth / banking
	sensitiveContexts := []string{
		"1Password - Password Manager",
		"Bitwarden Vault",
		"KeePass Password Safe",
		"Login - Microsoft Account",
		"Sign in to your account - Google Chrome",
		"Bank of America - Sign In",
		"PayPal: Summary",
		"Windows Security - Enter your PIN",
		"User Account Control",
		"Private Browsing - Firefox",
		"Incognito - Chrome",
	}

	for _, ctx := range sensitiveContexts {
		if !m.IsSensitive(ctx, "a") {
			t.Errorf("expected context %q to be marked sensitive", ctx)
		}
	}

	// Non-sensitive context
	m.Reset()
	if m.IsSensitive("Visual Studio Code - main.go", "a") {
		t.Errorf("normal VS Code keystroke should not be sensitive")
	}

	// Direct sensitive keyName
	if !m.IsSensitive("Notepad", "password") {
		t.Errorf("direct 'password' keyName should be marked sensitive")
	}

	// Consecutive pattern typing: "p", "a", "s", "s", "w", "o", "r", "d"
	m.Reset()
	keys := []string{"p", "a", "s", "s", "w", "o", "r"}
	for _, k := range keys {
		m.IsSensitive("Notepad", k)
	}
	// The 8th key completes "password"
	if !m.IsSensitive("Notepad", "d") {
		t.Errorf("completing trigger word 'password' should trigger sensitive mode")
	}
	// Subsequent typed characters should now be masked as private
	if !m.IsSensitive("Notepad", "x") {
		t.Errorf("keystroke following sensitive trigger should remain sensitive")
	}
	if !m.IsSensitive("Notepad", "y") {
		t.Errorf("keystroke following sensitive trigger should remain sensitive")
	}

	// Enter resets sensitive streak
	m.IsSensitive("Notepad", "Enter")
	if m.IsSensitive("Notepad", "z") {
		t.Errorf("keystroke after Enter should no longer be sensitive")
	}
}

func TestConcurrentMatcher(t *testing.T) {
	m := semantic.NewMatcher()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = m.MatchAction("Code.exe", hook.ModCtrl, "C")
				_ = m.IsSensitive("Notepad", "a")
			}
		}(i)
	}
	wg.Wait()
}
