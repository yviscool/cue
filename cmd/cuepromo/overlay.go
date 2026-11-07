package main

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// The frame is one transparent window painted in explicit layers.
//
// mygo paints flow children first and absolute children afterwards, in document
// order, so absolute siblings stack in the order they are built. That is the
// whole trick: build the aurora, then the vignette, then the product, then the
// type, then the flash, and the composition holds. Wrapping the product render
// in an absolute layer is what lets it sit *above* the backdrop — left in flow,
// it would have been painted before every absolute sibling and the background
// would have covered the keycaps.
const (
	// Type anchors. The frame is 1920x1040 and is cropped to 16:9 later, so the
	// safe band is roughly y 60..980.
	tagX, tagY = 96, 56 // act slate, top left
	headY      = 104    // act headline, centred
	bigY       = 392    // act 1 wordmark, centred
	subY       = 500    // act 1 second line, centred
	footY      = 846    // act footnote, centred
	dotsY      = 900    // theme strip, centred
)

// copySpec is everything the overlay prints in a given act.
type copySpec struct {
	tag  string // top-left slate: "01 / 05 — AZURE PULSE"
	head string // centred headline
	big  string // act 1 wordmark
	sub  string // act 1 second line
	foot string // centred footnote, above the theme strip
}

func (f *film) render(c *ui.Context) {
	now := time.Now()
	// Copy fades and flashes are time-driven, so keep the frame loop running.
	c.Invalidate()

	f.mu.Lock()
	mode, cp, copyAt, active := f.mode, f.cp, f.copyAt, f.active
	f.mu.Unlock()

	c.Root().Background(ui.Transparent)

	// 1. Backdrop and atmosphere.
	f.scene.paintBackdrop(c, now)
	paintVignette(c)

	if mode == modeCard {
		f.renderCard(c, now)
	} else {
		// 2. The product, in its own absolute layer so it clears the backdrop.
		ui.Column(c).Absolute().Left(0).Top(0).
			Width(float32(stageW)).Height(float32(stageH)).Children(func() {
			f.vState.Render(c)
		})

		// 3. Type.
		f.renderCopy(c, now, cp, copyAt)
		// 4. The theme strip, which is the one graphic that says "many looks"
		// without a word of copy.
		f.renderStrip(c, active)
	}

	// 5. Camera flash, over everything.
	if v, tint := f.scene.flashValue(now); v > 0.004 {
		ui.Box(c).Absolute().Left(0).Top(0).
			Size(float32(stageW), float32(stageH)).
			Background(tint).Opacity(v)
	}
}

func (f *film) renderCopy(c *ui.Context, now time.Time, cp copySpec, copyAt time.Time) {
	accent := f.scene.accentNow()

	// Type rises and fades in over ~380 ms whenever the copy changes, which is
	// what makes a hard cut between beats feel edited instead of switched.
	fade, rise := reveal(now, copyAt)

	if cp.tag != "" {
		ui.Column(c).Absolute().Left(tagX).Top(tagY).Opacity(fade).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(13).Children(func() {
				ui.Box(c).Size(7, 7).Radius(3.5).Background(accent.Alpha(0.9))
				ui.Text(c, cp.tag).FontSize(13).LetterSpacing(4.2).
					TextColor(ui.Hex("#8FB6D8").Alpha(0.85)).SingleLine()
			})
			ui.Box(c).Size(46, 1.5).Background(accent.Alpha(0.5))
		})
	}

	if cp.head != "" {
		centered(c, headY+rise, fade, func() {
			ui.Text(c, cp.head).FontSize(42).Bold().LetterSpacing(1.5).
				TextColor(ui.Hex("#EAF6FF")).SingleLine()
		})
	}

	if cp.big != "" {
		centered(c, bigY+rise, fade, func() {
			ui.Text(c, cp.big).FontSize(80).Bold().LetterSpacing(4).
				TextColor(ui.Hex("#F4FAFF")).SingleLine()
		})
	}

	if cp.sub != "" {
		centered(c, subY+rise, fade, func() {
			ui.Text(c, cp.sub).FontSize(26).LetterSpacing(2).
				TextColor(ui.Hex("#7E9BBD")).SingleLine()
		})
	}

	if cp.foot != "" {
		centered(c, footY+rise, fade*0.9, func() {
			ui.Text(c, cp.foot).FontSize(20).LetterSpacing(2.4).
				TextColor(ui.Hex("#5C7695")).SingleLine()
		})
	}
}

// renderStrip is a row of dots, one per theme, with the one on screen lit. It
// turns "19 themes" from a claim into something the viewer can count, and it
// gives the montage act a moving element to follow.
func (f *film) renderStrip(c *ui.Context, active int) {
	// Nothing lit means no theme is on screen yet — the cold open, and the beat
	// of air before the card. A row of nineteen unlit dots would read as an
	// unfinished element rather than as restraint.
	if len(f.swatches) == 0 || active < 0 {
		return
	}
	ui.Row(c).Absolute().Left(0).Top(dotsY).Width(float32(stageW)).
		Justify(ui.Center).AlignItems(ui.Center).Gap(21).Children(func() {
		for i, col := range f.swatches {
			size, alpha := float32(5), float32(0.30)
			if i == active {
				size, alpha = 10, 1
			}
			d := ui.Box(c).Size(size, size).Radius(size / 2).Background(col.Alpha(alpha))
			if i == active {
				d.Shadow(0, 0, 20, 4, col.Alpha(0.9))
			}
		}
	})
}

// renderCard is the end lockup: wordmark, promise, proof points, address.
func (f *film) renderCard(c *ui.Context, now time.Time) {
	accent := f.scene.accentNow()
	fade, _ := reveal(now, f.copyAt)

	ui.Column(c).Absolute().Left(0).Top(0).
		Width(float32(stageW)).Height(float32(stageH)).
		AlignItems(ui.Center).Justify(ui.Center).Gap(28).Opacity(fade).Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(22).Children(func() {
			ui.Box(c).Size(18, 18).Radius(9).Background(ui.Hex("#22D3EE")).
				Shadow(0, 0, 34, 5, ui.Hex("#22D3EE").Alpha(0.85))
			ui.Text(c, "Cue").FontSize(104).Bold().LetterSpacing(-2).
				TextColor(ui.Hex("#F5FBFF"))
		})
		ui.Text(c, "See what you press.").FontSize(31).LetterSpacing(2).
			TextColor(ui.Hex("#9CB4CF")).SingleLine()
		ui.Box(c).Size(300, 1.5).Background(accent.Alpha(0.35))
		ui.Row(c).AlignItems(ui.Center).Gap(34).Children(func() {
			for _, s := range []string{"19 themes", "4 languages", "Pure Go, zero deps", "One .exe"} {
				ui.Text(c, s).FontSize(19).LetterSpacing(1.5).
					TextColor(ui.Hex("#5A7292")).SingleLine()
			}
		})
		ui.Text(c, "github.com/yviscool/cue").FontSize(21).LetterSpacing(2).
			TextColor(ui.Hex("#3FA9C9")).SingleLine()
	})
}

// centered places a centred row at an absolute y offset in the frame.
func centered(c *ui.Context, top, fade float32, children func()) {
	ui.Column(c).Absolute().Left(0).Top(top).Width(float32(stageW)).
		AlignItems(ui.Center).Opacity(fade).Children(children)
}

// reveal returns the opacity and vertical offset of copy that has just been set.
func reveal(now time.Time, at time.Time) (float32, float32) {
	if at.IsZero() {
		return 1, 0
	}
	t := clamp01(float32(now.Sub(at).Seconds()) / 0.38)
	return easeOut(t), (1 - easeOut(t)) * 16
}
