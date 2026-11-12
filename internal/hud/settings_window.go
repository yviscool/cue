package hud

import (
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/assets"
	"github.com/yviscool/cue/internal/config"
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/engine"
	"github.com/yviscool/cue/internal/hook"
	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/semantic"
	"github.com/yviscool/cue/internal/theme"
	"github.com/yviscool/cue/internal/version"
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
	iconPlus        = strokeIcon(`<path d="M12 5v14M5 12h14"/>`)
	iconSparkles    = strokeIcon(`<path d="M12 3l1.6 4.4L18 9l-4.4 1.6L12 15l-1.6-4.4L6 9l4.4-1.6z"/><path d="M19 15l.8 2.2L22 18l-2.2.8L19 21l-.8-2.2L16 18l2.2-.8z"/>`)
	iconVolume      = strokeIcon(`<path d="M11 5 6 9H3v6h3l5 4z"/><path d="M15.5 8.5a5 5 0 0 1 0 7"/><path d="M18.5 5.5a9 9 0 0 1 0 13"/>`)
	iconFolder      = strokeIcon(`<path d="M3 7a2 2 0 0 1 2-2h5l2 2h7a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>`)
)

// themeNameKeys maps a preset theme name (lower-cased) to its localization key.
// The theme package keeps a Chinese-first DisplayName for documentation and for
// ThemeByName lookups, but the studio renders these localized labels so the
// interface never mixes languages when the user switches to English/Japanese.
var themeNameKeys = map[string]string{
	"liquid glass":          "theme_liquid_glass",
	"macbook scissor":       "theme_macbook_scissor",
	"dieter rams 1960":      "theme_dieter_rams",
	"cyberpunk neon":        "theme_cyberpunk_neon",
	"azure pulse":           "theme_azure_pulse",
	"classic mechanical":    "theme_classic_mechanical",
	"ibm model m 1984":      "theme_ibm_model_m",
	"minimal chalk":         "theme_minimal_chalk",
	"obsidian oled":         "theme_obsidian_oled",
	"tokyo blossom":         "theme_tokyo_blossom",
	"emerald matrix":        "theme_emerald_matrix",
	"sunset horizon":        "theme_sunset_horizon",
	"royal velvet":          "theme_royal_velvet",
	"sony walkman 1979":     "theme_sony_walkman",
	"teenage op-1":          "theme_teenage_op1",
	"amber cathode 1977":    "theme_amber_cathode",
	"swiss typography 1957": "theme_swiss_typography",
	"kyoto kintsugi":        "theme_kyoto_kintsugi",
	"blueprint cyanotype":   "theme_blueprint_cyanotype",
	"leica m":               "theme_leica_m",
}

// themeDisplayName returns the localized display name for a preset theme.
func themeDisplayName(name string) string {
	if key, ok := themeNameKeys[strings.ToLower(strings.TrimSpace(name))]; ok {
		return i18n.T(key)
	}
	th := theme.ThemeByName(name)
	if th.DisplayName != "" {
		return th.DisplayName
	}
	return name
}

// themeDescription returns the localized one-line description for a preset theme.
func themeDescription(name string) string {
	if key, ok := themeNameKeys[strings.ToLower(strings.TrimSpace(name))]; ok {
		return i18n.T(key + "_desc")
	}
	return theme.ThemeByName(name).Description
}

// mouseThemeKeys maps a mouse theme ID to its localization key.
var mouseThemeKeys = map[string]string{
	"neon_green":       "mouse_theme_neon_green",
	"aurora_bicolor":   "mouse_theme_aurora",
	"cyber_cyan":       "mouse_theme_cyber_cyan",
	"obsidian_crimson": "mouse_theme_obsidian_crimson",
	"deep_blue":        "mouse_theme_deep_blue",
	"minimal_white":    "mouse_theme_minimal_white",
	"sakura_mist":      "mouse_theme_sakura_mist",
	"arctic_ice":       "mouse_theme_arctic_ice",
	"graphite_mono":    "mouse_theme_graphite_mono",
	"amber_terminal":   "mouse_theme_amber_terminal",
	"acid_lime":        "mouse_theme_acid_lime",
	"kyoto_lacquer":    "mouse_theme_kyoto_lacquer",
}

// mouseThemeDisplayName returns the localized label for a mouse color theme.
func mouseThemeDisplayName(m theme.MouseThemeConfig) string {
	if key, ok := mouseThemeKeys[m.ID]; ok {
		return i18n.T(key)
	}
	if m.DisplayName != "" {
		return m.DisplayName
	}
	return m.Name
}

// physicalStyleKeys maps a keycap physical model ID to its localization key.
var physicalStyleKeys = map[string]string{
	"minimal": "ui_style_minimal",
	"laptop":  "ui_style_laptop",
	"pbt":     "ui_style_pbt",
	"glass":   "ui_style_glass",
	"neon":    "ui_style_neon",
	"retro":   "ui_style_retro",
	"outline": "ui_style_outline",
}

// physicalStyleShortKeys are the compact labels used by the model gallery,
// where a full name such as "Laptop Scissor (1.2 mm)" would wrap and compete
// with the preview itself.
var physicalStyleShortKeys = map[string]string{
	"lowprofile": "ui_keycap_style_short_lowprofile",
	"minimal":    "ui_keycap_style_short_minimal",
	"laptop":     "ui_keycap_style_short_laptop",
	"pbt":        "ui_keycap_style_short_pbt",
	"glass":      "ui_keycap_style_short_glass",
	"neon":       "ui_keycap_style_short_neon",
	"retro":      "ui_keycap_style_short_retro",
	"outline":    "ui_keycap_style_short_outline",
}

// physicalStyleShortName returns the gallery label for a model. An empty ID is
// the Auto choice, which is a mode rather than a ninth physical build.
func physicalStyleShortName(id string) string {
	if id == "" {
		return i18n.T("ui_keycap_style_short_auto")
	}
	if key, ok := physicalStyleShortKeys[id]; ok {
		return i18n.T(key)
	}
	return i18n.T("ui_keycap_style_short_lowprofile")
}

// physicalStyleName returns the localized name for a keycap physical model.
func physicalStyleName(id string) string {
	if key, ok := physicalStyleKeys[id]; ok {
		return i18n.T(key)
	}
	return i18n.T("ui_style_lowprofile")
}

// styleTravelDepth returns the travel the HUD will actually use for a model.
func styleTravelDepth(t theme.ThemeConfig, styleID string) float32 {
	return t.EffectiveTravelDepthFor(styleID)
}

// styleTravelLabel names the physical response: a compact distance for
// mechanical builds, a pulse for optical ones. The gallery rail already says
// "travel", so the word does not need to repeat in every card.
func styleTravelLabel(t theme.ThemeConfig, styleID string) string {
	depth := styleTravelDepth(t, styleID)
	if depth <= 0 {
		return i18n.T("ui_held_pulse")
	}
	return i18n.TF("ui_travel_mm", depth)
}

// travelDepthLabel returns the localized physical stroke depth or kinetic mode.
func travelDepthLabel(t theme.ThemeConfig) string {
	if t.HeldPhysics == theme.HeldPulse {
		return i18n.T("ui_held_pulse")
	}
	return i18n.TF("ui_travel_depth", t.EffectiveTravelDepth())
}

// styleStageLabel describes whether the physical model is inherited from the
// archetype or explicitly overridden.
func styleStageLabel(cfg config.AppConfig, t theme.ThemeConfig) string {
	if cfg.Appearance.Style == "" {
		return i18n.TF("ui_stage_preset", physicalStyleName(effectiveStyle(t, cfg)), travelDepthLabel(t))
	}
	return i18n.TF("ui_stage_override", physicalStyleName(cfg.Appearance.Style))
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
	store                   *config.Store
	engine                  *engine.Engine
	vState                  *ViewState
	activeTab               string
	window                  *mygo.Window
	logoSrc                 ui.ImageSource
	marginLinked            bool
	scrollState             ui.ScrollState
	showResetConfirm        bool
	showCustomFilterModal   bool
	showCustomShortcutModal bool
	filterTab               string
	previewNonce            int
	contextSource           func() context.AppContext
	audioPreview            func(profile string)
	// audioPulse drives the audition pill: 1 while the sound plays, then back
	// to 0. audioPulseAt gives the meter its phase, so the bars move with the
	// clock rather than with the frame count.
	audioPulse          float32
	audioPulseAt        time.Time
	customDrafts        []config.CustomShortcut
	excludedDraft       string
	excludedInitialized bool
	configNotice        string
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
		activeTab:    "Keycap",
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

// renderPercentSlider draws a percentage slider bound to a 0..1 configuration
// fraction.
//
// The slider is stepped to 5% so dragging it writes config.json at most ~20
// times instead of once per frame; Store.Update rewrites the file atomically on
// every call, so an unstepped slider would hammer the disk. The readout is drawn
// after the widget so it shows the value the user just dragged to.
func (sc *SettingsController) renderPercentSlider(
	c *ui.Context,
	key string,
	fraction, minPercent float64,
	apply func(next float64),
	p palette,
) {
	percent := fraction * 100
	if percent < minPercent {
		percent = minPercent
	}
	if percent > 100 {
		percent = 100
	}

	ui.Row(c).AlignItems(ui.Center).Gap(10).Width(240).Children(func() {
		slider := ui.Slider(c.Key(key), &percent, minPercent, 100).Step(5).Grow(1)
		ui.Textf(c, "%.0f%%", percent).FontSize(12.5).Bold().TextColor(p.text).Width(44)
		if slider.Changed() {
			apply(percent / 100)
		}
	})
}

// SetContextSource installs a foreground-application snapshot source used by
// the General tab's live "recognized application" preview.
func (sc *SettingsController) SetContextSource(fn func() context.AppContext) {
	sc.contextSource = fn
}

// SetAudioPreview installs a one-shot audition callback for the Sound tab. The
// settings package stays independent of the acoustic engine: main owns the
// player and decides how a preview is rendered.
func (sc *SettingsController) SetAudioPreview(fn func(profile string)) {
	sc.audioPreview = fn
}

// ellipsize shortens a label on rune boundaries so CJK titles keep their
// glyphs intact.
func ellipsize(s string, maxRunes int) string {
	s = strings.TrimSpace(s)
	if maxRunes <= 0 {
		return s
	}
	r := []rune(s)
	if len(r) <= maxRunes {
		return s
	}
	return string(r[:maxRunes-1]) + "…"
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
					ui.Text(c, version.Display()+" · Direct3D 11").FontSize(11).TextColor(p.muted)
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
					if sc.window != nil {
						sc.window.Invalidate()
					}
				}
			}

			navTab("Keycap", "tab_keycap", iconKeyboard)
			navTab("Stage", "tab_appearance", iconMonitor)
			navTab("Mouse", "tab_mouse", iconMouse)
			navTab("Intelligence", "tab_intelligence", iconSparkles)
			navTab("Sound", "tab_sound", iconVolume)

			ui.Spacer(c)

			navTab("App", "tab_general", iconSettings)
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
			// Persistent live showcase stage (pinned at top). Every page gets the
			// same stage, so a dropdown can never be the only evidence of what
			// was chosen: the answer appears above it, at product scale.
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
					case "Keycap":
						sc.renderKeycapTab(c, cfg, p, group, row, stacked)
					case "Stage":
						sc.renderStageTab(c, cfg, p, group, row, stacked)
					case "Mouse":
						sc.renderMouseTab(c, cfg, p, group, row)
					case "Intelligence":
						sc.renderIntelligenceTab(c, cfg, p, group, row)
					case "Sound":
						sc.renderSoundTab(c, cfg, p, group, row)
					case "App":
						sc.renderAppTab(c, cfg, p, group, row)
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
	if sc.showCustomShortcutModal {
		sc.renderCustomShortcutModal(c, cfg, p)
	}
}

// previewStageColor is the backdrop a theme was drawn against, so the preview
// stops lying about what the overlay will look like. An opaque cap carries its
// own value, and the stage steps one notch away from it; a translucent cap is a
// tint, and the surface that tint was mixed for is the theme's window colour.
// Neither extreme is ever flat: a cap needs a surface to be an object on.
func previewStageColor(t theme.ThemeConfig) ui.Color {
	cap := t.KeyCapBackground
	if cap.A >= 128 {
		if colorLuminance(cap) > 0.5 {
			return cap.Mix(ui.RGBA(0, 0, 0, 255), 0.12)
		}
		return cap.Mix(ui.RGBA(255, 255, 255, 255), 0.12)
	}

	stage := t.WindowBackground
	stage.A = 255
	switch {
	case colorLuminance(stage) > 0.86:
		return stage.Mix(ui.RGBA(0, 0, 0, 255), 0.06)
	case colorLuminance(stage) < 0.10:
		return stage.Mix(ui.RGBA(255, 255, 255, 255), 0.06)
	}
	return stage
}

// renderShowcaseStage renders an Apple Cupertino style live preview stage,
// displaying realistic active HUD keycaps ([Ctrl] + [C] "Copy" and [⌘] + [V] "Paste")
// in the live active theme, physical model, typography, and color palette.
// Zero textual noise, zero double-frame nesting; clicking the stage replays the spring transition.
func (sc *SettingsController) renderShowcaseStage(c *ui.Context, cfg config.AppConfig, p palette) {
	currentTheme := theme.ThemeByName(cfg.ThemePreset)

	// The mouse simulation needs a little more air than the keycap strip.
	stageHeight := float32(92)
	if sc.activeTab == "Mouse" {
		stageHeight = 128
	}

	// The stage is the theme's own backdrop, never a fixed black box: a light
	// theme has to be previewed on the light surface it was designed for.
	stageColor := previewStageColor(currentTheme)

	stage := ui.Box(c.Key("showcase-canvas")).
		Fill().
		Height(stageHeight).
		Radius(12).
		Background(stageColor).
		Border(1, p.hairline).
		Padding(12).
		AlignItems(ui.Center).
		Justify(ui.Center).
		Shadow(0, 4, 16, 0, ui.RGBA(0, 0, 0, 0.35))

	if sc.activeTab == "Mouse" {
		// A quiet vertical gradient reads as a lit display instead of a flat box.
		stage = stage.Gradient(stageColor.Mix(ui.RGBA(255, 255, 255, 255), 0.06), stageColor.Mix(ui.RGBA(0, 0, 0, 255), 0.35), 135)
	}

	if stage.Clicked() {
		sc.previewNonce++
		if sc.window != nil {
			sc.window.Invalidate()
		}
	}

	stage.Children(func() {
		// The stage follows the active page: a mouse page that previews keycaps
		// teaches nothing. Other pages keep the keycap sample because every theme
		// change is visible there.
		if sc.activeTab == "Mouse" {
			sc.renderMouseShowcase(c, cfg)
			return
		}
		ui.Row(c).Gap(28).AlignItems(ui.Center).Children(func() {
			// Sample Chord 1: [ Ctrl ] + [ C ] -> "Copy" intent subtitle
			sampleChord1 := engine.DisplayChord{
				ID: fmt.Sprintf("preview-%d-1", sc.previewNonce),
				Keys: []engine.KeyCap{
					{Label: "Ctrl", IsModifier: true},
					{Label: "C", IsModifier: false},
				},
				Action: semantic.ActionInfo{Title: "Copy"},
			}
			renderChord(c, sampleChord1, currentTheme, cfg)

			// Sample Chord 2: [ ⌘ Cmd ] + [ V ] -> "Paste" intent subtitle
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

// renderMouseShowcase is the Mouse page's live overlay simulation. It shows the
// two things this page actually produces on screen: the cursor ripple and the
// button indicator. The silhouette is the same SVG the real overlay uses, so
// the preview cannot drift away from reality, and clicking the stage cycles the
// active button to show every accent the theme defines.
func (sc *SettingsController) renderMouseShowcase(c *ui.Context, cfg config.AppConfig) {
	mth := theme.MouseThemeByID(cfg.Mouse.ThemePreset)
	if cfg.Mouse.Color != "" {
		mth.LeftColor = cfg.Mouse.Color
		mth.RingColor = cfg.Mouse.Color
	}
	ringColor := ui.Hex(mth.RingColor)
	opacity := float32(cfg.Mouse.EffectiveOpacity())

	// The ripple follows the configured cursor size, so moving the slider has a
	// visible consequence instead of editing a number in the abstract.
	ring := float32(cfg.Mouse.Size)
	if ring <= 0 {
		ring = 50
	}
	ring = clampFloat(ring*0.55, 54, 84)
	core := ring * 0.46
	halo := ring * 1.22

	bodyH := float32(cfg.Mouse.IndicatorSize)
	if bodyH <= 0 {
		bodyH = 40
	}
	bodyH = clampFloat(bodyH*0.92, 32, 52)
	bodyW := bodyH * 0.85

	states := []struct {
		button hook.MouseButton
		wheel  int
	}{
		{hook.ButtonLeft, 0},
		{hook.ButtonRight, 0},
		{hook.ButtonMiddle, 1},
	}
	state := states[sc.previewNonce%len(states)]

	ui.Row(c).Gap(52).AlignItems(ui.Center).Opacity(opacity).Children(func() {
		// Ripple: outer halo, contact ring, and a tactile core holding the pointer.
		ui.Box(c.Key("mouse-preview-halo")).
			Size(halo, halo).
			Radius(halo/2).
			Border(1, ringColor.Alpha(0.10)).
			Background(ringColor.Alpha(0.025)).
			AlignItems(ui.Center).
			Justify(ui.Center).
			Children(func() {
				ui.Box(c.Key("mouse-preview-ring")).
					Size(ring, ring).
					Radius(ring/2).
					Border(1.6, ringColor.Alpha(0.45)).
					Background(ringColor.Alpha(0.06)).
					AlignItems(ui.Center).
					Justify(ui.Center).
					Children(func() {
						ui.Box(c.Key("mouse-preview-core")).
							Size(core, core).
							Radius(core/2).
							Border(2.2, ringColor).
							Background(ringColor.Alpha(0.18)).
							Shadow(0, 0, 20, 0, ringColor.Alpha(0.50)).
							AlignItems(ui.Center).
							Justify(ui.Center).
							Children(func() {
								ui.Image(c, cursorArrowSVG()).Size(core*0.40, core*0.52)
							})
					})
			})

		// The same floating pill the overlay draws beside the pointer.
		ui.Box(c.Key("mouse-preview-indicator")).
			Padding(10, 12).
			Radius(16).
			Background(ui.RGBA(10, 12, 17, 0.86)).
			Border(1, ui.RGBA(255, 255, 255, 0.16)).
			Shadow(0, 10, 28, 0, ui.RGBA(0, 0, 0, 0.55)).
			Children(func() {
				ui.Image(c, mouseBodySvg(state.button, state.wheel, mth)).Size(bodyW, bodyH)
			})
	})
}

// cursorArrowSVG is the neutral pointer inside the preview ripple. It stays
// white with a dark outline because the overlay sits on top of the real system
// cursor rather than replacing it.
func cursorArrowSVG() *ui.SVG {
	const svg = `<svg width="22" height="28" viewBox="0 0 22 28" fill="none" xmlns="http://www.w3.org/2000/svg">
		<path d="M3 2.2L19 15.4L11.9 16.5L15.8 25.1L11.7 26.9L7.9 18.2L3 23.2Z" fill="#FFFFFF" stroke="#0B0F16" stroke-width="1.6" stroke-linejoin="round"/>
	</svg>`
	return ui.MustParseSVG([]byte(svg))
}

// clampFloat keeps a preview dimension inside the stage without letting a
// hand-edited config push it off the canvas.
func clampFloat(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// renderAppTab is the application-behaviour page: language, capture policy,
// history, the global hotkey, and launch/tray behaviour. Visual design lives on
// the Keycap and Stage pages; this page only changes what the app does.
func (sc *SettingsController) renderAppTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("general_title")).FontSize(20).Bold().TextColor(p.text)

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
			row(i18n.T("ui_row_custom_filter_list"), i18n.TF("ui_custom_filter_list_desc", len(cfg.General.AllowedKeys)), func() {
				if ui.Button(c.Key("open-custom-filter-btn"), i18n.T("ui_btn_configure_matrix")).
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

		row(i18n.T("ui_row_injected_keys"), i18n.T("ui_row_injected_keys_desc"), func() {
			on := cfg.General.ShowInjectedKeys
			if ui.Switch(c.Key("toggle-injected-keys"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.General.ShowInjectedKeys = on })
			}
		})
	})

	group(i18n.T("history_title"), func() {
		row(i18n.T("history_title"), i18n.T("history_desc"), func() {
			on := cfg.General.ShowEventHistory
			if ui.Switch(c.Key("toggle-history"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.General.ShowEventHistory = on })
			}
		})

		row(i18n.T("ui_row_word_grouping"), i18n.T("ui_row_word_grouping_desc"), func() {
			on := cfg.General.WordGrouping
			if ui.Switch(c.Key("toggle-word-grouping"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.General.WordGrouping = on })
			}
		})

		row(i18n.T("ui_row_stay_duration"), i18n.T("ui_row_stay_duration_desc"), func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("dur-dec"), "−0.5s").Padding(4, 8).Radius(6).Clicked() && cfg.General.LingerDurationMs > 1000 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.General.LingerDurationMs -= 500 })
				}
				ui.Textf(c, i18n.T("ui_seconds"), float64(cfg.General.LingerDurationMs)/1000.0).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("dur-inc"), "+0.5s").Padding(4, 8).Radius(6).Clicked() && cfg.General.LingerDurationMs < 10000 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.General.LingerDurationMs += 500 })
				}
			})
		})

		if cfg.General.ShowEventHistory {
			row(i18n.T("ui_row_history_groups"), i18n.T("ui_row_history_groups_desc"), func() {
				ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
					if ui.Button(c.Key("hist-dec"), "−").Padding(4, 10).Radius(6).Clicked() && cfg.General.MaxHistory > 2 {
						sc.updateConfig(func(ac *config.AppConfig) { ac.General.MaxHistory-- })
					}
					ui.Textf(c, i18n.T("ui_groups_count"), cfg.General.MaxHistory).FontSize(13).Bold().TextColor(p.text)
					if ui.Button(c.Key("hist-inc"), "+").Padding(4, 10).Radius(6).Clicked() && cfg.General.MaxHistory < 12 {
						sc.updateConfig(func(ac *config.AppConfig) { ac.General.MaxHistory++ })
					}
				})
			})
		} else {
			row(i18n.T("ui_row_history_groups"), i18n.T("ui_row_history_groups_single_desc"), func() {
				ui.Box(c).Padding(5, 10).Radius(6).Background(p.sidebar).Border(1, p.hairline).Children(func() {
					ui.Text(c, i18n.T("ui_groups_fixed_one")).FontSize(12.5).Bold().TextColor(p.muted)
				})
			})
		}
	})

	group(i18n.T("ui_group_global_hotkey"), func() {
		row(i18n.T("toggle_shortcut"), i18n.T("toggle_shortcut_desc"), func() {
			buf := formatShortcut(cfg.General.ToggleShortcut)
			input := ui.TextInput(c.Key("toggle-hotkey-input"), &buf).
				ReadOnly(true).
				Placeholder(i18n.T("ui_hotkey_placeholder")).
				Width(220)
			input.OnTextKey(func(k ui.TextKey) bool {
				tokens := shortcutTokens(k.Mods, k.Key)
				if len(tokens) == 0 {
					return false
				}
				sc.updateConfig(func(ac *config.AppConfig) { ac.General.ToggleShortcut = tokens })
				return true
			})
		})
	})

	group(i18n.T("ui_group_app_behavior"), func() {
		row(i18n.T("ui_row_start_at_login"), i18n.T("ui_row_start_at_login_desc"), func() {
			on := cfg.App.StartAtLogin
			if ui.Switch(c.Key("toggle-start-at-login"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.App.StartAtLogin = on })
			}
		})

		row(i18n.T("ui_row_start_hidden"), i18n.T("ui_row_start_hidden_desc"), func() {
			on := cfg.App.StartHidden
			if ui.Switch(c.Key("toggle-start-hidden"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.App.StartHidden = on })
			}
		})

		row(i18n.T("ui_row_tray_click"), i18n.T("ui_row_tray_click_desc"), func() {
			actions := []string{"settings", "toggle"}
			labels := []string{i18n.T("tray_click_settings"), i18n.T("tray_click_toggle")}
			idx := 0
			if cfg.App.TrayClickAction == "toggle" {
				idx = 1
			}
			if ui.Segmented(c.Key("tray-click-action"), &idx, labels...).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.App.TrayClickAction = actions[idx] })
			}
		})
	})
}

// renderStageTab controls the overlay as a physical stage: which display it
// occupies, where it anchors, how it moves, and how its backdrop reads. The
// keycap's own look belongs to the Keycap page.
func (sc *SettingsController) renderStageTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
	stacked func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("appearance_title")).FontSize(20).Bold().TextColor(p.text)

	group(i18n.T("ui_group_display_position"), func() {
		row(i18n.T("ui_row_monitor"), i18n.T("ui_row_monitor_desc"), func() {
			displays := mygo.Screen.Displays()
			choices := []string{""}
			labels := []string{i18n.T("ui_display_primary")}
			for i, d := range displays {
				label := strings.TrimSpace(d.Label)
				if label == "" {
					label = fmt.Sprintf("%s %d", i18n.T("ui_display"), i+1)
				}
				if d.Primary {
					label += " · " + i18n.T("ui_display_primary")
				}
				choices = append(choices, strconv.FormatInt(d.ID, 10))
				labels = append(labels, label)
			}
			selected := cfg.Appearance.Monitor
			idx := 0
			for i, choice := range choices {
				if choice == selected {
					idx = i
					break
				}
			}
			if ui.Segmented(c.Key("monitor-seg"), &idx, labels...).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.Monitor = choices[idx] })
			}
		})

		row(i18n.T("alignment_title"), i18n.T("alignment_desc"), func() {
			sc.renderAlignmentSelector(c, cfg, p)
		})
	})

	group(i18n.T("ui_group_layout"), func() {
		row(i18n.T("direction_title"), i18n.T("ui_direction_desc"), func() {
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
	})

	group(i18n.T("ui_group_margins_opacity"), func() {
		row(i18n.T("ui_row_margin_link"), i18n.T("ui_row_margin_link_desc"), func() {
			linked := sc.marginLinked
			if ui.Switch(c.Key("toggle-margin-link"), &linked).Changed() {
				sc.marginLinked = linked
				if linked {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.MarginY = ac.Appearance.MarginX })
				}
			}
		})

		row(i18n.T("ui_row_margin_x"), "", func() {
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

		row(i18n.T("ui_row_margin_y"), "", func() {
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

		row(i18n.T("ui_row_overlay_opacity"), i18n.T("ui_row_overlay_opacity_desc"), func() {
			sc.renderPercentSlider(c, "overlay-opacity", cfg.Appearance.EffectiveOpacity(), config.MinOpacity*100, func(next float64) {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.Opacity = next })
			}, p)
		})
	})

	group(i18n.T("anim_title"), func() {
		row(i18n.T("ui_row_anim_style"), i18n.T("ui_row_anim_style_desc"), func() {
			anims := []string{"spring", "fade", "slide", "float", "zoom", "none"}
			labels := []string{
				i18n.T("ui_anim_spring"),
				i18n.T("ui_anim_fade"),
				i18n.T("ui_anim_slide"),
				i18n.T("ui_anim_float"),
				i18n.T("ui_anim_zoom"),
				i18n.T("ui_anim_none"),
			}
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
		row(i18n.T("ui_row_press_fx"), i18n.T("ui_row_press_fx_desc"), func() {
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
			row(i18n.T("ui_row_anim_duration"), i18n.T("ui_row_anim_duration_desc"), func() {
				ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
					if ui.Button(c.Key("anim-dur-dec"), "−50ms").Padding(4, 8).Radius(6).Clicked() && cfg.Appearance.AnimationDuration > 0.10 {
						sc.updateConfig(func(ac *config.AppConfig) {
							ac.Appearance.AnimationDuration = clampAnimDuration(ac.Appearance.AnimationDuration - 0.05)
						})
						sc.previewNonce++
					}
					ui.Textf(c, "%d ms", int(cfg.Appearance.AnimationDuration*1000)).
						FontSize(13).Bold().TextColor(p.text)
					if ui.Button(c.Key("anim-dur-inc"), "+50ms").Padding(4, 8).Radius(6).Clicked() && cfg.Appearance.AnimationDuration < 1.50 {
						sc.updateConfig(func(ac *config.AppConfig) {
							ac.Appearance.AnimationDuration = clampAnimDuration(ac.Appearance.AnimationDuration + 0.05)
						})
						sc.previewNonce++
					}
				})
			})
		}
	})

	group(i18n.T("ui_group_background"), func() {
		row(i18n.T("ui_row_background"), i18n.T("ui_row_background_desc"), func() {
			on := cfg.Background.Enabled
			if ui.Switch(c.Key("toggle-background"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Background.Enabled = on })
			}
		})
		if cfg.Background.Enabled {
			row(i18n.T("ui_row_background_color"), "", func() {
				col := parseHexSafe(cfg.Background.Color, ui.RGBA(0, 0, 0, 0.35))
				well := ui.ColorWell(c.Key("background-color"), &col)
				if well.Changed() {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Background.Color = colorHex(col) })
				}
			})
		}
	})
}

// renderKeycapTab is the design page: the preset theme sits first because it
// changes the whole keycap, then the physical model, colors, typography, and
// the small structural toggles that make up a cap.
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

	// A preset is the fastest way to change what the keycaps look like, so it
	// belongs at the top of the page. Twenty designs wait in a menu that opens
	// on request instead of standing in a wall the eye has to inventory first:
	// the trigger already wears the cap you have on.
	group(i18n.T("ui_group_preset_themes"), func() {
		row(i18n.T("ui_row_theme_preset"), i18n.T("ui_row_theme_preset_desc"), func() {
			sc.renderThemeSelect(c, cfg, p, currentTheme)
		})
	})

	// Physical model: the second half of "what does a keycap look like". Auto
	// is the first entry because it is the default, and the builds follow from
	// the shallowest to the deepest, so the list reads as a spectrum.
	group(i18n.T("keycap_model"), func() {
		row(i18n.T("ui_row_physical_style"), i18n.T("ui_row_physical_style_desc"), func() {
			sc.renderModelSelect(c, cfg, p, currentTheme)
		})
	})

	// Color palette: the same cap, repainted. Fifteen swatches in a dropdown
	// force the user to try each one; a grid lets the eye choose.
	group(i18n.T("ui_group_color_palettes"), func() {
		sc.renderPaletteGallery(c, cfg, p, currentTheme)

		row(i18n.T("ui_row_gradient"), i18n.T("ui_row_gradient_desc"), func() {
			on := cfg.Color.UseGradient
			if ui.Switch(c.Key("toggle-gradient"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Color.UseGradient = on })
			}
		})
	})

	// Modifier accent is a color behaviour, so it lives beside the palette
	// instead of masquerading as a show/hide element.
	group(i18n.T("ui_group_modifier_colors"), func() {
		row(i18n.T("ui_row_mod_highlight"), i18n.T("ui_row_mod_highlight_desc"), func() {
			on := cfg.Modifier.Highlight
			if ui.Switch(c.Key("toggle-mod-hl"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Modifier.Highlight = on })
			}
		})

		if cfg.Modifier.Highlight {
			modColorRow := func(label, key string, pick func(*config.AppConfig) *string, fallback ui.Color) {
				row(label, "", func() {
					col := parseHexSafe(*pick(&cfg), fallback)
					well := ui.ColorWell(c.Key(key), &col)
					if well.Changed() {
						hex := colorHex(col)
						sc.updateConfig(func(ac *config.AppConfig) { *pick(ac) = hex })
					}
				})
			}
			modColorRow(i18n.T("ui_row_mod_bg"), "mod-bg-color", func(ac *config.AppConfig) *string { return &ac.Modifier.Color }, currentTheme.KeyCapBackground)
			modColorRow(i18n.T("ui_row_mod_border"), "mod-border-color", func(ac *config.AppConfig) *string { return &ac.Modifier.BorderColor }, currentTheme.KeyCapBorder)
			modColorRow(i18n.T("ui_row_mod_text"), "mod-text-color", func(ac *config.AppConfig) *string { return &ac.Modifier.TextColor }, currentTheme.KeyCapTextColor)
		}
	})

	// Group 3: Typography & Text Settings
	group(i18n.T("text_formatting"), func() {
		row(i18n.T("text_caps"), "", func() {
			caps := []string{"uppercase", "capitalize", "lowercase"}
			labels := []string{i18n.T("caps_upper"), i18n.T("caps_cap"), i18n.T("caps_lower")}
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
			labels := []string{i18n.T("variant_text"), i18n.T("variant_short"), i18n.T("variant_icon")}
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

		row(i18n.T("ui_row_font_size"), "", func() {
			size := cfg.Text.Size
			suffix := ""
			if size <= 0 {
				size = int(effectiveFontSize(currentTheme, cfg) + 0.5)
				suffix = i18n.T("ui_suffix_theme")
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

		row(i18n.T("ui_row_text_alignment"), i18n.T("ui_row_text_alignment_desc"), func() {
			aligns := []config.Alignment{config.AlignCenterLeft, config.AlignCenter, config.AlignCenterRight}
			labels := []string{i18n.T("ui_align_left"), i18n.T("ui_align_center"), i18n.T("ui_align_right")}
			idx := 1
			for i, al := range aligns {
				if al == cfg.Text.Alignment {
					idx = i
				}
			}
			if ui.Segmented(c.Key("text-align-seg"), &idx, labels...).Changed() {
				chosen := aligns[idx]
				sc.updateConfig(func(ac *config.AppConfig) { ac.Text.Alignment = chosen })
			}
		})

		row(i18n.T("ui_row_icon_alignment"), i18n.T("ui_row_icon_alignment_desc"), func() {
			labels := []string{i18n.T("ui_align_icon_before"), i18n.T("ui_align_icon_after")}
			idx := 0
			if cfg.Layout.IconAlignment == "flex-end" {
				idx = 1
			}
			if ui.Segmented(c.Key("icon-align-seg"), &idx, labels...).Changed() {
				value := "flex-start"
				if idx == 1 {
					value = "flex-end"
				}
				sc.updateConfig(func(ac *config.AppConfig) { ac.Layout.IconAlignment = value })
			}
		})
	})

	// Group 4: Layout & Elements
	group(i18n.T("display_elements"), func() {
		row(i18n.T("ui_row_key_icon"), i18n.T("ui_row_key_icon_desc"), func() {
			on := cfg.Layout.ShowIcon
			if ui.Switch(c.Key("toggle-icon"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Layout.ShowIcon = on })
			}
		})

		row(i18n.T("ui_row_special_symbol"), i18n.T("ui_row_special_symbol_desc"), func() {
			on := cfg.Layout.ShowSymbol
			if ui.Switch(c.Key("toggle-symbol"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Layout.ShowSymbol = on })
			}
		})

		row(i18n.T("ui_row_press_count"), i18n.T("ui_row_press_count_desc"), func() {
			on := cfg.Layout.ShowPressCount
			if ui.Switch(c.Key("toggle-press-count"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Layout.ShowPressCount = on })
			}
		})
	})

	// Group 5: Border & Corner Radius
	group(i18n.T("border_title"), func() {
		row(i18n.T("ui_row_keycap_border"), "", func() {
			on := cfg.Border.Enabled
			if ui.Switch(c.Key("toggle-border"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Border.Enabled = on })
			}
		})

		row(i18n.T("ui_row_border_width"), "", func() {
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

		row(i18n.T("ui_row_border_radius"), "", func() {
			radius := cfg.Border.Radius
			suffix := ""
			if radius <= 0 {
				radius = float64(effectiveRadius(currentTheme, cfg))
				suffix = i18n.T("ui_suffix_theme")
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
}

// renderMouseTab is the cursor page. Themes and custom color first, then the
// click ripple, the overlay's opacity, and the button indicator; drag threshold
// stays with gestures because that is the behaviour it tunes.
func (sc *SettingsController) renderMouseTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("mouse_title")).FontSize(20).Bold().TextColor(p.text)

	// Group 1: mouse theme and the optional custom accent
	group(i18n.T("ui_group_mouse_theme"), func() {
		currentMTheme := theme.MouseThemeByID(cfg.Mouse.ThemePreset)
		mChoice := currentMTheme.ID

		row(i18n.T("ui_row_preset_theme"), i18n.T("ui_row_preset_theme_desc"), func() {
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
					ui.Text(c, mouseThemeDisplayName(currentMTheme)).FontSize(13).Bold().TextColor(p.text)
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

				ui.Scroll(c).MaxHeight(368).Children(func() {
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
							ui.Text(c, mouseThemeDisplayName(themeItem)).FontSize(12.5).TextColor(p.text).Grow(1)
							if themeItem.ID == currentMTheme.ID {
								ui.Text(c, "✓").FontSize(13).Bold().TextColor(p.accent)
							}
						})
					}
				})
			})

			if sel.Trigger.Changed() {
				targetMth := theme.MouseThemeByID(mChoice)
				// An empty color means "follow the selected theme"; the old code
				// copied the pale green seed into every theme and silently
				// overrode it.
				sc.updateConfig(func(ac *config.AppConfig) {
					ac.Mouse.ThemePreset = targetMth.ID
					ac.Mouse.Color = ""
				})
			}
		})

		row(i18n.T("ui_row_mouse_color"), i18n.T("ui_row_mouse_color_desc"), func() {
			base := ui.Hex(currentMTheme.LeftColor)
			if cfg.Mouse.Color != "" {
				base = parseHexSafe(cfg.Mouse.Color, base)
			}
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				col := base
				well := ui.ColorWell(c.Key("mouse-custom-color"), &col)
				if well.Changed() {
					hex := colorHex(col)
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.Color = hex })
				}
				if cfg.Mouse.Color != "" {
					if ui.Button(c.Key("mouse-follow-theme"), i18n.T("ui_row_follow_theme")).Padding(5, 10).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.text).FontSize(11.5).Clicked() {
						sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.Color = "" })
					}
				}
			})
		})
	})

	// Group 2: click ripple behaviour
	group(i18n.T("mouse_click_ripple"), func() {
		row(i18n.T("ui_row_click_ring"), i18n.T("ui_row_click_ring_desc"), func() {
			on := cfg.Mouse.ShowClicks
			if ui.Switch(c.Key("toggle-mouse-clicks"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.ShowClicks = on })
			}
		})

		row(i18n.T("ui_row_always_ring"), i18n.T("ui_row_always_ring_desc"), func() {
			on := cfg.Mouse.KeepHighlight
			if ui.Switch(c.Key("toggle-always-hl"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.KeepHighlight = on })
			}
		})

		row(i18n.T("ui_row_click_linger"), i18n.T("ui_row_click_linger_desc"), func() {
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

	// Group 3: mouse overlay appearance, symmetric with the keycap opacity
	group(i18n.T("ui_group_mouse_appearance"), func() {
		row(i18n.T("ui_row_ring_size"), i18n.T("ui_row_ring_size_desc"), func() {
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

		row(i18n.T("ui_row_mouse_opacity"), i18n.T("ui_row_mouse_opacity_desc"), func() {
			sc.renderPercentSlider(c, "mouse-opacity", cfg.Mouse.EffectiveOpacity(), config.MinOpacity*100, func(next float64) {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.Opacity = next })
			}, p)
		})
	})

	// Group 4: button and wheel indicator
	group(i18n.T("mouse_indicator"), func() {
		row(i18n.T("ui_row_button_indicator"), i18n.T("ui_row_button_indicator_desc"), func() {
			on := cfg.Mouse.ShowIndicator
			if ui.Switch(c.Key("toggle-indicator"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.ShowIndicator = on })
			}
		})

		row(i18n.T("ui_row_keep_indicator"), i18n.T("ui_row_keep_indicator_desc"), func() {
			on := cfg.Mouse.KeepIndicator
			if ui.Switch(c.Key("toggle-keep-ind"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.KeepIndicator = on })
			}
		})

		row(i18n.T("ui_row_indicator_size"), "", func() {
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

		row(i18n.T("ui_row_indicator_offset_x"), "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("iox-dec"), "−5").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.IndicatorOffsetX > 0 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.IndicatorOffsetX -= 5 })
				}
				ui.Textf(c, "%d px", cfg.Mouse.IndicatorOffsetX).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("iox-inc"), "+5").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.IndicatorOffsetX < 300 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.IndicatorOffsetX += 5 })
				}
			})
		})

		row(i18n.T("ui_row_indicator_offset_y"), "", func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("ioy-dec"), "−5").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.IndicatorOffsetY > 0 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.IndicatorOffsetY -= 5 })
				}
				ui.Textf(c, "%d px", cfg.Mouse.IndicatorOffsetY).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("ioy-inc"), "+5").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.IndicatorOffsetY < 300 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.IndicatorOffsetY += 5 })
				}
			})
		})
	})

	// Group 5: gesture tuning
	group(i18n.T("ui_group_gesture_drag"), func() {
		row(i18n.T("ui_row_drag_threshold"), i18n.T("ui_row_drag_threshold_desc"), func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				if ui.Button(c.Key("dt-dec"), "−10").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.DragThreshold > 10 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.DragThreshold -= 10 })
				}
				ui.Textf(c, "%d px", cfg.Mouse.DragThreshold).FontSize(13).Bold().TextColor(p.text)
				if ui.Button(c.Key("dt-inc"), "+10").Padding(4, 8).Radius(6).Clicked() && cfg.Mouse.DragThreshold < 150 {
					sc.updateConfig(func(ac *config.AppConfig) { ac.Mouse.DragThreshold += 10 })
				}
			})
		})
	})
}

// renderAboutTab keeps the product story and the system facts in separate
// groups, and turns "configuration management" into actions that actually
// manage the file instead of a second copy of the factory reset.
func (sc *SettingsController) renderAboutTab(
	c *ui.Context,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("tab_about")).FontSize(20).Bold().TextColor(p.text)

	group(i18n.T("ui_group_system_info"), func() {
		row(i18n.T("app_name"), i18n.T("about_desc"), func() {})
		row(i18n.T("ui_row_version"), i18n.T("ui_row_version_desc"), func() {
			ui.Box(c).Padding(5, 10).Radius(6).Background(p.sidebar).Border(1, p.hairline).Children(func() {
				ui.Text(c, version.Display()).FontSize(12.5).Bold().TextColor(p.text)
			})
		})
		row(i18n.T("ui_row_renderer"), i18n.T("renderer_architecture"), func() {})
		row(i18n.T("ui_row_engine"), i18n.T("engine_architecture"), func() {})
	})

	group(i18n.T("ui_group_design_story"), func() {
		row(i18n.T("ui_row_inspiration"), i18n.T("original_credits"), func() {})
		row(i18n.T("ui_row_philosophy"), i18n.T("about_philosophy"), func() {})
	})

	group(i18n.T("ui_group_config_management"), func() {
		row(i18n.T("ui_row_open_config"), i18n.T("ui_row_open_config_desc"), func() {
			if ui.Button(c.Key("btn-open-config"), i18n.T("ui_btn_open")).Padding(6, 14).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.text).FontSize(12).Clicked() {
				if err := openPath(sc.store.Path()); err != nil {
					sc.configNotice = i18n.T("ui_config_open_failed")
				}
			}
		})

		row(i18n.T("ui_row_open_config_dir"), i18n.T("ui_row_open_config_dir_desc"), func() {
			if ui.Button(c.Key("btn-open-config-dir"), i18n.T("ui_btn_open")).Padding(6, 14).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.text).FontSize(12).Clicked() {
				if err := openPath(filepath.Dir(sc.store.Path())); err != nil {
					sc.configNotice = i18n.T("ui_config_open_failed")
				}
			}
		})

		row(i18n.T("ui_row_copy_config_path"), i18n.T("ui_row_copy_config_path_desc"), func() {
			if ui.Button(c.Key("btn-copy-config-path"), i18n.T("ui_btn_copy")).Padding(6, 14).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.text).FontSize(12).Clicked() {
				mygo.Clipboard.WriteText(sc.store.Path())
				sc.configNotice = i18n.T("ui_config_copied")
			}
		})

		row(i18n.T("ui_row_export_config"), i18n.T("ui_row_export_config_desc"), func() {
			if ui.Button(c.Key("btn-export-config"), i18n.T("ui_btn_copy")).Padding(6, 14).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.text).FontSize(12).Clicked() {
				data, err := config.Encode(sc.store.Get())
				if err != nil {
					sc.configNotice = i18n.T("ui_config_export_failed")
					return
				}
				mygo.Clipboard.WriteText(string(data))
				sc.configNotice = i18n.T("ui_config_exported")
			}
		})

		row(i18n.T("ui_row_import_config"), i18n.T("ui_row_import_config_desc"), func() {
			if ui.Button(c.Key("btn-import-config"), i18n.T("ui_btn_import")).Padding(6, 14).Radius(6).Background(p.accent).TextColor(p.onAccent).FontSize(12).Bold().Clicked() {
				text := mygo.Clipboard.ReadText()
				next, err := config.Decode([]byte(text))
				if err != nil {
					sc.configNotice = i18n.T("ui_config_import_failed")
					return
				}
				if err := sc.store.Update(func(ac *config.AppConfig) { *ac = next }); err != nil {
					sc.configNotice = i18n.T("ui_config_import_failed")
					return
				}
				sc.configNotice = i18n.T("ui_config_imported")
			}
		})

		if sc.configNotice != "" {
			ui.Row(c).Padding(10, 16).Children(func() {
				ui.Text(c, sc.configNotice).FontSize(12).TextColor(p.muted)
			})
		}
	})
}

// renderIntelligenceTab collects everything that interprets the keystroke:
// context awareness, subtitles, application badges, custom shortcut rules, and
// the privacy list that decides which apps may contribute context.
func (sc *SettingsController) renderIntelligenceTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("ctx_group_title")).FontSize(20).Bold().TextColor(p.text)

	group(i18n.T("ctx_group_title"), func() {
		row(i18n.T("ctx_smart_title"), i18n.T("ctx_smart_desc"), func() {
			on := cfg.SemanticContext
			if ui.Switch(c.Key("toggle-semantic-context"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.SemanticContext = on })
			}
		})

		row(i18n.T("ctx_sub_title"), i18n.T("ctx_sub_desc"), func() {
			on := cfg.SemanticSub
			if ui.Switch(c.Key("toggle-semantic-sub"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.SemanticSub = on })
			}
		})

		row(i18n.T("ctx_badge_title"), i18n.T("ctx_badge_desc"), func() {
			on := cfg.ShowAppBadge
			if ui.Switch(c.Key("toggle-app-badge"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.ShowAppBadge = on })
			}
		})

		appCtx := context.AppContext{}
		if sc.contextSource != nil {
			appCtx = sc.contextSource()
		}
		appLabel := appCtx.DisplayName()
		appDesc := appCtx.Exe
		if appCtx.Title != "" {
			if appDesc != "" {
				appDesc += " · "
			}
			appDesc += ellipsize(appCtx.Title, 42)
		}
		if appDesc == "" {
			appDesc = i18n.T("ctx_no_app_desc")
		}
		if appLabel == "" {
			appLabel = i18n.T("ctx_unknown_app")
		}
		row(i18n.T("ctx_current_app"), appDesc, func() {
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				if appCtx.Category != context.CategoryUnknown {
					ui.Box(c).Padding(3, 7).Radius(5).Background(p.sidebar).Border(1, p.hairline).Children(func() {
						ui.Text(c, appCtx.Category.DisplayName()).FontSize(10.5).Bold().TextColor(p.muted)
					})
				}
				ui.Box(c).Padding(5, 10).Radius(6).Background(p.sidebar).Border(1, p.hairline).Children(func() {
					ui.Text(c, appLabel).FontSize(12.5).Bold().SingleLine().TextColor(p.text)
				})
			})
		})
	})

	group(i18n.T("ui_group_custom_shortcuts"), func() {
		row(i18n.T("ui_row_custom_shortcuts"), i18n.T("ui_custom_shortcuts_desc"), func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				ui.Box(c).Padding(5, 10).Radius(6).Background(p.sidebar).Border(1, p.hairline).Children(func() {
					ui.Textf(c, "%d", len(cfg.CustomShortcuts)).FontSize(12.5).Bold().TextColor(p.text)
				})
				if ui.Button(c.Key("btn-edit-custom-shortcuts"), i18n.T("ui_btn_edit_custom_shortcuts")).Padding(6, 14).Radius(6).Background(p.accent).TextColor(p.onAccent).FontSize(12).Bold().Clicked() {
					sc.customDrafts = append([]config.CustomShortcut(nil), cfg.CustomShortcuts...)
					sc.showCustomShortcutModal = true
				}
			})
		})
	})

	group(i18n.T("ui_group_privacy"), func() {
		row(i18n.T("ui_row_privacy_shield"), i18n.T("ui_row_privacy_shield_desc"), func() {
			on := cfg.General.PrivacyMode
			if ui.Switch(c.Key("toggle-privacy-mode"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.General.PrivacyMode = on })
			}
		})

		if !sc.excludedInitialized {
			sc.excludedDraft = strings.Join(cfg.App.ExcludedApps, ", ")
			sc.excludedInitialized = true
		}
		row(i18n.T("ui_row_excluded_apps"), i18n.T("ui_row_excluded_apps_desc"), func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				ui.TextInput(c.Key("excluded-apps-input"), &sc.excludedDraft).Width(260).Placeholder(i18n.T("ui_excluded_placeholder"))
				if ui.Button(c.Key("btn-apply-exclusions"), i18n.T("ui_btn_apply_exclusions")).Padding(6, 12).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.text).FontSize(12).Clicked() {
					entries := splitList(sc.excludedDraft)
					sc.updateConfig(func(ac *config.AppConfig) { ac.App.ExcludedApps = entries })
				}
			})
		})
	})
}

// renderSoundTab is a single-purpose page: audition the material of every
// profile, then choose one. Previewing does not depend on the feature being
// enabled, so users can compare sounds before turning audio on.
func (sc *SettingsController) renderSoundTab(
	c *ui.Context,
	cfg config.AppConfig,
	p palette,
	group func(title string, rows func()),
	row func(label, desc string, control func()),
) {
	ui.Text(c, i18n.T("tab_sound")).FontSize(20).Bold().TextColor(p.text)

	group(i18n.T("ui_group_audio"), func() {
		row(i18n.T("ui_row_keystroke_audio"), i18n.T("ui_row_keystroke_audio_desc"), func() {
			on := cfg.TactileAudio
			if ui.Switch(c.Key("toggle-tactile-audio"), &on).Changed() {
				sc.updateConfig(func(ac *config.AppConfig) { ac.TactileAudio = on })
			}
		})

		if cfg.TactileAudio {
			row(i18n.T("ui_row_volume"), i18n.T("ui_row_volume_desc"), func() {
				sc.renderPercentSlider(c, "audio-volume", cfg.EffectiveVolume(), 0, func(next float64) {
					sc.updateConfig(func(ac *config.AppConfig) { ac.AudioVolume = next })
				}, p)
			})
		}

		row(i18n.T("sound_profile"), i18n.T("sound_profile_desc"), func() {
			options := soundProfileOptions()
			curProf := cfg.AudioProfile
			if curProf == "" {
				curProf = "auto"
			}
			curLabel := i18n.T("profile_auto")
			for _, opt := range options {
				if opt.id == curProf {
					curLabel = opt.label
					break
				}
			}

			selectedProf := curProf
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				sel := ui.SelectBase(c, &selectedProf)
				sel.Trigger.
					Height(32).
					PaddingX(12).
					Gap(8).
					Radius(6).
					Background(p.sidebar).
					Border(1, p.hairline).
					AlignItems(ui.Center).
					Children(func() {
						ui.Text(c, curLabel).FontSize(12.5).Bold().TextColor(p.text)
						ui.Icon(c, iconChevronDown).Size(12, 12).TextColor(p.muted)
					})

				sel.Popup(func(panel ui.Element) {
					panel.
						Width(240).
						Background(p.card).
						Border(1, p.hairline).
						Radius(8).
						Padding(4).
						Shadow(0, 10, 24, 0, ui.RGBA(0, 0, 0, 0.45))

					ui.Scroll(c).MaxHeight(320).Children(func() {
						for _, opt := range options {
							itemOpt := opt
							item := sel.Item(itemOpt.id).
								Height(32).
								PaddingX(10).
								Gap(8).
								Radius(6).
								AlignItems(ui.Center)
							if item.Highlighted() {
								item.Background(p.press)
							} else if itemOpt.id == curProf {
								item.Background(p.hover)
							}
							item.Children(func() {
								ui.Text(c, itemOpt.label).FontSize(12.5).TextColor(p.text).Grow(1)
								if itemOpt.id == curProf {
									ui.Text(c, "✓").FontSize(12.5).Bold().TextColor(p.accent)
								}
							})
						}
					})
				})

				if sel.Trigger.Changed() {
					chosen := selectedProf
					sc.updateConfig(func(ac *config.AppConfig) { ac.AudioProfile = chosen })
				}

				sc.renderAuditionButton(c, p, selectedProf)
			})
		})
	})
}

// renderAuditionButton is the sound audition. A click has to be felt even when
// the speakers are muted or the room is loud, so the pill fills with the accent
// and its meter rides the sound for four hundred milliseconds: the button
// answers, and it answers in the same shape as the keycap stage above it.
func (sc *SettingsController) renderAuditionButton(c *ui.Context, p palette, profile string) {
	btn := ui.ButtonBase(c.Key("btn-preview-sound")).
		Height(32).
		PaddingX(11).
		Gap(7).
		Radius(6).
		AlignItems(ui.Center).
		FocusRing(false)

	// Animate keeps frames coming while it moves, which is what makes the
	// meter advance at all: the window is otherwise event-driven.
	pulse := btn.Animate("audio-pulse", sc.audioPulse, 200*time.Millisecond)
	if btn.Arrived("audio-pulse") && sc.audioPulse > 0 {
		sc.audioPulse = 0
	}

	live := pulse > 0.02
	bg := p.sidebar
	border := p.hairline
	if btn.Hovered() {
		bg = p.hover
	}
	if live {
		bg, border = p.accent, p.accent
	}
	btn.Background(btn.AnimateColor("audio-bg", bg, ui.Animation{Duration: 140 * time.Millisecond})).
		Border(1, btn.AnimateColor("audio-border", border, ui.Animation{Duration: 140 * time.Millisecond}))

	fg := p.text
	if pulse > 0.6 {
		fg = p.onAccent
	}

	btn.Children(func() {
		// Three bars are the icon: at rest they read as a level meter, and
		// while the sound plays they move at about fourteen radians a second,
		// out of phase, so the group reads as sound rather than as a pulse.
		ui.Row(c).AlignItems(ui.End).Gap(2).Height(13).Children(func() {
			phase := float32(time.Since(sc.audioPulseAt).Seconds()) * 14
			for i := 0; i < 3; i++ {
				height := float32(4)
				if live {
					wave := 0.55 + 0.45*float32(math.Sin(float64(phase+float32(i)*2.1)))
					height = 3 + 9*pulse*wave
				}
				ui.Box(c).Width(2.5).Height(height).Radius(1.25).Background(fg)
			}
		})
		ui.Text(c, i18n.T("ui_btn_preview")).FontSize(12).TextColor(fg)
	})

	if btn.Clicked() {
		sc.audioPulse = 1
		sc.audioPulseAt = time.Now()
		if sc.audioPreview != nil {
			sc.audioPreview(profile)
		}
	}
}

// renderCustomShortcutModal is the visual editor for config.CustomShortcut.
// Previously this list could only be changed by hand-editing config.json.
func (sc *SettingsController) renderCustomShortcutModal(c *ui.Context, cfg config.AppConfig, p palette) {
	issues := shortcutDraftIssues(sc.customDrafts)

	ui.Box(c.Key("shortcut-modal-backdrop")).
		Absolute().
		Fill().
		Background(ui.RGBA(0, 0, 0, 0.55)).
		AlignItems(ui.Center).
		Justify(ui.Center).
		Children(func() {
			ui.Column(c.Key("shortcut-modal-card")).
				Width(720).
				MaxHeight(560).
				Background(p.card).
				Border(1, p.hairline).
				Radius(14).
				Padding(20).
				Gap(12).
				Shadow(0, 16, 40, 0, ui.RGBA(0, 0, 0, 0.50)).
				Children(func() {
					ui.Row(c).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
						ui.Column(c).Gap(2).Children(func() {
							ui.Text(c, i18n.T("ui_modal_custom_shortcuts_title")).FontSize(16).Bold().TextColor(p.text)
							ui.Text(c, i18n.T("ui_modal_custom_shortcuts_desc")).FontSize(11.5).TextColor(p.muted)
						})
						ui.Row(c).Gap(8).Children(func() {
							if ui.Button(c.Key("close-shortcut-modal"), i18n.T("ui_btn_cancel")).Padding(6, 14).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.text).Clicked() {
								sc.showCustomShortcutModal = false
							}
							save := ui.Button(c.Key("save-shortcut-modal"), i18n.T("ui_btn_save")).
								Padding(6, 14).
								Radius(6).
								Bold()
							if len(issues) > 0 {
								// A disabled-looking button is more honest than one
								// that silently refuses to save.
								save.Background(p.sidebar).Border(1, p.hairline).TextColor(p.muted)
							} else {
								save.Background(p.accent).TextColor(p.onAccent)
							}
							if save.Clicked() && len(issues) == 0 {
								cleaned := make([]config.CustomShortcut, 0, len(sc.customDrafts))
								for _, draft := range sc.customDrafts {
									if strings.TrimSpace(draft.Chord) == "" || strings.TrimSpace(draft.Label) == "" {
										continue
									}
									cleaned = append(cleaned, draft)
								}
								sc.updateConfig(func(ac *config.AppConfig) { ac.CustomShortcuts = cleaned })
								sc.showCustomShortcutModal = false
							}
						})
					})

					ui.Row(c).Gap(8).Children(func() {
						ui.Text(c, i18n.T("ui_column_app")).FontSize(11).Bold().TextColor(p.muted).Width(120)
						ui.Text(c, i18n.T("ui_column_chord")).FontSize(11).Bold().TextColor(p.muted).Width(130)
						ui.Text(c, i18n.T("ui_column_label")).FontSize(11).Bold().TextColor(p.muted).Width(140)
						ui.Text(c, i18n.T("ui_column_description")).FontSize(11).Bold().TextColor(p.muted).Grow(1)
						ui.Box(c).Width(24)
					})

					ui.Scroll(c).MaxHeight(340).Grow(1).Children(func() {
						for i := range sc.customDrafts {
							draft := &sc.customDrafts[i]
							ui.Row(c.Key(fmt.Sprintf("shortcut-row-%d", i))).Gap(8).Padding(2, 0).Children(func() {
								ui.TextInput(c.Key(fmt.Sprintf("shortcut-app-%d", i)), &draft.App).Width(120).Placeholder(i18n.T("ui_column_app_placeholder"))
								chord := ui.TextInput(c.Key(fmt.Sprintf("shortcut-chord-%d", i)), &draft.Chord).
									Width(130).
									ReadOnly(true).
									Placeholder(i18n.T("ui_hotkey_placeholder"))
								chord.OnTextKey(func(k ui.TextKey) bool {
									return recordShortcutChord(draft, k)
								})
								ui.TextInput(c.Key(fmt.Sprintf("shortcut-label-%d", i)), &draft.Label).Width(140).Placeholder(i18n.T("ui_column_label_placeholder"))
								ui.TextInput(c.Key(fmt.Sprintf("shortcut-desc-%d", i)), &draft.Description).Grow(1).Placeholder(i18n.T("ui_column_description_placeholder"))
								if ui.Button(c.Key(fmt.Sprintf("delete-shortcut-%d", i)), "×").Padding(3, 8).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.muted).Clicked() {
									sc.customDrafts = append(sc.customDrafts[:i], sc.customDrafts[i+1:]...)
								}
							})
						}
					})

					if len(issues) > 0 {
						ui.Textf(c, i18n.T("ui_shortcut_conflict_summary"), len(issues)).
							FontSize(11.5).
							TextColor(ui.Hex("#ff6b6b"))
					}

					ui.Row(c).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
						if ui.Button(c.Key("add-shortcut"), "+ "+i18n.T("ui_btn_add")).Padding(5, 12).Radius(6).Background(p.sidebar).Border(1, p.hairline).TextColor(p.text).FontSize(12).Clicked() {
							sc.customDrafts = append(sc.customDrafts, config.CustomShortcut{})
						}
						ui.Textf(c, "%d", len(sc.customDrafts)).FontSize(11.5).TextColor(p.muted)
					})
				})
		})
}

// soundProfileOptions is the single source of truth for the Sound page and
// keeps the order stable between the trigger label and the popup.
func soundProfileOptions() []struct {
	id    string
	label string
} {
	return []struct {
		id    string
		label string
	}{
		{"auto", i18n.T("profile_auto")},
		{"thock", i18n.T("profile_thock")},
		{"click", i18n.T("profile_click")},
		{"scissor", i18n.T("profile_scissor")},
		{"glass", i18n.T("profile_glass")},
		{"typewriter", i18n.T("profile_typewriter")},
		{"hinoki", i18n.T("profile_hinoki")},
		{"pebble", i18n.T("profile_pebble")},
		{"paper", i18n.T("profile_paper")},
		{"felt", i18n.T("profile_felt")},
		{"droplet", i18n.T("profile_droplet")},
		{"crystal", i18n.T("profile_crystal")},
		{"chip", i18n.T("profile_chip")},
		{"shutter", i18n.T("profile_shutter")},
	}
}

// colorHex serialises a ColorWell value back into the config's #RRGGBB or
// #RRGGBBAA form.
func colorHex(c ui.Color) string {
	if c.A >= 255 {
		return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}

// keyToken maps a UI key to the canonical token names understood by the hook
// matcher in main.go.
func keyToken(k ui.Key) string {
	fkeys := []ui.Key{ui.KeyF1, ui.KeyF2, ui.KeyF3, ui.KeyF4, ui.KeyF5, ui.KeyF6, ui.KeyF7, ui.KeyF8, ui.KeyF9, ui.KeyF10, ui.KeyF11, ui.KeyF12}
	for i, fk := range fkeys {
		if k == fk {
			return fmt.Sprintf("F%d", i+1)
		}
	}
	if k >= ui.KeyA && k <= ui.KeyZ {
		return string(rune('A' + int(k-ui.KeyA)))
	}
	if k >= ui.Key0 && k <= ui.Key9 {
		return string(rune('0' + int(k-ui.Key0)))
	}
	switch k {
	case ui.KeyEnter:
		return "Enter"
	case ui.KeyEscape:
		return "Esc"
	case ui.KeyBackspace:
		return "Backspace"
	case ui.KeyTab:
		return "Tab"
	case ui.KeySpace:
		return "Space"
	case ui.KeyDelete:
		return "Delete"
	case ui.KeyInsert:
		return "Insert"
	case ui.KeyHome:
		return "Home"
	case ui.KeyEnd:
		return "End"
	case ui.KeyPageUp:
		return "PageUp"
	case ui.KeyPageDown:
		return "PageDown"
	case ui.KeyLeft:
		return "Left"
	case ui.KeyRight:
		return "Right"
	case ui.KeyUp:
		return "Up"
	case ui.KeyDown:
		return "Down"
	case ui.KeyMinus:
		return "-"
	case ui.KeyEqual:
		return "="
	case ui.KeyComma:
		return ","
	case ui.KeyPeriod:
		return "."
	case ui.KeySlash:
		return "/"
	case ui.KeySemicolon:
		return ";"
	case ui.KeyQuote:
		return "'"
	case ui.KeyBracketLeft:
		return "["
	case ui.KeyBracketRight:
		return "]"
	case ui.KeyBackslash:
		return "\\"
	case ui.KeyBackquote:
		return "`"
	}
	return ""
}

// shortcutTokens converts a captured key press into the persisted token list.
// A modifier-only press has no key token, so it is ignored rather than saving
// a hotkey that can never fire.
func shortcutTokens(mods ui.Modifiers, key ui.Key) []string {
	token := keyToken(key)
	if token == "" {
		return nil
	}
	out := make([]string, 0, 5)
	if mods&ui.Ctrl != 0 {
		out = append(out, "ControlLeft")
	}
	if mods&ui.Shift != 0 {
		out = append(out, "ShiftLeft")
	}
	if mods&ui.Alt != 0 {
		out = append(out, "AltLeft")
	}
	if mods&ui.Super != 0 {
		out = append(out, "MetaLeft")
	}
	return append(out, token)
}

// recordShortcutChord turns one captured key press into the draft's chord.
// Backspace or Delete clears it, and a modifier-only press is ignored because
// it cannot become a chord on its own.
func recordShortcutChord(draft *config.CustomShortcut, k ui.TextKey) bool {
	if k.Key == ui.KeyBackspace || k.Key == ui.KeyDelete {
		draft.Chord = ""
		return true
	}
	tokens := shortcutTokens(k.Mods, k.Key)
	if len(tokens) == 0 {
		return false
	}
	draft.Chord = formatShortcut(tokens)
	return true
}

// shortcutDraftIssues lists the rules that cannot be saved: a chord the matcher
// cannot parse, or two rules claiming the same chord for the same application.
// Overriding a built-in shortcut is deliberately not an issue — user intent is
// meant to win over the registry.
func shortcutDraftIssues(drafts []config.CustomShortcut) []string {
	var issues []string
	seen := make(map[string]int, len(drafts))
	for i, draft := range drafts {
		chord := strings.TrimSpace(draft.Chord)
		if chord == "" {
			// An incomplete row is dropped by the save step, not an error.
			continue
		}
		if _, _, err := semantic.ParseChord(chord); err != nil {
			issues = append(issues, fmt.Sprintf("#%d", i+1))
			continue
		}
		scope := strings.ToLower(strings.TrimSpace(draft.App))
		normalized := strings.ToLower(strings.ReplaceAll(chord, " ", ""))
		key := scope + "\x00" + normalized
		if first, ok := seen[key]; ok {
			issues = append(issues, fmt.Sprintf("#%d/#%d", first+1, i+1))
			continue
		}
		seen[key] = i
	}
	return issues
}

// currentDisplayLabel resolves the persisted monitor selection to the label
// shown on the alignment mockup. Previously this read "Primary" even when the
// user had selected a secondary screen.
func currentDisplayLabel(selection string) string {
	displays := mygo.Screen.Displays()
	fallback := i18n.T("ui_display_primary")
	if len(displays) == 0 {
		return fallback
	}
	target := displays[0]
	for _, d := range displays {
		if d.Primary {
			target = d
			break
		}
	}
	sel := strings.TrimSpace(selection)
	if sel != "" && !strings.EqualFold(sel, "primary") {
		for _, d := range displays {
			if strconv.FormatInt(d.ID, 10) == sel || strings.EqualFold(strings.TrimSpace(d.Label), sel) {
				target = d
				break
			}
		}
	}
	if label := strings.TrimSpace(target.Label); label != "" {
		return label
	}
	return fallback
}

// splitList accepts the same separators a user is likely to paste: commas,
// semicolons, and newlines.
func splitList(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r'
	})
	out := make([]string, 0, len(fields))
	seen := make(map[string]bool)
	for _, field := range fields {
		value := strings.TrimSpace(field)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, value)
	}
	return out
}

// openPath hands a file or folder to its default Windows handler without
// flashing a console window.
func openPath(path string) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	return cmd.Start()
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
										ui.Text(c, currentDisplayLabel(cfg.Appearance.Monitor)).FontSize(8.5).Bold().TextColor(p.text)
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

// clampAnimDuration keeps the animation duration in a useful, non-zero band.
func clampAnimDuration(v float64) float64 {
	if v < 0.10 {
		return 0.10
	}
	if v > 1.50 {
		return 1.50
	}
	return v
}

// themeShortName drops the parenthetical original name so a gallery card can
// hold the label on one line. The full bilingual name still appears in the
// showcase stage and in the configuration file.
func themeShortName(name string) string {
	full := themeDisplayName(name)
	for _, sep := range []string{" (", "（"} {
		if i := strings.Index(full, sep); i > 0 {
			return strings.TrimSpace(full[:i])
		}
	}
	return full
}

// renderThemeSelect is the theme choice as a menu. A wall of twenty equal
// cards makes choosing feel like inventory work; a menu answers the question
// only when it is asked, and its trigger keeps the current cap in sight while
// it is closed.
func (sc *SettingsController) renderThemeSelect(c *ui.Context, cfg config.AppConfig, p palette, current theme.ThemeConfig) {
	choice := current.Name
	// Alternatives show the theme as designed; the trigger shows the live
	// design, so it can never contradict the stage above it.
	preview := themePreviewConfig(cfg)

	sel := ui.SelectBase(c, &choice)
	sel.Trigger.
		Width(232).
		Height(42).
		PaddingX(9).
		Gap(9).
		Radius(8).
		Background(p.sidebar).
		Border(1, p.hairline).
		AlignItems(ui.Center).
		Children(func() {
			// The swatch sits on the theme's own backdrop, so a white cap is
			// not swallowed by the row it happens to be drawn on.
			stage := ui.Box(c).
				Size(40, 34).
				Radius(6).
				Background(previewStageColor(current)).
				Border(1, p.hairline).
				AlignItems(ui.Center).
				Justify(ui.Center)
			stage.Children(func() {
				sc.renderKeycapPreview(c, cfg, current, cfg.Appearance.Style, "theme-select-cap", 36, 30, 0.58, false)
			})
			ui.Text(c, themeShortName(current.Name)).
				FontSize(12.5).
				FontWeight(600).
				TextColor(p.text).
				SingleLine().
				Grow(1)
			ui.Icon(c, iconChevronDown).Size(12, 12).TextColor(p.muted)
		})

	sel.Popup(func(panel ui.Element) {
		panel.
			Width(276).
			Background(p.card).
			Border(1, p.hairline).
			Radius(10).
			Padding(4).
			Shadow(0, 12, 32, 0, ui.RGBA(0, 0, 0, 0.50))

		ui.Scroll(c).MaxHeight(312).Children(func() {
			for _, preset := range theme.AllThemes {
				itemPreset := preset
				selected := strings.EqualFold(cfg.ThemePreset, itemPreset.Name)
				item := sel.Item(itemPreset.Name).
					Height(44).
					PaddingX(8).
					Gap(10).
					Radius(7).
					AlignItems(ui.Center)

				if item.Highlighted() {
					item.Background(p.press)
				} else if selected {
					item.Background(p.hover)
				}

				item.Children(func() {
					stage := ui.Box(c).
						Size(42, 34).
						Radius(6).
						Background(previewStageColor(itemPreset)).
						Border(1, p.hairline).
						AlignItems(ui.Center).
						Justify(ui.Center)
					stage.Children(func() {
						sc.renderKeycapPreview(c, preview, itemPreset, "", "theme-item-"+itemPreset.Name, 38, 30, 0.52, item.Hovered())
					})

					ui.Column(c).Gap(2).Grow(1).Shrink(1).MinWidth(0).Children(func() {
						ui.Text(c, themeShortName(itemPreset.Name)).FontSize(12.5).TextColor(p.text).SingleLine()
						ui.Text(c, styleItemDetail(itemPreset)).FontSize(10.5).TextColor(p.muted).SingleLine()
					})

					if selected {
						ui.Text(c, "✓").FontSize(12.5).Bold().TextColor(p.accent)
					}
				})
			}
		})
	})

	if sel.Trigger.Changed() {
		sc.selectTheme(theme.ThemeByName(choice))
	}
}

// styleItemDetail is the second line of a theme row: the build the theme would
// choose for itself, and how far that build travels. Two themes can share a
// palette family; they rarely share both.
func styleItemDetail(preset theme.ThemeConfig) string {
	style := preset.DefaultStyle
	if style == "" {
		style = "lowprofile"
	}
	return physicalStyleShortName(style) + " · " + styleTravelLabel(preset, "")
}

// selectTheme applies a preset. Theme-derived values return to inheritance
// instead of being frozen into the config as copies of today's colors.
func (sc *SettingsController) selectTheme(preset theme.ThemeConfig) {
	sc.updateConfig(func(ac *config.AppConfig) {
		ac.ThemePreset = preset.Name
		ac.Appearance.Style = ""
		ac.Color.Color = ""
		ac.Color.SecondaryColor = ""
		ac.Text.Color = ""
		ac.Border.Color = ""
		ac.Border.Radius = 0
	})
}

// renderPaletteGallery is one shelf of paint, not a wall of cards. Color is
// read faster than any label, so fifteen presets collapse into a single row of
// keycap-shaped swatches; the custom entry is a plus that opens three wells.
func (sc *SettingsController) renderPaletteGallery(c *ui.Context, cfg config.AppConfig, p palette, t theme.ThemeConfig) {
	customActive := paletteIsCustom(cfg)

	ui.Column(c).Padding(14, 16).Gap(14).Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			sc.renderPaletteChip(c, p, t, ColorSchemePreset{Name: i18n.T("ui_row_follow_theme")}, true, cfg.Color.Color == "")
			for _, sch := range MasterColorSchemes {
				sc.renderPaletteChip(c, p, t, sch, false, strings.EqualFold(cfg.Color.Color, sch.Primary))
			}
			ui.Box(c).Width(1).Height(16).Background(p.hairline)
			sc.renderCustomPaletteChip(c, p, customActive)
		})

		if customActive {
			ui.Row(c).FillWidth().Gap(16).Children(func() {
				sc.renderPaletteWell(c, p, i18n.T("ui_row_palette_bg"), "palette-bg", cfg, t.KeyCapBackground, func(ac *config.AppConfig) *string { return &ac.Color.Color })
				sc.renderPaletteWell(c, p, i18n.T("ui_row_palette_border"), "palette-border", cfg, t.KeyCapBorder, func(ac *config.AppConfig) *string { return &ac.Color.SecondaryColor })
				sc.renderPaletteWell(c, p, i18n.T("ui_row_palette_text"), "palette-text", cfg, t.KeyCapTextColor, func(ac *config.AppConfig) *string { return &ac.Text.Color })
			})
		}
	})
}

// renderPaletteChip is a keycap-shaped swatch: the fill, outline and label
// colors are the whole message, so there is nothing else on it.
func (sc *SettingsController) renderPaletteChip(c *ui.Context, p palette, t theme.ThemeConfig, sch ColorSchemePreset, inherit, selected bool) {
	bg := t.KeyCapBackground
	border := t.KeyCapBorder
	text := t.KeyCapTextColor
	if !inherit {
		bg = parseHexSafe(sch.Primary, bg)
		border = parseHexSafe(sch.Secondary, border)
		text = parseHexSafe(sch.Text, text)
	}

	key := "palette-chip-" + sch.Name
	if inherit {
		key = "palette-chip-inherit"
	}

	btn := ui.ButtonBase(c.Key(key)).
		Size(30, 28).
		Radius(7).
		AlignItems(ui.Center).
		Justify(ui.Center).
		FocusRing(false).
		Background(bg)

	switch {
	case selected:
		btn.Border(2, p.accent).Shadow(0, 0, 8, 0, p.accent.Alpha(0.35))
	case btn.Hovered():
		btn.Border(1, p.accent.Alpha(0.65))
	default:
		btn.Border(1, border)
	}

	btn.Children(func() {
		ui.Text(c, "A").FontSize(9.5).FontWeight(700).TextColor(text)
	})

	if btn.Clicked() {
		sc.selectPalette(sch, inherit)
	}
}

// renderCustomPaletteChip opens the three wells. It is a door, not a swatch:
// the colors it produces are shown below it the moment it is used.
func (sc *SettingsController) renderCustomPaletteChip(c *ui.Context, p palette, selected bool) {
	btn := ui.ButtonBase(c.Key("palette-chip-custom")).
		Size(30, 28).
		Radius(7).
		AlignItems(ui.Center).
		Justify(ui.Center).
		FocusRing(false).
		Background(p.sidebar)

	if selected {
		btn.Border(2, p.accent).Shadow(0, 0, 8, 0, p.accent.Alpha(0.35))
	} else if btn.Hovered() {
		btn.Border(1, p.accent.Alpha(0.65))
	} else {
		btn.Border(1, p.muted.Alpha(0.45))
	}

	btn.Children(func() {
		ui.Icon(c, iconPlus).Size(11, 11).TextColor(p.text.Alpha(0.55))
	})

	if btn.Clicked() && !selected {
		sc.materializePalette()
	}
}

// selectPalette applies a preset, or returns to the theme by clearing every
// override. One function for both keeps the two paths from drifting apart.
func (sc *SettingsController) selectPalette(sch ColorSchemePreset, inherit bool) {
	sc.updateConfig(func(ac *config.AppConfig) {
		if inherit {
			ac.Color.Color = ""
			ac.Color.SecondaryColor = ""
			ac.Text.Color = ""
			ac.Border.Color = ""
			return
		}
		ac.Color.Color = sch.Primary
		ac.Color.SecondaryColor = sch.Secondary
		ac.Text.Color = sch.Text
		ac.Border.Color = sch.Secondary
	})
}

// materializePalette copies the active theme colors into explicit values, so
// the wells start from what the user already sees instead of from nothing.
func (sc *SettingsController) materializePalette() {
	t := theme.ThemeByName(sc.store.Get().ThemePreset)
	sc.updateConfig(func(ac *config.AppConfig) {
		ac.Color.Color = colorHex(t.KeyCapBackground)
		ac.Color.SecondaryColor = colorHex(t.KeyCapBorder)
		ac.Text.Color = colorHex(t.KeyCapTextColor)
		ac.Border.Color = ""
	})
}

// renderPaletteWell is one labelled color well under the custom palette card.
func (sc *SettingsController) renderPaletteWell(c *ui.Context, p palette, label, key string, cfg config.AppConfig, fallback ui.Color, pick func(*config.AppConfig) *string) {
	ui.Column(c).Grow(1).Gap(6).AlignItems(ui.Start).Children(func() {
		ui.Text(c, label).FontSize(11).TextColor(p.muted).SingleLine()
		col := parseHexSafe(*pick(&cfg), fallback)
		well := ui.ColorWell(c.Key(key), &col)
		if well.Changed() {
			hex := colorHex(col)
			sc.updateConfig(func(ac *config.AppConfig) { *pick(ac) = hex })
		}
	})
}

// paletteIsCustom reports whether the palette is neither inherited nor one of
// the fifteen presets, i.e. the custom wells should be visible.
func paletteIsCustom(cfg config.AppConfig) bool {
	if cfg.Color.Color == "" {
		return false
	}
	for _, sch := range MasterColorSchemes {
		if strings.EqualFold(cfg.Color.Color, sch.Primary) {
			return false
		}
	}
	return true
}

// themePreviewConfig strips every user override that is not part of the theme
// itself, so a gallery card shows the design rather than the current tweaks.
func themePreviewConfig(cfg config.AppConfig) config.AppConfig {
	preview := cfg
	preview.Color.Color = ""
	preview.Color.SecondaryColor = ""
	preview.Text.Color = ""
	preview.Border.Color = ""
	preview.Border.Radius = 0
	preview.Border.Width = 0
	preview.Modifier.Highlight = false
	return preview
}

// modelChoices orders the physical builds as a spectrum: Auto first, then the
// shallowest mechanical build to the deepest, then the optical ones.
var modelChoices = []string{"auto", "laptop", "minimal", "lowprofile", "pbt", "retro", "glass", "neon", "outline"}

// modelStyleID turns a menu entry back into a config value: Auto is the empty
// string, which means "follow the theme".
func modelStyleID(choice string) string {
	if choice == "auto" {
		return ""
	}
	return choice
}

// renderModelSelect is the build menu. Auto is a mode rather than a ninth
// build, so it stays first and names what it resolves to; every other entry
// renders the same key with the real HUD renderer, which is the only honest
// comparison of one physical model against another.
func (sc *SettingsController) renderModelSelect(c *ui.Context, cfg config.AppConfig, p palette, t theme.ThemeConfig) {
	choice := cfg.Appearance.Style
	if choice == "" {
		choice = "auto"
	}
	style := cfg.Appearance.Style

	sel := ui.SelectBase(c, &choice)
	sel.Trigger.
		Width(232).
		Height(42).
		PaddingX(9).
		Gap(9).
		Radius(8).
		Background(p.sidebar).
		Border(1, p.hairline).
		AlignItems(ui.Center).
		Children(func() {
			stage := ui.Box(c).
				Size(40, 34).
				Radius(6).
				Background(previewStageColor(t)).
				Border(1, p.hairline).
				AlignItems(ui.Center).
				Justify(ui.Center)
			stage.Children(func() {
				sc.renderKeycapPreview(c, cfg, t, style, "model-select-cap", 36, 30, 0.58, false)
			})
			ui.Column(c).Gap(2).Grow(1).Shrink(1).MinWidth(0).Children(func() {
				ui.Text(c, physicalStyleShortName(style)).FontSize(12.5).FontWeight(600).TextColor(p.text).SingleLine()
				ui.Text(c, styleTravelLabel(t, style)).FontSize(10.5).TextColor(p.muted).SingleLine()
			})
			ui.Icon(c, iconChevronDown).Size(12, 12).TextColor(p.muted)
		})

	sel.Popup(func(panel ui.Element) {
		panel.
			Width(276).
			Background(p.card).
			Border(1, p.hairline).
			Radius(10).
			Padding(4).
			Shadow(0, 12, 32, 0, ui.RGBA(0, 0, 0, 0.50))

		// Short enough to open below the trigger on a 720-high window instead
		// of flipping over the stage.
		ui.Scroll(c).MaxHeight(264).Children(func() {
			for _, id := range modelChoices {
				itemStyle := modelStyleID(id)
				selected := style == itemStyle
				item := sel.Item(id).
					Height(44).
					PaddingX(8).
					Gap(10).
					Radius(7).
					AlignItems(ui.Center)

				if item.Highlighted() {
					item.Background(p.press)
				} else if selected {
					item.Background(p.hover)
				}

				item.Children(func() {
					stage := ui.Box(c).
						Size(42, 34).
						Radius(6).
						Background(previewStageColor(t)).
						Border(1, p.hairline).
						AlignItems(ui.Center).
						Justify(ui.Center)
					stage.Children(func() {
						sc.renderKeycapPreview(c, cfg, t, itemStyle, "model-item-"+id, 38, 30, 0.52, item.Hovered())
					})

					ui.Column(c).Gap(2).Grow(1).Shrink(1).MinWidth(0).Children(func() {
						ui.Text(c, physicalStyleShortName(itemStyle)).FontSize(12.5).TextColor(p.text).SingleLine()
						ui.Text(c, modelItemDetail(t, cfg, itemStyle)).FontSize(10.5).TextColor(p.muted).SingleLine()
					})

					if selected {
						ui.Text(c, "✓").FontSize(12.5).Bold().TextColor(p.accent)
					}
				})
			}
		})
	})

	if sel.Trigger.Changed() {
		next := modelStyleID(choice)
		sc.updateConfig(func(ac *config.AppConfig) { ac.Appearance.Style = next })
	}
}

// modelItemDetail names what Auto resolves to, and what every other build
// measures, so the list stays honest about inheritance.
func modelItemDetail(t theme.ThemeConfig, cfg config.AppConfig, styleID string) string {
	if styleID == "" {
		return physicalStyleShortName(effectiveStyle(t, cfg))
	}
	return styleTravelLabel(t, styleID)
}

// renderKeycapPreview renders the very same "A" cap the overlay draws, scaled
// into a gallery stage. The cap is not a hand-drawn illustration: it is the
// live renderer, so a preview can never drift away from the HUD.
func (sc *SettingsController) renderKeycapPreview(
	c *ui.Context,
	cfg config.AppConfig,
	t theme.ThemeConfig,
	styleID, keyID string,
	stageW, stageH, scale float32,
	hovered bool,
) {
	target := float32(0)
	if hovered {
		// The full advertised travel, scaled with the preview, not the overlay's
		// damped 70%: this demonstrates the physical distance.
		target = styleTravelDepth(t, styleID) * scale
	}

	holder := ui.Box(c).
		Width(stageW).
		Height(stageH).
		AlignItems(ui.Center).
		Justify(ui.Center)
	y := holder.AnimateTo("press", target, ui.Animation{Duration: 140 * time.Millisecond, Ease: ui.EaseOut})
	holder.Translate(0, y).Children(func() {
		ui.Row(c).Scale(scale, scale).Children(func() {
			renderKeyCap(c, keyID, engine.KeyCap{Label: "A"}, 1, 0, false, t, keycapPreviewConfig(cfg, styleID))
		})
	})
}

// keycapPreviewConfig isolates the physical model from unrelated user toggles.
// The gallery is about build, not borders or animation preferences, so it shows
// the theme's palette on a cap whose structure is defined by the model alone.
func keycapPreviewConfig(cfg config.AppConfig, styleID string) config.AppConfig {
	preview := cfg
	preview.Appearance.Style = styleID
	preview.Appearance.DisablePressAnimation = false
	preview.Text.Size = 19
	preview.Text.Variant = "text"
	preview.Layout.ShowIcon = false
	preview.Border.Enabled = true
	preview.Border.Width = 0
	return preview
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
							ui.Text(c, i18n.T("ui_reset_confirm_title")).FontSize(16).Bold().TextColor(p.text)
							ui.Text(c, i18n.T("ui_reset_confirm_sub")).FontSize(11.5).TextColor(p.muted)
						})
					})

					ui.Text(c, i18n.T("ui_reset_confirm_body")).
						FontSize(13).
						TextColor(p.muted).
						LineHeight(1.5)

					ui.Row(c).Justify(ui.End).Gap(10).Margin(8, 0, 0, 0).Children(func() {
						if ui.Button(c.Key("cancel-reset-btn"), i18n.T("ui_btn_cancel")).
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

						if ui.Button(c.Key("confirm-reset-btn"), i18n.T("ui_btn_confirm_reset")).
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
							ui.Text(c, i18n.T("ui_filter_modal_title")).FontSize(16).Bold().TextColor(p.text)
							ui.Text(c, i18n.T("ui_filter_modal_desc")).FontSize(11.5).TextColor(p.muted)
						})
						if ui.Button(c.Key("close-filter-modal"), i18n.T("ui_btn_done")).Padding(6, 14).Radius(6).Background(p.accent).TextColor(p.onAccent).Bold().Clicked() {
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
						tabLabels := []string{i18n.T("ui_filter_tab_keyboard"), i18n.T("ui_filter_tab_mouse"), i18n.T("ui_filter_tab_numpad")}
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

							if ui.Button(c.Key("filter-all-btn"), i18n.T("Select All")).Padding(4, 9).Radius(6).Background(p.sidebar).TextColor(p.text).FontSize(11).Clicked() {
								allCommon := append(append([]string{}, allAlpha...), allMods...)
								allCommon = append(allCommon, "Enter", "Space", "Backspace", "Tab", "Esc", "Insert", "Delete", "Home", "End", "PageUp", "PageDown", "Up", "Down", "Left", "Right", "Left Click", "Right Click", "Middle Click", "Drag", "Scroll Up", "Scroll Down")
								sc.updateConfig(func(ac *config.AppConfig) {
									ac.General.AllowedKeys = allCommon
								})
							}
							if ui.Button(c.Key("filter-clear-btn"), i18n.T("ui_filter_clear")).Padding(4, 9).Radius(6).Background(p.sidebar).TextColor(p.text).FontSize(11).Clicked() {
								sc.updateConfig(func(ac *config.AppConfig) {
									ac.General.AllowedKeys = []string{}
								})
							}
							if ui.Button(c.Key("filter-mods-btn"), i18n.T("ui_filter_mods_only")).Padding(4, 9).Radius(6).Background(p.sidebar).TextColor(p.text).FontSize(11).Clicked() {
								sc.updateConfig(func(ac *config.AppConfig) {
									ac.General.AllowedKeys = allMods
								})
							}
							if ui.Button(c.Key("filter-alpha-btn"), i18n.T("ui_filter_alpha_num")).Padding(4, 9).Radius(6).Background(p.sidebar).TextColor(p.text).FontSize(11).Clicked() {
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
									sc.renderKeycapButton(c, "Left Click", i18n.T("ui_mouse_label_left"), 64, 48, 0, cfg, p)

									// Center Column (Middle, Scroll, Drag)
									ui.Column(c).Gap(6).AlignItems(ui.Center).Children(func() {
										sc.renderKeycapButton(c, "Middle Click", i18n.T("ui_mouse_label_middle"), 80, 32, 0, cfg, p)
										sc.renderKeycapButton(c, "Scroll Up", i18n.T("ui_mouse_label_scroll_up"), 80, 30, 0, cfg, p)
										sc.renderKeycapButton(c, "Scroll Down", i18n.T("ui_mouse_label_scroll_down"), 80, 30, 0, cfg, p)
										sc.renderKeycapButton(c, "Drag", i18n.T("ui_mouse_label_drag"), 80, 32, 0, cfg, p)
									})

									// Right Click Button
									sc.renderKeycapButton(c, "Right Click", i18n.T("ui_mouse_label_right"), 64, 48, 0, cfg, p)
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
