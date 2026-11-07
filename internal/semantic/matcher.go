package semantic

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// DefaultMatcher implements Matcher with a curated, context-aware shortcut
// registry and context/keystroke-based privacy detection.
//
// Hot path: MatchAction performs at most a few map lookups. Expensive context
// analysis (profile resolution, sensitive-window keyword scan) is cached and
// recomputed only when the foreground AppContext changes.
type DefaultMatcher struct {
	mu        sync.Mutex
	registry  *Registry
	overrides overrideSet

	// Context cache, invalidated by AppContext.Seq.
	ctxSeq       uint64
	ctxKey       string
	ctxProfile   *Profile
	ctxCategory  context.Category
	ctxSensitive bool

	// Keystroke history and privacy state.
	history       []rune
	sensitiveLeft int
	lastInputTime time.Time
}

var _ Matcher = (*DefaultMatcher)(nil)

// NewMatcher creates a DefaultMatcher with the standard registry.
func NewMatcher() *DefaultMatcher {
	return &DefaultMatcher{
		registry: NewRegistry(),
		history:  make([]rune, 0, 64),
	}
}

// NewDefaultMatcher is an alias for NewMatcher.
func NewDefaultMatcher() *DefaultMatcher {
	return NewMatcher()
}

// Registry returns the immutable shortcut knowledge base.
func (m *DefaultMatcher) Registry() *Registry {
	return m.registry
}

// SetOverrides installs user-defined shortcut interpretations. Overrides are
// indexed once and then resolved with the same priority as built-in profiles,
// but above them (user intent always wins).
func (m *DefaultMatcher) SetOverrides(list []Override) {
	set := buildOverrideSet(list)
	m.mu.Lock()
	m.overrides = set
	m.mu.Unlock()
}

// ParseChord parses a human-readable chord such as "Ctrl+Shift+P" into a
// modifier mask and canonical primary key. Shifted symbols fold into the Shift
// modifier, so "Ctrl+Shift+1" and "Ctrl+!" parse identically.
func ParseChord(chord string) (hook.Modifiers, string, error) {
	s := strings.TrimSpace(chord)
	if s == "" {
		return hook.ModNone, "", fmt.Errorf("empty chord")
	}
	parts := strings.Split(s, "+")
	var mods hook.Modifiers
	var keys []string
	for i, raw := range parts {
		token := strings.TrimSpace(raw)
		if token == "" {
			// A trailing empty token represents a literal '+' key ("Ctrl++").
			if i > 0 && i == len(parts)-1 {
				keys = append(keys, "+")
			}
			continue
		}
		switch strings.ToLower(token) {
		case "ctrl", "control", "controlleft", "controlright":
			mods |= hook.ModCtrl
		case "shift", "shiftleft", "shiftright":
			mods |= hook.ModShift
		case "alt", "menu", "option", "altleft", "altright":
			mods |= hook.ModAlt
		case "win", "super", "cmd", "command", "meta", "metaleft", "metaright":
			mods |= hook.ModSuper
		default:
			keys = append(keys, token)
		}
	}
	if len(keys) != 1 {
		return hook.ModNone, "", fmt.Errorf("chord %q must contain exactly one non-modifier key", chord)
	}
	key, impliedShift := canonicalKey(keys[0])
	if key == "" {
		return hook.ModNone, "", fmt.Errorf("chord %q has an empty key", chord)
	}
	if impliedShift {
		mods |= hook.ModShift
	}
	return mods, key, nil
}

// ParseOverride converts a user configuration entry into an Override.
func ParseOverride(app, chord, label, description string) (Override, error) {
	mods, key, err := ParseChord(chord)
	if err != nil {
		return Override{}, err
	}
	label = strings.TrimSpace(label)
	if label == "" {
		return Override{}, fmt.Errorf("custom shortcut %q has an empty label", chord)
	}
	return Override{
		App:  strings.ToLower(strings.TrimSpace(app)),
		Mods: mods,
		Key:  key,
		Action: ActionInfo{
			Title:       label,
			Category:    "Custom",
			Description: description,
		},
	}, nil
}

// overrideSet is a pre-indexed, immutable set of user overrides.
type overrideSet struct {
	byApp map[string]map[comboKey]ActionInfo
	count int
}

func buildOverrideSet(list []Override) overrideSet {
	set := overrideSet{byApp: make(map[string]map[comboKey]ActionInfo)}
	for _, o := range list {
		key := strings.ToLower(strings.TrimSpace(o.App))
		if key == "*" {
			key = ""
		}
		ck := canonicalCombo(o.Mods, o.Key)
		if ck.key == "" {
			continue
		}
		if set.byApp[key] == nil {
			set.byApp[key] = make(map[comboKey]ActionInfo)
		}
		set.byApp[key][ck] = o.Action
		set.count++
	}
	return set
}

func (o overrideSet) lookup(app context.AppContext, profile *Profile, category context.Category, combo comboKey) (ActionInfo, bool) {
	if o.count == 0 {
		return ActionInfo{}, false
	}
	if exe := context.NormalizeExe(app.Exe); exe != "" {
		if m, ok := o.byApp[exe]; ok {
			if a, found := m[combo]; found {
				return a, true
			}
		}
	}
	if profile != nil {
		if a, found := o.lookupApp(profile.ID, combo); found {
			return a, true
		}
	}
	if category != "" {
		if a, found := o.lookupApp(string(category), combo); found {
			return a, true
		}
	}
	if a, found := o.lookupApp("", combo); found {
		return a, true
	}
	return ActionInfo{}, false
}

func (o overrideSet) lookupApp(app string, combo comboKey) (ActionInfo, bool) {
	m, ok := o.byApp[strings.ToLower(strings.TrimSpace(app))]
	if !ok {
		return ActionInfo{}, false
	}
	a, found := m[combo]
	return a, found
}

// prepareLocked resolves the AppContext once per foreground change.
func (m *DefaultMatcher) prepareLocked(app context.AppContext) {
	if app.Seq != 0 {
		if m.ctxSeq == app.Seq {
			return
		}
		m.ctxSeq = app.Seq
	} else {
		key := app.Key()
		if key == m.ctxKey {
			return
		}
		m.ctxKey = key
	}
	m.ctxProfile, m.ctxCategory = m.registry.Resolve(app)
	m.ctxSensitive = isSensitiveContext(app)
}

func resultFor(a ActionInfo, confidence Confidence, profileID, profileName string, app context.AppContext) MatchResult {
	appName := profileName
	if appName == "" {
		appName = app.DisplayName()
	}
	return MatchResult{
		Action:      a,
		Confidence:  confidence,
		ProfileID:   profileID,
		ProfileName: profileName,
		AppName:     appName,
	}
}

// MatchAction returns an ActionInfo if the combo corresponds to a known shortcut.
func (m *DefaultMatcher) MatchAction(app context.AppContext, mods hook.Modifiers, keyName string) MatchResult {
	normKey := NormalizeKey(keyName)
	if normKey == "" {
		return MatchResult{}
	}
	combo := canonicalCombo(mods, keyName)

	m.mu.Lock()
	m.prepareLocked(app)
	profile := m.ctxProfile
	category := m.ctxCategory
	if a, ok := m.overrides.lookup(app, profile, category, combo); ok {
		res := resultFor(a, ConfidenceUser, profileMetaID(profile), profileMetaName(profile), app)
		m.mu.Unlock()
		return res
	}
	a, confidence, profileID, profileName := m.registry.Lookup(profile, category, combo)
	m.mu.Unlock()

	return resultFor(a, confidence, profileID, profileName, app)
}

// MatchUniversal resolves a chord using only global user overrides and the
// universal dictionary. It deliberately bypasses the AppContext cache so it can
// coexist with IsSensitive(realContext) without invalidating privacy state.
func (m *DefaultMatcher) MatchUniversal(mods hook.Modifiers, keyName string) MatchResult {
	normKey := NormalizeKey(keyName)
	if normKey == "" {
		return MatchResult{}
	}
	combo := canonicalCombo(mods, keyName)

	m.mu.Lock()
	if a, ok := m.overrides.lookupApp("", combo); ok {
		m.mu.Unlock()
		return resultFor(a, ConfidenceUser, "", "", context.AppContext{})
	}
	a, confidence, _, _ := m.registry.Lookup(nil, context.CategoryUnknown, combo)
	m.mu.Unlock()

	return resultFor(a, confidence, "", "", context.AppContext{})
}

func profileMetaID(p *Profile) string {
	if p == nil {
		return ""
	}
	return p.ID
}

func profileMetaName(p *Profile) string {
	if p == nil {
		return ""
	}
	return p.Name
}

// MatchSequence interprets a two-chord shortcut. The engine enforces the time
// window; this method only resolves the data.
func (m *DefaultMatcher) MatchSequence(
	app context.AppContext,
	firstMods hook.Modifiers, firstKey string,
	secondMods hook.Modifiers, secondKey string,
) (MatchResult, bool) {
	first := canonicalCombo(firstMods, firstKey)
	second := canonicalCombo(secondMods, secondKey)
	if first.key == "" || second.key == "" {
		return MatchResult{}, false
	}
	seq := sequenceKey{first: first, second: second}

	m.mu.Lock()
	m.prepareLocked(app)
	profile := m.ctxProfile
	category := m.ctxCategory
	a, confidence, profileID, profileName := m.registry.LookupSequence(profile, category, seq)
	m.mu.Unlock()
	if a.Title == "" {
		return MatchResult{}, false
	}
	res := resultFor(a, confidence, profileID, profileName, app)
	res.IsSequence = true
	return res, true
}

// sensitiveContextKeywords contains substrings indicative of credentials, logins, or finance.
var sensitiveContextKeywords = []string{
	// Password managers
	"1password", "bitwarden", "keepass", "lastpass", "dashlane", "enpass",
	"nordpass", "keeper", "roboform", "authenticator",

	// Authentication & logins
	"password", "passwd", "passcode", "pin code", "pincode", "master password",
	"sign in", "signin", "log in", "login", "authenticate", "authentication",
	"credential", "two-factor", "2fa", "mfa", "verification code", "one-time password", "otp",

	// Financial / Banking / Payment
	"bank", "banking", "checkout", "paypal", "credit card", "debit card", "cvv",
	"wallet", "binance", "coinbase", "metamask",

	// Browsing privacy modes
	"incognito", "inprivate", "private browsing", "private window",

	// OS security prompts
	"user account control", "windows security", "credential manager", "pin entry", "sudo password",
}

// isSensitiveContext scans window/process metadata. It runs only when the
// foreground app changes, never on every keystroke.
func isSensitiveContext(app context.AppContext) bool {
	if app.Title == "" && app.Exe == "" && app.AppName == "" {
		return false
	}
	text := strings.ToLower(app.Title + "\x00" + app.Exe + "\x00" + app.AppName)
	for _, kw := range sensitiveContextKeywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

// sensitiveKeyTriggers contains trigger sequences in typed history.
var sensitiveKeyTriggers = []string{
	"password", "passwd", "pwd:", "pwd ", "pin:", "secret", "token",
	"api_key", "apikey", "bearer ", "private_key", "sudo ", "cvv",
}

// IsSensitive returns true if current window title / context indicates
// password/credentials/banking or if the user has typed sensitive patterns.
func (m *DefaultMatcher) IsSensitive(app context.AppContext, keyName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.prepareLocked(app)
	if m.ctxSensitive {
		return true
	}

	// Direct sensitive keyName.
	keyLower := strings.ToLower(keyName)
	if keyLower == "password" || keyLower == "secret" || keyLower == "token" ||
		keyLower == "credential" || keyLower == "pin" || strings.Contains(keyLower, "password") {
		return true
	}

	// Keystroke pattern tracking.
	now := time.Now()
	if now.Sub(m.lastInputTime) > 5*time.Second {
		m.history = m.history[:0]
	}
	m.lastInputTime = now

	if m.sensitiveLeft > 0 {
		m.sensitiveLeft--
		if norm := NormalizeKey(keyName); norm == "ENTER" || norm == "ESC" {
			m.sensitiveLeft = 0
			m.history = m.history[:0]
		}
		return true
	}

	norm := NormalizeKey(keyName)
	if norm == "BACKSPACE" {
		if len(m.history) > 0 {
			m.history = m.history[:len(m.history)-1]
		}
		return false
	}
	if norm == "ENTER" || norm == "ESC" {
		m.history = m.history[:0]
		return false
	}

	runes := []rune(keyName)
	if len(runes) == 1 {
		r := unicode.ToLower(runes[0])
		m.history = append(m.history, r)
		if len(m.history) > 64 {
			m.history = m.history[len(m.history)-64:]
		}
		histStr := string(m.history)
		for _, trig := range sensitiveKeyTriggers {
			if strings.HasSuffix(histStr, trig) || strings.Contains(histStr, trig) {
				m.sensitiveLeft = 32
				return true
			}
		}
		if hasConsecutiveDigits(histStr, 8) {
			m.sensitiveLeft = 16
			return true
		}
	}
	return false
}

// Reset clears the internal history and privacy state.
func (m *DefaultMatcher) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = m.history[:0]
	m.sensitiveLeft = 0
	m.lastInputTime = time.Time{}
}

// AllActionTitles returns every registered action title for localization
// completeness checks.
func (m *DefaultMatcher) AllActionTitles() []string {
	return m.registry.AllActionTitles()
}

func hasConsecutiveDigits(s string, minDigits int) bool {
	count := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			count++
			if count >= minDigits {
				return true
			}
		} else if r != '-' && r != ' ' { // ignore separators
			count = 0
		}
	}
	return false
}
