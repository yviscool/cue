package theme

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

// HeldPhysicsMode defines how the keycap behaves when held or pressed.
type HeldPhysicsMode string

const (
	// HeldSink represents mechanical physical descent (travel distance downward).
	HeldSink HeldPhysicsMode = "sink"
	// HeldSnappy represents ultra-responsive shallow scissor descent.
	HeldSnappy HeldPhysicsMode = "snappy"
	// HeldPulse represents luminous optical/photon compression without downward displacement.
	HeldPulse HeldPhysicsMode = "pulse"
)

// ThemeConfig holds the visual, physical, and acoustic archetype properties for an overlay theme.
type ThemeConfig struct {
	Name        string
	DisplayName string
	Description string

	// Signature Physical Archetype
	DefaultStyle        string          // Physical geometry: "glass", "laptop", "minimal", "pbt", "neon", "retro", "outline", "lowprofile"
	TravelDepth         float32         // Physical keypress displacement in pixels (0.0 to 4.0)
	HeldPhysics         HeldPhysicsMode // Dynamic held response: "sink", "snappy", "pulse"
	AcousticProfile     string          // Signature acoustic switch mapping: "thock", "click", "glass", "scissor"
	SpringBounce        float64         // Spring animation bounce factor (0.10 - 0.35)
	SpringDurationScale float64         // Multiplier for transition speed (0.75 - 1.25)

	// Window & Canvas geometry
	WindowBackground ui.Color
	WindowCorner     float32

	// KeyCap Surfaces & Borders
	KeyCapBackground     ui.Color
	KeyCapBorder         ui.Color
	KeyCapHeldBackground ui.Color
	KeyCapHeldBorder     ui.Color
	KeyCapTextColor      ui.Color
	KeyCapFontSize       float32
	KeyCapRadius         float32
	KeyCapShadowSpread   float32

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

// EffectiveTravelDepth returns the physical downward displacement in pixels.
func (t ThemeConfig) EffectiveTravelDepth() float32 {
	if t.HeldPhysics == HeldPulse {
		return 0.0
	}
	if t.TravelDepth > 0 {
		return t.TravelDepth
	}
	switch t.DefaultStyle {
	case "glass", "neon", "outline":
		return 0.0
	case "laptop":
		return 1.2
	case "minimal":
		return 1.8
	case "pbt":
		return 3.0
	case "retro":
		return 3.8
	default:
		return 2.2
	}
}

// EffectiveSpringBounce returns the spring bounce factor for motion animations.
func (t ThemeConfig) EffectiveSpringBounce() float64 {
	if t.SpringBounce > 0 {
		return t.SpringBounce
	}
	return 0.18
}

// EffectiveDurationScale returns the duration scaling factor for transitions.
func (t ThemeConfig) EffectiveDurationScale() float64 {
	if t.SpringDurationScale > 0 {
		return t.SpringDurationScale
	}
	return 1.0
}

// rgba creates a ui.Color from 0-255 values
func rgba(r, g, b, a uint8) ui.Color {
	return ui.Color{R: r, G: g, B: b, A: a}
}

// Curated keyboard themes
var (
	ThemeLiquidGlass = ThemeConfig{
		Name:                  "Liquid Glass",
		DisplayName:           "液态冰晶 (Liquid Glass)",
		Description:           "高透光微磨砂拟态玻璃，visionOS 悬浮通透质感。",
		DefaultStyle:          "glass",
		TravelDepth:           0.0,
		HeldPhysics:           HeldPulse,
		AcousticProfile:       "droplet",
		SpringBounce:          0.22,
		SpringDurationScale:   1.0,
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
		TravelDepth:           1.2,
		HeldPhysics:           HeldSnappy,
		AcousticProfile:       "scissor",
		SpringBounce:          0.15,
		SpringDurationScale:   0.80,
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
		TravelDepth:           2.0,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "pebble",
		SpringBounce:          0.18,
		SpringDurationScale:   0.95,
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
		TravelDepth:           0.0,
		HeldPhysics:           HeldPulse,
		AcousticProfile:       "crystal",
		SpringBounce:          0.28,
		SpringDurationScale:   0.90,
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

	// ThemeAzurePulse is the all-blue sibling of Cyberpunk Neon. It keeps the
	// same luminous navy keycap plate and the same azure used for the letter
	// keys, but paints the modifier keys in that identical azure instead of
	// violet, and carries the blue through to the action subtitle and streak
	// badge. Requested for people who like the neon look but read a purple Ctrl
	// as a second, unrelated accent colour.
	ThemeAzurePulse = ThemeConfig{
		Name:                  "Azure Pulse",
		DisplayName:           "深蓝脉冲 (Azure Pulse)",
		Description:           "统一冰蓝键帽：修饰键与字母键同为深空蓝荧光，深蓝底板一体化配色。",
		DefaultStyle:          "neon",
		TravelDepth:           0.0,
		HeldPhysics:           HeldPulse,
		AcousticProfile:       "droplet",
		SpringBounce:          0.28,
		SpringDurationScale:   0.90,
		WindowBackground:      rgba(8, 12, 22, 235),
		WindowCorner:          16,
		KeyCapBackground:      rgba(14, 20, 40, 245),
		KeyCapBorder:          rgba(6, 182, 212, 225),
		KeyCapHeldBackground:  rgba(14, 165, 233, 205),
		KeyCapHeldBorder:      rgba(125, 211, 252, 255),
		KeyCapTextColor:       rgba(34, 211, 238, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    10,
		ModCapBackground:      rgba(12, 18, 38, 242),
		ModCapBorder:          rgba(6, 182, 212, 225),
		ModCapTextColor:       rgba(34, 211, 238, 255),
		SubtitleBackground:    rgba(10, 18, 34, 242),
		SubtitleTextColor:     rgba(56, 189, 248, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(14, 165, 233, 245),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(56, 189, 248, 205),
	}

	ThemeClassicMechanical = ThemeConfig{
		Name:                  "Classic Mechanical",
		DisplayName:           "复古机械 (Classic Mechanical)",
		Description:           "PBT 磨砂双色注塑键帽与微凹指腹触感仿真。",
		DefaultStyle:          "pbt",
		TravelDepth:           3.0,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "thock",
		SpringBounce:          0.22,
		SpringDurationScale:   1.05,
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
		TravelDepth:           3.8,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "typewriter",
		SpringBounce:          0.26,
		SpringDurationScale:   1.10,
		WindowBackground:      rgba(235, 231, 224, 240),
		WindowCorner:          12,
		KeyCapBackground:      rgba(245, 242, 237, 255),
		KeyCapBorder:          rgba(202, 196, 184, 255),
		KeyCapHeldBackground:  rgba(222, 217, 208, 255),
		KeyCapHeldBorder:      rgba(175, 168, 156, 255),
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
		TravelDepth:           1.6,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "paper",
		SpringBounce:          0.16,
		SpringDurationScale:   0.90,
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
		TravelDepth:           0.0,
		HeldPhysics:           HeldPulse,
		AcousticProfile:       "pebble",
		SpringBounce:          0.14,
		SpringDurationScale:   0.85,
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
		TravelDepth:           2.2,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "felt",
		SpringBounce:          0.20,
		SpringDurationScale:   1.0,
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
		TravelDepth:           2.8,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "chip",
		SpringBounce:          0.24,
		SpringDurationScale:   1.0,
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
		TravelDepth:           2.0,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "thock",
		SpringBounce:          0.20,
		SpringDurationScale:   0.95,
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
		TravelDepth:           3.2,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "felt",
		SpringBounce:          0.22,
		SpringDurationScale:   1.05,
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

	ThemeSonyWalkman = ThemeConfig{
		Name:                  "Sony Walkman 1979",
		DisplayName:           "索尼经典 (Sony Walkman 1979)",
		Description:           "索尼首款 Walkman 金属银蓝与磁带仓亚克力，经典橙色 HOT LINE 标识。",
		DefaultStyle:          "laptop",
		TravelDepth:           2.4,
		HeldPhysics:           HeldSnappy,
		AcousticProfile:       "hinoki",
		SpringBounce:          0.18,
		SpringDurationScale:   0.90,
		WindowBackground:      rgba(20, 35, 60, 235),
		WindowCorner:          14,
		KeyCapBackground:      rgba(215, 222, 232, 255),
		KeyCapBorder:          rgba(148, 163, 184, 255),
		KeyCapHeldBackground:  rgba(255, 107, 0, 255),
		KeyCapHeldBorder:      rgba(234, 88, 12, 255),
		KeyCapTextColor:       rgba(15, 23, 42, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    5,
		ModCapBackground:      rgba(50, 75, 110, 245),
		ModCapBorder:          rgba(80, 115, 160, 255),
		ModCapTextColor:       rgba(224, 231, 255, 255),
		SubtitleBackground:    rgba(15, 23, 42, 230),
		SubtitleTextColor:     rgba(255, 107, 0, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(255, 107, 0, 255),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(255, 107, 0, 190),
	}

	ThemeTeenageOP1 = ThemeConfig{
		Name:                  "Teenage OP-1",
		DisplayName:           "合成器先锋 (Teenage OP-1)",
		Description:           "瑞典合成器工业美学：平头矮键帽、冷灰白底座与高饱和乐高三原色。",
		DefaultStyle:          "minimal",
		TravelDepth:           1.0,
		HeldPhysics:           HeldSnappy,
		AcousticProfile:       "chip",
		SpringBounce:          0.14,
		SpringDurationScale:   0.80,
		WindowBackground:      rgba(228, 231, 236, 240),
		WindowCorner:          12,
		KeyCapBackground:      rgba(255, 255, 255, 255),
		KeyCapBorder:          rgba(209, 213, 219, 255),
		KeyCapHeldBackground:  rgba(16, 185, 129, 255),
		KeyCapHeldBorder:      rgba(5, 150, 105, 255),
		KeyCapTextColor:       rgba(31, 41, 55, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    4,
		ModCapBackground:      rgba(243, 244, 246, 255),
		ModCapBorder:          rgba(14, 165, 233, 220),
		ModCapTextColor:       rgba(2, 132, 199, 255),
		SubtitleBackground:    rgba(31, 41, 55, 230),
		SubtitleTextColor:     rgba(16, 185, 129, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(244, 63, 94, 255),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(16, 185, 129, 200),
	}

	ThemeAmberCathode = ThemeConfig{
		Name:                  "Amber Cathode 1977",
		DisplayName:           "琥珀终端 (Amber Cathode 1977)",
		Description:           "DEC VT100 阴极射线管曲面屏，深黑背景与琥珀金单色荧光磷光体。",
		DefaultStyle:          "retro",
		TravelDepth:           3.4,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "thock",
		SpringBounce:          0.22,
		SpringDurationScale:   1.05,
		WindowBackground:      rgba(12, 10, 4, 245),
		WindowCorner:          12,
		KeyCapBackground:      rgba(30, 24, 10, 250),
		KeyCapBorder:          rgba(245, 158, 11, 200),
		KeyCapHeldBackground:  rgba(245, 158, 11, 230),
		KeyCapHeldBorder:      rgba(251, 191, 36, 255),
		KeyCapTextColor:       rgba(251, 191, 36, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    8,
		ModCapBackground:      rgba(20, 16, 6, 245),
		ModCapBorder:          rgba(217, 119, 6, 180),
		ModCapTextColor:       rgba(245, 158, 11, 255),
		SubtitleBackground:    rgba(20, 16, 6, 240),
		SubtitleTextColor:     rgba(251, 191, 36, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(217, 119, 6, 240),
		StreakBadgeTextColor:  rgba(0, 0, 0, 255),
		MouseRippleColor:      rgba(245, 158, 11, 190),
	}

	ThemeSwissTypography = ThemeConfig{
		Name:                  "Swiss Typography 1957",
		DisplayName:           "苏黎世平面 (Swiss Typography 1957)",
		Description:           "瑞士国际主义平面风格：纯白底板、纯黑粗体无衬线字体与瑞士红点。",
		DefaultStyle:          "minimal",
		TravelDepth:           1.4,
		HeldPhysics:           HeldSnappy,
		AcousticProfile:       "paper",
		SpringBounce:          0.12,
		SpringDurationScale:   0.85,
		WindowBackground:      rgba(250, 250, 250, 245),
		WindowCorner:          12,
		KeyCapBackground:      rgba(255, 255, 255, 255),
		KeyCapBorder:          rgba(0, 0, 0, 255),
		KeyCapHeldBackground:  rgba(220, 38, 38, 255),
		KeyCapHeldBorder:      rgba(185, 28, 28, 255),
		KeyCapTextColor:       rgba(0, 0, 0, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    4,
		ModCapBackground:      rgba(240, 240, 240, 255),
		ModCapBorder:          rgba(0, 0, 0, 220),
		ModCapTextColor:       rgba(0, 0, 0, 255),
		SubtitleBackground:    rgba(0, 0, 0, 240),
		SubtitleTextColor:     rgba(255, 255, 255, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(220, 38, 38, 255),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(220, 38, 38, 200),
	}

	ThemeKyotoKintsugi = ThemeConfig{
		Name:                  "Kyoto Kintsugi",
		DisplayName:           "京都金缮 (Kyoto Kintsugi)",
		Description:           "粗粝黑陶砂质感与流金漆纹，以侘寂物哀哲学治愈屏幕的冰冷。",
		DefaultStyle:          "pbt",
		TravelDepth:           2.6,
		HeldPhysics:           HeldSink,
		AcousticProfile:       "hinoki",
		SpringBounce:          0.18,
		SpringDurationScale:   1.0,
		WindowBackground:      rgba(26, 24, 22, 240),
		WindowCorner:          16,
		KeyCapBackground:      rgba(42, 39, 36, 250),
		KeyCapBorder:          rgba(217, 160, 43, 200),
		KeyCapHeldBackground:  rgba(180, 130, 25, 240),
		KeyCapHeldBorder:      rgba(245, 197, 74, 255),
		KeyCapTextColor:       rgba(245, 238, 220, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    6,
		ModCapBackground:      rgba(32, 29, 27, 245),
		ModCapBorder:          rgba(180, 130, 25, 160),
		ModCapTextColor:       rgba(217, 160, 43, 255),
		SubtitleBackground:    rgba(26, 24, 22, 240),
		SubtitleTextColor:     rgba(245, 197, 74, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(217, 160, 43, 240),
		StreakBadgeTextColor:  rgba(20, 18, 16, 255),
		MouseRippleColor:      rgba(217, 160, 43, 190),
	}

	ThemeBlueprintCyanotype = ThemeConfig{
		Name:                  "Blueprint Cyanotype",
		DisplayName:           "工程蓝晒 (Blueprint Cyanotype)",
		Description:           "普鲁士蓝底色与高精度细白工程测绘网格，理性严谨的工程师手稿。",
		DefaultStyle:          "outline",
		TravelDepth:           1.6,
		HeldPhysics:           HeldSnappy,
		AcousticProfile:       "paper",
		SpringBounce:          0.15,
		SpringDurationScale:   0.88,
		WindowBackground:      rgba(10, 30, 75, 240),
		WindowCorner:          14,
		KeyCapBackground:      rgba(15, 45, 105, 245),
		KeyCapBorder:          rgba(147, 197, 253, 220),
		KeyCapHeldBackground:  rgba(37, 99, 235, 240),
		KeyCapHeldBorder:      rgba(219, 234, 254, 255),
		KeyCapTextColor:       rgba(239, 246, 255, 255),
		KeyCapFontSize:        20,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    4,
		ModCapBackground:      rgba(10, 35, 85, 240),
		ModCapBorder:          rgba(96, 165, 250, 180),
		ModCapTextColor:       rgba(191, 219, 254, 255),
		SubtitleBackground:    rgba(10, 30, 75, 240),
		SubtitleTextColor:     rgba(147, 197, 253, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(37, 99, 235, 240),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(147, 197, 253, 190),
	}

	ThemeLeicaM = ThemeConfig{
		Name:                  "Leica M",
		DisplayName:           "银黑徕卡 (Leica M)",
		Description:           "银色镀铬顶盖与黑色皮革机身，快门落下的一刻被永远留住。",
		DefaultStyle:          "laptop",
		TravelDepth:           1.8,
		HeldPhysics:           HeldSnappy,
		AcousticProfile:       "shutter",
		SpringBounce:          0.14,
		SpringDurationScale:   0.85,
		WindowBackground:      rgba(18, 18, 18, 235),
		WindowCorner:          14,
		KeyCapBackground:      rgba(196, 199, 201, 250),
		KeyCapBorder:          rgba(238, 240, 242, 255),
		KeyCapHeldBackground:  rgba(120, 124, 128, 255),
		KeyCapHeldBorder:      rgba(255, 255, 255, 255),
		KeyCapTextColor:       rgba(24, 24, 24, 255),
		KeyCapFontSize:        19,
		KeyCapRadius:          12.5,
		KeyCapShadowSpread:    5,
		ModCapBackground:      rgba(34, 32, 30, 255),
		ModCapBorder:          rgba(70, 66, 62, 255),
		ModCapTextColor:       rgba(230, 226, 220, 255),
		SubtitleBackground:    rgba(20, 20, 20, 235),
		SubtitleTextColor:     rgba(229, 231, 235, 255),
		SubtitleFontSize:      12,
		StreakBadgeBackground: rgba(200, 16, 46, 235),
		StreakBadgeTextColor:  rgba(255, 255, 255, 255),
		MouseRippleColor:      rgba(200, 16, 46, 170),
	}
)

var AllThemes = []ThemeConfig{
	ThemeLiquidGlass,
	ThemeMacBookScissor,
	ThemeDieterRams,
	ThemeCyberpunkNeon,
	ThemeAzurePulse,
	ThemeClassicMechanical,
	ThemeIBMModelM,
	ThemeMinimalChalk,
	ThemeObsidianOLED,
	ThemeTokyoBlossom,
	ThemeEmeraldMatrix,
	ThemeSunsetHorizon,
	ThemeRoyalVelvet,
	ThemeSonyWalkman,
	ThemeTeenageOP1,
	ThemeAmberCathode,
	ThemeSwissTypography,
	ThemeKyotoKintsugi,
	ThemeBlueprintCyanotype,
	ThemeLeicaM,
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
	case "sony walkman 1979":
		return "laptop"
	case "teenage op-1":
		return "minimal"
	case "amber cathode 1977":
		return "retro"
	case "swiss typography 1957":
		return "minimal"
	case "kyoto kintsugi":
		return "pbt"
	case "blueprint cyanotype":
		return "outline"
	case "leica m":
		return "laptop"
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

// MouseThemeConfig holds the visual color palette and aesthetics for mouse cursor overlay and indicators.
type MouseThemeConfig struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	LeftColor   string `json:"leftColor"`  // Left click button color
	RightColor  string `json:"rightColor"` // Right click button color
	WheelColor  string `json:"wheelColor"` // Middle click & scroll wheel color
	RingColor   string `json:"ringColor"`  // Contact puck & cursor ripple color
	BadgeColor  string `json:"badgeColor"` // Multi-click pill badge color
}

// AllMouseThemes provides vibrant preset mouse themes.
var AllMouseThemes = []MouseThemeConfig{
	{
		ID:          "neon_green",
		Name:        "Keyviz Neon Green",
		DisplayName: "原版荧光绿 (Neon Green)",
		LeftColor:   "#00FF6A",
		RightColor:  "#00FF6A",
		WheelColor:  "#00FF6A",
		RingColor:   "#00FF6A",
		BadgeColor:  "#ef4444",
	},
	{
		ID:          "aurora_bicolor",
		Name:        "Aurora Dual Tone",
		DisplayName: "极光双色 (Emerald & Ice)",
		LeftColor:   "#00FF6A",
		RightColor:  "#38bdf8",
		WheelColor:  "#a3e635",
		RingColor:   "#00FF6A",
		BadgeColor:  "#ef4444",
	},
	{
		ID:          "cyber_cyan",
		Name:        "Cyber Cyan",
		DisplayName: "赛博霓虹 (Cyber Neon)",
		LeftColor:   "#00f0ff",
		RightColor:  "#ff007f",
		WheelColor:  "#facc15",
		RingColor:   "#00f0ff",
		BadgeColor:  "#ff007f",
	},
	{
		ID:          "obsidian_crimson",
		Name:        "Obsidian Crimson",
		DisplayName: "黑曜炽红 (Obsidian Crimson)",
		LeftColor:   "#ff3b30",
		RightColor:  "#ff9500",
		WheelColor:  "#ffd60a",
		RingColor:   "#ff3b30",
		BadgeColor:  "#ff3b30",
	},
	{
		ID:          "deep_blue",
		Name:        "Deep Space Blue",
		DisplayName: "深空幻蓝 (Deep Space Blue)",
		LeftColor:   "#38bdf8",
		RightColor:  "#818cf8",
		WheelColor:  "#60a5fa",
		RingColor:   "#38bdf8",
		BadgeColor:  "#38bdf8",
	},
	{
		ID:          "minimal_white",
		Name:        "Minimal Chalk White",
		DisplayName: "纯白极简 (Chalk White)",
		LeftColor:   "#ffffff",
		RightColor:  "#cbd5e1",
		WheelColor:  "#f1f5f9",
		RingColor:   "#ffffff",
		BadgeColor:  "#ef4444",
	},
}

// MouseThemeByID returns the configured mouse theme or the default Neon Green.
func MouseThemeByID(id string) MouseThemeConfig {
	for _, m := range AllMouseThemes {
		if strings.EqualFold(m.ID, id) {
			return m
		}
	}
	return AllMouseThemes[0] // default: neon_green
}
