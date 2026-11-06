package main

import (
	"strconv"
	"strings"

	"github.com/egoist/mygo"
)

// resolveDisplay turns the persisted monitor preference into a concrete
// display. A saved ID wins; if IDs changed after a reboot the label is used as
// a fallback. Missing displays always fall back to the primary screen so the
// HUD can never disappear off-screen after a monitor is unplugged.
func resolveDisplay(displays []mygo.Display, selection string) (mygo.Display, bool) {
	if len(displays) == 0 {
		return mygo.Display{}, false
	}

	selection = strings.TrimSpace(selection)
	if selection != "" && !strings.EqualFold(selection, "primary") {
		for _, d := range displays {
			if strconv.FormatInt(d.ID, 10) == selection {
				return d, true
			}
		}
		for _, d := range displays {
			if strings.EqualFold(strings.TrimSpace(d.Label), selection) {
				return d, true
			}
		}
	}

	for _, d := range displays {
		if d.Primary {
			return d, true
		}
	}
	return displays[0], true
}
