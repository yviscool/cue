package engine

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/semantic"
)

// mockMatcher implements semantic.Matcher for testing.
type mockMatcher struct {
	actions map[string]semantic.ActionInfo
	secrets map[string]bool
}

func newMockMatcher() *mockMatcher {
	return &mockMatcher{
		actions: make(map[string]semantic.ActionInfo),
		secrets: make(map[string]bool),
	}
}

func (m *mockMatcher) MatchAction(appContext string, mods hook.Modifiers, keyName string) (semantic.ActionInfo, bool) {
	key := mods.String() + "+" + keyName
	if mods.IsEmpty() {
		key = keyName
	}
	act, ok := m.actions[key]
	return act, ok
}

func (m *mockMatcher) IsSensitive(appContext string, keyName string) bool {
	return m.secrets[keyName]
}

// mockSoundPlayer records acoustic feedback calls.
type mockSoundPlayer struct {
	mu     sync.Mutex
	sounds []SoundType
}

func (p *mockSoundPlayer) PlaySound(s SoundType) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sounds = append(p.sounds, s)
}

func (p *mockSoundPlayer) count(s SoundType) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	c := 0
	for _, snd := range p.sounds {
		if snd == s {
			c++
		}
	}
	return c
}

func TestModifierKeyChordAggregation(t *testing.T) {
	eng := NewEngine(WithFilterMode(FilterSmart))
	defer eng.Stop()

	// 1. User holds Ctrl
	eng.HandleEvent(hook.InputEvent{
		Type:       hook.EventKeyDown,
		KeyName:    "Ctrl",
		IsModifier: true,
		Modifiers:  hook.ModCtrl,
	})

	// Modifier alone does not create a display chord
	if len(eng.Chords()) != 0 {
		t.Fatalf("expected 0 chords for modifier alone, got %d", len(eng.Chords()))
	}
	if !eng.IsKeyHeld("Ctrl") {
		t.Fatalf("expected Ctrl to be held")
	}

	// 2. User holds Shift as well
	eng.HandleEvent(hook.InputEvent{
		Type:       hook.EventKeyDown,
		KeyName:    "Shift",
		IsModifier: true,
		Modifiers:  hook.ModCtrl | hook.ModShift,
	})
	if !eng.IsKeyHeld("Shift") {
		t.Fatalf("expected Shift to be held")
	}

	// 3. User taps P
	eng.HandleEvent(hook.InputEvent{
		Type:       hook.EventKeyDown,
		KeyName:    "P",
		IsModifier: false,
		Modifiers:  hook.ModCtrl | hook.ModShift,
	})

	chords := eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 chord, got %d", len(chords))
	}

	chord := chords[0]
	if len(chord.Keys) != 3 {
		t.Fatalf("expected 3 keys [Ctrl, Shift, P], got %d", len(chord.Keys))
	}

	// Verify keycaps: [Ctrl, Shift, P]
	if chord.Keys[0].Label != "Ctrl" || !chord.Keys[0].IsModifier || !chord.Keys[0].IsHeld {
		t.Errorf("expected Key 0 to be Ctrl (held modifier), got %+v", chord.Keys[0])
	}
	if chord.Keys[1].Label != "Shift" || !chord.Keys[1].IsModifier || !chord.Keys[1].IsHeld {
		t.Errorf("expected Key 1 to be Shift (held modifier), got %+v", chord.Keys[1])
	}
	if chord.Keys[2].Label != "P" || chord.Keys[2].IsModifier || !chord.Keys[2].IsHeld {
		t.Errorf("expected Key 2 to be P (held primary key), got %+v", chord.Keys[2])
	}

	if !chord.IsActive {
		t.Errorf("expected chord to be active while keys are held")
	}
	if chord.Streak != 1 {
		t.Errorf("expected streak 1, got %d", chord.Streak)
	}
}

func TestKeyHoldStateTracking(t *testing.T) {
	eng := NewEngine(WithFilterMode(FilterSmart))
	defer eng.Stop()

	// Press Ctrl + Shift + P
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Ctrl", IsModifier: true, Modifiers: hook.ModCtrl})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Shift", IsModifier: true, Modifiers: hook.ModCtrl | hook.ModShift})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "P", Modifiers: hook.ModCtrl | hook.ModShift})

	chords := eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 chord, got %d", len(chords))
	}
	if !chords[0].IsActive {
		t.Fatalf("chord should be active")
	}

	// Release primary key P
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyUp, KeyName: "P"})
	if eng.IsKeyHeld("P") {
		t.Errorf("expected P to not be held")
	}

	chords = eng.Chords()
	if chords[0].Keys[2].IsHeld {
		t.Errorf("expected P keycap to have IsHeld=false")
	}
	if !chords[0].IsActive {
		t.Errorf("chord should still be active because modifiers are still held")
	}

	// Release Shift
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyUp, KeyName: "Shift", IsModifier: true})
	if eng.IsKeyHeld("Shift") {
		t.Errorf("expected Shift to not be held")
	}
	chords = eng.Chords()
	if chords[0].Keys[1].IsHeld {
		t.Errorf("expected Shift keycap to have IsHeld=false")
	}
	if !chords[0].IsActive {
		t.Errorf("chord should still be active because Ctrl is still held")
	}

	// Release Ctrl
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyUp, KeyName: "Ctrl", IsModifier: true})
	if eng.IsKeyHeld("Ctrl") {
		t.Errorf("expected Ctrl to not be held")
	}
	chords = eng.Chords()
	if chords[0].Keys[0].IsHeld {
		t.Errorf("expected Ctrl keycap to have IsHeld=false")
	}
	if chords[0].IsActive {
		t.Errorf("expected chord to be inactive once all keys are released")
	}
}

func TestRepeatStreakDetection(t *testing.T) {
	eng := NewEngine(
		WithFilterMode(FilterSmart),
		WithStreakThreshold(1200*time.Millisecond),
	)
	defer eng.Stop()

	baseTime := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)

	// Press Ctrl+Z #1
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Z",
		Modifiers: hook.ModCtrl,
		Timestamp: baseTime,
	})

	chords := eng.Chords()
	if len(chords) != 1 || chords[0].Streak != 1 {
		t.Fatalf("expected 1 chord with streak 1, got %d with streak %d", len(chords), chords[0].Streak)
	}

	// Press Ctrl+Z #2 (after 200ms)
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Z",
		Modifiers: hook.ModCtrl,
		Timestamp: baseTime.Add(200 * time.Millisecond),
	})

	chords = eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected in-place update, got %d chords", len(chords))
	}
	if chords[0].Streak != 2 {
		t.Errorf("expected streak 2, got %d", chords[0].Streak)
	}

	// Press Ctrl+Z #3 (after 400ms)
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Z",
		Modifiers: hook.ModCtrl,
		Timestamp: baseTime.Add(400 * time.Millisecond),
	})

	chords = eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected in-place update, got %d chords", len(chords))
	}
	if chords[0].Streak != 3 {
		t.Errorf("expected streak 3, got %d", chords[0].Streak)
	}

	// Press Ctrl+Y (different key) -> should create new chord
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Y",
		Modifiers: hook.ModCtrl,
		Timestamp: baseTime.Add(500 * time.Millisecond),
	})

	chords = eng.Chords()
	if len(chords) != 2 {
		t.Fatalf("expected 2 chords after new key, got %d", len(chords))
	}
	if chords[0].Streak != 3 {
		t.Errorf("first chord streak should remain 3, got %d", chords[0].Streak)
	}
	if chords[1].Streak != 1 {
		t.Errorf("second chord streak should be 1, got %d", chords[1].Streak)
	}

	// Press Ctrl+Y after 2 seconds (> 1.2s streak window) -> should create new chord
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Y",
		Modifiers: hook.ModCtrl,
		Timestamp: baseTime.Add(3000 * time.Millisecond),
	})

	chords = eng.Chords()
	if len(chords) != 3 {
		t.Fatalf("expected 3 chords due to streak expiration, got %d", len(chords))
	}
	if chords[2].Streak != 1 {
		t.Errorf("expired streak should reset to 1, got %d", chords[2].Streak)
	}
}

func TestBackspaceAbsorption(t *testing.T) {
	eng := NewEngine(WithFilterMode(FilterAll))
	defer eng.Stop()

	// Type "a", "b", "c"
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "a"})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "b"})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "c"})

	chords := eng.Chords()
	if len(chords) != 3 {
		t.Fatalf("expected 3 chords, got %d", len(chords))
	}
	if chords[2].Keys[0].Label != "c" {
		t.Errorf("expected last chord to be 'c', got %s", chords[2].Keys[0].Label)
	}

	// Hit Backspace -> should absorb 'c'
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Backspace"})
	chords = eng.Chords()
	if len(chords) != 2 {
		t.Fatalf("expected 2 chords after absorbing 'c', got %d", len(chords))
	}
	if chords[1].Keys[0].Label != "b" {
		t.Errorf("expected last chord to be 'b', got %s", chords[1].Keys[0].Label)
	}

	// Hit Backspace again -> should absorb 'b'
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Backspace"})
	chords = eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 chord after absorbing 'b', got %d", len(chords))
	}
	if chords[0].Keys[0].Label != "a" {
		t.Errorf("expected last chord to be 'a', got %s", chords[0].Keys[0].Label)
	}

	// Hit Backspace again -> should absorb 'a'
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Backspace"})
	chords = eng.Chords()
	if len(chords) != 0 {
		t.Fatalf("expected 0 chords after absorbing 'a', got %d", len(chords))
	}

	// Hit Backspace when buffer is empty -> should append Backspace key
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Backspace"})
	chords = eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 Backspace chord when buffer empty, got %d", len(chords))
	}
	if chords[0].Keys[0].Label != "Backspace" {
		t.Errorf("expected Backspace chord, got %s", chords[0].Keys[0].Label)
	}

	// Test absorption on chord with streak
	eng.Clear()
	baseTime := time.Now()
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "x", Timestamp: baseTime})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "x", Timestamp: baseTime.Add(100 * time.Millisecond)})
	chords = eng.Chords()
	if len(chords) != 1 || chords[0].Streak != 2 {
		t.Fatalf("expected 1 chord with streak 2, got %d chords streak %d", len(chords), chords[0].Streak)
	}

	// First backspace decrements streak to 1
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Backspace"})
	chords = eng.Chords()
	if len(chords) != 1 || chords[0].Streak != 1 {
		t.Fatalf("expected streak to decrement to 1, got len %d streak %d", len(chords), chords[0].Streak)
	}

	// Second backspace removes the chord entirely
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Backspace"})
	chords = eng.Chords()
	if len(chords) != 0 {
		t.Fatalf("expected 0 chords, got %d", len(chords))
	}

	// Test non-absorbable chord (e.g. shortcut Ctrl+C)
	eng.Clear()
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "c", Modifiers: hook.ModCtrl})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Backspace"})
	chords = eng.Chords()
	if len(chords) != 2 {
		t.Fatalf("expected Ctrl+C not to be absorbed, expected 2 chords, got %d", len(chords))
	}
}

func TestFilterModes(t *testing.T) {
	// 1. FilterSilent
	engSilent := NewEngine(WithFilterMode(FilterSilent))
	defer engSilent.Stop()

	engSilent.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "a"})
	engSilent.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Enter"})
	engSilent.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "C", Modifiers: hook.ModCtrl})
	if len(engSilent.Chords()) != 0 {
		t.Errorf("FilterSilent should produce 0 chords, got %d", len(engSilent.Chords()))
	}

	// 2. FilterSmart
	engSmart := NewEngine(WithFilterMode(FilterSmart), WithMaxChords(10))
	defer engSmart.Stop()

	// Plain typing: ignored
	engSmart.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "a"})
	engSmart.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "z"})
	if len(engSmart.Chords()) != 0 {
		t.Errorf("FilterSmart should ignore plain alphanumeric keys, got %d", len(engSmart.Chords()))
	}

	// Special keys: displayed
	engSmart.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Enter"})
	engSmart.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Esc"})
	engSmart.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Tab"})
	engSmart.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "F5"})
	engSmart.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "ArrowUp"})
	if len(engSmart.Chords()) != 5 {
		t.Errorf("FilterSmart should display special keys, expected 5, got %d", len(engSmart.Chords()))
	}

	// Shortcut with modifier: displayed
	engSmart.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "C", Modifiers: hook.ModCtrl})
	if len(engSmart.Chords()) != 6 {
		t.Errorf("FilterSmart should display shortcuts, expected 6, got %d", len(engSmart.Chords()))
	}

	// 3. FilterAll
	engAll := NewEngine(WithFilterMode(FilterAll))
	defer engAll.Stop()

	engAll.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "h"})
	engAll.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "i"})
	if len(engAll.Chords()) != 2 {
		t.Errorf("FilterAll should visualize all keystrokes, expected 2, got %d", len(engAll.Chords()))
	}

	// 4. FilterCustom
	engCustom := NewEngine(
		WithFilterMode(FilterCustom),
		WithAllowedKeys([]string{"Ctrl", "A", "Enter", "Left Click"}),
	)
	defer engCustom.Stop()

	// Allowed primary key: displayed
	engCustom.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "A"})
	if len(engCustom.Chords()) != 1 {
		t.Errorf("FilterCustom should allow 'A', got %d chords", len(engCustom.Chords()))
	}

	// Denied primary key: filtered
	engCustom.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "B"})
	if len(engCustom.Chords()) != 1 {
		t.Errorf("FilterCustom should ignore 'B', got %d chords", len(engCustom.Chords()))
	}

	// Allowed shortcut with modifier: displayed
	engCustom.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "C", Modifiers: hook.ModCtrl})
	if len(engCustom.Chords()) != 2 {
		t.Errorf("FilterCustom should allow Ctrl+C because Ctrl is allowed, got %d chords", len(engCustom.Chords()))
	}

	// Allowed mouse click: displayed
	engCustom.HandleEvent(hook.InputEvent{Type: hook.EventMouseDown, Button: hook.ButtonLeft})
	if len(engCustom.Chords()) != 3 {
		t.Errorf("FilterCustom should allow Left Click, got %d chords", len(engCustom.Chords()))
	}

	// Denied mouse click: filtered
	engCustom.HandleEvent(hook.InputEvent{Type: hook.EventMouseDown, Button: hook.ButtonRight})
	if len(engCustom.Chords()) != 3 {
		t.Errorf("FilterCustom should block Right Click, got %d chords", len(engCustom.Chords()))
	}

	// Update AllowedKeys dynamically
	engCustom.SetAllowedKeys([]string{"Ctrl", "A", "B", "Enter", "Left Click"})
	engCustom.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "B"})
	if len(engCustom.Chords()) != 4 {
		t.Errorf("FilterCustom should now allow 'B' after SetAllowedKeys, got %d chords", len(engCustom.Chords()))
	}
}

func TestSetLingerDurationControlsExpiry(t *testing.T) {
	eng := NewEngine(WithFilterMode(FilterAll))
	defer eng.Stop()

	eng.SetLingerDuration(1200)
	now := time.Now()
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "A", Timestamp: now})

	chords := eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 chord, got %d", len(chords))
	}
	life := chords[0].ExpiresAt.Sub(chords[0].Timestamp)
	if life < 1100*time.Millisecond || life > 1400*time.Millisecond {
		t.Errorf("expected ~1.2s linger, got %v", life)
	}
}

func TestSetMaxChordsTrimsHistory(t *testing.T) {
	eng := NewEngine(WithFilterMode(FilterAll))
	defer eng.Stop()

	eng.SetMaxChords(2)
	now := time.Now()
	for i, key := range []string{"A", "B", "C", "D"} {
		eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: key, Timestamp: now.Add(time.Duration(i) * 10 * time.Millisecond)})
	}

	chords := eng.Chords()
	if len(chords) != 2 {
		t.Fatalf("expected 2 chords after cap, got %d", len(chords))
	}
	if chords[0].Keys[0].Label != "C" || chords[1].Keys[0].Label != "D" {
		t.Errorf("expected newest chords C,D; got %s,%s", chords[0].Keys[0].Label, chords[1].Keys[0].Label)
	}
}

func TestDecayRetiresStaleMouseFeedback(t *testing.T) {
	eng := NewEngine(WithShowMouseClicks(true))
	defer eng.Stop()

	now := time.Now()
	eng.HandleEvent(hook.InputEvent{Type: hook.EventMouseDown, Button: hook.ButtonLeft, X: 100, Y: 100, Timestamp: now})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventMouseUp, Button: hook.ButtonLeft, X: 100, Y: 100, Timestamp: now.Add(40 * time.Millisecond)})

	if eng.State().MouseState.Timestamp.IsZero() {
		t.Fatal("expected mouse feedback timestamp after click")
	}

	// Before the retirement window it remains visible for the exit animation.
	eng.DecayCheck(now.Add(200 * time.Millisecond))
	if eng.State().MouseState.Timestamp.IsZero() {
		t.Fatal("mouse feedback retired too early")
	}

	// After ~320 ms the frame must be cleared, otherwise the ring/indicator
	// stays frozen on screen until the next unrelated input.
	eng.DecayCheck(now.Add(500 * time.Millisecond))
	st := eng.State().MouseState
	if !st.Timestamp.IsZero() || st.Button != hook.ButtonNone || st.Wheel != 0 || st.Active {
		t.Errorf("stale mouse feedback was not retired: %+v", st)
	}
}

func TestPrivacyMasking(t *testing.T) {
	matcher := newMockMatcher()
	matcher.secrets["secret_key"] = true

	eng := NewEngine(
		WithFilterMode(FilterAll),
		WithMatcher(matcher),
	)
	defer eng.Stop()

	// Normal key
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "A"})
	chords := eng.Chords()
	if chords[0].IsPrivate || chords[0].Keys[0].Label != "A" {
		t.Errorf("normal key should not be private")
	}

	// Sensitive key
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "secret_key"})
	chords = eng.Chords()
	if len(chords) != 2 {
		t.Fatalf("expected 2 chords, got %d", len(chords))
	}

	privChord := chords[1]
	if !privChord.IsPrivate {
		t.Errorf("expected IsPrivate=true for sensitive key")
	}
	if len(privChord.Keys) != 1 || privChord.Keys[0].Label != "[🔒 Private]" {
		t.Errorf("expected key to be replaced by [🔒 Private], got %+v", privChord.Keys)
	}

	// Releasing sensitive key marks chord inactive
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyUp, KeyName: "secret_key"})
	chords = eng.Chords()
	if chords[1].IsActive {
		t.Errorf("expected private chord to become inactive on release")
	}
}

func TestDecayExpirationCleanup(t *testing.T) {
	// Background decay test
	eng := NewEngine(
		WithFilterMode(FilterSmart),
		WithLingerShort(35*time.Millisecond),
		WithDecayInterval(10*time.Millisecond),
	)
	defer eng.Stop()

	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Enter"})
	if len(eng.Chords()) != 1 {
		t.Fatalf("expected 1 chord initially")
	}

	// Sleep beyond expiration
	time.Sleep(80 * time.Millisecond)
	if len(eng.Chords()) != 0 {
		t.Errorf("expected chord to be removed after decay interval, got %d", len(eng.Chords()))
	}

	// Synchronous DecayCheck test
	now := time.Now()
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Tab", Timestamp: now})
	if len(eng.Chords()) != 1 {
		t.Fatalf("expected 1 chord")
	}

	removed := eng.DecayCheck(now.Add(5 * time.Second))
	if removed != 1 {
		t.Errorf("expected DecayCheck to remove 1 expired chord, removed %d", removed)
	}
	if len(eng.Chords()) != 0 {
		t.Errorf("expected 0 chords after DecayCheck, got %d", len(eng.Chords()))
	}

	// Clear test
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "F1"})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "F2"})
	if len(eng.Chords()) != 2 {
		t.Fatalf("expected 2 chords before Clear")
	}
	eng.Clear()
	if len(eng.Chords()) != 0 {
		t.Errorf("expected 0 chords after Clear, got %d", len(eng.Chords()))
	}
}

func TestSemanticActionMatchAndAdaptiveLinger(t *testing.T) {
	matcher := newMockMatcher()
	matcher.actions["Ctrl+C"] = semantic.ActionInfo{
		Title:    "Copy",
		Category: "Clipboard",
	}

	eng := NewEngine(
		WithFilterMode(FilterSmart),
		WithMatcher(matcher),
		WithLingerShort(1800*time.Millisecond),
		WithLingerLong(2500*time.Millisecond),
	)
	defer eng.Stop()

	start := time.Now()
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "C",
		Modifiers: hook.ModCtrl,
		Timestamp: start,
	})

	chords := eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 chord, got %d", len(chords))
	}
	chord := chords[0]
	if chord.Action.Title != "Copy" {
		t.Errorf("expected Action.Title 'Copy', got '%s'", chord.Action.Title)
	}

	// Check adaptive linger duration (should be ~2.5s for matched action)
	expectedExpiry := start.Add(2500 * time.Millisecond)
	diff := chord.ExpiresAt.Sub(expectedExpiry)
	if diff < -50*time.Millisecond || diff > 50*time.Millisecond {
		t.Errorf("expected adaptive linger ~2.5s, got %v", chord.ExpiresAt.Sub(start))
	}
}

func TestAcousticFeedback(t *testing.T) {
	player := &mockSoundPlayer{}
	eng := NewEngine(
		WithFilterMode(FilterAll),
		WithSoundPlayer(player),
		WithSoundActive(true),
	)
	defer eng.Stop()

	// Tap modifier -> Thock
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Ctrl", IsModifier: true})
	if player.count(SoundThock) != 1 {
		t.Errorf("expected 1 Thock sound for modifier press, got %d", player.count(SoundThock))
	}

	// Tap normal key -> Click
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "A"})
	if player.count(SoundClick) != 1 {
		t.Errorf("expected 1 Click sound for regular key press, got %d", player.count(SoundClick))
	}

	// Mute sound
	eng.SetSoundActive(false)
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "B"})
	if player.count(SoundClick) != 1 {
		t.Errorf("expected sound count to remain 1 when muted, got %d", player.count(SoundClick))
	}
}

func TestMouseVisualState(t *testing.T) {
	eng := NewEngine()
	defer eng.Stop()

	now := time.Now()
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventMouseDown,
		Button:    hook.ButtonLeft,
		X:         100,
		Y:         200,
		Modifiers: hook.ModCtrl,
		Timestamp: now,
	})

	ms := eng.State().MouseState
	if ms.Button != hook.ButtonLeft || ms.X != 100 || ms.Y != 200 || !ms.Active || ms.Modifiers != hook.ModCtrl {
		t.Errorf("unexpected mouse state on MouseDown: %+v", ms)
	}

	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventMouseUp,
		Button:    hook.ButtonLeft,
		X:         105,
		Y:         205,
		Timestamp: now.Add(50 * time.Millisecond),
	})

	ms = eng.State().MouseState
	if ms.Active {
		t.Errorf("expected MouseState.Active to be false on MouseUp")
	}
	if ms.X != 105 || ms.Y != 205 {
		t.Errorf("expected updated mouse coords, got (%d, %d)", ms.X, ms.Y)
	}
}

func TestSubscriptionAndUIState(t *testing.T) {
	eng := NewEngine(WithFilterMode(FilterSmart))
	defer eng.Stop()

	sub := eng.Subscribe()
	defer eng.Unsubscribe(sub)

	// Subscriber receives initial snapshot immediately
	select {
	case initial := <-sub:
		if len(initial.Chords) != 0 {
			t.Errorf("expected 0 chords in initial state, got %d", len(initial.Chords))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timed out waiting for initial state")
	}

	// Dispatch an event
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Enter"})

	select {
	case updated := <-sub:
		if len(updated.Chords) != 1 {
			t.Errorf("expected 1 chord in updated state, got %d", len(updated.Chords))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timed out waiting for state update")
	}
}

func TestMaxChordsLimit(t *testing.T) {
	eng := NewEngine(
		WithFilterMode(FilterAll),
		WithMaxChords(3),
	)
	defer eng.Stop()

	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "1"})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "2"})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "3"})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "4"})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "5"})

	chords := eng.Chords()
	if len(chords) != 3 {
		t.Fatalf("expected max 3 chords, got %d", len(chords))
	}

	// Oldest 1 and 2 should have been pushed out
	expected := []string{"3", "4", "5"}
	for i, exp := range expected {
		if chords[i].Keys[0].Label != exp {
			t.Errorf("chord %d expected %s, got %s", i, exp, chords[i].Keys[0].Label)
		}
	}
}

func TestMouseActionsAndWheel(t *testing.T) {
	matcher := semantic.NewDefaultMatcher()
	eng := NewEngine(
		WithFilterMode(FilterSmart),
		WithMatcher(matcher),
		WithShowMouseClicks(true),
	)
	defer eng.Stop()

	// 1. Mouse left click alone when showMouseClicks is true
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventMouseDown,
		Button:    hook.ButtonLeft,
		Timestamp: time.Now(),
	})

	chords := eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 chord for left click, got %d", len(chords))
	}
	if chords[0].Keys[0].Label != "Left Click" {
		t.Fatalf("expected 'Left Click', got %q", chords[0].Keys[0].Label)
	}

	// Release left click
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventMouseUp,
		Button:    hook.ButtonLeft,
		Timestamp: time.Now(),
	})

	// 2. Ctrl + Scroll Up -> Zoom In
	eng.HandleEvent(hook.InputEvent{
		Type:       hook.EventMouseWheel,
		WheelDelta: 120,
		Modifiers:  hook.ModCtrl,
		Timestamp:  time.Now(),
	})

	chords = eng.Chords()
	if len(chords) != 2 {
		t.Fatalf("expected 2 chords, got %d", len(chords))
	}
	last := chords[len(chords)-1]
	if last.Action.Title != "Zoom In" {
		t.Fatalf("expected 'Zoom In', got %q", last.Action.Title)
	}
	if len(last.Keys) != 2 || last.Keys[0].Label != "Ctrl" || last.Keys[1].Label != "Scroll Up" {
		t.Fatalf("expected [Ctrl, Scroll Up], got %+v", last.Keys)
	}
}

func TestMouseDisabledByDefault(t *testing.T) {
	eng := NewEngine(WithFilterMode(FilterAll))
	defer eng.Stop()

	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventMouseDown,
		Button:    hook.ButtonLeft,
		Timestamp: time.Now(),
	})
	if len(eng.Chords()) != 0 {
		t.Fatalf("expected 0 chords for mouse click when showMouseClicks is false, got %d", len(eng.Chords()))
	}

	eng.HandleEvent(hook.InputEvent{
		Type:       hook.EventMouseWheel,
		WheelDelta: 120,
		Timestamp:  time.Now(),
	})
	if len(eng.Chords()) != 0 {
		t.Fatalf("expected 0 chords for mouse wheel when showMouseClicks is false, got %d", len(eng.Chords()))
	}
}

func TestWordGrouping_BasicAndBackspace(t *testing.T) {
	eng := NewEngine(WithWordGrouping(true))
	defer eng.Stop()

	if !eng.WordGrouping() {
		t.Fatalf("expected WordGrouping to be enabled")
	}

	baseTime := time.Date(2026, 11, 5, 14, 0, 0, 0, time.UTC)

	// Type "h", "e", "l", "l", "o"
	chars := []string{"h", "e", "l", "l", "o"}
	for i, ch := range chars {
		eng.HandleEvent(hook.InputEvent{
			Type:      hook.EventKeyDown,
			KeyName:   ch,
			Timestamp: baseTime.Add(time.Duration(i*50) * time.Millisecond),
		})
	}

	chords := eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 word chord, got %d", len(chords))
	}
	if len(chords[0].Keys) != 1 {
		t.Fatalf("expected 1 KeyCap in word chord, got %d", len(chords[0].Keys))
	}
	if chords[0].Keys[0].Label != "hello" {
		t.Errorf("expected chord label 'hello', got %q", chords[0].Keys[0].Label)
	}
	if !chords[0].IsActive || !chords[0].Keys[0].IsHeld {
		t.Errorf("expected active and held while 'o' is down")
	}

	// Release 'o'
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyUp,
		KeyName:   "o",
		Timestamp: baseTime.Add(260 * time.Millisecond),
	})
	chords = eng.Chords()
	if chords[0].IsActive || chords[0].Keys[0].IsHeld {
		t.Errorf("expected inactive and unheld after release")
	}

	// Backspace #1: "hello" -> "hell"
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Backspace",
		Timestamp: baseTime.Add(300 * time.Millisecond),
	})
	chords = eng.Chords()
	if len(chords) != 1 || chords[0].Keys[0].Label != "hell" {
		t.Fatalf("expected 'hell' after 1 backspace, got %+v", chords)
	}

	// Backspace #2: "hell" -> "hel"
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Backspace",
		Timestamp: baseTime.Add(350 * time.Millisecond),
	})
	chords = eng.Chords()
	if chords[0].Keys[0].Label != "hel" {
		t.Errorf("expected 'hel', got %q", chords[0].Keys[0].Label)
	}

	// Backspace #3: "hel" -> "he"
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Backspace",
		Timestamp: baseTime.Add(400 * time.Millisecond),
	})
	chords = eng.Chords()
	if chords[0].Keys[0].Label != "he" {
		t.Errorf("expected 'he', got %q", chords[0].Keys[0].Label)
	}

	// Backspace #4: "he" -> "h"
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Backspace",
		Timestamp: baseTime.Add(450 * time.Millisecond),
	})
	chords = eng.Chords()
	if chords[0].Keys[0].Label != "h" {
		t.Errorf("expected 'h', got %q", chords[0].Keys[0].Label)
	}

	// Backspace #5: "h" -> removed entirely
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Backspace",
		Timestamp: baseTime.Add(500 * time.Millisecond),
	})
	chords = eng.Chords()
	if len(chords) != 0 {
		t.Errorf("expected chord removed after deleting last character, got %d chords", len(chords))
	}
}

func TestWordGrouping_UppercaseShiftAndSpace(t *testing.T) {
	eng := NewEngine(WithWordGrouping(true))
	defer eng.Stop()

	baseTime := time.Date(2026, 11, 5, 14, 0, 0, 0, time.UTC)

	// Press Shift+H (uppercase letter)
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "H",
		Modifiers: hook.ModShift,
		Timestamp: baseTime,
	})

	// Followed by lowercase "i", Space, "4", "2"
	steps := []struct {
		key  string
		mods hook.Modifiers
	}{
		{"i", hook.ModNone},
		{" ", hook.ModNone},
		{"4", hook.ModNone},
		{"2", hook.ModNone},
	}

	for i, step := range steps {
		eng.HandleEvent(hook.InputEvent{
			Type:      hook.EventKeyDown,
			KeyName:   step.key,
			Modifiers: step.mods,
			Timestamp: baseTime.Add(time.Duration((i+1)*50) * time.Millisecond),
		})
	}

	chords := eng.Chords()
	// Continuous alphanumeric groups: "Hi" is separated from "42" by non-alphanumeric space
	if len(chords) < 2 {
		t.Fatalf("expected at least 2 distinct word/action chords, got %d", len(chords))
	}
	if chords[0].Keys[0].Label != "Hi" {
		t.Errorf("expected first chord label 'Hi', got %q", chords[0].Keys[0].Label)
	}
	lastIdx := len(chords) - 1
	if chords[lastIdx].Keys[0].Label != "42" {
		t.Errorf("expected last chord label '42', got %q", chords[lastIdx].Keys[0].Label)
	}
}

func TestWordGrouping_CadenceThreshold(t *testing.T) {
	eng := NewEngine(WithWordGrouping(true))
	defer eng.Stop()

	baseTime := time.Date(2026, 11, 5, 14, 0, 0, 0, time.UTC)

	// Type "h" at t=0
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "h",
		Timestamp: baseTime,
	})

	// Wait 800ms (> 600ms wordThreshold) and type "i"
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "i",
		Timestamp: baseTime.Add(800 * time.Millisecond),
	})

	chords := eng.Chords()
	if len(chords) != 2 {
		t.Fatalf("expected 2 chords due to threshold timeout, got %d", len(chords))
	}
	if chords[0].Keys[0].Label != "h" {
		t.Errorf("expected first chord 'h', got %q", chords[0].Keys[0].Label)
	}
	if chords[1].Keys[0].Label != "i" {
		t.Errorf("expected second chord 'i', got %q", chords[1].Keys[0].Label)
	}
}

func TestWordGrouping_ModifierShortcutInterruption(t *testing.T) {
	eng := NewEngine(WithWordGrouping(true))
	defer eng.Stop()

	baseTime := time.Date(2026, 11, 5, 14, 0, 0, 0, time.UTC)

	// Type "h", "e", "l"
	for i, ch := range []string{"h", "e", "l"} {
		eng.HandleEvent(hook.InputEvent{
			Type:      hook.EventKeyDown,
			KeyName:   ch,
			Timestamp: baseTime.Add(time.Duration(i*50) * time.Millisecond),
		})
	}

	// Press Ctrl+C shortcut
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "C",
		Modifiers: hook.ModCtrl,
		Timestamp: baseTime.Add(200 * time.Millisecond),
	})

	chords := eng.Chords()
	if len(chords) != 2 {
		t.Fatalf("expected 2 chords (word + shortcut), got %d", len(chords))
	}
	if chords[0].Keys[0].Label != "hel" {
		t.Errorf("expected word chord 'hel', got %q", chords[0].Keys[0].Label)
	}
	if len(chords[1].Keys) != 2 || chords[1].Keys[0].Label != "Ctrl" || chords[1].Keys[1].Label != "C" {
		t.Errorf("expected [Ctrl, C] shortcut chord, got %+v", chords[1].Keys)
	}

	// Resume typing "p" -> should start a new word chord
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "p",
		Timestamp: baseTime.Add(300 * time.Millisecond),
	})

	chords = eng.Chords()
	if len(chords) != 3 {
		t.Fatalf("expected 3 chords, got %d", len(chords))
	}
	if chords[2].Keys[0].Label != "p" {
		t.Errorf("expected third chord 'p', got %q", chords[2].Keys[0].Label)
	}
}

func TestWordGrouping_SetWordGroupingToggle(t *testing.T) {
	eng := NewEngine(WithFilterMode(FilterSmart))
	defer eng.Stop()

	if eng.WordGrouping() {
		t.Fatalf("expected WordGrouping to be false by default")
	}

	// When WordGrouping is false in FilterSmart mode, plain typing is ignored
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "a"})
	if len(eng.Chords()) != 0 {
		t.Fatalf("plain typing should be ignored in Smart mode when wordGrouping is false")
	}

	// Enable WordGrouping at runtime
	eng.SetWordGrouping(true)
	if !eng.WordGrouping() {
		t.Fatalf("expected WordGrouping to be true after SetWordGrouping(true)")
	}

	// Now plain typing is grouped and displayed
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "a"})
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "b"})
	chords := eng.Chords()
	if len(chords) != 1 || chords[0].Keys[0].Label != "ab" {
		t.Fatalf("expected 1 word chord 'ab', got %+v", chords)
	}
}

func TestSpatialSoundPlayerDispatch(t *testing.T) {
	var playedKeys []string
	var playedTypes []SoundType

	spatialMock := SpatialSoundPlayerFunc(func(s SoundType, key string) {
		playedTypes = append(playedTypes, s)
		playedKeys = append(playedKeys, key)
	})

	eng := NewEngine(
		WithFilterMode(FilterAll),
		WithSoundPlayer(spatialMock),
		WithSoundActive(true),
	)

	// Modifier key should trigger Thock with key name
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Ctrl", IsModifier: true})
	// Alphanumeric key should trigger Click with key name
	eng.HandleEvent(hook.InputEvent{Type: hook.EventKeyDown, KeyName: "Enter"})

	if len(playedKeys) != 2 {
		t.Fatalf("expected 2 sound events, got %d", len(playedKeys))
	}
	if playedTypes[0] != SoundThock || playedKeys[0] != "ctrl" {
		t.Errorf("first sound mismatch: got (%v, %q), want (Thock, 'ctrl')", playedTypes[0], playedKeys[0])
	}
	if playedTypes[1] != SoundClick || playedKeys[1] != "Enter" {
		t.Errorf("second sound mismatch: got (%v, %q), want (Click, 'Enter')", playedTypes[1], playedKeys[1])
	}
}

func TestCadenceVelocityAndFinisher(t *testing.T) {
	eng := NewEngine(
		WithFilterMode(FilterAll),
	)

	now := time.Now()

	// 1. Single slow keypress -> Velocity 0 (Adagio)
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "A",
		Timestamp: now,
	})
	chords := eng.Chords()
	if len(chords) != 1 {
		t.Fatalf("expected 1 chord, got %d", len(chords))
	}
	if chords[0].Velocity != 0 {
		t.Errorf("single press should be Velocity 0 (Adagio), got %d", chords[0].Velocity)
	}
	if chords[0].IsFinisher {
		t.Errorf("'A' should not be marked as finisher")
	}

	// 2. High-speed burst -> Velocity 2 (Allegro, >75 WPM)
	// 10 keystrokes in 0.5s = (10 * 12) / 0.5 = 240 WPM
	burstTime := now.Add(1 * time.Second)
	for i := 0; i < 10; i++ {
		eng.HandleEvent(hook.InputEvent{
			Type:      hook.EventKeyDown,
			KeyName:   fmt.Sprintf("k%d", i),
			Timestamp: burstTime.Add(time.Duration(i*40) * time.Millisecond),
		})
	}
	chords = eng.Chords()
	last := chords[len(chords)-1]
	if last.Velocity != 2 {
		t.Errorf("high cadence burst should be Velocity 2 (Allegro), got %d (wpm=%d)", last.Velocity, eng.CadenceWPM())
	}

	// 3. Finisher key -> Enter
	eng.HandleEvent(hook.InputEvent{
		Type:      hook.EventKeyDown,
		KeyName:   "Enter",
		Timestamp: burstTime.Add(500 * time.Millisecond),
	})
	chords = eng.Chords()
	finisher := chords[len(chords)-1]
	if !finisher.IsFinisher {
		t.Errorf("Enter chord should have IsFinisher=true")
	}
}
