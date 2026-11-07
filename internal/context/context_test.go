package context

import "testing"

func TestNormalizeExe(t *testing.T) {
	cases := map[string]string{
		`C:\Program Files (x86)\Dev-Cpp\devcpp.exe`: "devcpp.exe",
		"C:/Tools/Code.EXE":                         "code.exe",
		"  pythonw.exe  ":                           "pythonw.exe",
		"":                                          "",
	}
	for input, want := range cases {
		if got := NormalizeExe(input); got != want {
			t.Errorf("NormalizeExe(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestAppContextDisplayAndProfile(t *testing.T) {
	ctx := AppContext{Exe: "devcpp.exe", Title: "main.cpp - Dev-C++"}
	if got := ctx.DisplayName(); got != "devcpp" {
		t.Errorf("expected exe-derived name, got %q", got)
	}

	ctx = ctx.WithProfile("devcpp", "Dev-C++", CategoryEditor)
	if got := ctx.DisplayName(); got != "Dev-C++" {
		t.Errorf("expected profile name, got %q", got)
	}
	if ctx.Category != CategoryEditor {
		t.Errorf("expected editor category, got %q", ctx.Category)
	}
	if CategoryEditor.DisplayName() != "Editor" {
		t.Errorf("unexpected category display name %q", CategoryEditor.DisplayName())
	}
}

func TestStaticProviderPublishesOnChangeOnly(t *testing.T) {
	p := NewStatic(AppContext{Exe: "code.exe", Title: "main.go"})
	defer p.Stop()

	var events []string
	unsub := p.Subscribe(func(ctx AppContext) {
		events = append(events, ctx.Title)
	})
	defer unsub()

	// Same effective context: Key() unchanged, so no notification.
	p.Set(AppContext{Exe: "code.exe", Title: "main.go"})
	if len(events) != 0 {
		t.Fatalf("expected no event for identical context, got %v", events)
	}

	p.Set(AppContext{Exe: "code.exe", Title: "other.go"})
	if len(events) != 1 || events[0] != "other.go" {
		t.Fatalf("expected one change event, got %v", events)
	}

	// Seq must be monotonic and non-zero for provider-backed contexts.
	first := p.Current().Seq
	p.Set(AppContext{Exe: "devenv.exe", Title: "solution.sln"})
	if p.Current().Seq <= first {
		t.Fatalf("expected increasing Seq, got %d then %d", first, p.Current().Seq)
	}

	unsub()
	p.Set(AppContext{Exe: "chrome.exe", Title: "example"})
	if len(events) != 2 {
		t.Fatalf("unsubscribed callback still fired: %v", events)
	}
}
