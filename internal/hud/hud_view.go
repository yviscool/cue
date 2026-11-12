package hud

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/engine"
	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/theme"
)

// ViewState holds state needed by the HUD view rendering loop.
type ViewState struct {
	mu           sync.RWMutex
	currentTheme theme.ThemeConfig
	cfg          config.AppConfig
	state        engine.UIState
	lastRepaint  time.Time
	scale        float64
	originX      int
	originY      int
	width        int
	height       int
}

// NewViewState creates a new HUD view state holder.
func NewViewState() *ViewState {
	return &ViewState{
		currentTheme: theme.ThemeLiquidGlass,
		cfg:          config.DefaultConfig(),
		state: engine.UIState{
			FilterMode:  engine.FilterSmart,
			SoundActive: true,
			ThemeName:   theme.ThemeLiquidGlass.Name,
		},
		lastRepaint: time.Now(),
		scale:       1.0,
	}
}

// SetScale sets the monitor DPI scale factor.
func (v *ViewState) SetScale(s float64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if s <= 0 {
		s = 1.0
	}
	v.scale = s
}

// GetScale returns the DPI scale factor.
func (v *ViewState) GetScale() float64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.scale
}

// SetDisplayGeometry sets the active display origin, dimensions, and DPI scale factor.
func (v *ViewState) SetDisplayGeometry(originX, originY, width, height int, scale float64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if scale <= 0 {
		scale = 1.0
	}
	v.originX = originX
	v.originY = originY
	v.width = width
	v.height = height
	v.scale = scale
}

// DisplayGeometry returns the active display origin, dimensions, and DPI scale factor.
func (v *ViewState) DisplayGeometry() (originX, originY, width, height int, scale float64) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.originX, v.originY, v.width, v.height, v.scale
}

// SetTheme updates the current theme.
func (v *ViewState) SetTheme(t theme.ThemeConfig) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.currentTheme = t
	v.state.ThemeName = t.Name
}

// GetTheme returns the active theme.
func (v *ViewState) GetTheme() theme.ThemeConfig {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.currentTheme
}

// SetConfig updates the live configuration driving visual rendering.
func (v *ViewState) SetConfig(c config.AppConfig) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.cfg = c
}

// GetConfig returns a copy of active configuration.
func (v *ViewState) GetConfig() config.AppConfig {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.cfg
}

// UpdateState updates the current engine state snapshot.
func (v *ViewState) UpdateState(s engine.UIState) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.state = s
	v.lastRepaint = time.Now()
}

// Snapshot returns the current UI state.
func (v *ViewState) Snapshot() engine.UIState {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.state
}

// CalcHUDPosition computes the (x, y) coordinates for the HUD window based on display work area,
// alignment, margins, and window dimensions.
func CalcHUDPosition(disp mygo.Display, winW, winH int, al config.Alignment, mx, my int) (int, int) {
	wa := disp.WorkArea
	var x, y int

	switch al {
	case config.AlignTopLeft:
		x = wa.X + mx
		y = wa.Y + my
	case config.AlignTopCenter:
		x = wa.X + (wa.Width-winW)/2
		y = wa.Y + my
	case config.AlignTopRight:
		x = wa.X + wa.Width - winW - mx
		y = wa.Y + my
	case config.AlignCenterLeft:
		x = wa.X + mx
		y = wa.Y + (wa.Height-winH)/2
	case config.AlignCenter:
		x = wa.X + (wa.Width-winW)/2
		y = wa.Y + (wa.Height-winH)/2
	case config.AlignCenterRight:
		x = wa.X + wa.Width - winW - mx
		y = wa.Y + (wa.Height-winH)/2
	case config.AlignBottomLeft:
		x = wa.X + mx
		y = wa.Y + wa.Height - winH - my
	case config.AlignBottomCenter:
		x = wa.X + (wa.Width-winW)/2
		y = wa.Y + wa.Height - winH - my
	case config.AlignBottomRight:
		x = wa.X + wa.Width - winW - mx
		y = wa.Y + wa.Height - winH - my
	default:
		x = wa.X + (wa.Width-winW)/2
		y = wa.Y + wa.Height - winH - my
	}

	// Clamp within work area bounds
	if x < wa.X {
		x = wa.X
	}
	if x+winW > wa.X+wa.Width {
		x = wa.X + wa.Width - winW
	}
	if y < wa.Y {
		y = wa.Y
	}
	if y+winH > wa.Y+wa.Height {
		y = wa.Y + wa.Height - winH
	}

	return x, y
}

// effectiveStyle resolves the physical keycap model. Empty means "use the
// signature style of the active theme", which is what makes a newly selected
// preset look like its preview instead of inheriting a stale low-profile look.
func effectiveStyle(t theme.ThemeConfig, cfg config.AppConfig) string {
	if cfg.Appearance.Style != "" {
		return cfg.Appearance.Style
	}
	if t.DefaultStyle != "" {
		return t.DefaultStyle
	}
	return "lowprofile"
}

// effectiveFontSize resolves typography: an explicit user size wins, otherwise
// the active theme supplies its own cap height.
func effectiveFontSize(t theme.ThemeConfig, cfg config.AppConfig) float32 {
	if cfg.Text.Size > 0 {
		return float32(cfg.Text.Size)
	}
	if t.KeyCapFontSize > 0 {
		return t.KeyCapFontSize
	}
	return 20
}

// effectiveRadius resolves the keycap corner radius with theme inheritance.
func effectiveRadius(t theme.ThemeConfig, cfg config.AppConfig) float32 {
	if cfg.Border.Radius > 0 {
		return float32(cfg.Border.Radius)
	}
	if t.KeyCapRadius > 0 {
		return t.KeyCapRadius
	}
	return 8
}

// slideVector points toward the screen edge nearest the HUD anchor. Enter comes
// from that edge and Exit returns to it, so the motion explains where the cue
// lives rather than decorating it with an arbitrary direction.
func slideVector(a config.Alignment) (float32, float32) {
	const d float32 = 32
	switch a {
	case config.AlignTopLeft, config.AlignCenterLeft, config.AlignBottomLeft:
		return -d, 0
	case config.AlignTopRight, config.AlignCenterRight, config.AlignBottomRight:
		return d, 0
	case config.AlignTopCenter:
		return 0, -d
	case config.AlignBottomCenter:
		return 0, d
	default: // center / center-center: no nearer edge, use a gentle horizontal reveal
		return d, 0
	}
}

// Render builds the native UI tree for the HUD overlay window.
func (v *ViewState) Render(c *ui.Context) {
	v.mu.RLock()
	t := v.currentTheme
	snap := v.state
	cfg := v.cfg
	scale := v.scale
	originX := v.originX
	originY := v.originY
	v.mu.RUnlock()

	// Ensure whole canvas root is 100% transparent
	c.Root().Background(ui.Transparent)

	hasMouse := snap.MouseState.Active || !snap.MouseState.Timestamp.IsZero() || snap.MouseState.Wheel != 0
	shouldRenderMouse := cfg.Mouse.KeepHighlight || cfg.Mouse.KeepIndicator || (hasMouse && (cfg.Mouse.ShowIndicator || cfg.Mouse.ShowClicks))

	// 1. Mouse Overlay (Cursor Ring & Mouse Indicator at cursor coordinates)
	if shouldRenderMouse {
		RenderMouseOverlay(c, snap.MouseState, cfg, t, scale, originX, originY)
	}

	// 2. Key Overlay (Chords container aligned to 9-grid screen edge with Margins)
	// We ALWAYS keep this container mounted in the MyGo scene graph so that:
	// 1) The element's parent is already born before a chord appears, triggering MyGo Enter transitions.
	// 2) The parent remains present while expired chords animate their Exit ghost transitions.
	var hAlign, vAlign ui.Align

	switch cfg.Appearance.Alignment {
	case config.AlignTopLeft:
		hAlign, vAlign = ui.Start, ui.Start
	case config.AlignTopCenter:
		hAlign, vAlign = ui.Center, ui.Start
	case config.AlignTopRight:
		hAlign, vAlign = ui.End, ui.Start
	case config.AlignCenterLeft:
		hAlign, vAlign = ui.Start, ui.Center
	case config.AlignCenter:
		hAlign, vAlign = ui.Center, ui.Center
	case config.AlignCenterRight:
		hAlign, vAlign = ui.End, ui.Center
	case config.AlignBottomLeft:
		hAlign, vAlign = ui.Start, ui.End
	case config.AlignBottomCenter:
		hAlign, vAlign = ui.Center, ui.End
	case config.AlignBottomRight:
		hAlign, vAlign = ui.End, ui.End
	default:
		hAlign, vAlign = ui.Center, ui.End
	}

	renderContent := func() {
		for _, chord := range snap.Chords {
			renderChord(c, chord, t, cfg)
		}
	}

	mx := cfg.Appearance.MarginX
	my := cfg.Appearance.MarginY
	if mx < 0 {
		mx = 0
	}
	if my < 0 {
		my = 0
	}

	fmx := float32(mx)
	fmy := float32(my)

	// The whole key overlay shares one opacity multiplier, so the user can fade
	// the keycaps into the presentation without touching individual colors. It
	// multiplies down the tree, composing with each chord's own enter/exit fade.
	overlayOpacity := float32(cfg.Appearance.EffectiveOpacity())

	if cfg.Appearance.FlexDirection == "col" || cfg.Appearance.FlexDirection == "column" {
		ui.Column(c.Key("key-overlay-container")).
			Fill().
			Justify(vAlign).
			AlignItems(hAlign).
			Padding(fmy, fmx, fmy, fmx).
			Gap(12).
			Opacity(overlayOpacity).
			Children(renderContent)
	} else {
		ui.Row(c.Key("key-overlay-container")).
			Fill().
			Justify(hAlign).
			AlignItems(vAlign).
			Padding(fmy, fmx, fmy, fmx).
			Gap(14).
			Opacity(overlayOpacity).
			Children(renderContent)
	}
}

// renderChord renders an individual key combo cluster with dynamic motion and subtitles.
func renderChord(c *ui.Context, chord engine.DisplayChord, t theme.ThemeConfig, cfg config.AppConfig) {
	animDuration := time.Duration(cfg.Appearance.AnimationDuration*1000) * time.Millisecond
	if animDuration <= 0 {
		animDuration = 200 * time.Millisecond
	}
	if animDuration > 2000*time.Millisecond {
		animDuration = 2000 * time.Millisecond
	}

	// "none" must be truly instant. Leaving ElementTransition at Duration 0
	// would ask MyGo for its 200 ms default, which is the opposite of the
	// promise the label makes, so no transition is attached at all.
	var trans *ui.ElementTransition
	switch cfg.Appearance.Animation {
	case "none":
		trans = nil
	case "fade":
		dur := animDuration
		if chord.Velocity == 2 {
			dur = animDuration * 65 / 100
		} else if chord.Velocity == 0 {
			dur = animDuration * 135 / 100
		}
		trans = &ui.ElementTransition{
			Duration: dur,
			Ease:     ui.EaseOut,
			Enter:    &ui.Motion{Opacity: 0},
			Exit:     &ui.Motion{Opacity: 0},
		}
	case "zoom":
		dur := animDuration
		bounce := 0.18
		if chord.Velocity == 2 {
			dur = animDuration * 65 / 100
			bounce = 0.24
		} else if chord.Velocity == 0 {
			dur = animDuration * 135 / 100
			bounce = 0.12
		}
		if chord.IsFinisher {
			dur = animDuration * 110 / 100
			bounce = 0.22
		}
		springAnim := ui.Animation{
			Spring: ui.SpringOptions{
				Duration: dur,
				Bounce:   bounce,
			},
		}
		trans = &ui.ElementTransition{
			Animation: &springAnim,
			Enter:     &ui.Motion{Opacity: 0, Scale: 0.68},
			Exit:      &ui.Motion{Opacity: 0, Scale: 0.85},
		}
	case "float":
		dur := animDuration
		if chord.Velocity == 2 {
			dur = animDuration * 65 / 100
		} else if chord.Velocity == 0 {
			dur = animDuration * 135 / 100
		}
		trans = &ui.ElementTransition{
			Duration: dur,
			Ease:     ui.EaseOut,
			Enter:    &ui.Motion{Y: 18, Opacity: 0},
			Exit:     &ui.Motion{Y: -12, Opacity: 0},
		}
	case "slide":
		// Slide from, and return to, the nearest screen edge with clear 32px displacement.
		dx, dy := slideVector(cfg.Appearance.Alignment)
		dur := animDuration
		if chord.Velocity == 2 {
			dur = animDuration * 65 / 100
		} else if chord.Velocity == 0 {
			dur = animDuration * 135 / 100
		}
		trans = &ui.ElementTransition{
			Duration: dur,
			Ease:     ui.EaseOut,
			Enter:    &ui.Motion{X: dx, Y: dy, Opacity: 0},
			Exit:     &ui.Motion{X: dx, Y: dy, Opacity: 0},
		}
	default: // "spring" or default
		dur := animDuration
		if t.EffectiveDurationScale() > 0 {
			dur = time.Duration(float64(dur) * t.EffectiveDurationScale())
		}
		// Natural elastic spring pop: lively, organic, and tactile
		bounce := t.EffectiveSpringBounce()
		if bounce < 0.20 {
			bounce = 0.24
		}
		enterY := float32(8)
		exitY := float32(-6)
		enterScale := float32(0.88)

		if chord.Velocity == 2 {
			dur = dur * 65 / 100
			bounce = bounce * 1.15
			enterY = 5
			exitY = -3
			enterScale = 0.90
		} else if chord.Velocity == 0 {
			dur = dur * 135 / 100
			bounce = bounce * 0.80
			enterY = 10
			exitY = -8
			enterScale = 0.92
		}

		if chord.IsFinisher {
			dur = dur * 110 / 100
			bounce = bounce * 1.15
			enterScale = 1.05
			enterY = 4
		}

		springAnim := ui.Animation{
			Spring: ui.SpringOptions{
				Duration: dur,
				Bounce:   bounce,
			},
		}
		trans = &ui.ElementTransition{
			Animation: &springAnim,
			Enter:     &ui.Motion{Y: enterY, Opacity: 0, Scale: enterScale},
			Exit:      &ui.Motion{Y: exitY, Opacity: 0, Scale: 0.95},
		}
	}

	// Clean, natural physical keycap gap preventing crowded or colliding keys
	capsuleGap := float32(10)
	switch chord.Velocity {
	case 2: // Allegro: tightly drawn together in unified momentum
		capsuleGap = 8
	case 1: // Andante: standard balanced gap
		capsuleGap = 10
	default: // Adagio: spacious, contemplative gap
		capsuleGap = 10
	}

	chordColumn := ui.Column(c).
		Key(chord.ID).
		Center().
		Gap(6)
	if cfg.Background.Enabled {
		chordColumn = chordColumn.
			Padding(12, 16).
			Radius(14).
			Background(parseHexSafe(cfg.Background.Color, ui.RGBA(0, 0, 0, 0.35)))
	}
	if trans != nil {
		chordColumn = chordColumn.Transition(*trans)
	}
	chordColumn.Children(func() {
		// Row containing the keycaps with dynamic cadence gap
		ui.Row(c).Center().Gap(capsuleGap).AlignItems(ui.Center).Children(func() {
			if chord.IsPrivate {
				renderPrivacyShield(c, t)
				return
			}

			for i, key := range chord.Keys {
				streak := 1
				if i == len(chord.Keys)-1 && cfg.Layout.ShowPressCount {
					// A wheel gesture carries its distance instead of a press count:
					// ten notches is one motion, and the notch total is the only
					// number that describes it. Round it so a trackpad's fractional
					// deltas still read as whole notches.
					if notches := int(math.Round(float64(chord.ScrollAmount))); notches > 1 {
						streak = notches
					} else if chord.Streak > 1 {
						streak = chord.Streak
					}
				}
				renderKeyCap(c, fmt.Sprintf("%s-key-%d", chord.ID, i), key, streak, chord.Velocity, chord.IsFinisher, t, cfg)
			}
		})

		// If matched with semantic action subtitle and SemanticSub enabled, render intent pill
		if chord.Action.Title != "" && !chord.IsPrivate && cfg.SemanticSub {
			renderActionSubtitle(c, chord, t, cfg)
		}
	})
}

// colorLuminance is the perceived brightness of a colour, 0 for black and 1
// for white. Alpha is ignored: it answers "which way does this colour lean",
// which is what every contrast decision here needs.
func colorLuminance(c ui.Color) float32 {
	return (float32(c.R)*0.299 + float32(c.G)*0.587 + float32(c.B)*0.114) / 255.0
}

func absFloat32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

// renderKeyCap renders a single tactile keycap driven by live configuration styles and physical models.
func renderKeyCap(c *ui.Context, keyId string, key engine.KeyCap, streak int, velocity int, isFinisher bool, t theme.ThemeConfig, cfg config.AppConfig) {
	bg := t.KeyCapBackground
	borderColor := t.KeyCapBorder
	textColor := t.KeyCapTextColor
	fontSize := effectiveFontSize(t, cfg)
	radius := effectiveRadius(t, cfg)
	style := effectiveStyle(t, cfg)

	// Apply custom primary/secondary color schemes if set
	if cfg.Color.Color != "" {
		bg = parseHexSafe(cfg.Color.Color, bg)
	}
	if cfg.Color.SecondaryColor != "" {
		borderColor = parseHexSafe(cfg.Color.SecondaryColor, borderColor)
	}
	if cfg.Border.Color != "" {
		borderColor = parseHexSafe(cfg.Border.Color, borderColor)
	}
	if cfg.Text.Color != "" {
		textColor = parseHexSafe(cfg.Text.Color, textColor)
	}

	if key.IsModifier {
		bg = t.ModCapBackground
		borderColor = t.ModCapBorder
		textColor = t.ModCapTextColor

		if cfg.Modifier.Highlight && cfg.Modifier.Color != "" {
			bg = parseHexSafe(cfg.Modifier.Color, bg)
			if cfg.Modifier.BorderColor != "" {
				borderColor = parseHexSafe(cfg.Modifier.BorderColor, borderColor)
			}
			if cfg.Modifier.TextColor != "" {
				textColor = parseHexSafe(cfg.Modifier.TextColor, textColor)
			}
		}
	}

	// Tactile hold physics: dynamic key travel descent and press highlight FX.
	// An explicit model owns its travel and held response, so choosing PBT on a
	// pulse theme still sinks like a mechanical switch.
	travelDepth := t.EffectiveTravelDepthFor(cfg.Appearance.Style)
	heldPhysics := t.EffectiveHeldPhysicsFor(cfg.Appearance.Style)
	yOffset := float32(0)
	pressFxEnabled := cfg.Appearance.Animation != "none" && !cfg.Appearance.DisablePressAnimation
	if key.IsHeld && pressFxEnabled {
		bg = t.KeyCapHeldBackground
		borderColor = t.KeyCapHeldBorder

		// Automatic contrast guard: guarantee crisp readability in held state
		heldLum := (float32(bg.R)*0.299 + float32(bg.G)*0.587 + float32(bg.B)*0.114) / 255.0
		textLum := (float32(textColor.R)*0.299 + float32(textColor.G)*0.587 + float32(textColor.B)*0.114) / 255.0
		if heldLum < 0.45 && textLum < 0.45 {
			textColor = ui.Hex("#f8fafc")
		} else if heldLum > 0.65 && textLum > 0.65 {
			textColor = ui.Hex("#0f172a")
		}

		// Pulse physics mode keeps physical position zero; mechanical descent uses restrained travel depth.
		if heldPhysics != theme.HeldPulse {
			yOffset = travelDepth * 0.70
		}
	}

	// Border width: 0 if disabled
	borderWidth := float32(1.5)
	if !cfg.Border.Enabled {
		borderWidth = 0
		borderColor = ui.Transparent
	} else if cfg.Border.Width > 0 {
		borderWidth = float32(cfg.Border.Width)
	}

	// Physical Style Variations (Cue lowprofile, minimal, laptop, pbt, glass, neon, retro, outline)
	// Pure geometry and vertical drop physics: shadowSpread is kept strictly at 0 to prevent
	// shadows from bleeding horizontally and fusing adjacent keys together.
	padX := float32(14)
	padY := float32(8)
	shadowY := float32(3)
	shadowSpread := float32(0)
	shadowBlur := float32(6)
	shadowColor := borderColor.Alpha(0.25)
	// Silhouette scaling: a model has to remain recognisable across the entire
	// keyboard, not only on keys with long labels, so the minimum cap box is
	// part of the physical signature instead of a fixed 44x38 for every build.
	styleW := float32(1)
	styleH := float32(1)
	insetGlow := ui.Transparent
	insetBlur := float32(7)

	switch style {
	case "minimal":
		padX = 14
		padY = 6
		borderWidth = 0
		borderColor = ui.Transparent
		shadowY = 0
		shadowBlur = 0
		shadowSpread = 0
		shadowColor = ui.Transparent
		radius = radius * 0.5
		styleH = 0.9
	case "laptop":
		padX = 13
		padY = 4
		borderWidth = 1.0
		borderColor = borderColor.Alpha(0.4)
		shadowY = 2
		shadowBlur = 3
		shadowSpread = 0
		shadowColor = ui.RGBA(0, 0, 0, 0.22)
		radius = radius * 0.45
		styleW = 1.02
		styleH = 0.78
	case "pbt":
		padX = 16
		padY = 11
		borderWidth = borderWidth + 1.0
		shadowY = 5
		shadowBlur = 7
		shadowSpread = 0
		shadowColor = ui.RGBA(0, 0, 0, 0.38)
		radius = radius * 1.25
		styleW = 1.06
		styleH = 1.2
		insetGlow = ui.RGBA(255, 255, 255, 0.12)
	case "glass":
		// Frosted glass is a surface, not a hole. Three things make it read as
		// glass, and none of them is "make it transparent": a fill that keeps
		// enough of the theme's tint for the label to keep its contrast, an
		// edge drawn from the label colour so it catches light on a dark
		// backdrop and reads as a rule on a light one, and a real drop shadow
		// that lifts the cap off the surface. Thirty percent of white with a
		// white glow around it was none of those: on a dark stage it became a
		// grey tile with dark text, which is the one thing glass never is.
		padX = 15
		padY = 9
		if bg.A > 200 {
			// An opaque theme has no glass in it yet: give it a body you can
			// still see through, and let the edge do the rest.
			bg = bg.Alpha(0.55)
		}
		borderWidth = 1.2
		borderColor = textColor.Alpha(0.26)
		shadowY = 5
		shadowBlur = 18
		shadowSpread = 0
		shadowColor = ui.RGBA(0, 0, 0, 0.30)
		radius = radius * 1.25
		// A crisp inner top edge, not a blurred wash: this is the highlight
		// that makes a flat rectangle read as a lit pane.
		insetGlow = ui.RGBA(255, 255, 255, 0.45)
		insetBlur = 2
		if key.IsHeld {
			shadowBlur = 22
			shadowColor = ui.RGBA(0, 0, 0, 0.36)
		}
	case "neon":
		padX = 14
		padY = 8
		borderWidth = 2.2
		shadowY = 0
		shadowBlur = 18
		shadowSpread = 1
		shadowColor = borderColor.Alpha(0.70)
		insetGlow = borderColor.Alpha(0.40)
		if key.IsHeld {
			shadowBlur = 22
			shadowColor = borderColor.Alpha(0.90)
		}
	case "retro":
		padX = 15
		padY = 12
		borderWidth = 1.6
		shadowY = 5
		shadowBlur = 0
		shadowSpread = 0 // Zero spread: hard bevel drops strictly downwards, zero horizontal spill
		shadowColor = borderColor.Alpha(0.80)
		radius = radius * 0.25
		styleW = 1.04
		styleH = 1.24
		insetGlow = ui.RGBA(255, 255, 255, 0.14)
	case "outline":
		// A wireframe needs a stroke that exists. Some themes carry a black
		// rule, others a translucent white one, and some draw their border in
		// the same value as the cap itself, which leaves no line at all. When
		// the border cannot draw one, the label ink does: it is the only colour
		// in the theme guaranteed to be readable on the cap.
		stroke := borderColor
		if stroke.A < 60 || absFloat32(colorLuminance(stroke)-colorLuminance(bg)) < 0.25 {
			stroke = textColor
		}
		padX = 15
		padY = 9
		// The body stays, faint. A wireframe that throws away the cap's own
		// fill also throws away the contrast its label was chosen for, which is
		// how a dark modifier cap ended up as white text on white.
		if bg.A > 200 {
			bg = bg.Alpha(0.55)
		}
		borderColor = stroke.Alpha(0.88)
		borderWidth = 1.6
		shadowY = 0
		shadowBlur = 0
		shadowSpread = 0
		shadowColor = ui.Transparent
		radius = radius * 0.55
		styleH = 0.98
	default: // "lowprofile"
		shadowY = 2
		shadowBlur = 5
		shadowSpread = 0
		radius = radius * 0.9
		styleH = 0.97
	}

	// Velocity Dynamics: Supersonic Radiant Glow (Allegro) & Finisher Settle Accent
	if pressFxEnabled {
		if velocity == 2 {
			shadowY = 0
			shadowBlur = 12
			shadowSpread = 1
			glowColor := borderColor.Alpha(0.55)
			if style == "neon" || (borderColor.R < 30 && borderColor.G > 200) {
				glowColor = ui.RGBA(0, 240, 255, 0.65) // Radiant cyan neon aura
			}
			shadowColor = glowColor
			borderWidth = borderWidth + 0.3
		} else if isFinisher {
			shadowBlur = 10
			shadowSpread = 1
			shadowColor = borderColor.Alpha(0.35)
			borderWidth = borderWidth + 0.3
		}
	}

	// Text label formatting taking variant, layout symbols, and caps into account
	label := formatKeycapLabel(key.Label, key.IsModifier, cfg)

	icon := getKeyIcon(key.Label, key.IsModifier, cfg)
	isMouse := isMouseKey(key.Label)
	isArrow := isArrowKey(key.Label)
	isPureIcon := (cfg.Text.Variant == "icon" && icon != nil) || (isArrow && icon != nil && cfg.Text.Variant != "text")
	shouldShowIcon := icon != nil && (isPureIcon || (cfg.Layout.ShowIcon && (key.IsModifier || isArrow || isMouse)) || isMouse)

	// Mouse event accent: mouse keycaps glow with the active mouse theme accent color (e.g. vibrant green)
	mth := theme.MouseThemeByID(cfg.Mouse.ThemePreset)
	mouseAccent := ui.Hex(mth.LeftColor)
	if cfg.Mouse.Color != "" {
		mouseAccent = parseHexSafe(cfg.Mouse.Color, mouseAccent)
	}
	if isMouse && borderColor.A > 0 {
		borderColor = mouseAccent.Alpha(0.65)
	}

	// Standard keycap sizing: guarantee natural aspect ratio so 1u alphanumeric keys
	// never collapse into skinny or squashed rectangles.
	minW := float32(44)
	minH := float32(38)
	if isPureIcon {
		minW = 38
		minH = 38
	} else if key.IsModifier {
		minW = 54
	}

	// The model shapes the silhouette, not only its border. Without this the
	// 44x38 floor would flatten laptop, PBT, and retro into the same rectangle.
	minW *= styleW
	minH *= styleH

	justify := ui.Center
	switch cfg.Text.Alignment {
	case config.AlignCenterLeft:
		justify = ui.Start
	case config.AlignCenterRight:
		justify = ui.End
	}

	box := ui.Box(c).
		Key(keyId).
		Background(bg).
		Border(borderWidth, borderColor).
		Radius(radius).
		PaddingX(padX).
		PaddingY(padY).
		MinWidth(minW).
		MinHeight(minH).
		AlignItems(ui.Center).
		Justify(justify).
		Shadow(0, shadowY, shadowBlur, shadowSpread, shadowColor)

	if insetGlow.A > 0 {
		// A restrained inner light gives thick PBT, glass, neon, and retro caps
		// a readable surface instead of a flat sticker on the overlay. Glass
		// asks for a tight one: a wide wash reads as a smudge, not an edge.
		box = box.InsetShadow(0, 1, insetBlur, 0, insetGlow)
	}

	if cfg.Color.UseGradient && bg.A > 0 {
		to := borderColor
		if to.A == 0 {
			to = bg.Mix(ui.RGBA(255, 255, 255, 0.35), 0.35)
		}
		box.Gradient(bg, to, 135)
	}

	if yOffset > 0 {
		box.Translate(0, yOffset)
	}

	box.Children(func() {
		iconColor := textColor
		if isMouse {
			iconColor = mouseAccent
		}
		if isPureIcon {
			iconSize := fontSize * 1.05
			ui.Icon(c, icon).Size(iconSize, iconSize).TextColor(iconColor)
		} else if shouldShowIcon {
			iconSize := fontSize * 0.88
			showIcon := func() {
				ui.Icon(c, icon).Size(iconSize, iconSize).TextColor(iconColor)
			}
			showText := func() {
				if label != "" && !isArrow {
					ui.Text(c, label).
						FontSize(fontSize).
						Bold().
						SingleLine().
						TextColor(textColor)
				}
			}
			ui.Row(c).AlignItems(ui.Center).Justify(ui.Center).Gap(6).Children(func() {
				if cfg.Layout.IconAlignment == "flex-end" {
					showText()
					showIcon()
					return
				}
				showIcon()
				showText()
			})
		} else {
			ui.Text(c, label).
				FontSize(fontSize).
				Bold().
				SingleLine().
				TextColor(textColor)
		}

		if streak > 1 {
			renderStreakBadge(c, fmt.Sprintf("%s-streak", keyId), streak, t, cfg)
		}
	})
}

// formatKeycapLabel formats keycap text according to variant (text, text-short, icon), layout flags, and caps.
func formatKeycapLabel(label string, isMod bool, cfg config.AppConfig) string {
	if isMouseKey(label) {
		return formatMouseLabel(label)
	}

	norm := strings.ToLower(strings.TrimSpace(label))
	v := cfg.Text.Variant
	if v == "" {
		v = "text-short"
	}

	// 1. Variant == "icon"
	if v == "icon" {
		switch norm {
		case "ctrl", "control":
			return "⌃"
		case "shift":
			return "⇧"
		case "alt", "alternate", "menu":
			return "⌥"
		case "win", "super", "cmd", "command":
			if cfg.Layout.ShowIcon {
				return "⊞"
			}
			return "⌘"
		case "esc", "escape":
			return "⎋"
		case "enter", "return", "↵":
			return "↵"
		case "backspace", "⌫", "back", "bksp":
			return "⌫"
		case "tab", "⇥":
			return "⇥"
		case "space", "␣":
			return "␣"
		case "delete", "del":
			return "⌦"
		case "up", "arrowup", "↑":
			return "↑"
		case "down", "arrowdown", "↓":
			return "↓"
		case "left", "arrowleft", "←":
			return "←"
		case "right", "arrowright", "→":
			return "→"
		case "capslock":
			return "⇪"
		case "pageup", "pgup":
			return "⇞"
		case "pagedown", "pgdn":
			return "⇟"
		}
	}

	// 2. Variant == "text" (Full Text)
	if v == "text" {
		switch norm {
		case "ctrl", "control":
			return "Control"
		case "shift":
			return "Shift"
		case "alt", "alternate", "menu":
			return "Alt"
		case "win", "super", "cmd", "command":
			return "Windows"
		case "esc", "escape":
			return "Escape"
		case "enter", "return", "↵":
			return "Enter"
		case "backspace", "⌫", "back", "bksp":
			return "Backspace"
		case "tab", "⇥":
			return "Tab"
		case "space", "␣":
			return "Space"
		case "delete", "del":
			return "Delete"
		case "up", "arrowup", "↑":
			return "Up Arrow"
		case "down", "arrowdown", "↓":
			return "Down Arrow"
		case "left", "arrowleft", "←":
			return "Left Arrow"
		case "right", "arrowright", "→":
			return "Right Arrow"
		case "capslock":
			return "Caps Lock"
		case "pageup", "pgup":
			return "Page Up"
		case "pagedown", "pgdn":
			return "Page Down"
		}
	}

	// 3. Variant == "text-short" (Default)
	res := label
	switch norm {
	case "ctrl", "control":
		res = "Ctrl"
	case "shift":
		res = "Shift"
	case "alt", "alternate", "menu":
		res = "Alt"
	case "win", "super", "cmd", "command":
		if runtime.GOOS == "darwin" {
			res = "Cmd"
		} else {
			res = "Win"
		}
	case "esc", "escape":
		res = "Esc"
	case "enter", "return":
		if cfg.Layout.ShowSymbol {
			res = "↵"
		} else {
			res = "Enter"
		}
	case "backspace", "back", "bksp":
		if cfg.Layout.ShowSymbol {
			res = "⌫"
		} else {
			res = "Bksp"
		}
	case "tab":
		if cfg.Layout.ShowSymbol {
			res = "⇥"
		} else {
			res = "Tab"
		}
	case "space":
		if cfg.Layout.ShowSymbol {
			res = "␣"
		} else {
			res = "Space"
		}
	case "delete", "del":
		res = "Del"
	case "up", "arrowup":
		if cfg.Layout.ShowSymbol {
			res = "↑"
		} else {
			res = "Up"
		}
	case "down", "arrowdown":
		if cfg.Layout.ShowSymbol {
			res = "↓"
		} else {
			res = "Down"
		}
	case "left", "arrowleft":
		if cfg.Layout.ShowSymbol {
			res = "←"
		} else {
			res = "Left"
		}
	case "right", "arrowright":
		if cfg.Layout.ShowSymbol {
			res = "→"
		} else {
			res = "Right"
		}
	}

	// If showSymbol is explicitly false and res is a symbol, convert to text
	if !cfg.Layout.ShowSymbol {
		switch res {
		case "↵":
			res = "Enter"
		case "⌫":
			res = "Bksp"
		case "⇥":
			res = "Tab"
		case "␣":
			res = "Space"
		case "↑":
			res = "Up"
		case "↓":
			res = "Down"
		case "←":
			res = "Left"
		case "→":
			res = "Right"
		}
	}

	// 4. Capitalization
	if cfg.Text.Caps == "uppercase" {
		res = strings.ToUpper(res)
	} else if cfg.Text.Caps == "lowercase" {
		res = strings.ToLower(res)
	}

	return res
}

// parseHexSafe parses a hex string safely without panicking.
func parseHexSafe(hexStr string, fallback ui.Color) ui.Color {
	if hexStr == "" {
		return fallback
	}
	defer func() {
		_ = recover()
	}()
	return ui.Hex(hexStr)
}

// renderStreakBadge renders a Keynote/iOS style multiplier badge anchored to the keycap top-right corner.
func renderStreakBadge(c *ui.Context, badgeId string, streak int, t theme.ThemeConfig, cfg config.AppConfig) {
	bg := t.StreakBadgeBackground
	textColor := t.StreakBadgeTextColor
	if bg.A == 0 {
		bg = ui.Hex("#ef4444")
		textColor = ui.Hex("#ffffff")
	}

	// The badge is drawn in cap-relative units so it keeps its proportions when
	// the user scales the caps up: at the default 20 px it is exactly the size it
	// has always been, and at 62 px it is no longer a speck stuck to the corner of
	// a key it no longer matches. The curve is deliberately shallower than the
	// cap's own, though — a badge that grew in lockstep with a 3x cap would
	// swallow the letter it is anchored to.
	badgeScale := 1 + (effectiveFontSize(t, cfg)/20-1)*0.3
	if badgeScale < 1 {
		badgeScale = 1
	}

	badge := ui.Box(c).
		Key(fmt.Sprintf("%s-%d", badgeId, streak)).
		Attach(ui.AnchorTopRight, ui.AnchorCenter).
		Background(bg).
		Radius(10*badgeScale).
		PaddingX(6*badgeScale).
		PaddingY(1.5*badgeScale).
		MinWidth(18*badgeScale).
		MinHeight(18*badgeScale).
		AlignItems(ui.Center).
		Justify(ui.Center).
		Border(1.2*badgeScale, ui.RGBA(255, 255, 255, 0.45)).
		Shadow(0, 3*badgeScale, 8*badgeScale, 0, bg.Alpha(0.50))

	if cfg.Appearance.Animation != "none" {
		badge = badge.Transition(ui.ElementTransition{
			Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.18}},
			Enter:     &ui.Motion{Scale: 0.82, Opacity: 0},
			Exit:      &ui.Motion{Scale: 0.88, Opacity: 0},
		})
	}

	badge.Children(func() {
		ui.Textf(c, "×%d", streak).
			FontSize(11 * badgeScale).
			Bold().
			SingleLine().
			TextColor(textColor)
	})
}

// renderActionSubtitle renders the contextual action subtitle pill,
// dynamically adapting to the active theme and custom color palette.
func renderActionSubtitle(c *ui.Context, chord engine.DisplayChord, t theme.ThemeConfig, cfg config.AppConfig) {
	translated := i18n.T(chord.Action.Title)

	// The intent pill always takes the theme's dedicated subtitle pair. Reusing
	// keycap colors made "Copy/Paste" disappear on light keycaps or light-on-light
	// themes; the theme author already chose a high-contrast pair for this job.
	bg := t.SubtitleBackground
	textColor := t.SubtitleTextColor
	fontSize := t.SubtitleFontSize
	if fontSize <= 0 {
		fontSize = 11.5
	}
	radius := effectiveRadius(t, cfg) * 0.55
	if radius < 6 {
		radius = 6
	}
	if radius > 12 {
		radius = 12
	}

	ui.Box(c).
		Background(bg).
		Radius(radius).
		PaddingX(11).
		PaddingY(3).
		Border(1, textColor.Alpha(0.35)).
		Shadow(0, 4, 12, 0, bg.Alpha(0.40)).
		Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(7).Children(func() {
				// Application badge: makes the interpretation auditable at a
				// glance ("why does Ctrl+P say Quick Open?" because VS Code).
				if cfg.ShowAppBadge && chord.AppName != "" {
					ui.Box(c).
						Padding(2, 7).
						Radius(5).
						Background(textColor.Alpha(0.12)).
						Children(func() {
							ui.Text(c, chord.AppName).
								FontSize(fontSize - 1).
								Bold().
								SingleLine().
								TextColor(textColor.Alpha(0.82))
						})
				}
				if chord.IsSequence {
					ui.Text(c, "⇢").
						FontSize(fontSize + 1).
						Bold().
						TextColor(textColor.Alpha(0.55))
				}
				ui.Text(c, translated).
					FontSize(fontSize).
					Bold().
					SingleLine().
					TextColor(textColor)
			})
		})
}

func isMouseKey(name string) bool {
	n := strings.ToLower(name)
	return strings.Contains(n, "click") || strings.Contains(n, "scroll") ||
		strings.Contains(n, "drag") || strings.Contains(n, "mouse") ||
		strings.Contains(n, "wheel")
}

func formatMouseLabel(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "left click", "mouse left", "left":
		return i18n.T("Left Click")
	case "right click", "mouse right", "right":
		return i18n.T("Right Click")
	case "middle click", "mouse middle", "middle":
		return i18n.T("Middle Click")
	case "scroll up", "wheel up":
		return i18n.T("Scroll Up")
	case "scroll down", "wheel down":
		return i18n.T("Scroll Down")
	case "drag", "mouse drag":
		return i18n.T("Drag")
	default:
		return i18n.T(name)
	}
}

// renderPrivacyShield renders the [🔒 Private] badge when sensitive data is detected.
func renderPrivacyShield(c *ui.Context, t theme.ThemeConfig) {
	ui.Box(c).
		Background(t.SubtitleBackground).
		Border(1.5, t.StreakBadgeBackground).
		Radius(12).
		PaddingX(16).
		PaddingY(8).
		Children(func() {
			ui.Text(c, "🔒  Private Keystroke").
				FontSize(16).
				Bold().
				TextColor(t.SubtitleTextColor)
		})
}
