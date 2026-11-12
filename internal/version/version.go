// Package version exposes the build identity injected by the release pipeline.
package version

import "strings"

// Version is overridden at link time:
//
//	go build -ldflags "-X github.com/yviscool/cue/internal/version.Version=v1.4.0"
//
// It intentionally stays "dev" for local builds so the About page never claims
// to be a release that it is not.
var Version = "dev"

// Display returns the build version with a leading v, or "dev" when the build
// pipeline did not inject anything.
func Display() string {
	v := strings.TrimSpace(Version)
	if v == "" {
		return "dev"
	}
	if strings.HasPrefix(v, "v") || v == "dev" {
		return v
	}
	return "v" + v
}
