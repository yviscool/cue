package theme

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

// ThemeConfig holds the visual and physical properties for an overlay theme.
type ThemeConfig struct {
	Name        string
	DisplayName string
	Description string

	// Physical & Layout attributes
	DefaultStyle       string
	WindowBackground   ui.Color
	WindowCorner       float32
	KeyCapBackground   ui.Color
	KeyCapBorder       ui.Color
	KeyCapHeldBackground ui.Color
	KeyCapHeldBorder   ui.Color
	KeyCapTextColor    ui.Color
	KeyCapFontSize     float32
	KeyCapRadius       float32
	KeyCapShadowSpread float32

	// Modifier keycaps
	ModCapBackground ui.Color
	ModCapBorder     ui.Color
	ModCapTextColor  ui.Color

	// Subtitle & Action Intent Pill
	SubtitleBackground ui.Color
	SubtitleTextColor  ui.Color
	SubtitleFontSize   float32

	// Streak Badge (e.g. ×3)
	StreakBadgeBackground ui.Color
	StreakBadgeTextColor  ui.Color

	// Mouse Ripple
	MouseRippleColor ui.Color
}

// rgba creates a ui.Color from 0-255 values
func rgba(r, g, b, a uint8) ui.Color {
	return ui.Color{R: r, G: g, B: b, A: a}
}

// 12 Master Themes
var (
	ThemeLiquidGlass = ThemeConfig{
		Name:                  "Liquid Glass",
		DisplayName:           "液态冰晶 (Liquid Glass)",
		Description:           "高透光微磨砂拟态玻璃，visionOS 悬浮通透质感。",
		DefaultStyle:          "glass",
		WindowBackground:      rgba(15, 23, 42, 140),
		WindowCorner:          20,
		KeyCapBackground:      rgba(255, 255, 255, 38),
		KeyCapBorder:          rgba(255, 255, 255, 70),
		KeyCapHeldBackground:  rgba(56, 189, 248, 120),
		KeyCapHeldBorder:      rgba(186, 230, 253, 200),
		KeyCapTextColor:       rgba(248, 250, 252, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    8,
		ModCapBackground:      rgba(255, 255, 255, 25),
		ModCapBorder:          rgba(255, 255, 255, 50),
		ModCapTextColor:       rgba(226, 232, 240, 255),
		SubtitleBackground:    rgba(15, 23, 42, 210),
		SubtitleTextColor:     rgba(56, 189, 248, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(14, 165, 233, 220),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(56, 189, 248, 160),
	}

	ThemeMacBookScissor = ThemeConfig{
		Name:                  "MacBook Scissor",
		DisplayName:           "深空剪刀脚 (MacBook Scissor)",
		Description:           "Apple 经典深空灰磨砂氧化铝与低键程剪刀脚质感。",
		DefaultStyle:          "laptop",
		WindowBackground:      rgba(24, 24, 27, 230),
		WindowCorner:          16,
		KeyCapBackground:      rgba(39, 39, 42, 245),
		KeyCapBorder:          rgba(63, 63, 70, 220),
		KeyCapHeldBackground:  rgba(82, 82, 91, 245),
		KeyCapHeldBorder:      rgba(161, 161, 170, 220),
		KeyCapTextColor:       rgba(244, 244, 245, 255),
		KeyCapFontSize:        19,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    4,
		ModCapBackground:      rgba(24, 24, 27, 240),
		ModCapBorder:          rgba(63, 63, 70, 200),
		ModCapTextColor:       rgba(161, 161, 170, 255),
		SubtitleBackground:    rgba(39, 39, 42, 230),
		SubtitleTextColor:     rgba(250, 250, 250, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(161, 161, 170, 220),
		StreakBadgeTextColor:  rgba(9, 9, 11, 255),
		MouseRippleColor:      rgba(244, 244, 245, 140),
	}

	ThemeDieterRams = ThemeConfig{
		Name:                  "Dieter Rams 1960",
		DisplayName:           "博朗经典 (Dieter Rams 1960)",
		Description:           "博朗 60 年代机电工业设计：少即是多，标志性暖灰与博朗橙。",
		DefaultStyle:          "minimal",
		WindowBackground:      rgba(245, 245, 244, 240),
		WindowCorner:          14,
		KeyCapBackground:      rgba(255, 255, 255, 255),
		KeyCapBorder:          rgba(214, 211, 209, 255),
		KeyCapHeldBackground:  rgba(255, 107, 0, 255),
		KeyCapHeldBorder:      rgba(234, 88, 12, 255),
		KeyCapTextColor:       rgba(28, 25, 23, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    6,
		ModCapBackground:      rgba(231, 229, 228, 255),
		ModCapBorder:          rgba(214, 211, 209, 255),
		ModCapTextColor:       rgba(87, 83, 78, 255),
		SubtitleBackground:    rgba(28, 25, 23, 230),
		SubtitleTextColor:     rgba(255, 255, 255, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(255, 107, 0, 255),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(255, 107, 0, 180),
	}

	ThemeCyberpunkNeon = ThemeConfig{
		Name:                  "Cyberpunk Neon",
		DisplayName:           "赛博霓虹 (Cyberpunk Neon)",
		Description:           "高反差赛博暗黑底色与霓虹荧光脉冲边缘直绘。",
		DefaultStyle:          "neon",
		WindowBackground:      rgba(10, 10, 20, 230),
		WindowCorner:          16,
		KeyCapBackground:      rgba(18, 18, 35, 245),
		KeyCapBorder:          rgba(6, 182, 212, 220),
		KeyCapHeldBackground:  rgba(236, 72, 153, 200),
		KeyCapHeldBorder:      rgba(244, 114, 182, 255),
		KeyCapTextColor:       rgba(34, 211, 238, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    10,
		ModCapBackground:      rgba(12, 12, 28, 240),
		ModCapBorder:          rgba(168, 85, 247, 200),
		ModCapTextColor:       rgba(192, 132, 252, 255),
		SubtitleBackground:    rgba(15, 15, 30, 240),
		SubtitleTextColor:     rgba(236, 72, 153, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(236, 72, 153, 240),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(6, 182, 212, 200),
	}

	ThemeClassicMechanical = ThemeConfig{
		Name:                  "Classic Mechanical",
		DisplayName:           "复古机械 (Classic Mechanical)",
		Description:           "PBT 磨砂双色注塑键帽与微凹指腹触感仿真。",
		DefaultStyle:          "pbt",
		WindowBackground:      rgba(30, 41, 59, 230),
		WindowCorner:          16,
		KeyCapBackground:      rgba(51, 65, 85, 250),
		KeyCapBorder:          rgba(71, 85, 105, 255),
		KeyCapHeldBackground:  rgba(249, 115, 22, 230),
		KeyCapHeldBorder:      rgba(251, 146, 60, 255),
		KeyCapTextColor:       rgba(241, 245, 249, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    6,
		ModCapBackground:      rgba(30, 41, 59, 245),
		ModCapBorder:          rgba(51, 65, 85, 255),
		ModCapTextColor:       rgba(148, 163, 184, 255),
		SubtitleBackground:    rgba(15, 23, 42, 230),
		SubtitleTextColor:     rgba(249, 115, 22, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(249, 115, 22, 240),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(249, 115, 22, 180),
	}

	ThemeIBMModelM = ThemeConfig{
		Name:                  "IBM Model M 1984",
		DisplayName:           "传奇屈曲轴 (IBM Model M 1984)",
		Description:           "工业复古米白与青灰，致敬打字机黄金时代的屈曲弹簧轴。",
		DefaultStyle:          "retro",
		WindowBackground:      rgba(235, 231, 224, 240),
		WindowCorner:          12,
		KeyCapBackground:      rgba(245, 242, 237, 255),
		KeyCapBorder:          rgba(202, 196, 184, 255),
		KeyCapHeldBackground:  rgba(90, 107, 124, 240),
		KeyCapHeldBorder:      rgba(67, 83, 99, 255),
		KeyCapTextColor:       rgba(44, 42, 41, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    6,
		ModCapBackground:      rgba(218, 213, 203, 255),
		ModCapBorder:          rgba(188, 181, 168, 255),
		ModCapTextColor:       rgba(74, 71, 69, 255),
		SubtitleBackground:    rgba(44, 42, 41, 230),
		SubtitleTextColor:     rgba(245, 242, 237, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(90, 107, 124, 255),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(90, 107, 124, 180),
	}

	ThemeMinimalChalk = ThemeConfig{
		Name:                  "Minimal Chalk",
		DisplayName:           "极简粉笔白 (Minimal Chalk)",
		Description:           "纯粹白墨与轻柔轮廓线，克制而优雅的演示风格。",
		DefaultStyle:          "minimal",
		WindowBackground:      rgba(18, 18, 18, 200),
		WindowCorner:          16,
		KeyCapBackground:      rgba(255, 255, 255, 240),
		KeyCapBorder:          rgba(229, 231, 235, 255),
		KeyCapHeldBackground:  rgba(220, 38, 38, 220),
		KeyCapHeldBorder:      rgba(239, 68, 68, 255),
		KeyCapTextColor:       rgba(17, 24, 39, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    5,
		ModCapBackground:      rgba(243, 244, 246, 240),
		ModCapBorder:          rgba(209, 213, 219, 255),
		ModCapTextColor:       rgba(75, 85, 99, 255),
		SubtitleBackground:    rgba(17, 24, 39, 230),
		SubtitleTextColor:     rgba(255, 255, 255, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(220, 38, 38, 240),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(220, 38, 38, 180),
	}

	ThemeObsidianOLED = ThemeConfig{
		Name:                  "Obsidian OLED",
		DisplayName:           "黑曜极光 (Obsidian OLED)",
		Description:           "纯黑 OLED 像素深度与极窄高亮微光切边。",
		DefaultStyle:          "outline",
		WindowBackground:      rgba(0, 0, 0, 240),
		WindowCorner:          16,
		KeyCapBackground:      rgba(15, 15, 15, 255),
		KeyCapBorder:          rgba(60, 60, 60, 255),
		KeyCapHeldBackground:  rgba(40, 40, 40, 255),
		KeyCapHeldBorder:      rgba(255, 255, 255, 255),
		KeyCapTextColor:       rgba(255, 255, 255, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    4,
		ModCapBackground:      rgba(5, 5, 5, 255),
		ModCapBorder:          rgba(40, 40, 40, 255),
		ModCapTextColor:       rgba(160, 160, 160, 255),
		SubtitleBackground:    rgba(20, 20, 20, 240),
		SubtitleTextColor:     rgba(255, 255, 255, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(255, 255, 255, 230),
		StreakBadgeTextColor:  rgba(0, 0, 0, 255),
		MouseRippleColor:      rgba(255, 255, 255, 160),
	}

	ThemeTokyoBlossom = ThemeConfig{
		Name:                  "Tokyo Blossom",
		DisplayName:           "东京落樱 (Tokyo Blossom)",
		Description:           "浅绯色樱花与深紫夜色交织，日系和风优雅美学。",
		DefaultStyle:          "lowprofile",
		WindowBackground:      rgba(38, 20, 36, 230),
		WindowCorner:          18,
		KeyCapBackground:      rgba(66, 36, 62, 245),
		KeyCapBorder:          rgba(244, 114, 182, 160),
		KeyCapHeldBackground:  rgba(244, 63, 94, 210),
		KeyCapHeldBorder:      rgba(251, 113, 133, 255),
		KeyCapTextColor:       rgba(253, 242, 248, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    7,
		ModCapBackground:      rgba(48, 26, 45, 240),
		ModCapBorder:          rgba(244, 114, 182, 120),
		ModCapTextColor:       rgba(249, 168, 212, 255),
		SubtitleBackground:    rgba(38, 20, 36, 240),
		SubtitleTextColor:     rgba(244, 114, 182, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(244, 63, 94, 230),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(244, 114, 182, 180),
	}

	ThemeEmeraldMatrix = ThemeConfig{
		Name:                  "Emerald Matrix",
		DisplayName:           "代码矩阵 (Emerald Matrix)",
		Description:           "黑客帝国经典荧光绿代码流与深渊控制台。",
		DefaultStyle:          "retro",
		WindowBackground:      rgba(6, 20, 12, 240),
		WindowCorner:          12,
		KeyCapBackground:      rgba(12, 35, 20, 250),
		KeyCapBorder:          rgba(34, 197, 94, 180),
		KeyCapHeldBackground:  rgba(34, 197, 94, 220),
		KeyCapHeldBorder:      rgba(134, 239, 172, 255),
		KeyCapTextColor:       rgba(74, 222, 128, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    8,
		ModCapBackground:      rgba(8, 26, 15, 245),
		ModCapBorder:          rgba(34, 197, 94, 120),
		ModCapTextColor:       rgba(134, 239, 172, 255),
		SubtitleBackground:    rgba(6, 20, 12, 240),
		SubtitleTextColor:     rgba(74, 222, 128, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(34, 197, 94, 230),
		StreakBadgeTextColor:  rgba(0, 20, 10, 255),
		MouseRippleColor:      rgba(34, 197, 94, 200),
	}

	ThemeSunsetHorizon = ThemeConfig{
		Name:                  "Sunset Horizon",
		DisplayName:           "暮光地平线 (Sunset Horizon)",
		Description:           "加州落日橙黄至紫红渐层色彩，温暖而富有张力。",
		DefaultStyle:          "lowprofile",
		WindowBackground:      rgba(30, 18, 38, 230),
		WindowCorner:          18,
		KeyCapBackground:      rgba(55, 30, 65, 245),
		KeyCapBorder:          rgba(249, 115, 22, 180),
		KeyCapHeldBackground:  rgba(234, 88, 12, 220),
		KeyCapHeldBorder:      rgba(253, 186, 116, 255),
		KeyCapTextColor:       rgba(255, 247, 237, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    7,
		ModCapBackground:      rgba(40, 22, 48, 240),
		ModCapBorder:          rgba(249, 115, 22, 130),
		ModCapTextColor:       rgba(253, 186, 116, 255),
		SubtitleBackground:    rgba(30, 18, 38, 240),
		SubtitleTextColor:     rgba(249, 115, 22, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(234, 88, 12, 230),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(249, 115, 22, 190),
	}

	ThemeRoyalVelvet = ThemeConfig{
		Name:                  "Royal Velvet",
		DisplayName:           "皇家丝绒 (Royal Velvet)",
		Description:           "深邃典雅的帝王靛蓝与香槟流金点缀。",
		DefaultStyle:          "pbt",
		WindowBackground:      rgba(18, 20, 48, 235),
		WindowCorner:          16,
		KeyCapBackground:      rgba(30, 36, 85, 250),
		KeyCapBorder:          rgba(234, 179, 8, 180),
		KeyCapHeldBackground:  rgba(202, 138, 4, 220),
		KeyCapHeldBorder:      rgba(254, 240, 138, 255),
		KeyCapTextColor:       rgba(254, 249, 195, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    6,
		ModCapBackground:      rgba(22, 26, 62, 245),
		ModCapBorder:          rgba(234, 179, 8, 120),
		ModCapTextColor:       rgba(250, 204, 21, 255),
		SubtitleBackground:    rgba(18, 20, 48, 240),
		SubtitleTextColor:     rgba(250, 204, 21, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(202, 138, 4, 230),
		StreakBadgeTextColor:  rgba(15, 23, 42, 255),
		MouseRippleColor:      rgba(234, 179, 8, 180),
	}
)

var AllThemes = []ThemeConfig{
	ThemeLiquidGlass,
	ThemeMacBookScissor,
	ThemeDieterRams,
	ThemeCyberpunkNeon,
	ThemeClassicMechanical,
	ThemeIBMModelM,
	ThemeMinimalChalk,
	ThemeObsidianOLED,
	ThemeTokyoBlossom,
	ThemeEmeraldMatrix,
	ThemeSunsetHorizon,
	ThemeRoyalVelvet,
}

// ThemeByName retrieves a theme by name with fuzzy/case-insensitive matching.
func ThemeByName(name string) ThemeConfig {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	if cleaned == "" {
		return ThemeLiquidGlass
	}

	// Legacy alias compatibility
	switch cleaned {
	case "cupertino-pro":
		return ThemeMacBookScissor
	case "braun-tactile":
		return ThemeDieterRams
	}

	for _, t := range AllThemes {
		if strings.EqualFold(t.Name, name) ||
			strings.EqualFold(strings.ReplaceAll(t.Name, " ", "-"), cleaned) ||
			strings.EqualFold(t.DisplayName, name) {
			return t
		}
	}
	return ThemeLiquidGlass
}

// GetTheme retrieves a theme by name with fallback to Liquid Glass (alias of ThemeByName).
func GetTheme(name string) ThemeConfig {
	return ThemeByName(name)
}

// DefaultStyleForTheme returns the signature physical keycap style for the theme.
func DefaultStyleForTheme(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "liquid glass":
		return "glass"
	case "macbook scissor":
		return "laptop"
	case "dieter rams 1960":
		return "minimal"
	case "cyberpunk neon":
		return "neon"
	case "classic mechanical":
		return "pbt"
	case "ibm model m 1984":
		return "retro"
	case "minimal chalk":
		return "minimal"
	case "obsidian oled":
		return "outline"
	case "tokyo blossom":
		return "lowprofile"
	case "emerald matrix":
		return "retro"
	case "sunset horizon":
		return "lowprofile"
	case "royal velvet":
		return "pbt"
	default:
		return "lowprofile"
	}
}

// ColorToHex converts a ui.Color into a canonical #rrggbb or #rrggbbaa hex string.
func ColorToHex(c ui.Color) string {
	if c.A == 255 {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}
