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
                      +-----------------------------+
                      |      Cue State Engine       |
                      | - Cadence Decay Loop        |
                      | - Streak & Multi-Click      |
                      | - Word Grouping & Backspace |
                      | - Semantic Intent Matching  |
                      +-------+--------------+------+
                              |              |
             Snapshot Emit    |              | Acoustic Trigger
                              v              v
     +--------------------------+          +-------------------------+
     |   MyGo UI Scene Graph    |          | Spatial Acoustic Engine |
     | - Retained Tree Elements |          | - Stereo Stage Panning  |
     | - Spring Physics Motions |          | - Dynamic Thock/Clicks  |
     | - D3D11 Pass-Through     |          +-------------------------+
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

### 2.4 Semantic Shortcut Intent Matching (`internal/semantic`)
- Inspects active chord combinations against an intelligent dictionary of universal shortcuts.
- Translates raw key combinations into human-understandable intents (e.g., `Ctrl + C` $\rightarrow$ **Copy**, `Ctrl + V` $\rightarrow$ **Paste**, `Alt + Tab` $\rightarrow$ **Switch Window**, `Ctrl + Z` $\rightarrow$ **Undo**).
- **Privacy Shield**: Detects sensitive passwords and authentication patterns, automatically masking private chords with a lock glyph (`[🔒 Private]`).

### 2.5 Spatial Acoustic Stage (`internal/acoustic`)
- Replaces generic beeps with a real-time synthesized acoustic model.
- Keystrokes generate dynamic mechanical "Thock" (modifier press) or crisp "Click" (character actuation) sound waves.
- Stereo panning maps the horizontal origin of keys across the virtual soundstage based on keyboard geometry.
