# Cue Configuration Guide

Cue stores user settings in a human-readable JSON/YAML configuration file located at `%APPDATA%\Cue\config.json`. The configuration engine supports hot-reloading: all modifications take effect immediately without restarting the application.

---

## Configuration Schema

```json
{
  "version": 3,
  "language": "en",
  "themePreset": "Liquid Glass",
  "tactileAudio": true,
  "semanticSub": true,
  "semanticContext": true,
  "showAppBadge": true,
  "customShortcuts": [
    { "app": "devcpp.exe", "chord": "Ctrl+Alt+B", "label": "Build My Project" }
  ],
  "general": {
    "filter": "smart",
    "allowedKeys": [],
    "lingerDurationMs": 2000,
    "showEventHistory": true,
    "maxHistory": 5,
    "wordGrouping": true,
    "dragThreshold": 50,
    "toggleShortcut": ["Shift", "F10"]
  },
  "appearance": {
    "alignment": "bottom-right",
    "monitor": 0,
    "offsetX": 32,
    "offsetY": 32,
    "keyCapStyle": "liquid-glass",
    "animation": "spring",
    "scale": 1.0
  },
  "mouse": {
    "enabled": true,
    "showClicks": true,
    "size": 50,
    "color": "#00FF6A",
    "themePreset": "neon-green",
    "clickLingerMs": 180,
    "keepHighlight": false,
    "showIndicator": true,
    "keepIndicator": false,
    "indicatorSize": 40,
    "indicatorOffsetX": 24,
    "indicatorOffsetY": 24
  }
}
```

---

## Detailed Section Descriptions

### 1. Root & General

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `language` | `string` | `"zh-CN"` / `"en"` | Interface locale (`en`, `zh-CN`, `zh-TW`, `ja`). |
| `themePreset` | `string` | `"Liquid Glass"` | Active keyboard theme preset. |
| `tactileAudio` | `bool` | `true` | Toggle spatial binaural mechanical sound synthesis. |
| `audioProfile` | `string` | `"auto"` | Keystroke sound: `auto`, `thock`, `click`, `scissor`, `glass`, `typewriter`, `hinoki`, `pebble`, `paper`, `felt`, `droplet`, `crystal`, `chip`, `shutter`. |
| `semanticSub` | `bool` | `true` | Show human-readable action subtitles under keycaps. |
| `semanticContext` | `bool` | `true` | Interpret shortcuts using the foreground application. When disabled, universal decoding still works but app profiles are skipped. Privacy shielding remains active. |
| `showAppBadge` | `bool` | `true` | Include the recognized application name in action subtitles. |
| `customShortcuts` | `array` | `[]` | User-defined shortcut interpretations. See below. |
| `general.filter` | `string` | `"smart"` | Keystroke capture policy: `smart` (shortcuts & special keys), `all` (all typing), `custom`, `silent`. |
| `general.lingerDurationMs`| `int` | `2000` | Lifetime (in ms) keycap chords linger on screen. |
| `general.showEventHistory`| `bool` | `true` | When enabled, keeps a vertical stack of recent chords. |
| `general.maxHistory` | `int` | `5` | Maximum number of concurrent chords displayed on stage (2–12). |
| `general.wordGrouping` | `bool` | `true` | Groups consecutive alphanumeric typing into words. |
| `general.dragThreshold` | `int` | `50` | Minimum cursor travel (in px) before recognizing a drag. |
| `general.toggleShortcut` | `[]string` | `["Shift", "F10"]` | Global hotkey to toggle overlay visibility. |

#### Custom shortcut entries

Each `customShortcuts` item has the following shape:

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `app` | `string` | no | Executable (`devcpp.exe`), profile id (`vscode`), category (`editor`), or `*` / omitted for a global rule. |
| `chord` | `string` | yes | One chord, e.g. `Ctrl+Shift+F9`. Modifiers: `Ctrl`, `Shift`, `Alt`, `Win`/`Super`/`Cmd`. |
| `label` | `string` | yes | Subtitle text. Also acts as the translation key; unknown labels display verbatim. |
| `description` | `string` | no | Optional longer explanation for future UI surfaces. |

User overrides take priority over every built-in interpretation. Invalid
entries are logged and skipped; they never disable the rest of the dictionary.

> **Migration:** configuration files written before version 3 are upgraded
> automatically. `semanticContext` and `showAppBadge` default to `true` for
> existing installations, and the file is rewritten in place.

---

### 2. Appearance & Positioning

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `appearance.alignment` | `string` | `"bottom-right"` | Screen edge anchoring (`top-left`, `top-center`, `top-right`, `center-left`, `center`, `center-right`, `bottom-left`, `bottom-center`, `bottom-right`). |
| `appearance.monitor` | `int` | `0` | Target display index (0 = Primary monitor). |
| `appearance.offsetX` / `Y`| `int` | `32` | Distance (px) inward from the anchored display boundary. |
| `appearance.keyCapStyle` | `string` | `"liquid-glass"` | Physical keycap rendering aesthetic (`liquid-glass`, `minimal`, `laptop`, `pbt-mechanical`, `cyber-neon`, `retro-typewriter`, `outline`). |
| `appearance.animation` | `string` | `"spring"` | Transition motion curve (`spring`, `slide`, `fade`, `none`). |
| `appearance.scale` | `float` | `1.0` | Global scale factor (0.5 to 2.5). |

---

### 3. Mouse Feedback & Gestures

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `mouse.showClicks` | `bool` | `true` | Visualizes contact ripple rings around the cursor. |
| `mouse.themePreset` | `string` | `"neon-green"` | Mouse theme preset: `neon-green`, `aurora`, `cyber-cyan`, `crimson`, `deep-blue`, `minimal-white`. |
| `mouse.clickLingerMs` | `int` | `180` | Milliseconds the button lobe remains highlighted after release. |
| `mouse.keepHighlight` | `bool` | `false` | Keeps a steady spotlight halo around the cursor when idling. |
| `mouse.showIndicator` | `bool` | `true` | Shows the companion mouse silhouette beside the cursor during interaction. |
| `mouse.keepIndicator` | `bool` | `false` | Keeps the mouse companion silhouette permanently visible. |
| `mouse.indicatorSize` | `int` | `40` | Height/scale of the companion mouse silhouette. |
| `mouse.indicatorOffsetX` / `Y` | `int` | `24` | Offset distance (px) from the cursor tip to the mouse silhouette. |
