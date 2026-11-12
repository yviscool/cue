package engine

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/semantic"
)

const (
	defaultDecayInterval    = 50 * time.Millisecond
	defaultLingerShort      = 1800 * time.Millisecond
	defaultLingerLong       = 2500 * time.Millisecond
	defaultMouseClickLinger = 180 * time.Millisecond
	defaultMouseBadgeLinger = 1200 * time.Millisecond
	defaultMaxChords        = 5
	defaultThemeName        = "Default"
)

// The five windows that decide when two inputs are one gesture. They are kept
// together because tuning one in isolation is how you break another: a word
// window longer than the streak window makes letters merge into a single press,
// and a scroll window as long as the streak window lets a flick survive a pause
// that should have ended it.
const (
	// wordGroupingWindow: how soon a typed character may follow the previous one
	// and still join the same word chip.
	wordGroupingWindow = 600 * time.Millisecond

	// defaultStreakThreshold: how soon the same key may be pressed again and still
	// read as a repeat of the same action rather than a new cue.
	defaultStreakThreshold = 1200 * time.Millisecond

	// sequenceWindow: how long the second chord of a two-chord shortcut may take.
	// Ctrl+K then Ctrl+S is one intent only while the first chord is still fresh.
	sequenceWindow = semantic.DefaultSequenceWindow

	// doubleClickWindow: how far apart two clicks may be and still count as a
	// double click. Matches the Windows default so the badge agrees with the OS.
	doubleClickWindow = 400 * time.Millisecond

	// scrollMergeWindow: how soon a wheel notch may follow the previous one and
	// still belong to the same scroll gesture. Deliberately shorter than the
	// streak window — a flick is one motion, two separate scrolls are two.
	scrollMergeWindow = 250 * time.Millisecond
)

// Engine is the central input state machine and cadence decay engine.
type Engine struct {
	mu sync.RWMutex

	// Active configuration
	filterMode       FilterMode
	allowedKeys      map[string]bool // normalized set of allowed keys for FilterCustom mode
	decayInterval    time.Duration
	streakThreshold  time.Duration
	lingerShort      time.Duration
	lingerLong       time.Duration
	maxChords        int
	soundPlayer      SoundPlayer
	soundActive      bool
	matcher          semantic.Matcher
	appContext       context.AppContext
	semanticContext  bool
	themeName        string
	showMouseClicks  bool
	mouseClickLinger time.Duration
	mouseBadgeLinger time.Duration
	dragThreshold    int
	wordGrouping     bool
	wordBreakPending bool
	privacyMode      bool

	// Current state
	heldKeys        map[string]bool // normalized key names currently pressed
	chords          []DisplayChord  // active display chords on screen
	mouseState      MouseVisualState
	pressedButton   hook.MouseButton
	lastClickTime   time.Time
	lastClickButton hook.MouseButton
	clickCount      int
	dragStartX      int
	dragStartY      int
	isDragging      bool
	lastChordKey    string
	lastChordMods   hook.Modifiers
	chordSeq        uint64
	recentStrokes   []time.Time // timestamps of recent keystrokes within sliding window
	currentWPM      int         // instantaneous words per minute

	// Two-chord sequence tracking (e.g. Ctrl+K Ctrl+S). The first chord is
	// remembered for a short window; if the second arrives in time and the
	// matcher knows the pair, the two cues collapse into one.
	seqMods hook.Modifiers
	seqKey  string
	seqAt   time.Time

	// Mouse retirement timers
	mouseRetireTimer *time.Timer
	mouseBadgeTimer  *time.Timer

	// Decay cadence loop control
	stopCh    chan struct{}
	startOnce sync.Once
	stopOnce  sync.Once
	wg        sync.WaitGroup

	// Subscription management
	subMu       sync.RWMutex
	subscribers []chan UIState
}

// Option allows configuring Engine upon creation.
type Option func(*Engine)

// WithFilterMode sets initial filter mode.
func WithFilterMode(m FilterMode) Option {
	return func(e *Engine) {
		e.filterMode = m
	}
}

// WithAllowedKeys sets the initial whitelist of allowed keys for FilterCustom mode.
func WithAllowedKeys(keys []string) Option {
	return func(e *Engine) {
		e.allowedKeys = make(map[string]bool, len(keys))
		for _, k := range keys {
			norm := normalizeAllowedKey(k)
			if norm != "" {
				e.allowedKeys[norm] = true
			}
		}
	}
}

// WithDecayInterval sets periodic decay interval.
func WithDecayInterval(d time.Duration) Option {
	return func(e *Engine) {
		if d > 0 {
			e.decayInterval = d
		}
	}
}

// WithStreakThreshold sets streak repeat window.
func WithStreakThreshold(d time.Duration) Option {
	return func(e *Engine) {
		if d > 0 {
			e.streakThreshold = d
		}
	}
}

// WithLingerShort sets linger time for single short keypresses.
func WithLingerShort(d time.Duration) Option {
	return func(e *Engine) {
		if d > 0 {
			e.lingerShort = d
		}
	}
}

// WithLingerLong sets linger time for shortcuts with action subtitles.
func WithLingerLong(d time.Duration) Option {
	return func(e *Engine) {
		if d > 0 {
			e.lingerLong = d
		}
	}
}

// WithMaxChords sets maximum number of on-screen display chords.
func WithMaxChords(n int) Option {
	return func(e *Engine) {
		if n > 0 {
			e.maxChords = n
		}
	}
}

// WithSoundPlayer attaches an acoustic feedback player.
func WithSoundPlayer(p SoundPlayer) Option {
	return func(e *Engine) {
		e.soundPlayer = p
	}
}

// WithSoundActive enables or disables acoustic feedback.
func WithSoundActive(active bool) Option {
	return func(e *Engine) {
		e.soundActive = active
	}
}

// WithMatcher attaches a semantic matcher for shortcuts and privacy.
func WithMatcher(m semantic.Matcher) Option {
	return func(e *Engine) {
		e.matcher = m
	}
}

// WithAppContext sets initial application context.
func WithAppContext(ctx context.AppContext) Option {
	return func(e *Engine) {
		e.appContext = ctx
	}
}

// WithSemanticContext enables or disables context-aware shortcut interpretation.
func WithSemanticContext(enabled bool) Option {
	return func(e *Engine) {
		e.semanticContext = enabled
	}
}

// WithThemeName sets initial UI theme name.
func WithThemeName(name string) Option {
	return func(e *Engine) {
		e.themeName = name
	}
}

// WithShowMouseClicks configures whether mouse clicks produce display chords.
func WithShowMouseClicks(show bool) Option {
	return func(e *Engine) {
		e.showMouseClicks = show
	}
}

// WithMouseClickLinger sets the linger duration for mouse button click highlight and ripple.
func WithMouseClickLinger(d time.Duration) Option {
	return func(e *Engine) {
		if d > 0 {
			e.mouseClickLinger = d
		}
	}
}

// WithMouseBadgeLinger sets the linger duration for multi-click streak badge.
func WithMouseBadgeLinger(d time.Duration) Option {
	return func(e *Engine) {
		if d > 0 {
			e.mouseBadgeLinger = d
		}
	}
}

// WithDragThreshold sets drag movement sensitivity threshold.
func WithDragThreshold(th int) Option {
	return func(e *Engine) {
		e.dragThreshold = th
	}
}

// WithWordGrouping enables or disables word grouping for plain typed text.
func WithWordGrouping(grouping bool) Option {
	return func(e *Engine) {
		e.wordGrouping = grouping
	}
}

// WithPrivacyMode enables the privacy shield: keystrokes typed into password,
// login, or banking windows are replaced by a single placeholder. Off by
// default — masking is a strong claim about what the overlay may reveal, so it
// is opt-in rather than something a new install silently does.
func WithPrivacyMode(enabled bool) Option {
	return func(e *Engine) {
		e.privacyMode = enabled
	}
}

// NewEngine creates and initializes a Engine instance.
func NewEngine(opts ...Option) *Engine {
	e := &Engine{
		filterMode:       FilterSmart,
		decayInterval:    defaultDecayInterval,
		streakThreshold:  defaultStreakThreshold,
		lingerShort:      defaultLingerShort,
		lingerLong:       defaultLingerLong,
		mouseClickLinger: defaultMouseClickLinger,
		mouseBadgeLinger: defaultMouseBadgeLinger,
		maxChords:        defaultMaxChords,
		soundActive:      true,
		matcher:          nil,
		semanticContext:  true,
		themeName:        defaultThemeName,
		showMouseClicks:  false,
		dragThreshold:    50,
		heldKeys:         make(map[string]bool),
		allowedKeys:      make(map[string]bool),
		stopCh:           make(chan struct{}),
	}

	for _, opt := range opts {
		opt(e)
	}

	e.Start()
	return e
}

// SetShowMouseClicks toggles on-screen visualization of mouse clicks.
func (e *Engine) SetShowMouseClicks(show bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.showMouseClicks = show
}

// SetMouseClickLinger updates the linger duration for mouse button highlight and ripple.
func (e *Engine) SetMouseClickLinger(ms int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if ms <= 0 {
		ms = int(defaultMouseClickLinger / time.Millisecond)
	}
	if ms < 40 {
		ms = 40
	}
	if ms > 2000 {
		ms = 2000
	}
	e.mouseClickLinger = time.Duration(ms) * time.Millisecond
}

// SetMouseBadgeLinger updates the linger duration for multi-click streak badge.
func (e *Engine) SetMouseBadgeLinger(ms int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if ms <= 0 {
		ms = int(defaultMouseBadgeLinger / time.Millisecond)
	}
	if ms < 300 {
		ms = 300
	}
	if ms > 5000 {
		ms = 5000
	}
	e.mouseBadgeLinger = time.Duration(ms) * time.Millisecond
}

// SetDragThreshold sets minimum displacement pixels before triggering mouse drag.
func (e *Engine) SetDragThreshold(th int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.dragThreshold = th
}

// SetWordGrouping enables or disables word grouping at runtime and emits updated UIState.
func (e *Engine) SetWordGrouping(grouping bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.wordGrouping = grouping
	e.emitStateLocked()
}

// SetPrivacyMode turns the privacy shield on or off at runtime. Changing it
// affects subsequent keystrokes only; cues already on screen keep whatever they
// were showing, because rewriting history would be its own privacy leak.
func (e *Engine) SetPrivacyMode(enabled bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.privacyMode = enabled
	e.emitStateLocked()
}

// SetLingerDuration applies the configured on-screen lifetime in milliseconds.
// Short keys use the configured duration; chords with a semantic action stay a
// little longer so viewers have time to read the intent.
func (e *Engine) SetLingerDuration(ms int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if ms <= 0 {
		ms = int(defaultLingerShort / time.Millisecond)
	}
	if ms < 300 {
		ms = 300
	}
	if ms > 15000 {
		ms = 15000
	}
	e.lingerShort = time.Duration(ms) * time.Millisecond
	e.lingerLong = e.lingerShort + 700*time.Millisecond
	if e.lingerLong > 20*time.Second {
		e.lingerLong = 20 * time.Second
	}
	e.emitStateLocked()
}

// SetMaxChords caps how many chord clusters stay on screen at once.
func (e *Engine) SetMaxChords(n int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if n < 1 {
		n = 1
	}
	if n > 20 {
		n = 20
	}
	e.maxChords = n
	if len(e.chords) > n {
		e.chords = e.chords[len(e.chords)-n:]
		if len(e.chords) == 0 {
			e.lastChordKey = ""
			e.lastChordMods = hook.ModNone
		} else {
			last := e.chords[len(e.chords)-1]
			if len(last.Keys) > 0 {
				e.lastChordKey = last.Keys[len(last.Keys)-1].Label
			}
			e.lastChordMods = extractModifiers(last)
		}
	}
	e.emitStateLocked()
}

// WordGrouping returns whether word grouping is enabled.
func (e *Engine) WordGrouping() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.wordGrouping
}

// Start launches the cadence decay background loop. Safe to call multiple times.
func (e *Engine) Start() {
	e.startOnce.Do(func() {
		if e.stopCh == nil {
			e.stopCh = make(chan struct{})
		}
		e.wg.Add(1)
		go e.decayLoop()
	})
}

// Stop shuts down the background cadence decay loop and releases resources.
func (e *Engine) Stop() {
	e.stopOnce.Do(func() {
		if e.stopCh != nil {
			close(e.stopCh)
		}
		e.wg.Wait()

		e.mu.Lock()
		if e.mouseRetireTimer != nil {
			e.mouseRetireTimer.Stop()
			e.mouseRetireTimer = nil
		}
		if e.mouseBadgeTimer != nil {
			e.mouseBadgeTimer.Stop()
			e.mouseBadgeTimer = nil
		}
		e.mu.Unlock()

		e.subMu.Lock()
		for _, ch := range e.subscribers {
			close(ch)
		}
		e.subscribers = nil
		e.subMu.Unlock()
	})
}

// decayLoop runs periodically to clean up expired chords.
func (e *Engine) decayLoop() {
	defer e.wg.Done()
	ticker := time.NewTicker(e.decayInterval)
	defer ticker.Stop()

	for {
		select {
		case <-e.stopCh:
			return
		case now := <-ticker.C:
			e.DecayCheck(now)
		}
	}
}

// scheduleMouseRetirementLocked programs precise high-resolution timers for:
// 1) Click linger expiration (turns off lobe highlight and ripple)
// 2) Multi-click badge expiration (dismisses the ×2 streak badge and mouse silhouette)
func (e *Engine) scheduleMouseRetirementLocked() {
	if e.mouseRetireTimer != nil {
		e.mouseRetireTimer.Stop()
		e.mouseRetireTimer = nil
	}
	if e.mouseBadgeTimer != nil {
		e.mouseBadgeTimer.Stop()
		e.mouseBadgeTimer = nil
	}

	linger := e.mouseClickLinger
	if linger <= 0 {
		linger = defaultMouseClickLinger
	}

	ts := e.mouseState.Timestamp
	clickCount := e.mouseState.ClickCount
	isWheel := e.mouseState.Wheel != 0

	if clickCount > 1 && !isWheel {
		// Multi-click streak:
		// Step 1: When click linger expires, refresh so the button lobe returns to idle outline
		e.mouseRetireTimer = time.AfterFunc(linger, func() {
			e.mu.Lock()
			defer e.mu.Unlock()
			if e.mouseState.Timestamp.Equal(ts) && !e.mouseState.Active {
				e.emitStateLocked()
			}
		})

		// Step 2: When badge linger expires, fully retire mouse state and badge
		badgeLinger := e.mouseBadgeLinger
		if badgeLinger <= linger {
			badgeLinger = defaultMouseBadgeLinger
		}
		e.mouseBadgeTimer = time.AfterFunc(badgeLinger, func() {
			e.mu.Lock()
			defer e.mu.Unlock()
			if e.mouseState.Timestamp.Equal(ts) && !e.mouseState.Active {
				e.mouseState.Active = false
				e.mouseState.Dragging = false
				e.mouseState.Wheel = 0
				e.mouseState.Button = hook.ButtonNone
				e.mouseState.Timestamp = time.Time{}
				e.mouseState.ClickCount = 0
				e.emitStateLocked()
			}
		})
	} else {
		// Single click or wheel scroll: retire completely at click linger
		e.mouseRetireTimer = time.AfterFunc(linger, func() {
			e.mu.Lock()
			defer e.mu.Unlock()
			if e.mouseState.Timestamp.Equal(ts) && !e.mouseState.Active {
				e.mouseState.Active = false
				e.mouseState.Dragging = false
				e.mouseState.Wheel = 0
				e.mouseState.Button = hook.ButtonNone
				e.mouseState.Timestamp = time.Time{}
				e.mouseState.ClickCount = 0
				e.emitStateLocked()
			}
		})
	}
}

// DecayCheck evaluates all on-screen chords against current time and prunes expired ones.
// It also retires stale mouse feedback so click rings and wheel indicators do not stay
// frozen on screen after the user stops interacting. Returns the number of removed chords.
func (e *Engine) DecayCheck(now time.Time) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	removed := 0
	if len(e.chords) > 0 {
		active := e.chords[:0]
		for _, chord := range e.chords {
			if now.After(chord.ExpiresAt) {
				removed++
			} else {
				active = append(active, chord)
			}
		}
		if removed > 0 {
			e.chords = active
			if len(e.chords) == 0 {
				e.lastChordKey = ""
				e.lastChordMods = hook.ModNone
			} else {
				last := e.chords[len(e.chords)-1]
				if len(last.Keys) > 0 {
					e.lastChordKey = last.Keys[len(last.Keys)-1].Label
				}
				e.lastChordMods = extractModifiers(last)
			}
		}
	}

	// Mouse feedback is time-based: a click ring lives mouseClickLinger (e.g. 180 ms)
	// and a multi-click streak badge lives mouseBadgeLinger (e.g. 1200 ms).
	// Without this cleanup the last frame would remain painted until the next input event.
	mouseRetired := false
	if !e.mouseState.Timestamp.IsZero() && !e.mouseState.Active {
		age := now.Sub(e.mouseState.Timestamp)
		linger := e.mouseClickLinger
		if linger <= 0 {
			linger = defaultMouseClickLinger
		}
		badgeLinger := e.mouseBadgeLinger
		if badgeLinger <= linger {
			badgeLinger = defaultMouseBadgeLinger
		}
		expireDuration := linger
		if e.mouseState.ClickCount > 1 && e.mouseState.Wheel == 0 {
			expireDuration = badgeLinger
		}
		if age >= expireDuration {
			e.mouseState.Active = false
			e.mouseState.Dragging = false
			e.mouseState.Wheel = 0
			e.mouseState.Button = hook.ButtonNone
			e.mouseState.Timestamp = time.Time{}
			e.mouseState.ClickCount = 0
			mouseRetired = true
		}
	}

	wpmBefore := e.currentWPM
	e.decayCadenceLocked(now)
	wpmDecayed := wpmBefore > 0 && e.currentWPM == 0

	if removed > 0 || mouseRetired || wpmDecayed {
		e.emitStateLocked()
	}

	return removed
}

// Clear immediately clears all on-screen keys and mouse visual states.
func (e *Engine) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.mouseRetireTimer != nil {
		e.mouseRetireTimer.Stop()
		e.mouseRetireTimer = nil
	}
	if e.mouseBadgeTimer != nil {
		e.mouseBadgeTimer.Stop()
		e.mouseBadgeTimer = nil
	}

	e.chords = nil
	e.mouseState = MouseVisualState{}
	e.clickCount = 0
	e.lastChordKey = ""
	e.lastChordMods = hook.ModNone
	e.seqAt = time.Time{}
	e.seqKey = ""
	e.seqMods = hook.ModNone
	e.emitStateLocked()
}

// HandleEvent processes normalized raw input events from hook.
func (e *Engine) HandleEvent(ev hook.InputEvent) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch ev.Type {
	case hook.EventKeyDown:
		e.handleKeyDown(ev)
	case hook.EventKeyUp:
		e.handleKeyUp(ev)
	case hook.EventMouseDown, hook.EventMouseUp, hook.EventMouseMove, hook.EventMouseWheel:
		e.handleMouseEvent(ev)
	}
}

// handleKeyDown processes key press events.
func (e *Engine) handleKeyDown(ev hook.InputEvent) {
	now := ev.Timestamp
	if now.IsZero() {
		now = time.Now()
	}

	normKey := normalizeKeyName(ev.KeyName)
	isMod := ev.IsModifier || isModifierKey(ev.KeyName)

	// If key is a modifier alone, track as active held modifier and trigger Thock feedback
	if isMod {
		e.heldKeys[normKey] = true
		e.playAcoustic(SoundThock, normKey)
		return
	}

	// Non-modifier key is pressed
	e.heldKeys[normKey] = true

	// Effective modifiers calculation (hook snapshot or derived from held keys)
	mods := ev.Modifiers
	if mods.IsEmpty() {
		mods = e.activeModifiersFromHeld()
	}

	// Filter checks
	if e.filterMode == FilterSilent {
		return
	}
	if e.filterMode == FilterCustom && !e.isCustomAllowedLocked(ev.KeyName, ev.VKCode, mods) {
		return
	}

	// Word Grouping: Backspace handling
	if e.wordGrouping && isBackspaceKey(ev.KeyName, ev.VKCode) && mods == hook.ModNone {
		if len(e.chords) > 0 && isTextWordChord(e.chords[len(e.chords)-1]) {
			lastChord := &e.chords[len(e.chords)-1]
			runes := []rune(lastChord.Keys[0].Label)
			if len(runes) > 1 {
				lastChord.Keys[0].Label = string(runes[:len(runes)-1])
				lastChord.ExpiresAt = now.Add(e.lingerShort)
				lastChord.Timestamp = now
				e.lastChordKey = lastChord.Keys[0].Label
			} else {
				e.chords = e.chords[:len(e.chords)-1]
				if len(e.chords) == 0 {
					e.lastChordKey = ""
					e.lastChordMods = hook.ModNone
				} else {
					last := e.chords[len(e.chords)-1]
					if len(last.Keys) > 0 {
						e.lastChordKey = last.Keys[len(last.Keys)-1].Label
					}
					e.lastChordMods = extractModifiers(last)
				}
			}
			e.recordStrokeLocked(now)
			e.playAcoustic(SoundClick, "Backspace")
			e.emitStateLocked()
			return
		}
	}

	// Word Grouping: Alphanumeric grouping
	if e.wordGrouping {
		if charStr, ok := isWordGroupingChar(ev.KeyName, ev.VKCode, mods); ok {
			if len(e.chords) > 0 && isTextWordChord(e.chords[len(e.chords)-1]) && !e.wordBreakPending {
				lastChord := &e.chords[len(e.chords)-1]
				if now.Sub(lastChord.Timestamp) <= wordGroupingWindow {
					lastChord.Keys[0].Label = appendWordRune(lastChord.Keys[0].Label, charStr)
					lastChord.ExpiresAt = now.Add(e.lingerShort)
					lastChord.Timestamp = now
					lastChord.IsActive = true
					lastChord.Keys[0].IsHeld = true
					wpm := e.recordStrokeLocked(now)
					lastChord.Velocity = velocityForWPM(wpm)
					e.lastChordKey = lastChord.Keys[0].Label
					e.lastChordMods = hook.ModNone
					e.playAcoustic(SoundClick, charStr)
					e.emitStateLocked()
					return
				}
			}
			hadBreak := e.wordBreakPending
			e.wordBreakPending = false
			if hadBreak {
				// The key that ended the previous word may never have reached the
				// display — the smart filter drops spaces and plain punctuation —
				// but it still ended the word. Clearing the repeat chain is what
				// keeps an invisible key from leaving a visible trace: without it
				// `x=x` collapses into a single `x ×2` cue.
				e.lastChordKey = ""
				e.lastChordMods = hook.ModNone
			}
			// Initial character or new word after break: push as new word chord without modifier badges
			e.pushChordLocked(now, hook.ModNone, charStr, true, ev.IsKeyRepeat)
			return
		} else {
			// Non-alphanumeric character encountered (space, symbol, enter, punctuation, etc.)
			// Signals that the previous alphanumeric word is completed and the next will start afresh.
			e.wordBreakPending = true
		}
	}

	if e.filterMode == FilterSmart {
		if mods == hook.ModNone && !isSpecialKey(ev.KeyName, ev.VKCode) {
			return
		}
	}
	// FilterAll displays all keys

	// Backspace absorption in FilterAll mode
	if e.filterMode == FilterAll && isBackspaceKey(ev.KeyName, ev.VKCode) && mods == hook.ModNone {
		if len(e.chords) > 0 && isAbsorbableCharacter(e.chords[len(e.chords)-1]) {
			lastChord := &e.chords[len(e.chords)-1]
			if lastChord.Streak > 1 {
				lastChord.Streak--
				lastChord.ExpiresAt = now.Add(e.lingerShort)
			} else {
				e.chords = e.chords[:len(e.chords)-1]
				if len(e.chords) == 0 {
					e.lastChordKey = ""
					e.lastChordMods = hook.ModNone
				} else {
					last := e.chords[len(e.chords)-1]
					if len(last.Keys) > 0 {
						e.lastChordKey = last.Keys[len(last.Keys)-1].Label
					}
					e.lastChordMods = extractModifiers(last)
				}
			}
			e.recordStrokeLocked(now)
			e.playAcoustic(SoundClick, "Backspace")
			e.emitStateLocked()
			return
		}
	}

	e.pushChordLocked(now, mods, ev.KeyName, true, ev.IsKeyRepeat)
}

// pushChordLocked constructs a new DisplayChord or increments streak for repetitive inputs.
// isRepeat marks an OS auto-repeat event: the user pressed once and held, so the
// repeat must keep the cue alive without pretending it was another press.
func (e *Engine) pushChordLocked(now time.Time, mods hook.Modifiers, keyName string, isHeld, isRepeat bool) {
	wpm := e.recordStrokeLocked(now)
	velocity := velocityForWPM(wpm)
	isFinisher := isFinisherKey(keyName)

	// Privacy is evaluated once per keystroke: it feeds the sequence guard, the
	// placeholder mask, and the matcher's typed-pattern history.
	isPrivate := false
	if e.privacyMode && e.matcher != nil {
		isPrivate = e.matcher.IsSensitive(e.appContext, keyName)
	}

	// ---- Two-chord shortcut sequences -------------------------------------
	// Ctrl+K Ctrl+S style shortcuts are common in editors. If the previous
	// chord is still inside the sequence window and the matcher knows the pair,
	// collapse both chords into one cue instead of showing two unrelated ones.
	//
	// Both chords must carry a modifier: FilterSmart drops bare keystrokes
	// before they reach this point, so a sequence ending in a plain key (VS
	// Code's Ctrl+K S) cannot be collapsed yet.
	if !isPrivate && e.semanticContext && e.matcher != nil && mods != hook.ModNone &&
		!e.seqAt.IsZero() && now.Sub(e.seqAt) >= 0 && now.Sub(e.seqAt) <= sequenceWindow {
		if res, ok := e.matcher.MatchSequence(e.appContext, e.seqMods, e.seqKey, mods, keyName); ok {
			keys := buildKeyCaps(e.seqMods, e.seqKey, true)
			keys = append(keys, KeyCap{Label: "→", IsModifier: true})
			keys = append(keys, buildKeyCaps(mods, keyName, isHeld)...)

			// Drop the pending first-chord cue so the combined form replaces it.
			if n := len(e.chords); n > 0 {
				last := e.chords[n-1]
				if !last.Timestamp.Before(e.seqAt.Add(-5 * time.Millisecond)) {
					e.chords = e.chords[:n-1]
				}
			}

			chord := DisplayChord{
				ID:         fmt.Sprintf("chord-%d-%d", now.UnixNano(), atomic.AddUint64(&e.chordSeq, 1)),
				Keys:       keys,
				Streak:     1,
				Action:     res.Action,
				Timestamp:  now,
				ExpiresAt:  now.Add(e.lingerLong),
				IsActive:   isHeld,
				Velocity:   velocity,
				IsFinisher: isFinisher,
				Confidence: res.Confidence,
				AppName:    res.AppName,
				IsSequence: true,
			}
			e.chords = append(e.chords, chord)
			if len(e.chords) > e.maxChords {
				e.chords = e.chords[len(e.chords)-e.maxChords:]
			}
			e.lastChordKey = keyName
			e.lastChordMods = mods
			e.seqAt = time.Time{}
			e.seqKey = ""
			e.seqMods = hook.ModNone
			e.playAcoustic(SoundClick, keyName)
			e.emitStateLocked()
			return
		}
	}

	// Streak check: identical to the most recent chord (same modifiers + same key) within threshold.
	// Plain word characters never reach this point — the grouping branch above
	// returns for them — so nothing here needs to be excluded from streaking.
	// The guard that used to sit here (`!e.wordGrouping || mods != ModNone`)
	// therefore excluded exactly the wrong keys: with word grouping on, every
	// unmodified key lost its repeat count, so three Right presses or three
	// Backspaces showed as three separate cues instead of one ×3.
	// NOTE: Scroll wheel events are continuous fluid motions and must never accumulate streak badges.
	// NOTE: Auto-repeat is the machine pressing the key, not the user. It refreshes
	// the cue (so a held key stays lit) but never raises the count, otherwise a
	// two-second hold on an arrow reads as ×60.
	isScroll := strings.EqualFold(keyName, "Scroll Down") || strings.EqualFold(keyName, "Scroll Up")
	if len(e.chords) > 0 && e.lastChordMods == mods && strings.EqualFold(e.lastChordKey, keyName) {
		lastChord := &e.chords[len(e.chords)-1]
		if now.Sub(lastChord.Timestamp) <= e.streakThreshold {
			if !isScroll && !isRepeat {
				lastChord.Streak++
			}
			lastChord.Velocity = velocity
			duration := e.lingerShort
			if lastChord.Action.Title != "" {
				duration = e.lingerLong
			}
			lastChord.ExpiresAt = now.Add(duration)
			lastChord.Timestamp = now
			lastChord.IsActive = isHeld
			for i := range lastChord.Keys {
				lastChord.Keys[i].IsHeld = isHeld
			}
			e.playAcoustic(SoundClick, keyName)
			e.emitStateLocked()
			return
		}
	}

	// Construct DisplayChord
	keys := buildKeyCaps(mods, keyName, isHeld)

	// Resolve the semantic interpretation for one chord. When context
	// intelligence is disabled, universal decoding and global overrides still
	// work; only app profiles and categories are skipped.
	var result semantic.MatchResult
	if e.matcher != nil {
		if e.semanticContext {
			result = e.matcher.MatchAction(e.appContext, mods, keyName)
		} else {
			result = e.matcher.MatchUniversal(mods, keyName)
		}
	}
	action := result.Action

	// Privacy mask check: sensitive contexts replace the physical key with a
	// single protected placeholder.
	if isPrivate {
		keys = []KeyCap{{Label: "[🔒 Private]", IsModifier: false, IsHeld: isHeld}}
	}

	// Adaptive decay linger duration
	duration := e.lingerShort
	if action.Title != "" {
		duration = e.lingerLong
	}
	expiresAt := now.Add(duration)

	chordID := fmt.Sprintf("chord-%d-%d", now.UnixNano(), atomic.AddUint64(&e.chordSeq, 1))
	chord := DisplayChord{
		ID:         chordID,
		Keys:       keys,
		Streak:     1,
		Action:     action,
		Timestamp:  now,
		ExpiresAt:  expiresAt,
		IsActive:   isHeld,
		IsPrivate:  isPrivate,
		Velocity:   velocity,
		IsFinisher: isFinisher,
		Confidence: result.Confidence,
		AppName:    result.AppName,
	}

	// Append and enforce max limit
	e.chords = append(e.chords, chord)
	if len(e.chords) > e.maxChords {
		e.chords = e.chords[len(e.chords)-e.maxChords:]
	}

	e.lastChordKey = keyName
	e.lastChordMods = mods

	// Remember this chord as a possible first half of a sequence. Plain keys
	// and finishers reset the pending state, so a stray key can never combine
	// with a much later shortcut.
	if e.semanticContext && e.matcher != nil && mods != hook.ModNone && !isPrivate {
		e.seqMods, e.seqKey, e.seqAt = mods, keyName, now
	} else {
		e.seqAt = time.Time{}
	}

	// Acoustic feedback: Click for regular key/shortcut with spatial stage placement
	e.playAcoustic(SoundClick, keyName)

	e.emitStateLocked()
}

// pushScrollLocked folds a wheel gesture into a single cue. A direction change
// starts a new cue, and so does a pause longer than scrollMergeWindow: the user
// stopped scrolling, so the next flick is a separate motion. The count is the
// accumulated distance rather than the number of events, because a trackpad
// reports a dozen small deltas for the same motion a wheel reports in one.
func (e *Engine) pushScrollLocked(now time.Time, mods hook.Modifiers, keyName string, delta int) {
	notches := float32(delta) / 120
	if notches < 0 {
		notches = -notches
	}
	if notches < 0.05 {
		notches = 0.05
	}

	if len(e.chords) > 0 && e.lastChordMods == mods && strings.EqualFold(e.lastChordKey, keyName) {
		last := &e.chords[len(e.chords)-1]
		if now.Sub(last.Timestamp) <= scrollMergeWindow {
			last.ScrollAmount += notches
			last.Timestamp = now
			last.ExpiresAt = now.Add(e.lingerShort)
			e.playAcoustic(SoundClick, keyName)
			e.emitStateLocked()
			return
		}
	}

	// A new gesture: drop the repeat chain so the generic path appends a fresh cue
	// instead of refreshing the previous one, then seed it with the first notch.
	e.lastChordKey = ""
	e.lastChordMods = hook.ModNone
	e.pushChordLocked(now, mods, keyName, false, false)
	if n := len(e.chords); n > 0 {
		e.chords[n-1].ScrollAmount = notches
	}
}

// buildKeyCaps renders the modifier prefix and primary key for a chord.
func buildKeyCaps(mods hook.Modifiers, keyName string, held bool) []KeyCap {
	keys := make([]KeyCap, 0, 5)
	if mods.HasCtrl() {
		keys = append(keys, KeyCap{Label: "Ctrl", IsModifier: true, IsHeld: held})
	}
	if mods.HasShift() {
		keys = append(keys, KeyCap{Label: "Shift", IsModifier: true, IsHeld: held})
	}
	if mods.HasAlt() {
		keys = append(keys, KeyCap{Label: "Alt", IsModifier: true, IsHeld: held})
	}
	if mods.HasSuper() {
		keys = append(keys, KeyCap{Label: "Win", IsModifier: true, IsHeld: held})
	}
	return append(keys, KeyCap{Label: keyName, IsHeld: held})
}

func mouseButtonName(btn hook.MouseButton) string {
	switch btn {
	case hook.ButtonLeft:
		return "Left Click"
	case hook.ButtonRight:
		return "Right Click"
	case hook.ButtonMiddle:
		return "Middle Click"
	case hook.ButtonExtra1:
		return "Mouse 4"
	case hook.ButtonExtra2:
		return "Mouse 5"
	default:
		return ""
	}
}

func (e *Engine) releaseChordKey(keyName string) {
	normKey := normalizeKeyName(keyName)
	delete(e.heldKeys, normKey)

	changed := false
	for i := range e.chords {
		chord := &e.chords[i]
		allReleased := true
		for j := range chord.Keys {
			if strings.EqualFold(chord.Keys[j].Label, keyName) {
				chord.Keys[j].IsHeld = false
				changed = true
			}
			if chord.Keys[j].IsHeld {
				allReleased = false
			}
		}
		if allReleased && chord.IsActive {
			chord.IsActive = false
			changed = true
		}
	}
	if changed {
		e.emitStateLocked()
	}
}

// handleKeyUp processes key release events.
func (e *Engine) handleKeyUp(ev hook.InputEvent) {
	normKey := normalizeKeyName(ev.KeyName)
	e.heldKeys[normKey] = false

	isMod := ev.IsModifier || isModifierKey(ev.KeyName)
	changed := false

	for i := range e.chords {
		chord := &e.chords[i]
		if !chord.IsActive {
			continue
		}

		if chord.IsPrivate {
			chord.Keys[0].IsHeld = false
			chord.IsActive = false
			changed = true
			continue
		}

		for j := range chord.Keys {
			if keyMatches(chord.Keys[j], ev.KeyName, isMod) || (isTextWordChord(*chord) && wordKeyMatches(chord.Keys[j].Label, ev.KeyName)) {
				if chord.Keys[j].IsHeld {
					chord.Keys[j].IsHeld = false
					changed = true
				}
			}
		}

		// If all keys in a chord are released, mark chord inactive
		allReleased := true
		for _, k := range chord.Keys {
			if k.IsHeld {
				allReleased = false
				break
			}
		}
		if allReleased && chord.IsActive {
			chord.IsActive = false
			changed = true
		}
	}

	if changed {
		e.emitStateLocked()
	}
}

// handleMouseEvent processes mouse button, scroll, and movement drag events.
func (e *Engine) handleMouseEvent(ev hook.InputEvent) {
	now := ev.Timestamp
	if now.IsZero() {
		now = time.Now()
	}

	active := false
	wheelDelta := 0
	btn := hook.ButtonNone
	clickTimestamp := e.lastClickTime

	switch ev.Type {
	case hook.EventMouseDown:
		if e.mouseRetireTimer != nil {
			e.mouseRetireTimer.Stop()
			e.mouseRetireTimer = nil
		}
		if e.mouseBadgeTimer != nil {
			e.mouseBadgeTimer.Stop()
			e.mouseBadgeTimer = nil
		}

		active = true
		if !e.lastClickTime.IsZero() && now.Sub(e.lastClickTime) <= doubleClickWindow && e.lastClickButton == ev.Button {
			e.clickCount++
		} else {
			e.clickCount = 1
		}
		e.pressedButton = ev.Button
		e.lastClickButton = ev.Button
		e.lastClickTime = now
		clickTimestamp = now
		btn = ev.Button
		e.dragStartX = ev.X
		e.dragStartY = ev.Y
		e.isDragging = false

		btnName := mouseButtonName(ev.Button)
		if btnName != "" {
			e.heldKeys[normalizeKeyName(btnName)] = true
			mods := ev.Modifiers
			if mods.IsEmpty() {
				mods = e.activeModifiersFromHeld()
			}
			allowed := false
			if e.filterMode == FilterSilent {
				allowed = false
			} else if e.filterMode == FilterCustom {
				allowed = e.isCustomAllowedLocked(btnName, 0, mods)
			} else if e.showMouseClicks || mods != hook.ModNone {
				allowed = true
			}
			if allowed {
				if e.clickCount > 1 && len(e.chords) > 0 && e.lastChordKey == btnName && e.lastChordMods == mods {
					lastChord := &e.chords[len(e.chords)-1]
					lastChord.Streak = e.clickCount
					lastChord.ExpiresAt = now.Add(e.lingerShort)
					lastChord.Timestamp = now
					lastChord.IsActive = true
				} else {
					e.pushChordLocked(now, mods, btnName, true, false)
				}
			}
		}

	case hook.EventMouseUp:
		active = false
		e.pressedButton = hook.ButtonNone
		e.lastClickTime = now
		clickTimestamp = now
		btn = e.lastClickButton

		btnName := mouseButtonName(ev.Button)
		if e.isDragging {
			e.isDragging = false
			e.releaseChordKey("Drag")
		} else if btnName != "" {
			e.releaseChordKey(btnName)
		}

	case hook.EventMouseMove:
		active = e.pressedButton != hook.ButtonNone
		btn = e.pressedButton
		clickTimestamp = e.lastClickTime
		if active && !e.isDragging && e.dragThreshold > 0 {
			dx := ev.X - e.dragStartX
			dy := ev.Y - e.dragStartY
			if dx*dx+dy*dy > e.dragThreshold*e.dragThreshold {
				e.isDragging = true
				mods := ev.Modifiers
				if mods.IsEmpty() {
					mods = e.activeModifiersFromHeld()
				}
				allowed := false
				if e.filterMode == FilterSilent {
					allowed = false
				} else if e.filterMode == FilterCustom {
					allowed = e.isCustomAllowedLocked("Drag", 0, mods)
				} else if e.showMouseClicks || mods != hook.ModNone {
					allowed = true
				}
				if allowed {
					e.pushChordLocked(now, mods, "Drag", true, false)
				}
			}
		}

	case hook.EventMouseWheel:
		active = false
		e.clickCount = 0 // Scroll wheel is not a mouse click; clear clickCount to prevent multi-click badges
		clickTimestamp = now
		wheelDelta = ev.WheelDelta
		wheelKey := "Scroll Down"
		if ev.WheelDelta > 0 {
			wheelKey = "Scroll Up"
		}
		mods := ev.Modifiers
		if mods.IsEmpty() {
			mods = e.activeModifiersFromHeld()
		}
		allowed := false
		if e.filterMode == FilterSilent {
			allowed = false
		} else if e.filterMode == FilterCustom {
			allowed = e.isCustomAllowedLocked(wheelKey, 0, mods)
		} else if e.showMouseClicks || mods != hook.ModNone {
			allowed = true
		}
		if allowed {
			e.pushScrollLocked(now, mods, wheelKey, ev.WheelDelta)
		}
	}

	e.mouseState = MouseVisualState{
		Button:     btn,
		X:          ev.X,
		Y:          ev.Y,
		Active:     active,
		Dragging:   e.isDragging,
		Wheel:      wheelDelta,
		Modifiers:  ev.Modifiers,
		Timestamp:  clickTimestamp,
		ClickCount: e.clickCount,
	}

	if ev.Type == hook.EventMouseUp || ev.Type == hook.EventMouseWheel {
		e.scheduleMouseRetirementLocked()
	}

	e.emitStateLocked()
}

// Subscribe registers a subscriber channel for UIState snapshots.
// Returns a buffered channel that immediately receives the initial UIState.
func (e *Engine) Subscribe() <-chan UIState {
	e.subMu.Lock()
	defer e.subMu.Unlock()

	ch := make(chan UIState, 32)
	e.subscribers = append(e.subscribers, ch)

	e.mu.RLock()
	st := e.snapshot()
	e.mu.RUnlock()

	ch <- st
	return ch
}

// Unsubscribe unregisters and closes a subscriber channel.
func (e *Engine) Unsubscribe(target <-chan UIState) {
	e.subMu.Lock()
	defer e.subMu.Unlock()

	for i, ch := range e.subscribers {
		if ch == target {
			e.subscribers = append(e.subscribers[:i], e.subscribers[i+1:]...)
			close(ch)
			return
		}
	}
}

// State returns the current immutable UIState snapshot.
func (e *Engine) State() UIState {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.snapshot()
}

// Chords returns an immutable copy of currently visible display chords.
func (e *Engine) Chords() []DisplayChord {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.snapshot().Chords
}

// IsKeyHeld returns whether the given key is currently pressed down.
func (e *Engine) IsKeyHeld(keyName string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.heldKeys[normalizeKeyName(keyName)]
}

// HeldKeys returns a copy of currently pressed keys.
func (e *Engine) HeldKeys() map[string]bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make(map[string]bool, len(e.heldKeys))
	for k, v := range e.heldKeys {
		if v {
			out[k] = true
		}
	}
	return out
}

// SetFilterMode sets the filter mode and emits updated UIState.
func (e *Engine) SetFilterMode(m FilterMode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.filterMode = m
	e.emitStateLocked()
}

// FilterMode returns the active filter mode.
func (e *Engine) FilterMode() FilterMode {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.filterMode
}

// SetAllowedKeys updates the allowed keys whitelist for FilterCustom mode.
func (e *Engine) SetAllowedKeys(keys []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.allowedKeys = make(map[string]bool, len(keys))
	for _, k := range keys {
		norm := normalizeAllowedKey(k)
		if norm != "" {
			e.allowedKeys[norm] = true
		}
	}
	e.emitStateLocked()
}

// AllowedKeys returns a list of active allowed keys in FilterCustom mode.
func (e *Engine) AllowedKeys() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]string, 0, len(e.allowedKeys))
	for k := range e.allowedKeys {
		out = append(out, k)
	}
	return out
}

// normalizeAllowedKey maps varying key label aliases to a consistent key token.
func normalizeAllowedKey(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	switch s {
	case "controlleft", "controlright", "ctrl", "control", "lctrl", "rctrl":
		return "ctrl"
	case "shiftleft", "shiftright", "shift", "lshift", "rshift":
		return "shift"
	case "alt", "altleft", "altright", "menu", "lalt", "ralt":
		return "alt"
	case "metaleft", "metaright", "win", "windows", "super", "cmd", "command":
		return "win"
	case "left click", "left", "mouseleft":
		return "left click"
	case "right click", "right", "mouseright":
		return "right click"
	case "middle click", "middle", "mousemiddle":
		return "middle click"
	case "scroll up", "scrollup", "wheelup":
		return "scroll up"
	case "scroll down", "scrolldown", "wheeldown":
		return "scroll down"
	case "drag", "mousedrag":
		return "drag"
	case "return", "enter", "↵":
		return "enter"
	case "escape", "esc":
		return "esc"
	case "space", "␣":
		return "space"
	case "backspace", "⌫", "bksp":
		return "backspace"
	case "tab", "⇥":
		return "tab"
	default:
		return s
	}
}

// isCustomAllowedLocked checks whether the given key/modifier/action is allowed by allowedKeys in FilterCustom mode.
func (e *Engine) isCustomAllowedLocked(keyName string, vkCode uint32, mods hook.Modifiers) bool {
	if len(e.allowedKeys) == 0 {
		return false
	}
	// 1. Check if any held modifier is allowed
	if mods.HasCtrl() && e.allowedKeys["ctrl"] {
		return true
	}
	if mods.HasShift() && e.allowedKeys["shift"] {
		return true
	}
	if mods.HasAlt() && e.allowedKeys["alt"] {
		return true
	}
	if mods.HasSuper() && e.allowedKeys["win"] {
		return true
	}
	// 2. Check primary key name
	norm := normalizeAllowedKey(keyName)
	if e.allowedKeys[norm] {
		return true
	}
	if e.allowedKeys[strings.ToLower(strings.TrimSpace(keyName))] {
		return true
	}
	// 3. Check modifier by keyName if pressing modifier itself
	modKey := normalizeModifier(keyName)
	if modKey == "ctrl" || modKey == "shift" || modKey == "alt" || modKey == "win" {
		if e.allowedKeys[modKey] {
			return true
		}
	}
	return false
}

// SetSoundActive toggles acoustic feedback and emits updated UIState.
func (e *Engine) SetSoundActive(active bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.soundActive = active
	e.emitStateLocked()
}

// IsSoundActive returns whether acoustic feedback is active.
func (e *Engine) IsSoundActive() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.soundActive
}

// SetSoundPlayer attaches or detaches a SoundPlayer.
func (e *Engine) SetSoundPlayer(p SoundPlayer) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.soundPlayer = p
}

// SetMatcher attaches or updates the semantic matcher.
func (e *Engine) SetMatcher(m semantic.Matcher) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.matcher = m
}

// SetAppContext updates the cached foreground application snapshot. It is
// intentionally cheap: the context provider only publishes when the foreground
// window actually changes, so this never runs on every keystroke.
func (e *Engine) SetAppContext(ctx context.AppContext) {
	e.mu.Lock()
	changed := e.appContext.Key() != ctx.Key() || e.appContext.Seq != ctx.Seq
	// Application identity is the executable (plus the resolved profile), not the
	// window title: a browser retitling itself on every tab must not interrupt a
	// run of keystrokes.
	identityChanged := !strings.EqualFold(e.appContext.Exe, ctx.Exe) || e.appContext.AppID != ctx.AppID
	e.appContext = ctx
	// A foreground switch invalidates any pending sequence first chord.
	if changed {
		e.seqAt = time.Time{}
	}
	if identityChanged {
		// Different application, different vocabulary: → means "forward" in a
		// browser and "move the caret" in an editor. Two presses on either side of
		// a switch are two cues, never one ×2.
		e.lastChordKey = ""
		e.lastChordMods = hook.ModNone
	}
	e.mu.Unlock()
}

// AppContext returns the current foreground application snapshot.
func (e *Engine) AppContext() context.AppContext {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.appContext
}

// SetSemanticContext enables or disables context-aware action subtitles. Privacy
// masking remains active regardless of this switch.
func (e *Engine) SetSemanticContext(enabled bool) {
	e.mu.Lock()
	if e.semanticContext != enabled {
		e.seqAt = time.Time{}
	}
	e.semanticContext = enabled
	e.mu.Unlock()
}

// SemanticContext reports whether contextual interpretation is enabled.
func (e *Engine) SemanticContext() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.semanticContext
}

// SetThemeName updates the UI theme name and emits updated UIState.
func (e *Engine) SetThemeName(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.themeName = name
	e.emitStateLocked()
}

// ThemeName returns the active UI theme name.
func (e *Engine) ThemeName() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.themeName
}

// snapshot produces a deep-copied, immutable UIState.
func (e *Engine) snapshot() UIState {
	chordsCopy := make([]DisplayChord, len(e.chords))
	for i, c := range e.chords {
		keysCopy := make([]KeyCap, len(c.Keys))
		copy(keysCopy, c.Keys)
		c.Keys = keysCopy
		chordsCopy[i] = c
	}

	return UIState{
		Chords:       chordsCopy,
		MouseState:   e.mouseState,
		FilterMode:   e.filterMode,
		SoundActive:  e.soundActive,
		ThemeName:    e.themeName,
		WordGrouping: e.wordGrouping,
		CadenceWPM:   e.currentWPM,
	}
}

// CadenceWPM returns current typing velocity in words per minute.
func (e *Engine) CadenceWPM() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentWPM
}

// recordStrokeLocked tracks typing timestamps in a sliding window and calculates instantaneous WPM.
func (e *Engine) recordStrokeLocked(now time.Time) int {
	cutoff := now.Add(-2500 * time.Millisecond)
	kept := e.recentStrokes[:0]
	for _, t := range e.recentStrokes {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	kept = append(kept, now)
	e.recentStrokes = kept

	n := len(e.recentStrokes)
	if n < 2 {
		e.currentWPM = 0
		return 0
	}
	dt := now.Sub(e.recentStrokes[0]).Seconds()
	if dt < 0.25 {
		dt = 0.25
	}
	wpm := int(math.Round((float64(n) * 12.0) / dt))
	e.currentWPM = wpm
	return wpm
}

func (e *Engine) decayCadenceLocked(now time.Time) {
	if len(e.recentStrokes) > 0 {
		cutoff := now.Add(-2500 * time.Millisecond)
		if e.recentStrokes[len(e.recentStrokes)-1].Before(cutoff) {
			e.recentStrokes = nil
			e.currentWPM = 0
		}
	}
}

func velocityForWPM(wpm int) int {
	if wpm >= 75 {
		return 2 // Allegro / Supersonic
	}
	if wpm >= 35 {
		return 1 // Andante / Flow
	}
	return 0 // Adagio / Contemplation
}

func isFinisherKey(key string) bool {
	k := strings.ToLower(strings.TrimSpace(key))
	switch k {
	case "enter", "return", "↵", "tab", "⇥", "esc", "escape":
		return true
	default:
		return false
	}
}

// playAcoustic dispatches key acoustic feedback with spatial stage awareness.
func (e *Engine) playAcoustic(s SoundType, key string) {
	if !e.soundActive || e.soundPlayer == nil {
		return
	}
	if sp, ok := e.soundPlayer.(SpatialSoundPlayer); ok {
		sp.PlaySpatialSound(s, key)
	} else {
		e.soundPlayer.PlaySound(s)
	}
}

// emitStateLocked notifies all subscribers with a fresh immutable UIState snapshot.
func (e *Engine) emitStateLocked() {
	st := e.snapshot()

	e.subMu.RLock()
	defer e.subMu.RUnlock()

	for _, ch := range e.subscribers {
		select {
		case ch <- st:
		default:
			// Non-blocking drain old item and push new item
			select {
			case <-ch:
			default:
			}
			select {
			case ch <- st:
			default:
			}
		}
	}
}

// activeModifiersFromHeld checks heldKeys to synthesize active modifiers.
func (e *Engine) activeModifiersFromHeld() hook.Modifiers {
	var mods hook.Modifiers
	if e.heldKeys["ctrl"] {
		mods |= hook.ModCtrl
	}
	if e.heldKeys["shift"] {
		mods |= hook.ModShift
	}
	if e.heldKeys["alt"] {
		mods |= hook.ModAlt
	}
	if e.heldKeys["win"] {
		mods |= hook.ModSuper
	}
	return mods
}

// extractModifiers reconstructs hook.Modifiers from a DisplayChord's modifier keycaps.
func extractModifiers(chord DisplayChord) hook.Modifiers {
	var mods hook.Modifiers
	for _, k := range chord.Keys {
		if k.IsModifier {
			switch normalizeModifier(k.Label) {
			case "ctrl":
				mods |= hook.ModCtrl
			case "shift":
				mods |= hook.ModShift
			case "alt":
				mods |= hook.ModAlt
			case "win":
				mods |= hook.ModSuper
			}
		}
	}
	return mods
}

// normalizeKeyName standardizes key names for heldKeys indexing.
func normalizeKeyName(name string) string {
	mod := normalizeModifier(name)
	if mod == "ctrl" || mod == "shift" || mod == "alt" || mod == "win" {
		return mod
	}
	return strings.ToLower(strings.TrimSpace(name))
}

// normalizeModifier maps various modifier key name aliases to canonical tokens.
func normalizeModifier(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	switch s {
	case "ctrl", "control", "lctrl", "rctrl", "lcontrol", "rcontrol", "controlkey", "leftcontrol", "rightcontrol":
		return "ctrl"
	case "shift", "lshift", "rshift", "shiftkey", "leftshift", "rightshift":
		return "shift"
	case "alt", "lalt", "ralt", "menu", "lmenu", "rmenu", "altkey", "leftalt", "rightalt":
		return "alt"
	case "win", "windows", "lwin", "rwin", "super", "cmd", "command", "meta", "leftwin", "rightwin":
		return "win"
	default:
		return s
	}
}

// isModifierKey checks if a given key name is a modifier key.
func isModifierKey(name string) bool {
	m := normalizeModifier(name)
	return m == "ctrl" || m == "shift" || m == "alt" || m == "win"
}

// keyMatches compares a KeyCap with an event key name taking modifiers and aliases into account.
func keyMatches(k KeyCap, evKeyName string, isMod bool) bool {
	if k.IsModifier || isMod {
		return normalizeModifier(k.Label) == normalizeModifier(evKeyName)
	}
	return strings.EqualFold(k.Label, evKeyName) || isKeyAliasMatch(k.Label, evKeyName)
}

// isKeyAliasMatch checks common keyboard symbol aliases.
func isKeyAliasMatch(label, evKeyName string) bool {
	l1 := strings.ToLower(strings.TrimSpace(label))
	l2 := strings.ToLower(strings.TrimSpace(evKeyName))
	if l1 == l2 {
		return true
	}
	if (l1 == "enter" || l1 == "return" || l1 == "↵" || l1 == "⏎") && (l2 == "enter" || l2 == "return" || l2 == "↵" || l2 == "⏎") {
		return true
	}
	if (l1 == "backspace" || l1 == "⌫" || l1 == "back" || l1 == "bksp") && (l2 == "backspace" || l2 == "⌫" || l2 == "back" || l2 == "bksp") {
		return true
	}
	if (l1 == "tab" || l1 == "⇥") && (l2 == "tab" || l2 == "⇥") {
		return true
	}
	if (l1 == "esc" || l1 == "escape") && (l2 == "esc" || l2 == "escape") {
		return true
	}
	if (l1 == "up" || l1 == "arrowup" || l1 == "↑") && (l2 == "up" || l2 == "arrowup" || l2 == "↑") {
		return true
	}
	if (l1 == "down" || l1 == "arrowdown" || l1 == "↓") && (l2 == "down" || l2 == "arrowdown" || l2 == "↓") {
		return true
	}
	if (l1 == "left" || l1 == "arrowleft" || l1 == "←") && (l2 == "left" || l2 == "arrowleft" || l2 == "←") {
		return true
	}
	if (l1 == "right" || l1 == "arrowright" || l1 == "→") && (l2 == "right" || l2 == "arrowright" || l2 == "→") {
		return true
	}
	return false
}

// isBackspaceKey checks if key is a Backspace variant.
func isBackspaceKey(name string, vkCode uint32) bool {
	if vkCode == 0x08 {
		return true
	}
	s := strings.ToLower(strings.TrimSpace(name))
	return s == "backspace" || s == "⌫" || s == "back" || s == "bksp"
}

// isSpecialKey checks whether the key is a special navigation/control/function key.
func isSpecialKey(name string, vkCode uint32) bool {
	switch vkCode {
	case 0x08: // VK_BACK
		return true
	case 0x09: // VK_TAB
		return true
	case 0x0D: // VK_RETURN
		return true
	case 0x1B: // VK_ESCAPE
		return true
	case 0x14: // VK_CAPITAL
		return true
	case 0x21, 0x22, 0x23, 0x24: // VK_PRIOR, VK_NEXT, VK_END, VK_HOME
		return true
	case 0x25, 0x26, 0x27, 0x28: // VK_LEFT, VK_UP, VK_RIGHT, VK_DOWN
		return true
	case 0x2C, 0x2D, 0x2E: // VK_SNAPSHOT, VK_INSERT, VK_DELETE
		return true
	case 0x90, 0x91: // VK_NUMLOCK, VK_SCROLL
		return true
	}
	if vkCode >= 0x70 && vkCode <= 0x87 { // VK_F1 to VK_F24
		return true
	}

	s := strings.ToLower(strings.TrimSpace(name))
	switch s {
	case "enter", "return", "↵", "⏎",
		"esc", "escape",
		"tab", "⇥",
		"backspace", "⌫", "back", "del", "delete", "⌦",
		"up", "down", "left", "right", "arrowup", "arrowdown", "arrowleft", "arrowright", "↑", "↓", "←", "→",
		"home", "end", "pageup", "pagedown", "pgup", "pgdn", "insert", "ins",
		"capslock", "numlock", "scrolllock", "printscreen", "prtscn", "pause", "break":
		return true
	}

	if len(s) >= 2 && s[0] == 'f' {
		if n, err := strconv.Atoi(s[1:]); err == nil && n >= 1 && n <= 24 {
			return true
		}
	}

	return false
}

// isAbsorbableCharacter determines whether chord is a single plain typed character.
func isAbsorbableCharacter(chord DisplayChord) bool {
	if chord.IsPrivate || chord.Action.Title != "" {
		return false
	}
	if len(chord.Keys) != 1 {
		return false
	}
	k := chord.Keys[0]
	if k.IsModifier {
		return false
	}
	if isBackspaceKey(k.Label, 0) || isSpecialKey(k.Label, 0) {
		return false
	}
	runes := []rune(k.Label)
	return len(runes) == 1
}

// maxWordRunes caps how long a single word chip can grow. Typing — or holding a
// key and letting the OS repeat it — otherwise extends one label without bound,
// and the HUD has no truncation, so a long run widens the card until it walks
// off screen. 24 runes covers the longest everyday words with room to spare;
// past that the word is finished and an ellipsis says so.
const maxWordRunes = 24

// appendWordRune adds one character to a word chip, sealing it at the cap.
func appendWordRune(label, char string) string {
	if strings.HasSuffix(label, "…") {
		return label
	}
	if utf8.RuneCountInString(label) >= maxWordRunes {
		return label + "…"
	}
	return label + char
}

// isWordGroupingChar checks if key is a single alphanumeric character (A-Z, a-z, 0-9).
// Non-alphanumeric keys (including Space, Enter, symbols, punctuation) return false,
// causing words to break and wrap into distinct chords as per user specification.
func isWordGroupingChar(keyName string, vkCode uint32, mods hook.Modifiers) (string, bool) {
	// Single rune check
	runes := []rune(keyName)
	if len(runes) == 1 {
		r := runes[0]
		// Lowercase letter 'a'-'z'
		if r >= 'a' && r <= 'z' {
			if mods == hook.ModNone {
				return string(r), true
			}
			if mods == hook.ModShift {
				return string(r - ('a' - 'A')), true
			}
			return "", false
		}
		// Uppercase letter 'A'-'Z'
		if r >= 'A' && r <= 'Z' {
			if mods == hook.ModShift {
				return string(r), true
			}
			if mods == hook.ModNone {
				return strings.ToLower(string(r)), true
			}
			return "", false
		}
		// Digit '0'-'9'
		if r >= '0' && r <= '9' {
			if mods == hook.ModNone {
				return string(r), true
			}
			return "", false
		}
	}

	// VK code fallback for A-Z (0x41 to 0x5A) and 0-9 (0x30 to 0x39)
	if vkCode >= 0x41 && vkCode <= 0x5A {
		ch := rune(vkCode) // 'A'-'Z'
		if mods == hook.ModShift {
			return string(ch), true
		}
		if mods == hook.ModNone {
			return string(ch + ('a' - 'A')), true
		}
		return "", false
	}
	if vkCode >= 0x30 && vkCode <= 0x39 {
		if mods == hook.ModNone {
			return string(rune(vkCode)), true
		}
		return "", false
	}

	return "", false
}

// isTextWordChord determines whether chord is a single plain text word chord of alphanumeric characters.
func isTextWordChord(chord DisplayChord) bool {
	if chord.IsPrivate || chord.Action.Title != "" {
		return false
	}
	if len(chord.Keys) != 1 {
		return false
	}
	k := chord.Keys[0]
	if k.IsModifier {
		return false
	}
	if isBackspaceKey(k.Label, 0) || isSpecialKey(k.Label, 0) {
		return false
	}
	// Must consist only of letters and digits. A word that hit the length cap
	// carries a trailing ellipsis and is still the same word, so it keeps
	// absorbing keys instead of starting a fresh chip.
	for _, r := range strings.TrimSuffix(k.Label, "…") {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

// wordKeyMatches checks if evKeyName matches the label or last character of a word chord.
func wordKeyMatches(label, evKeyName string) bool {
	if strings.EqualFold(label, evKeyName) {
		return true
	}
	l1 := strings.ToLower(label)
	l2 := strings.ToLower(evKeyName)
	if l2 == "space" || l2 == "␣" || l2 == " " {
		return strings.HasSuffix(l1, " ")
	}
	return strings.HasSuffix(l1, l2)
}
