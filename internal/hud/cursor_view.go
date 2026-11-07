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
func mouseBodySvg(btn hook.MouseButton, wheel int, mth theme.MouseThemeConfig) *ui.SVG {
	leftFill := "none"
	leftStroke := "rgba(255,255,255,0.70)"
	rightFill := "none"
	rightStroke := "rgba(255,255,255,0.70)"
	wheelStroke := "rgba(255,255,255,0.70)"
	wheelFill := "none"

	switch btn {
	case hook.ButtonLeft:
		leftFill = mth.LeftColor
		leftStroke = mth.LeftColor
	case hook.ButtonRight:
		rightFill = mth.RightColor
		rightStroke = mth.RightColor
	case hook.ButtonMiddle:
		wheelStroke = mth.WheelColor
		wheelFill = mth.WheelColor
	}

	centerElement := ""
	if wheel > 0 {
		// High-contrast Scroll Up indicator with refined directional arrow head
		centerElement = fmt.Sprintf(`
			<rect x="11.5" y="6" width="10" height="17" rx="5" fill="#0c0c12" fill-opacity="0.88" stroke="%s" stroke-width="1.4"/>
			<line x1="16.5" y1="19" x2="16.5" y2="12.6" stroke="%s" stroke-width="2.2" stroke-linecap="round"/>
			<polygon points="16.5,8.2 13.2,12.6 19.8,12.6" fill="%s"/>
		`, mth.WheelColor, mth.WheelColor, mth.WheelColor)
	} else if wheel < 0 {
		// High-contrast Scroll Down indicator with refined directional arrow head
		centerElement = fmt.Sprintf(`
			<rect x="11.5" y="6" width="10" height="17" rx="5" fill="#0c0c12" fill-opacity="0.88" stroke="%s" stroke-width="1.4"/>
			<line x1="16.5" y1="10" x2="16.5" y2="16.4" stroke="%s" stroke-width="2.2" stroke-linecap="round"/>
			<polygon points="16.5,20.8 13.2,16.4 19.8,16.4" fill="%s"/>
		`, mth.WheelColor, mth.WheelColor, mth.WheelColor)
	} else {
		// Idle or Middle Click scroll wheel capsule
		centerElement = fmt.Sprintf(`<rect x="13.5" y="8.5" width="6" height="11" rx="3" stroke="%s" stroke-width="2.2" fill="%s" fill-opacity="0.95"/>`, wheelStroke, wheelFill)
	}

	svg := fmt.Sprintf(`<svg width="33" height="43" viewBox="0 0 33 43" fill="none" xmlns="http://www.w3.org/2000/svg">
		<!-- Mouse Body Outer Shell -->
		<path d="M16.5 41.5C28.5 41.5 31.5 32.48 31.5 21.5C31.5 10.52 28.5 1.5 16.5 1.5C4.5 1.5 1.5 10.52 1.5 21.5C1.5 32.48 4.5 41.5 16.5 41.5Z" stroke="white" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round"/>
		<!-- Top Center Divider -->
		<path d="M16.5 9.5V1.5" stroke="white" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"/>
		<!-- Left Button Lobe -->
		<path d="M16.5 1.5C9 1.5 3 7 2.5 18H16.5V1.5Z" fill="%s" fill-opacity="0.92" stroke="%s" stroke-width="2.0"/>
		<!-- Right Button Lobe -->
		<path d="M16.5 1.5C24 1.5 30 7 30.5 18H16.5V1.5Z" fill="%s" fill-opacity="0.92" stroke="%s" stroke-width="2.0"/>
		<!-- Center Scroll Element (Idle capsule or Bold Directional Arrow) -->
		%s
	</svg>`, leftFill, leftStroke, rightFill, rightStroke, centerElement)

	return ui.MustParseSVG([]byte(svg))
}

// RenderMouseOverlay renders the complete Cue-replicated mouse visualization at cursor (x, y):
// 1. Click Ripple FX (tactile contact puck, expanding translucent water ripple, and outer ambient halo)
//   - Left click: Primary color (default vibrant neon #00FF6A)
//   - Right click: Secondary high-contrast Ice Blue (#38bdf8)
//   - Middle click: Amber Gold (#facc15)
//   - KeepHighlight: Steady cursor spotlight halo when idling
//
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

	clickLinger := time.Duration(cfg.Mouse.ClickLingerMs) * time.Millisecond
	if clickLinger <= 0 {
		clickLinger = 180 * time.Millisecond
	}

	now := time.Now()
	clickAge := 999 * time.Second
	if !state.Timestamp.IsZero() {
		clickAge = now.Sub(state.Timestamp)
	}
	isRecentlyClicked := clickAge < clickLinger
	showClickRipple := cfg.Mouse.ShowClicks && (state.Active || isRecentlyClicked)
	showHighlight := cfg.Mouse.KeepHighlight || showClickRipple

	// Primary and secondary button colors from MouseTheme
	mth := theme.MouseThemeByID(cfg.Mouse.ThemePreset)
	if cfg.Mouse.Color != "" && cfg.Mouse.ThemePreset == "" {
		mth.LeftColor = cfg.Mouse.Color
		mth.RingColor = cfg.Mouse.Color
	}
	leftColor := ui.Hex(mth.LeftColor)
	rightColor := ui.Hex(mth.RightColor)
	middleColor := ui.Hex(mth.WheelColor)
	ringColor := ui.Hex(mth.RingColor)

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
				Left(cx-coreSize/2).
				Top(cy-coreSize/2).
				Size(coreSize, coreSize).
				Radius(coreSize/2).
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
				Left(cx-rippleSize/2).
				Top(cy-rippleSize/2).
				Size(rippleSize, rippleSize).
				Radius(rippleSize/2).
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
				Left(cx-haloSize/2).
				Top(cy-haloSize/2).
				Size(haloSize, haloSize).
				Radius(haloSize/2).
				Border(1.4, clickColor.Alpha(0.45)).
				Shadow(0, 0, 22, 0, clickColor.Alpha(0.25)).
				Transition(ui.ElementTransition{
					Duration: 340 * time.Millisecond,
					Ease:     ui.EaseOut,
					Enter:    &ui.Motion{Scale: 0.45, Opacity: 0.50},
					Exit:     &ui.Motion{Scale: 2.10, Opacity: 0},
				})

			// Layer 4: Echo Ripple for Double / Multi-click (Concentric Pulse)
			if state.ClickCount > 1 {
				echoSize := baseSize * 1.85
				ui.Box(c.Key(fmt.Sprintf("mouse-echo-%d-%d", clickKey, state.ClickCount))).
					Absolute().
					Left(cx-echoSize/2).
					Top(cy-echoSize/2).
					Size(echoSize, echoSize).
					Radius(echoSize/2).
					Border(2.5, clickColor).
					Shadow(0, 0, 18, 0, clickColor.Alpha(0.60)).
					Transition(ui.ElementTransition{
						Duration: 240 * time.Millisecond,
						Ease:     ui.EaseOut,
						Enter:    &ui.Motion{Scale: 0.25, Opacity: 1.0},
						Exit:     &ui.Motion{Scale: 1.85, Opacity: 0},
					})
			}
		} else if cfg.Mouse.KeepHighlight {
			// Steady idle cursor spotlight ring when KeepHighlight is on
			steadySize := baseSize
			borderWidth := baseSize / 18.0
			if borderWidth < 2.0 {
				borderWidth = 2.0
			}
			ui.Box(c.Key("mouse-cursor-steady-ring")).
				Absolute().
				Left(cx-steadySize/2).
				Top(cy-steadySize/2).
				Size(steadySize, steadySize).
				Radius(steadySize/2).
				Border(borderWidth, ringColor.Alpha(0.70)).
				Shadow(0, 0, 12, 0, ringColor.Alpha(0.35)).
				Transition(ui.ElementTransition{
					Duration: 180 * time.Millisecond,
					Ease:     ui.EaseOut,
					Enter:    &ui.Motion{Scale: 0.75, Opacity: 0},
					Exit:     &ui.Motion{Scale: 1.25, Opacity: 0},
				})
		}
	}

	// 2. Button & Scroll Indicator (Cue mouse silhouette)
	badgeLinger := 1200 * time.Millisecond
	if cfg.General.LingerDurationMs > 0 {
		badgeLinger = time.Duration(cfg.General.LingerDurationMs) * time.Millisecond
		if badgeLinger > 1500*time.Millisecond {
			badgeLinger = 1500 * time.Millisecond
		} else if badgeLinger < 800*time.Millisecond {
			badgeLinger = 800 * time.Millisecond
		}
	}

	hasActiveBadge := state.ClickCount > 1 && clickAge < badgeLinger && state.Wheel == 0
	isRecentlyActiveIndicator := clickAge < clickLinger || hasActiveBadge
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

		boxRadius := indSize * 0.22
		if boxRadius < 6 {
			boxRadius = 6
		} else if boxRadius > 24 {
			boxRadius = 24
		}

		indicatorBox := ui.Box(c.Key("mouse-btn-indicator")).
			Absolute().
			Left(cx+ox).
			Top(cy+oy).
			Width(w+14).
			Height(h+14).
			Background(ui.RGBA(12, 12, 16, 0.78)).
			Radius(boxRadius).
			Border(1, ui.RGBA(255, 255, 255, 0.16)).
			Shadow(0, 6, 20, 0, ui.RGBA(0, 0, 0, 0.50)).
			AlignItems(ui.Center).
			Justify(ui.Center)

		if cfg.Appearance.Animation != "none" {
			indicatorBox = indicatorBox.Transition(ui.ElementTransition{
				Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.18}},
				Enter:     &ui.Motion{Scale: 0.86, Opacity: 0},
				Exit:      &ui.Motion{Scale: 0.92, Opacity: 0},
			})
		}

		indicatorBox.Children(func() {
			btn := state.Button
			// Lobe highlight turns off after clickLinger even while multi-click badge remains visible
			if !state.Active && clickAge >= clickLinger {
				btn = hook.ButtonNone
			}
			wheel := state.Wheel
			if !state.Active && clickAge >= clickLinger {
				wheel = 0
			}
			ui.Image(c, mouseBodySvg(btn, wheel, mth)).
				Size(w, h)

			// Multi-click Retina Pill Badge (e.g. ×2 for double-click, ×3 for triple-click)
			// Dynamically anchored to the physical button pressed: Left Click -> Top Left, Right -> Top Right, Middle -> Top Center
			if hasActiveBadge {
				badgeBg := ui.Hex(mth.BadgeColor)
				if badgeBg.A == 0 {
					badgeBg = t.StreakBadgeBackground
					if badgeBg.A == 0 {
						badgeBg = ui.Hex("#ef4444")
					}
				}
				badgeTextColor := ui.Hex("#ffffff")

				badgeFont := 10.5 * (0.8 + 0.2*(indSize/40.0))
				if badgeFont > 22 {
					badgeFont = 22
				}
				badgeRadius := 8 * (0.8 + 0.2*(indSize/40.0))
				if badgeRadius > 16 {
					badgeRadius = 16
				}
				badgePadX := 5 * (0.8 + 0.2*(indSize/40.0))
				badgePadY := 1.5 * (0.8 + 0.2*(indSize/40.0))

				// Anchor directly on top of the button lobe being struck
				anchor := ui.AnchorTopLeft
				if state.Button == hook.ButtonRight {
					anchor = ui.AnchorTopRight
				} else if state.Button == hook.ButtonMiddle {
					anchor = ui.AnchorTop
				}

				badgeBox := ui.Box(c.Key(fmt.Sprintf("mouse-badge-%d", state.ClickCount))).
					Attach(anchor, ui.AnchorCenter).
					Background(badgeBg).
					Radius(badgeRadius).
					PaddingX(badgePadX).
					PaddingY(badgePadY).
					MinWidth(badgeFont*1.6).
					MinHeight(badgeFont*1.4).
					AlignItems(ui.Center).
					Justify(ui.Center).
					Border(1, ui.RGBA(255, 255, 255, 0.40)).
					Shadow(0, 2, 6, 0, badgeBg.Alpha(0.50))

				if cfg.Appearance.Animation != "none" {
					badgeBox = badgeBox.Transition(ui.ElementTransition{
						Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.22}},
						Enter:     &ui.Motion{Scale: 0.65, Opacity: 0},
						Exit:      &ui.Motion{Scale: 0.80, Opacity: 0},
					})
				}

				badgeBox.Children(func() {
					ui.Textf(c, "×%d", state.ClickCount).
						FontSize(badgeFont).
						Bold().
						SingleLine().
						TextColor(badgeTextColor)
				})
			}
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
