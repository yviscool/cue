package hud

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/engine"
	"github.com/yviscool/cue/internal/theme"
)

// labelContrast is the strongest difference between a glyph and the surface it
// is printed on, measured inside the middle of the cap: the most common value
// there is the cap's own surface, and the glyphs are the pixels that depart
// from it. Measuring against the stage instead would fail a light label on a
// dark modifier cap, which is perfectly readable.
func labelContrast(img *image.RGBA) float32 {
	hist := make([]int, 51)
	lums := make([]float32, 0, 64*64)
	for y := img.Bounds().Dy() * 2 / 5; y < img.Bounds().Dy()*3/5; y++ {
		for x := img.Bounds().Dx() * 2 / 5; x < img.Bounds().Dx()*3/5; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lum := (float32(r>>8)*0.299 + float32(g>>8)*0.587 + float32(b>>8)*0.114) / 255.0
			lums = append(lums, lum)
			hist[int(lum*50+0.5)]++
		}
	}
	surface := 0
	for i, n := range hist {
		if n > hist[surface] {
			surface = i
		}
	}
	body := float32(surface) / 50

	worst := float32(0)
	for _, lum := range lums {
		if d := absFloat32(lum - body); d > worst {
			worst = d
		}
	}
	return worst
}

// TestFragileModelsStayVisible is the regression guard for the two models that
// take their ink from the theme instead of inventing it: glass and outline.
// Both used to assume a dark desktop, so a light theme drew black rules on a
// black preview stage and the caps disappeared entirely. A label that cannot be
// told from the surface it sits on is the bug this catches.
func TestFragileModelsStayVisible(t *testing.T) {
	// A modifier cap is its own colour pair in every theme: dark on a light
	// theme, light on a dark one. That pair is where a model that throws away
	// the cap's own fill breaks first, so both are probed.
	probes := []engine.KeyCap{
		{Label: "A"},
		{Label: "Ctrl", IsModifier: true},
	}
	for _, th := range theme.AllThemes {
		probe := th
		for _, style := range []string{"glass", "outline"} {
			style := style
			for _, key := range probes {
				key := key
				stage := previewStageColor(probe)
				cfg := config.AppConfig{
					ThemePreset: probe.Name,
					Appearance:  config.AppearanceSettings{Style: style},
				}

				tt := ui.NewTester(func(c *ui.Context) {
					ui.Box(c).
						Fill().
						Background(stage).
						AlignItems(ui.Center).
						Justify(ui.Center).
						Children(func() {
							renderKeyCap(c, "probe", key, 1, 0, false, probe, cfg)
						})
				}, 180, 90)
				tt.Frame()

				if got := labelContrast(tt.Image()); got < 0.30 {
					t.Errorf("%s / %s / %s: label barely differs from the cap surface (contrast %.2f)", probe.Name, style, key.Label, got)
				}
			}
		}
	}
}

// TestLabelContrastMetricIsSensitive proves the guard above can fail: a flat
// surface has no label, so it must measure as no contrast at all.
func TestLabelContrastMetricIsSensitive(t *testing.T) {
	flat := image.NewRGBA(image.Rect(0, 0, 180, 90))
	draw.Draw(flat, flat.Bounds(), &image.Uniform{C: color.RGBA{R: 120, G: 120, B: 120, A: 255}}, image.Point{}, draw.Src)
	if got := labelContrast(flat); got > 0.02 {
		t.Errorf("a blank surface measured %.2f contrast, want ~0", got)
	}
}

// TestPreviewStageFollowsTheme pins the reason the preview is honest: a light
// theme gets a light stage, a dark theme a dark one, and neither is pure.
func TestPreviewStageFollowsTheme(t *testing.T) {
	light := previewStageColor(theme.ThemeByName("Swiss Typography 1957"))
	if colorLuminance(light) < 0.7 {
		t.Errorf("light theme stage is dark: %#v", light)
	}
	dark := previewStageColor(theme.ThemeByName("Liquid Glass"))
	if colorLuminance(dark) > 0.3 {
		t.Errorf("dark theme stage is light: %#v", dark)
	}
	for _, th := range theme.AllThemes {
		lum := colorLuminance(previewStageColor(th))
		if lum < 0.02 || lum > 0.98 {
			t.Errorf("%s: stage is a flat extreme (%.2f), caps would lose their edge", th.Name, lum)
		}
	}
}

// TestArrowKeysGetArrowIcons guards a lookup-order bug: the mouse table matched
// the bare words "left" and "right", which are also the arrow key labels, so
// every Left/Right press drew a mouse instead of an arrow. Mouse events use the
// qualified names, so the bare words must fall through to the navigation table.
func TestArrowKeysGetArrowIcons(t *testing.T) {
	cfg := config.DefaultConfig()

	cases := []struct {
		label string
		want  *ui.SVG
	}{
		{"Right", iconKeyArrowRight},
		{"Left", iconKeyArrowLeft},
		{"Up", iconKeyArrowUp},
		{"Down", iconKeyArrowDown},
		{"Right Click", iconKeyMouseRight},
		{"Left Click", iconKeyMouseLeft},
	}

	for _, tc := range cases {
		got := getKeyIcon(tc.label, false, cfg)
		if got != tc.want {
			t.Errorf("%q resolved to the wrong icon", tc.label)
		}
	}
	if getKeyIcon("Right", false, cfg) == getKeyIcon("Right Click", false, cfg) {
		t.Error("the Right arrow and the right mouse button must not share an icon")
	}
}
