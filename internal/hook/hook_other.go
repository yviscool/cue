//go:build !windows

package hook

import (
	"errors"
	"sync"
)

// ErrUnsupportedOS is returned when input hook is started on non-Windows platforms.
var ErrUnsupportedOS = errors.New("input hook is not supported on this operating system")

// QueryModifiers stub for non-Windows platforms.
func QueryModifiers() Modifiers {
	return ModNone
}

type stubHook struct {
	mu         sync.Mutex
	running    bool
	hasStarted bool
	events     chan InputEvent
	stoppedCh  chan struct{}
}

// New creates a new Hook instance for non-Windows environments.
func New() Hook {
	return NewHook(2048)
}

// NewHook creates a new Hook instance with the specified buffer size.
func NewHook(bufferSize ...int) Hook {
	size := 2048
	if len(bufferSize) > 0 && bufferSize[0] > 0 {
		size = bufferSize[0]
	}
	return &stubHook{
		events:    make(chan InputEvent, size),
		stoppedCh: make(chan struct{}),
	}
}

func (h *stubHook) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.running {
		return errors.New("hook is already running")
	}
	if h.hasStarted {
		return errors.New("hook instance cannot be reused after stop; create a new instance with New()")
	}
	h.running = true
	h.hasStarted = true
	return ErrUnsupportedOS
}

func (h *stubHook) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.running {
		return
	}
	h.running = false
	close(h.events)
	close(h.stoppedCh)
}

func (h *stubHook) Events() <-chan InputEvent {
	return h.events
}

// SetIncludeInjected is a no-op: non-Windows stubs never emit events at all.
func (h *stubHook) SetIncludeInjected(bool) {}
