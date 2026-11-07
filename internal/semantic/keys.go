package semantic

import (
	"strings"

	"github.com/yviscool/cue/internal/hook"
)

// This file owns the single canonical spelling of a primary key. Every layer
// (registry, profiles, user overrides, matcher) funnels through canonicalKey,
// so the same physical key always produces the same comboKey no matter how it
// was written down.
//
// The problem it solves: a keystroke arrives as (modifiers, symbol). On a US
// layout Ctrl+Shift+1 arrives as (Ctrl|Shift, "!"). A profile author may write
// it as Ctrl+Shift+1, Ctrl+!, or Ctrl+Shift+Exclamation. All three must resolve
// to one map entry, and a user's customShortcuts entry must match the built-in
// registry even when the two spell the chord differently.

// keyAliases canonicalizes alternate spellings of a primary key: hook glyphs
// (↑, ↵, ␣), Win32 names (PGDN, PRTSC), and friendly aliases.
var keyAliases = map[string]string{
	// Arrows
	"↑": "UP", "ARROWUP": "UP", "UPARROW": "UP",
	"↓": "DOWN", "ARROWDOWN": "DOWN", "DOWNARROW": "DOWN",
	"←": "LEFT", "ARROWLEFT": "LEFT", "LEFTARROW": "LEFT",
	"→": "RIGHT", "ARROWRIGHT": "RIGHT", "RIGHTARROW": "RIGHT",

	// Editing keys
	"↵": "ENTER", "RETURN": "ENTER",
	"␣": "SPACE", "SPACEBAR": "SPACE",
	"⌫": "BACKSPACE", "BKSP": "BACKSPACE", "BS": "BACKSPACE",
	"DEL":    "DELETE",
	"⇥":      "TAB",
	"ESCAPE": "ESC",
	"PGUP":   "PAGEUP",
	"PGDN":   "PAGEDOWN",
	"INS":    "INSERT",
	"PRTSC":  "PRINTSCREEN", "PRTSCN": "PRINTSCREEN",

	// Mouse
	"LBUTTON": "LEFT CLICK", "MOUSE LEFT": "LEFT CLICK",
	"RBUTTON": "RIGHT CLICK", "MOUSE RIGHT": "RIGHT CLICK",
	"MBUTTON": "MIDDLE CLICK", "MOUSE MIDDLE": "MIDDLE CLICK",
	"SCROLLUP": "SCROLL UP", "WHEEL UP": "SCROLL UP", "WHEELUP": "SCROLL UP",
	"SCROLLDOWN": "SCROLL DOWN", "WHEEL DOWN": "SCROLL DOWN", "WHEELDOWN": "SCROLL DOWN",

	// Spelled-out punctuation, so hand-written configuration and profiles can
	// say Ctrl+Slash instead of Ctrl+/. Shifted spellings resolve to the shifted
	// symbol, which then folds into Shift via shiftedKeys.
	"GRAVE": "`", "BACKQUOTE": "`", "BACKTICK": "`", "TILDE": "~",
	"SLASH": "/", "FORWARDSLASH": "/", "BACKSLASH": "\\", "PIPE": "|",
	"PLUS": "+", "ADD": "+", "MINUS": "-", "SUBTRACT": "-",
	"EQUAL": "=", "EQUALS": "=", "UNDERSCORE": "_",
	"PERIOD": ".", "DOT": ".", "COMMA": ",",
	"SEMICOLON": ";", "COLON": ":",
	"APOSTROPHE": "'", "QUOTE": "'", "QUOTEDBL": "\"",
	"LBRACKET": "[", "LEFTBRACKET": "[", "BRACKETLEFT": "[",
	"RBRACKET": "]", "RIGHTBRACKET": "]", "BRACKETRIGHT": "]",
	"LBRACE": "{", "RBRACE": "}",
	"EXCLAMATION": "!", "AT": "@", "HASH": "#", "NUMBERSIGN": "#",
	"DOLLAR": "$", "PERCENT": "%", "CARET": "^", "AMPERSAND": "&",
	"ASTERISK": "*", "STAR": "*",
	"LEFTPAREN": "(", "RIGHTPAREN": ")",
	"LESS": "<", "GREATER": ">", "QUESTION": "?",
}

// shiftedKeys maps a shifted symbol back to the physical key that produces it
// on a standard layout. Registering "Ctrl+!" and "Ctrl+Shift+1" therefore
// produce the same comboKey, which is what makes shifted punctuation safe to
// use in profiles.
var shiftedKeys = map[string]string{
	"!": "1", "@": "2", "#": "3", "$": "4", "%": "5",
	"^": "6", "&": "7", "*": "8", "(": "9", ")": "0",

	"~": "`", "_": "-", "+": "=",
	"{": "[", "}": "]", "|": "\\",
	":": ";", "\"": "'",
	"<": ",", ">": ".", "?": "/",
}

// unshiftedKeys is the inverse of shiftedKeys, used when rendering a chord for
// documentation.
var unshiftedKeys = func() map[string]string {
	m := make(map[string]string, len(shiftedKeys))
	for shifted, base := range shiftedKeys {
		m[base] = shifted
	}
	return m
}()

// canonicalKey resolves a raw key name into its physical base key and reports
// whether the symbol itself implies a held Shift.
func canonicalKey(keyName string) (base string, impliedShift bool) {
	// A single space is a real key and would be trimmed away otherwise.
	if keyName == " " {
		return "SPACE", false
	}
	k := strings.ToUpper(strings.TrimSpace(keyName))
	if k == "" {
		return "", false
	}
	if alias, ok := keyAliases[k]; ok {
		k = alias
	}
	if physical, ok := shiftedKeys[k]; ok {
		return physical, true
	}
	return k, false
}

// NormalizeKey standardizes key representations (glyphs, aliases, casing) and
// folds shifted symbols onto their physical key. It is the public spelling used
// by configuration parsing; use canonicalCombo when a modifier mask matters.
func NormalizeKey(keyName string) string {
	base, _ := canonicalKey(keyName)
	return base
}

// canonicalCombo builds the registry key for a chord, folding a shifted symbol
// into the Shift modifier so all spellings agree.
func canonicalCombo(mods hook.Modifiers, key string) comboKey {
	base, impliedShift := canonicalKey(key)
	if impliedShift {
		mods |= hook.ModShift
	}
	return comboKey{mods: mods, key: base}
}

// keyLabel renders the physical key for documentation. When Shift is held the
// shifted symbol is shown alongside the base key so the reader can recognise
// either spelling: "Ctrl + Shift + 1 (!)".
func keyLabel(combo comboKey) string {
	if symbol, ok := unshiftedKeys[combo.key]; ok && combo.mods.HasShift() {
		return combo.key + " (" + symbol + ")"
	}
	return combo.key
}

// Format renders a chord in human-readable form, e.g. "Ctrl + Shift + P".
func (c comboKey) Format() string {
	parts := make([]string, 0, 5)
	if c.mods.HasSuper() {
		parts = append(parts, "Win")
	}
	if c.mods.HasCtrl() {
		parts = append(parts, "Ctrl")
	}
	if c.mods.HasAlt() {
		parts = append(parts, "Alt")
	}
	if c.mods.HasShift() {
		parts = append(parts, "Shift")
	}
	parts = append(parts, keyLabel(c))
	return strings.Join(parts, " + ")
}

// Format renders a two-chord sequence, e.g. "Ctrl + K, Ctrl + B".
func (s sequenceKey) Format() string {
	return s.first.Format() + ", " + s.second.Format()
}
