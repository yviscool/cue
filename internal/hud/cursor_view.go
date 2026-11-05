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
// highlighting the clicked button or scroll direction in vibrant neon green (#00FF6A).
func mouseBodySvg(btn hook.MouseButton, wheel int) *ui.SVG {
	// Base stroke color: white (#ffffff)
	// Active highlight color: #00FF6A
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
		leftFill = "#00FF6A"
		leftStroke = "#00FF6A"
	case hook.ButtonRight:
		rightFill = "#00FF6A"
		rightStroke = "#00FF6A"
	case hook.ButtonMiddle:
		wheelStroke = "#00FF6A"
		wheelFill = "#00FF6A"
	}

	if wheel > 0 {
		wheelStroke = "#00FF6A"
		upArrow = `<path d="M16.5 17V9M13 12.5L16.5 9L20 12.5" stroke="#00FF6A" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>`
	} else if wheel < 0 {
		wheelStroke = "#00FF6A"
		downArrow = `<path d="M16.5 9V17M13 13.5L16.5 17L20 13.5" stroke="#00FF6A" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>`
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
// 1. Cursor Ring Highlight (animated scale, color, size, keep-highlight)
// 2. Button & Scroll Indicator (pill badge beside cursor showing mouse silhouette with green active button)
func RenderMouseOverlay(c *ui.Context, state engine.MouseVisualState, cfg config.AppConfig, t theme.ThemeConfig, scale float64) {
	if scale <= 0 {
		scale = 1.0
	}
	cx := float32(state.X) / float32(scale)
	cy := float32(state.Y) / float32(scale)

	// Don't render if cursor position is 0,0 or negative
	if cx <= 0 && cy <= 0 {
		return
	}

	now := time.Now()
	isRecentlyClicked := !state.Timestamp.IsZero() && now.Sub(state.Timestamp) < 220*time.Millisecond
	showClickRing := cfg.Mouse.ShowClicks && (state.Active || isRecentlyClicked)
	showHighlight := cfg.Mouse.KeepHighlight || showClickRing

	// 1. Cursor Ring Highlight
	if showHighlight {
		ringSize := float32(cfg.Mouse.Size)
		if ringSize <= 0 {
			ringSize = 50
		}

		// Ring color
		ringHex := cfg.Mouse.Color
		if ringHex == "" {
			ringHex = "#00FF6A"
		}
		ringColor := ui.Hex(ringHex)

		// Pulse scale: when clicking or active, shrink slightly (0.75x) for tactile punch
		renderSize := ringSize
		borderWidth := ringSize / 18.0
		if borderWidth < 2.0 {
			borderWidth = 2.0
		}
		if state.Active || isRecentlyClicked {
			renderSize = ringSize * 0.75
			borderWidth = borderWidth * 1.2
		}

		ui.Box(c.Key("mouse-cursor-ring")).
			Absolute().
			Left(cx - renderSize/2).
			Top(cy - renderSize/2).
			Size(renderSize, renderSize).
			Radius(renderSize / 2).
			Border(borderWidth, ringColor).
			Shadow(0, 0, 14, 0, ringColor.Alpha(0.55))
	}

	// 2. Button & Scroll Indicator (Cue mouse silhouette)
	isRecentlyActiveIndicator := !state.Timestamp.IsZero() && now.Sub(state.Timestamp) < 260*time.Millisecond
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
			Children(func() {
				btn := state.Button
				if !state.Active && !isRecentlyActiveIndicator {
					btn = hook.ButtonNone
				}
				wheel := state.Wheel
				if !isRecentlyActiveIndicator {
					wheel = 0
				}
				ui.Icon(c, mouseBodySvg(btn, wheel)).
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
