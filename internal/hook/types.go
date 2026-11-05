package hook

import (
	"fmt"
	"strings"
	"time"
)

// EventType distinguishes the category of input events.
type EventType int

const (
	EventKeyDown EventType = iota
	EventKeyUp
	EventMouseDown
	EventMouseUp
	EventMouseMove
	EventMouseWheel
)

func (t EventType) String() string {
	switch t {
	case EventKeyDown:
		return "KeyDown"
	case EventKeyUp:
		return "KeyUp"
	case EventMouseDown:
		return "MouseDown"
	case EventMouseUp:
		return "MouseUp"
	case EventMouseMove:
		return "MouseMove"
	case EventMouseWheel:
		return "MouseWheel"
	default:
		return fmt.Sprintf("Event(%d)", t)
	}
}

// MouseButton identifies a mouse button.
type MouseButton int

const (
	ButtonNone MouseButton = iota
	ButtonLeft
	ButtonRight
	ButtonMiddle
	ButtonExtra1
	ButtonExtra2
)

func (b MouseButton) String() string {
	switch b {
	case ButtonLeft:
		return "Left"
	case ButtonRight:
		return "Right"
	case ButtonMiddle:
		return "Middle"
	case ButtonExtra1:
		return "Extra1"
	case ButtonExtra2:
		return "Extra2"
	default:
		return "None"
	}
}

// Modifiers bitmask for active modifier keys.
type Modifiers uint32

const (
	ModNone  Modifiers = 0
	ModCtrl  Modifiers = 1 << 0
	ModShift Modifiers = 1 << 1
	ModAlt   Modifiers = 1 << 2
	ModSuper Modifiers = 1 << 3 // Win on Windows, Cmd on macOS
)

func (m Modifiers) HasCtrl() bool  { return (m & ModCtrl) != 0 }
func (m Modifiers) HasShift() bool { return (m & ModShift) != 0 }
func (m Modifiers) HasAlt() bool   { return (m & ModAlt) != 0 }
func (m Modifiers) HasSuper() bool { return (m & ModSuper) != 0 }

func (m Modifiers) IsEmpty() bool { return m == 0 }

func (m Modifiers) String() string {
	var parts []string
	if m.HasCtrl() {
		parts = append(parts, "Ctrl")
	}
	if m.HasAlt() {
		parts = append(parts, "Alt")
	}
	if m.HasShift() {
		parts = append(parts, "Shift")
	}
	if m.HasSuper() {
		parts = append(parts, "Win")
	}
	return strings.Join(parts, "+")
}

// InputEvent is the normalized raw input event dispatched by hook.
type InputEvent struct {
	Type      EventType
	Timestamp time.Time

	// Keyboard fields
	VKCode      uint32    // Virtual Key Code (Windows VK or platform key)
	ScanCode    uint32    // Hardware scan code
	KeyName     string    // Human readable key symbol (e.g. "A", "Enter", "F1", "#")
	Modifiers   Modifiers // Active modifier keys snapshot
	IsModifier  bool      // True if the key itself is Ctrl/Shift/Alt/Win
	IsKeyRepeat bool      // True if auto-repeat from keeping key pressed

	// Mouse fields
	Button     MouseButton
	X          int // Absolute screen X
	Y          int // Absolute screen Y
	WheelDelta int // Positive = up/forward, Negative = down/backward
}

// Hook is the interface for listening to OS-level inputs.
type Hook interface {
	Start() error
	Stop()
	Events() <-chan InputEvent
}
