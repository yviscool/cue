package semantic

import (
	"sort"
	"strings"

	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// comboKey identifies one chord: modifier mask + canonical physical key.
type comboKey struct {
	mods hook.Modifiers
	key  string
}

// sequenceKey identifies a two-chord shortcut (e.g. Ctrl+K Ctrl+S).
type sequenceKey struct {
	first  comboKey
	second comboKey
}

// Registry is the immutable shortcut knowledge base. It is built once and read
// concurrently afterwards; the matcher only performs map lookups per keystroke.
type Registry struct {
	universal     map[comboKey]ActionInfo
	universalSeq  map[sequenceKey]ActionInfo
	categories    map[context.Category]map[comboKey]ActionInfo
	categorySeq   map[context.Category]map[sequenceKey]ActionInfo
	profiles      []*Profile
	titleProfiles []*Profile
	byExe         map[string][]*Profile

	// duplicates records every silent overwrite performed while building, so
	// Audit can report it without panicking in production.
	duplicates []Duplicate
}

func act(title string) ActionInfo {
	return ActionInfo{Title: title}
}

func actC(title, category string) ActionInfo {
	return ActionInfo{Title: title, Category: category}
}

// NewRegistry builds the complete shortcut knowledge base.
func NewRegistry() *Registry {
	r := &Registry{
		universal:    make(map[comboKey]ActionInfo),
		universalSeq: make(map[sequenceKey]ActionInfo),
		categories:   make(map[context.Category]map[comboKey]ActionInfo),
		categorySeq:  make(map[context.Category]map[sequenceKey]ActionInfo),
		byExe:        make(map[string][]*Profile),
	}
	r.initUniversal()
	r.initCategoryDictionaries()
	r.initProfiles()
	r.index()
	return r
}

func (r *Registry) setUniversal(combo comboKey, a ActionInfo) {
	if prev, ok := r.universal[combo]; ok && prev.Title != a.Title {
		r.duplicates = append(r.duplicates, Duplicate{
			Scope: "universal", Combo: combo.Format(), Kept: a.Title, Dropped: prev.Title,
		})
	}
	r.universal[combo] = a
}

func (r *Registry) addU(mods hook.Modifiers, key, title string) {
	r.setUniversal(canonicalCombo(mods, key), act(title))
}

func (r *Registry) addUc(mods hook.Modifiers, key, title, category string) {
	r.setUniversal(canonicalCombo(mods, key), actC(title, category))
}

func (r *Registry) addCat(cat context.Category, mods hook.Modifiers, key, title string) {
	if r.categories[cat] == nil {
		r.categories[cat] = make(map[comboKey]ActionInfo)
	}
	combo := canonicalCombo(mods, key)
	if prev, ok := r.categories[cat][combo]; ok && prev.Title != title {
		r.duplicates = append(r.duplicates, Duplicate{
			Scope: "category:" + string(cat), Combo: combo.Format(), Kept: title, Dropped: prev.Title,
		})
	}
	r.categories[cat][combo] = act(title)
}

func (r *Registry) addCatSeq(cat context.Category, firstMods hook.Modifiers, firstKey string, secondMods hook.Modifiers, secondKey, title string) {
	if r.categorySeq[cat] == nil {
		r.categorySeq[cat] = make(map[sequenceKey]ActionInfo)
	}
	r.categorySeq[cat][sequenceKey{
		first:  canonicalCombo(firstMods, firstKey),
		second: canonicalCombo(secondMods, secondKey),
	}] = act(title)
}

func (r *Registry) addProfile(p *Profile) {
	r.profiles = append(r.profiles, p)
	if len(p.TitleHints) > 0 {
		r.titleProfiles = append(r.titleProfiles, p)
	}
}

func (r *Registry) index() {
	for _, p := range r.profiles {
		for _, exe := range p.Exe {
			norm := context.NormalizeExe(exe)
			if norm == "" {
				continue
			}
			r.byExe[norm] = append(r.byExe[norm], p)
		}
	}
}

// initProfiles registers every application family. Files are split by domain so
// the data stays reviewable: adding a shortcut should never mean scrolling
// through a thousand lines of unrelated mappings.
func (r *Registry) initProfiles() {
	// Editors & IDEs
	r.initBrowsers()
	r.initTerminals()
	r.initVSCodeFamily()
	r.initJetBrains()
	r.initEclipse()
	r.initVisualStudio()
	r.initSublime()
	r.initDevCpp()
	r.initIDLE()
	r.initNotepadPP()
	r.initZed()

	// Office, design, media, chat, notes, shell
	r.initOffice()
	r.initDesign()
	r.initMedia()
	r.initChat()
	r.initNotes()
	r.initExplorer()
	r.initTotalCommander()
}

// Resolve picks the most specific profile for an app context.
func (r *Registry) Resolve(app context.AppContext) (*Profile, context.Category) {
	exe := context.NormalizeExe(app.Exe)
	if candidates := r.byExe[exe]; len(candidates) > 0 {
		// Unambiguous executables are claimed immediately.
		for _, p := range candidates {
			if !p.requiresTitle(exe) {
				return p, p.Category
			}
		}
		// Ambiguous executables (python.exe, …) only match with a title hint;
		// this keeps ordinary Python scripts from being labelled IDLE.
		for _, p := range candidates {
			if p.requiresTitle(exe) && p.matchesTitle(app.Title) {
				return p, p.Category
			}
		}
		// Ambiguous process and no hint matched: stay conservative and fall
		// back to the supplied category (if any).
		return nil, app.Category
	}

	// Unknown executable (e.g. a renamed IDLE build): allow title-hinted
	// profiles to claim it, but never guess from a generic title.
	for _, p := range r.titleProfiles {
		if p.matchesTitle(app.Title) {
			return p, p.Category
		}
	}
	return nil, app.Category
}

func profileMeta(p *Profile) (string, string) {
	if p == nil {
		return "", ""
	}
	return p.ID, p.Name
}

// Lookup resolves one chord through user-independent layers:
// profile -> category -> universal -> none.
func (r *Registry) Lookup(profile *Profile, category context.Category, combo comboKey) (ActionInfo, Confidence, string, string) {
	id, name := profileMeta(profile)
	if profile != nil {
		if a, ok := profile.Shortcuts[combo]; ok {
			return a, ConfidenceProfile, id, name
		}
	}
	if category != "" {
		if a, ok := r.categories[category][combo]; ok {
			return a, ConfidenceCategory, id, name
		}
	}
	if a, ok := r.universal[combo]; ok {
		return a, ConfidenceUniversal, id, name
	}
	return ActionInfo{}, ConfidenceNone, id, name
}

// LookupSequence resolves a two-chord sequence through the same layers.
func (r *Registry) LookupSequence(profile *Profile, category context.Category, seq sequenceKey) (ActionInfo, Confidence, string, string) {
	id, name := profileMeta(profile)
	if profile != nil {
		if a, ok := profile.Sequences[seq]; ok {
			return a, ConfidenceProfile, id, name
		}
	}
	if category != "" {
		if a, ok := r.categorySeq[category][seq]; ok {
			return a, ConfidenceCategory, id, name
		}
	}
	if a, ok := r.universalSeq[seq]; ok {
		return a, ConfidenceUniversal, id, name
	}
	return ActionInfo{}, ConfidenceNone, id, name
}

// AllActionTitles returns every distinct action title in the registry. It is
// used by localization completeness tests and documentation tooling.
func (r *Registry) AllActionTitles() []string {
	seen := make(map[string]struct{}, 512)
	add := func(a ActionInfo) {
		if a.Title != "" {
			seen[a.Title] = struct{}{}
		}
	}
	for _, a := range r.universal {
		add(a)
	}
	for _, a := range r.universalSeq {
		add(a)
	}
	for _, m := range r.categories {
		for _, a := range m {
			add(a)
		}
	}
	for _, m := range r.categorySeq {
		for _, a := range m {
			add(a)
		}
	}
	for _, p := range r.profiles {
		for _, a := range p.Shortcuts {
			add(a)
		}
		for _, a := range p.Sequences {
			add(a)
		}
	}
	out := make([]string, 0, len(seen))
	for title := range seen {
		out = append(out, title)
	}
	sort.Strings(out)
	return out
}

// ProfileByApp finds a profile by executable, id, or category. Used to expand
// user-defined shortcuts consistently with the built-in registry.
func (r *Registry) ProfileByApp(app string) *Profile {
	key := strings.ToLower(strings.TrimSpace(app))
	if key == "" {
		return nil
	}
	for _, p := range r.profiles {
		if strings.EqualFold(p.ID, key) || strings.EqualFold(p.Name, key) {
			return p
		}
	}
	if candidates := r.byExe[context.NormalizeExe(key)]; len(candidates) > 0 {
		return candidates[0]
	}
	return nil
}

// ProfileByID returns a profile by its stable identifier.
func (r *Registry) ProfileByID(id string) *Profile {
	for _, p := range r.profiles {
		if strings.EqualFold(p.ID, id) {
			return p
		}
	}
	return nil
}

// Profiles returns the profiles in registration order.
func (r *Registry) Profiles() []*Profile {
	out := make([]*Profile, len(r.profiles))
	copy(out, r.profiles)
	return out
}
