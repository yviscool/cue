package hook

import (
	"fmt"
	"runtime"
)

// Virtual Key Code constants for Windows and cross-platform compatibility.
const (
	VK_LBUTTON  = 0x01
	VK_RBUTTON  = 0x02
	VK_CANCEL   = 0x03
	VK_MBUTTON  = 0x04
	VK_XBUTTON1 = 0x05
	VK_XBUTTON2 = 0x06
	VK_BACK     = 0x08 // Backspace
	VK_TAB      = 0x09 // Tab
	VK_CLEAR    = 0x0C
	VK_RETURN   = 0x0D // Enter
	VK_SHIFT    = 0x10
	VK_CONTROL  = 0x11
	VK_MENU     = 0x12 // Alt
	VK_PAUSE    = 0x13
	VK_CAPITAL  = 0x14 // Caps Lock
	VK_ESCAPE   = 0x1B // Esc
	VK_SPACE    = 0x20 // Space
	VK_PRIOR    = 0x21 // Page Up
	VK_NEXT     = 0x22 // Page Down
	VK_END      = 0x23 // End
	VK_HOME     = 0x24 // Home
	VK_LEFT     = 0x25 // Left Arrow
	VK_UP       = 0x26 // Up Arrow
	VK_RIGHT    = 0x27 // Right Arrow
	VK_DOWN     = 0x28 // Down Arrow
	VK_SELECT   = 0x29
	VK_PRINT    = 0x2A
	VK_EXECUTE  = 0x2B
	VK_SNAPSHOT = 0x2C // Print Screen
	VK_INSERT   = 0x2D // Insert
	VK_DELETE   = 0x2E // Delete
	VK_HELP     = 0x2F

	// Standard alphanumeric keys
	VK_0 = 0x30
	VK_1 = 0x31
	VK_2 = 0x32
	VK_3 = 0x33
	VK_4 = 0x34
	VK_5 = 0x35
	VK_6 = 0x36
	VK_7 = 0x37
	VK_8 = 0x38
	VK_9 = 0x39

	VK_A = 0x41
	VK_B = 0x42
	VK_C = 0x43
	VK_D = 0x44
	VK_E = 0x45
	VK_F = 0x46
	VK_G = 0x47
	VK_H = 0x48
	VK_I = 0x49
	VK_J = 0x4A
	VK_K = 0x4B
	VK_L = 0x4C
	VK_M = 0x4D
	VK_N = 0x4E
	VK_O = 0x4F
	VK_P = 0x50
	VK_Q = 0x51
	VK_R = 0x52
	VK_S = 0x53
	VK_T = 0x54
	VK_U = 0x55
	VK_V = 0x56
	VK_W = 0x57
	VK_X = 0x58
	VK_Y = 0x59
	VK_Z = 0x5A

	// System / OS keys
	VK_LWIN  = 0x5B // Left Windows / Super
	VK_RWIN  = 0x5C // Right Windows / Super
	VK_APPS  = 0x5D // Applications / Context menu
	VK_SLEEP = 0x5F

	// Numpad keys
	VK_NUMPAD0   = 0x60
	VK_NUMPAD1   = 0x61
	VK_NUMPAD2   = 0x62
	VK_NUMPAD3   = 0x63
	VK_NUMPAD4   = 0x64
	VK_NUMPAD5   = 0x65
	VK_NUMPAD6   = 0x66
	VK_NUMPAD7   = 0x67
	VK_NUMPAD8   = 0x68
	VK_NUMPAD9   = 0x69
	VK_MULTIPLY  = 0x6A // Num*
	VK_ADD       = 0x6B // Num+
	VK_SEPARATOR = 0x6C // NumSep
	VK_SUBTRACT  = 0x6D // Num-
	VK_DECIMAL   = 0x6E // Num.
	VK_DIVIDE    = 0x6F // Num/

	// Function keys
	VK_F1  = 0x70
	VK_F2  = 0x71
	VK_F3  = 0x72
	VK_F4  = 0x73
	VK_F5  = 0x74
	VK_F6  = 0x75
	VK_F7  = 0x76
	VK_F8  = 0x77
	VK_F9  = 0x78
	VK_F10 = 0x79
	VK_F11 = 0x7A
	VK_F12 = 0x7B
	VK_F13 = 0x7C
	VK_F14 = 0x7D
	VK_F15 = 0x7E
	VK_F16 = 0x7F
	VK_F17 = 0x80
	VK_F18 = 0x81
	VK_F19 = 0x82
	VK_F20 = 0x83
	VK_F21 = 0x84
	VK_F22 = 0x85
	VK_F23 = 0x86
	VK_F24 = 0x87

	VK_NUMLOCK = 0x90
	VK_SCROLL  = 0x91

	// Left and Right modifier variants
	VK_LSHIFT   = 0xA0
	VK_RSHIFT   = 0xA1
	VK_LCONTROL = 0xA2
	VK_RCONTROL = 0xA3
	VK_LMENU    = 0xA4
	VK_RMENU    = 0xA5

	// Media and browser keys
	VK_BROWSER_BACK        = 0xA6
	VK_BROWSER_FORWARD     = 0xA7
	VK_BROWSER_REFRESH     = 0xA8
	VK_BROWSER_STOP        = 0xA9
	VK_BROWSER_SEARCH      = 0xAA
	VK_BROWSER_FAVORITES   = 0xAB
	VK_BROWSER_HOME        = 0xAC
	VK_VOLUME_MUTE         = 0xAD
	VK_VOLUME_DOWN         = 0xAE
	VK_VOLUME_UP           = 0xAF
	VK_MEDIA_NEXT_TRACK    = 0xB0
	VK_MEDIA_PREV_TRACK    = 0xB1
	VK_MEDIA_STOP          = 0xB2
	VK_MEDIA_PLAY_PAUSE    = 0xB3
	VK_LAUNCH_MAIL         = 0xB4
	VK_LAUNCH_MEDIA_SELECT = 0xB5
	VK_LAUNCH_APP1         = 0xB6
	VK_LAUNCH_APP2         = 0xB7

	// OEM Punctuation keys
	VK_OEM_1      = 0xBA // ; :
	VK_OEM_PLUS   = 0xBB // = +
	VK_OEM_COMMA  = 0xBC // , <
	VK_OEM_MINUS  = 0xBD // - _
	VK_OEM_PERIOD = 0xBE // . >
	VK_OEM_2      = 0xBF // / ?
	VK_OEM_3      = 0xC0 // ` ~
	VK_OEM_4      = 0xDB // [ {
	VK_OEM_5      = 0xDC // \ |
	VK_OEM_6      = 0xDD // ] }
	VK_OEM_7      = 0xDE // ' "
	VK_OEM_8      = 0xDF
	VK_OEM_102    = 0xE2 // <> or \|
)

// shiftSymbols maps base characters (numbers & punctuation) to their shifted symbols on standard keyboards.
var shiftSymbols = map[string]string{
	// Numbers with Shift
	"1": "!",
	"2": "@",
	"3": "#",
	"4": "$",
	"5": "%",
	"6": "^",
	"7": "&",
	"8": "*",
	"9": "(",
	"0": ")",

	// Punctuation with Shift
	"`":  "~",
	"-":  "_",
	"=":  "+",
	"[":  "{",
	"]":  "}",
	"\\": "|",
	";":  ":",
	"'":  "\"",
	",":  "<",
	".":  ">",
	"/":  "?",
}

// vkShiftMap maps virtual keycodes directly to their shifted symbols.
var vkShiftMap = map[uint32]string{
	VK_1: "!",
	VK_2: "@",
	VK_3: "#",
	VK_4: "$",
	VK_5: "%",
	VK_6: "^",
	VK_7: "&",
	VK_8: "*",
	VK_9: "(",
	VK_0: ")",

	VK_OEM_3:      "~",
	VK_OEM_MINUS:  "_",
	VK_OEM_PLUS:   "+",
	VK_OEM_4:      "{",
	VK_OEM_6:      "}",
	VK_OEM_5:      "|",
	VK_OEM_1:      ":",
	VK_OEM_7:      "\"",
	VK_OEM_COMMA:  "<",
	VK_OEM_PERIOD: ">",
	VK_OEM_2:      "?",
}

// superKeyName returns "⌘" on macOS and "Win" on other systems.
func superKeyName() string {
	if runtime.GOOS == "darwin" {
		return "⌘"
	}
	return "Win"
}

// baseKeyNames maps virtual key codes to clean, human-readable glyphs and names.
var baseKeyNames = map[uint32]string{
	// Special keys with clean glyph symbols
	VK_RETURN: "↵",
	VK_BACK:   "⌫",
	VK_TAB:    "⇥",
	VK_SPACE:  "␣",
	VK_ESCAPE: "Esc",
	VK_UP:     "↑",
	VK_DOWN:   "↓",
	VK_LEFT:   "←",
	VK_RIGHT:  "→",

	// Modifiers
	VK_CONTROL:  "Ctrl",
	VK_LCONTROL: "Ctrl",
	VK_RCONTROL: "Ctrl",
	VK_MENU:     "Alt",
	VK_LMENU:    "Alt",
	VK_RMENU:    "Alt",
	VK_SHIFT:    "Shift",
	VK_LSHIFT:   "Shift",
	VK_RSHIFT:   "Shift",
	VK_LWIN:     superKeyName(),
	VK_RWIN:     superKeyName(),

	// Function keys F1-F12
	VK_F1:  "F1",
	VK_F2:  "F2",
	VK_F3:  "F3",
	VK_F4:  "F4",
	VK_F5:  "F5",
	VK_F6:  "F6",
	VK_F7:  "F7",
	VK_F8:  "F8",
	VK_F9:  "F9",
	VK_F10: "F10",
	VK_F11: "F11",
	VK_F12: "F12",
	VK_F13: "F13",
	VK_F14: "F14",
	VK_F15: "F15",
	VK_F16: "F16",
	VK_F17: "F17",
	VK_F18: "F18",
	VK_F19: "F19",
	VK_F20: "F20",
	VK_F21: "F21",
	VK_F22: "F22",
	VK_F23: "F23",
	VK_F24: "F24",

	// Numpad keys
	VK_NUMPAD0:   "Num0",
	VK_NUMPAD1:   "Num1",
	VK_NUMPAD2:   "Num2",
	VK_NUMPAD3:   "Num3",
	VK_NUMPAD4:   "Num4",
	VK_NUMPAD5:   "Num5",
	VK_NUMPAD6:   "Num6",
	VK_NUMPAD7:   "Num7",
	VK_NUMPAD8:   "Num8",
	VK_NUMPAD9:   "Num9",
	VK_MULTIPLY:  "Num*",
	VK_ADD:       "Num+",
	VK_SEPARATOR: "NumSep",
	VK_SUBTRACT:  "Num-",
	VK_DECIMAL:   "Num.",
	VK_DIVIDE:    "Num/",
	VK_NUMLOCK:   "NumLock",

	// Navigation & Editing
	VK_HOME:     "Home",
	VK_END:      "End",
	VK_PRIOR:    "PgUp",
	VK_NEXT:     "PgDn",
	VK_DELETE:   "Delete",
	VK_INSERT:   "Insert",
	VK_CAPITAL:  "CapsLock",
	VK_SNAPSHOT: "PrintScreen",
	VK_SCROLL:   "ScrollLock",
	VK_PAUSE:    "Pause",
	VK_APPS:     "Menu",

	// Numbers
	VK_0: "0",
	VK_1: "1",
	VK_2: "2",
	VK_3: "3",
	VK_4: "4",
	VK_5: "5",
	VK_6: "6",
	VK_7: "7",
	VK_8: "8",
	VK_9: "9",

	// Letters
	VK_A: "A",
	VK_B: "B",
	VK_C: "C",
	VK_D: "D",
	VK_E: "E",
	VK_F: "F",
	VK_G: "G",
	VK_H: "H",
	VK_I: "I",
	VK_J: "J",
	VK_K: "K",
	VK_L: "L",
	VK_M: "M",
	VK_N: "N",
	VK_O: "O",
	VK_P: "P",
	VK_Q: "Q",
	VK_R: "R",
	VK_S: "S",
	VK_T: "T",
	VK_U: "U",
	VK_V: "V",
	VK_W: "W",
	VK_X: "X",
	VK_Y: "Y",
	VK_Z: "Z",

	// Punctuation (Base)
	VK_OEM_3:      "`",
	VK_OEM_MINUS:  "-",
	VK_OEM_PLUS:   "=",
	VK_OEM_4:      "[",
	VK_OEM_6:      "]",
	VK_OEM_5:      "\\",
	VK_OEM_1:      ";",
	VK_OEM_7:      "'",
	VK_OEM_COMMA:  ",",
	VK_OEM_PERIOD: ".",
	VK_OEM_2:      "/",
	VK_OEM_102:    "\\",

	// Media & Browser
	VK_VOLUME_MUTE:      "VolMute",
	VK_VOLUME_DOWN:      "VolDown",
	VK_VOLUME_UP:        "VolUp",
	VK_MEDIA_NEXT_TRACK: "MediaNext",
	VK_MEDIA_PREV_TRACK: "MediaPrev",
	VK_MEDIA_STOP:       "MediaStop",
	VK_MEDIA_PLAY_PAUSE: "MediaPlayPause",
}

// VKToKeyName translates a virtual keycode into its canonical base human-readable string.
func VKToKeyName(vk uint32) string {
	if name, ok := baseKeyNames[vk]; ok {
		return name
	}
	// Fallback to formatted hex name for unknown keycodes
	return fmt.Sprintf("VK_%02X", vk)
}

// ResolveKeyName translates a virtual keycode into its human-readable string,
// taking into account whether Shift is active for number and punctuation keys.
func ResolveKeyName(vk uint32, shift bool) string {
	if shift {
		if sym, ok := vkShiftMap[vk]; ok {
			return sym
		}
	}
	return VKToKeyName(vk)
}

// ShiftSymbol returns the shifted representation of a base character if defined,
// or the original string if no shift mapping exists.
func ShiftSymbol(key string) string {
	if sym, ok := shiftSymbols[key]; ok {
		return sym
	}
	return key
}

// UnshiftSymbol returns the unshifted base character for a given shifted symbol,
// or the original string if no reverse mapping exists.
func UnshiftSymbol(sym string) string {
	for base, s := range shiftSymbols {
		if s == sym {
			return base
		}
	}
	return sym
}

// ShiftRune returns the shifted rune if defined, or the rune unchanged.
func ShiftRune(r rune) rune {
	s := ShiftSymbol(string(r))
	runes := []rune(s)
	if len(runes) == 1 {
		return runes[0]
	}
	return r
}

// IsModifierKey reports whether the given virtual keycode represents a modifier key (Ctrl, Shift, Alt, Win).
func IsModifierKey(vk uint32) bool {
	switch vk {
	case VK_SHIFT, VK_LSHIFT, VK_RSHIFT,
		VK_CONTROL, VK_LCONTROL, VK_RCONTROL,
		VK_MENU, VK_LMENU, VK_RMENU,
		VK_LWIN, VK_RWIN:
		return true
	default:
		return false
	}
}

// ModifierFromVK returns the corresponding Modifiers bit flag for a modifier virtual keycode,
// or ModNone if the key is not a modifier.
func ModifierFromVK(vk uint32) Modifiers {
	switch vk {
	case VK_SHIFT, VK_LSHIFT, VK_RSHIFT:
		return ModShift
	case VK_CONTROL, VK_LCONTROL, VK_RCONTROL:
		return ModCtrl
	case VK_MENU, VK_LMENU, VK_RMENU:
		return ModAlt
	case VK_LWIN, VK_RWIN:
		return ModSuper
	default:
		return ModNone
	}
}
