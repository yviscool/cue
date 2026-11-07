// Package context tracks the foreground application so semantic subtitles can
// adapt to the app the presenter is actually using.
//
// The provider intentionally lives outside the low-level input hook: Win32
// foreground/process queries are comparatively expensive and must never run on
// the keyboard hook callback thread. The semantic matcher only reads the cached
// AppContext, which keeps the per-keystroke path allocation-free and in the
// nanosecond range.
package context

import (
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Category is a coarse application family used as a fallback dictionary scope.
type Category string

const (
	CategoryUnknown  Category = ""
	CategoryEditor   Category = "editor"
	CategoryBrowser  Category = "browser"
	CategoryTerminal Category = "terminal"
	CategoryOffice   Category = "office"
	CategoryDesign   Category = "design"
	CategoryMedia    Category = "media"
	CategoryChat     Category = "chat"
	CategoryNotes    Category = "notes"
	CategorySystem   Category = "system"
)

// DisplayName returns a human-readable category name.
func (c Category) DisplayName() string {
	switch c {
	case CategoryEditor:
		return "Editor"
	case CategoryBrowser:
		return "Browser"
	case CategoryTerminal:
		return "Terminal"
	case CategoryOffice:
		return "Office"
	case CategoryDesign:
		return "Design"
	case CategoryMedia:
		return "Media"
	case CategoryChat:
		return "Chat"
	case CategoryNotes:
		return "Notes"
	case CategorySystem:
		return "System"
	default:
		return "Unknown"
	}
}

// AppContext is an immutable snapshot of the foreground application.
type AppContext struct {
	HWND       uintptr   `json:"-"`
	PID        uint32    `json:"pid,omitempty"`
	Exe        string    `json:"exe,omitempty"`        // normalized executable name, e.g. "code.exe"
	Title      string    `json:"title,omitempty"`      // foreground window title
	AppID      string    `json:"appId,omitempty"`      // resolved profile id, e.g. "vscode"
	AppName    string    `json:"appName,omitempty"`    // resolved display name, e.g. "VS Code"
	Category   Category  `json:"category,omitempty"`   // coarse family
	Platform   string    `json:"platform,omitempty"`   // "windows", "darwin", "linux"
	Seq        uint64    `json:"seq,omitempty"`        // monotonically increasing change counter
	CapturedAt time.Time `json:"capturedAt,omitempty"` // when this snapshot was taken
}

// IsZero reports whether the snapshot carries no application information.
func (a AppContext) IsZero() bool {
	return a.HWND == 0 && a.PID == 0 && a.Exe == "" && a.Title == "" && a.AppID == ""
}

// Key is a cheap cache key that changes when the meaningful context changes.
func (a AppContext) Key() string {
	if a.Seq != 0 {
		// Seq is the authoritative change marker for provider-backed contexts.
		return a.Exe + "\x00" + a.Title
	}
	return a.Exe + "\x00" + a.Title + "\x00" + a.AppID
}

// DisplayName returns the best available label for the active application.
func (a AppContext) DisplayName() string {
	if a.AppName != "" {
		return a.AppName
	}
	if a.Exe != "" {
		return strings.TrimSuffix(a.Exe, ".exe")
	}
	if a.Title != "" {
		return a.Title
	}
	return ""
}

// NormalizeExe canonicalizes an executable path/name for exact map lookups.
func NormalizeExe(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return ""
	}
	name = filepath.Base(name)
	return name
}

// Provider supplies foreground application snapshots and change notifications.
// Subscribers are invoked from the provider goroutine and must not block; they
// should only copy the snapshot or hand it to a buffered channel.
type Provider interface {
	Start() error
	Stop()
	Current() AppContext
	// Subscribe registers a callback and returns an unsubscribe function.
	Subscribe(fn func(AppContext)) func()
}

// NewProvider creates the platform-native foreground application provider.
// interval controls how often the provider polls; values <= 0 use 250ms.
func NewProvider(interval time.Duration) Provider {
	return newPlatformProvider(interval)
}

// StaticProvider is a deterministic provider used by tests, previews, and
// embedding scenarios. It is safe for concurrent use.
type StaticProvider struct {
	mu      sync.RWMutex
	current AppContext
	seq     uint64
	nextID  uint64
	subs    map[uint64]func(AppContext)
}

// NewStatic creates a provider seeded with an initial snapshot.
func NewStatic(app AppContext) *StaticProvider {
	p := &StaticProvider{subs: make(map[uint64]func(AppContext))}
	p.Set(app)
	return p
}

// Start implements Provider.
func (p *StaticProvider) Start() error { return nil }

// Stop implements Provider.
func (p *StaticProvider) Stop() {
	p.mu.Lock()
	p.subs = make(map[uint64]func(AppContext))
	p.mu.Unlock()
}

// Current implements Provider.
func (p *StaticProvider) Current() AppContext {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.current
}

// Set updates the snapshot and notifies subscribers when it changed.
func (p *StaticProvider) Set(app AppContext) {
	p.mu.Lock()
	p.seq++
	app.Seq = p.seq
	if app.CapturedAt.IsZero() {
		app.CapturedAt = time.Now()
	}
	if p.current.Key() == app.Key() {
		p.mu.Unlock()
		return
	}
	p.current = app
	subs := make([]func(AppContext), 0, len(p.subs))
	for _, fn := range p.subs {
		subs = append(subs, fn)
	}
	p.mu.Unlock()
	for _, fn := range subs {
		fn(app)
	}
}

// Subscribe implements Provider.
func (p *StaticProvider) Subscribe(fn func(AppContext)) func() {
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

// WithProfile decorates a snapshot with a resolved profile/category. It is used
// by the semantic registry so the UI and privacy layer can share one snapshot.
func (a AppContext) WithProfile(id, name string, category Category) AppContext {
	a.AppID = id
	a.AppName = name
	a.Category = category
	return a
}
