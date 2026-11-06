# Universal Shortcut Decoders & Privacy Shield

Cue incorporates a semantic interpretation engine that translates raw physical key chords into intuitive human-readable operation labels, making presentations instantly understandable for any audience.

---

## 1. Universal Shortcut Decoders

When a standard chord is pressed, Cue attaches a secondary action subtitle underneath the keycaps:

| Chord | Recognized Action (EN) | 识别意图 (ZH) | 認識された意図 (JA) |
| :--- | :--- | :--- | :--- |
| `Ctrl + C` | **Copy** | 复制 | コピー |
| `Ctrl + V` | **Paste** | 粘贴 | 貼り付け |
| `Ctrl + X` | **Cut** | 剪切 | 切り取り |
| `Ctrl + Z` | **Undo** | 撤销 | 元に戻す |
| `Ctrl + Y` / `Ctrl + Shift + Z` | **Redo** | 重做 | やり直し |
| `Ctrl + A` | **Select All** | 全选 | すべて選択 |
| `Ctrl + S` | **Save** | 保存 | 保存 |
| `Ctrl + F` | **Find** | 查找 | 検索 |
| `Ctrl + H` | **Replace** | 替换 | 置換 |
| `Ctrl + P` | **Print / Quick Open** | 打印 / 命令面板 | 印刷 / コマンドパレット |
| `Ctrl + Shift + P` / `F1` | **Command Palette** | 命令面板 | コマンドパレット |
| `Alt + Tab` | **Switch Windows** | 切换窗口 | ウィンドウ切替 |
| `Win + D` | **Show Desktop** | 显示桌面 | デスクトップ表示 |
| `Win + E` | **File Explorer** | 文件资源管理器 | ファイルエクスプローラー |
| `Win + L` | **Lock Screen** | 锁定屏幕 | 画面ロック |
| `Ctrl + +` / `Ctrl + =` | **Zoom In** | 放大 | 拡大 |
| `Ctrl + -` | **Zoom Out** | 缩小 | 縮小 |
| `Ctrl + 0` | **Reset Zoom** | 还原缩放 | 拡大リセット |

---

## 2. Privacy Shield

To protect presenter security during live webcasts or video recordings, Cue monitors active system input contexts:

- **Sensitive Password Fields**: When typing into known password fields or credential dialogs, Cue suppresses the physical character keys and displays a single protected placeholder:
  ```
  [🔒 Private]
  ```
- **Silent & Custom Filter Modes**:
  - `Smart Mode` (Default): Captures combinations with modifiers (`Ctrl`, `Alt`, `Shift`, `Win`) and special navigation keys while suppressing pure typing.
  - `Filter All`: Captures all keystrokes, with automatic word grouping and backspace absorption.
  - `Silent Mode`: Suppresses all key overlays (only mouse indicators appear if enabled).
  - `Custom Mode`: Whitelist only specific keys defined in `%APPDATA%\Cue\config.json`.
