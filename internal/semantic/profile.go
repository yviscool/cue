package semantic

import (
	"strings"

	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// Profile is an application family with an exact executable match.
//
// Source and VerifiedAt exist so every mapping can be traced back to an
// authoritative keymap instead of a blog post: a wrong chord is worse than a
// missing one, because it teaches the user the wrong thing.
type Profile struct {
	ID       string
	Name     string
	Category context.Category

	// Source is the authoritative reference the mappings were transcribed from
	// (official keymap file or documentation URL).
	Source string
	// VerifiedAt names the exact build the keymap was checked against, e.g.
	// "Sublime Text 4 (Windows)".
	VerifiedAt string

	Exe        []string
	TitleHints []string
	// Ambiguous lists executables that are shared by several applications and
	// therefore need a matching title hint before this profile may claim them.
	// Example: python.exe is used by IDLE and by ordinary scripts.
	Ambiguous []string

	Shortcuts map[comboKey]ActionInfo
	Sequences map[sequenceKey]ActionInfo

	duplicates []Duplicate
}

func newProfile(id, name string, category context.Category, exe ...string) *Profile {
	return &Profile{
		ID:        id,
		Name:      name,
		Category:  category,
		Exe:       exe,
		Shortcuts: make(map[comboKey]ActionInfo),
		Sequences: make(map[sequenceKey]ActionInfo),
	}
}

func (p *Profile) titleHints(hints ...string) *Profile {
	p.TitleHints = append(p.TitleHints, hints...)
	return p
}

func (p *Profile) ambiguous(exe ...string) *Profile {
	p.Ambiguous = append(p.Ambiguous, exe...)
	return p
}

// source records provenance for the whole profile.
func (p *Profile) source(reference, verifiedAt string) *Profile {
	p.Source = reference
	p.VerifiedAt = verifiedAt
	return p
}

// set installs one shortcut and records an overwrite when it silently replaces
// a different action. Overwrites are legal at build time but must never ship:
// TestRegistryHasNoDuplicates asserts the list stays empty.
func (p *Profile) set(combo comboKey, a ActionInfo) {
	if prev, ok := p.Shortcuts[combo]; ok && prev.Title != a.Title {
		p.duplicates = append(p.duplicates, Duplicate{
			Scope:   "profile:" + p.ID,
			Combo:   combo.Format(),
			Kept:    a.Title,
			Dropped: prev.Title,
		})
	}
	p.Shortcuts[combo] = a
}

// k registers a shortcut with an inferred category.
func (p *Profile) k(mods hook.Modifiers, key, title string) *Profile {
	p.set(canonicalCombo(mods, key), act(title))
	return p
}

// kc registers a shortcut with an explicit category.
func (p *Profile) kc(mods hook.Modifiers, key, title, category string) *Profile {
	p.set(canonicalCombo(mods, key), actC(title, category))
	return p
}

// seq registers a two-chord sequence such as Ctrl+K Ctrl+B.
func (p *Profile) seq(firstMods hook.Modifiers, firstKey string, secondMods hook.Modifiers, secondKey, title string) *Profile {
	key := sequenceKey{
		first:  canonicalCombo(firstMods, firstKey),
		second: canonicalCombo(secondMods, secondKey),
	}
	if prev, ok := p.Sequences[key]; ok && prev.Title != title {
		p.duplicates = append(p.duplicates, Duplicate{
			Scope:   "profile:" + p.ID,
			Combo:   key.Format(),
			Kept:    title,
			Dropped: prev.Title,
		})
	}
	p.Sequences[key] = act(title)
	return p
}

// requiresTitle reports whether this executable can only be claimed when one of
// the profile's title hints matches.
func (p *Profile) requiresTitle(exe string) bool {
	if len(p.Ambiguous) == 0 {
		return false
	}
	norm := context.NormalizeExe(exe)
	for _, candidate := range p.Ambiguous {
		if context.NormalizeExe(candidate) == norm {
			return true
		}
	}
	return false
}

func (p *Profile) matchesTitle(title string) bool {
	if len(p.TitleHints) == 0 {
		return false
	}
	t := strings.ToLower(title)
	for _, hint := range p.TitleHints {
		if strings.Contains(t, strings.ToLower(hint)) {
			return true
		}
	}
	return false
}

// ShortcutCount reports how many mappings the profile carries, sequences
// included.
func (p *Profile) ShortcutCount() int {
	return len(p.Shortcuts) + len(p.Sequences)
}
