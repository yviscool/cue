# Cue Configuration Guide

Cue keeps its settings in a human-readable JSON file:

```
%USERPROFILE%\.cue\config.json
```

The file is hot-reloaded: changes take effect as soon as they are valid, without
restarting the app. The settings window edits this same file atomically, so a
hand edit and the UI can never silently disagree.

Every field below can be changed from the settings window. Direct JSON editing
remains useful for bulk edits and for sharing configurations.

---

## Full example

```json
{
  "version": 6,
  "language": "en-US",
  "themePreset": "Swiss Typography 1957",
  "tactileAudio": true,
  "audioProfile": "chip",
  "audioVolume": 1.0,
  "semanticSub": true,
  "semanticContext": true,
  "showAppBadge": true,
  "customShortcuts": [
    { "app": "devcpp.exe", "chord": "Ctrl+Alt+B", "label": "Build My Project" }
  ],
  "general": {
    "filter": "none",
    "allowedKeys": [],
    "showEventHistory": true,
    "maxHistory": 5,
    "lingerDurationMs": 4500,
    "toggleShortcut": ["ShiftLeft", "F10"],
    "wordGrouping": true,
    "privacyMode": false,
    "showInjectedKeys": true
  },
  "appearance": {
    "monitor": "",
    "flexDirection": "column",
    "alignment": "bottom-right",
    "marginX": 100,
    "marginY": 100,
    "animation": "spring",
    "animationDuration": 0.25,
    "disablePressAnimation": true,
    "style": "",
    "opacity": 1.0
  },
  "layout": {
    "showIcon": true,
    "showSymbol": true,
    "showPressCount": true,
    "iconAlignment": "center"
  },
  "color": { "color": "", "secondaryColor": "", "useGradient": false },
  "modifier": {
    "highlight": false,
    "color": "",
    "secondaryColor": "",
    "textColor": "",
    "borderColor": ""
  },
  "text": { "size": 0, "color": "", "caps": "capitalize", "variant": "text-short", "alignment": "center" },
  "border": { "enabled": true, "color": "", "width": 2, "radius": 0 },
  "background": { "enabled": false, "color": "" },
  "mouse": {
    "showClicks": false,
    "size": 40,
    "color": "",
    "themePreset": "neon_green",
    "keepHighlight": false,
    "showIndicator": true,
    "keepIndicator": true,
    "indicatorSize": 135,
    "indicatorOffsetX": 35,
    "indicatorOffsetY": 35,
    "clickLingerMs": 60,
    "dragThreshold": 50,
    "opacity": 1.0
  },
  "app": {
    "startAtLogin": false,
    "startHidden": false,
    "trayClickAction": "settings",
    "excludedApps": []
  }
}
```

Empty color strings mean "inherit from the active theme". That is why a fresh
profile contains almost no colors: switching a design theme recolors the whole
overlay instead of leaving stale copies behind.

---

## Root

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `version` | `int` | `6` | Schema version. Older files are migrated automatically on load. |
| `language` | `string` | system | `zh-CN`, `en-US`, `zh-TW`, `ja-JP`. Empty follows the OS. |
| `themePreset` | `string` | `"Swiss Typography 1957"` | Design theme for keycaps, subtitles, and audio character. |
| `tactileAudio` | `bool` | `true` | Enable synthesized keystroke feedback. |
| `audioProfile` | `string` | `"chip"` | `auto`, `thock`, `click`, `scissor`, `glass`, `typewriter`, `hinoki`, `pebble`, `paper`, `felt`, `droplet`, `crystal`, `chip`, `shutter`. |
| `audioVolume` | `float` | `1.0` | Feedback volume, `0.0`–`1.0`. |
| `semanticSub` | `bool` | `true` | Show the human-readable action subtitle under each chord. |
| `semanticContext` | `bool` | `true` | Use the foreground app to decode shortcuts more precisely. |
| `showAppBadge` | `bool` | `true` | Include the recognized app name next to the subtitle. |
| `customShortcuts` | `array` | `[]` | User-defined shortcut interpretations. See below. |

### Custom shortcut entries

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `app` | `string` | no | Executable (`devcpp.exe`), profile id (`vscode`), category (`editor`), or `*`/omitted for a global rule. |
| `chord` | `string` | yes | One chord, for example `Ctrl+Shift+F9`. Modifiers: `Ctrl`, `Shift`, `Alt`, `Win`/`Super`/`Cmd`. |
| `label` | `string` | yes | Subtitle text. Unknown labels display verbatim. |
| `description` | `string` | no | Optional longer explanation. |

Rules are editable in **Settings › Intelligence › Custom Shortcuts**. User
rules take priority over every built-in interpretation. Invalid entries are
logged and skipped; they never disable the rest of the dictionary.

---

## General

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `general.filter` | `string` | `"none"` | Capture policy: `none` (all keys), `modifiers` (hotkeys and special keys only), `custom` (the allowed-key matrix), `silent` (paused). |
| `general.allowedKeys` | `string[]` | `[]` | Keys captured when `filter` is `custom`. |
| `general.showEventHistory` | `bool` | `true` | Keep a stack of recent chords on screen. |
| `general.maxHistory` | `int` | `5` | Maximum retained chords, `2`–`12`. |
| `general.lingerDurationMs` | `int` | `4500` | How long a chord stays visible, in milliseconds. |
| `general.toggleShortcut` | `string[]` | `["ShiftLeft", "F10"]` | Global hotkey that shows or hides the overlay. Record it in Settings instead of editing tokens by hand. |
| `general.wordGrouping` | `bool` | `true` | Merge consecutive alphanumeric keystrokes into words. |
| `general.privacyMode` | `bool` | `false` | Mask keystrokes typed into password, login, and banking windows. Off by default. |
| `general.showInjectedKeys` | `bool` | `true` | Also show keys synthesized by another program (AutoHotkey, remote desktop, key remappers). Needed for Total Commander's Vim mode, which is an AutoHotkey layer sending `{Down}` instead of `j`. Turn it off if a script floods the overlay. |

The deprecated `general.dragThreshold` is read only by the v6 migration; new
configurations store the value under `mouse.dragThreshold`.

`general.showInjectedKeys` only helps when the helper actually sends keystrokes.
Layers built on window messages — Total Commander's `cm_*` commands, for example
— never produce a key event at all, so no keyboard hook can observe them.

### Timing windows

Five durations decide when two inputs count as one gesture. They are listed
together because they interact: a word window longer than the repeat window
merges separate presses into one word, and a scroll window as long as the repeat
window lets a flick survive a pause that should have ended it.

| Window | Value | What it governs |
| :--- | :--- | :--- |
| Word grouping | `600 ms` | How soon a typed character may follow the previous one and still join the same word chip. |
| Repeat count | `1200 ms` | How soon the same key may be pressed again and still read as one action with a `×N` badge. |
| Shortcut sequence | `1500 ms` | How long the second chord of a two-chord shortcut (`Ctrl+K`, `Ctrl+S`) may take. |
| Double click | `400 ms` | How far apart two clicks may be and still count as a double click. Matches the Windows default. |
| Scroll gesture | `250 ms` | How soon a wheel notch may follow the previous one and still belong to the same gesture. Shorter than the repeat window on purpose: a flick is one motion, two separate scrolls are two. |

Scroll distance is measured in notches (`120` wheel units = one notch), so a
trackpad that reports a dozen small deltas and a wheel that reports one notch per
click both show an honest `×N`.

---

## Appearance (Stage)

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `appearance.monitor` | `string` | `""` | Target display: empty or `"primary"` for the primary screen, otherwise a display ID or label. |
| `appearance.flexDirection` | `string` | `"column"` | Chord flow: `row` or `column`. |
| `appearance.alignment` | `string` | `"bottom-center"` | Nine-grid anchor: `top-left`, `top-center`, `top-right`, `center-left`, `center`, `center-right`, `bottom-left`, `bottom-center`, `bottom-right`. |
| `appearance.marginX` / `marginY` | `int` | `100` | Inward margin from the anchored edge, in pixels. |
| `appearance.animation` | `string` | `"fade"` | `spring`, `fade`, `slide`, `float`, `zoom`, `none`. |
| `appearance.animationDuration` | `float` | `0.25` | Transition duration in seconds, `0.10`–`1.50`. |
| `appearance.disablePressAnimation` | `bool` | `false` | Disable the keycap depression and press highlight. |
| `appearance.style` | `string` | `""` | Physical keycap model: empty (theme default), `lowprofile`, `minimal`, `laptop`, `pbt`, `glass`, `neon`, `retro`, `outline`. |
| `appearance.opacity` | `float` | `1.0` | Keycap overlay opacity, `0.10`–`1.0`. |

---

## Layout, colors, typography, border, backdrop

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `layout.showIcon` | `bool` | `true` | Show modifier glyphs. |
| `layout.showSymbol` | `bool` | `true` | Show symbols such as Enter and Backspace glyphs. |
| `layout.showPressCount` | `bool` | `true` | Show the streak multiplier. |
| `layout.iconAlignment` | `string` | `"center"` | Icon position: `flex-start` or `center` put the icon before the label, `flex-end` puts it after. |
| `color.color` | `string` | `""` | Keycap background. Empty inherits the theme; a hex value overrides it. |
| `color.secondaryColor` | `string` | `""` | Keycap border/secondary color. |
| `color.useGradient` | `bool` | `false` | Blend the keycap background toward its accent color. |
| `modifier.highlight` | `bool` | `false` | Apply the modifier accent to modifier keys. |
| `modifier.color` | `string` | `""` | Modifier background. |
| `modifier.borderColor` | `string` | `""` | Modifier border. |
| `modifier.textColor` | `string` | `""` | Modifier label color. |
| `modifier.secondaryColor` | `string` | `""` | Reserved for future two-tone modifier accents. |
| `text.size` | `int` | `0` | Label size in pixels. `0` inherits the theme. |
| `text.color` | `string` | `""` | Label color. |
| `text.caps` | `string` | `"capitalize"` | `uppercase`, `capitalize`, `lowercase`. |
| `text.variant` | `string` | `"text-short"` | `text`, `text-short`, `icon`. |
| `text.alignment` | `string` | `"center"` | Horizontal label anchor inside the cap: `center-left`, `center`, `center-right`. |
| `border.enabled` | `bool` | `true` | Draw the keycap border. |
| `border.color` | `string` | `""` | Border color. |
| `border.width` | `int` | `2` | Border width in pixels. |
| `border.radius` | `float` | `0` | Corner radius. `0` inherits the theme. |
| `background.enabled` | `bool` | `false` | Draw a soft panel behind each chord. |
| `background.color` | `string` | `""` | Backdrop color; empty uses a translucent black. |

---

## Mouse

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `mouse.showClicks` | `bool` | `false` | Draw a ripple ring on button presses. |
| `mouse.size` | `int` | `150` | Cursor ring diameter in pixels. |
| `mouse.color` | `string` | `""` | Custom accent for the ring and buttons. Empty follows `themePreset`. |
| `mouse.themePreset` | `string` | `"neon_green"` | `neon_green`, `aurora_bicolor`, `cyber_cyan`, `obsidian_crimson`, `deep_blue`, `minimal_white`, `sakura_mist`, `arctic_ice`, `graphite_mono`, `amber_terminal`, `acid_lime`, `kyoto_lacquer`. |
| `mouse.keepHighlight` | `bool` | `false` | Keep a ring around the cursor while idle. |
| `mouse.showIndicator` | `bool` | `true` | Show button and wheel state beside the cursor. |
| `mouse.keepIndicator` | `bool` | `true` | Keep the indicator visible without input. |
| `mouse.indicatorSize` | `int` | `50` | Indicator size in pixels. |
| `mouse.indicatorOffsetX` / `Y` | `int` | `50` | Independent offset from the cursor tip, in pixels. |
| `mouse.clickLingerMs` | `int` | `180` | How long a button or wheel event stays visible. |
| `mouse.dragThreshold` | `int` | `50` | Cursor travel that counts as a drag, in pixels. |
| `mouse.opacity` | `float` | `1.0` | Mouse overlay opacity, `0.10`–`1.0`. |

---

## App behavior

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `app.startAtLogin` | `bool` | `false` | Register Cue in `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`. |
| `app.startHidden` | `bool` | `false` | Launch without showing the overlay; the global hotkey can reveal it. |
| `app.trayClickAction` | `string` | `"settings"` | What a single tray click does: `settings` or `toggle`. |
| `app.excludedApps` | `string[]` | `[]` | Executables that never contribute foreground context. Matching ignores case and `.exe`. |

Excluded apps still get universal shortcut decoding; only the window title and
application identity are withheld.

---

## Migration notes

- **v6** moves `general.dragThreshold` to `mouse.dragThreshold`, keeps the
  previously inert backdrop and gradient off, centers text and icon defaults,
  clears the legacy neon mouse color so the selected mouse theme applies, and
  fills the new tray action.
- **v5** adds the independent mouse overlay opacity.
- **v4** adds presentation controls: keycap opacity and keystroke volume.
- **v3** enables context intelligence defaults for existing installations.

A migration is written back to disk immediately after it runs.

---

## Managing the file from the app

**Settings › About › Configuration** can open the file, reveal its folder, copy
its path, export the whole configuration to the clipboard as JSON, and import a
pasted configuration. Import runs the same migrations as a file load, so a
shared v5 file is safe to paste into a v6 install. Invalid JSON is rejected
without touching the current settings.

The only reset control lives in the settings sidebar, so there is exactly one
way to restore factory defaults.
