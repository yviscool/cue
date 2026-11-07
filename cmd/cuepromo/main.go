// Command cuepromo shoots the Cue promotional film.
//
// It is a marketing build tool, not part of the shipped app. It renders the real
// HUD in a real window and drives a scripted input timeline through
// engine.HandleEvent instead of the low-level hook, because the hook
// deliberately drops injected events (LLKHF_INJECTED). Feeding the engine
// directly also makes every take deterministic: the same script always produces
// the same film.
//
// The frame is one transparent window painted in layers — backdrop, product,
// type, flash — which is the only arrangement that works: mygo paints absolute
// elements after flow content, so a backdrop drawn as an absolute sibling would
// cover the keycaps, and the keycaps drawn as flow content would be covered by
// every absolute sibling after them. See overlay.go.
//
// Usage:
//
//	PROMO_READY=ready PROMO_GO=go ./cuepromo.exe
package main

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	"unicode"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/engine"
	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/hud"
	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/semantic"
	"github.com/yviscool/cue/internal/theme"
)

var (
	stageW = envInt("PROMO_W", 1920)
	stageH = envInt("PROMO_H", 1040)
)

const (
	modeBeat = iota
	modeCard
)

// stripNames is the theme walk: every shipped theme, in the order the product
// lists them. The dots at the foot of the frame are printed from it and the
// montage visits them in the same order, so the lit dot always travels left to
// right instead of jumping around.
func stripNames() []string {
	names := make([]string, 0, len(theme.AllThemes))
	for _, t := range theme.AllThemes {
		names = append(names, t.Name)
	}
	return names
}

// film owns the overlay state and the copy that is on screen at any moment.
type film struct {
	vState  *hud.ViewState
	eng     *engine.Engine
	cfg     config.AppConfig
	scene   *scene
	overlay *mygo.Window

	strip    []string
	swatches []ui.Color // one per theme in strip

	mu     sync.Mutex
	mode   int
	cp     copySpec
	active int // index into swatches of the theme on screen, or -1
	copyAt time.Time
}

// ---- Direction ------------------------------------------------------------

// slate sets every line of copy at once. The moment is recorded so the copy can
// animate in.
func (f *film) slate(cp copySpec) {
	f.mu.Lock()
	f.cp, f.copyAt = cp, time.Now()
	f.mu.Unlock()
	f.invalidate()
}

func (f *film) showCard() {
	f.mu.Lock()
	f.mode, f.copyAt = modeCard, time.Now()
	f.mu.Unlock()
	f.invalidate()
}

// accentColor is the colour the type borrows from the theme on screen, which is
// what ties the slate and the rule under it to the keycaps beside them.
func (f *film) accentColor() ui.Color {
	return f.scene.accentNow()
}

// setTheme swaps the palette live. Existing chords re-render in the new theme on
// the next frame, which is what makes the montage possible without re-firing a
// single keystroke; the stage follows the same colour on a 700 ms dissolve.
func (f *film) setTheme(name string) {
	th := theme.ThemeByName(name)
	th.SubtitleFontSize = float32(envInt("PROMO_SUB_FONT", 28))

	f.vState.SetTheme(th)
	f.cfg.ThemePreset = name
	f.vState.SetConfig(f.cfg)
	f.eng.SetThemeName(name)
	f.scene.follow(time.Now(), th)

	active := -1
	for i, n := range f.strip {
		if n == name {
			active = i
			break
		}
	}
	f.mu.Lock()
	f.active = active
	f.mu.Unlock()

	f.invalidate()
}

func (f *film) invalidate() {
	if f.overlay != nil {
		f.overlay.Update(func() {})
	}
}

// ---- Scripted input -------------------------------------------------------

type modKey struct {
	name string
	vk   uint32
}

func modifierKeys(mods hook.Modifiers) []modKey {
	var keys []modKey
	if mods.HasCtrl() {
		keys = append(keys, modKey{"Ctrl", 0x11})
	}
	if mods.HasShift() {
		keys = append(keys, modKey{"Shift", 0x10})
	}
	if mods.HasAlt() {
		keys = append(keys, modKey{"Alt", 0x12})
	}
	if mods.HasSuper() {
		keys = append(keys, modKey{"Win", 0x5B})
	}
	return keys
}

// press types a full chord: modifiers down, key tap, key up, modifiers up.
//
// Releasing the modifiers matters. The engine clears a keycap's held flag only
// when that key's own release arrives, so a chord built from a bare keydown with
// a modifier snapshot keeps its Ctrl cap stuck in the pressed (filled) style for
// the rest of its life — a permanent highlight that reads as a bug on camera.
func (f *film) press(mods hook.Modifiers, name string, vk uint32) {
	down := modifierKeys(mods)
	for _, m := range down {
		f.eng.HandleEvent(hook.InputEvent{
			Type: hook.EventKeyDown, Timestamp: time.Now(), KeyName: m.name, VKCode: m.vk, IsModifier: true,
		})
		time.Sleep(12 * time.Millisecond)
	}

	f.eng.HandleEvent(hook.InputEvent{
		Type: hook.EventKeyDown, Timestamp: time.Now(), KeyName: name, VKCode: vk, Modifiers: mods,
	})
	time.Sleep(70 * time.Millisecond)
	f.eng.HandleEvent(hook.InputEvent{
		Type: hook.EventKeyUp, Timestamp: time.Now(), KeyName: name, VKCode: vk, Modifiers: mods,
	})
	time.Sleep(30 * time.Millisecond)

	for i := len(down) - 1; i >= 0; i-- {
		m := down[i]
		f.eng.HandleEvent(hook.InputEvent{
			Type: hook.EventKeyUp, Timestamp: time.Now(), KeyName: m.name, VKCode: m.vk, IsModifier: true,
		})
		time.Sleep(12 * time.Millisecond)
	}
}

func (f *film) typeWord(word string) {
	for _, r := range word {
		up := unicode.ToUpper(r)
		name := string(up)
		vk := uint32(up)
		f.eng.HandleEvent(hook.InputEvent{
			Type: hook.EventKeyDown, Timestamp: time.Now(), KeyName: name, VKCode: vk,
		})
		time.Sleep(70 * time.Millisecond)
		f.eng.HandleEvent(hook.InputEvent{
			Type: hook.EventKeyUp, Timestamp: time.Now(), KeyName: name, VKCode: vk,
		})
		time.Sleep(45 * time.Millisecond)
	}
}

// click presses and releases one mouse button.
func (f *film) click(btn hook.MouseButton, x, y int) {
	f.eng.HandleEvent(hook.InputEvent{
		Type: hook.EventMouseDown, Timestamp: time.Now(), Button: btn, X: x, Y: y,
	})
	time.Sleep(80 * time.Millisecond)
	f.eng.HandleEvent(hook.InputEvent{
		Type: hook.EventMouseUp, Timestamp: time.Now(), Button: btn, X: x, Y: y,
	})
}

// wheel sends one notch. The engine never lets wheel notches stack into a
// repeat badge — they are fluid motion, not repeated intent — so a run of them
// reads as a single scroll.
func (f *film) wheel(delta, x, y int) {
	f.eng.HandleEvent(hook.InputEvent{
		Type: hook.EventMouseWheel, Timestamp: time.Now(), WheelDelta: delta, X: x, Y: y,
	})
}

// ---- The film -------------------------------------------------------------

func (f *film) waitForSignal() {
	if ready := os.Getenv("PROMO_READY"); ready != "" {
		_ = os.WriteFile(ready, []byte("ready"), 0o644)
	}
	goFile := os.Getenv("PROMO_GO")
	if goFile == "" {
		time.Sleep(1500 * time.Millisecond)
		return
	}
	for i := 0; i < 1200; i++ {
		if _, err := os.Stat(goFile); err == nil {
			time.Sleep(350 * time.Millisecond)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// beat is one themed chapter: a cut, a held breath, then the input.
//
// The 700 ms of quiet before the action is deliberate. A cut that lands on the
// same frame as the keystroke gives the eye nowhere to rest; letting the empty
// stage sit for a beat makes the keycaps arriving feel like an entrance.
func (f *film) beat(themeName, ghost string, cp copySpec, run func(), hold time.Duration) {
	f.eng.Clear()
	f.setTheme(themeName)
	f.scene.setGhost(time.Now(), ghost)
	f.slate(cp)
	f.scene.flash(time.Now(), 0.07, f.accentColor())
	time.Sleep(700 * time.Millisecond)
	if run != nil {
		run()
	}
	time.Sleep(hold)
}

// shotList is the edit.
func (f *film) shotList() {
	f.waitForSignal()

	// Act 1 — cold open. Nothing but light, then the promise.
	f.slate(copySpec{})
	time.Sleep(1200 * time.Millisecond)

	f.slate(copySpec{big: "See what you press."})
	time.Sleep(1800 * time.Millisecond)

	f.slate(copySpec{
		big: "See what you press.",
		sub: "Every shortcut, spelled out the moment you hit it.",
	})
	time.Sleep(2200 * time.Millisecond)

	// Act 2 — five beats, five themes. The palette changes; the promise does
	// not.
	f.beat("Azure Pulse", "01",
		copySpec{
			tag:  "01 / 06 — AZURE PULSE",
			head: "One chord, one plain-English line",
			foot: "Ice-blue caps — modifiers and letters in one hue",
		},
		func() { f.press(hook.ModCtrl, "C", 0x43) }, 2600*time.Millisecond)

	f.beat("Cyberpunk Neon", "02",
		copySpec{
			tag:  "02 / 06 — CYBERPUNK NEON",
			head: "Three keys, still obvious",
			foot: "Violet modifiers against a cyan accent",
		},
		func() { f.press(hook.ModCtrl|hook.ModShift, "P", 0x50) }, 2600*time.Millisecond)

	f.beat("Emerald Matrix", "03",
		copySpec{
			tag:  "03 / 06 — EMERALD MATRIX",
			head: "Two stages, one gesture",
			foot: "Terminal green, both stages on one line",
		},
		func() {
			f.press(hook.ModCtrl, "K", 0x4B)
			time.Sleep(420 * time.Millisecond)
			f.press(hook.ModCtrl, "C", 0x43)
		}, 2800*time.Millisecond)

	// The repeat badge. Three presses inside the engine's 1200 ms streak window
	// fold into one cue and the last keycap grows a ×2 then ×3 multiplier.
	f.beat("Amber Cathode 1977", "04",
		copySpec{
			tag:  "04 / 06 — AMBER CATHODE 1977",
			head: "Press it again — it counts",
			foot: "Phosphor amber on warm charcoal",
		},
		func() {
			for i := 0; i < 3; i++ {
				f.press(hook.ModCtrl, "C", 0x43)
				time.Sleep(220 * time.Millisecond)
			}
		}, 2600*time.Millisecond)

	f.beat("Tokyo Blossom", "05",
		copySpec{
			tag:  "05 / 06 — TOKYO BLOSSOM",
			head: "Type fast and it groups itself",
			foot: "Sakura pink, grouped word by word",
		},
		func() { f.typeWord("hello") }, 2400*time.Millisecond)

	// Act 3 — one word, four palettes. The same five keys, typed again in each
	// theme, so the point is not "it can type" but "the look is yours".
	f.eng.Clear()
	f.scene.setGhost(time.Now(), "Aa")
	f.slate(copySpec{
		tag:  "TYPING — SAME WORD, FOUR THEMES",
		head: "Same five keys, four palettes",
		foot: "h · e · l · l · o — grouped into one word every time",
	})
	f.scene.flash(time.Now(), 0.07, f.accentColor())
	time.Sleep(600 * time.Millisecond)

	for _, name := range []string{
		"Azure Pulse", "Emerald Matrix", "Amber Cathode 1977", "Blueprint Cyanotype",
	} {
		f.eng.Clear()
		f.setTheme(name)
		f.scene.flash(time.Now(), 0.05, f.accentColor())
		time.Sleep(260 * time.Millisecond)
		f.typeWord("hello")
		time.Sleep(700 * time.Millisecond)
	}

	// Act 4 — the mouse. Left, right, wheel, and a double click, which the
	// cursor grows its own ×2 badge for.
	f.eng.Clear()
	f.scene.setGhost(time.Now(), "")
	f.setTheme("Liquid Glass")
	f.slate(copySpec{
		tag:  "MOUSE — LIQUID GLASS",
		head: "Clicks and scrolls, too",
		foot: "Left · Right · Wheel · Double-click",
	})
	f.scene.flash(time.Now(), 0.07, f.accentColor())
	time.Sleep(700 * time.Millisecond)

	// One of each, in the order a hand would do it. Exactly one wheel notch per
	// direction on purpose: the engine pushes a fresh cue for every scroll event
	// rather than folding them, so a run of notches would stack four identical
	// "Scroll Up" caps and read as a glitch.
	mx, my := stageW/2, stageH/2+210
	f.click(hook.ButtonLeft, mx, my)
	time.Sleep(700 * time.Millisecond)
	f.click(hook.ButtonRight, mx+120, my-40)
	time.Sleep(700 * time.Millisecond)
	f.wheel(120, mx+240, my-80)
	time.Sleep(500 * time.Millisecond)
	f.wheel(-120, mx+240, my-80)
	time.Sleep(900 * time.Millisecond)

	// A double click folds into the same cue and grows the same ×N badge the
	// keyboard gets, plus a badge on the cursor itself.
	f.click(hook.ButtonLeft, mx-140, my-60)
	time.Sleep(70 * time.Millisecond)
	f.click(hook.ButtonLeft, mx-140, my-60)
	time.Sleep(1500 * time.Millisecond)

	// Act 5 — the montage. One chord, every palette, no re-triggering: the chord
	// simply stays on stage while the theme underneath it changes. The flash on
	// each change is a strobe, and the stage colour chasing it is what makes a
	// colour swap read as a lighting change rather than a repaint.
	f.eng.Clear()
	f.eng.SetLingerDuration(9000)
	f.scene.setGhost(time.Now(), "19")
	f.setTheme(f.strip[0])
	f.slate(copySpec{
		tag:  "19 THEMES",
		head: "Same shortcuts, a different mood",
		foot: "Every preset ships with the binary",
	})
	f.scene.flash(time.Now(), 0.08, f.accentColor())
	time.Sleep(700 * time.Millisecond)
	f.press(hook.ModCtrl, "C", 0x43)
	time.Sleep(800 * time.Millisecond)

	for _, name := range f.strip {
		f.setTheme(name)
		f.scene.flash(time.Now(), 0.04, f.accentColor())
		time.Sleep(320 * time.Millisecond)
	}

	// Act 6 — the card.
	f.eng.Clear()
	f.scene.setGhost(time.Now(), "")
	f.slate(copySpec{})
	time.Sleep(500 * time.Millisecond)
	f.showCard()
	time.Sleep(6000 * time.Millisecond)

	// Tell the recorder the take is over, then keep the last frame on screen
	// until the capture window closes. Exiting here would drop the window and
	// leave the recorder trailing a few seconds of bare desktop, which then has
	// to be trimmed by eye.
	if done := os.Getenv("PROMO_DONE"); done != "" {
		_ = os.WriteFile(done, []byte("done"), 0o644)
		time.Sleep(120 * time.Second)
	}
	os.Exit(0)
}

// ---- Lifecycle ------------------------------------------------------------

func run() {
	f := &film{scene: &scene{}, active: -1}
	f.cfg = promoConfig()

	i18n.Instance().SetLanguage(i18n.Language(f.cfg.Language))

	f.strip = stripNames()
	for _, name := range f.strip {
		f.swatches = append(f.swatches, theme.ThemeByName(name).KeyCapHeldBorder)
	}

	f.vState = hud.NewViewState()
	f.vState.SetDisplayGeometry(0, 0, stageW, stageH, 1.0)
	f.vState.SetConfig(f.cfg)

	// Seed the aurora with the opening theme so the cold open already has its
	// colour instead of fading up from transparent black.
	f.scene.init(theme.ThemeByName(f.cfg.ThemePreset))

	appCtx := context.AppContext{
		Exe:      "code.exe",
		AppID:    "vscode",
		AppName:  "VS Code",
		Category: context.CategoryEditor,
	}

	f.eng = engine.NewEngine(
		engine.WithFilterMode(engine.FilterAll),
		engine.WithMatcher(semantic.NewDefaultMatcher()),
		engine.WithAppContext(appCtx),
		engine.WithSemanticContext(true),
		engine.WithWordGrouping(true),
		engine.WithLingerShort(3000*time.Millisecond),
		engine.WithLingerLong(3800*time.Millisecond),
		engine.WithMaxChords(4),
		engine.WithSoundActive(false),
		engine.WithShowMouseClicks(true),
		engine.WithMouseClickLinger(900*time.Millisecond),
	)
	f.eng.Start()

	f.overlay = mygo.NewWindow(mygo.WindowOptions{
		Title:         "Cue Promo",
		Panel:         false,
		Frameless:     true,
		Transparent:   true,
		AlwaysOnTop:   true,
		SkipTaskbar:   true,
		DisableResize: true,
		DisableShadow: true,
		Width:         stageW,
		Height:        stageH,
		X:             0,
		Y:             0,
		Content:       ui.View(f.render),
	})
	f.overlay.SetHasShadow(false)
	f.overlay.Show()

	sub := f.eng.Subscribe()
	go func() {
		for snap := range sub {
			f.vState.UpdateState(snap)
			f.vState.SetConfig(f.cfg)
			f.invalidate()
		}
	}()

	go f.shotList()
}

func promoConfig() config.AppConfig {
	cfg := config.DefaultConfig()

	cfg.ThemePreset = "Azure Pulse"
	// English on purpose: the film's copy is English, and the semantic subtitles
	// the engine draws under each chord come from the same localiser, so the two
	// have to agree or the frame reads as half-translated.
	cfg.Language = "en-US"
	cfg.Appearance.Animation = "spring"
	cfg.Appearance.Alignment = config.AlignCenter
	cfg.Appearance.MarginX = 0
	cfg.Appearance.MarginY = 0

	cfg.General.Filter = "none"
	cfg.General.LingerDurationMs = 3000
	cfg.General.MaxHistory = 3
	cfg.General.WordGrouping = true

	// Presentation scale: the shipped themes target a 1080p desktop where the
	// overlay occupies a corner. A promo frame is watched at a fraction of that
	// size, so the caps are drawn one step larger. Same knob a user turns with
	// the text-size setting, not a mock-up.
	cfg.Text.Size = envInt("PROMO_FONT", 62)

	cfg.Layout.ShowPressCount = true

	cfg.SemanticSub = true
	cfg.SemanticContext = true
	cfg.ShowAppBadge = true
	cfg.TactileAudio = false

	cfg.Mouse.ShowClicks = true
	cfg.Mouse.ShowIndicator = true
	cfg.Mouse.Size = 110
	cfg.Mouse.ThemePreset = "cyber_cyan"

	return cfg
}

func envInt(key string, fallback int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func main() {
	mygo.App.SetName("Cue Promo")
	mygo.App.OnWindowAllClosed(func() { os.Exit(0) })
	mygo.App.WhenReady(run)

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
