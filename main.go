package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/acoustic"
	"github.com/yviscool/cue/internal/assets"
	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/engine"
	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/hud"
	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/semantic"
	"github.com/yviscool/cue/internal/theme"
)

type app struct {
	mu         sync.Mutex
	store      *config.Store
	eng        *engine.Engine
	hk         hook.Hook
	player     *acoustic.Player
	matcher    *semantic.DefaultMatcher
	vState     *hud.ViewState
	settings   *hud.SettingsController
	hudWindow  *mygo.Window
	tray       *mygo.Tray
	hudVisible bool
	hudMu      sync.Mutex
}

// toggleHUD switches the overlay window between visible and hidden states.
func (a *app) toggleHUD() {
	a.hudMu.Lock()
	defer a.hudMu.Unlock()
	if a.hudWindow == nil {
		return
	}
	if a.hudVisible {
		a.hudVisible = false
		a.hudWindow.Hide()
		if a.eng != nil {
			a.eng.Clear()
		}
	} else {
		a.hudVisible = true
		a.hudWindow.ShowInactive()
	}
}

// isToggleShortcutMatch tests whether an input event matches the configured global toggle hotkey.
func isToggleShortcutMatch(ev hook.InputEvent, shortcutTokens []string) bool {
	if ev.Type != hook.EventKeyDown {
		return false
	}
	shortcutStr := strings.Join(shortcutTokens, "+")
	if len(shortcutTokens) == 0 {
		shortcutStr = "Shift + F10"
	}
	s := strings.ToLower(strings.TrimSpace(shortcutStr))
	if s == "" || s == "shift + f10" || s == "shift+f10" || s == "shiftleft+f10" {
		return ev.Modifiers.HasShift() && (strings.EqualFold(ev.KeyName, "F10") || ev.VKCode == 0x79)
	}

	parts := strings.Split(s, "+")
	var reqMods hook.Modifiers
	var reqKey string
	for _, p := range parts {
		token := strings.TrimSpace(p)
		switch token {
		case "ctrl", "control", "controlleft", "controlright":
			reqMods |= hook.ModCtrl
		case "shift", "shiftleft", "shiftright":
			reqMods |= hook.ModShift
		case "alt", "menu", "altleft", "altright":
			reqMods |= hook.ModAlt
		case "win", "super", "cmd", "metaleft", "metaright":
			reqMods |= hook.ModSuper
		default:
			reqKey = token
		}
	}
	if ev.Modifiers != reqMods {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(ev.KeyName), reqKey)
}

// setFilter is the only write path for the tray capture-mode items. It persists
// the choice and lets Store.OnChange drive the engine through applyConfig.
func (a *app) setFilter(filter string) {
	if a.store == nil {
		return
	}
	if err := a.store.Update(func(c *config.AppConfig) { c.General.Filter = filter }); err != nil {
		log.Printf("save filter %q: %v", filter, err)
	}
}

func (a *app) refreshTrayMenu() {
	if a.tray != nil {
		a.tray.SetMenu(a.buildTrayMenu())
	}
}

// reportConfigError keeps the last good configuration active and tells the user
// through the tray tooltip, where they will look when a manual edit appears to
// do nothing.
func (a *app) reportConfigError(err error) {
	log.Printf("config reload rejected, keeping last good config: %v", err)
	if a.tray == nil {
		return
	}
	a.tray.SetToolTip("Cue · 配置有误，已保留上次设置")
	time.AfterFunc(6*time.Second, func() {
		if a.tray != nil {
			a.tray.SetToolTip("Cue · 演示与舞台信号")
		}
	})
}

// buildTrayMenu constructs the context menu dynamically with active locale strings.
func (a *app) buildTrayMenu() *mygo.Menu {
	curLang := i18n.Instance().GetLanguage()
	curFilter := ""
	if a.store != nil {
		curFilter = a.store.Get().General.Filter
	}

	var langItems []*mygo.MenuItem
	for _, l := range i18n.SupportedLanguages {
		langCode := l.Code
		label := l.Name
		if curLang == langCode {
			label = "✓ " + label
		}
		targetLang := langCode
		langItems = append(langItems, &mygo.MenuItem{
			Label: label,
			Click: func(*mygo.MenuItem, *mygo.Window) {
				i18n.Instance().SetLanguage(targetLang)
				_ = a.store.Update(func(c *config.AppConfig) { c.Language = string(targetLang) })
				if a.tray != nil {
					a.tray.SetMenu(a.buildTrayMenu())
				}
			},
		})
	}

	items := []*mygo.MenuItem{
		{
			Label: i18n.T("tray_settings"),
			Click: func(*mygo.MenuItem, *mygo.Window) { a.settings.Open() },
		},
		mygo.Separator(),
		{
			ID:      "filter-smart",
			Label:   i18n.T("tray_smart"),
			Type:    mygo.MenuItemRadio,
			Checked: curFilter == "modifiers" || curFilter == "",
			Click: func(*mygo.MenuItem, *mygo.Window) {
				a.setFilter("modifiers")
			},
		},
		{
			ID:      "filter-all",
			Label:   i18n.T("tray_all"),
			Type:    mygo.MenuItemRadio,
			Checked: curFilter == "none",
			Click: func(*mygo.MenuItem, *mygo.Window) {
				a.setFilter("none")
			},
		},
		{
			ID:      "filter-silent",
			Label:   i18n.T("tray_silent"),
			Type:    mygo.MenuItemRadio,
			Checked: curFilter == "silent",
			Click: func(*mygo.MenuItem, *mygo.Window) {
				a.setFilter("silent")
			},
		},
		mygo.Separator(),
		{
			Label:   i18n.T("language"),
			Submenu: langItems,
		},
		mygo.Separator(),
		{
			Label: i18n.T("tray_clear"),
			Click: func(*mygo.MenuItem, *mygo.Window) {
				a.eng.Clear()
			},
		},
		mygo.Separator(),
		{
			Label: i18n.T("tray_quit"),
			Click: func(*mygo.MenuItem, *mygo.Window) {
				mygo.App.Quit()
			},
		},
	}
	return mygo.NewMenu(items)
}

func (a *app) ready() {
	var err error
	a.store, err = config.NewStore("")
	if err != nil {
		log.Printf("store init warning: %v", err)
	}
	cfg := a.store.Get()

	// Apply user preferred language if configured
	if cfg.Language != "" {
		i18n.Instance().SetLanguage(i18n.Language(cfg.Language))
	}

	// Set application icon for all windows & taskbar in pure Go
	_ = mygo.App.Dock.SetIcon(assets.AppIconPNG)

	a.vState = hud.NewViewState()
	a.player = acoustic.NewPlayer()
	a.matcher = semantic.NewDefaultMatcher()

	soundAdapter := engine.SpatialSoundPlayerFunc(func(s engine.SoundType, key string) {
		st := acoustic.SoundThock
		if s == engine.SoundClick {
			st = acoustic.SoundClick
		}
		a.player.PlaySpatial(st, key)
	})

	a.eng = engine.NewEngine(
		engine.WithFilterMode(filterModeFor(cfg.General.Filter)),
		engine.WithAllowedKeys(cfg.General.AllowedKeys),
		engine.WithSoundPlayer(soundAdapter),
		engine.WithSoundActive(cfg.TactileAudio),
		engine.WithMatcher(a.matcher),
		engine.WithThemeName(cfg.ThemePreset),
		engine.WithShowMouseClicks(cfg.Mouse.ShowClicks),
		engine.WithMouseClickLinger(time.Duration(cfg.Mouse.ClickLingerMs)*time.Millisecond),
		engine.WithDragThreshold(cfg.General.DragThreshold),
		engine.WithWordGrouping(cfg.General.WordGrouping),
		engine.WithLingerShort(time.Duration(cfg.General.LingerDurationMs)*time.Millisecond),
		engine.WithMaxChords(historyChordLimit(cfg)),
	)
	a.applyConfig(cfg)

	a.settings = hud.NewSettingsController(a.store, a.eng, a.vState)

	a.hudVisible = true

	// 1. Setup HUD overlay window with full screen bounds for Cue overlay & mouse tracking
	d, _ := resolveDisplay(mygo.Screen.Displays(), cfg.Appearance.Monitor)
	if d.Bounds.Width <= 0 || d.Bounds.Height <= 0 {
		d = mygo.Screen.PrimaryDisplay()
	}
	winW := d.Bounds.Width
	winH := d.Bounds.Height
	posX := d.Bounds.X
	posY := d.Bounds.Y
	if winW <= 0 || winH <= 0 {
		winW, winH = 1920, 1080
	}
	a.vState.SetDisplayGeometry(posX, posY, winW, winH, d.ScaleFactor)

	a.hudWindow = mygo.NewWindow(mygo.WindowOptions{
		Title:         "Cue HUD",
		Panel:         false,
		Frameless:     true,
		Transparent:   true,
		AlwaysOnTop:   true,
		SkipTaskbar:   true,
		DisableResize: true,
		DisableShadow: true,
		Width:         winW,
		Height:        winH,
		X:             posX,
		Y:             posY,
		Content:       ui.View(a.vState.Render),
	})

	// Make overlay pass-through to clicks and disable any DWM border or shadow
	a.hudWindow.SetHasShadow(false)
	a.hudWindow.SetIgnoreMouseEvents(true)
	disableWindowBorder(a.hudWindow.NativeHandle())
	a.hudWindow.ShowInactive()

	// 2. Dynamic config change listener: one apply path keeps the HUD, engine,
	// theme, and overlay geometry in sync with every settings change.
	a.store.OnChange(func(updated config.AppConfig) {
		a.applyConfig(updated)
		a.refreshTrayMenu()
		if a.hudWindow != nil {
			a.hudWindow.Update(func() {})
		}
	})
	a.store.OnError(a.reportConfigError)
	a.store.Watch(400 * time.Millisecond)

	// 3. Subscribe HUD view to state engine
	sub := a.eng.Subscribe()
	go func() {
		for snap := range sub {
			a.vState.UpdateState(snap)
			a.vState.SetConfig(a.store.Get())
			a.hudWindow.Update(func() {})
		}
	}()

	// 3. Setup System Tray with official Cue icon
	a.tray, err = mygo.NewTray(mygo.TrayOptions{
		ID:             "cue-app",
		Icon:           assets.TrayIconPNG,
		IconIsTemplate: true,
		ToolTip:        "Cue · 演示与舞台信号",
		Menu:           a.buildTrayMenu(),
	})
	if err != nil {
		log.Printf("tray warning: %v", err)
	} else {
		// Single or double click directly opens Settings Studio (just like cloudmusic)
		a.tray.OnClick(func() {
			a.settings.Open()
		})
		if msg := a.store.LastError(); msg != "" {
			a.reportConfigError(fmt.Errorf("%s", msg))
		}
	}

	// Keep the overlay on the selected display when monitors are added, removed,
	// rearranged, or their resolution / DPI changes.
	mygo.Screen.OnDisplaysChanged(func() {
		a.applyConfig(a.store.Get())
		if a.settings != nil {
			a.settings.Invalidate()
		}
	})

	// 4. Start OS Low-Level Hook with Global Toggle Shortcut support
	a.hk = hook.NewHook()
	go func() {
		for ev := range a.hk.Events() {
			if isToggleShortcutMatch(ev, a.store.Get().General.ToggleShortcut) {
				a.toggleHUD()
				continue
			}
			a.eng.HandleEvent(ev)
		}
	}()

	if err := a.hk.Start(); err != nil {
		log.Printf("hook start error: %v", err)
	}

	// 5. Clean shutdown
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) {
		if a.store != nil {
			a.store.StopWatching()
		}
		if a.hk != nil {
			a.hk.Stop()
		}
		if a.eng != nil {
			a.eng.Stop()
		}
		if a.player != nil {
			a.player.Close()
		}
		if a.tray != nil {
			a.tray.Destroy()
		}
	})
}

// filterModeFor maps the persisted filter policy to the engine mode.
func filterModeFor(filter string) engine.FilterMode {
	switch filter {
	case "none":
		return engine.FilterAll
	case "custom":
		return engine.FilterCustom
	case "silent":
		return engine.FilterSilent
	default:
		return engine.FilterSmart
	}
}

// historyChordLimit turns "show event history" into a chord limit. With history
// off the HUD is intentionally a single-cue stage: one current chord at a time.
func historyChordLimit(cfg config.AppConfig) int {
	if !cfg.General.ShowEventHistory {
		return 1
	}
	n := cfg.General.MaxHistory
	if n < 2 {
		n = 2
	}
	if n > 12 {
		n = 12
	}
	return n
}

// applyConfig is the single synchronization point for live configuration. Both
// startup and Store.OnChange use it, so a setting cannot be wired to one path
// and silently ignored by the other.
func (a *app) applyConfig(cfg config.AppConfig) {
	if a.vState != nil {
		a.vState.SetTheme(theme.ThemeByName(cfg.ThemePreset))
		a.vState.SetConfig(cfg)
	}
	a.repositionHUD(cfg)
	if a.player != nil {
		a.player.SetEnabled(cfg.TactileAudio)
		prof := cfg.AudioProfile
		if prof == "" {
			prof = "auto"
		}
		a.player.SetProfile(acoustic.SoundProfile(prof))
		thm := theme.ThemeByName(cfg.ThemePreset)
		a.player.SetThemeProfile(thm.AcousticProfile)
	}
	if a.eng == nil {
		return
	}
	a.eng.SetFilterMode(filterModeFor(cfg.General.Filter))
	a.eng.SetAllowedKeys(cfg.General.AllowedKeys)
	a.eng.SetShowMouseClicks(cfg.Mouse.ShowClicks)
	a.eng.SetMouseClickLinger(cfg.Mouse.ClickLingerMs)
	a.eng.SetDragThreshold(cfg.General.DragThreshold)
	a.eng.SetSoundActive(cfg.TactileAudio)
	a.eng.SetWordGrouping(cfg.General.WordGrouping)
	a.eng.SetLingerDuration(cfg.General.LingerDurationMs)
	a.eng.SetMaxChords(historyChordLimit(cfg))
	a.eng.SetThemeName(cfg.ThemePreset)
}

func (a *app) repositionHUD(cfg config.AppConfig) {
	if a.hudWindow == nil {
		return
	}

	d, ok := resolveDisplay(mygo.Screen.Displays(), cfg.Appearance.Monitor)
	if !ok || d.Bounds.Width <= 0 || d.Bounds.Height <= 0 {
		d = mygo.Screen.PrimaryDisplay()
	}

	bounds := d.Bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		bounds = mygo.Rectangle{Width: 1920, Height: 1080}
	}

	// Remember the display-local coordinate space so mouse feedback is drawn
	// under the real cursor on any monitor, including negative origins.
	a.vState.SetDisplayGeometry(bounds.X, bounds.Y, bounds.Width, bounds.Height, d.ScaleFactor)

	if a.hudWindow.Bounds() != bounds {
		// SetBounds moves and resizes in one step, avoiding intermediate states
		// where a mixed-DPI window is briefly sized for the old monitor.
		a.hudWindow.SetBounds(bounds)
	}
}

func disableWindowBorder(hwnd uintptr) {
	if runtime.GOOS != "windows" || hwnd == 0 {
		return
	}
	dwmapi := syscall.NewLazyDLL("dwmapi.dll")
	user32 := syscall.NewLazyDLL("user32.dll")

	procDwmSetWindowAttribute := dwmapi.NewProc("DwmSetWindowAttribute")
	if procDwmSetWindowAttribute != nil {
		// DWMWA_WINDOW_CORNER_PREFERENCE = 33, DWMWCP_DONOTROUND = 1
		pref := int32(1)
		procDwmSetWindowAttribute.Call(hwnd, 33, uintptr(unsafe.Pointer(&pref)), 4)

		// DWMWA_BORDER_COLOR = 34, DWMWA_COLOR_NONE = 0xFFFFFFFE (prevents DWM drawing any border)
		borderColor := uint32(0xFFFFFFFE)
		procDwmSetWindowAttribute.Call(hwnd, 34, uintptr(unsafe.Pointer(&borderColor)), 4)
	}

	procGetWindowLong := user32.NewProc("GetWindowLongW")
	procSetWindowLong := user32.NewProc("SetWindowLongW")
	procSetWindowPos := user32.NewProc("SetWindowPos")

	if procGetWindowLong != nil && procSetWindowLong != nil && procSetWindowPos != nil {
		const (
			wsBorder        = 0x00800000
			wsThickFrame    = 0x00040000
			wsCaption       = 0x00C00000
			wsExLayered     = 0x00080000
			wsExTransparent = 0x00000020
			wsExToolWindow  = 0x00000080
			wsExTopmost     = 0x00000008
			wsExNoActivate  = 0x08000000
			swpNoMove       = 0x0002
			swpNoSize       = 0x0001
			swpNoZOrder     = 0x0004
			swpFrameChanged = 0x0020
		)
		n16 := int32(-16)
		n20 := int32(-20)
		gwlStyle := uintptr(uint32(n16))
		gwlExStyle := uintptr(uint32(n20))

		style, _, _ := procGetWindowLong.Call(hwnd, gwlStyle)
		newStyle := style &^ uintptr(wsBorder|wsThickFrame|wsCaption)
		procSetWindowLong.Call(hwnd, gwlStyle, newStyle)

		exStyle, _, _ := procGetWindowLong.Call(hwnd, gwlExStyle)
		newExStyle := exStyle | uintptr(wsExLayered|wsExTransparent|wsExToolWindow|wsExTopmost|wsExNoActivate)
		procSetWindowLong.Call(hwnd, gwlExStyle, newExStyle)

		procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, uintptr(swpNoMove|swpNoSize|swpNoZOrder|swpFrameChanged))
	}
}

func hideConsoleOnWindows() {
	if runtime.GOOS != "windows" {
		return
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	procGetConsole := kernel32.NewProc("GetConsoleWindow")
	procShowWindow := user32.NewProc("ShowWindow")
	if procGetConsole != nil && procShowWindow != nil {
		hwnd, _, _ := procGetConsole.Call()
		if hwnd != 0 {
			const swHide = 0
			procShowWindow.Call(hwnd, uintptr(swHide))
		}
	}
}

func main() {
	hideConsoleOnWindows()
	a := &app{}
	mygo.App.SetName("Cue")
	mygo.App.OnWindowAllClosed(func() {})
	mygo.App.WhenReady(a.ready)

	err := mygo.App.Run()
	log.Printf("mygo.App.Run() returned: %v", err)
}
