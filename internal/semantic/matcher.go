package semantic

import (
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/yviscool/cue/internal/hook"
)

// comboKey represents a unique modifier combination and normalized key name.
type comboKey struct {
	mods hook.Modifiers
	key  string
}

// DefaultMatcher implements Matcher with a curated shortcut dictionary
// and context/keystroke-based privacy detection.
type DefaultMatcher struct {
	mu           sync.Mutex
	shortcuts    map[comboKey]ActionInfo
	appShortcuts map[string]map[comboKey]ActionInfo

	// Keystroke history and privacy state
	history       []rune
	sensitiveLeft int
	lastInputTime time.Time
}

var _ Matcher = (*DefaultMatcher)(nil)

// NewMatcher creates a DefaultMatcher with standard curated shortcuts.
func NewMatcher() *DefaultMatcher {
	m := &DefaultMatcher{
		shortcuts:    make(map[comboKey]ActionInfo),
		appShortcuts: make(map[string]map[comboKey]ActionInfo),
		history:      make([]rune, 0, 64),
	}
	m.initShortcuts()
	return m
}

// NewDefaultMatcher is an alias for NewMatcher.
func NewDefaultMatcher() *DefaultMatcher {
	return NewMatcher()
}

// NormalizeKey standardizes key representations (e.g. "+"/"=", "ArrowUp"/"Up").
func NormalizeKey(keyName string) string {
	k := strings.ToUpper(strings.TrimSpace(keyName))
	switch k {
	case "+", "=", "PLUS", "ADD":
		return "+"
	case "-", "_", "MINUS", "SUBTRACT":
		return "-"
	case "`", "~", "BACKQUOTE", "GRAVE", "TILDE":
		return "`"
	case "/", "SLASH":
		return "/"
	case "ARROWUP", "UPARROW":
		return "UP"
	case "ARROWDOWN", "DOWNARROW":
		return "DOWN"
	case "ARROWLEFT", "LEFTARROW":
		return "LEFT"
	case "ARROWRIGHT", "RIGHTARROW":
		return "RIGHT"
	case "ESCAPE":
		return "ESC"
	case "RETURN":
		return "ENTER"
	case "SPACEBAR", " ":
		return "SPACE"
	case "BKSP", "BS":
		return "BACKSPACE"
	case "DEL":
		return "DELETE"
	case "LEFT CLICK", "MOUSE LEFT", "LBUTTON":
		return "LEFT CLICK"
	case "RIGHT CLICK", "MOUSE RIGHT", "RBUTTON":
		return "RIGHT CLICK"
	case "MIDDLE CLICK", "MOUSE MIDDLE", "MBUTTON":
		return "MIDDLE CLICK"
	case "SCROLL UP", "SCROLLUP", "WHEEL UP", "WHEELUP":
		return "SCROLL UP"
	case "SCROLL DOWN", "SCROLLDOWN", "WHEEL DOWN", "WHEELDOWN":
		return "SCROLL DOWN"
	case "PGUP":
		return "PAGEUP"
	case "PGDN":
		return "PAGEDOWN"
	case "INS":
		return "INSERT"
	default:
		return k
	}
}

func (m *DefaultMatcher) register(mods hook.Modifiers, key string, action ActionInfo) {
	norm := NormalizeKey(key)
	m.shortcuts[comboKey{mods: mods, key: norm}] = action
}

func (m *DefaultMatcher) registerApp(appKey string, mods hook.Modifiers, key string, action ActionInfo) {
	appKey = strings.ToLower(appKey)
	if _, ok := m.appShortcuts[appKey]; !ok {
		m.appShortcuts[appKey] = make(map[comboKey]ActionInfo)
	}
	norm := NormalizeKey(key)
	m.appShortcuts[appKey][comboKey{mods: mods, key: norm}] = action
}

func (m *DefaultMatcher) initShortcuts() {
	// ==========================================
	// Universal / Windows Shortcuts
	// ==========================================
	// Edit operations
	m.register(hook.ModCtrl, "C", ActionInfo{Title: "Copy", Category: "Edit", Description: "Copy selection to clipboard"})
	m.register(hook.ModCtrl, "X", ActionInfo{Title: "Cut", Category: "Edit", Description: "Cut selection to clipboard"})
	m.register(hook.ModCtrl, "V", ActionInfo{Title: "Paste", Category: "Edit", Description: "Paste clipboard contents"})
	m.register(hook.ModCtrl, "Z", ActionInfo{Title: "Undo", Category: "Edit", Description: "Undo previous action"})
	m.register(hook.ModCtrl, "Y", ActionInfo{Title: "Redo", Category: "Edit", Description: "Redo undone action"})
	m.register(hook.ModCtrl|hook.ModShift, "Z", ActionInfo{Title: "Redo", Category: "Edit", Description: "Redo undone action"})
	m.register(hook.ModCtrl, "A", ActionInfo{Title: "Select All", Category: "Edit", Description: "Select all items or text"})
	m.register(hook.ModCtrl, "S", ActionInfo{Title: "Save", Category: "File", Description: "Save current document or state"})
	m.register(hook.ModCtrl, "F", ActionInfo{Title: "Find", Category: "Search", Description: "Open find bar or search query"})
	m.register(hook.ModCtrl, "H", ActionInfo{Title: "Replace", Category: "Search", Description: "Open find and replace"})
	m.register(hook.ModCtrl, "P", ActionInfo{Title: "Quick Open", Category: "Navigation", Description: "Quickly open file or jump to target"})
	m.register(hook.ModCtrl|hook.ModShift, "P", ActionInfo{Title: "Command Palette", Category: "System", Description: "Open command palette"})

	// Window & System management
	m.register(hook.ModAlt, "TAB", ActionInfo{Title: "Switch App", Category: "Window", Description: "Switch between active applications"})
	m.register(hook.ModAlt|hook.ModShift, "TAB", ActionInfo{Title: "Switch App", Category: "Window", Description: "Switch between active applications (reverse)"})
	m.register(hook.ModSuper, "D", ActionInfo{Title: "Show Desktop", Category: "System", Description: "Minimize all windows to display desktop"})
	m.register(hook.ModSuper, "E", ActionInfo{Title: "File Explorer", Category: "System", Description: "Open Windows File Explorer"})
	m.register(hook.ModSuper, "L", ActionInfo{Title: "Lock PC", Category: "System", Description: "Lock Windows workstation"})
	m.register(hook.ModSuper, "V", ActionInfo{Title: "Clipboard History", Category: "System", Description: "Open Windows clipboard history popup"})
	m.register(hook.ModSuper|hook.ModShift, "S", ActionInfo{Title: "Snipping Tool", Category: "System", Description: "Capture screen snip or screenshot"})
	m.register(hook.ModAlt, "F4", ActionInfo{Title: "Close Window", Category: "Window", Description: "Close active window or application"})

	// Tabs & Navigation
	m.register(hook.ModCtrl, "W", ActionInfo{Title: "Close Tab", Category: "Tab", Description: "Close active tab or document"})
	m.register(hook.ModCtrl, "T", ActionInfo{Title: "New Tab", Category: "Tab", Description: "Open new tab"})
	m.register(hook.ModCtrl|hook.ModShift, "T", ActionInfo{Title: "Reopen Closed Tab", Category: "Tab", Description: "Restore recently closed tab"})
	m.register(hook.ModCtrl, "R", ActionInfo{Title: "Refresh", Category: "Navigation", Description: "Reload active page or view"})
	m.register(hook.ModNone, "F5", ActionInfo{Title: "Refresh", Category: "Navigation", Description: "Reload active page or view"})
	m.register(hook.ModCtrl, "F5", ActionInfo{Title: "Refresh", Category: "Navigation", Description: "Hard reload page, ignoring cache"})

	// Zoom Controls (handling +, =, and shift variations, as well as mouse wheel)
	zoomIn := ActionInfo{Title: "Zoom In", Category: "View", Description: "Increase magnification"}
	m.register(hook.ModCtrl, "+", zoomIn)
	m.register(hook.ModCtrl, "=", zoomIn)
	m.register(hook.ModCtrl|hook.ModShift, "+", zoomIn)
	m.register(hook.ModCtrl|hook.ModShift, "=", zoomIn)
	m.register(hook.ModCtrl, "SCROLL UP", zoomIn)

	zoomOut := ActionInfo{Title: "Zoom Out", Category: "View", Description: "Decrease magnification"}
	m.register(hook.ModCtrl, "-", zoomOut)
	m.register(hook.ModCtrl|hook.ModShift, "-", zoomOut)
	m.register(hook.ModCtrl, "SCROLL DOWN", zoomOut)

	m.register(hook.ModCtrl, "0", ActionInfo{Title: "Reset Zoom", Category: "View", Description: "Reset magnification to default (100%)"})

	// Additional Universal Shortcuts
	m.register(hook.ModCtrl, "N", ActionInfo{Title: "New Window", Category: "File", Description: "Open a new window or document"})
	m.register(hook.ModCtrl, "O", ActionInfo{Title: "Open File", Category: "File", Description: "Open an existing file"})
	m.register(hook.ModCtrl|hook.ModShift, "S", ActionInfo{Title: "Save As", Category: "File", Description: "Save current document under a new name"})
	m.register(hook.ModSuper, "R", ActionInfo{Title: "Run Dialog", Category: "System", Description: "Open Windows Run prompt"})
	m.register(hook.ModSuper, "I", ActionInfo{Title: "Settings", Category: "System", Description: "Open Windows Settings"})
	m.register(hook.ModSuper, "A", ActionInfo{Title: "Action Center", Category: "System", Description: "Open Windows Quick Settings and Notifications"})
	m.register(hook.ModCtrl, "TAB", ActionInfo{Title: "Next Tab", Category: "Tab", Description: "Switch to next tab"})
	m.register(hook.ModCtrl|hook.ModShift, "TAB", ActionInfo{Title: "Previous Tab", Category: "Tab", Description: "Switch to previous tab"})

	// ==========================================
	// Code Editors (VS Code / JetBrains / Zed / Sublime)
	// ==========================================
	m.register(hook.ModCtrl, "B", ActionInfo{Title: "Toggle Sidebar", Category: "Editor", Description: "Show or hide editor primary sidebar"})
	m.register(hook.ModCtrl, "`", ActionInfo{Title: "Toggle Terminal", Category: "Editor", Description: "Show or hide integrated terminal"})
	m.register(hook.ModCtrl, "/", ActionInfo{Title: "Toggle Comment", Category: "Editor", Description: "Comment or uncomment line or selection"})
	m.register(hook.ModAlt, "UP", ActionInfo{Title: "Move Line Up", Category: "Editor", Description: "Move active line or selection up"})
	m.register(hook.ModAlt, "DOWN", ActionInfo{Title: "Move Line Down", Category: "Editor", Description: "Move active line or selection down"})
	m.register(hook.ModShift|hook.ModAlt, "DOWN", ActionInfo{Title: "Duplicate Line Down", Category: "Editor", Description: "Duplicate active line or selection downwards"})
	m.register(hook.ModShift|hook.ModAlt, "UP", ActionInfo{Title: "Duplicate Line Up", Category: "Editor", Description: "Duplicate active line or selection upwards"})
	m.register(hook.ModCtrl, "D", ActionInfo{Title: "Add Selection to Next Find Match", Category: "Editor", Description: "Add next match of selected word to multi-cursor"})
	m.register(hook.ModNone, "F12", ActionInfo{Title: "Go to Definition", Category: "Navigation", Description: "Jump to symbol declaration or definition"})
	m.register(hook.ModShift, "F12", ActionInfo{Title: "Find References", Category: "Search", Description: "Find all usages of symbol across workspace"})
	m.register(hook.ModNone, "F2", ActionInfo{Title: "Rename Symbol", Category: "Refactor", Description: "Rename symbol across workspace"})
	m.register(hook.ModCtrl|hook.ModShift, "O", ActionInfo{Title: "Go to Symbol", Category: "Navigation", Description: "Navigate to symbol in current file"})
	m.register(hook.ModCtrl, "G", ActionInfo{Title: "Go to Line", Category: "Navigation", Description: "Jump to specified line number"})

	// Additional Editor Shortcuts
	m.register(hook.ModCtrl|hook.ModShift, "F", ActionInfo{Title: "Find in Files", Category: "Search", Description: "Global search across entire workspace"})
	m.register(hook.ModCtrl|hook.ModShift, "H", ActionInfo{Title: "Replace in Files", Category: "Search", Description: "Global replace across entire workspace"})
	m.register(hook.ModCtrl|hook.ModShift, "K", ActionInfo{Title: "Delete Line", Category: "Editor", Description: "Delete active line"})
	m.register(hook.ModCtrl, "J", ActionInfo{Title: "Toggle Panel", Category: "Editor", Description: "Show or hide bottom panel"})
	m.register(hook.ModCtrl, "SPACE", ActionInfo{Title: "Trigger Suggestion", Category: "Editor", Description: "Open IntelliSense autocomplete list"})

	// ==========================================
	// App-Specific Context Overrides (Browser & Terminal)
	// ==========================================
	// Web Browsers (Chrome, Edge, Firefox, Brave)
	m.registerApp("browser", hook.ModCtrl, "D", ActionInfo{Title: "Bookmark Tab", Category: "Browser", Description: "Bookmark current page"})
	m.registerApp("browser", hook.ModCtrl|hook.ModShift, "B", ActionInfo{Title: "Toggle Bookmarks Bar", Category: "Browser", Description: "Show or hide bookmarks bar"})
	m.registerApp("browser", hook.ModCtrl, "J", ActionInfo{Title: "Downloads", Category: "Browser", Description: "Open browser downloads tab"})
	m.registerApp("browser", hook.ModCtrl, "H", ActionInfo{Title: "History", Category: "Browser", Description: "Open browsing history"})
	m.registerApp("browser", hook.ModNone, "F12", ActionInfo{Title: "Developer Tools", Category: "Browser", Description: "Open Developer Tools console"})
	m.registerApp("browser", hook.ModCtrl|hook.ModShift, "I", ActionInfo{Title: "Developer Tools", Category: "Browser", Description: "Open Developer Tools console"})

	// Terminal / Command Line
	m.registerApp("terminal", hook.ModCtrl, "C", ActionInfo{Title: "Interrupt Process", Category: "Terminal", Description: "Send SIGINT to interrupt active process"})
	m.registerApp("terminal", hook.ModCtrl, "D", ActionInfo{Title: "Exit / EOF", Category: "Terminal", Description: "Send EOF or close session"})
	m.registerApp("terminal", hook.ModCtrl, "L", ActionInfo{Title: "Clear Terminal", Category: "Terminal", Description: "Clear terminal screen buffer"})
}

func matchAppCategory(appContext string) string {
	app := strings.ToLower(appContext)
	switch {
	case strings.Contains(app, "chrome"),
		strings.Contains(app, "edge"),
		strings.Contains(app, "firefox"),
		strings.Contains(app, "brave"),
		strings.Contains(app, "opera"),
		strings.Contains(app, "safari"),
		strings.Contains(app, "vivaldi"):
		return "browser"
	case strings.Contains(app, "cmd.exe"),
		strings.Contains(app, "powershell"),
		strings.Contains(app, "windowsterminal"),
		strings.Contains(app, "wt.exe"),
		strings.Contains(app, "bash"),
		strings.Contains(app, "alacritty"),
		strings.Contains(app, "wezterm"),
		strings.Contains(app, "kitty"):
		return "terminal"
	default:
		return ""
	}
}

// MatchAction returns an ActionInfo if the combo corresponds to a known shortcut.
func (m *DefaultMatcher) MatchAction(appContext string, mods hook.Modifiers, keyName string) (ActionInfo, bool) {
	normKey := NormalizeKey(keyName)
	if normKey == "" {
		return ActionInfo{}, false
	}

	key := comboKey{mods: mods, key: normKey}

	// 1. Check app-specific overrides if appContext matches known app category
	if appContext != "" {
		category := matchAppCategory(appContext)
		if category != "" {
			if appMap, ok := m.appShortcuts[category]; ok {
				if action, found := appMap[key]; found {
					return action, true
				}
			}
		}
	}

	// 2. Check universal / default dictionary
	if action, found := m.shortcuts[key]; found {
		return action, true
	}

	return ActionInfo{}, false
}

// sensitiveContextKeywords contains substrings indicative of credentials, logins, or finance.
var sensitiveContextKeywords = []string{
	// Password managers
	"1password", "bitwarden", "keepass", "lastpass", "dashlane", "enpass",
	"nordpass", "keeper", "roboform", "authenticator",

	// Authentication & logins
	"password", "passwd", "passcode", "pin code", "pincode", "master password",
	"sign in", "signin", "log in", "login", "authenticate", "authentication",
	"credential", "two-factor", "2fa", "mfa", "verification code", "one-time password", "otp",

	// Financial / Banking / Payment
	"bank", "banking", "checkout", "paypal", "credit card", "debit card", "cvv",
	"wallet", "binance", "coinbase", "metamask",

	// Browsing privacy modes
	"incognito", "inprivate", "private browsing", "private window",

	// OS security prompts
	"user account control", "windows security", "credential manager", "pin entry", "sudo password",
}

// sensitiveKeyTriggers contains trigger sequences in typed history.
var sensitiveKeyTriggers = []string{
	"password", "passwd", "pwd:", "pwd ", "pin:", "secret", "token",
	"api_key", "apikey", "bearer ", "private_key", "sudo ", "cvv",
}

// IsSensitive returns true if current window title / context indicates
// password/credentials/banking or if the user has typed sensitive patterns.
func (m *DefaultMatcher) IsSensitive(appContext string, keyName string) bool {
	// 1. Check appContext (window title / process name)
	appLower := strings.ToLower(appContext)
	for _, kw := range sensitiveContextKeywords {
		if strings.Contains(appLower, kw) {
			return true
		}
	}

	// 2. Check keyName itself
	keyLower := strings.ToLower(keyName)
	if keyLower == "password" || keyLower == "secret" || keyLower == "token" ||
		keyLower == "credential" || keyLower == "pin" || strings.Contains(keyLower, "password") {
		return true
	}

	// 3. Keystroke pattern tracking
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	// Reset stale history after 5 seconds of inactivity
	if now.Sub(m.lastInputTime) > 5*time.Second {
		m.history = m.history[:0]
	}
	m.lastInputTime = now

	// If currently within a sensitive stream (e.g. following a trigger word)
	if m.sensitiveLeft > 0 {
		m.sensitiveLeft--
		if norm := NormalizeKey(keyName); norm == "ENTER" || norm == "ESC" {
			m.sensitiveLeft = 0
			m.history = m.history[:0]
		}
		return true
	}

	// Handle Backspace
	if norm := NormalizeKey(keyName); norm == "BACKSPACE" {
		if len(m.history) > 0 {
			m.history = m.history[:len(m.history)-1]
		}
		return false
	}

	// Handle Enter / Esc
	if norm := NormalizeKey(keyName); norm == "ENTER" || norm == "ESC" {
		m.history = m.history[:0]
		return false
	}

	// Record printable keystroke
	runes := []rune(keyName)
	if len(runes) == 1 {
		r := unicode.ToLower(runes[0])
		m.history = append(m.history, r)
		if len(m.history) > 64 {
			m.history = m.history[len(m.history)-64:]
		}

		histStr := string(m.history)
		for _, trig := range sensitiveKeyTriggers {
			if strings.HasSuffix(histStr, trig) || strings.Contains(histStr, trig) {
				m.sensitiveLeft = 32 // Shield the next 32 keystrokes
				return true
			}
		}

		// Detect consecutive digits (like credit cards / SSN / PINs)
		if hasConsecutiveDigits(histStr, 8) {
			m.sensitiveLeft = 16
			return true
		}
	}

	return false
}

// Reset clears the internal history and privacy state.
func (m *DefaultMatcher) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = m.history[:0]
	m.sensitiveLeft = 0
}

func hasConsecutiveDigits(s string, minDigits int) bool {
	count := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			count++
			if count >= minDigits {
				return true
			}
		} else if r != '-' && r != ' ' { // ignore separators
			count = 0
		}
	}
	return false
}
