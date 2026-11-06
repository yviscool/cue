package engine

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/semantic"
)

const (
	defaultDecayInterval   = 50 * time.Millisecond
	defaultStreakThreshold = 1200 * time.Millisecond
	defaultLingerShort     = 1800 * time.Millisecond
	defaultLingerLong      = 2500 * time.Millisecond
	defaultMaxChords       = 5
	defaultThemeName       = "Default"
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
	appContext       string
	themeName        string
	showMouseClicks  bool
	dragThreshold    int
	wordGrouping     bool
	wordBreakPending bool

	// Current state
	heldKeys        map[string]bool // normalized key names currently pressed
	chords          []DisplayChord  // active display chords on screen
	mouseState      MouseVisualState
	pressedButton   hook.MouseButton
	lastClickTime   time.Time
	lastClickButton hook.MouseButton
	dragStartX      int
	dragStartY      int
	isDragging      bool
	lastChordKey    string
	lastChordMods   hook.Modifiers
	chordSeq        uint64

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
func WithAppContext(ctx string) Option {
	return func(e *Engine) {
		e.appContext = ctx
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

// NewEngine creates and initializes a Engine instance.
func NewEngine(opts ...Option) *Engine {
	e := &Engine{
		filterMode:      FilterSmart,
		decayInterval:   defaultDecayInterval,
		streakThreshold: defaultStreakThreshold,
		lingerShort:     defaultLingerShort,
		lingerLong:      defaultLingerLong,
		maxChords:       defaultMaxChords,
		soundActive:     true,
		themeName:       defaultThemeName,
		showMouseClicks: false,
		dragThreshold:   50,
		heldKeys:        make(map[string]bool),
		allowedKeys:     make(map[string]bool),
		stopCh:          make(chan struct{}),
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

	// Mouse feedback is time-based: a click ring lives ~220 ms and the button/wheel
	// indicator ~260 ms. Without this cleanup the last frame would remain painted
	// until the next input event, leaving a permanently stuck ring or scroll arrow.
	mouseRetired := false
	if !e.mouseState.Timestamp.IsZero() && !e.mouseState.Active &&
		now.Sub(e.mouseState.Timestamp) > 320*time.Millisecond {
		e.mouseState.Active = false
		e.mouseState.Dragging = false
		e.mouseState.Wheel = 0
		e.mouseState.Button = hook.ButtonNone
		e.mouseState.Timestamp = time.Time{}
		mouseRetired = true
	}

	if removed > 0 || mouseRetired {
		e.emitStateLocked()
	}

	return removed
}

// Clear immediately clears all on-screen keys and mouse visual states.
func (e *Engine) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.chords = nil
	e.mouseState = MouseVisualState{}
	e.lastChordKey = ""
	e.lastChordMods = hook.ModNone
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
		if e.soundActive && e.soundPlayer != nil {
			e.soundPlayer.PlaySound(SoundThock)
		}
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
			if e.soundActive && e.soundPlayer != nil {
				e.soundPlayer.PlaySound(SoundClick)
			}
			e.emitStateLocked()
			return
		}
	}

	// Word Grouping: Alphanumeric grouping
	if e.wordGrouping {
		if charStr, ok := isWordGroupingChar(ev.KeyName, ev.VKCode, mods); ok {
			const wordThreshold = 600 * time.Millisecond
			if len(e.chords) > 0 && isTextWordChord(e.chords[len(e.chords)-1]) && !e.wordBreakPending {
				lastChord := &e.chords[len(e.chords)-1]
				if now.Sub(lastChord.Timestamp) <= wordThreshold {
					lastChord.Keys[0].Label += charStr
					lastChord.ExpiresAt = now.Add(e.lingerShort)
					lastChord.Timestamp = now
					lastChord.IsActive = true
					lastChord.Keys[0].IsHeld = true
					e.lastChordKey = lastChord.Keys[0].Label
					e.lastChordMods = hook.ModNone
					if e.soundActive && e.soundPlayer != nil {
						e.soundPlayer.PlaySound(SoundClick)
					}
					e.emitStateLocked()
					return
				}
			}
			e.wordBreakPending = false
			// Initial character or new word after break: push as new word chord without modifier badges
			e.pushChordLocked(now, hook.ModNone, charStr, true)
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
			if e.soundActive && e.soundPlayer != nil {
				e.soundPlayer.PlaySound(SoundClick)
			}
			e.emitStateLocked()
			return
		}
	}

	e.pushChordLocked(now, mods, ev.KeyName, true)
}

// pushChordLocked constructs a new DisplayChord or increments streak for repetitive inputs.
func (e *Engine) pushChordLocked(now time.Time, mods hook.Modifiers, keyName string, isHeld bool) {
	// Streak check: identical to the most recent chord (same modifiers + same key) within threshold
	if len(e.chords) > 0 && (!e.wordGrouping || mods != hook.ModNone) && e.lastChordMods == mods && strings.EqualFold(e.lastChordKey, keyName) {
		lastChord := &e.chords[len(e.chords)-1]
		if now.Sub(lastChord.Timestamp) <= e.streakThreshold {
			lastChord.Streak++
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
			if e.soundActive && e.soundPlayer != nil {
				e.soundPlayer.PlaySound(SoundClick)
			}
			e.emitStateLocked()
			return
		}
	}

	// Construct DisplayChord
	var keys []KeyCap
	if mods.HasCtrl() {
		keys = append(keys, KeyCap{Label: "Ctrl", IsModifier: true, IsHeld: true})
	}
	if mods.HasShift() {
		keys = append(keys, KeyCap{Label: "Shift", IsModifier: true, IsHeld: true})
	}
	if mods.HasAlt() {
		keys = append(keys, KeyCap{Label: "Alt", IsModifier: true, IsHeld: true})
	}
	if mods.HasSuper() {
		keys = append(keys, KeyCap{Label: "Win", IsModifier: true, IsHeld: true})
	}

	primaryKeyCap := KeyCap{
		Label:      keyName,
		IsModifier: false,
		IsHeld:     isHeld,
	}
	keys = append(keys, primaryKeyCap)

	// Call semantic matcher MatchAction
	var action semantic.ActionInfo
	if e.matcher != nil {
		if act, ok := e.matcher.MatchAction(e.appContext, mods, keyName); ok {
			action = act
		}
	}

	// Privacy mask check
	isPrivate := false
	if e.matcher != nil && e.matcher.IsSensitive(e.appContext, keyName) {
		isPrivate = true
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
		ID:        chordID,
		Keys:      keys,
		Streak:    1,
		Action:    action,
		Timestamp: now,
		ExpiresAt: expiresAt,
		IsActive:  isHeld,
		IsPrivate: isPrivate,
	}

	// Append and enforce max limit
	e.chords = append(e.chords, chord)
	if len(e.chords) > e.maxChords {
		e.chords = e.chords[len(e.chords)-e.maxChords:]
	}

	e.lastChordKey = keyName
	e.lastChordMods = mods

	// Acoustic feedback: Click for regular key/shortcut
	if e.soundActive && e.soundPlayer != nil {
		e.soundPlayer.PlaySound(SoundClick)
	}

	e.emitStateLocked()
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
		active = true
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
				e.pushChordLocked(now, mods, btnName, true)
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
					e.pushChordLocked(now, mods, "Drag", true)
				}
			}
		}

	case hook.EventMouseWheel:
		active = false
		e.lastClickTime = now
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
			e.pushChordLocked(now, mods, wheelKey, false)
		}
	}

	e.mouseState = MouseVisualState{
		Button:    btn,
		X:         ev.X,
		Y:         ev.Y,
		Active:    active,
		Dragging:  e.isDragging,
		Wheel:     wheelDelta,
		Modifiers: ev.Modifiers,
		Timestamp: clickTimestamp,
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

// SetAppContext updates the active application context for semantic matching.
func (e *Engine) SetAppContext(ctx string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.appContext = ctx
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
	// Must consist only of letters and digits
	for _, r := range k.Label {
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
