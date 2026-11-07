//go:build windows

package acoustic

import (
	"encoding/binary"
	"testing"
)

// newTestMixer builds a stream engine that never touched an audio device. The
// mixer only reads e.volume and e.voices, so gain behaviour is testable without
// waveOut being available.
func newTestMixer(volume float64, samples []int16) (*windowsStreamEngine, []byte) {
	e := &windowsStreamEngine{volume: volume}
	if len(samples) > 0 {
		e.voices = append(e.voices, activeVoice{samples: samples, pos: 0})
	}
	return e, make([]byte, bufSizeBytes)
}

func frameAt(buf []byte, i int) (int16, int16) {
	off := i * 4
	return int16(binary.LittleEndian.Uint16(buf[off:])), int16(binary.LittleEndian.Uint16(buf[off+2:]))
}

// TestMixBufferAppliesMasterGain is the core check for the volume slider: the
// synthesized samples are pre-rendered, so the only place a volume change can
// take effect is the mixer.
func TestMixBufferAppliesMasterGain(t *testing.T) {
	cases := []struct {
		name   string
		volume float64
		wantL  int16
		wantR  int16
	}{
		{"unity gain is unchanged", 1.0, 1000, 2000},
		{"half volume halves the signal", 0.5, 500, 1000},
		{"quarter volume", 0.25, 250, 500},
		{"near silence", 0.01, 10, 20},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e, buf := newTestMixer(tc.volume, []int16{1000, 2000})
			e.mixBuffer(buf)
			l, r := frameAt(buf, 0)
			if l != tc.wantL || r != tc.wantR {
				t.Errorf("got (%d, %d), want (%d, %d)", l, r, tc.wantL, tc.wantR)
			}
		})
	}
}

// TestMixBufferMuteDropsVoices checks that muting produces silence and discards
// pending voices, so turning the volume back up cannot replay a burst of keys
// that were pressed while muted.
func TestMixBufferMuteDropsVoices(t *testing.T) {
	e, buf := newTestMixer(0, []int16{1000, 2000, 3000, 4000})
	e.mixBuffer(buf)

	for i := 0; i < framesPerBuf; i++ {
		if l, r := frameAt(buf, i); l != 0 || r != 0 {
			t.Fatalf("frame %d is not silent: (%d, %d)", i, l, r)
		}
	}
	if len(e.voices) != 0 {
		t.Errorf("expected muted voices to be dropped, %d remain", len(e.voices))
	}

	// Raising the volume afterwards must not resurrect them.
	e.SetVolume(1.0)
	e.mixBuffer(buf)
	for i := 0; i < framesPerBuf; i++ {
		if l, r := frameAt(buf, i); l != 0 || r != 0 {
			t.Fatalf("frame %d resurrected audio after unmute: (%d, %d)", i, l, r)
		}
	}
}

func TestMixBufferIdleIsSilent(t *testing.T) {
	e, buf := newTestMixer(1.0, nil)
	for i := range buf {
		buf[i] = 0xff
	}
	e.mixBuffer(buf)
	for i := range buf {
		if buf[i] != 0 {
			t.Fatalf("idle buffer byte %d = %#x, want 0", i, buf[i])
		}
	}
}

func TestSetVolumeClamps(t *testing.T) {
	e, _ := newTestMixer(1.0, nil)
	for _, tc := range []struct {
		in   float64
		want float64
	}{
		{-2, 0},
		{0, 0},
		{0.42, 0.42},
		{1, 1},
		{7, 1},
	} {
		e.SetVolume(tc.in)
		if e.volume != tc.want {
			t.Errorf("SetVolume(%v) left %v, want %v", tc.in, e.volume, tc.want)
		}
	}
}

// TestPlayerVolumeRoundTrip covers the public surface used by the settings UI
// and the config apply path.
func TestPlayerVolumeRoundTrip(t *testing.T) {
	p := NewPlayer()
	defer p.Close()

	if got := p.GetVolume(); got != 1.0 {
		t.Errorf("new player should be at unity gain, got %v", got)
	}
	p.SetVolume(0.4)
	if got := p.GetVolume(); got != 0.4 {
		t.Errorf("GetVolume() = %v, want 0.4", got)
	}
	p.SetVolume(-1)
	if got := p.GetVolume(); got != 0 {
		t.Errorf("negative volume should clamp to 0, got %v", got)
	}
	p.SetVolume(9)
	if got := p.GetVolume(); got != 1 {
		t.Errorf("oversized volume should clamp to 1, got %v", got)
	}
}
