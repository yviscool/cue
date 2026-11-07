# Context Intelligence, Shortcut Decoders & Privacy Shield

Cue's semantic layer translates raw physical key chords into intuitive,
human-readable operation labels. It goes beyond a static lookup table: the
foreground application is detected continuously, so `Ctrl + P` can say
**Print** in a document, **Quick Open** in VS Code, and **Command Palette** in
Obsidian—without any configuration.

---

## 1. How Context Intelligence Works

```
foreground HWND/title  ──►  context.Provider (250 ms poll)
        │                         │  only on HWND/title change
        │                         ▼
        │                  semantic.Registry
        │                  ├─ user overrides
        │                  ├─ exact app profile   (Dev-C++, IDLE, VS Code…)
        │                  ├─ app category        (editor, browser, terminal…)
        │                  └─ universal dictionary
        ▼
  Engine ──► HUD subtitle + optional application badge
```

Design guarantees:

- **Hook-safe.** Win32 foreground/process queries run on a dedicated provider
  goroutine. The low-level keyboard hook never performs a window or process
  query.
- **Cheap per keystroke.** Profile resolution and sensitive-window keyword
  scanning are cached per foreground context (`AppContext.Seq`); a keystroke
  costs only a few map lookups.
- **Deterministic priority.**
  `User override > exact app profile > app category > universal > no label`.
- **Exact IDLE detection.** `python.exe` / `pythonw.exe` is ambiguous, so the
  IDLE profile additionally requires an `IDLE` title hint. A normal Python
  script is never mislabeled as IDLE.
- **Privacy always on.** Disabling contextual interpretation removes application
  profiles only; universal decoding and the privacy shield stay active.

### Confidence levels

| Level | Meaning |
| :--- | :--- |
| **User** | Your `customShortcuts` entry wins over everything. |
| **Profile** | Exact application match (e.g. Dev-C++ `F9` → Compile). |
| **Category** | Family fallback (e.g. any editor: `Ctrl + G` → Go to Line). |
| **Universal** | Meaning shared by virtually every Windows app. |

---

## 2. Universal Decoders

| Chord | Recognized Action (EN) | 识别意图 (ZH) | 認識された意図 (JA) |
| :--- | :--- | :--- | :--- |
| `Ctrl + C` | **Copy** | 复制 | コピー |
| `Ctrl + V` | **Paste** | 粘贴 | 貼り付け |
| `Ctrl + X` | **Cut** | 剪切 | 切り取り |
| `Ctrl + Z` | **Undo** | 撤销 | 元に戻す |
| `Ctrl + Y` / `Ctrl + Shift + Z` | **Redo** | 重做 | やり直し |
| `Ctrl + A` | **Select All** | 全选 | すべて選択 |
| `Ctrl + S` | **Save** | 保存 | 保存 |
| `Ctrl + Shift + S` | **Save As** | 另存为 | 名前を付けて保存 |
| `Ctrl + F` | **Find** | 查找 | 検索 |
| `Ctrl + H` | **Replace** | 替换 | 置換 |
| `Ctrl + P` | **Print** (apps may override) | 打印 | 印刷 |
| `Alt + Tab` | **Switch App** | 切换应用 | アプリ切替 |
| `Alt + F4` | **Close Window** | 关闭窗口 | ウィンドウを閉じる |
| `Win + D` | **Show Desktop** | 显示桌面 | デスクトップ表示 |
| `Win + E` | **File Explorer** | 文件资源管理器 | ファイルエクスプローラー |
| `Win + L` | **Lock PC** | 锁定电脑 | 画面ロック |
| `Win + V` | **Clipboard History** | 剪贴板历史 | クリップボード履歴 |
| `Win + Shift + S` | **Snipping Tool** | 截图工具 | 切り取りツール |
| `Win + ←/→/↑/↓` | **Snap / Maximize / Restore** | 贴靠 / 最大化 / 还原 | スナップ / 最大化 / 元に戻す |
| `Ctrl + Shift + Esc` | **Task Manager** | 任务管理器 | タスクマネージャー |
| `Ctrl + +` / `Ctrl + =` | **Zoom In** | 放大 | 拡大 |
| `Ctrl + -` | **Zoom Out** | 缩小 | 縮小 |
| `Ctrl + 0` | **Reset Zoom** | 还原缩放 | 拡大リセット |
| `PrintScreen` | **Screenshot** | 截屏 | スクリーンショット |

---

## 3. Built-in Application Profiles

### Dev-C++ (`devcpp.exe`)

| Chord | Action |
| :--- | :--- |
| `F9` | Compile |
| `F10` | Run |
| `F11` | Compile & Run |
| `F5` | Debug |
| `F6` | Stop Execution |
| `Ctrl + Shift + F` | Find in Files |
| `Ctrl + /` | Toggle Comment |
| `Ctrl + Space` | Code Completion |

### IDLE (`python.exe` / `pythonw.exe` + “IDLE” in the title)

| Chord | Action |
| :--- | :--- |
| `F5` | Run Module |
| `Ctrl + N` | New File |
| `Ctrl + O` | Open File |
| `Alt + G` | Go to Line |
| `Alt + 3` | Comment Region |
| `Alt + 4` | Uncomment Region |
| `Ctrl + ]` | Indent Region |
| `Ctrl + [` | Dedent Region |

### Other first-class profiles

| Application | Example interpretation |
| :--- | :--- |
| **VS Code / Cursor / VSCodium / Windsurf** | `Ctrl + Shift + P` → Command Palette, `Ctrl + P` → Quick Open, `Alt + ↑/↓` → Move Line |
| **JetBrains IDEs** | `Ctrl + B` → Go to Declaration, `Ctrl + Alt + L` → Reformat Code |
| **Visual Studio** | `F5` → Start Debugging, `Ctrl + Shift + B` → Build Solution |
| **Eclipse** | `F3` → Open Declaration, `Ctrl + Shift + T` → Open Type |
| **Sublime Text / Zed / Notepad++** | Command Palette, Goto Anything, duplicate/delete line |
| **Browsers** (Chrome, Edge, Firefox, Brave…) | `Ctrl + T` → New Tab, `Ctrl + D` → Bookmark Tab, `F12` → Developer Tools |
| **Terminals** (Windows Terminal, PowerShell, cmd, WezTerm…) | `Ctrl + C` → Interrupt Process, `Ctrl + L` → Clear Terminal |
| **Microsoft Office / WPS** | Word formatting, Excel `Ctrl + 1` → Format Cells, PowerPoint `F5` → Start Slideshow |
| **Figma / Photoshop / Illustrator** | Group/Ungroup, Duplicate, Free Transform, Clipping Mask |
| **Premiere / VLC** | Export Media, Play/Pause, Seek, Volume |
| **Slack / Teams / Discord / WeChat** | Quick Switcher, Toggle Mute, Raise Hand |
| **Obsidian / Notion** | Quick Open, Command Palette, Reading View |
| **File Explorer** | New Folder, Focus Address Bar, Rename, Up One Level |

---

## 4. Two-Chord Sequences

Editors commonly use chorded shortcuts. Cue matches the pair within **1.5
seconds** and collapses both chords into one cue with a `→` separator:

| Sequence | Action |
| :--- | :--- |
| `Ctrl + K` `Ctrl + S` | Save All (VS Code) |
| `Ctrl + K` `Ctrl + C` | Add Line Comment (VS Code) |
| `Ctrl + K` `Ctrl + U` | Remove Line Comment (VS Code) |
| `Ctrl + K` `Ctrl + D` | Format Document (Visual Studio) |
| `Ctrl + K` `Ctrl + C` | Comment Selection (Visual Studio) |

If the second chord arrives after the window, or the pair is unknown, Cue
simply shows the two chords normally.

---

## 5. Application Badge

When **Show application badge** is enabled, each subtitle includes the name of
the application that produced the interpretation (for example `VS Code` or
`Dev-C++`). This makes the context auditable at a glance and is especially
useful when demonstrating multi-app workflows.

Both **Context intelligence** and **Show application badge** can be toggled in
**Settings → General → Action Intelligence**. The same panel shows the
currently recognized application in real time.

---

## 6. Custom Shortcuts

Any built-in interpretation can be overridden from `config.json`. `app` accepts
an executable name, a profile id, a category, or `*` / empty for a global rule.
User overrides always win.

```json
{
  "customShortcuts": [
    {
      "app": "devcpp.exe",
      "chord": "Ctrl+Alt+B",
      "label": "Build My Project",
      "description": "Team build script"
    },
    {
      "app": "editor",
      "chord": "Ctrl+Shift+Q",
      "label": "Quick Note"
    },
    {
      "chord": "Ctrl+Alt+P",
      "label": "Presenter Mode"
    }
  ]
}
```

Valid tokens: `Ctrl`/`Control`, `Shift`, `Alt`, `Win`/`Super`/`Cmd`, plus one
primary key (`A–Z`, `0–9`, `F1–F12`, `Space`, `Enter`, `Esc`, arrows, `+`, `-`,
`/`, `` ` `` …). Invalid entries are logged and ignored without affecting the
rest of the dictionary.

---

## 7. Privacy Shield

To protect presenters during live webcasts or recordings, Cue monitors both the
foreground context and the typed keystroke pattern:

- **Sensitive windows**: password managers, sign-in pages, banking/payment
  pages, UAC/Windows Security dialogs, private/incognito browsing windows.
- **Sensitive typed patterns**: `password`, `passwd`, `pwd:`, `secret`,
  `token`, `api_key`, `bearer`, `private_key`, `sudo`, `cvv`, and long digit
  runs (8+ digits).

When triggered, Cue suppresses the physical character keys and displays a
single protected placeholder:

```
[🔒 Private]
```

The context keyword scan is cached per foreground change, so the privacy layer
adds no measurable cost to typing.

### Capture modes

- **Smart Mode** (default): shows modifier combinations and special navigation
  keys while suppressing plain typing.
- **Filter All**: shows every keystroke with word grouping and backspace
  absorption.
- **Silent Mode**: suppresses all key overlays (mouse indicators still work).
- **Custom Mode**: shows only the keys in `general.allowedKeys`.

---

## 8. Performance Notes

| Path | Cost |
| :--- | :--- |
| Foreground poll | 250 ms, two cheap Win32 calls when unchanged |
| Process/executable query | only when HWND or title changes |
| Per-keystroke lookup | a few map reads; no allocation, no window scan |
| Sensitive-window scan | cached per `AppContext.Seq` |
| Context switch | clears pending sequence state |

The provider clamps its interval to 50 ms–2 s and ignores Cue's own windows, so
opening Settings never changes the interpretation of the app behind it.

---

## 9. Full Shortcut Inventory

This page is the curated tour. The complete, always-current list of every chord
Cue can interpret lives in **[shortcut-reference.md](shortcut-reference.md)**,
which is generated from the registry itself:

| Layer | What it holds |
| :--- | :--- |
| Universal | 102 chords that mean the same thing in nearly every Windows app |
| Category | Editor / browser / terminal fallbacks for apps without an exact profile |
| Profiles | 27 application families, from VS Code to OBS Studio |

Every profile records the keymap it was transcribed from and the build it was
verified against, so a wrong chord can be traced and fixed instead of being
believed forever.

Regenerate the reference after changing the registry:

```bash
CUE_UPDATE_REFERENCE=1 go test ./internal/semantic -run TestShortcutReference
```

`TestShortcutReference` fails when the checked-in document drifts from the code,
and `TestProfileCoverageFloors` fails when a profile silently gets thinner.

### Writing a chord

Shifted symbols are normalised, so these are equivalent:

| Written as | Also matches |
| :--- | :--- |
| `Ctrl+Shift+1` | `Ctrl+!` |
| `Ctrl+Shift+\` | `Ctrl+\|` |
| `Ctrl+Shift+;` | `Ctrl+:` |

Spelled-out punctuation is accepted too (`Ctrl+Slash`, `Ctrl+Period`,
`Ctrl+Backtick`), which keeps hand-written `customShortcuts` readable.
