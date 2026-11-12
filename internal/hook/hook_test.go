package hook

import (
	"runtime"
	"testing"
	"time"
	"unsafe"
)

// TestKeymap_VirtualKeyCodes verifies virtual keycode translation into canonical strings and glyphs.
func TestKeymap_VirtualKeyCodes(t *testing.T) {
	tests := []struct {
		vk       uint32
		expected string
	}{
		// Clean glyph symbols for special keys
		{VK_RETURN, "↵"},
		{VK_BACK, "⌫"},
		{VK_TAB, "⇥"},
		{VK_SPACE, "␣"},
		{VK_ESCAPE, "Esc"},
		{VK_UP, "↑"},
		{VK_DOWN, "↓"},
		{VK_LEFT, "←"},
		{VK_RIGHT, "→"},

		// Modifiers
		{VK_CONTROL, "Ctrl"},
		{VK_LCONTROL, "Ctrl"},
		{VK_RCONTROL, "Ctrl"},
		{VK_MENU, "Alt"},
		{VK_LMENU, "Alt"},
		{VK_RMENU, "Alt"},
		{VK_SHIFT, "Shift"},
		{VK_LSHIFT, "Shift"},
		{VK_RSHIFT, "Shift"},

		// Function keys F1-F12
		{VK_F1, "F1"},
		{VK_F2, "F2"},
		{VK_F3, "F3"},
		{VK_F4, "F4"},
		{VK_F5, "F5"},
		{VK_F6, "F6"},
		{VK_F7, "F7"},
		{VK_F8, "F8"},
		{VK_F9, "F9"},
		{VK_F10, "F10"},
		{VK_F11, "F11"},
		{VK_F12, "F12"},

		// Numpad keys
		{VK_NUMPAD0, "Num0"},
		{VK_NUMPAD1, "Num1"},
		{VK_NUMPAD2, "Num2"},
		{VK_NUMPAD3, "Num3"},
		{VK_NUMPAD4, "Num4"},
		{VK_NUMPAD5, "Num5"},
		{VK_NUMPAD6, "Num6"},
		{VK_NUMPAD7, "Num7"},
		{VK_NUMPAD8, "Num8"},
		{VK_NUMPAD9, "Num9"},
		{VK_MULTIPLY, "Num*"},
		{VK_ADD, "Num+"},
		{VK_SUBTRACT, "Num-"},
		{VK_DECIMAL, "Num."},
		{VK_DIVIDE, "Num/"},
		{VK_NUMLOCK, "NumLock"},

		// Navigation & Editing
		{VK_HOME, "Home"},
		{VK_END, "End"},
		{VK_PRIOR, "PgUp"},
		{VK_NEXT, "PgDn"},
		{VK_DELETE, "Delete"},
		{VK_INSERT, "Insert"},
		{VK_CAPITAL, "CapsLock"},
		{VK_SNAPSHOT, "PrintScreen"},
		{VK_SCROLL, "ScrollLock"},
		{VK_PAUSE, "Pause"},

		// Numbers
		{VK_0, "0"},
		{VK_1, "1"},
		{VK_2, "2"},
		{VK_3, "3"},
		{VK_4, "4"},
		{VK_5, "5"},
		{VK_6, "6"},
		{VK_7, "7"},
		{VK_8, "8"},
		{VK_9, "9"},

		// Letters
		{VK_A, "A"},
		{VK_B, "B"},
		{VK_C, "C"},
		{VK_Z, "Z"},

		// Punctuation (Base)
		{VK_OEM_3, "`"},
		{VK_OEM_MINUS, "-"},
		{VK_OEM_PLUS, "="},
		{VK_OEM_4, "["},
		{VK_OEM_6, "]"},
		{VK_OEM_5, "\\"},
		{VK_OEM_1, ";"},
		{VK_OEM_7, "'"},
		{VK_OEM_COMMA, ","},
		{VK_OEM_PERIOD, "."},
		{VK_OEM_2, "/"},
	}

	for _, tt := range tests {
		actual := VKToKeyName(tt.vk)
		if actual != tt.expected {
			t.Errorf("VKToKeyName(0x%02X) = %q, expected %q", tt.vk, actual, tt.expected)
		}
	}

	// Super/Win key check
	winExpected := "Win"
	if runtime.GOOS == "darwin" {
		winExpected = "⌘"
	}
	if actual := VKToKeyName(VK_LWIN); actual != winExpected {
		t.Errorf("VKToKeyName(VK_LWIN) = %q, expected %q", actual, winExpected)
	}
	if actual := VKToKeyName(VK_RWIN); actual != winExpected {
		t.Errorf("VKToKeyName(VK_RWIN) = %q, expected %q", actual, winExpected)
	}

	// Fallback check
	if actual := VKToKeyName(0xFF); actual != "VK_FF" {
		t.Errorf("VKToKeyName(0xFF) = %q, expected %q", actual, "VK_FF")
	}
}

// TestKeymap_ShiftSymbolMapping verifies the comprehensive Shift symbol mappings.
func TestKeymap_ShiftSymbolMapping(t *testing.T) {
	// Numbers with Shift: '1' -> '!', '2' -> '@', '3' -> '#', '4' -> '$', '5' -> '%', '6' -> '^', '7' -> '&', '8' -> '*', '9' -> '(', '0' -> ')'
	numTests := []struct {
		vk       uint32
		base     string
		expected string
	}{
		{VK_1, "1", "!"},
		{VK_2, "2", "@"},
		{VK_3, "3", "#"},
		{VK_4, "4", "$"},
		{VK_5, "5", "%"},
		{VK_6, "6", "^"},
		{VK_7, "7", "&"},
		{VK_8, "8", "*"},
		{VK_9, "9", "("},
		{VK_0, "0", ")"},
	}

	for _, tt := range numTests {
		// Test ShiftSymbol helper
		if sym := ShiftSymbol(tt.base); sym != tt.expected {
			t.Errorf("ShiftSymbol(%q) = %q, expected %q", tt.base, sym, tt.expected)
		}
		// Test UnshiftSymbol helper
		if base := UnshiftSymbol(tt.expected); base != tt.base {
			t.Errorf("UnshiftSymbol(%q) = %q, expected %q", tt.expected, base, tt.base)
		}
		// Test ShiftRune
		if r := ShiftRune([]rune(tt.base)[0]); string(r) != tt.expected {
			t.Errorf("ShiftRune(%q) = %q, expected %q", tt.base, string(r), tt.expected)
		}
		// Test ResolveKeyName with Shift
		if res := ResolveKeyName(tt.vk, true); res != tt.expected {
			t.Errorf("ResolveKeyName(0x%02X, true) = %q, expected %q", tt.vk, res, tt.expected)
		}
		// Test ResolveKeyName without Shift
		if res := ResolveKeyName(tt.vk, false); res != tt.base {
			t.Errorf("ResolveKeyName(0x%02X, false) = %q, expected %q", tt.vk, res, tt.base)
		}
	}

	// Punctuation with Shift: '`' -> '~', '-' -> '_', '=' -> '+', '[' -> '{', ']' -> '}', '\' -> '|', ';' -> ':', '\'' -> '"', ',' -> '<', '.' -> '>', '/' -> '?'
	puncTests := []struct {
		vk       uint32
		base     string
		expected string
	}{
		{VK_OEM_3, "`", "~"},
		{VK_OEM_MINUS, "-", "_"},
		{VK_OEM_PLUS, "=", "+"},
		{VK_OEM_4, "[", "{"},
		{VK_OEM_6, "]", "}"},
		{VK_OEM_5, "\\", "|"},
		{VK_OEM_1, ";", ":"},
		{VK_OEM_7, "'", "\""},
		{VK_OEM_COMMA, ",", "<"},
		{VK_OEM_PERIOD, ".", ">"},
		{VK_OEM_2, "/", "?"},
	}

	for _, tt := range puncTests {
		// Test ShiftSymbol helper
		if sym := ShiftSymbol(tt.base); sym != tt.expected {
			t.Errorf("ShiftSymbol(%q) = %q, expected %q", tt.base, sym, tt.expected)
		}
		// Test UnshiftSymbol helper
		if base := UnshiftSymbol(tt.expected); base != tt.base {
			t.Errorf("UnshiftSymbol(%q) = %q, expected %q", tt.expected, base, tt.base)
		}
		// Test ShiftRune
		if r := ShiftRune([]rune(tt.base)[0]); string(r) != tt.expected {
			t.Errorf("ShiftRune(%q) = %q, expected %q", tt.base, string(r), tt.expected)
		}
		// Test ResolveKeyName with Shift
		if res := ResolveKeyName(tt.vk, true); res != tt.expected {
			t.Errorf("ResolveKeyName(0x%02X, true) = %q, expected %q", tt.vk, res, tt.expected)
		}
		// Test ResolveKeyName without Shift
		if res := ResolveKeyName(tt.vk, false); res != tt.base {
			t.Errorf("ResolveKeyName(0x%02X, false) = %q, expected %q", tt.vk, res, tt.base)
		}
	}

	// Letters with Shift should remain unchanged uppercase letter
	if res := ResolveKeyName(VK_A, true); res != "A" {
		t.Errorf("ResolveKeyName(VK_A, true) = %q, expected %q", res, "A")
	}
	if res := ShiftSymbol("A"); res != "A" {
		t.Errorf("ShiftSymbol(\"A\") = %q, expected %q", res, "A")
	}
}

// TestModifiers_Bitmask verifies the behavior of the Modifiers bitmask.
func TestModifiers_Bitmask(t *testing.T) {
	var m Modifiers
	if !m.IsEmpty() {
		t.Errorf("empty Modifiers should be empty")
	}
	if m.HasCtrl() || m.HasShift() || m.HasAlt() || m.HasSuper() {
		t.Errorf("empty Modifiers should have no flags set")
	}
	if m.String() != "" {
		t.Errorf("empty Modifiers String() = %q, expected \"\"", m.String())
	}

	m |= ModCtrl
	if !m.HasCtrl() || m.HasShift() || m.HasAlt() || m.HasSuper() {
		t.Errorf("m with ModCtrl has incorrect flags: %v", m)
	}
	if m.String() != "Ctrl" {
		t.Errorf("m.String() = %q, expected \"Ctrl\"", m.String())
	}

	m |= ModShift
	if !m.HasCtrl() || !m.HasShift() || m.HasAlt() || m.HasSuper() {
		t.Errorf("m with ModCtrl|ModShift has incorrect flags: %v", m)
	}

	m |= ModAlt
	m |= ModSuper
	if !m.HasCtrl() || !m.HasShift() || !m.HasAlt() || !m.HasSuper() {
		t.Errorf("m with all flags has incorrect flags: %v", m)
	}
	if m.String() != "Ctrl+Alt+Shift+Win" {
		t.Errorf("m.String() = %q, expected \"Ctrl+Alt+Shift+Win\"", m.String())
	}
}

// TestIsModifierKey_And_ModifierFromVK verifies modifier key recognition.
func TestIsModifierKey_And_ModifierFromVK(t *testing.T) {
	modifiers := []struct {
		vk   uint32
		mods Modifiers
	}{
		{VK_SHIFT, ModShift},
		{VK_LSHIFT, ModShift},
		{VK_RSHIFT, ModShift},
		{VK_CONTROL, ModCtrl},
		{VK_LCONTROL, ModCtrl},
		{VK_RCONTROL, ModCtrl},
		{VK_MENU, ModAlt},
		{VK_LMENU, ModAlt},
		{VK_RMENU, ModAlt},
		{VK_LWIN, ModSuper},
		{VK_RWIN, ModSuper},
	}

	for _, tt := range modifiers {
		if !IsModifierKey(tt.vk) {
			t.Errorf("IsModifierKey(0x%02X) = false, expected true", tt.vk)
		}
		if m := ModifierFromVK(tt.vk); m != tt.mods {
			t.Errorf("ModifierFromVK(0x%02X) = %v, expected %v", tt.vk, m, tt.mods)
		}
	}

	nonModifiers := []uint32{VK_A, VK_0, VK_RETURN, VK_SPACE, VK_F1, VK_DELETE}
	for _, vk := range nonModifiers {
		if IsModifierKey(vk) {
			t.Errorf("IsModifierKey(0x%02X) = true, expected false", vk)
		}
		if m := ModifierFromVK(vk); m != ModNone {
			t.Errorf("ModifierFromVK(0x%02X) = %v, expected ModNone", vk, m)
		}
	}
}

// TestHook_Lifecycle verifies Start, Stop, and channel lifecycle.
func TestHook_Lifecycle(t *testing.T) {
	h := New()
	if h.Events() == nil {
		t.Fatal("h.Events() channel should not be nil")
	}

	err := h.Start()
	if runtime.GOOS != "windows" {
		if err != ErrUnsupportedOS {
			t.Fatalf("expected ErrUnsupportedOS on non-windows, got %v", err)
		}
		return
	}

	if err != nil {
		t.Fatalf("h.Start() failed: %v", err)
	}

	// Double start should return an error
	if err := h.Start(); err == nil {
		t.Error("second h.Start() should return an error")
	}

	// Stop cleanly
	h.Stop()

	// Calling Stop again should be safe and idempotent
	h.Stop()

	// Events channel should be closed after Stop
	select {
	case _, ok := <-h.Events():
		if ok {
			t.Error("Events channel should be closed after Stop")
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for Events channel to close")
	}
}

// TestHook_KeyboardSimulation verifies event generation, repeat detection, and modifier handling.
func TestHook_KeyboardSimulation(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific simulation on non-Windows")
	}

	wh := NewHook(64).(*windowsHook)

	// 1. Simulate keydown for 'A'
	wh.processKeyboard(WM_KEYDOWN, VK_A, 0x1E, false)

	select {
	case ev := <-wh.Events():
		if ev.Type != EventKeyDown {
			t.Errorf("expected EventKeyDown, got %v", ev.Type)
		}
		if ev.VKCode != VK_A {
			t.Errorf("expected VKCode VK_A, got 0x%02X", ev.VKCode)
		}
		if ev.KeyName != "A" {
			t.Errorf("expected KeyName \"A\", got %q", ev.KeyName)
		}
		if ev.IsKeyRepeat {
			t.Errorf("first key down should not be repeat")
		}
		if ev.IsModifier {
			t.Errorf("A should not be a modifier")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for event")
	}

	// 2. Simulate repeat keydown for 'A'
	wh.processKeyboard(WM_KEYDOWN, VK_A, 0x1E, false)
	select {
	case ev := <-wh.Events():
		if !ev.IsKeyRepeat {
			t.Errorf("second consecutive keydown should have IsKeyRepeat = true")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for repeat event")
	}

	// 3. Simulate keyup for 'A'
	wh.processKeyboard(WM_KEYUP, VK_A, 0x1E, false)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventKeyUp {
			t.Errorf("expected EventKeyUp, got %v", ev.Type)
		}
		if ev.IsKeyRepeat {
			t.Errorf("keyup should not be repeat")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for keyup event")
	}

	// 4. Simulate Shift keydown, then '1' keydown -> should produce '!'
	wh.processKeyboard(WM_KEYDOWN, VK_LSHIFT, 0x2A, false)
	select {
	case ev := <-wh.Events():
		if ev.KeyName != "Shift" || !ev.IsModifier || !ev.Modifiers.HasShift() {
			t.Errorf("Shift event invalid: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for Shift event")
	}

	wh.processKeyboard(WM_KEYDOWN, VK_1, 0x02, false)
	select {
	case ev := <-wh.Events():
		if ev.KeyName != "!" {
			t.Errorf("expected KeyName '!' with Shift, got %q", ev.KeyName)
		}
		if !ev.Modifiers.HasShift() {
			t.Errorf("expected Modifiers to include Shift")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for '!' event")
	}

	// 5. Release Shift, then press '1' -> should produce '1'
	wh.processKeyboard(WM_KEYUP, VK_LSHIFT, 0x2A, false)
	select {
	case ev := <-wh.Events():
		if ev.Modifiers.HasShift() {
			t.Errorf("Shift should not be active after keyup")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for Shift keyup event")
	}

	wh.processKeyboard(WM_KEYDOWN, VK_1, 0x02, false)
	select {
	case ev := <-wh.Events():
		if ev.KeyName != "1" {
			t.Errorf("expected KeyName '1' without Shift, got %q", ev.KeyName)
		}
		if ev.Modifiers.HasShift() {
			t.Errorf("Modifiers should not have Shift")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for '1' event")
	}

	// 6. Test special keys: Enter ("↵") and Backspace ("⌫")
	wh.processKeyboard(WM_KEYDOWN, VK_RETURN, 0x1C, false)
	select {
	case ev := <-wh.Events():
		if ev.KeyName != "↵" {
			t.Errorf("expected Enter glyph '↵', got %q", ev.KeyName)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for Enter event")
	}
}

// TestHook_MouseSimulation verifies mouse button, coordinates, wheel deltas, and event types.
func TestHook_MouseSimulation(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific simulation on non-Windows")
	}

	wh := NewHook(64).(*windowsHook)

	// 1. Mouse move
	wh.processMouse(WM_MOUSEMOVE, 100, 200, 0)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseMove || ev.X != 100 || ev.Y != 200 || ev.Button != ButtonNone {
			t.Errorf("invalid mouse move event: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for mouse move event")
	}

	// 2. Left click down and up
	wh.processMouse(WM_LBUTTONDOWN, 150, 250, 0)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseDown || ev.Button != ButtonLeft || ev.X != 150 || ev.Y != 250 {
			t.Errorf("invalid mouse left down event: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for mouse left down event")
	}

	wh.processMouse(WM_LBUTTONUP, 150, 250, 0)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseUp || ev.Button != ButtonLeft {
			t.Errorf("invalid mouse left up event: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for mouse left up event")
	}

	// 3. Right click
	wh.processMouse(WM_RBUTTONDOWN, 150, 250, 0)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseDown || ev.Button != ButtonRight {
			t.Errorf("invalid mouse right down event: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for mouse right down event")
	}

	// 4. Middle click
	wh.processMouse(WM_MBUTTONDOWN, 150, 250, 0)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseDown || ev.Button != ButtonMiddle {
			t.Errorf("invalid mouse middle down event: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for mouse middle down event")
	}

	// 5. Extra mouse buttons (XBUTTON1 & XBUTTON2)
	wh.processMouse(WM_XBUTTONDOWN, 0, 0, uint32(XBUTTON1)<<16)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseDown || ev.Button != ButtonExtra1 {
			t.Errorf("invalid mouse Extra1 event: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for mouse Extra1 event")
	}

	wh.processMouse(WM_XBUTTONUP, 0, 0, uint32(XBUTTON2)<<16)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseUp || ev.Button != ButtonExtra2 {
			t.Errorf("invalid mouse Extra2 event: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for mouse Extra2 event")
	}

	// 6. Mouse wheel: positive (120) and negative (-120)
	wh.processMouse(WM_MOUSEWHEEL, 0, 0, uint32(uint16(120))<<16)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseWheel || ev.WheelDelta != 120 {
			t.Errorf("invalid wheel up event: delta=%d", ev.WheelDelta)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for wheel up event")
	}

	downDelta := int16(-120)
	wh.processMouse(WM_MOUSEWHEEL, 0, 0, uint32(uint16(downDelta))<<16)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventMouseWheel || ev.WheelDelta != -120 {
			t.Errorf("invalid wheel down event: delta=%d", ev.WheelDelta)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for wheel down event")
	}
}

// TestHook_NonBlockingChannel verifies that events are dropped without deadlock when buffer is full.
func TestHook_NonBlockingChannel(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific simulation on non-Windows")
	}

	// Create hook with buffer size of 2
	wh := NewHook(2).(*windowsHook)

	// Dispatch 10 events rapidly without consuming
	done := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			wh.processKeyboard(WM_KEYDOWN, VK_A, 0, false)
		}
		close(done)
	}()

	select {
	case <-done:
		// Succeeded without blocking
	case <-time.After(1 * time.Second):
		t.Fatal("processKeyboard deadlocked or blocked on full channel buffer")
	}

	// Should have exactly 2 events in channel
	count := 0
	for len(wh.Events()) > 0 {
		<-wh.Events()
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 buffered events, got %d", count)
	}
}

// TestHook_InjectedFilter verifies that injected synthetic events (LLKHF_INJECTED) are ignored.
func TestHook_InjectedFilter(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific simulation on non-Windows")
	}

	wh := NewHook(16).(*windowsHook)

	// 1. Simulate injected key event (e.g. Win+Esc injecting Ctrl+A)
	injectedKbd := kbdllHookStruct{
		vkCode:   VK_A,
		scanCode: 0x1E,
		flags:    LLKHF_INJECTED,
	}
	wh.handleKeyboard(WM_KEYDOWN, unsafe.Pointer(&injectedKbd))

	select {
	case ev := <-wh.Events():
		t.Fatalf("unexpected event received for injected key: %+v", ev)
	case <-time.After(50 * time.Millisecond):
		// Expected: event was ignored
	}

	// 2. Normal non-injected key event should pass through
	normalKbd := kbdllHookStruct{
		vkCode:   VK_A,
		scanCode: 0x1E,
		flags:    0,
	}
	wh.handleKeyboard(WM_KEYDOWN, unsafe.Pointer(&normalKbd))

	select {
	case ev := <-wh.Events():
		if ev.VKCode != VK_A {
			t.Errorf("expected VK_A, got 0x%02X", ev.VKCode)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for normal event")
	}
}

// TestHook_IncludeInjected verifies the opt-in path: once enabled, keys sent by
// another program are reported and flagged, because AutoHotkey remaps (Total
// Commander's Vim mode, PowerToys, remote desktop) never produce physical
// keystrokes and would otherwise be invisible to the overlay.
func TestHook_IncludeInjected(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific simulation on non-Windows")
	}

	wh := NewHook(16).(*windowsHook)
	wh.SetIncludeInjected(true)

	injectedKbd := kbdllHookStruct{
		vkCode:   VK_DOWN,
		scanCode: 0x50,
		flags:    LLKHF_INJECTED,
	}
	wh.handleKeyboard(WM_KEYDOWN, unsafe.Pointer(&injectedKbd))

	select {
	case ev := <-wh.Events():
		if ev.VKCode != VK_DOWN {
			t.Fatalf("expected VK_DOWN, got 0x%02X", ev.VKCode)
		}
		if !ev.IsInjected {
			t.Error("event should carry IsInjected=true")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for injected event while includeInjected is on")
	}

	// Turning it back off must restore the drop.
	wh.SetIncludeInjected(false)
	wh.handleKeyboard(WM_KEYDOWN, unsafe.Pointer(&injectedKbd))

	select {
	case ev := <-wh.Events():
		t.Fatalf("injected key leaked after disabling: %+v", ev)
	case <-time.After(50 * time.Millisecond):
	}
}

// TestHook_CallbacksRegister verifies both hook thunks can actually be built.
// The hook procedure carries lParam as an unsafe.Pointer so the KBDLLHOOKSTRUCT
// pointer reaches the handler without a uintptr round trip, and
// syscall.NewCallback rejects a signature it cannot marshal by panicking at
// registration time. Without this test that panic would only ever surface on a
// real user's machine, on the first keystroke after launch.
func TestHook_CallbacksRegister(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific hook registration on non-Windows")
	}

	initCallbacks()

	if kbdCallbackPtr == 0 {
		t.Error("keyboard hook callback was not registered")
	}
	if mouseCallbackPtr == 0 {
		t.Error("mouse hook callback was not registered")
	}
}

// TestHook_ModifierAliasCleanup verifies that releasing or changing a modifier properly deletes all alias VKs.
func TestHook_ModifierAliasCleanup(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific simulation on non-Windows")
	}

	wh := NewHook(16).(*windowsHook)

	// 1. Test Ctrl aliases: press VK_LCONTROL and inject VK_CONTROL alias
	wh.processKeyboard(WM_KEYDOWN, VK_LCONTROL, 0x1D, false)
	wh.mu.Lock()
	wh.pressedKeys[VK_CONTROL] = true
	wh.mu.Unlock()

	// Release VK_LCONTROL -> should clear VK_CONTROL, VK_LCONTROL, and VK_RCONTROL
	wh.processKeyboard(WM_KEYUP, VK_LCONTROL, 0x1D, false)
	wh.mu.Lock()
	if wh.pressedKeys[VK_CONTROL] || wh.pressedKeys[VK_LCONTROL] || wh.pressedKeys[VK_RCONTROL] {
		t.Errorf("Ctrl aliases not cleaned up: %+v", wh.pressedKeys)
	}
	if wh.activeMods.HasCtrl() {
		t.Errorf("activeMods should not have ModCtrl after keyup")
	}
	wh.mu.Unlock()

	// 2. Test Alt aliases
	wh.processKeyboard(WM_KEYDOWN, VK_LMENU, 0x38, true)
	wh.mu.Lock()
	wh.pressedKeys[VK_MENU] = true
	wh.mu.Unlock()

	wh.processKeyboard(WM_KEYUP, VK_LMENU, 0x38, false)
	wh.mu.Lock()
	if wh.pressedKeys[VK_MENU] || wh.pressedKeys[VK_LMENU] || wh.pressedKeys[VK_RMENU] {
		t.Errorf("Alt aliases not cleaned up: %+v", wh.pressedKeys)
	}
	if wh.activeMods.HasAlt() {
		t.Errorf("activeMods should not have ModAlt after keyup")
	}
	wh.mu.Unlock()

	// 3. Test Shift aliases
	wh.processKeyboard(WM_KEYDOWN, VK_LSHIFT, 0x2A, false)
	wh.mu.Lock()
	wh.pressedKeys[VK_SHIFT] = true
	wh.mu.Unlock()

	wh.processKeyboard(WM_KEYUP, VK_LSHIFT, 0x2A, false)
	wh.mu.Lock()
	if wh.pressedKeys[VK_SHIFT] || wh.pressedKeys[VK_LSHIFT] || wh.pressedKeys[VK_RSHIFT] {
		t.Errorf("Shift aliases not cleaned up: %+v", wh.pressedKeys)
	}
	if wh.activeMods.HasShift() {
		t.Errorf("activeMods should not have ModShift after keyup")
	}
	wh.mu.Unlock()

	// 4. Test Win aliases
	wh.processKeyboard(WM_KEYDOWN, VK_LWIN, 0x5B, false)
	wh.mu.Lock()
	wh.pressedKeys[VK_RWIN] = true
	wh.mu.Unlock()

	wh.processKeyboard(WM_KEYUP, VK_LWIN, 0x5B, false)
	wh.mu.Lock()
	if wh.pressedKeys[VK_LWIN] || wh.pressedKeys[VK_RWIN] {
		t.Errorf("Win aliases not cleaned up: %+v", wh.pressedKeys)
	}
	if wh.activeMods.HasSuper() {
		t.Errorf("activeMods should not have ModSuper after keyup")
	}
	wh.mu.Unlock()
}

// TestHook_SysKeyAltModifiers verifies that WM_SYSKEYDOWN and WM_SYSKEYUP properly emit ModAlt.
func TestHook_SysKeyAltModifiers(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific simulation on non-Windows")
	}

	wh := NewHook(16).(*windowsHook)

	// Simulate Alt + E key combination via WM_SYSKEYDOWN and WM_SYSKEYUP
	wh.processKeyboard(WM_SYSKEYDOWN, VK_E, 0x12, false)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventKeyDown {
			t.Errorf("expected EventKeyDown, got %v", ev.Type)
		}
		if ev.KeyName != "E" {
			t.Errorf("expected KeyName 'E', got %q", ev.KeyName)
		}
		if !ev.Modifiers.HasAlt() {
			t.Errorf("expected Modifiers to include ModAlt on WM_SYSKEYDOWN, got %v", ev.Modifiers)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for WM_SYSKEYDOWN event")
	}

	wh.processKeyboard(WM_SYSKEYUP, VK_E, 0x12, false)
	select {
	case ev := <-wh.Events():
		if ev.Type != EventKeyUp {
			t.Errorf("expected EventKeyUp, got %v", ev.Type)
		}
		if ev.KeyName != "E" {
			t.Errorf("expected KeyName 'E', got %q", ev.KeyName)
		}
		if !ev.Modifiers.HasAlt() {
			t.Errorf("expected Modifiers to include ModAlt on WM_SYSKEYUP, got %v", ev.Modifiers)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for WM_SYSKEYUP event")
	}
}

// TestHook_DesyncedModifierPhysicalSync verifies that physical hardware checks clean up desynced modifiers.
func TestHook_DesyncedModifierPhysicalSync(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific simulation on non-Windows")
	}

	wh := NewHook(16).(*windowsHook)
	wh.mu.Lock()
	wh.running = true // Enable physical check
	// Manually inject stuck modifiers
	wh.pressedKeys[VK_LCONTROL] = true
	wh.pressedKeys[VK_LMENU] = true
	wh.mu.Unlock()

	// Process normal key press when physical modifiers are NOT held on physical keyboard
	wh.processKeyboard(WM_KEYDOWN, VK_A, 0x1E, false)

	wh.mu.Lock()
	if wh.pressedKeys[VK_LCONTROL] || wh.pressedKeys[VK_CONTROL] {
		t.Errorf("desynced Ctrl was not cleaned up by physical check")
	}
	if wh.pressedKeys[VK_LMENU] || wh.pressedKeys[VK_MENU] {
		t.Errorf("desynced Alt was not cleaned up by physical check")
	}
	wh.mu.Unlock()
}
