<div align="center">

<img src="assets/logo.png" alt="Cue Logo" width="128" height="128" />

# Cue

**Every keystroke, on center stage.**  
**极简高质感按键与手势舞台提示工具，让每一次敲击与手势都清晰可见。**

[![Go Version](https://img.shields.io/badge/Go-1.27+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows%2010%2B%20(Direct3D%2011)-0078D6?style=flat-square&logo=windows)](https://microsoft.com/windows)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/yviscool/cue?style=flat-square)](https://github.com/yviscool/cue/releases)

[English](#english) · [简体中文](#简体中文)

</div>

---

<a name="english"></a>
## English

### What is Cue?

**Cue** (提示 / 信号 / 舞台提示) is an ultra-fast, minimalist keystroke and mouse stage cue tool built from scratch in **Pure Go (Zero CGO)** with **Direct3D 11 GPU hardware acceleration**.

Designed for educators, presenters, streamers, and developers, Cue illuminates presenter actions in real time—instantly turning physical keyboard shortcuts and mouse clicks into elegant, responsive stage cues that feel native, fluid, and distraction-free.

### Key Features

- **⚡ Pure Go & Zero CGO**: Built with pure Go runtime performance. Single lightweight binary (< 12 MB) with zero external DLL dependencies.
- **🎮 Direct3D 11 GPU Direct Rendering**: Ultra-low latency render pipeline running smoothly at 60 / 120 / 144+ FPS. Zero tearing, transparent pass-through overlay.
- **🎨 12 Master Preset Themes**:
  - *Liquid Glass* (visionOS frosted translucent refractive glass)
  - *MacBook Scissor* (Apple Space Gray scissor-switch)
  - *Dieter Rams 1960* (Braun warm gray & iconic orange)
  - *Cyberpunk Neon* (High-contrast cyan & magenta neon pulse)
  - *Classic Mechanical* (PBT double-shot keycaps)
  - *IBM Model M 1984* (Industrial retro buckling spring)
  - *Minimal Chalk*, *Obsidian OLED*, *Tokyo Blossom*, *Emerald Matrix*, *Sunset Horizon*, *Royal Velvet*
- **⌨️ 8 Physical Keycap Styles**: Low-Profile, Minimal, Laptop, PBT Mechanical, Liquid Glass, Cyber Neon, Retro Typewriter, Outline.
- **🧠 Semantic Action Recognition**: Intelligently detects and displays semantic human-readable action labels (e.g. `Ctrl + C` → *Copy*, `Alt + Tab` → *Switch Windows*, `Ctrl + +` → *Zoom In*).
- **🔤 Continuous Alphanumeric Word Grouping**: Smoothly groups typing strings into unified words while cleanly breaking on punctuation, space, and boundary keys.
- **📍 9-Grid Desktop Anchoring**: Dual-mode positioning (Classic 9-Grid selector + Interactive desktop mini-map with real-time X/Y offset sliders).
- **🖱️ Mouse Interaction Spotlight**: Real-time left/right/middle click indicators, drag visualization, and wheel scroll direction. Disabled by default for unobtrusive presentations.
- **🔒 Privacy Shield**: Automatically detects sensitive input dialogs and provides one-click privacy masking.
- **🌐 Native Multi-Language**: Built-in dynamic switcher for English, 简体中文, 繁體中文, and 日本語.

### Quick Start

#### Download Binary
Download the latest `cue.exe` directly from the [GitHub Releases](https://github.com/yviscool/cue/releases) page and run.

#### Build from Source
```bash
git clone https://github.com/yviscool/cue.git
cd cue
go build -ldflags="-H=windowsgui -s -w" -o cue.exe .
```

---

<a name="简体中文"></a>
## 简体中文

### 品牌理念

**Cue** —— 提示、信号、舞台提示。

让观众清晰看到演示者正在做什么。专为教学、技术分享、软件演示、录屏与直播场景倾力打造。短小、有极强记忆点，以极高审美标准重新定义按键可视化工具。

### 核心亮点

- **⚡ 原生 Pure Go 直驱 (Zero CGO)**：纯粹 Go 语言底层直连 Windows Win32 / Direct3D 11 API，单文件绿色免安装，启动毫秒级响应，内存占用极低。
- **🎮 GPU 硬件加速 direct 渲染**：基于 Direct3D 11 硬件管线渲染，支持 60Hz / 120Hz / 144Hz 高刷新率，透明无边框悬浮窗完全点击穿透，不抢占任何前台焦点。
- **🎨 12 款官方精选大师主题**：包括《液态冰晶 (Liquid Glass)》、《深空剪刀脚 (MacBook Scissor)》、《博朗经典 (Dieter Rams 1960)》、《赛博霓虹 (Cyberpunk Neon)》等全套高质感配色。
- **⌨️ 8 种实体键帽质感**：提供矮轴、极简粉笔、笔记本剪刀脚、PBT 双色注塑、拟态毛玻璃、荧光霓虹、复古机械与极细微光切边等多种键帽物理风格。
- **🧠 意图语义与快捷键识别**：内置 Universal 快捷键意图识别引擎，自动在悬浮窗下方展示操作意图（如“复制”、“粘贴”、“缩放”、“切换窗口”等）。
- **🔤 连续字母数字智能单词聚合**：输入长单词时自动凝聚为单一单词气泡，遇到空格、标点或控制键时利落断句分行。
- **📍 9 宫格屏幕吸附定位系统**：提供直观的九宫格方位选择器与全比例微缩屏幕定位小地图，支持精确边距像素微调。
- **🖱️ 鼠标手势与光标信号**：支持左键、右键、中键点击、滚轮滚动与拖拽手势动效；出厂默认静音关闭，随时按需开启。
- **🌐 完整四国语言原生动态切换**：支持简体中文、English、繁體中文与日本語，界面与悬浮窗字幕即时响应无缝刷新。

### 快速使用

1. 从 [Releases 页面](https://github.com/yviscool/cue/releases) 下载最新打包好的单文件 `cue.exe`。
2. 双击直接运行，系统托盘将出现黑白自适应 Cue 图标。
3. 默认全局显示/隐藏热键：`Shift + F10`。
4. 单击托盘图标或右键选择“通用设置”唤出现代化设置中心。

---

## 许可证 (License)

本项目采用 [MIT License](LICENSE) 开源协议。
