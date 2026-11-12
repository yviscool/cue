package hud

import (
	"runtime"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/config"
)

var (
	// Modifier SVG Icons optically centered and scaled to text cap-height
	iconKeyCtrl  = strokeIcon(`<path d="M6 14.5L12 7.5L18 14.5"/>`)
	iconKeyShift = strokeIcon(`<path d="M12 4L5 11H9V19H15V11H19L12 4Z"/>`)
	iconKeyAlt   = strokeIcon(`<path d="M3 4h6l6 16h6M14 4h7"/>`)
	iconKeyWin   = strokeIcon(`<rect width="17" height="17" x="3.5" y="3.5" rx="1.5"/><path d="M3.5 12h17"/><path d="M12 3.5v17"/>`)
	iconKeyCmd   = strokeIcon(`<path d="M15 6v12a3 3 0 1 0 3-3H6a3 3 0 1 0 3 3V6a3 3 0 1 0-3 3h12a3 3 0 1 0-3-3"/>`)

	// Navigation & Special Key Icons matching keyviz-source
	iconKeyReturn    = strokeIcon(`<path d="M9 10L5 14L9 18"/><path d="M5 14H15C17.5 14 19.5 12 19.5 9.5V6"/>`)
	iconKeyBackspace = strokeIcon(`<path d="M9 5a2 2 0 0 0-1.4.6L2 12l5.6 6.4A2 2 0 0 0 9 19h11a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2H9z"/><path d="m11 9 6 6m0-6-6 6"/>`)
	iconKeyTab       = strokeIcon(`<path d="M4 8h15M15 4l4 4-4 4M20 16H5M9 12l-4 4 4 4"/>`)
	iconKeySpace     = strokeIcon(`<path d="M4 9v6a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1V9"/>`)
	iconKeyEscape    = strokeIcon(`<circle cx="12" cy="12" r="9"/><path d="M9 9l6 6M15 9l-6 6"/>`)
	iconKeyCapsLock  = strokeIcon(`<path d="M12 4L6 10H10V16H14V10H18L12 4Z"/><line x1="6" x2="18" y1="19" y2="19"/>`)

	// Arrow Key Icons matching Lucide
	iconKeyArrowUp    = strokeIcon(`<path d="m5 12 7-7 7 7"/><path d="M12 19V5"/>`)
	iconKeyArrowDown  = strokeIcon(`<path d="M12 5v14"/><path d="m19 12-7 7-7-7"/>`)
	iconKeyArrowLeft  = strokeIcon(`<path d="m12 19-7-7 7-7"/><path d="M19 12H5"/>`)
	iconKeyArrowRight = strokeIcon(`<path d="M5 12h14"/><path d="m12 5 7 7-7 7"/>`)

	// Navigation
	iconKeyPageUp   = strokeIcon(`<path d="m18 15-6-6-6 6"/><path d="M6 5h12"/>`)
	iconKeyPageDown = strokeIcon(`<path d="m6 9 6 6 6-6"/><path d="M6 19h12"/>`)
	iconKeyHome     = strokeIcon(`<path d="M8 8V4"/><path d="M4 8h4"/><path d="m4 4 10 10"/>`)
	iconKeyEnd      = strokeIcon(`<path d="M16 16v4"/><path d="M20 16h-4"/><path d="m20 20-10-10"/>`)

	// Mouse Event SVG Icons exactly replicating keyviz-source custom SVGs
	iconKeyMouseLeft       = strokeIcon(`<path d="M5 11L5 15C5 18.866 8.13401 22 12 22C15.866 22 19 18.866 19 15V9C19 5.13401 15.866 2 12 2C10.9264 2 9.90926 2.24169 9 2.67363"/><path d="M12 6V10"/><circle cx="5" cy="6" r="2" fill="currentColor"/>`)
	iconKeyMouseRight      = strokeIcon(`<path d="M19 11V15C19 18.866 15.866 22 12 22C8.13401 22 5 18.866 5 15V9C5 5.13401 8.13401 2 12 2C13.0736 2 14.0907 2.24169 15 2.67363"/><path d="M12 6V10"/><circle cx="19" cy="6" r="2" fill="currentColor"/>`)
	iconKeyMouseMiddle     = strokeIcon(`<path d="M19 9C19 5.13401 15.866 2 12 2C8.13401 2 5 5.13401 5 9V15C5 18.866 8.13401 22 12 22C15.866 22 19 18.866 19 15V9Z"/><circle cx="12" cy="8" r="2" fill="currentColor"/>`)
	iconKeyMouseScrollUp   = strokeIcon(`<path d="M19 9C19 5.13401 15.866 2 12 2C8.13401 2 5 5.13401 5 9V15C5 18.866 8.13401 22 12 22C15.866 22 19 18.866 19 15V9Z"/><path d="M12 7V12"/><path d="M14 8L12 6L10 8"/>`)
	iconKeyMouseScrollDown = strokeIcon(`<path d="M19 9C19 5.13401 15.866 2 12 2C8.13401 2 5 5.13401 5 9V15C5 18.866 8.13401 22 12 22C15.866 22 19 18.866 19 15V9Z"/><path d="M12 7V12"/><path d="M10 10L12 12L14 10"/>`)
	iconKeyMouseDrag       = strokeIcon(`<path d="M5.1851 18.9941C9.48005 21.4312 12.2743 19.1116 14.3687 15.5464C16.463 11.9811 17.1098 8.44303 12.8149 6.00594C8.51993 3.56885 5.72575 5.8884 3.63136 9.45367C1.53697 13.0189 0.890156 16.557 5.1851 18.9941Z"/><path d="M12 8L12.7192 6.70551C13.6233 5.07824 14.0753 4.26461 14.8427 4.05095C15.61 3.83729 16.393 4.30704 17.9589 5.24654L20.0351 6.49216C20.7231 6.90492 21.6028 6.65997 22 5.94505"/>`)
)

// getKeyIcon returns the vector SVG icon corresponding to the key label and modifier state.
func getKeyIcon(label string, isMod bool, cfg config.AppConfig) *ui.SVG {
	norm := strings.ToLower(strings.TrimSpace(label))

	// 1. Mouse Keys
	// The bare words "left" and "right" are deliberately absent: they are arrow
	// keys, and matching them here drew a mouse on every Left/Right press. Mouse
	// events arrive as "Left Click" / "Right Click", so the qualified names are
	// all this table needs.
	switch norm {
	case "left click", "mouse left":
		return iconKeyMouseLeft
	case "right click", "mouse right":
		return iconKeyMouseRight
	case "middle click", "mouse middle", "middle":
		return iconKeyMouseMiddle
	case "scroll up", "wheel up":
		return iconKeyMouseScrollUp
	case "scroll down", "wheel down":
		return iconKeyMouseScrollDown
	case "drag", "mouse drag":
		return iconKeyMouseDrag
	}

	// 2. Modifiers
	switch norm {
	case "ctrl", "control":
		return iconKeyCtrl
	case "shift":
		return iconKeyShift
	case "alt", "alternate", "opt", "option":
		return iconKeyAlt
	case "win", "windows", "super", "cmd", "command", "meta":
		if runtime.GOOS == "darwin" {
			return iconKeyCmd
		}
		return iconKeyWin
	}

	// 3. Special & Navigation Keys
	switch norm {
	case "enter", "return", "↵":
		return iconKeyReturn
	case "backspace", "⌫", "back", "bksp", "delete", "del":
		return iconKeyBackspace
	case "tab", "⇥":
		return iconKeyTab
	case "space", "␣":
		return iconKeySpace
	case "esc", "escape", "⎋":
		return iconKeyEscape
	case "capslock", "caps lock":
		return iconKeyCapsLock
	case "up", "arrowup", "↑":
		return iconKeyArrowUp
	case "down", "arrowdown", "↓":
		return iconKeyArrowDown
	case "left", "arrowleft", "←":
		return iconKeyArrowLeft
	case "right", "arrowright", "→":
		return iconKeyArrowRight
	case "pageup", "pgup":
		return iconKeyPageUp
	case "pagedown", "pgdn":
		return iconKeyPageDown
	case "home":
		return iconKeyHome
	case "end":
		return iconKeyEnd
	}

	return nil
}

// isArrowKey reports whether the key label represents an arrow direction key.
func isArrowKey(label string) bool {
	norm := strings.ToLower(strings.TrimSpace(label))
	return norm == "up" || norm == "down" || norm == "left" || norm == "right" ||
		norm == "arrowup" || norm == "arrowdown" || norm == "arrowleft" || norm == "arrowright" ||
		norm == "↑" || norm == "↓" || norm == "←" || norm == "→"
}
