package hud

import (
	"testing"

	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/theme"
)

func TestThemeInheritanceForVisuals(t *testing.T) {
	cfg := config.DefaultConfig()
	th := theme.ThemeLiquidGlass

	if got := effectiveStyle(th, cfg); got != "glass" {
		t.Errorf("expected Liquid Glass to inherit glass style, got %q", got)
	}
	if got := effectiveFontSize(th, cfg); got != th.KeyCapFontSize {
		t.Errorf("expected theme font size %.1f, got %.1f", th.KeyCapFontSize, got)
	}
	if got := effectiveRadius(th, cfg); got != th.KeyCapRadius {
		t.Errorf("expected theme radius %.1f, got %.1f", th.KeyCapRadius, got)
	}

	// Explicit user choices must still win over the theme.
	cfg.Appearance.Style = "neon"
	cfg.Text.Size = 40
	cfg.Border.Radius = 4
	if got := effectiveStyle(th, cfg); got != "neon" {
		t.Errorf("expected explicit style neon, got %q", got)
	}
	if got := effectiveFontSize(th, cfg); got != 40 {
		t.Errorf("expected explicit font size 40, got %.1f", got)
	}
	if got := effectiveRadius(th, cfg); got != 4 {
		t.Errorf("expected explicit radius 4, got %.1f", got)
	}
}

func TestSlideVectorFollowsScreenAnchor(t *testing.T) {
	cases := []struct {
		align config.Alignment
		x, y  float32
	}{
		{config.AlignBottomCenter, 0, 28},
		{config.AlignTopCenter, 0, -28},
		{config.AlignCenterLeft, -28, 0},
		{config.AlignTopRight, 28, 0},
		{config.AlignCenter, 28, 0},
	}
	for _, tc := range cases {
		x, y := slideVector(tc.align)
		if x != tc.x || y != tc.y {
			t.Errorf("slideVector(%s) = (%.0f, %.0f), want (%.0f, %.0f)", tc.align, x, y, tc.x, tc.y)
		}
	}
}

func TestViewStateDisplayGeometry(t *testing.T) {
	v := NewViewState()

	// Initial default scale
	if got := v.GetScale(); got != 1.0 {
		t.Errorf("expected default scale 1.0, got %f", got)
	}

	// Multi-monitor geometry with negative origin
	v.SetDisplayGeometry(-1920, 0, 1920, 1080, 1.5)
	ox, oy, w, h, s := v.DisplayGeometry()
	if ox != -1920 || oy != 0 || w != 1920 || h != 1080 || s != 1.5 {
		t.Errorf("DisplayGeometry mismatch: got (%d, %d, %d, %d, %f)", ox, oy, w, h, s)
	}
}
