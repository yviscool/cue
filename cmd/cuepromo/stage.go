package main

import (
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/yviscool/cue/internal/theme"
)

// scene is the film's look: an aurora of drifting colour fields whose hue
// follows the theme currently on screen, plus the huge ghost numeral that gives
// each act depth.
//
// Everything here is a pure function of the wall clock rather than an
// accumulated simulation, so a dropped frame or a slow build never makes the
// background jump or drift out of sync with the script.
type scene struct {
	mu   sync.Mutex
	t0   time.Time
	acc  ramp // primary accent, taken from the theme's held border
	acc2 ramp // secondary accent, taken from the theme's cap text

	flashAt   time.Time
	flashPeak float32
	flashTint ui.Color

	ghost   string
	ghostAt time.Time
}

// ramp eases one colour to another, so a theme change reads as a dissolve
// rather than a cut.
type ramp struct {
	from, to ui.Color
	start    time.Time
	dur      time.Duration
}

func (r *ramp) value(now time.Time) ui.Color {
	if r.dur <= 0 || r.start.IsZero() {
		return r.to
	}
	t := float32(now.Sub(r.start)) / float32(r.dur)
	switch {
	case t >= 1:
		return r.to
	case t <= 0:
		return r.from
	}
	t = t * t * (3 - 2*t) // smoothstep: no visible corner at either end
	return mixColor(r.from, r.to, t)
}

func (r *ramp) set(now time.Time, c ui.Color, dur time.Duration) {
	r.from = r.value(now)
	r.to = c
	r.start = now
	r.dur = dur
}

// mixColor blends two colours componentwise. Alpha is carried along because the
// ramps are sometimes fed translucent theme colours.
func mixColor(a, b ui.Color, t float32) ui.Color {
	lerp := func(x, y uint8) uint8 {
		return uint8(float32(x) + (float32(y)-float32(x))*t + 0.5)
	}
	aa, ba := float32(a.A)/255, float32(b.A)/255
	return ui.RGBA(lerp(a.R, b.R), lerp(a.G, b.G), lerp(a.B, b.B), aa+(ba-aa)*t)
}

// init seeds both ramps so the opening frame already carries the theme colour
// instead of fading up from transparent black.
func (s *scene) init(th theme.ThemeConfig) {
	s.mu.Lock()
	s.t0 = time.Now()
	s.acc = ramp{to: th.KeyCapHeldBorder}
	s.acc2 = ramp{to: shiftHue(th.KeyCapTextColor, 148)}
	s.mu.Unlock()
}

// accentNow is the accent as it is being drawn right now, so the type on the
// overlay always matches the haze behind it.
func (s *scene) accentNow() ui.Color {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.acc.value(time.Now())
}

// follow points the aurora at a theme.
func (s *scene) follow(now time.Time, th theme.ThemeConfig) {
	s.mu.Lock()
	s.acc.set(now, th.KeyCapHeldBorder, 700*time.Millisecond)
	// The rim light is the same hue turned most of the way round the wheel, so
	// the pair always contrasts no matter which of the nineteen themes is up.
	s.acc2.set(now, shiftHue(th.KeyCapTextColor, 148), 900*time.Millisecond)
	s.mu.Unlock()
}

// flash queues a light hit. Peak is the veil's alpha; it decays with a short
// time constant so it reads as a camera flash, not a fade.
func (s *scene) flash(now time.Time, peak float32, tint ui.Color) {
	s.mu.Lock()
	s.flashAt, s.flashPeak, s.flashTint = now, peak, tint
	s.mu.Unlock()
}

func (s *scene) flashValue(now time.Time) (float32, ui.Color) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.flashAt.IsZero() {
		return 0, ui.Color{}
	}
	d := now.Sub(s.flashAt).Seconds()
	if d < 0 || d > 0.8 {
		return 0, ui.Color{}
	}
	return s.flashPeak * float32(math.Exp(-d/0.13)), s.flashTint
}

// setGhost places the giant background numeral. An empty string clears it.
func (s *scene) setGhost(now time.Time, text string) {
	s.mu.Lock()
	s.ghost, s.ghostAt = text, now
	s.mu.Unlock()
}

// ---- Painting -------------------------------------------------------------

// paintBackdrop draws the aurora and the ghost numeral. It is layer one of the
// frame: everything the product draws goes on top of it.
func (s *scene) paintBackdrop(c *ui.Context, now time.Time) {
	s.mu.Lock()
	if s.t0.IsZero() {
		s.t0 = now
	}
	accent := s.acc.value(now)
	accent2 := s.acc2.value(now)
	ghost, ghostAt := s.ghost, s.ghostAt
	s.mu.Unlock()

	t := float32(now.Sub(s.t0).Seconds())
	w, h := float32(stageW), float32(stageH)

	// Base wash: a touch of blue at the top falling away to near-black, which
	// keeps the edges from going flat grey once the video is compressed.
	ui.Box(c).Absolute().Left(0).Top(0).Size(w, h).LinearGradient(ui.LinearGradient{
		Angle: 168,
		Stops: []ui.GradientStop{
			{Color: ui.Hex("#081124"), At: 0},
			{Color: ui.Hex("#040A15"), At: 0.45},
			{Color: ui.Hex("#010309"), At: 1},
		},
	})

	// Colour fields on long, incommensurate orbits. The periods share no common
	// multiple, so the composition never visibly repeats within a take.
	//
	// The second field is deliberately the accent rotated round the wheel: a
	// key light and a rim light of two temperatures is what separates a shot
	// from a flat fill, and doing it from the theme's own hue means every one of
	// the nineteen themes gets its own complementary pairing for free.
	glow(c, 960+210*sin(t*0.23), 470+95*cos(t*0.17), 1320, accent, 0.38)
	glow(c, 400+160*cos(t*0.15+1.2), 250+120*sin(t*0.21+0.4), 1000, accent2, 0.24)
	glow(c, 1560+180*sin(t*0.12+2.4), 850+85*cos(t*0.16+1.1), 1150, accent, 0.16)

	// Floor: a wide field mostly below the frame, so only its upper arc shows.
	// It gives the keycaps something to stand on instead of floating in a void.
	glow(c, 960+70*sin(t*0.09), 1210, 1900, accent, 0.20)

	// Ghost numeral, inside the haze rather than behind it.
	if ghost != "" {
		fade := easeOut(clamp01(float32(now.Sub(ghostAt).Seconds()) / 0.6))
		ui.Column(c).Absolute().Left(0).Top(0).Width(w).Height(h).
			AlignItems(ui.Center).Justify(ui.Center).Children(func() {
			ui.Text(c, ghost).FontSize(780).Bold().
				TextColor(accent.Alpha(0.055 * fade)).SingleLine()
		})
	}
}

// paintVignette darkens the corners. The single cheapest thing that makes a
// frame look shot rather than rendered, so it sits above the aurora and below
// the subject.
func paintVignette(c *ui.Context) {
	w, h := float32(stageW), float32(stageH)
	ui.Box(c).Absolute().Left(0).Top(0).Size(w, h).RadialGradient(ui.RadialGradient{
		Center: ui.Point{X: 0.5, Y: 0.46},
		Radius: 1320,
		Stops: []ui.GradientStop{
			{Color: ui.RGBA(0, 0, 0, 0), At: 0.28},
			{Color: ui.RGBA(0, 0, 0, 0.40), At: 0.70},
			{Color: ui.RGBA(0, 0, 0, 0.86), At: 1},
		},
	})
}

// glow paints one soft circular colour field.
//
// The radius is set explicitly to half the box instead of left at mygo's
// farthest-corner default: with the default the gradient is still partly opaque
// where it meets the box edge, and every glow would show a hard rectangular
// seam — which is exactly what a soft light must never do.
func glow(c *ui.Context, cx, cy, diameter float32, col ui.Color, alpha float32) {
	r := diameter / 2
	ui.Box(c).Absolute().Left(cx-r).Top(cy-r).Size(diameter, diameter).
		RadialGradient(ui.RadialGradient{
			Center: ui.Point{X: 0.5, Y: 0.5},
			Radius: r,
			Stops: []ui.GradientStop{
				{Color: col.Alpha(alpha), At: 0},
				{Color: col.Alpha(alpha * 0.46), At: 0.40},
				{Color: col.Alpha(alpha * 0.13), At: 0.72},
				{Color: col.Alpha(0), At: 1},
			},
		})
}

func sin(v float32) float32 { return float32(math.Sin(float64(v))) }
func cos(v float32) float32 { return float32(math.Cos(float64(v))) }

// shiftHue rotates a colour around the hue wheel.
//
// Saturation is floored and lightness is pulled into a mid band, so a nearly
// grey theme still produces a rim light with actual colour in it instead of a
// second helping of the same wash.
func shiftHue(c ui.Color, degrees float32) ui.Color {
	r, g, b := float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	mx := maxF(r, maxF(g, b))
	mn := minF(r, minF(g, b))
	l := (mx + mn) / 2

	var h, s float32
	if d := mx - mn; d > 0 {
		if l > 0.5 {
			s = d / (2 - mx - mn)
		} else {
			s = d / (mx + mn)
		}
		switch mx {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		default:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	h = float32(math.Mod(float64(h+degrees/360+1), 1))
	if s < 0.55 {
		s = 0.55
	}
	if l < 0.42 {
		l = 0.42
	}
	if l > 0.72 {
		l = 0.72
	}

	q := l + s - l*s
	if l < 0.5 {
		q = l * (1 + s)
	}
	p := 2*l - q

	hue := func(t float32) float32 {
		if t < 0 {
			t++
		}
		if t > 1 {
			t--
		}
		switch {
		case t < 1.0/6:
			return p + (q-p)*6*t
		case t < 0.5:
			return q
		case t < 2.0/3:
			return p + (q-p)*(2.0/3-t)*6
		}
		return p
	}

	to8 := func(v float32) uint8 {
		if v < 0 {
			v = 0
		}
		if v > 1 {
			v = 1
		}
		return uint8(v*255 + 0.5)
	}
	return ui.RGBA(to8(hue(h+1.0/3)), to8(hue(h)), to8(hue(h-1.0/3)), float32(c.A)/255)
}

func maxF(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func minF(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func clamp01(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func easeOut(t float32) float32 {
	u := 1 - t
	return 1 - u*u*u
}
