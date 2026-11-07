package hud

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/yviscool/cue/internal/assets"
	"github.com/yviscool/cue/internal/engine"
	"github.com/yviscool/cue/internal/i18n"
	"github.com/yviscool/cue/internal/theme"
)

// FlyoutController coordinates settings changes from the flyout UI.
type FlyoutController struct {
	ViewState *ViewState
	Engine    *engine.Engine
	logoSrc   ui.ImageSource
}

// RenderFlyout renders the tray control center flyout.
func (fc *FlyoutController) RenderFlyout(c *ui.Context) {
	t := c.Theme()
	vState := fc.ViewState
	activeTheme := vState.GetTheme()
	snap := vState.Snapshot()

	if fc.logoSrc == nil {
		fc.logoSrc, _ = assets.GetAppLogoSource()
		if fc.logoSrc == nil {
			fc.logoSrc, _ = assets.GetAppIconSource()
		}
	}

	ui.Column(c).Padding(16).Gap(14).Children(func() {
		// Header with Cue logo
		ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
			if fc.logoSrc != nil {
				ui.Image(c, fc.logoSrc).Size(24, 24)
			}
			ui.Text(c, i18n.T("app_name")).FontSize(17).Bold()
		})
		ui.Text(c, i18n.T("app_subtitle")).FontSize(11).TextColor(t.TextMuted)

		ui.Divider(c)

		// 1. Industrial Design Theme Section
		ui.Text(c, i18n.T("flyout_design_themes")).FontSize(12).Bold().TextColor(t.Text)
		ui.Column(c).Gap(6).Children(func() {
			for _, thm := range theme.AllThemes {
				themeName := thm.Name
				isSelected := activeTheme.Name == themeName
				bg := ui.Transparent
				if isSelected {
					bg = t.Accent.Alpha(0.15)
				}

				row := ui.Row(c).
					Key("theme-"+themeName).
					Padding(8, 12).
					Radius(8).
					Background(bg).
					Border(1, t.Border).
					AlignItems(ui.Center).
					Gap(10)

				row.Children(func() {
					radioText := "○"
					if isSelected {
						radioText = "●"
					}
					ui.Text(c, radioText).FontSize(14).TextColor(t.Accent).Bold()
					ui.Column(c).Gap(2).Children(func() {
						ui.Text(c, themeDisplayName(thm.Name)).FontSize(13).Bold()
						ui.Text(c, themeDescription(thm.Name)).FontSize(10).TextColor(t.TextMuted)
					})
				})

				if row.Clicked() {
					vState.SetTheme(thm)
					if fc.Engine != nil {
						fc.Engine.SetThemeName(themeName)
					}
				}
			}
		})

		ui.Divider(c)

		// 2. Filter Mode Section
		ui.Text(c, i18n.T("flyout_capture_mode")).FontSize(12).Bold().TextColor(t.Text)
		ui.Row(c).Gap(6).Children(func() {
			modes := []struct {
				mode  engine.FilterMode
				label string
			}{
				{engine.FilterSmart, i18n.T("flyout_mode_smart")},
				{engine.FilterAll, i18n.T("flyout_mode_all")},
				{engine.FilterSilent, i18n.T("flyout_mode_silent")},
			}

			for _, m := range modes {
				modeVal := m.mode
				isSelected := snap.FilterMode == modeVal
				btn := ui.Button(c.Key("mode-"+m.label), m.label).Grow(1)
				if isSelected {
					btn.Background(t.Accent).TextColor(t.AccentText)
				}
				if btn.Clicked() {
					if fc.Engine != nil {
						fc.Engine.SetFilterMode(modeVal)
					}
				}
			}
		})

		ui.Divider(c)

		// 3. Audio & Interaction Section
		ui.Row(c).AlignItems(ui.Center).Gap(10).Children(func() {
			soundLabel := i18n.T("tactile_audio_off")
			if snap.SoundActive {
				soundLabel = i18n.T("tactile_audio_on")
			}
			if ui.Button(c.Key("sound-toggle"), soundLabel).Grow(1).Clicked() {
				if fc.Engine != nil {
					fc.Engine.SetSoundActive(!fc.Engine.IsSoundActive())
				}
			}
		})

		ui.Divider(c)

		// 4. Footer & Actions
		ui.Row(c).AlignItems(ui.Center).Children(func() {
			ui.Text(c, i18n.T("flyout_alt_drag")).FontSize(11).TextColor(t.TextMuted).Grow(1)
			if ui.Button(c.Key("quit-app"), i18n.T("btn_quit")).Clicked() {
				mygo.App.Quit()
			}
		})
	})
}
