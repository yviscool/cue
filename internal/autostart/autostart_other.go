//go:build !windows

// Package autostart registers Cue in the current user's login items. The
// Windows build implements this with the per-user Run key; other platforms keep
// the setting visible but inert so the configuration schema stays portable.
package autostart

// Set is a no-op outside Windows.
func Set(enabled bool) error { return nil }
