package semantic

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yviscool/cue/internal/context"
)

// Duplicate is one silent overwrite: a chord registered twice with two
// different meanings inside the same scope. The later registration wins, which
// is fine at build time but almost always a data bug.
type Duplicate struct {
	Scope   string // "universal", "category:editor", "profile:sublime"
	Combo   string // human-readable chord
	Kept    string // action that survived (registered last)
	Dropped string // action that was silently replaced
}

func (d Duplicate) String() string {
	return fmt.Sprintf("%s: %s kept %q, dropped %q", d.Scope, d.Combo, d.Kept, d.Dropped)
}

// Audit reports every silent overwrite across all scopes. Production callers
// log it; tests assert it is empty.
func (r *Registry) Audit() []Duplicate {
	out := make([]Duplicate, 0, len(r.duplicates))
	out = append(out, r.duplicates...)
	for _, p := range r.profiles {
		out = append(out, p.duplicates...)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Scope != out[j].Scope {
			return out[i].Scope < out[j].Scope
		}
		return out[i].Combo < out[j].Combo
	})
	return out
}

// Stats quantifies the knowledge base. Every number is derived from the live
// maps, never hand-maintained, so documentation cannot drift from the code.
type Stats struct {
	Universal  int
	Categories int
	ProfileMap int
	Sequences  int
	Total      int
	Titles     int
	Profiles   int
	Duplicates int
}

// Stats computes the current coverage counters.
func (r *Registry) Stats() Stats {
	s := Stats{
		Universal: len(r.universal) + len(r.universalSeq),
		Profiles:  len(r.profiles),
		Titles:    len(r.AllActionTitles()),
	}
	for _, m := range r.categories {
		s.Categories += len(m)
	}
	for _, m := range r.categorySeq {
		s.Categories += len(m)
	}
	for _, p := range r.profiles {
		s.ProfileMap += len(p.Shortcuts)
		s.Sequences += len(p.Sequences)
	}
	s.Total = s.Universal + s.Categories + s.ProfileMap + s.Sequences
	s.Duplicates = len(r.Audit())
	return s
}

// Entry is one documented mapping.
type Entry struct {
	Chord    string
	Action   string
	Category string
}

// Section is one documentation table.
type Section struct {
	Heading string
	Note    string
	// ShowGroup renders the "Group" column. Only the universal table carries a
	// per-row group; category and profile sections are already grouped by their
	// heading, so repeating it in every row would be noise.
	ShowGroup bool
	Entries   []Entry
}

func sortedEntries(m map[comboKey]ActionInfo) []Entry {
	out := make([]Entry, 0, len(m))
	for combo, a := range m {
		out = append(out, Entry{Chord: combo.Format(), Action: a.Title, Category: a.Category})
	}
	sortEntries(out)
	return out
}

func sortedSequenceEntries(m map[sequenceKey]ActionInfo) []Entry {
	out := make([]Entry, 0, len(m))
	for seq, a := range m {
		out = append(out, Entry{Chord: seq.Format(), Action: a.Title, Category: a.Category})
	}
	sortEntries(out)
	return out
}

func sortEntries(out []Entry) {
	sort.Slice(out, func(i, j int) bool {
		if out[i].Action != out[j].Action {
			return out[i].Action < out[j].Action
		}
		return out[i].Chord < out[j].Chord
	})
}

// categoryOrder keeps the documentation stable and readable.
var categoryOrder = []context.Category{
	context.CategoryEditor,
	context.CategoryOffice,
	context.CategoryDesign,
	context.CategoryMedia,
	context.CategoryChat,
	context.CategoryNotes,
	context.CategoryBrowser,
	context.CategoryTerminal,
	context.CategorySystem,
}

func categoryRank(c context.Category) int {
	for i, candidate := range categoryOrder {
		if candidate == c {
			return i
		}
	}
	return len(categoryOrder)
}

// Snapshot returns the documentation model: universal dictionary, category
// dictionaries and every profile, with sequences merged into their section.
func (r *Registry) Snapshot() []Section {
	sections := make([]Section, 0, 32)

	sections = append(sections, Section{
		Heading:   "Universal",
		Note:      "Meaning shared by virtually every Windows application. Profiles and categories override these.",
		ShowGroup: true,
		Entries:   append(sortedEntries(r.universal), sortedSequenceEntries(r.universalSeq)...),
	})

	for _, cat := range categoryOrder {
		entries := sortedEntries(r.categories[cat])
		entries = append(entries, sortedSequenceEntries(r.categorySeq[cat])...)
		if len(entries) == 0 {
			continue
		}
		sortEntries(entries)
		sections = append(sections, Section{
			Heading: "Category · " + cat.DisplayName(),
			Note:    "Fallback for applications in the " + string(cat) + " family that have no exact profile.",
			Entries: entries,
		})
	}

	profiles := make([]*Profile, len(r.profiles))
	copy(profiles, r.profiles)
	sort.SliceStable(profiles, func(i, j int) bool {
		ri, rj := categoryRank(profiles[i].Category), categoryRank(profiles[j].Category)
		if ri != rj {
			return ri < rj
		}
		return profiles[i].Name < profiles[j].Name
	})

	for _, p := range profiles {
		entries := sortedEntries(p.Shortcuts)
		entries = append(entries, sortedSequenceEntries(p.Sequences)...)
		if len(entries) == 0 {
			continue
		}
		sortEntries(entries)
		sections = append(sections, Section{
			Heading: "Profile · " + p.Name,
			Note:    profileNote(p, len(entries)),
			Entries: entries,
		})
	}
	return sections
}

func profileNote(p *Profile, count int) string {
	parts := []string{fmt.Sprintf("%d mappings", count)}
	if p.VerifiedAt != "" {
		parts = append(parts, "verified against "+p.VerifiedAt)
	}
	if p.Source != "" {
		parts = append(parts, "source: "+p.Source)
	}
	return strings.Join(parts, " · ")
}

// RenderReference renders the full shortcut reference as Markdown. The output
// is deterministic so a test can assert the checked-in file matches the code.
func (r *Registry) RenderReference() string {
	s := r.Stats()
	var b strings.Builder

	b.WriteString("# Shortcut Reference\n\n")
	b.WriteString("Complete, generated inventory of every shortcut Cue can interpret. ")
	b.WriteString("This file is produced from the semantic registry itself, so it can never\n")
	b.WriteString("describe shortcuts the app does not know about.\n\n")
	b.WriteString("> Regenerate with `CUE_UPDATE_REFERENCE=1 go test ./internal/semantic -run TestShortcutReference`.\n")
	b.WriteString("> Do not edit by hand; `TestShortcutReference` fails when this file drifts from the code.\n\n")

	b.WriteString("## Coverage\n\n")
	b.WriteString("| Layer | Mappings |\n| :--- | ---: |\n")
	fmt.Fprintf(&b, "| Universal | %d |\n", s.Universal)
	fmt.Fprintf(&b, "| Category dictionaries | %d |\n", s.Categories)
	fmt.Fprintf(&b, "| Application profiles | %d |\n", s.ProfileMap)
	fmt.Fprintf(&b, "| Two-chord sequences | %d |\n", s.Sequences)
	fmt.Fprintf(&b, "| **Total** | **%d** |\n", s.Total)
	fmt.Fprintf(&b, "| Distinct action titles | %d |\n", s.Titles)
	fmt.Fprintf(&b, "| Profiles | %d |\n\n", s.Profiles)

	b.WriteString("Priority when resolving a keystroke:\n")
	b.WriteString("`user override > exact app profile > app category > universal > no label`.\n\n")

	for _, section := range r.Snapshot() {
		fmt.Fprintf(&b, "## %s\n\n", section.Heading)
		if section.Note != "" {
			b.WriteString(section.Note + "\n\n")
		}
		if section.ShowGroup {
			b.WriteString("| Shortcut | Action | Group |\n| :--- | :--- | :--- |\n")
			for _, e := range section.Entries {
				fmt.Fprintf(&b, "| `%s` | %s | %s |\n", e.Chord, e.Action, e.Category)
			}
		} else {
			b.WriteString("| Shortcut | Action |\n| :--- | :--- |\n")
			for _, e := range section.Entries {
				fmt.Fprintf(&b, "| `%s` | %s |\n", e.Chord, e.Action)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
