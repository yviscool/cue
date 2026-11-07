# Cue System Architecture

Cue is an ultra-fast, minimalist keystroke and gesture visualizer designed specifically for live presentations, screencasts, and developer streaming. Built with a **Pure Go (Zero CGO)** runtime, it eliminates external native compilation toolchains while delivering Direct3D 11 hardware-accelerated graphics at full display refresh rate.

---

## 1. High-Level Architecture Overview

```
                      +-----------------------------+
                      |   Windows OS Event Hook     |
                      |  (WH_KEYBOARD_LL / MOUSE)   |
                      +--------------+--------------+
                                     |
                                     v
                      +-----------------------------+
                      |      Input Normalizer       |
                      |  (Key Names, Mods, Drag)    |
                      +--------------+--------------+
                                     |
                                     v
  +----------------+      +-----------------------------+
  | App Context    |      |      Cue State Engine       |
  | Provider       |----->| - Cadence Decay Loop        |
  | - 250 ms poll  |      | - Streak & Multi-Click      |
  | - change-driven|      | - Word Grouping & Backspace |
  |   process query|      | - Sequence Collapsing       |
  +----------------+      | - Semantic Intent Matching  |
                          +-------+--------------+------+
                                  |              |
                     Snapshot Emit|              | Acoustic Trigger
                                  v              v
                 +--------------------------+  +-------------------------+
                 |   MyGo UI Scene Graph    |  | Spatial Acoustic Engine |
                 | - Retained Tree Elements |  | - Stereo Stage Panning  |
                 | - Spring Physics Motions |  | - Dynamic Thock/Clicks  |
                 | - Subtitle + App Badge   |  +-------------------------+
                 | - D3D11 Pass-Through     |
                 +--------------------------+
```

---

## 2. Core Subsystems

### 2.1 Pure Go & Direct3D 11 GPU Rendering Pipeline
- **Zero CGO**: Cue interfaces directly with Windows system libraries (`user32.dll`, `d3d11.dll`, `dwmapi.dll`) via pure Go syscalls without requiring GCC or MinGW.
- **Hardware Acceleration**: Window surfaces are backed by Direct3D 11 swapchains, allowing transparent overlays to draw at 60Hz, 120Hz, or 144Hz with zero frame tearing and negligible CPU overhead.
- **Pass-Through Compositing**: Windows extended styles (`WS_EX_LAYERED`, `WS_EX_TRANSPARENT`, `WS_EX_NOACTIVATE`) ensure Cue never steals input focus or disrupts the presenter's active application.

### 2.2 Low-Level OS Hook (`internal/hook`)
- Operates on a dedicated OS message pump capturing `WH_KEYBOARD_LL` and `WH_MOUSE_LL` events.
- Synthesizes physical modifier states (`Shift`, `Ctrl`, `Alt`, `Super`) and normalizes key symbols across multiple international keyboard layouts.
- Detects drag actions when cursor displacement exceeds the configured `DragThreshold`.

### 2.3 Cue State Machine & Cadence Decay Engine (`internal/engine`)
- **Word Grouping**: Successive alphanumeric keystrokes are seamlessly aggregated into a single typing chord until a word boundary (space, punctuation, enter) is struck.
- **Backspace Absorption**: Typing backspace automatically trims the trailing character of an active word chord, keeping the presentation mirror clean.
- **Streak & Multi-Click Counter**: Rapid repeated presses of the same key or mouse button increment a streak badge (`×2`, `×3`, `×4`) without flooding the screen with duplicate capsules.
- **High-Resolution Retirement Timers**: Click contact highlights dim dynamically at `ClickLingerMs`, while multi-click badges persist for the full viewer recognition window (`1200ms`).
- **Two-Chord Sequences**: A chord may be remembered as the first half of a sequence for `1.5s`. If the matcher recognizes the pair (e.g. `Ctrl+K Ctrl+S`), both cues collapse into one combined chord with a `→` separator instead of showing two unrelated key clusters. Any foreground change, plain key, or privacy trigger clears the pending state.
- **Privacy-First Ordering**: Sensitivity is evaluated once per keystroke, before sequence matching, so a private context can never be merged into a shortcut hint.

### 2.4 Foreground Application Perception (`internal/context`)
- **Polling, not hooking**: A dedicated goroutine samples `GetForegroundWindow` every `250ms` (clamped to 50ms–2s). It never runs on the low-level hook thread, so input latency is unaffected.
- **Change-driven work**: Window titles are compared cheaply on every tick. The expensive process query (`OpenProcess` + `QueryFullProcessImageNameW`) only runs when the HWND actually changes.
- **UWP resolution**: `ApplicationFrameHost` windows are resolved to their real child `Windows.UI.Core.CoreWindow` process, so packaged apps are identified correctly.
- **Self-filtering**: Cue's own settings and tray windows publish an empty context, so opening Settings never changes the interpretation of the app behind it.
- **Snapshot contract**: Subscribers receive an immutable `AppContext` (`HWND`, `PID`, `Exe`, `Title`, `Category`, `Seq`). `Seq` increments on every real change and is the cache key used by the semantic layer.
- **Static provider**: `context.NewStatic` provides the same contract for tests and non-Windows builds.

### 2.5 Contextual Semantic Shortcut Matching (`internal/semantic`)
- **Deterministic resolution order**: user override → exact application profile → application category → universal dictionary → no label. The chosen level is exposed as `Confidence` and can be shown in the HUD.
- **Registry-driven profiles**: `internal/semantic/profiles.go` is the single data source for VS Code/Cursor, JetBrains, Visual Studio, Eclipse, Sublime, Zed, Notepad++, **Dev-C++**, **IDLE**, Office/WPS, browsers, terminals, Figma/Photoshop/Illustrator, Premiere/VLC, Slack/Teams/Discord/WeChat, Obsidian/Notion and File Explorer, plus category-wide and universal fallbacks.
- **Ambiguous executables**: `python.exe` / `pythonw.exe` are shared by IDLE and ordinary scripts, so the IDLE profile additionally requires an `idle` or ` shell` title hint. `idle.exe` is unambiguous. A normal Python process is never mislabeled.
- **Caching**: profile resolution and sensitive-window keyword scans are cached per `AppContext.Seq`; a keystroke costs only a few map lookups and allocates nothing on the hot path.
- **User overrides**: `customShortcuts` entries are parsed once on config apply and pre-indexed by app + chord, so `ParseOverride` errors can never break built-in matching.
- **Localization**: every action title emitted by the registry is translated into all four supported languages by `internal/i18n/actions.go`; a completeness test fails the build if a new action lacks a translation.
- **Privacy Shield**: Detects sensitive windows (password managers, sign-in and banking pages, UAC/Windows Security, private/incognito) and typed patterns (`password`, `token`, `api_key`, long digit runs…), automatically masking private chords with a lock glyph (`[🔒 Private]`). Disabling context intelligence only affects subtitles; the shield always stays on.

### 2.6 Spatial Acoustic Stage (`internal/acoustic`)
- Replaces generic beeps with a real-time synthesized acoustic model.
- Keystrokes generate dynamic mechanical "Thock" (modifier press) or crisp "Click" (character actuation) sound waves.
- Stereo panning maps the horizontal origin of keys across the virtual soundstage based on keyboard geometry.
