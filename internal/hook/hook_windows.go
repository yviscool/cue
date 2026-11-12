//go:build windows

package hook

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procSetWindowsHookExW   = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procPostThreadMessageW  = user32.NewProc("PostThreadMessageW")
	procPeekMessageW        = user32.NewProc("PeekMessageW")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState")

	procGetCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	procGetModuleHandleW   = kernel32.NewProc("GetModuleHandleW")
)

const (
	WH_KEYBOARD_LL = 13
	WH_MOUSE_LL    = 14

	WM_QUIT = 0x0012

	WM_KEYDOWN    = 0x0100
	WM_KEYUP      = 0x0101
	WM_SYSKEYDOWN = 0x0104
	WM_SYSKEYUP   = 0x0105

	WM_MOUSEMOVE   = 0x0200
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
	WM_MOUSEWHEEL  = 0x020A
	WM_XBUTTONDOWN = 0x020B
	WM_XBUTTONUP   = 0x020C
	WM_MOUSEHWHEEL = 0x020E

	PM_NOREMOVE = 0x0000

	// LLKHF flags in KBDLLHOOKSTRUCT
	LLKHF_EXTENDED = 0x01
	LLKHF_INJECTED = 0x10
	LLKHF_ALTDOWN  = 0x20
	LLKHF_UP       = 0x80

	// XBUTTON identifiers in MSLLHOOKSTRUCT mouseData HIWORD
	XBUTTON1 = 0x0001
	XBUTTON2 = 0x0002

	DefaultBufferSize = 2048
)

// ErrUnsupportedOS is returned on non-Windows operating systems where low-level hook is unsupported.
var ErrUnsupportedOS = errors.New("input hook is not supported on this operating system")

type point struct {
	x int32
	y int32
}

type kbdllHookStruct struct {
	vkCode      uint32
	scanCode    uint32
	flags       uint32
	time        uint32
	dwExtraInfo uintptr
}

type msllHookStruct struct {
	pt          point
	mouseData   uint32
	flags       uint32
	time        uint32
	_           uint32 // 4-byte padding on 64-bit for 8-byte alignment of dwExtraInfo
	dwExtraInfo uintptr
}

type winMsg struct {
	hwnd    uintptr
	message uint32
	_       uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      point
	_       uint32
}

// Global hook registry to avoid exhausting Windows callback slots.
var (
	registryMu       sync.RWMutex
	hookRegistry     = make(map[*windowsHook]struct{})
	initOnce         sync.Once
	kbdCallbackPtr   uintptr
	mouseCallbackPtr uintptr
)

// initCallbacks registers the two low-level hook thunks.
//
// lParam is declared as unsafe.Pointer rather than uintptr on purpose: Windows
// hands the hook procedure a pointer to a KBDLLHOOKSTRUCT / MSLLHOOKSTRUCT, and
// carrying it as a real pointer from the thunk down to the handler avoids the
// uintptr -> unsafe.Pointer round trip that go vet (correctly) flags as unsafe,
// since the GC cannot see through a uintptr and could in principle move the
// object underneath it.
func initCallbacks() {
	initOnce.Do(func() {
		kbdCallbackPtr = syscall.NewCallback(func(nCode int, wParam uintptr, lParam unsafe.Pointer) uintptr {
			if nCode >= 0 {
				registryMu.RLock()
				for h := range hookRegistry {
					h.handleKeyboard(wParam, lParam)
				}
				registryMu.RUnlock()
			}
			r, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, uintptr(lParam))
			return r
		})

		mouseCallbackPtr = syscall.NewCallback(func(nCode int, wParam uintptr, lParam unsafe.Pointer) uintptr {
			if nCode >= 0 {
				registryMu.RLock()
				for h := range hookRegistry {
					h.handleMouse(wParam, lParam)
				}
				registryMu.RUnlock()
			}
			r, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, uintptr(lParam))
			return r
		})
	})
}

func registerHook(h *windowsHook) {
	registryMu.Lock()
	hookRegistry[h] = struct{}{}
	registryMu.Unlock()
}

func unregisterHook(h *windowsHook) {
	registryMu.Lock()
	delete(hookRegistry, h)
	registryMu.Unlock()
}

// isKeyDown checks if a virtual key is currently physically pressed down.
func isKeyDown(vk int) bool {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return int16(r) < 0
}

// QueryModifiers queries the active modifier keys physically pressed down in the OS.
func QueryModifiers() Modifiers {
	var mods Modifiers
	if isKeyDown(VK_CONTROL) || isKeyDown(VK_LCONTROL) || isKeyDown(VK_RCONTROL) {
		mods |= ModCtrl
	}
	if isKeyDown(VK_SHIFT) || isKeyDown(VK_LSHIFT) || isKeyDown(VK_RSHIFT) {
		mods |= ModShift
	}
	if isKeyDown(VK_MENU) || isKeyDown(VK_LMENU) || isKeyDown(VK_RMENU) {
		mods |= ModAlt
	}
	if isKeyDown(VK_LWIN) || isKeyDown(VK_RWIN) {
		mods |= ModSuper
	}
	return mods
}

// windowsHook implements the Hook interface for Windows using pure Win32 API calls.
type windowsHook struct {
	mu         sync.Mutex
	running    bool
	stoppedCh  chan struct{}
	stopOnce   sync.Once
	hasStarted bool
	events     chan InputEvent
	bufferSize int

	threadID uintptr
	hhkKbd   uintptr
	hhkMouse uintptr

	pressedKeys       map[uint32]bool
	activeMods        Modifiers
	lastMouseMoveTime time.Time

	// includeInjected mirrors config.General.ShowInjectedKeys. It is read on the
	// hook thread and written from the UI thread, hence the atomic.
	includeInjected atomic.Bool
}

// New creates a new Windows input Hook with default buffer size.
func New() Hook {
	return NewHook(DefaultBufferSize)
}

// NewHook creates a new Windows input Hook with a specified buffer size.
func NewHook(bufferSize ...int) Hook {
	size := DefaultBufferSize
	if len(bufferSize) > 0 && bufferSize[0] > 0 {
		size = bufferSize[0]
	}
	return &windowsHook{
		bufferSize:  size,
		events:      make(chan InputEvent, size),
		stoppedCh:   make(chan struct{}),
		pressedKeys: make(map[uint32]bool),
	}
}

// Events returns the receive-only channel of dispatched input events.
func (h *windowsHook) Events() <-chan InputEvent {
	return h.events
}

// Start installs the global Windows keyboard and mouse hooks and runs the message loop.
func (h *windowsHook) Start() error {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return errors.New("hook is already running")
	}
	if h.hasStarted {
		h.mu.Unlock()
		return errors.New("hook instance cannot be reused after stop; create a new instance with New()")
	}
	h.running = true
	h.hasStarted = true

	// Sync initial modifier states from the OS
	h.syncInitialModifiersLocked()
	h.mu.Unlock()

	initCallbacks()

	ready := make(chan error, 1)
	go h.messageLoop(ready)

	err := <-ready
	if err != nil {
		h.mu.Lock()
		h.running = false
		h.mu.Unlock()
		return err
	}
	return nil
}

// Stop safely unhooks Windows hooks and terminates the message pump thread.
func (h *windowsHook) Stop() {
	h.stopOnce.Do(func() {
		h.mu.Lock()
		if !h.running {
			h.mu.Unlock()
			return
		}
		h.running = false
		threadID := h.threadID
		h.mu.Unlock()

		if threadID != 0 {
			// Post WM_QUIT to terminate the dedicated message pump
			procPostThreadMessageW.Call(threadID, WM_QUIT, 0, 0)
		}

		// Wait for message loop to exit and clean up
		<-h.stoppedCh
	})
}

// messageLoop runs on a dedicated OS thread locked with runtime.LockOSThread().
func (h *windowsHook) messageLoop(ready chan<- error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tID, _, _ := procGetCurrentThreadId.Call()
	h.mu.Lock()
	h.threadID = tID
	h.mu.Unlock()

	// Force creation of the message queue for this thread
	var dummy winMsg
	procPeekMessageW.Call(uintptr(unsafe.Pointer(&dummy)), 0, 0, 0, PM_NOREMOVE)

	hMod, _, _ := procGetModuleHandleW.Call(0)

	// Install WH_KEYBOARD_LL hook
	hhkKbd, _, errKbd := procSetWindowsHookExW.Call(
		WH_KEYBOARD_LL,
		kbdCallbackPtr,
		hMod,
		0,
	)
	if hhkKbd == 0 {
		ready <- fmt.Errorf("failed to install keyboard hook: %v", errKbd)
		close(h.stoppedCh)
		return
	}
	h.hhkKbd = hhkKbd

	// Install WH_MOUSE_LL hook
	hhkMouse, _, errMouse := procSetWindowsHookExW.Call(
		WH_MOUSE_LL,
		mouseCallbackPtr,
		hMod,
		0,
	)
	if hhkMouse == 0 {
		procUnhookWindowsHookEx.Call(hhkKbd)
		h.hhkKbd = 0
		ready <- fmt.Errorf("failed to install mouse hook: %v", errMouse)
		close(h.stoppedCh)
		return
	}
	h.hhkMouse = hhkMouse

	registerHook(h)
	ready <- nil

	// Win32 message pump
	var m winMsg
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if ret == 0 || int32(ret) == -1 {
			// ret == 0 means WM_QUIT; -1 means error
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}

	// Teardown
	unregisterHook(h)

	if h.hhkKbd != 0 {
		procUnhookWindowsHookEx.Call(h.hhkKbd)
		h.hhkKbd = 0
	}
	if h.hhkMouse != 0 {
		procUnhookWindowsHookEx.Call(h.hhkMouse)
		h.hhkMouse = 0
	}

	// Safely close event stream
	close(h.events)
	close(h.stoppedCh)
}

// syncInitialModifiersLocked syncs modifier state with current OS key states.
func (h *windowsHook) syncInitialModifiersLocked() {
	if isKeyDown(VK_LCONTROL) || isKeyDown(VK_CONTROL) {
		h.pressedKeys[VK_LCONTROL] = true
	}
	if isKeyDown(VK_RCONTROL) {
		h.pressedKeys[VK_RCONTROL] = true
	}
	if isKeyDown(VK_LSHIFT) || isKeyDown(VK_SHIFT) {
		h.pressedKeys[VK_LSHIFT] = true
	}
	if isKeyDown(VK_RSHIFT) {
		h.pressedKeys[VK_RSHIFT] = true
	}
	if isKeyDown(VK_LMENU) || isKeyDown(VK_MENU) {
		h.pressedKeys[VK_LMENU] = true
	}
	if isKeyDown(VK_RMENU) {
		h.pressedKeys[VK_RMENU] = true
	}
	if isKeyDown(VK_LWIN) {
		h.pressedKeys[VK_LWIN] = true
	}
	if isKeyDown(VK_RWIN) {
		h.pressedKeys[VK_RWIN] = true
	}
	h.recomputeModifiersLocked()
}

// recomputeModifiersLocked updates activeMods according to pressedKeys.
func (h *windowsHook) recomputeModifiersLocked() {
	var mods Modifiers
	if h.pressedKeys[VK_LCONTROL] || h.pressedKeys[VK_RCONTROL] || h.pressedKeys[VK_CONTROL] {
		mods |= ModCtrl
	}
	if h.pressedKeys[VK_LSHIFT] || h.pressedKeys[VK_RSHIFT] || h.pressedKeys[VK_SHIFT] {
		mods |= ModShift
	}
	if h.pressedKeys[VK_LMENU] || h.pressedKeys[VK_RMENU] || h.pressedKeys[VK_MENU] {
		mods |= ModAlt
	}
	if h.pressedKeys[VK_LWIN] || h.pressedKeys[VK_RWIN] {
		mods |= ModSuper
	}
	h.activeMods = mods
}

// deleteModifierAliasesLocked removes all virtual key aliases for the modifier group.
func (h *windowsHook) deleteModifierAliasesLocked(vk uint32) {
	switch vk {
	case VK_CONTROL, VK_LCONTROL, VK_RCONTROL:
		delete(h.pressedKeys, VK_CONTROL)
		delete(h.pressedKeys, VK_LCONTROL)
		delete(h.pressedKeys, VK_RCONTROL)
	case VK_MENU, VK_LMENU, VK_RMENU:
		delete(h.pressedKeys, VK_MENU)
		delete(h.pressedKeys, VK_LMENU)
		delete(h.pressedKeys, VK_RMENU)
	case VK_SHIFT, VK_LSHIFT, VK_RSHIFT:
		delete(h.pressedKeys, VK_SHIFT)
		delete(h.pressedKeys, VK_LSHIFT)
		delete(h.pressedKeys, VK_RSHIFT)
	case VK_LWIN, VK_RWIN:
		delete(h.pressedKeys, VK_LWIN)
		delete(h.pressedKeys, VK_RWIN)
	}
}

// syncDesyncedModifiersLocked calls GetAsyncKeyState physical checks to clean up
// any desynced modifier states if physical hardware is up.
func (h *windowsHook) syncDesyncedModifiersLocked() {
	if !h.running {
		return
	}
	// Check Ctrl
	if h.pressedKeys[VK_CONTROL] || h.pressedKeys[VK_LCONTROL] || h.pressedKeys[VK_RCONTROL] {
		if !isKeyDown(VK_CONTROL) && !isKeyDown(VK_LCONTROL) && !isKeyDown(VK_RCONTROL) {
			delete(h.pressedKeys, VK_CONTROL)
			delete(h.pressedKeys, VK_LCONTROL)
			delete(h.pressedKeys, VK_RCONTROL)
		}
	}
	// Check Alt
	if h.pressedKeys[VK_MENU] || h.pressedKeys[VK_LMENU] || h.pressedKeys[VK_RMENU] {
		if !isKeyDown(VK_MENU) && !isKeyDown(VK_LMENU) && !isKeyDown(VK_RMENU) {
			delete(h.pressedKeys, VK_MENU)
			delete(h.pressedKeys, VK_LMENU)
			delete(h.pressedKeys, VK_RMENU)
		}
	}
	// Check Shift
	if h.pressedKeys[VK_SHIFT] || h.pressedKeys[VK_LSHIFT] || h.pressedKeys[VK_RSHIFT] {
		if !isKeyDown(VK_SHIFT) && !isKeyDown(VK_LSHIFT) && !isKeyDown(VK_RSHIFT) {
			delete(h.pressedKeys, VK_SHIFT)
			delete(h.pressedKeys, VK_LSHIFT)
			delete(h.pressedKeys, VK_RSHIFT)
		}
	}
	// Check Win
	if h.pressedKeys[VK_LWIN] || h.pressedKeys[VK_RWIN] {
		if !isKeyDown(VK_LWIN) && !isKeyDown(VK_RWIN) {
			delete(h.pressedKeys, VK_LWIN)
			delete(h.pressedKeys, VK_RWIN)
		}
	}
}

// updateKeyState updates key state and returns repetition status, modifier snapshot, and isModifier.
func (h *windowsHook) updateKeyState(vk uint32, isDown bool, altFlag bool) (bool, Modifiers, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	isRepeat := false
	isMod := IsModifierKey(vk)

	if isDown {
		if h.pressedKeys[vk] {
			isRepeat = true
		} else {
			if isMod {
				h.deleteModifierAliasesLocked(vk)
			}
			h.pressedKeys[vk] = true
		}
	} else {
		delete(h.pressedKeys, vk)
		if isMod {
			h.deleteModifierAliasesLocked(vk)
		}
	}

	h.syncDesyncedModifiersLocked()
	h.recomputeModifiersLocked()
	mods := h.activeMods

	// If ALT flag is set in KBDLLHOOKSTRUCT or from SYS key, ensure ModAlt is active
	// (except when Alt key itself is being released)
	isAltKey := (vk == VK_MENU || vk == VK_LMENU || vk == VK_RMENU)
	if altFlag && !(isAltKey && !isDown) {
		mods |= ModAlt
	}

	return isRepeat, mods, isMod
}

// currentModifiers returns a snapshot of active modifiers.
func (h *windowsHook) currentModifiers() Modifiers {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.syncDesyncedModifiersLocked()
	h.recomputeModifiersLocked()
	return h.activeMods
}

// handleKeyboard processes keyboard hook messages and dispatches InputEvents without blocking.
// lParam is the KBDLLHOOKSTRUCT pointer supplied by the hook thunk.
func (h *windowsHook) handleKeyboard(wParam uintptr, lParam unsafe.Pointer) {
	if lParam == nil {
		return
	}
	kbd := (*kbdllHookStruct)(lParam)
	// Injected synthetic events are dropped unless the user asked for them.
	// Windows itself sometimes synthesizes keys (the Start menu injecting Ctrl+A
	// while Win+Esc is held), which is why they are off by default — but the same
	// flag marks every AutoHotkey remap, so Total Commander's Vim layer, PowerToys
	// key remaps and remote-desktop keystrokes are invisible while it is off.
	injected := (kbd.flags & LLKHF_INJECTED) != 0
	if injected && !h.includeInjected.Load() {
		return
	}
	h.processKeyboardInjected(wParam, kbd.vkCode, kbd.scanCode, (kbd.flags&LLKHF_ALTDOWN) != 0, injected)
}

// SetIncludeInjected reports keys synthesized by other programs when enabled.
func (h *windowsHook) SetIncludeInjected(enabled bool) {
	h.includeInjected.Store(enabled)
}

// processKeyboard handles a physically pressed key; synthetic keys go through
// processKeyboardInjected directly.
func (h *windowsHook) processKeyboard(wParam uintptr, vk uint32, scan uint32, altFlag bool) {
	h.processKeyboardInjected(wParam, vk, scan, altFlag, false)
}

// processKeyboardInjected performs event mapping, repeat checking, and non-blocking event emission.
func (h *windowsHook) processKeyboardInjected(wParam uintptr, vk uint32, scan uint32, altFlag bool, injected bool) {
	var isDown bool
	var evType EventType

	switch wParam {
	case WM_KEYDOWN, WM_SYSKEYDOWN:
		isDown = true
		evType = EventKeyDown
	case WM_KEYUP, WM_SYSKEYUP:
		isDown = false
		evType = EventKeyUp
	default:
		return
	}

	// For WM_SYSKEYDOWN and WM_SYSKEYUP (e.g. Alt + key combinations like Alt + E),
	// ensure Alt is recognized as active even if VK_MENU keydown was missed or message loop intercepted it.
	isSys := (wParam == WM_SYSKEYDOWN || wParam == WM_SYSKEYUP)
	if isSys && vk != VK_F10 {
		altFlag = true
	}

	isRepeat, mods, isMod := h.updateKeyState(vk, isDown, altFlag)

	var keyName string
	if isMod {
		keyName = VKToKeyName(vk)
	} else {
		keyName = ResolveKeyName(vk, mods.HasShift())
	}

	ev := InputEvent{
		Type:        evType,
		Timestamp:   time.Now(),
		VKCode:      vk,
		ScanCode:    scan,
		KeyName:     keyName,
		Modifiers:   mods,
		IsModifier:  isMod,
		IsKeyRepeat: isRepeat,
		IsInjected:  injected,
	}

	// Priority dispatch for keyboard events: never drop keystrokes due to mouse event flooding.
	select {
	case h.events <- ev:
		return
	default:
	}

	// Channel is full; drain one stale event to ensure keystrokes are always recorded.
	select {
	case <-h.events:
	default:
	}

	select {
	case h.events <- ev:
	default:
	}
}

// handleMouse processes mouse hook messages and dispatches InputEvents without blocking.
// lParam is the MSLLHOOKSTRUCT pointer supplied by the hook thunk.
func (h *windowsHook) handleMouse(wParam uintptr, lParam unsafe.Pointer) {
	if lParam == nil {
		return
	}
	ms := (*msllHookStruct)(lParam)
	h.processMouse(wParam, ms.pt.x, ms.pt.y, ms.mouseData)
}

// processMouse normalizes mouse movement, clicks, and wheel events.
func (h *windowsHook) processMouse(wParam uintptr, x int32, y int32, mouseData uint32) {
	// Fast path for high-frequency WM_MOUSEMOVE: throttle rate and avoid heavy modifier polling.
	if wParam == WM_MOUSEMOVE {
		now := time.Now()
		h.mu.Lock()
		if now.Sub(h.lastMouseMoveTime) < 8*time.Millisecond {
			h.mu.Unlock()
			return // Throttle mouse movement to ~125Hz max
		}
		h.lastMouseMoveTime = now
		mods := h.activeMods
		h.mu.Unlock()

		ev := InputEvent{
			Type:       EventMouseMove,
			Timestamp:  now,
			Button:     ButtonNone,
			X:          int(x),
			Y:          int(y),
			WheelDelta: 0,
			Modifiers:  mods,
		}
		select {
		case h.events <- ev:
		default:
		}
		return
	}

	mods := h.currentModifiers()

	var evType EventType
	var btn MouseButton
	var wheelDelta int

	switch wParam {
	case WM_LBUTTONDOWN:
		evType = EventMouseDown
		btn = ButtonLeft
	case WM_LBUTTONUP:
		evType = EventMouseUp
		btn = ButtonLeft

	case WM_RBUTTONDOWN:
		evType = EventMouseDown
		btn = ButtonRight
	case WM_RBUTTONUP:
		evType = EventMouseUp
		btn = ButtonRight

	case WM_MBUTTONDOWN:
		evType = EventMouseDown
		btn = ButtonMiddle
	case WM_MBUTTONUP:
		evType = EventMouseUp
		btn = ButtonMiddle

	case WM_XBUTTONDOWN:
		evType = EventMouseDown
		xBtn := (mouseData >> 16) & 0xFFFF
		if xBtn == XBUTTON1 {
			btn = ButtonExtra1
		} else if xBtn == XBUTTON2 {
			btn = ButtonExtra2
		}

	case WM_XBUTTONUP:
		evType = EventMouseUp
		xBtn := (mouseData >> 16) & 0xFFFF
		if xBtn == XBUTTON1 {
			btn = ButtonExtra1
		} else if xBtn == XBUTTON2 {
			btn = ButtonExtra2
		}

	case WM_MOUSEWHEEL, WM_MOUSEHWHEEL:
		evType = EventMouseWheel
		btn = ButtonNone
		wheelDelta = int(int16(mouseData >> 16))

	default:
		return
	}

	ev := InputEvent{
		Type:       evType,
		Timestamp:  time.Now(),
		Button:     btn,
		X:          int(x),
		Y:          int(y),
		WheelDelta: wheelDelta,
		Modifiers:  mods,
	}

	// Priority dispatch for mouse clicks/wheel
	select {
	case h.events <- ev:
		return
	default:
	}

	select {
	case <-h.events:
	default:
	}

	select {
	case h.events <- ev:
	default:
	}
}
