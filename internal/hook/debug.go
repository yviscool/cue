package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Temporary diagnostics: set CUE_DEBUG_INPUT=1 to append raw hook events to
// CUE_DEBUG_LOG (default: %TEMP%/cue-input.log). Remove after debugging.

var (
	debugOnce sync.Once
	debugPath string
)

// DebugLogf appends one line to the diagnostics log when enabled.
func DebugLogf(format string, args ...any) {
	debugOnce.Do(func() {
		if os.Getenv("CUE_DEBUG_INPUT") == "" {
			return
		}
		debugPath = os.Getenv("CUE_DEBUG_LOG")
		if debugPath == "" {
			debugPath = filepath.Join(os.TempDir(), "cue-input.log")
		}
	})
	if debugPath == "" {
		return
	}
	f, err := os.OpenFile(debugPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, time.Now().Format("15:04:05.000")+" "+format+"\n", args...)
}

func debugLogf(format string, args ...any) { DebugLogf(format, args...) }
