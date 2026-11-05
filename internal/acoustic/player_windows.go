//go:build windows

package acoustic

import (
	"syscall"
	"unsafe"
)

var (
	winmmDLL      = syscall.NewLazyDLL("winmm.dll")
	procPlaySound = winmmDLL.NewProc("PlaySoundW")
)

const (
	sndSync      = 0x0000
	sndAsync     = 0x0001
	sndNodefault = 0x0002
	sndMemory    = 0x0004
)

// playWAVMemory plays in-memory WAV audio using Windows winmm.dll PlaySoundW.
// Flags used: SND_MEMORY | SND_ASYNC | SND_NODEFAULT (0x0004 | 0x0001 | 0x0002).
// If audio fails or winmm is missing, it silently falls back without panic.
func playWAVMemory(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	defer func() {
		_ = recover() // Catch unexpected low-level panic
	}()

	if err := procPlaySound.Find(); err != nil {
		return err
	}

	flags := uintptr(sndMemory | sndAsync | sndNodefault)
	r1, _, err := procPlaySound.Call(
		uintptr(unsafe.Pointer(&data[0])),
		0,
		flags,
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// stopSoundMemory stops any active asynchronous waveform playback.
func stopSoundMemory() {
	defer func() {
		_ = recover()
	}()
	if err := procPlaySound.Find(); err == nil {
		_, _, _ = procPlaySound.Call(0, 0, 0)
	}
}
