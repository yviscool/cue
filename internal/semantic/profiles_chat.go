package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// Chat, meeting and streaming-audience tools. The mute/video/hand-raise chords
// matter most here: they are what a live audience actually watches you press.

func (r *Registry) initChat() {
	r.initSlack()
	r.initTeams()
	r.initDiscord()
	r.initWeChat()
	r.initZoom()
}

func (r *Registry) initSlack() {
	slack := newProfile("slack", "Slack", context.CategoryChat,
		"slack.exe",
	).source("Slack keyboard shortcuts (Windows)", "Slack 4.x desktop (Windows)")

	// Navigation
	slack.k(hook.ModCtrl, "K", "Quick Switcher").
		k(hook.ModCtrl, "N", "New Message").
		k(hook.ModCtrl|hook.ModShift, "K", "Browse DMs").
		k(hook.ModCtrl|hook.ModShift, "L", "Browse Channels").
		k(hook.ModAlt, "UP", "Previous Channel").
		k(hook.ModAlt, "DOWN", "Next Channel").
		k(hook.ModCtrl|hook.ModShift, "M", "Activity").
		k(hook.ModCtrl|hook.ModShift, "A", "All Unreads").
		k(hook.ModCtrl|hook.ModShift, "T", "Threads").
		k(hook.ModCtrl|hook.ModShift, "E", "People").
		k(hook.ModShift, "ESC", "Mark All as Read")

	// Search
	slack.k(hook.ModCtrl, "F", "Search").
		k(hook.ModCtrl|hook.ModShift, "F", "Search All Channels")

	// Messages
	slack.k(hook.ModNone, "ENTER", "Send Message").
		k(hook.ModShift, "ENTER", "New Line").
		k(hook.ModCtrl, "UP", "Edit Last Message").
		k(hook.ModCtrl, "B", "Bold").
		k(hook.ModCtrl, "I", "Italic").
		k(hook.ModCtrl|hook.ModShift, "C", "Code Block").
		k(hook.ModCtrl|hook.ModShift, ">", "Block Quote").
		k(hook.ModCtrl, "Z", "Undo").
		k(hook.ModCtrl|hook.ModShift, "Z", "Redo").
		k(hook.ModCtrl|hook.ModShift, "Y", "Set Status").
		k(hook.ModCtrl, ",", "Preferences")

	r.addProfile(slack)
}

func (r *Registry) initTeams() {
	teams := newProfile("teams", "Microsoft Teams", context.CategoryChat,
		"ms-teams.exe", "teams.exe",
	).source("Microsoft Teams keyboard shortcuts", "Microsoft Teams (work or school, Windows)")

	// Search & navigation
	teams.k(hook.ModCtrl, "K", "Search/Command").
		k(hook.ModCtrl, "N", "New Chat").
		k(hook.ModCtrl, "1", "Go to Activity").
		k(hook.ModCtrl, "2", "Go to Chat").
		k(hook.ModCtrl, "3", "Go to Teams").
		k(hook.ModCtrl, "4", "Go to Calendar").
		k(hook.ModCtrl, "5", "Go to Calls").
		k(hook.ModCtrl, "6", "Go to Files")

	// Meeting controls
	teams.k(hook.ModCtrl|hook.ModShift, "M", "Toggle Mute").
		k(hook.ModCtrl|hook.ModShift, "O", "Toggle Video").
		k(hook.ModCtrl|hook.ModShift, "K", "Raise Hand").
		k(hook.ModCtrl|hook.ModShift, "E", "Start Screen Share").
		k(hook.ModCtrl|hook.ModShift, "A", "Accept Video Call").
		k(hook.ModCtrl|hook.ModShift, "D", "Decline Video Call").
		k(hook.ModCtrl|hook.ModShift, "H", "Hang Up")

	// Composing
	teams.k(hook.ModNone, "ENTER", "Send Message").
		k(hook.ModShift, "ENTER", "New Line").
		k(hook.ModCtrl, "B", "Bold").
		k(hook.ModCtrl, "I", "Italic").
		k(hook.ModCtrl, "U", "Underline").
		k(hook.ModCtrl, "Z", "Undo").
		k(hook.ModCtrl, "Y", "Redo")

	r.addProfile(teams)
}

func (r *Registry) initDiscord() {
	discord := newProfile("discord", "Discord", context.CategoryChat,
		"discord.exe",
	).source("Discord keyboard shortcuts", "Discord desktop (Windows)")

	discord.k(hook.ModCtrl, "K", "Quick Switcher").
		k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "/", "Keyboard Shortcuts").
		k(hook.ModCtrl|hook.ModShift, "M", "Toggle Mute").
		k(hook.ModCtrl|hook.ModShift, "D", "Toggle Deafen").
		k(hook.ModAlt, "UP", "Previous Channel").
		k(hook.ModAlt, "DOWN", "Next Channel").
		k(hook.ModAlt|hook.ModShift, "UP", "Previous Server").
		k(hook.ModAlt|hook.ModShift, "DOWN", "Next Server").
		k(hook.ModShift, "ESC", "Mark Server Read").
		k(hook.ModNone, "ESC", "Mark Channel Read").
		k(hook.ModCtrl, "B", "Bold").
		k(hook.ModCtrl, "I", "Italic").
		k(hook.ModCtrl, "U", "Underline").
		k(hook.ModNone, "ENTER", "Send Message").
		k(hook.ModShift, "ENTER", "New Line")

	r.addProfile(discord)
}

func (r *Registry) initWeChat() {
	wechat := newProfile("wechat", "WeChat", context.CategoryChat,
		"wechat.exe", "weixin.exe",
	).source("WeChat for Windows shortcuts", "WeChat 3.9 (Windows)")

	wechat.k(hook.ModCtrl, "F", "Search").
		k(hook.ModAlt, "A", "Screenshot").
		k(hook.ModCtrl|hook.ModAlt, "W", "Open WeChat").
		k(hook.ModNone, "ENTER", "Send Message").
		k(hook.ModCtrl, "ENTER", "New Line").
		k(hook.ModNone, "ESC", "Close Window")

	r.addProfile(wechat)
}

func (r *Registry) initZoom() {
	zoom := newProfile("zoom", "Zoom", context.CategoryChat,
		"zoom.exe",
	).source("Zoom keyboard shortcuts (Windows)", "Zoom Workplace 6.x (Windows)")

	// The four chords a live audience cares about most.
	zoom.k(hook.ModAlt, "A", "Toggle Mute").
		k(hook.ModAlt, "V", "Toggle Video").
		k(hook.ModAlt, "H", "Raise Hand").
		k(hook.ModAlt, "S", "Start Screen Share")

	// Recording
	zoom.k(hook.ModAlt, "R", "Start Recording").
		k(hook.ModAlt, "C", "Start Cloud Recording").
		k(hook.ModAlt, "P", "Pause Recording")

	// Meeting management
	zoom.k(hook.ModAlt, "M", "Mute All").
		k(hook.ModAlt, "Q", "End Meeting").
		k(hook.ModAlt, "U", "Participants Panel").
		k(hook.ModAlt, "N", "Switch Camera").
		k(hook.ModAlt, "F", "Enter Full Screen").
		k(hook.ModNone, "ESC", "Exit Full Screen")

	r.addProfile(zoom)
}
