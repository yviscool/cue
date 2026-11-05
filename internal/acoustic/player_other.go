//go:build !windows

package acoustic

// playWAVMemory is a no-op fallback on non-Windows platforms.
func playWAVMemory(data []byte) error {
	return nil
}

// stopSoundMemory is a no-op fallback on non-Windows platforms.
func stopSoundMemory() {}
