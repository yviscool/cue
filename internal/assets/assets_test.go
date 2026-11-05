package assets

import (
	"testing"
)

func TestEmbeddedAssets(t *testing.T) {
	if len(AppIconICO) == 0 {
		t.Fatal("AppIconICO is empty")
	}
	if len(AppIconPNG) == 0 {
		t.Fatal("AppIconPNG is empty")
	}
	if len(LogoPNG) == 0 {
		t.Fatal("LogoPNG is empty")
	}
	if len(LogoSVG) == 0 {
		t.Fatal("LogoSVG is empty")
	}
	if len(LogoProSVG) == 0 {
		t.Fatal("LogoProSVG is empty")
	}

	src, err := GetAppLogoSource()
	if err != nil {
		t.Fatalf("failed to decode LogoSVG: %v", err)
	}
	if src == nil {
		t.Fatal("decoded LogoSVG is nil")
	}

	iconSrc, err := GetAppIconSource()
	if err != nil {
		t.Fatalf("failed to decode AppIconPNG: %v", err)
	}
	if iconSrc == nil {
		t.Fatal("decoded AppIconPNG is nil")
	}
}
