package hud

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/engine"
	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/theme"
)

// mouseBodySvg returns SVG icon data replicating Cue's mouse indicator body,
// highlighting the clicked button or scroll direction with dedicated colors:
// Left Click: Primary theme/custom color (#00FF6A by default)
// Right Click: High-contrast Sky Blue (#38bdf8)
// Middle Click: Amber Gold (#facc15)
func mouseBodySvg(btn hook.MouseButton, wheel int, primaryColorHex string) *ui.SVG {
	if primaryColorHex == "" {
		primaryColorHex = "#00FF6A"
	}
	leftFill := "none"
	leftStroke := "white"
	rightFill := "none"
	rightStroke := "white"
	wheelStroke := "white"
	wheelFill := "none"
	upArrow := ""
	downArrow := ""

	switch btn {
	case hook.ButtonLeft:
		leftFill = primaryColorHex
		leftStroke = primaryColorHex
	case hook.ButtonRight:
		rightFill = "#38bdf8"
		rightStroke = "#38bdf8"
	case hook.ButtonMiddle:
		wheelStroke = "#facc15"
		wheelFill = "#facc15"
	}

	if wheel > 0 {
		wheelStroke = primaryColorHex
		upArrow = fmt.Sprintf(`<path d="M16.5 17V9M13 12.5L16.5 9L20 12.5" stroke="%s" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>`, primaryColorHex)
	} else if wheel < 0 {
		wheelStroke = primaryColorHex
		downArrow = fmt.Sprintf(`<path d="M16.5 9V17M13 13.5L16.5 17L20 13.5" stroke="%s" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>`, primaryColorHex)
	}

	svg := fmt.Sprintf(`<svg width="33" height="43" viewBox="0 0 33 43" fill="none" xmlns="http://www.w3.org/2000/svg">
		<!-- Mouse Body Outer Shell -->
		<path d="M16.5 41.5C28.5 41.5 31.5 32.48 31.5 21.5C31.5 10.52 28.5 1.5 16.5 1.5C4.5 1.5 1.5 10.52 1.5 21.5C1.5 32.48 4.5 41.5 16.5 41.5Z" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
		<!-- Top Center Divider -->
		<path d="M16.5 9.5V1.5" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
		<!-- Left Button Lobe -->
		<path d="M16.5 1.5C9 1.5 3 7 2.5 18H16.5V1.5Z" fill="%s" fill-opacity="0.85" stroke="%s" stroke-width="1.5"/>
		<!-- Right Button Lobe -->
		<path d="M16.5 1.5C24 1.5 30 7 30.5 18H16.5V1.5Z" fill="%s" fill-opacity="0.85" stroke="%s" stroke-width="1.5"/>
		<!-- Center Scroll Wheel -->
		<rect x="13.5" y="9.5" width="6" height="10" rx="3" stroke="%s" stroke-width="2.2" fill="%s" fill-opacity="0.9"/>
		%s
		%s
	</svg>`, leftFill, leftStroke, rightFill, rightStroke, wheelStroke, wheelFill, upArrow, downArrow)

	return ui.MustParseSVG([]byte(svg))
}

// RenderMouseOverlay renders the complete Cue-replicated mouse visualization at cursor (x, y):
// 1. Click Ripple FX (tactile contact puck, expanding translucent water ripple, and outer ambient halo)
//    - Left click: Primary color (default vibrant neon #00FF6A)
//    - Right click: Secondary high-contrast Ice Blue (#38bdf8)
//    - Middle click: Amber Gold (#facc15)
//    - KeepHighlight: Steady cursor spotlight halo when idling
// 2. Button & Scroll Indicator (pill badge beside cursor showing mouse silhouette with color-coded active lobes)
func RenderMouseOverlay(c *ui.Context, state engine.MouseVisualState, cfg config.AppConfig, t theme.ThemeConfig, scale float64, originX, originY int) {
	if scale <= 0 {
		scale = 1.0
	}
	cx := float32(state.X-originX) / float32(scale)
	cy := float32(state.Y-originY) / float32(scale)

	// Don't render if cursor position is off-screen
	if cx <= 0 && cy <= 0 {
		return
	}

	now := time.Now()
	clickAge := 999 * time.Second
	if !state.Timestamp.IsZero() {
		clickAge = now.Sub(state.Timestamp)
	}
	isRecentlyClicked := clickAge < 260*time.Millisecond
	showClickRipple := cfg.Mouse.ShowClicks && (state.Active || isRecentlyClicked)
	showHighlight := cfg.Mouse.KeepHighlight || showClickRipple

	// Primary and secondary button colors
	primaryHex := cfg.Mouse.Color
	if primaryHex == "" {
		primaryHex = "#00FF6A"
	}
	leftColor := ui.Hex(primaryHex)
	rightColor := ui.Hex("#38bdf8") // Distinctive ice blue for right click
	middleColor := ui.Hex("#facc15") // Radiant amber gold for middle click

	// Determine active click color based on button pressed
	clickColor := leftColor
	switch state.Button {
	case hook.ButtonRight:
		clickColor = rightColor
	case hook.ButtonMiddle:
		clickColor = middleColor
	default:
		clickColor = leftColor
	}

	// 1. Click Ripple FX & Cursor Ring Highlight
	if showHighlight {
		baseSize := float32(cfg.Mouse.Size)
		if baseSize <= 0 {
			baseSize = 50
		}

		if showClickRipple {
			clickKey := state.Timestamp.UnixMilli()

			// Layer 1: Core Contact Puck (Tactile punch right at cursor point)
			coreSize := baseSize * 0.44
			ui.Box(c.Key(fmt.Sprintf("mouse-core-%d", clickKey))).
				Absolute().
				Left(cx - coreSize/2).
				Top(cy - coreSize/2).
				Size(coreSize, coreSize).
				Radius(coreSize / 2).
				Background(clickColor.Alpha(0.24)).
				Border(2.0, clickColor).
				Shadow(0, 0, 12, 0, clickColor.Alpha(0.65)).
				Transition(ui.ElementTransition{
					Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 200 * time.Millisecond, Bounce: 0.25}},
					Enter:     &ui.Motion{Scale: 0.40, Opacity: 0.95},
					Exit:      &ui.Motion{Scale: 0.85, Opacity: 0},
				})

			// Layer 2: Expanding Primary Water Ripple (Radiates outward and fades)
			rippleSize := baseSize * 1.55
			ui.Box(c.Key(fmt.Sprintf("mouse-ripple-%d", clickKey))).
				Absolute().
				Left(cx - rippleSize/2).
				Top(cy - rippleSize/2).
				Size(rippleSize, rippleSize).
				Radius(rippleSize / 2).
				Border(2.2, clickColor.Alpha(0.85)).
				Shadow(0, 0, 16, 0, clickColor.Alpha(0.40)).
				Transition(ui.ElementTransition{
					Duration: 280 * time.Millisecond,
					Ease:     ui.EaseOut,
					Enter:    &ui.Motion{Scale: 0.35, Opacity: 0.85},
					Exit:     &ui.Motion{Scale: 1.65, Opacity: 0},
				})

			// Layer 3: Secondary Ambient Water Halo (Soft outer expansion wave)
			haloSize := baseSize * 2.15
			ui.Box(c.Key(fmt.Sprintf("mouse-halo-%d", clickKey))).
				Absolute().
				Left(cx - haloSize/2).
				Top(cy - haloSize/2).
				Size(haloSize, haloSize).
				Radius(haloSize / 2).
				Border(1.4, clickColor.Alpha(0.45)).
				Shadow(0, 0, 22, 0, clickColor.Alpha(0.25)).
				Transition(ui.ElementTransition{
					Duration: 340 * time.Millisecond,
					Ease:     ui.EaseOut,
					Enter:    &ui.Motion{Scale: 0.45, Opacity: 0.50},
					Exit:     &ui.Motion{Scale: 2.10, Opacity: 0},
				})
		} else if cfg.Mouse.KeepHighlight {
			// Steady idle cursor spotlight ring when KeepHighlight is on
			steadySize := baseSize
			borderWidth := baseSize / 18.0
			if borderWidth < 2.0 {
				borderWidth = 2.0
			}
			ui.Box(c.Key("mouse-cursor-steady-ring")).
				Absolute().
				Left(cx - steadySize/2).
				Top(cy - steadySize/2).
				Size(steadySize, steadySize).
				Radius(steadySize / 2).
				Border(borderWidth, leftColor.Alpha(0.70)).
				Shadow(0, 0, 12, 0, leftColor.Alpha(0.35)).
				Transition(ui.ElementTransition{
					Duration: 180 * time.Millisecond,
					Ease:     ui.EaseOut,
					Enter:    &ui.Motion{Scale: 0.75, Opacity: 0},
					Exit:     &ui.Motion{Scale: 1.25, Opacity: 0},
				})
		}
	}

	// 2. Button & Scroll Indicator (Cue mouse silhouette)
	isRecentlyActiveIndicator := clickAge < 260*time.Millisecond
	showIndicator := cfg.Mouse.KeepIndicator || (cfg.Mouse.ShowIndicator && (state.Active || state.Wheel != 0 || isRecentlyActiveIndicator))

	if showIndicator {
		indSize := float32(cfg.Mouse.IndicatorSize)
		if indSize <= 0 {
			indSize = 40
		}
		ox := float32(cfg.Mouse.IndicatorOffsetX)
		oy := float32(cfg.Mouse.IndicatorOffsetY)
		if ox == 0 && oy == 0 {
			ox, oy = 24, 24
		}

		w := indSize * 0.85
		h := indSize

		ui.Box(c.Key("mouse-btn-indicator")).
			Absolute().
			Left(cx + ox).
			Top(cy + oy).
			Width(w + 14).
			Height(h + 14).
			Background(ui.RGBA(12, 12, 16, 0.78)).
			Radius(10).
			Border(1, ui.RGBA(255, 255, 255, 0.16)).
			Shadow(0, 6, 20, 0, ui.RGBA(0, 0, 0, 0.50)).
			AlignItems(ui.Center).
			Justify(ui.Center).
			Transition(ui.ElementTransition{
				Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.18}},
				Enter:     &ui.Motion{Scale: 0.86, Opacity: 0},
				Exit:      &ui.Motion{Scale: 0.92, Opacity: 0},
			}).
			Children(func() {
				btn := state.Button
				if !state.Active && !isRecentlyActiveIndicator {
					btn = hook.ButtonNone
				}
				wheel := state.Wheel
				if !isRecentlyActiveIndicator {
					wheel = 0
				}
				ui.Icon(c, mouseBodySvg(btn, wheel, primaryHex)).
					Size(w, h).
					TextColor(ui.Hex("#ffffff"))
			})
	}
}

// RenderMouseHalo renders legacy mouse interaction badge in the key chord stream if needed.
func RenderMouseHalo(c *ui.Context, state engine.MouseVisualState, t theme.ThemeConfig) {
	if state.Button == hook.ButtonNone || time.Since(state.Timestamp) > 800*time.Millisecond {
		return
	}

	btnKey := "Left Click"
	switch state.Button {
	case hook.ButtonLeft:
		btnKey = "Left Click"
	case hook.ButtonRight:
		btnKey = "Right Click"
	case hook.ButtonMiddle:
		btnKey = "Middle Click"
	case hook.ButtonExtra1:
		btnKey = "Mouse 4"
	case hook.ButtonExtra2:
		btnKey = "Mouse 5"
	default:
		btnKey = state.Button.String()
	}

	btnName := "🖱️ " + i18n.T(btnKey)
	label := btnName
	if !state.Modifiers.IsEmpty() {
		label = fmt.Sprintf("%s + %s", state.Modifiers.String(), btnName)
	}

	ui.Box(c.Key("mouse-halo")).
		Background(t.SubtitleBackground).
		Border(1.5, t.MouseRippleColor).
		Radius(16).
		PaddingX(14).
		PaddingY(6).
		Shadow(0, 4, 12, 0, t.MouseRippleColor.Alpha(0.35)).
		Children(func() {
			ui.Text(c, label).
				FontSize(13).
				Bold().
				TextColor(t.MouseRippleColor)
		})
}
