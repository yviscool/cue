package i18n

import (
	"testing"
)

func TestI18nTranslations(t *testing.T) {
	mgr := Instance()

	// 1. Test Chinese (Default)
	mgr.SetLanguage(LangZhCN)
	if mgr.GetLanguage() != LangZhCN {
		t.Fatalf("expected LangZhCN, got %v", mgr.GetLanguage())
	}
	if text := mgr.T("tab_general"); text != "通用设置" {
		t.Fatalf("unexpected zh-CN translation: %q", text)
	}

	// 2. Test English
	mgr.SetLanguage(LangEnUS)
	if text := mgr.T("tab_general"); text != "General" {
		t.Fatalf("unexpected en-US translation: %q", text)
	}
	if text := mgr.T("reset_defaults"); text != "Reset Defaults" {
		t.Fatalf("unexpected en-US translation: %q", text)
	}

	// 3. Test Traditional Chinese
	mgr.SetLanguage(LangZhTW)
	if text := mgr.T("tab_general"); text != "一般設定" {
		t.Fatalf("unexpected zh-TW translation: %q", text)
	}

	// 4. Test Japanese
	mgr.SetLanguage(LangJaJP)
	if text := mgr.T("tab_general"); text != "全般設定" {
		t.Fatalf("unexpected ja-JP translation: %q", text)
	}

	// 5. Fallback test for missing key
	if text := mgr.T("non_existent_key_xyz"); text != "non_existent_key_xyz" {
		t.Fatalf("expected fallback to key itself, got %q", text)
	}

	// 6. Test Action subtitles & Mouse across all locales
	mgr.SetLanguage(LangZhCN)
	if text := mgr.T("Select All"); text != "全选" {
		t.Fatalf("expected '全选', got %q", text)
	}
	if text := mgr.T("Left Click"); text != "鼠标左键" {
		t.Fatalf("expected '鼠标左键', got %q", text)
	}
	if text := mgr.T("Close Tab"); text != "关闭标签页" {
		t.Fatalf("expected '关闭标签页', got %q", text)
	}

	mgr.SetLanguage(LangZhTW)
	if text := mgr.T("Select All"); text != "全選" {
		t.Fatalf("expected '全選', got %q", text)
	}
	if text := mgr.T("Cut"); text != "剪下" {
		t.Fatalf("expected '剪下', got %q", text)
	}

	mgr.SetLanguage(LangJaJP)
	if text := mgr.T("Select All"); text != "すべて選択" {
		t.Fatalf("expected 'すべて選択', got %q", text)
	}
	if text := mgr.T("Left Click"); text != "左クリック" {
		t.Fatalf("expected '左クリック', got %q", text)
	}

	mgr.SetLanguage(LangEnUS)
	if text := mgr.T("Select All"); text != "Select All" {
		t.Fatalf("expected 'Select All', got %q", text)
	}
}

func TestI18nFormatting(t *testing.T) {
	mgr := Instance()
	mgr.SetLanguage(LangZhCN)
	formatted := TF("stay_duration", 3.5)
	if formatted != "屏幕留存时长: 3.5 秒" {
		t.Fatalf("unexpected formatted string: %q", formatted)
	}

	mgr.SetLanguage(LangEnUS)
	formattedEn := TF("stay_duration", 3.5)
	if formattedEn != "Stay on screen: 3.5 seconds" {
		t.Fatalf("unexpected formatted en string: %q", formattedEn)
	}
}
