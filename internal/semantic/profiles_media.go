package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// Media and streaming tools.

func (r *Registry) initMedia() {
	r.initPremiere()
	r.initVLC()
	r.initOBS()
}

func (r *Registry) initPremiere() {
	premiere := newProfile("premiere", "Adobe Premiere Pro", context.CategoryMedia,
		"adobe premiere pro.exe",
	).source("Premiere Pro default keyboard shortcuts", "Premiere Pro 2024 (Windows)")

	// Playback
	premiere.k(hook.ModNone, "SPACE", "Play/Pause").
		k(hook.ModNone, "J", "Play Backward").
		k(hook.ModNone, "K", "Pause").
		k(hook.ModNone, "L", "Play Forward").
		k(hook.ModNone, "LEFT", "Step Back").
		k(hook.ModNone, "RIGHT", "Step Forward").
		k(hook.ModShift, "LEFT", "Step Back 5 Frames").
		k(hook.ModShift, "RIGHT", "Step Forward 5 Frames")

	// Editing
	premiere.k(hook.ModCtrl, "K", "Add Edit").
		k(hook.ModCtrl|hook.ModShift, "K", "Add Edit to All Tracks").
		k(hook.ModCtrl, "D", "Apply Video Transition").
		k(hook.ModCtrl|hook.ModShift, "D", "Apply Audio Transition").
		k(hook.ModCtrl, "R", "Speed/Duration").
		k(hook.ModCtrl|hook.ModShift, "E", "Export Frame")

	// Markers
	premiere.k(hook.ModNone, "I", "Mark In").
		k(hook.ModNone, "O", "Mark Out").
		k(hook.ModNone, "X", "Mark Clip").
		k(hook.ModNone, "M", "Add Marker").
		k(hook.ModShift, "M", "Next Marker").
		k(hook.ModCtrl|hook.ModShift, "M", "Previous Marker")

	// Project & export
	premiere.k(hook.ModCtrl, "M", "Export Media").
		k(hook.ModCtrl, "N", "New Sequence").
		k(hook.ModCtrl, "I", "Import Media").
		k(hook.ModCtrl, "S", "Save")

	// Timeline zoom
	premiere.k(hook.ModNone, "+", "Zoom In Timeline").
		k(hook.ModNone, "-", "Zoom Out Timeline").
		k(hook.ModNone, "\\", "Zoom to Sequence")

	r.addProfile(premiere)
}

func (r *Registry) initVLC() {
	vlc := newProfile("vlc", "VLC", context.CategoryMedia,
		"vlc.exe",
	).source("VLC hotkeys (Preferences > Hotkeys)", "VLC 3.0 (Windows)")

	// Playback
	vlc.k(hook.ModNone, "SPACE", "Play/Pause").
		k(hook.ModNone, "S", "Stop").
		k(hook.ModNone, "N", "Next Track").
		k(hook.ModNone, "P", "Previous Track").
		k(hook.ModCtrl, "RIGHT", "Seek Forward").
		k(hook.ModCtrl, "LEFT", "Seek Backward").
		k(hook.ModNone, "E", "Next Frame")

	// Volume & speed
	vlc.k(hook.ModCtrl, "UP", "Volume Up").
		k(hook.ModCtrl, "DOWN", "Volume Down").
		k(hook.ModNone, "M", "Mute").
		k(hook.ModNone, "+", "Speed Up").
		k(hook.ModNone, "-", "Slow Down").
		k(hook.ModNone, "=", "Normal Speed")

	// Tracks & subtitles
	vlc.k(hook.ModNone, "V", "Cycle Subtitles").
		k(hook.ModNone, "B", "Cycle Audio Track").
		k(hook.ModNone, "A", "Cycle Aspect Ratio").
		k(hook.ModNone, "C", "Cycle Crop").
		k(hook.ModNone, "T", "Show Time").
		k(hook.ModCtrl, "T", "Go to Time")
	vlc.k(hook.ModNone, "I", "Media Information")

	// Display
	vlc.k(hook.ModNone, "F", "Full Screen").
		k(hook.ModNone, "ESC", "Exit Full Screen")

	// Files
	vlc.k(hook.ModCtrl, "O", "Open File").
		k(hook.ModCtrl, "N", "Open Network").
		k(hook.ModCtrl, "L", "Toggle Playlist").
		k(hook.ModCtrl, "P", "Preferences")

	r.addProfile(vlc)
}

func (r *Registry) initOBS() {
	// OBS ships almost no default hotkeys: recording, streaming and scene
	// switching must be bound by the user. Only the editor-level defaults are
	// listed here, which is why this profile is intentionally small.
	obs := newProfile("obs", "OBS Studio", context.CategoryMedia,
		"obs64.exe", "obs32.exe",
	).source("OBS Studio default hotkeys (Settings > Hotkeys)", "OBS Studio 30.x (Windows)")

	obs.k(hook.ModCtrl, "S", "Save Scene Collection").
		k(hook.ModCtrl, "Z", "Undo").
		k(hook.ModCtrl, "Y", "Redo").
		k(hook.ModCtrl|hook.ModShift, "Z", "Redo").
		k(hook.ModCtrl, "D", "Duplicate").
		k(hook.ModCtrl, "C", "Copy").
		k(hook.ModCtrl, "V", "Paste").
		k(hook.ModNone, "DELETE", "Remove").
		k(hook.ModNone, "F2", "Rename")

	r.addProfile(obs)
}
