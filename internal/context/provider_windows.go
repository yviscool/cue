//go:build windows

package context

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

const (
	defaultPollInterval = 250 * time.Millisecond
	minPollInterval     = 50 * time.Millisecond
	maxPollInterval     = 2 * time.Second

	processQueryLimitedInformation = 0x1000
	processQueryInformation        = 0x0400
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetForegroundWindow       = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW            = user32.NewProc("GetWindowTextW")
	procGetWindowThreadProcessID  = user32.NewProc("GetWindowThreadProcessId")
	procGetClassNameW             = user32.NewProc("GetClassNameW")
	procEnumChildWindows          = user32.NewProc("EnumChildWindows")
	procOpenProcess               = kernel32.NewProc("OpenProcess")
	procCloseHandle               = kernel32.NewProc("CloseHandle")
	procQueryFullProcessImageName = kernel32.NewProc("QueryFullProcessImageNameW")
)

// UWP child-window enumeration is synchronous, so a single guarded slot is
// enough to carry the result out of the Win32 callback. Using a slot instead of
// lParam avoids an unsafe pointer conversion in the callback.
var (
	enumChildMu  sync.Mutex
	enumChildPID uint32
)

var enumChildCallback = syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
	if enumChildPID != 0 {
		return 0
	}
	var class [128]uint16
	n, _, _ := procGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&class[0])), uintptr(len(class)))
	if n == 0 {
		return 1
	}
	if strings.EqualFold(syscall.UTF16ToString(class[:n]), "Windows.UI.Core.CoreWindow") {
		var pid uint32
		procGetWindowThreadProcessID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
		if pid != 0 {
			enumChildPID = pid
			return 0 // stop enumeration
		}
	}
	return 1
})

type windowsProvider struct {
	interval time.Duration

	mu      sync.RWMutex
	current AppContext
	subs    map[uint64]func(AppContext)
	nextID  uint64
	seq     uint64

	lastHWND  uintptr
	lastTitle string
	lastPID   uint32
	lastExe   string
	ownPID    uint32

	stopCh    chan struct{}
	doneCh    chan struct{}
	startOnce sync.Once
	stopOnce  sync.Once
	started   atomic.Bool
}

func newPlatformProvider(interval time.Duration) Provider {
	if interval <= 0 {
		interval = defaultPollInterval
	}
	if interval < minPollInterval {
		interval = minPollInterval
	}
	if interval > maxPollInterval {
		interval = maxPollInterval
	}
	return &windowsProvider{
		interval: interval,
		subs:     make(map[uint64]func(AppContext)),
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
}

// Start begins foreground polling. It is idempotent.
func (p *windowsProvider) Start() error {
	p.startOnce.Do(func() {
		p.ownPID = uint32(os.Getpid())
		p.started.Store(true)
		p.refresh()
		go p.loop()
	})
	return nil
}

func (p *windowsProvider) loop() {
	defer close(p.doneCh)
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.refresh()
		}
	}
}

// Stop ends polling and waits for the worker goroutine to exit.
func (p *windowsProvider) Stop() {
	p.stopOnce.Do(func() {
		if !p.started.Swap(false) {
			return
		}
		close(p.stopCh)
		<-p.doneCh
	})
}

// Current returns the most recent captured snapshot.
func (p *windowsProvider) Current() AppContext {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.current
}

// Subscribe registers a change callback.
func (p *windowsProvider) Subscribe(fn func(AppContext)) func() {
	if fn == nil {
		return func() {}
	}
	p.mu.Lock()
	id := atomic.AddUint64(&p.nextID, 1)
	p.subs[id] = fn
	p.mu.Unlock()
	return func() {
		p.mu.Lock()
		delete(p.subs, id)
		p.mu.Unlock()
	}
}

func (p *windowsProvider) publish(ctx AppContext) {
	p.mu.Lock()
	p.current = ctx
	subs := make([]func(AppContext), 0, len(p.subs))
	for _, fn := range p.subs {
		subs = append(subs, fn)
	}
	p.mu.Unlock()
	for _, fn := range subs {
		fn(ctx)
	}
}

// refresh compares the foreground window cheaply and only performs the more
// expensive process query when the window handle actually changes.
func (p *windowsProvider) refresh() {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		if p.lastHWND != 0 {
			p.lastHWND = 0
			p.lastTitle = ""
			p.lastPID = 0
			p.lastExe = ""
			p.seq++
			p.publish(AppContext{Platform: "windows", CapturedAt: time.Now(), Seq: p.seq})
		}
		return
	}
	title := windowTitle(hwnd)
	if hwnd == p.lastHWND && title == p.lastTitle {
		return
	}

	// The executable query is the expensive step. It is only repeated when the
	// window handle changes; title-only changes reuse the cached process.
	var pid uint32
	var exe string
	if hwnd == p.lastHWND {
		pid, exe = p.lastPID, p.lastExe
	} else {
		pid, exe = queryForegroundProcess(hwnd)
		p.lastPID, p.lastExe = pid, exe
	}
	p.lastHWND = hwnd
	p.lastTitle = title

	ctx := AppContext{
		HWND:       hwnd,
		PID:        pid,
		Exe:        exe,
		Title:      title,
		Platform:   "windows",
		CapturedAt: time.Now(),
	}
	// Never let Cue's own settings/tray windows shadow the real presenter app.
	if pid != 0 && pid == p.ownPID {
		ctx = AppContext{Platform: "windows", CapturedAt: time.Now()}
	}

	p.seq++
	ctx.Seq = p.seq
	p.publish(ctx)
}

func windowTitle(hwnd uintptr) string {
	var buf [1024]uint16
	n, _, _ := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:n])
}

func queryForegroundProcess(hwnd uintptr) (uint32, string) {
	var pid uint32
	procGetWindowThreadProcessID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return 0, ""
	}

	exe := processExe(pid)
	// UWP apps are hosted by ApplicationFrameHost; resolve the real app process
	// from the CoreWindow child so Office/Store apps still get a profile.
	if strings.EqualFold(filepath.Base(exe), "applicationframehost.exe") {
		if uwpPID := findUWPProcess(hwnd); uwpPID != 0 {
			if realExe := processExe(uwpPID); realExe != "" {
				return uwpPID, NormalizeExe(realExe)
			}
		}
	}
	return pid, NormalizeExe(exe)
}

func findUWPProcess(hwnd uintptr) uint32 {
	enumChildMu.Lock()
	defer enumChildMu.Unlock()
	enumChildPID = 0
	procEnumChildWindows.Call(hwnd, enumChildCallback, 0)
	pid := enumChildPID
	enumChildPID = 0
	return pid
}

func processExe(pid uint32) string {
	handle, _, _ := procOpenProcess.Call(processQueryLimitedInformation, 0, uintptr(pid))
	if handle == 0 {
		// Fall back to the legacy access right for older or elevated targets.
		handle, _, _ = procOpenProcess.Call(processQueryInformation, 0, uintptr(pid))
	}
	if handle == 0 {
		return ""
	}
	defer procCloseHandle.Call(handle)

	var buf [1024]uint16
	size := uint32(len(buf))
	ok, _, _ := procQueryFullProcessImageName.Call(
		handle,
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ok == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:size])
}
