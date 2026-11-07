//go:build !windows

package context

import (
	"runtime"
	"time"
)

func newPlatformProvider(interval time.Duration) Provider {
	return NewStatic(AppContext{Platform: runtime.GOOS})
}
