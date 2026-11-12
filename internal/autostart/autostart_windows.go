//go:build windows

// Package autostart registers Cue in the current user's login items.
package autostart

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const (
	runKey = `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	value  = "Cue"
)

// Set enables or disables launch-at-login for the current executable. It uses
// the per-user Run key, so it never needs administrator rights and is easy for
// users to audit with regedit.
func Set(enabled bool) error {
	if enabled {
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		return run("add", runKey, "/v", value, "/t", "REG_SZ", "/d", `"`+exe+`"`, "/f")
	}
	// Deleting a value that does not exist is not an error for the user's intent.
	err := run("delete", runKey, "/v", value, "/f")
	if err != nil {
		return nil
	}
	return nil
}

func run(args ...string) error {
	cmd := exec.Command("reg", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("reg %s: %w: %s", args[0], err, string(out))
	}
	return nil
}
