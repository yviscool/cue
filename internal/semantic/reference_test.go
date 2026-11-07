package semantic_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yviscool/cue/internal/semantic"
)

// referencePath is the generated inventory committed to the repository.
var referencePath = filepath.Join("..", "..", "docs", "shortcut-reference.md")

// TestShortcutReference keeps docs/shortcut-reference.md in lockstep with the
// registry. Regenerate it with:
//
//	CUE_UPDATE_REFERENCE=1 go test ./internal/semantic -run TestShortcutReference
func TestShortcutReference(t *testing.T) {
	registry := semantic.NewRegistry()
	want := registry.RenderReference()

	if os.Getenv("CUE_UPDATE_REFERENCE") == "1" {
		if err := os.WriteFile(referencePath, []byte(want), 0o644); err != nil {
			t.Fatalf("write reference: %v", err)
		}
		t.Logf("wrote %s (%d bytes)", referencePath, len(want))
		return
	}

	got, err := os.ReadFile(referencePath)
	if err != nil {
		t.Fatalf("read %s: %v\nRegenerate with CUE_UPDATE_REFERENCE=1 go test ./internal/semantic -run TestShortcutReference", referencePath, err)
	}
	if string(got) != want {
		t.Fatalf("docs/shortcut-reference.md is stale (%d bytes on disk, %d rendered).\n"+
			"Regenerate with: CUE_UPDATE_REFERENCE=1 go test ./internal/semantic -run TestShortcutReference",
			len(got), len(want))
	}
}

// TestReferenceCoversEveryProfile asserts the rendered document actually lists
// every profile that carries shortcuts — a doc that silently drops a section is
// as misleading as no doc at all.
func TestReferenceCoversEveryProfile(t *testing.T) {
	registry := semantic.NewRegistry()
	doc := registry.RenderReference()

	for _, p := range registry.Profiles() {
		if p.ShortcutCount() == 0 {
			continue
		}
		if !strings.Contains(doc, "## Profile · "+p.Name) {
			t.Errorf("reference is missing the %q profile section", p.Name)
		}
	}
	if !strings.Contains(doc, "## Universal") {
		t.Error("reference is missing the universal section")
	}
}
