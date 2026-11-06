package hud

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/assets"
	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/engine"
	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/semantic"
	"github.com/yviscool/cue/internal/theme"
)

// strokeIcon parses an inline SVG path with standard 1.75 stroke.
func strokeIcon(paths string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">` + paths + `</svg>`))
}

var (
	iconSettings    = strokeIcon(`<path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/>`)
	iconMonitor     = strokeIcon(`<rect width="20" height="14" x="2" y="3" rx="2"/><line x1="8" x2="16" y1="21" y2="21"/><line x1="12" x2="12" y1="17" y2="21"/>`)
	iconKeyboard    = strokeIcon(`<rect width="20" height="16" x="2" y="4" rx="2"/><path d="M6 8h.01M10 8h.01M14 8h.01M18 8h.01M8 12h.01M12 12h.01M16 12h.01M7 16h10"/>`)
	iconMouse       = strokeIcon(`<rect x="5" y="2" width="14" height="20" rx="7"/><path d="M12 6v4"/>`)
	iconInfo        = strokeIcon(`<circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/>`)
	iconReset       = strokeIcon(`<path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/>`)
	iconSearch      = strokeIcon(`<circle cx="11" cy="11" r="7"/><path d="m20 20-3.5-3.5"/>`)
	iconChevronDown = strokeIcon(`<path d="m6 9 6 6 6-6"/>`)
)

// themeDisplayName returns an intuitive bilingual Chinese-first display name for preset themes.
func themeDisplayName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "liquid glass":
		return "液态冰晶 (Liquid Glass)"
	case "macbook scissor":
		return "深空剪刀脚 (MacBook Scissor)"
	case "dieter rams 1960":
		return "博朗经典 (Dieter Rams)"
	case "cyberpunk neon":
		return "赛博霓虹 (Cyberpunk Neon)"
	case "classic mechanical":
		return "复古机械 (Classic Mechanical)"
	case "ibm model m 1984":
		return "传奇屈曲轴 (IBM Model M)"
	case "minimal chalk":
		return "极简粉笔白 (Minimal Chalk)"
	case "obsidian oled":
		return "黑曜极光 (Obsidian OLED)"
	case "tokyo blossom":
		return "东京落樱 (Tokyo Blossom)"
	case "emerald matrix":
		return "代码矩阵 (Emerald Matrix)"
	case "sunset horizon":
		return "暮光地平线 (Sunset Horizon)"
	case "royal velvet":
		return "皇家丝绒 (Royal Velvet)"
	case "sony walkman 1979":
		return "索尼经典 (Sony Walkman 1979)"
	case "teenage op-1":
		return "合成器先锋 (Teenage OP-1)"
	case "amber cathode 1977":
		return "琥珀终端 (Amber Cathode 1977)"
	case "swiss typography 1957":
		return "苏黎世平面 (Swiss Typography 1957)"
	case "kyoto kintsugi":
		return "京都金缮 (Kyoto Kintsugi)"
	case "blueprint cyanotype":
		return "工程蓝晒 (Blueprint Cyanotype)"
	default:
		th := theme.ThemeByName(name)
		if th.DisplayName != "" {
			return th.DisplayName
		}
		return name
	}
}

// physicalStyleName returns a Chinese-first name for a keycap physical model.
func physicalStyleName(id string) string {
	switch id {
	case "minimal":
		return "极简扁平"
	case "laptop":
		return "剪刀脚"
	case "pbt":
		return "PBT 机械"
	case "glass":
		return "拟物玻璃"
	case "neon":
		return "霓虹发光"
	case "retro":
		return "复古终端"
	case "outline":
		return "极简线框"
	default:
		return "矮轴物理"
	}
}

// travelDepthLabel returns human-readable physical stroke depth or kinetic mode.
func travelDepthLabel(t theme.ThemeConfig) string {
	if t.HeldPhysics == theme.HeldPulse {
		return "光压脉冲"
	}
	return fmt.Sprintf("%.1fmm 轴体行程", t.EffectiveTravelDepth())
}

// styleStageLabel describes whether the physical model is inherited from the archetype or overridden.
func styleStageLabel(cfg config.AppConfig, t theme.ThemeConfig) string {
	if cfg.Appearance.Style == "" {
		return fmt.Sprintf("%s · 原型预设 (%s)", physicalStyleName(effectiveStyle(t, cfg)), travelDepthLabel(t))
	}
	return physicalStyleName(cfg.Appearance.Style) + " · 自定义覆盖"
}

// formatShortcut renders the real configured global hotkey rather than a
// hard-coded "default" label that could contradict the saved configuration.
func formatShortcut(tokens []string) string {
	if len(tokens) == 0 {
		return "Shift + F10"
	}
	out := make([]string, 0, len(tokens))
	for _, raw := range tokens {
		switch strings.ToLower(strings.TrimSpace(raw)) {
		case "shiftleft", "shiftright", "shift":
			out = append(out, "Shift")
		case "controlleft", "controlright", "ctrl", "control":
			out = append(out, "Ctrl")
		case "altleft", "altright", "alt", "menu":
			out = append(out, "Alt")
		case "metaleft", "metaright", "win", "super", "cmd":
			out = append(out, "Win")
		default:
			out = append(out, raw)
		}
	}
	return strings.Join(out, " + ")
}

// Palette holds Cupertino / Apple Human Interface theme colors.
type palette struct {
	dark      bool
	bg        ui.Color // card background
	card      ui.Color
	sidebar   ui.Color
	canvas    ui.Color
	text      ui.Color
	muted     ui.Color
	faint     ui.Color
	hairline  ui.Color
	hover     ui.Color
	press     ui.Color
	accent    ui.Color
	onAccent  ui.Color
	cardShade ui.Color
}

func paletteFor(t *ui.Theme) palette {
	if t.Dark {
		cardColor := ui.Hex("#242428")
		return palette{
			dark:      true,
			bg:        cardColor,
			card:      cardColor,
			sidebar:   ui.Hex("#131316"),
			canvas:    ui.Hex("#1c1c1f"),
			text:      ui.Hex("#f5f5f7"),
			muted:     ui.RGBA(235, 235, 245, 0.60),
			faint:     ui.RGBA(235, 235, 245, 0.28),
			hairline:  ui.RGBA(255, 255, 255, 0.08),
			hover:     ui.RGBA(255, 255, 255, 0.06),
			press:     ui.RGBA(255, 255, 255, 0.12),
			accent:    ui.Hex("#0a84ff"),
			onAccent:  ui.Hex("#ffffff"),
			cardShade: ui.RGBA(0, 0, 0, 0.45),
		}
	}
	cardColor := ui.Hex("#ffffff")
	return palette{
		bg:        cardColor,
		card:      cardColor,
		sidebar:   ui.Hex("#ebecef"),
		canvas:    ui.Hex("#f6f7f9"),
		text:      ui.Hex("#1d1d1f"),
		muted:     ui.RGBA(60, 60, 67, 0.65),
		faint:     ui.RGBA(60, 60, 67, 0.28),
		hairline:  ui.RGBA(0, 0, 0, 0.08),
		hover:     ui.RGBA(0, 0, 0, 0.04),
		press:     ui.RGBA(0, 0, 0, 0.08),
		accent:    ui.Hex("#007aff"),
		onAccent:  ui.Hex("#ffffff"),
		cardShade: ui.RGBA(0, 0, 0, 0.06),
	}
}

// MasterColorSchemes holds original 14 master color schemes from cue.
type ColorSchemePreset struct {
	Name      string
	Primary   string
	Secondary string
	Text      string
}

var MasterColorSchemes = []ColorSchemePreset{
	{"Silver", "#f8f8f8", "#dcdcdc", "#000000"},
	{"Stone", "#606060", "#4b4b4b", "#f8f8f8"},
	{"Lime", "#606060", "#4b4b4b", "#D6ED17"},
	{"Cyber", "#00B1D2", "#008ea8", "#FDDB27"},
	{"Turquoise", "#42EADD", "#2ec4b8", "#ffffff"},
	{"Blue", "#2196f3", "#1976d2", "#ffffff"},
	{"Yellow", "#FDDB27", "#dfc019", "#000000"},
	{"Green", "#66bb6a", "#43a047", "#ffffff"},
	{"Pink", "#f06292", "#d81b60", "#ffffff"},
	{"Red", "#ef5350", "#c62828", "#ffffff"},
	{"Pansy", "#673ab7", "#4527a0", "#ffc107"},
	{"Eclipse", "#343148", "#252333", "#D7C49E"},
	{"Bumblebee", "#404040", "#2e2e2e", "#FDDB27"},
	{"Charcoal", "#404040", "#2e2e2e", "#FFFFFF"},
}

// SettingsController drives the full-fledged Cue Settings Client studio window.
type SettingsController struct {
	store                 *config.Store
	engine                *engine.Engine
	vState                *ViewState
	activeTab             string
	window                *mygo.Window
	logoSrc               ui.ImageSource
	marginLinked          bool
	scrollState           ui.ScrollState
	showResetConfirm      bool
	showCustomFilterModal bool
	themeDropdownOpen     bool
	styleDropdownOpen     bool
	colorDropdownOpen     bool
	filterTab             string
	previewNonce          int
}

// NewSettingsController creates a controller for the settings client window.
func NewSettingsController(store *config.Store, eng *engine.Engine, vState *ViewState) *SettingsController {
	logo, _ := assets.GetAppLogoSource()
	if logo == nil {
		logo, _ = assets.GetAppIconSource()
	}
	sc := &SettingsController{
		store:        store,
		engine:       eng,
		vState:       vState,
		activeTab:    "General",
		logoSrc:      logo,
		marginLinked: true,
		filterTab:    "Keyboard",
	}
	if store != nil {
		store.OnChange(func(config.AppConfig) {
			sc.Invalidate()
		})
	}
	return sc
}

// updateConfig updates configuration in store and automatically invalidates window.
func (sc *SettingsController) updateConfig(fn func(*config.AppConfig)) {
	err := sc.store.Update(fn)
	_ = err
	sc.Invalidate()
}

// Invalidate repaints the settings window when it is open.
func (sc *SettingsController) Invalidate() {
	if sc.window != nil {
		sc.window.Invalidate()
	}
}

// Open shows or focuses the settings window (100% thread-safe & deadlock-free).
func (sc *SettingsController) Open() {
	if sc.window != nil {
		if sc.window.IsMinimized() {
			sc.window.Restore()
		}
		sc.window.Show()
		sc.window.Focus()
		mygo.App.Focus()
		return
	}

	sc.window = mygo.NewWindow(mygo.WindowOptions{
		Title:           "Cue — Settings",
		Width:           960,
		Height:          720,
		MinWidth:        800,
		MinHeight:       560,
		StateKey:        "settings",
		DisableMaximize: true,
		AutoHideMenuBar: true,
		BackgroundColor: "light-dark(#f6f7f9, #1c1c1f)",
		Content:         ui.View(sc.render),
	})
	_ = sc.window.SetIcon(assets.AppIconPNG)

	sc.window.OnClosed(func() {
		sc.window = nil
	})
}

// render renders the entire settings studio client matching mygo-cloudmusic design language.
func (sc *SettingsController) render(c *ui.Context) {
	t := c.Theme()
	p := paletteFor(t)
	c.Root().Background(p.canvas)

	activeTab := sc.activeTab
	cfg := sc.store.Get()
	logo := sc.logoSrc

	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		// ───────────── LEFT SIDEBAR (Apple Ventura/Jobsian Aesthetic, 1/5 narrower) ─────────────
		ui.Column(c).Width(208).Background(p.sidebar).Padding(16, 10, 14, 10).Gap(3).BorderWidth(0, 1, 0, 0).BorderColor(p.hairline).Children(func() {
			// 1. Brand Header with App Icon & Version
			ui.Row(c).AlignItems(ui.Center).Gap(10).Padding(4, 6, 10, 6).Children(func() {
				if logo != nil {
					ui.Image(c, logo).Size(36, 36).Radius(8)
				} else {
					ui.Box(c).Size(36, 36).Radius(8).Background(p.accent).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
						ui.Icon(c, iconKeyboard).Size(20, 20).TextColor(p.onAccent)
					})
				}
				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, i18n.T("app_name")).Bold().FontSize(14.5).TextColor(p.text)
					ui.Text(c, "v1.0 · Direct3D 11").FontSize(11).TextColor(p.muted)
				})
			})

			ui.Box(c).Height(1).Background(p.hairline).Margin(2, 4, 8, 4)

			// Navigation Tabs
			navTab := func(tabID, labelKey string, icon *ui.SVG) {
				isSelected := activeTab == tabID
				b := ui.ButtonBase(c.Key("nav-" + tabID)).
					Height(36).
					PaddingX(10).
					Gap(8).
					Radius(7).
					AlignItems(ui.Center).
					Justify(ui.Start).
					FocusRing(false)

				bg := ui.Transparent
				switch {
				case isSelected:
					bg = p.press
				case b.Hovered():
					bg = p.hover
				}
				b.Background(b.AnimateColor("bg", bg, ui.Animation{Duration: 100 * time.Millisecond}))

				iconColor := p.muted
				textColor := p.text
				if isSelected {
					iconColor = p.accent
				}

				b.Children(func() {
					if isSelected {
						ui.Box(c).Width(3).Height(18).Radius(1.5).Background(p.accent)
					} else {
						ui.Box(c).Width(3).Height(18)
					}
					ui.Icon(c, icon).Size(16, 16).TextColor(iconColor)
					ui.Text(c, i18n.T(labelKey)).FontSize(13).FontWeight(func() int {
						if isSelected {
							return 600
						}
						return 400
					}()).TextColor(textColor).Grow(1)
				})

				if b.Clicked() {
					sc.activeTab = tabID
					sc.themeDropdownOpen = false
					sc.styleDropdownOpen = false
					sc.colorDropdownOpen = false
					if sc.window != nil {
						sc.window.Invalidate()
					}
				}
			}

			navTab("General", "tab_general", iconSettings)
			navTab("Appearance", "tab_appearance", iconMonitor)
			navTab("Keycap", "tab_keycap", iconKeyboard)
			navTab("Mouse", "tab_mouse", iconMouse)

			ui.Spacer(c)

			navTab("About", "tab_about", iconInfo)

			ui.Box(c).Height(1).Background(p.hairline).Margin(2, 4, 4, 4)

			// Revert to Default configuration
			resetBtn := ui.ButtonBase(c.Key("reset-defaults")).
				Height(30).
				PaddingX(8).
				Gap(6).
				Radius(6).
				AlignItems(ui.Center).
				Justify(ui.Start).
				FocusRing(false)

			bgReset := ui.Transparent
			if resetBtn.Hovered() {
				bgReset = p.hover
			}
			resetBtn.Background(resetBtn.AnimateColor("bgReset", bgReset, ui.Animation{Duration: 100 * time.Millisecond}))

			resetBtn.Children(func() {
				ui.Icon(c, iconReset).Size(13, 13).TextColor(p.muted)
				ui.Text(c, i18n.T("reset_defaults")).FontSize(11.5).TextColor(p.muted)
			})

			if resetBtn.Clicked() {
				sc.showResetConfirm = true
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})

		ui.Box(c).Width(1).Background(p.hairline)

		// ───────────── RIGHT STAGE & CONTENT (Persistent Live Showcase Stage + Config Studio) ─────────────
		ui.Column(c.Key("main-content-area")).Grow(1).Background(p.canvas).Children(func() {
			// Persistent Live Showcase Stage (Always pinned at top, never scrolls away for General, Appearance, Keycap, Mouse)
			if activeTab != "About" {
				ui.Box(c.Key("persistent-showcase-container")).
					Padding(20, 32, 6, 32).
					Children(func() {
						sc.renderShowcaseStage(c, cfg, p)
					})
			}

			// Scrollable Settings Content Area
			ui.Scroll(c.Key("settings-scroll-"+activeTab)).
				Grow(1).
				Scrollbars(true).
				Padding(10, 32, 48, 32).
				Gap(14).
				Children(func() {
					// group renders a cloudmusic-style card group
					group := func(title string, rows func()) {
						ui.Text(c, title).FontSize(12).FontWeight(600).TextColor(p.muted).Margin(14, 4, 6, 4)
						ui.Column(c).Radius(12).Background(p.card).Border(1, p.hairline).Dividers(1, p.hairline).Children(rows)
					}

					// row renders a standard settings line
					row := func(label, desc string, control func()) {
						ui.Row(c).Padding(12, 16).Gap(16).AlignItems(ui.Center).MinHeight(48).Children(func() {
							ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Gap(3).Children(func() {
								ui.Text(c, label).FontSize(13.5).FontWeight(500).TextColor(p.text)
								if desc != "" {
									ui.Text(c, desc).FontSize(12).TextColor(p.muted).LineHeight(1.4)
								}
							})
							control()
						})
					}

					// stacked renders a full-width item (e.g. alignment grid or swatches)
					stacked := func(label, desc string, control func()) {
						ui.Column(c).Padding(14, 16).Gap(10).Children(func() {
							ui.Column(c).Gap(3).Children(func() {
								ui.Text(c, label).FontSize(13.5).FontWeight(500).TextColor(p.text)
								if desc != "" {
									ui.Text(c, desc).FontSize(12).TextColor(p.muted).LineHeight(1.4)
								}
							})
							control()
						})
					}

					switch activeTab {
					case "General":
						sc.renderGeneralTab(c, cfg, p, group, row)
					case "Appearance":
						sc.renderAppearanceTab(c, cfg, p, group, row, stacked)
					case "Keycap":
						sc.renderKeycapTab(c, cfg, p, group, row, stacked)
					case "Mouse":
						sc.renderMouseTab(c, cfg, p, group, row)
					case "About":
						sc.renderAboutTab(c, p, group, row)
					}
				})
		})
	})

	// Modal Dialog Overlays
	if sc.showResetConfirm {
		sc.renderResetConfirmDialog(c, p)
	}
	if sc.showCustomFilterModal {
		sc.renderCustomFilterModal(c, cfg, p)
	}
}

// renderShowcaseStage renders an Apple Cupertino style live preview stage,
// displaying realistic active HUD keycaps ([Ctrl] + [C] "Copy" and [⌘] + [V] "Paste")
// in the live active theme, physical model, typography, and color palette.
// Zero textual noise, zero double-frame nesting; clicking the stage replays the spring transition.
func (sc *SettingsController) renderShowcaseStage(c *ui.Context, cfg config.AppConfig, p palette) {
	currentTheme := theme.ThemeByName(cfg.ThemePreset)

	stage := ui.Box(c.Key("showcase-canvas")).
		Fill().
		Height(92).
		Radius(10).
		Background(ui.Hex("#080c14")). // deep obsidian simulated desktop backdrop
		Border(1, p.hairline).
		Padding(12).
		AlignItems(ui.Center).
		Justify(ui.Center).
		Shadow(0, 4, 16, 0, ui.RGBA(0, 0, 0, 0.35))

	if stage.Clicked() {
		sc.previewNonce++
		if sc.window != nil {
			sc.window.Invalidate()
		}
	}

	stage.Children(func() {
		ui.Row(c).Gap(28).AlignItems(ui.Center).Children(func() {
			// Sample Chord 1: [ Ctrl ] + [ C ] -> [ 复制 Copy ]
			sampleChord1 := engine.DisplayChord{
				ID: fmt.Sprintf("preview-%d-1", sc.previewNonce),
				Keys: []engine.KeyCap{
					{Label: "Ctrl", IsModifier: true},
					{Label: "C", IsModifier: false},
				},
				Action: semantic.ActionInfo{Title: "Copy"},
			}
			renderChord(c, sampleChord1, currentTheme, cfg)

			// Sample Chord 2: [ ⌘ Cmd ] + [ V ] -> [ 粘贴 Paste ]
			sampleChord2 := engine.DisplayChord{
				ID: fmt.Sprintf("preview-%d-2", sc.previewNonce),
				Keys: []engine.KeyCap{
					{Label: "⌘", IsModifier: true},
					{Label: "V", IsModifier: false},
				},
				Action: semantic.ActionInfo{Title: "Paste"},
			}
			renderChord(c, sampleChord2, currentTheme, cfg)
		})
	})
}

// renderGeneralTab mirrors Cue General settings with clean card groups and 12 preset themes.
func (sc *SettingsController) renderGeneralTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("general_title")).FontSize(20).Bold().TextColor(p.text)

	// Group 1: 预设主题 (Preset Themes Floating Web Dropdown)
	group("预设主题", func() {
		currentTheme := theme.ThemeByName(cfg.ThemePreset)
		currentDisplay := themeDisplayName(currentTheme.Name)
		themeChoice := cfg.ThemePreset

		row("预设原型", "设计风格与物理微几何方案", func() {
			sel := ui.SelectBase(c, &themeChoice)
			sel.Trigger.
				Height(36).
				PaddingX(12).
				Gap(10).
				Radius(8).
				Background(p.sidebar).
				Border(1, p.hairline).
				AlignItems(ui.Center).
				Children(func() {
					ui.Box(c).Size(22, 20).Radius(5).
						Background(currentTheme.KeyCapBackground).
						Border(1.2, currentTheme.KeyCapBorder).
						Shadow(0, 2, 6, 0, ui.RGBA(0, 0, 0, 0.25)).
						AlignItems(ui.Center).Justify(ui.Center).
						Children(func() {
							ui.Text(c, "⌘").FontSize(9.5).Bold().TextColor(currentTheme.KeyCapTextColor)
						})
					ui.Text(c, currentDisplay).FontSize(13).Bold().TextColor(p.text)
					ui.Icon(c, iconChevronDown).Size(12, 12).TextColor(p.muted)
				})

			sel.Popup(func(panel ui.Element) {
				panel.
					Width(260).
					Background(p.card).
					Border(1, p.hairline).
					Radius(10).
					Padding(4).
					Shadow(0, 12, 32, 0, ui.RGBA(0, 0, 0, 0.50))

				ui.Scroll(c).MaxHeight(260).Children(func() {
					for _, th := range theme.AllThemes {
						thConfig := th
						displayName := themeDisplayName(thConfig.Name)
						item := sel.Item(thConfig.Name).
							Height(36).
							PaddingX(10).
							Gap(10).
							Radius(6).
							AlignItems(ui.Center)

						if item.Highlighted() {
							item.Background(p.press)
						} else if thConfig.Name == cfg.ThemePreset {
							item.Background(p.hover)
						}

						item.Children(func() {
							ui.Box(c).Size(22, 20).Radius(5).
								Background(thConfig.KeyCapBackground).
								Border(1.2, thConfig.KeyCapBorder).
								Shadow(0, 2, 6, 0, ui.RGBA(0, 0, 0, 0.25)).
								AlignItems(ui.Center).Justify(ui.Center).
								Children(func() {
									ui.Text(c, "⌘").FontSize(9.5).Bold().TextColor(thConfig.KeyCapTextColor)
								})
							ui.Text(c, displayName).FontSize(12.5).TextColor(p.text).Grow(1)
							if thConfig.Name == cfg.ThemePreset {
								ui.Text(c, "✓").FontSize(13).Bold().TextColor(p.accent)
							}
						})
					}
				})
			})

			if sel.Trigger.Changed() {
				th := theme.ThemeByName(themeChoice)
				// Picking a preset means "give me that design", so theme-derived
				// values return to inheritance instead of being frozen into the
				// config as copies of today's colors.
				sc.updateConfig(func(ac *config.AppConfig) {
					ac.ThemePreset = th.Name
					ac.Appearance.Style = ""
					ac.Color.Color = ""
					ac.Color.SecondaryColor = ""
					ac.Text.Color = ""
					ac.Border.Color = ""
					ac.Border.Radius = 0
				})
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})
	})

	// Group 2: Display Language & System
	group(i18n.T("lang_section_title"), func() {
		row(i18n.T("lang_title"), i18n.T("lang_desc"), func() {
			langCodes := []i18n.Language{i18n.LangZhCN, i18n.LangEnUS, i18n.LangZhTW, i18n.LangJaJP}
			labels := []string{"简体中文", "English", "繁體中文", "日本語"}
			curLang := i18n.Instance().GetLanguage()
			selIdx := 0
			for i, code := range langCodes {
				if curLang == code {
					selIdx = i
					break
				}
			}
			if ui.Segmented(c.Key("lang-segmented"), &selIdx, labels...).Changed() {
				target := langCodes[selIdx]
				i18n.Instance().SetLanguage(target)
				sc.updateConfig(func(ac *config.AppConfig) {
					ac.Language = string(target)
				})
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})
	})

	// Group 3: Filter Policy & Custom Filter Matrix
	group(i18n.T("filter_title"), func() {
		row(i18n.T("filter_title"), i18n.T("filter_desc"), func() {
			modes := []string{"none", "modifiers", "custom", "silent"}
			labels := []string{
				i18n.T("filter_none"),
				i18n.T("filter_modifiers"),
				i18n.T("filter_custom"),
				i18n.T("filter_silent"),
			}
			selIdx := 1
			for i, m := range modes {
				if m == cfg.General.Filter {
					selIdx = i
				}
			}
			if ui.Segmented(c.Key("filter-mode"), &selIdx, labels...).Changed() {
				target := modes[selIdx]
				sc.updateConfig(func(ac *config.AppConfig) { ac.General.Filter = target })
			}
		})

		if cfg.General.Filter == "custom" {
			row("自定义过滤名单", fmt.Sprintf("已启用 %d 个按键，点击打开自定义按键过滤矩阵", len(cfg.General.AllowedKeys)), func() {
				if ui.Button(c.Key("open-custom-filter-btn"), "配置按键矩阵...").
					Padding(6, 14).
					Radius(6).
					Background(p.accent).
					TextColor(p.onAccent).
					Clicked() {
					sc.showCustomFilterModal = true
					if sc.window != nil {
						sc.window.Invalidate()
					}
				}
			})
		}
	})

	// Group 4: History Trail, Word Grouping & Flow Direction
	group(i18n.T("history_title"), func() {
		row(i18n.T("history_title"), i18n.T("history_desc"), func() {
			on := cfg.General.ShowEventHistory
			if ui.Switch(c.Key("toggle-history"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.General.ShowEventHistory = on })
			}
		})

		row("单词合并", "连续键入字母合并为单词", func() {
			on := cfg.General.WordGrouping
			if ui.Switch(c.Key("toggle-word-grouping"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.General.WordGrouping = on })
			}
		})

		row(i18n.T("direction_title"), "历史按键排列方向", func() {
			dirIdx := 1
			if cfg.Appearance.FlexDirection == "row" {
				dirIdx = 0
			}
			if ui.Segmented(c.Key("flow-dir"), &dirIdx, i18n.T("direction_row"), i18n.T("direction_column")).Changed() {
				dirVal := "column"
				if dirIdx == 0 {
					dirVal = "row"
				}
				sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.FlexDirection = dirVal })
			}
		})

		if cfg.General.ShowEventHistory {
			row("留存组数", "保留按键组数", func() {
				ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
					if ui.Button(c.Key("hist-dec"), "−").Padding(4, 10).Radius(6).Clicked() && cfg.General.MaxHistory > 2 {
						sc.updateConfig(func(ac *config.AppConfig) { ac.General.MaxHistory-- })
					}
					ui.Textf(c, "%d 组", cfg.General.MaxHistory).FontSize(13).Bold().TextColor(p.text)
					if ui.Button(c.Key("hist-inc"), "+").Padding(4, 10).Radius(6).Clicked() && cfg.General.MaxHistory < 12 {
						sc.updateConfig(func(ac *config.AppConfig) { ac.General.MaxHistory++ })
					}
				})
			})
		} else {
			row("留存组数", "仅显示当前按键", func() {
				ui.Box(c).Padding(5, 10).Radius(6).Background(p.sidebar).Border(1, p.hairline).Children(func() {
					ui.Text(c, "固定 1 组").FontSize(12.5).Bold().TextColor(p.muted)
				})
			})
		}
	})

	// Group 5: Shortcuts & Feedback
	group("快捷键与音效", func() {
		row(i18n.T("toggle_shortcut"), i18n.T("toggle_shortcut_desc"), func() {
			ui.Box(c).Background(p.sidebar).Border(1, p.hairline).Radius(6).Padding(6, 12).Children(func() {
				ui.Text(c, formatShortcut(cfg.General.ToggleShortcut)).Bold().FontSize(12.5).TextColor(p.text)
			})
		})

		row("机械音效", "敲击按键时播放轻微微动音效", func() {
			on := cfg.TactileAudio
			if ui.Switch(c.Key("toggle-tactile-audio"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.TactileAudio = on })
			}
		})
	})
}

// renderAppearanceTab mirrors Cue Appearance tab with 9-Grid and smooth stepper controls.
func (sc *SettingsController) renderAppearanceTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
	stacked func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("appearance_title")).FontSize(20).Bold().TextColor(p.text)

	// Group 1: Alignment & Position
	group(i18n.T("alignment_title"), func() {
		row(i18n.T("alignment_title"), i18n.T("alignment_desc"), func() {
			sc.renderAlignmentSelector(c, cfg, p)
		})
	})

	// Group 2: Linger Duration & Margins
	group(i18n.T("linger_title"), func() {
		row("停留时长", "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("dur-dec"), "−0.5s").Padding(4, 8).Radius(6).Clicked() && cfg.General.LingerDurationMs > 1000 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.General.LingerDurationMs -= 500 })
				}
				ui.Textf(c, "%.1f 秒", float64(cfg.General.LingerDurationMs)/1000.0).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("dur-inc"), "+0.5s").Padding(4, 8).Radius(6).Clicked() && cfg.General.LingerDurationMs < 10000 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.General.LingerDurationMs += 500 })
				}
			})
		})

		row("锁定等比边距", "水平与垂直边距联动", func() {
			linked := sc.marginLinked
			if ui.Switch(c.Key("toggle-margin-link"), &linked).Changed() {
				sc.marginLinked = linked
				if linked {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.MarginY = ac.Appearance.MarginX })
				}
			}
		})

		row("水平边距 (X)", "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("mx-dec"), "−10").Padding(4, 8).Radius(6).Clicked() && cfg.Appearance.MarginX > 0 {
					sc.updateConfig(func(ac *config.AppConfig) {
						ac.Appearance.MarginX -= 10
						if sc.marginLinked {
							ac.Appearance.MarginY = ac.Appearance.MarginX
						}
					})
				}
				ui.Textf(c, "%d px", cfg.Appearance.MarginX).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("mx-inc"), "+10").Padding(4, 8).Radius(6).Clicked() && cfg.Appearance.MarginX < 300 {
					sc.updateConfig(func(ac *config.AppConfig) {
						ac.Appearance.MarginX += 10
						if sc.marginLinked {
							ac.Appearance.MarginY = ac.Appearance.MarginX
						}
					})
				}
			})
		})

		row("垂直边距 (Y)", "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("my-dec"), "−10").Padding(4, 8).Radius(6).Clicked() && cfg.Appearance.MarginY > 0 {
					sc.updateConfig(func(ac *config.AppConfig) {
						ac.Appearance.MarginY -= 10
						if sc.marginLinked {
							ac.Appearance.MarginX = ac.Appearance.MarginY
						}
					})
				}
				ui.Textf(c, "%d px", cfg.Appearance.MarginY).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("my-inc"), "+10").Padding(4, 8).Radius(6).Clicked() && cfg.Appearance.MarginY < 300 {
					sc.updateConfig(func(ac *config.AppConfig) {
						ac.Appearance.MarginY += 10
						if sc.marginLinked {
							ac.Appearance.MarginX = ac.Appearance.MarginY
						}
					})
				}
			})
		})
	})

	// Group 3: Transition Animations & Physics Toggle
	group(i18n.T("anim_title"), func() {
		row("动效形态", "入场与退场运动风格", func() {
			anims := []string{"spring", "fade", "slide", "float", "zoom", "none"}
			labels := []string{"弹性微动", "优雅渐隐", "顺势滑入", "浮动漂移", "聚焦缩放", "纯粹静态"}
			animIdx := 0
			cur := cfg.Appearance.Animation
			if cur == "" {
				cur = "spring"
			}
			for i, an := range anims {
				if an == cur {
					animIdx = i
					break
				}
			}
			if ui.Segmented(c.Key("anim-seg"), &animIdx, labels...).Changed() {
				chosen := anims[animIdx]
				sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.Animation = chosen })
				sc.previewNonce++
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})

		pressFxOn := !cfg.Appearance.DisablePressAnimation && cfg.Appearance.Animation != "none"
		row("按键行程与高亮", "轴体下潜与材质微受力反馈", func() {
			if ui.Switch(c.Key("press-fx-toggle"), &pressFxOn).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) {
					ac.Appearance.DisablePressAnimation = !pressFxOn
				})
				sc.previewNonce++
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})

		if cfg.Appearance.Animation != "none" {
			row("响应节奏", "过渡动画持续时间", func() {
				ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
					if ui.Button(c.Key("anim-dur-dec"), "−50ms").Padding(4, 8).Radius(6).Clicked() && cfg.Appearance.AnimationDuration > 0.10 {
						sc.updateConfig(func(ac *config.AppConfig) {
							ac.Appearance.AnimationDuration = clampAnimDuration(ac.Appearance.AnimationDuration - 0.05)
						})
						sc.previewNonce++
						if sc.window != nil {
							sc.window.Invalidate()
						}
					}
					ui.Textf(c, "%d ms", int(cfg.Appearance.AnimationDuration*1000)).
						FontSize(13).Bold().TextColor(p.text)
					if ui.Button(c.Key("anim-dur-inc"), "+50ms").Padding(4, 8).Radius(6).Clicked() && cfg.Appearance.AnimationDuration < 1.50 {
						sc.updateConfig(func(ac *config.AppConfig) {
							ac.Appearance.AnimationDuration = clampAnimDuration(ac.Appearance.AnimationDuration + 0.05)
						})
						sc.previewNonce++
						if sc.window != nil {
							sc.window.Invalidate()
						}
					}
				})
			})
		}
	})
}

// clampAnimDuration keeps the animation duration in a useful, non-zero band.
func clampAnimDuration(v float64) float64 {
	if v < 0.05 {
		return 0.05
	}
	if v > 2.0 {
		return 2.0
	}
	return v
}

// renderKeycapTab mirrors Cue Keycap tab with full 14 color palettes and typography.
func (sc *SettingsController) renderKeycapTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
	stacked func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("keycap_title")).FontSize(20).Bold().TextColor(p.text)
	currentTheme := theme.ThemeByName(cfg.ThemePreset)

	// Group 1: Physical Model (Floating Popover Web Dropdown with Preview)
	group(i18n.T("keycap_model"), func() {
		styles := []struct {
			id   string
			name string
		}{
			{"", "智能跟随原型 (Signature Archetype)"},
			{"lowprofile", "矮轴物理 (Low Profile · 2.0mm)"},
			{"minimal", "极简扁平 (Minimal Flat · 1.6mm)"},
			{"laptop", "剪刀脚 (Laptop Scissor · 1.2mm)"},
			{"pbt", "PBT 机械 (PBT Mechanical · 3.0mm)"},
			{"glass", "拟物玻璃 (Liquid Glass · 光脉冲)"},
			{"neon", "霓虹发光 (Neon Glow · 光晕脉冲)"},
			{"retro", "复古终端 (Retro CRT · 3.8mm)"},
			{"outline", "极简线框 (Wireframe Outline · 极简)"},
		}

		currentStyleName := fmt.Sprintf("%s · 原型预设 (%s)", physicalStyleName(effectiveStyle(currentTheme, cfg)), travelDepthLabel(currentTheme))
		styleChoice := cfg.Appearance.Style
		for _, st := range styles {
			if st.id == cfg.Appearance.Style && st.id != "" {
				currentStyleName = st.name + " · 自定义"
				break
			}
		}

		row(i18n.T("keycap_model"), "微调物理微几何与轴体行程", func() {
			sel := ui.SelectBase(c, &styleChoice)
			sel.Trigger.
				Height(36).
				PaddingX(12).
				Gap(10).
				Radius(8).
				Background(p.sidebar).
				Border(1, p.hairline).
				AlignItems(ui.Center).
				Children(func() {
					sc.renderStylePreviewIcon(c, effectiveStyle(currentTheme, cfg), p)
					ui.Text(c, currentStyleName).FontSize(13).Bold().TextColor(p.text)
					ui.Icon(c, iconChevronDown).Size(12, 12).TextColor(p.muted)
				})

			sel.Popup(func(panel ui.Element) {
				panel.
					Width(260).
					Background(p.card).
					Border(1, p.hairline).
					Radius(10).
					Padding(4).
					Shadow(0, 12, 32, 0, ui.RGBA(0, 0, 0, 0.50))

				ui.Scroll(c).MaxHeight(260).Children(func() {
					for _, st := range styles {
						styleItem := st
						item := sel.Item(styleItem.id).
							Height(36).
							PaddingX(10).
							Gap(10).
							Radius(6).
							AlignItems(ui.Center)

						if item.Highlighted() {
							item.Background(p.press)
						} else if styleItem.id == cfg.Appearance.Style {
							item.Background(p.hover)
						}

						item.Children(func() {
							previewStyle := styleItem.id
							if previewStyle == "" {
								previewStyle = effectiveStyle(currentTheme, cfg)
							}
							sc.renderStylePreviewIcon(c, previewStyle, p)
							ui.Text(c, styleItem.name).FontSize(12.5).TextColor(p.text).Grow(1)
							if styleItem.id == cfg.Appearance.Style {
								ui.Text(c, "✓").FontSize(13).Bold().TextColor(p.accent)
							}
						})
					}
				})
			})

			if sel.Trigger.Changed() {
				sc.updateConfig(func(ac *config.AppConfig) {
					ac.Appearance.Style = styleChoice
				})
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})
	})

	// Group 2: 14 Color Palettes (Floating Popover Web Dropdown with Preview)
	group("预设色盘", func() {
		// An empty palette means "inherit the theme". It is a first-class choice,
		// not an unexplained "custom" state on a fresh installation.
		paletteChoices := append([]ColorSchemePreset{{Name: "跟随主题 (Theme Default)"}}, MasterColorSchemes...)
		currentSchemeName := "跟随主题 (Theme Default)"
		paletteChoice := cfg.Color.Color
		if cfg.Color.Color != "" {
			currentSchemeName = "自定义配色"
			for _, sch := range MasterColorSchemes {
				if strings.EqualFold(cfg.Color.Color, sch.Primary) {
					currentSchemeName = sch.Name
					paletteChoice = sch.Primary
					break
				}
			}
		}

		row("键帽配色", "", func() {
			sel := ui.SelectBase(c, &paletteChoice)
			sel.Trigger.
				Height(36).
				PaddingX(12).
				Gap(10).
				Radius(8).
				Background(p.sidebar).
				Border(1, p.hairline).
				AlignItems(ui.Center).
				Children(func() {
					ui.Box(c).Size(22, 20).Radius(5).
						Background(parseHexSafe(cfg.Color.Color, currentTheme.KeyCapBackground)).
						Border(1.2, parseHexSafe(cfg.Color.SecondaryColor, currentTheme.KeyCapBorder)).
						Shadow(0, 2, 6, 0, ui.RGBA(0, 0, 0, 0.25)).
						AlignItems(ui.Center).Justify(ui.Center).
						Children(func() {
							ui.Text(c, "A").FontSize(9.5).Bold().TextColor(parseHexSafe(cfg.Text.Color, currentTheme.KeyCapTextColor))
						})
					ui.Text(c, currentSchemeName).FontSize(13).Bold().TextColor(p.text)
					ui.Icon(c, iconChevronDown).Size(12, 12).TextColor(p.muted)
				})

			sel.Popup(func(panel ui.Element) {
				panel.
					Width(260).
					Background(p.card).
					Border(1, p.hairline).
					Radius(10).
					Padding(4).
					Shadow(0, 12, 32, 0, ui.RGBA(0, 0, 0, 0.50))

				ui.Scroll(c).MaxHeight(260).Children(func() {
					for _, choice := range paletteChoices {
						sch := choice
						isInherit := sch.Primary == ""
						isSelected := (isInherit && cfg.Color.Color == "") ||
							(!isInherit && strings.EqualFold(cfg.Color.Color, sch.Primary))
						item := sel.Item(sch.Primary).
							Height(36).
							PaddingX(10).
							Gap(10).
							Radius(6).
							AlignItems(ui.Center)

						if item.Highlighted() {
							item.Background(p.press)
						} else if isSelected {
							item.Background(p.hover)
						}

						item.Children(func() {
							bgPreview := currentTheme.KeyCapBackground
							borderPreview := currentTheme.KeyCapBorder
							textPreview := currentTheme.KeyCapTextColor
							if !isInherit {
								bgPreview = parseHexSafe(sch.Primary, bgPreview)
								borderPreview = parseHexSafe(sch.Secondary, borderPreview)
								textPreview = parseHexSafe(sch.Text, textPreview)
							}
							ui.Box(c).Size(22, 20).Radius(5).
								Background(bgPreview).
								Border(1.2, borderPreview).
								Shadow(0, 2, 6, 0, ui.RGBA(0, 0, 0, 0.25)).
								AlignItems(ui.Center).Justify(ui.Center).
								Children(func() {
									ui.Text(c, "A").FontSize(9.5).Bold().TextColor(textPreview)
								})
							ui.Text(c, sch.Name).FontSize(12.5).TextColor(p.text).Grow(1)
							if isSelected {
								ui.Text(c, "✓").FontSize(13).Bold().TextColor(p.accent)
							}
						})
					}
				})
			})

			if sel.Trigger.Changed() {
				if paletteChoice == "" {
					sc.updateConfig(func(ac *config.AppConfig) {
						ac.Color.Color = ""
						ac.Color.SecondaryColor = ""
						ac.Text.Color = ""
						ac.Border.Color = ""
					})
				} else {
					for _, sch := range MasterColorSchemes {
						if strings.EqualFold(sch.Primary, paletteChoice) {
							sc.updateConfig(func(ac *config.AppConfig) {
								ac.Color.Color = sch.Primary
								ac.Color.SecondaryColor = sch.Secondary
								ac.Text.Color = sch.Text
								ac.Border.Color = sch.Secondary
							})
							break
						}
					}
				}
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})
	})

	// Group 3: Typography & Text Settings
	group(i18n.T("text_formatting"), func() {
		row(i18n.T("text_caps"), "", func() {
			caps := []string{"uppercase", "capitalize", "lowercase"}
			labels := []string{"全大写 (AA)", "首字母大写 (Aa)", "全小写 (aa)"}
			capIdx := 1
			for i, cp := range caps {
				if cp == cfg.Text.Caps {
					capIdx = i
				}
			}
			if ui.Segmented(c.Key("caps-seg"), &capIdx, labels...).Changed() {
				chosen := caps[capIdx]
				sc.updateConfig(func(ac *config.AppConfig) { ac.Text.Caps = chosen })
			}
		})

		row(i18n.T("text_variant"), "", func() {
			vars := []string{"text", "text-short", "icon"}
			labels := []string{"完整文字", "简写文字", "纯图标"}
			varIdx := 1
			for i, v := range vars {
				if v == cfg.Text.Variant {
					varIdx = i
				}
			}
			if ui.Segmented(c.Key("var-seg"), &varIdx, labels...).Changed() {
				chosen := vars[varIdx]
				sc.updateConfig(func(ac *config.AppConfig) { ac.Text.Variant = chosen })
			}
		})

		row("字体字号", "", func() {
			size := cfg.Text.Size
			suffix := ""
			if size <= 0 {
				size = int(effectiveFontSize(currentTheme, cfg) + 0.5)
				suffix = " · 主题"
			}
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("ts-dec"), "−").Padding(4, 8).Radius(6).Clicked() && size > 16 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Text.Size = size - 2 })
				}
				ui.Textf(c, "%d px%s", size, suffix).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("ts-inc"), "+").Padding(4, 8).Radius(6).Clicked() && size < 56 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Text.Size = size + 2 })
				}
			})
		})
	})

	// Group 4: Layout & Elements
	group(i18n.T("display_elements"), func() {
		row("按键图标", "显示修饰键图标", func() {
			on := cfg.Layout.ShowIcon
			if ui.Switch(c.Key("toggle-icon"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Layout.ShowIcon = on })
			}
		})

		row("特殊符号", "显示回车、退格等特殊符号", func() {
			on := cfg.Layout.ShowSymbol
			if ui.Switch(c.Key("toggle-symbol"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Layout.ShowSymbol = on })
			}
		})

		row("连击倍率", "连续敲击同一按键显示次数", func() {
			on := cfg.Layout.ShowPressCount
			if ui.Switch(c.Key("toggle-press-count"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Layout.ShowPressCount = on })
			}
		})

		row("修饰键高亮", "为修饰键应用强调色", func() {
			on := cfg.Modifier.Highlight
			if ui.Switch(c.Key("toggle-mod-hl"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Modifier.Highlight = on })
			}
		})
	})

	// Group 5: Border & Corner Radius
	group(i18n.T("border_title"), func() {
		row("键帽边框", "", func() {
			on := cfg.Border.Enabled
			if ui.Switch(c.Key("toggle-border"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Border.Enabled = on })
			}
		})

		row("边框粗细", "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("bw-dec"), "−").Padding(4, 8).Radius(6).Clicked() && cfg.Border.Width > 1 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Border.Width-- })
				}
				ui.Textf(c, "%d px", cfg.Border.Width).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("bw-inc"), "+").Padding(4, 8).Radius(6).Clicked() && cfg.Border.Width < 6 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Border.Width++ })
				}
			})
		})

		row("倒角圆角", "", func() {
			radius := cfg.Border.Radius
			suffix := ""
			if radius <= 0 {
				radius = float64(effectiveRadius(currentTheme, cfg))
				suffix = " · 主题"
			}
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("br-dec"), "−").Padding(4, 8).Radius(6).Clicked() && radius > 2 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Border.Radius = radius - 2 })
				}
				ui.Textf(c, "%.1f px%s", radius, suffix).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("br-inc"), "+").Padding(4, 8).Radius(6).Clicked() && radius < 24 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Border.Radius = radius + 2 })
				}
			})
		})
	})

	// Group 6: Intelligent Action Subtitles
	group("动作提示", func() {
		row("动作提示字幕", "识别复制、粘贴等常用操作意图", func() {
			on := cfg.SemanticSub
			if ui.Switch(c.Key("toggle-semantic-sub"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.SemanticSub = on })
			}
		})
	})
}

// renderMouseTab mirrors Cue Mouse settings with clean Cupertino card groups.
func (sc *SettingsController) renderMouseTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("mouse_title")).FontSize(20).Bold().TextColor(p.text)

	// Group 1: 鼠标预设主题 (Preset Mouse Themes)
	group("鼠标主题", func() {
		currentMTheme := theme.MouseThemeByID(cfg.Mouse.ThemePreset)
		mChoice := currentMTheme.ID

		row("预设主题", "光圈与按键指示色彩风格", func() {
			sel := ui.SelectBase(c, &mChoice)
			sel.Trigger.
				Height(36).
				PaddingX(12).
				Gap(10).
				Radius(8).
				Background(p.sidebar).
				Border(1, p.hairline).
				AlignItems(ui.Center).
				Children(func() {
					ui.Row(c).Gap(3).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(9, 14).Radius(3).Background(ui.Hex(currentMTheme.LeftColor))
						ui.Box(c).Size(6, 10).Radius(2).Background(ui.Hex(currentMTheme.WheelColor))
						ui.Box(c).Size(9, 14).Radius(3).Background(ui.Hex(currentMTheme.RightColor))
					})
					ui.Text(c, currentMTheme.DisplayName).FontSize(13).Bold().TextColor(p.text)
					ui.Icon(c, iconChevronDown).Size(12, 12).TextColor(p.muted)
				})

			sel.Popup(func(panel ui.Element) {
				panel.
					Width(260).
					Background(p.card).
					Border(1, p.hairline).
					Radius(10).
					Padding(4).
					Shadow(0, 12, 32, 0, ui.RGBA(0, 0, 0, 0.50))

				ui.Scroll(c).MaxHeight(260).Children(func() {
					for _, mth := range theme.AllMouseThemes {
						themeItem := mth
						item := sel.Item(themeItem.ID).
							Height(36).
							PaddingX(10).
							Gap(10).
							Radius(6).
							AlignItems(ui.Center)

						if item.Highlighted() {
							item.Background(p.press)
						} else if themeItem.ID == currentMTheme.ID {
							item.Background(p.hover)
						}

						item.Children(func() {
							ui.Row(c).Gap(3).AlignItems(ui.Center).Children(func() {
								ui.Box(c).Size(9, 14).Radius(3).Background(ui.Hex(themeItem.LeftColor))
								ui.Box(c).Size(6, 10).Radius(2).Background(ui.Hex(themeItem.WheelColor))
								ui.Box(c).Size(9, 14).Radius(3).Background(ui.Hex(themeItem.RightColor))
							})
							ui.Text(c, themeItem.DisplayName).FontSize(12.5).TextColor(p.text).Grow(1)
							if themeItem.ID == currentMTheme.ID {
								ui.Text(c, "✓").FontSize(13).Bold().TextColor(p.accent)
							}
						})
					}
				})
			})

			if sel.Trigger.Changed() {
				targetMth := theme.MouseThemeByID(mChoice)
				sc.updateConfig(func(ac *config.AppConfig) {
					ac.Mouse.ThemePreset = targetMth.ID
					ac.Mouse.Color = targetMth.LeftColor
				})
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})
	})

	// Group 2: Cursor Highlight & Ring
	group(i18n.T("mouse_click_ripple"), func() {
		row("点击光圈", "按下按键时呈现涟漪光圈", func() {
			on := cfg.Mouse.ShowClicks
			if ui.Switch(c.Key("toggle-mouse-clicks"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.ShowClicks = on })
			}
		})

		row("常驻光圈", "始终在光标周围显示", func() {
			on := cfg.Mouse.KeepHighlight
			if ui.Switch(c.Key("toggle-always-hl"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.KeepHighlight = on })
			}
		})

		row("光圈尺寸", "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("ms-dec"), "−10").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.Size > 40 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.Size -= 10 })
				}
				ui.Textf(c, "%d px", cfg.Mouse.Size).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("ms-inc"), "+10").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.Size < 300 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.Size += 10 })
				}
			})
		})
	})

	// Group 2: Mouse Button & Wheel Indicator
	group(i18n.T("mouse_indicator"), func() {
		row("按键指示器", "在光标旁显示按键与滚轮状态", func() {
			on := cfg.Mouse.ShowIndicator
			if ui.Switch(c.Key("toggle-indicator"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.ShowIndicator = on })
			}
		})

		row("常驻显示", "无按键时仍保持常驻", func() {
			on := cfg.Mouse.KeepIndicator
			if ui.Switch(c.Key("toggle-keep-ind"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.KeepIndicator = on })
			}
		})

		row("指示器大小", "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("is-dec"), "−5").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.IndicatorSize > 20 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.IndicatorSize -= 5 })
				}
				ui.Textf(c, "%d px", cfg.Mouse.IndicatorSize).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("is-inc"), "+5").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.IndicatorSize < 300 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.IndicatorSize += 5 })
				}
			})
		})

		row("光标偏移", "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("io-dec"), "−5").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.IndicatorOffsetX > 10 {
					sc.updateConfig(func(ac *config.AppConfig) {
						ac.Mouse.IndicatorOffsetX -= 5
						ac.Mouse.IndicatorOffsetY -= 5
					})
				}
				ui.Textf(c, "%d px", cfg.Mouse.IndicatorOffsetX).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("io-inc"), "+5").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.IndicatorOffsetX < 150 {
					sc.updateConfig(func(ac *config.AppConfig) {
						ac.Mouse.IndicatorOffsetX += 5
						ac.Mouse.IndicatorOffsetY += 5
					})
				}
			})
		})

		row("高亮保留时长", "按键与滚轮触发后的展示时长", func() {
			linger := cfg.Mouse.ClickLingerMs
			if linger <= 0 {
				linger = 180
			}
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("cl-dec"), "−20ms").Padding(4, 8).Radius(6).Clicked() && linger > 60 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.ClickLingerMs = linger - 20 })
				}
				ui.Textf(c, "%d ms", linger).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("cl-inc"), "+20ms").Padding(4, 8).Radius(6).Clicked() && linger < 1000 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.ClickLingerMs = linger + 20 })
				}
			})
		})
	})

	// Group 3: Event & Drag
	group("手势与拖拽", func() {
		row("拖拽阈值", "移动超过该距离识别为拖拽", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("dt-dec"), "−10").Padding(4, 8).Radius(6).Clicked() && cfg.General.DragThreshold > 10 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.General.DragThreshold -= 10 })
				}
				ui.Textf(c, "%d px", cfg.General.DragThreshold).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("dt-inc"), "+10").Padding(4, 8).Radius(6).Clicked() && cfg.General.DragThreshold < 150 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.General.DragThreshold += 10 })
				}
			})
		})
	})
}

// renderAboutTab displays clean architecture and system information.
func (sc *SettingsController) renderAboutTab(
	c *ui.Context,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("tab_about")).FontSize(20).Bold().TextColor(p.text)

	group("系统信息", func() {
		row(i18n.T("app_name"), i18n.T("about_desc"), func() {})
		row("渲染后端", i18n.T("renderer_architecture"), func() {})
		row("技术底座", i18n.T("engine_architecture"), func() {})
		row("灵感来源", i18n.T("original_credits"), func() {})
		row("设计理念", i18n.T("about_philosophy"), func() {})
	})

	group("配置管理", func() {
		row("恢复出厂设置", "还原所有设置项至初始参数", func() {
			if ui.Button(c.Key("btn-reset-factory"), i18n.T("reset_defaults")).
				Padding(8, 16).
				Radius(8).
				Background(p.accent).
				TextColor(p.onAccent).
				Clicked() {
				sc.showResetConfirm = true
				if sc.window != nil {
					sc.window.Invalidate()
				}
			}
		})
	})
}

// renderAlignmentSelector renders Cue dual positioning controls:
// 1. Tactile 3x3 directional button matrix on the left (quick, tactile arrow clicks)
// 2. 16:10 interactive monitor screen mini-map on the right (with glass display border and active HUD indicator)
// Both are seamlessly synchronized and clickable!
func (sc *SettingsController) renderAlignmentSelector(c *ui.Context, cfg config.AppConfig, p palette) {
	items := [][]struct {
		align config.Alignment
		arrow string
	}{
		{{config.AlignTopLeft, "↖"}, {config.AlignTopCenter, "↑"}, {config.AlignTopRight, "↗"}},
		{{config.AlignCenterLeft, "←"}, {config.AlignCenter, "•"}, {config.AlignCenterRight, "→"}},
		{{config.AlignBottomLeft, "↙"}, {config.AlignBottomCenter, "↓"}, {config.AlignBottomRight, "↘"}},
	}

	ui.Row(c.Key("align-dual-control")).Gap(14).AlignItems(ui.Center).Children(func() {
		// 1. Tactile 9-Grid Arrow Buttons
		ui.Column(c.Key("align-tactile-9grid")).Gap(4).Children(func() {
			for r, rowAligns := range items {
				ui.Row(c.Key(fmt.Sprintf("grid-row-%d", r))).Gap(4).Children(func() {
					for _, item := range rowAligns {
						al := item.align
						isSelected := cfg.Appearance.Alignment == al
						btn := ui.ButtonBase(c.Key("tactile-align-"+string(al))).
							Size(34, 30).
							Radius(6).
							AlignItems(ui.Center).
							Justify(ui.Center).
							FocusRing(false)

						if isSelected {
							btn.Background(p.accent)
						} else if btn.Hovered() {
							btn.Background(p.press)
						} else {
							btn.Background(p.sidebar)
						}
						btn.Border(1, p.hairline)

						textColor := p.muted
						if isSelected {
							textColor = p.onAccent
						}

						btn.Children(func() {
							ui.Text(c, item.arrow).FontSize(13).Bold().TextColor(textColor)
						})

						if btn.Clicked() {
							sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.Alignment = al })
							if sc.window != nil {
								sc.window.Invalidate()
							}
						}
					}
				})
			}
		})

		// 2. Minimalist & Clear Display Screen Mockup
		ui.Column(c.Key("align-computer-stage")).AlignItems(ui.Center).Children(func() {
			// Monitor Display Unit
			ui.Column(c.Key("display-unit")).
				Width(196).
				Height(126).
				Radius(8).
				Background(ui.Hex("#131417")).
				Border(1.2, ui.Hex("#2e3038")).
				Padding(4).
				Shadow(0, 8, 20, 0, ui.RGBA(0, 0, 0, 0.40)).
				Children(func() {
					// Inner Screen: Display Canvas
					ui.Column(c.Key("desktop-screen")).
						Fill().
						Radius(5).
						Background(ui.Hex("#0b0d14")).
						Padding(3).
						Gap(3).
						Children(func() {
							// Clean Screen Header Bar (Clearly indicating Display / Screen status)
							ui.Row(c.Key("screen-header")).
								Height(16).
								Fill().
								PaddingX(6).
								AlignItems(ui.Center).
								Justify(ui.SpaceBetween).
								Background(ui.RGBA(255, 255, 255, 0.05)).
								Radius(3).
								Children(func() {
									ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
										ui.Icon(c, iconMonitor).Size(10, 10).TextColor(p.muted)
										ui.Text(c, "主显示器").FontSize(8.5).Bold().TextColor(p.text)
									})
									ui.Text(c, "100%").FontSize(7.5).TextColor(p.muted)
								})

							// Interactive 3x3 Desktop Grid Matrix
							ui.Column(c.Key("desktop-grid")).Fill().Gap(2).Children(func() {
								for r, rowAligns := range items {
									ui.Row(c.Key(fmt.Sprintf("screen-row-%d", r))).Fill().Gap(2).Children(func() {
										for _, item := range rowAligns {
											al := item.align
											isSelected := cfg.Appearance.Alignment == al
											cell := ui.ButtonBase(c.Key("screen-cell-" + string(al))).
												Grow(1).
												Fill().
												Radius(3).
												AlignItems(ui.Center).
												Justify(ui.Center).
												FocusRing(false)

											if isSelected {
												cell.Background(p.accent.Alpha(0.22))
												cell.Border(1.0, p.accent)
											} else if cell.Hovered() {
												cell.Background(ui.RGBA(255, 255, 255, 0.08))
											} else {
												cell.Background(ui.Transparent)
											}

											cell.Children(func() {
												if isSelected {
													// Crisp illuminated Keyviz HUD position pill
													ui.Row(c).
														Padding(2, 6).
														Radius(3).
														Background(p.accent).
														Shadow(0, 2, 6, 0, p.accent.Alpha(0.5)).
														AlignItems(ui.Center).
														Children(func() {
															ui.Text(c, "HUD").FontSize(8).Bold().TextColor(p.onAccent)
														})
												} else {
													ui.Box(c).Size(3, 3).Radius(1.5).Background(ui.RGBA(255, 255, 255, 0.18))
												}
											})

											if cell.Clicked() {
												sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.Alignment = al })
												if sc.window != nil {
													sc.window.Invalidate()
												}
											}
										}
									})
								}
							})
						})
				})

			// Minimalist Display Stand: Neck
			ui.Box(c.Key("stand-neck")).
				Width(24).
				Height(12).
				Background(ui.Hex("#3f3f46")).
				Border(0.8, ui.Hex("#27272a"))

			// Minimalist Display Stand: Sleek Base Plate
			ui.Box(c.Key("stand-base")).
				Width(68).
				Height(3.5).
				Radius(2).
				Background(ui.Hex("#52525b")).
				Border(0.8, ui.Hex("#3f3f46")).
				Shadow(0, 3, 6, 0, ui.RGBA(0, 0, 0, 0.30))
		})
	})
}

// renderStylePreviewIcon renders a high-fidelity miniature keycap previewing physical styling models.
func (sc *SettingsController) renderStylePreviewIcon(c *ui.Context, styleId string, p palette) {
	bg := p.press
	border := p.hairline
	textColor := p.text
	radius := float32(5)
	borderW := float32(1.0)
	shadow := float32(0)

	switch styleId {
	case "lowprofile":
		radius = 6
		bg = p.press
	case "minimal":
		radius = 3
		bg = p.canvas
		borderW = 0
	case "laptop":
		radius = 4
		bg = p.sidebar
		border = p.muted.Alpha(0.4)
	case "pbt":
		radius = 6
		bg = p.sidebar
		borderW = 2.0
		border = p.accent.Alpha(0.6)
		shadow = 4
	case "glass":
		radius = 6
		bg = ui.RGBA(255, 255, 255, 0.18)
		border = ui.RGBA(255, 255, 255, 0.50)
	case "neon":
		radius = 5
		bg = ui.Hex("#09090b")
		border = ui.Hex("#00FF6A")
		textColor = ui.Hex("#00FF6A")
	case "retro":
		radius = 2
		bg = ui.Hex("#1c1917")
		border = ui.Hex("#22c55e")
		textColor = ui.Hex("#22c55e")
	case "outline":
		radius = 5
		bg = ui.Transparent
		border = p.accent
		textColor = p.accent
	}

	box := ui.Box(c).Size(22, 20).Radius(radius).
		Background(bg).
		AlignItems(ui.Center).Justify(ui.Center)

	if borderW > 0 {
		box.Border(borderW, border)
	}
	if shadow > 0 {
		box.Shadow(0, 2, shadow, 0, ui.RGBA(0, 0, 0, 0.35))
	}

	box.Children(func() {
		ui.Text(c, "A").FontSize(9.5).Bold().TextColor(textColor)
	})
}

// renderResetConfirmDialog renders the Jobsian confirmation modal before restoring default parameters.
func (sc *SettingsController) renderResetConfirmDialog(c *ui.Context, p palette) {
	ui.Box(c.Key("modal-backdrop")).
		Absolute().
		Fill().
		Background(ui.RGBA(0, 0, 0, 0.50)).
		AlignItems(ui.Center).
		Justify(ui.Center).
		Children(func() {
			ui.Column(c.Key("reset-confirm-box")).
				Width(420).
				Background(p.card).
				Border(1, p.hairline).
				Radius(14).
				Padding(24).
				Gap(16).
				Shadow(0, 16, 40, 0, ui.RGBA(0, 0, 0, 0.45)).
				Children(func() {
					ui.Row(c).AlignItems(ui.Center).Gap(12).Children(func() {
						ui.Box(c).Size(40, 40).Radius(20).Background(ui.Hex("#ff3b30").Alpha(0.12)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
							ui.Icon(c, iconReset).Size(22, 22).TextColor(ui.Hex("#ff3b30"))
						})
						ui.Column(c).Gap(2).Children(func() {
							ui.Text(c, "确认恢复出厂设置？").FontSize(16).Bold().TextColor(p.text)
							ui.Text(c, "此操作不可撤销").FontSize(11.5).TextColor(p.muted)
						})
					})

					ui.Text(c, "将重置所有自定义键帽颜色、9 宫格定位、过滤策略、快捷键与动画过渡风格，立即还原为 Cue 官方出厂标准参数。").
						FontSize(13).
						TextColor(p.muted).
						LineHeight(1.5)

					ui.Row(c).Justify(ui.End).Gap(10).Margin(8, 0, 0, 0).Children(func() {
						if ui.Button(c.Key("cancel-reset-btn"), "取消").
							Padding(8, 16).
							Radius(8).
							Background(p.sidebar).
							TextColor(p.text).
							Clicked() {
							sc.showResetConfirm = false
							if sc.window != nil {
								sc.window.Invalidate()
							}
						}

						if ui.Button(c.Key("confirm-reset-btn"), "确认恢复出厂").
							Padding(8, 18).
							Radius(8).
							Background(ui.Hex("#ff3b30")).
							TextColor(ui.Hex("#ffffff")).
							Bold().
							Clicked() {
							sc.showResetConfirm = false
							// One source of truth: restore DefaultConfig and let the
							// central apply path re-sync engine, theme, and overlay.
							sc.updateConfig(func(ac *config.AppConfig) {
								*ac = config.DefaultConfig()
							})
							sc.previewNonce++
							if sc.window != nil {
								sc.window.Invalidate()
							}
						}
					})
				})
		})
}

// isAllowedKey checks if keyName is contained in the allowed list (case-insensitive).
func isAllowedKey(key string, list []string) bool {
	for _, k := range list {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}

// toggleAllowedKey toggles membership of a key in the allowed list.
func (sc *SettingsController) toggleAllowedKey(key string) {
	sc.updateConfig(func(ac *config.AppConfig) {
		if isAllowedKey(key, ac.General.AllowedKeys) {
			var rem []string
			for _, k := range ac.General.AllowedKeys {
				if !strings.EqualFold(k, key) {
					rem = append(rem, k)
				}
			}
			ac.General.AllowedKeys = rem
		} else {
			ac.General.AllowedKeys = append(ac.General.AllowedKeys, key)
		}
	})
	if sc.window != nil {
		sc.window.Invalidate()
	}
}

// renderKeycapButton renders a tactile, interactive virtual key for the filter matrix.
func (sc *SettingsController) renderKeycapButton(c *ui.Context, rawKey, label string, w, h float32, grow float32, cfg config.AppConfig, p palette) {
	allowed := isAllowedKey(rawKey, cfg.General.AllowedKeys)
	btn := ui.ButtonBase(c.Key("vkey-" + rawKey)).
		Height(h).
		Radius(6).
		AlignItems(ui.Center).
		Justify(ui.Center).
		FocusRing(false)

	if w > 0 {
		btn.Width(w)
	}
	if grow > 0 {
		btn.Grow(grow)
	}

	bg := p.sidebar
	textColor := p.muted
	border := p.hairline
	if allowed {
		bg = p.accent.Alpha(0.24)
		border = p.accent
		textColor = p.accent
	}

	btn.Background(bg)
	btn.Border(1, border)

	btn.Children(func() {
		ui.Text(c, label).FontSize(11).Bold().TextColor(textColor)
	})

	if btn.Clicked() {
		sc.toggleAllowedKey(rawKey)
	}
}

// renderCustomFilterModal mirrors Cue original custom-filter.tsx:
// a realistic interactive virtual ANSI keyboard, mouse layout, and numpad matrix.
func (sc *SettingsController) renderCustomFilterModal(c *ui.Context, cfg config.AppConfig, p palette) {
	ui.Box(c.Key("filter-modal-backdrop")).
		Absolute().
		Fill().
		Background(ui.RGBA(0, 0, 0, 0.55)).
		AlignItems(ui.Center).
		Justify(ui.Center).
		Children(func() {
			ui.Column(c.Key("filter-modal-card")).
				Width(740).
				MaxHeight(560).
				Background(p.card).
				Border(1, p.hairline).
				Radius(14).
				Padding(20).
				Gap(14).
				Shadow(0, 16, 40, 0, ui.RGBA(0, 0, 0, 0.50)).
				Children(func() {
					// Header
					ui.Row(c).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
						ui.Column(c).Gap(2).Children(func() {
							ui.Text(c, "自定义按键过滤面板 (Custom Filter)").FontSize(16).Bold().TextColor(p.text)
							ui.Text(c, "高亮键位将在屏幕上正常呈现，暗色键位将被静音过滤").FontSize(11.5).TextColor(p.muted)
						})
						if ui.Button(c.Key("close-filter-modal"), "✕ 完成").Padding(6, 14).Radius(6).Background(p.accent).TextColor(p.onAccent).Bold().Clicked() {
							sc.showCustomFilterModal = false
							if sc.window != nil {
								sc.window.Invalidate()
							}
						}
					})

					// Tab Switcher & Quick Actions Row
					ui.Row(c).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
						// 3 Virtual Tabs
						tabs := []string{"Keyboard", "Mouse", "Numpad"}
						tabLabels := []string{"键盘 (Keyboard)", "鼠标 (Mouse)", "数字键盘 (Numpad)"}
						tabIdx := 0
						for i, tb := range tabs {
							if tb == sc.filterTab {
								tabIdx = i
							}
						}
						if ui.Segmented(c.Key("filter-tab-seg"), &tabIdx, tabLabels...).Changed() {
							sc.filterTab = tabs[tabIdx]
							if sc.window != nil {
								sc.window.Invalidate()
							}
						}

						// Quick Presets
						ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
							allAlpha := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
							allMods := []string{"Ctrl", "Shift", "Alt", "Win"}

							if ui.Button(c.Key("filter-all-btn"), "全选").Padding(4, 9).Radius(6).Background(p.sidebar).TextColor(p.text).FontSize(11).Clicked() {
								allCommon := append(append([]string{}, allAlpha...), allMods...)
								allCommon = append(allCommon, "Enter", "Space", "Backspace", "Tab", "Esc", "Insert", "Delete", "Home", "End", "PageUp", "PageDown", "Up", "Down", "Left", "Right", "Left Click", "Right Click", "Middle Click", "Drag", "Scroll Up", "Scroll Down")
								sc.updateConfig(func(ac *config.AppConfig) {
									ac.General.AllowedKeys = allCommon
								})
							}
							if ui.Button(c.Key("filter-clear-btn"), "清空").Padding(4, 9).Radius(6).Background(p.sidebar).TextColor(p.text).FontSize(11).Clicked() {
								sc.updateConfig(func(ac *config.AppConfig) {
									ac.General.AllowedKeys = []string{}
								})
							}
							if ui.Button(c.Key("filter-mods-btn"), "仅修饰键").Padding(4, 9).Radius(6).Background(p.sidebar).TextColor(p.text).FontSize(11).Clicked() {
								sc.updateConfig(func(ac *config.AppConfig) {
									ac.General.AllowedKeys = allMods
								})
							}
							if ui.Button(c.Key("filter-alpha-btn"), "字母与数字").Padding(4, 9).Radius(6).Background(p.sidebar).TextColor(p.text).FontSize(11).Clicked() {
								sc.updateConfig(func(ac *config.AppConfig) {
									ac.General.AllowedKeys = allAlpha
								})
							}
						})
					})

					// Main Interactive Body
					ui.Column(c.Key("filter-body-container")).
						Fill().
						Padding(12).
						Radius(10).
						Background(p.sidebar).
						Border(1, p.hairline).
						AlignItems(ui.Center).
						Justify(ui.Center).
						Children(func() {
							switch sc.filterTab {
							case "Mouse":
								// Mouse Visual Matrix (Cue original mouse silhouette)
								ui.Row(c.Key("mouse-filter-matrix")).Gap(18).AlignItems(ui.Center).Children(func() {
									// Left Click Button
									sc.renderKeycapButton(c, "Left Click", "左键 Left", 64, 48, 0, cfg, p)

									// Center Column (Middle, Scroll, Drag)
									ui.Column(c).Gap(6).AlignItems(ui.Center).Children(func() {
										sc.renderKeycapButton(c, "Middle Click", "中键 Middle", 80, 32, 0, cfg, p)
										sc.renderKeycapButton(c, "Scroll Up", "滚轮向上 ↑", 80, 30, 0, cfg, p)
										sc.renderKeycapButton(c, "Scroll Down", "滚轮向下 ↓", 80, 30, 0, cfg, p)
										sc.renderKeycapButton(c, "Drag", "鼠标拖拽 Drag", 80, 32, 0, cfg, p)
									})

									// Right Click Button
									sc.renderKeycapButton(c, "Right Click", "右键 Right", 64, 48, 0, cfg, p)
								})

							case "Numpad":
								// 4x5 Numeric Keypad
								ui.Column(c.Key("numpad-filter-matrix")).Gap(5).Width(240).Children(func() {
									// Row 1: NumLock / * -
									ui.Row(c).Gap(5).Children(func() {
										sc.renderKeycapButton(c, "NumLock", "Num", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "/", "/", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "*", "*", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "-", "-", 52, 34, 0, cfg, p)
									})
									// Row 2: 7 8 9 +
									ui.Row(c).Gap(5).Children(func() {
										sc.renderKeycapButton(c, "7", "7", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "8", "8", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "9", "9", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "+", "+", 52, 34, 0, cfg, p)
									})
									// Row 3: 4 5 6
									ui.Row(c).Gap(5).Children(func() {
										sc.renderKeycapButton(c, "4", "4", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "5", "5", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "6", "6", 52, 34, 0, cfg, p)
										ui.Box(c).Size(52, 34) // placeholder
									})
									// Row 4: 1 2 3 Enter
									ui.Row(c).Gap(5).Children(func() {
										sc.renderKeycapButton(c, "1", "1", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "2", "2", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "3", "3", 52, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "Enter", "Enter", 52, 34, 0, cfg, p)
									})
									// Row 5: 0 .
									ui.Row(c).Gap(5).Children(func() {
										sc.renderKeycapButton(c, "0", "0", 109, 34, 0, cfg, p)
										sc.renderKeycapButton(c, ".", ".", 52, 34, 0, cfg, p)
										ui.Box(c).Size(52, 34) // placeholder
									})
								})

							default: // "Keyboard"
								// Realistic 6-Row ANSI Keyboard (Cue custom-filter.tsx layout)
								ui.Column(c.Key("keyboard-filter-matrix")).Gap(5).Fill().Children(func() {
									// Row 1: Esc, F1-F12, Ins, Del
									ui.Row(c).Gap(4).Fill().Children(func() {
										sc.renderKeycapButton(c, "Esc", "Esc", 42, 32, 0, cfg, p)
										for i := 1; i <= 12; i++ {
											fk := fmt.Sprintf("F%d", i)
											sc.renderKeycapButton(c, fk, fk, 36, 32, 1, cfg, p)
										}
										sc.renderKeycapButton(c, "Insert", "Ins", 38, 32, 0, cfg, p)
										sc.renderKeycapButton(c, "Delete", "Del", 38, 32, 0, cfg, p)
									})

									// Row 2: ~ 1 2 3 4 5 6 7 8 9 0 - = Backspace Home
									ui.Row(c).Gap(4).Fill().Children(func() {
										sc.renderKeycapButton(c, "`", "~", 36, 32, 0, cfg, p)
										for _, d := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "="} {
											sc.renderKeycapButton(c, d, d, 36, 32, 1, cfg, p)
										}
										sc.renderKeycapButton(c, "Backspace", "⌫", 58, 32, 0, cfg, p)
										sc.renderKeycapButton(c, "Home", "Home", 38, 32, 0, cfg, p)
									})

									// Row 3: Tab Q W E R T Y U I O P [ ] \ End
									ui.Row(c).Gap(4).Fill().Children(func() {
										sc.renderKeycapButton(c, "Tab", "Tab ⇥", 50, 32, 0, cfg, p)
										for _, q := range []string{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "[", "]", "\\"} {
											sc.renderKeycapButton(c, q, q, 36, 32, 1, cfg, p)
										}
										sc.renderKeycapButton(c, "End", "End", 38, 32, 0, cfg, p)
									})

									// Row 4: Caps A S D F G H J K L ; ' Enter PgUp
									ui.Row(c).Gap(4).Fill().Children(func() {
										sc.renderKeycapButton(c, "CapsLock", "Caps ⇪", 58, 32, 0, cfg, p)
										for _, a := range []string{"A", "S", "D", "F", "G", "H", "J", "K", "L", ";", "'"} {
											sc.renderKeycapButton(c, a, a, 36, 32, 1, cfg, p)
										}
										sc.renderKeycapButton(c, "Enter", "Enter ↵", 66, 32, 0, cfg, p)
										sc.renderKeycapButton(c, "PageUp", "PgUp", 38, 32, 0, cfg, p)
									})

									// Row 5: Shift Z X C V B N M , . / Shift Up PgDn
									ui.Row(c).Gap(4).Fill().Children(func() {
										sc.renderKeycapButton(c, "Shift", "Shift ⇧", 72, 32, 0, cfg, p)
										for _, z := range []string{"Z", "X", "C", "V", "B", "N", "M", ",", ".", "/"} {
											sc.renderKeycapButton(c, z, z, 36, 32, 1, cfg, p)
										}
										sc.renderKeycapButton(c, "Shift", "Shift ⇧", 56, 32, 0, cfg, p)
										sc.renderKeycapButton(c, "Up", "↑", 36, 32, 0, cfg, p)
										sc.renderKeycapButton(c, "PageDown", "PgDn", 38, 32, 0, cfg, p)
									})

									// Row 6: Ctrl Alt Win Space Ctrl Left Down Right
									ui.Row(c).Gap(4).Fill().Children(func() {
										sc.renderKeycapButton(c, "Ctrl", "Ctrl ⌃", 46, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "Alt", "Alt ⌥", 42, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "Win", "Win ⌘", 42, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "Space", "Space ␣", 0, 34, 4, cfg, p)
										sc.renderKeycapButton(c, "Ctrl", "Ctrl", 46, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "Left", "←", 36, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "Down", "↓", 36, 34, 0, cfg, p)
										sc.renderKeycapButton(c, "Right", "→", 36, 34, 0, cfg, p)
									})
								})
							}
						})
				})
		})
}
