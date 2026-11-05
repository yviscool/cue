package assets

import (
	_ "embed"
	"github.com/egoist/mygo/ui"
)

//go:embed app.ico
var AppIconICO []byte

//go:embed app.png
var AppIconPNG []byte

//go:embed tray.png
var TrayIconPNG []byte

//go:embed logo.svg
var LogoSVG []byte

//go:embed logo-pro.svg
var LogoProSVG []byte

//go:embed logo.png
var LogoPNG []byte

// GetAppLogoSource returns the decoded UI image source for Cue logo.
func GetAppLogoSource() (ui.ImageSource, error) {
	return ui.DecodeImage(LogoPNG)
}

// GetAppIconSource returns the decoded UI image source for Cue bitmap icon.
func GetAppIconSource() (ui.ImageSource, error) {
	return ui.DecodeImage(AppIconPNG)
}
