//go:build !windows

package acoustic

func newAudioEngine() audioEngine {
	return &stubAudioEngine{}
}

// playWAVMemory is a no-op fallback on non-Windows platforms.
func playWAVMemory(data []byte) error {
	return nil
}

// stopSoundMemory is a no-op fallback on non-Windows platforms.
func stopSoundMemory() {}
