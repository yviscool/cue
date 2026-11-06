package hud

import (
	"fmt"
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
	const d float32 = 28
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
	v.mu.RUnlock()

	// Ensure whole canvas root is 100% transparent
	c.Root().Background(ui.Transparent)

	hasMouse := snap.MouseState.Active || (!snap.MouseState.Timestamp.IsZero() && time.Since(snap.MouseState.Timestamp) < 300*time.Millisecond) || snap.MouseState.Wheel != 0
	shouldRenderMouse := cfg.Mouse.KeepHighlight || cfg.Mouse.KeepIndicator || (hasMouse && (cfg.Mouse.ShowIndicator || cfg.Mouse.ShowClicks))

	// 1. Mouse Overlay (Cursor Ring & Mouse Indicator at cursor coordinates)
	if shouldRenderMouse {
		RenderMouseOverlay(c, snap.MouseState, cfg, t, scale)
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

	if cfg.Appearance.FlexDirection == "col" || cfg.Appearance.FlexDirection == "column" {
		ui.Column(c.Key("key-overlay-container")).
			Fill().
			Justify(vAlign).
			AlignItems(hAlign).
			Padding(fmy, fmx, fmy, fmx).
			Gap(12).
			Children(renderContent)
	} else {
		ui.Row(c.Key("key-overlay-container")).
			Fill().
			Justify(hAlign).
			AlignItems(vAlign).
			Padding(fmy, fmx, fmy, fmx).
			Gap(14).
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
		trans = &ui.ElementTransition{
			Duration: animDuration,
			Ease:     ui.EaseOut,
			Enter:    &ui.Motion{Opacity: 0},
			Exit:     &ui.Motion{Opacity: 0},
		}
	case "zoom":
		springAnim := ui.Animation{
			Spring: ui.SpringOptions{
				Duration: animDuration,
				Bounce:   0.22,
			},
		}
		trans = &ui.ElementTransition{
			Animation: &springAnim,
			Enter:     &ui.Motion{Opacity: 0, Scale: 0.60},
			Exit:      &ui.Motion{Opacity: 0, Scale: 0.70},
		}
	case "float":
		trans = &ui.ElementTransition{
			Duration: animDuration,
			Ease:     ui.EaseOut,
			Enter:    &ui.Motion{Y: 20, Opacity: 0},
			Exit:     &ui.Motion{Y: -16, Opacity: 0},
		}
	case "slide":
		// Slide from, and return to, the nearest screen edge. A cue anchored at
		// the bottom should rise from below, not fly in from an unrelated side.
		dx, dy := slideVector(cfg.Appearance.Alignment)
		trans = &ui.ElementTransition{
			Duration: animDuration,
			Ease:     ui.EaseOut,
			Enter:    &ui.Motion{X: dx, Y: dy, Opacity: 0},
			Exit:     &ui.Motion{X: dx, Y: dy, Opacity: 0},
		}
	default: // "spring" or default
		springAnim := ui.Animation{
			Spring: ui.SpringOptions{
				Duration: animDuration,
				Bounce:   0.15,
			},
		}
		trans = &ui.ElementTransition{
			Animation: &springAnim,
			Enter:     &ui.Motion{Y: 12, Opacity: 0, Scale: 0.90},
			Exit:      &ui.Motion{Y: -8, Opacity: 0, Scale: 0.95},
		}
	}

	chordColumn := ui.Column(c).
		Key(chord.ID).
		Center().
		Gap(6)
	if trans != nil {
		chordColumn = chordColumn.Transition(*trans)
	}
	chordColumn.Children(func() {
		// Row containing the keycaps
		ui.Row(c).Center().Gap(8).AlignItems(ui.Center).Children(func() {
			if chord.IsPrivate {
				renderPrivacyShield(c, t)
				return
			}

			for i, key := range chord.Keys {
				renderKeyCap(c, fmt.Sprintf("%s-key-%d", chord.ID, i), key, t, cfg)
			}

			// If streak count > 1 and ShowPressCount enabled, render multiplier pill
			if chord.Streak > 1 && cfg.Layout.ShowPressCount {
				renderStreakPill(c, fmt.Sprintf("%s-streak", chord.ID), chord.Streak, t)
			}
		})

		// If matched with semantic action subtitle and SemanticSub enabled, render intent pill
		if chord.Action.Title != "" && !chord.IsPrivate && cfg.SemanticSub {
			renderActionSubtitle(c, chord.Action.Title, t, cfg)
		}
	})
}

// renderKeyCap renders a single tactile keycap driven by live configuration styles and physical models.
func renderKeyCap(c *ui.Context, keyId string, key engine.KeyCap, t theme.ThemeConfig, cfg config.AppConfig) {
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

	// Tactile hold physics: sink 3px downward on press
	yOffset := float32(0)
	if key.IsHeld {
		bg = t.KeyCapHeldBackground
		borderColor = t.KeyCapHeldBorder
		yOffset = 3.0
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
	padX := float32(14)
	padY := float32(8)
	shadowY := float32(4)
	shadowSpread := t.KeyCapShadowSpread
	shadowBlur := float32(8)
	shadowColor := borderColor.Alpha(0.28)

	switch style {
	case "minimal":
		borderWidth = 0
		borderColor = ui.Transparent
		shadowY = 1
		shadowBlur = 3
		shadowColor = ui.RGBA(0, 0, 0, 0.15)
		radius = radius * 0.7
	case "laptop":
		padX = 12
		padY = 6
		borderWidth = 1.0
		borderColor = borderColor.Alpha(0.4)
		shadowY = 2
		shadowBlur = 4
		shadowColor = ui.RGBA(0, 0, 0, 0.22)
		radius = radius * 0.6
	case "pbt":
		padX = 16
		padY = 10
		borderWidth = borderWidth + 1.0
		shadowY = 6
		shadowBlur = 10
		shadowColor = ui.RGBA(0, 0, 0, 0.45)
		radius = radius * 1.15
	case "glass":
		padX = 15
		padY = 8
		// A theme may already supply a translucent "glass" fill (Liquid Glass
		// does). Multiplying that by 0.45 again would wash the keycap away.
		if bg.A >= 128 {
			bg = bg.Alpha(0.45)
		}
		borderWidth = 1.5
		borderColor = ui.RGBA(255, 255, 255, 0.50)
		shadowY = 8
		shadowBlur = 16
		shadowColor = ui.RGBA(0, 0, 0, 0.35)
	case "neon":
		padX = 14
		padY = 8
		bg = ui.RGBA(10, 15, 25, 0.85)
		borderWidth = 2.0
		borderColor = ui.Hex("#00f0ff")
		textColor = ui.Hex("#00f0ff")
		shadowY = 0
		shadowBlur = 14
		shadowSpread = 2
		shadowColor = ui.RGBA(0, 240, 255, 0.65)
	case "retro":
		padX = 15
		padY = 9
		bg = ui.Hex("#ded5c4")
		borderWidth = 2.0
		borderColor = ui.Hex("#a89f8d")
		textColor = ui.Hex("#1c1c1a")
		shadowY = 5
		shadowBlur = 0
		shadowColor = ui.Hex("#786e5c")
		radius = radius * 0.5
	case "outline":
		padX = 14
		padY = 8
		bg = ui.RGBA(0, 0, 0, 0.25)
		borderWidth = 2.0
		borderColor = textColor
		shadowY = 0
		shadowBlur = 0
		shadowColor = ui.Transparent
	default: // "lowprofile"
		// Standard physical model
	}

	// Text label formatting taking variant, layout symbols, and caps into account
	label := formatKeycapLabel(key.Label, key.IsModifier, cfg)

	box := ui.Box(c).
		Key(keyId).
		Background(bg).
		Border(borderWidth, borderColor).
		Radius(radius).
		PaddingX(padX).
		PaddingY(padY).
		AlignItems(ui.Center).
		Justify(ui.Center).
		Shadow(0, shadowY, shadowBlur, shadowSpread, shadowColor)

	if yOffset > 0 {
		box.Translate(0, yOffset)
	}

	icon := getKeyIcon(key.Label, key.IsModifier, cfg)
	isMouse := isMouseKey(key.Label)
	isArrow := isArrowKey(key.Label)
	isPureIcon := (cfg.Text.Variant == "icon" && icon != nil) || (isArrow && icon != nil && cfg.Text.Variant != "text")
	shouldShowIcon := icon != nil && (isPureIcon || (cfg.Layout.ShowIcon && (key.IsModifier || isArrow || isMouse)) || isMouse)

	box.Children(func() {
		if isPureIcon {
			iconSize := fontSize * 1.05
			ui.Icon(c, icon).Size(iconSize, iconSize).TextColor(textColor)
		} else if shouldShowIcon {
			iconSize := fontSize * 0.88
			ui.Row(c).AlignItems(ui.Center).Justify(ui.Center).Gap(6).Children(func() {
				ui.Icon(c, icon).Size(iconSize, iconSize).TextColor(textColor)
				if label != "" && !isArrow {
					ui.Text(c, label).
						FontSize(fontSize).
						Bold().
						SingleLine().
						TextColor(textColor)
				}
			})
		} else {
			ui.Text(c, label).
				FontSize(fontSize).
				Bold().
				SingleLine().
				TextColor(textColor)
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

// renderStreakPill renders an Apple/Keynote-style multiplier pill badge.
func renderStreakPill(c *ui.Context, badgeId string, streak int, t theme.ThemeConfig) {
	// The multiplier badge is a reward signal, so it lands with a small spring
	// instead of appearing as another static chip.
	ui.Box(c).
		Key(badgeId).
		Background(t.StreakBadgeBackground).
		Radius(12).
		PaddingX(8).
		PaddingY(4).
		Shadow(0, 2, 6, 0, t.StreakBadgeBackground.Alpha(0.4)).
		Transition(ui.ElementTransition{
			Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.30}},
			Enter:     &ui.Motion{Scale: 0.55, Opacity: 0},
			Exit:      &ui.Motion{Scale: 0.70, Opacity: 0},
		}).
		Children(func() {
			ui.Textf(c, "×%d", streak).
				FontSize(13).
				Bold().
				TextColor(t.StreakBadgeTextColor)
		})
}

// renderActionSubtitle renders the contextual action subtitle pill,
// dynamically adapting to the active theme and custom color palette.
func renderActionSubtitle(c *ui.Context, title string, t theme.ThemeConfig, cfg config.AppConfig) {
	translated := i18n.T(title)

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
			ui.Text(c, translated).
				FontSize(fontSize).
				Bold().
				TextColor(textColor)
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
