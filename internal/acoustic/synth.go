package acoustic

import (
	"encoding/binary"
	"math"
	"sync"
)

// SoundType identifies the procedural sound character.
type SoundType int

const (
	// SoundThock is a low-frequency, deep, warm mechanical switch bottom-out (~180Hz-220Hz, ~35ms).
	SoundThock SoundType = iota
	// SoundClick is a crisp, bright tactile click transient (~1200Hz-2400Hz, ~15ms).
	SoundClick
)

func (s SoundType) String() string {
	switch s {
	case SoundThock:
		return "Thock"
	case SoundClick:
		return "Click"
	default:
		return "Unknown"
	}
}

// SampleRate is the standard CD-quality audio sample rate in Hz.
const SampleRate = 44100

// Player manages audio playback and maintains pre-rendered sound buffers in memory.
type Player struct {
	mu       sync.RWMutex
	enabled  bool
	thockWAV []byte
	clickWAV []byte
}

// NewPlayer initializes the acoustic player with synthesized mechanical sound buffers.
func NewPlayer() *Player {
	thock := SynthesizeWAV(SoundThock)
	click := SynthesizeWAV(SoundClick)

	p := &Player{
		enabled:  true,
		thockWAV: thock,
		clickWAV: click,
	}
	return p
}

// SetEnabled toggles audio playback on or off.
func (p *Player) SetEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = enabled
}

// IsEnabled reports whether acoustic feedback is currently enabled.
func (p *Player) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// PlayThock plays a low-frequency, deep mechanical thock.
func (p *Player) PlayThock() {
	if !p.IsEnabled() {
		return
	}
	_ = playWAVMemory(p.thockWAV)
}

// PlayClick plays a crisp, bright tactile switch click.
func (p *Player) PlayClick() {
	if !p.IsEnabled() {
		return
	}
	_ = playWAVMemory(p.clickWAV)
}

// PlaySound plays the specified sound type.
func (p *Player) PlaySound(st SoundType) {
	switch st {
	case SoundThock:
		p.PlayThock()
	case SoundClick:
		p.PlayClick()
	}
}

// GetWAV returns the raw WAV byte slice for a given sound type.
func (p *Player) GetWAV(st SoundType) []byte {
	switch st {
	case SoundThock:
		return p.thockWAV
	case SoundClick:
		return p.clickWAV
	default:
		return nil
	}
}

// Stop terminates any currently active asynchronous sound playback.
func (p *Player) Stop() {
	stopSoundMemory()
}

// SynthesizeWAV procedurally synthesizes in-memory RIFF WAV bytes for the given SoundType.
func SynthesizeWAV(st SoundType) []byte {
	samples := SynthesizeSamples(st)
	return EncodeWAV(samples, SampleRate)
}

// SynthesizeSamples generates raw 16-bit PCM samples for a sound type.
func SynthesizeSamples(st SoundType) []int16 {
	switch st {
	case SoundThock:
		return synthesizeThock()
	case SoundClick:
		return synthesizeClick()
	default:
		return synthesizeThock()
	}
}

// synthesizeThock generates a low-frequency, deep, warm mechanical thock.
// Frequency sweeps from ~220Hz down to ~180Hz over ~35ms with exponential decay
// and a subtle plastic contact transient.
func synthesizeThock() []int16 {
	durationSec := 0.035 // 35 ms
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	// Deterministic pseudo-random number generator for plastic impact transient
	var prngState uint32 = 0x1337BEEF
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phase := 0.0
	attackSec := 0.0015 // 1.5ms fast attack
	decayTau := 0.010   // 10ms exponential decay time constant

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)

		// Frequency glide: starts at 220Hz dropping towards 180Hz
		freq := 180.0 + 40.0*math.Exp(-t/0.008)
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		// Harmonic body: fundamental + warm 2nd harmonic
		body := math.Sin(phase) + 0.30*math.Sin(2.0*phase)

		// Plastic keycap/switch housing contact impact transient
		transientNoise := randFloat() * math.Exp(-t/0.002) * 0.25

		// Amplitude envelope: fast attack then exponential decay
		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		// Smooth fadeout at very end of sample buffer to avoid DC pop
		remaining := numSamples - 1 - i
		if remaining < 30 {
			env *= float64(remaining) / 30.0
		}

		// Soft clipping / saturation
		val := (body + transientNoise) * env * 24000.0
		if val > 32767.0 {
			val = 32767.0
		} else if val < -32768.0 {
			val = -32768.0
		}

		samples[i] = int16(val)
	}

	return samples
}

// synthesizeClick generates a crisp, bright tactile click transient.
// Transient click sweeps from ~2400Hz down to ~1200Hz with rapid decay ~15ms.
func synthesizeClick() []int16 {
	durationSec := 0.015 // 15 ms
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	var prngState uint32 = 0x98765432
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phase := 0.0
	attackSec := 0.0006 // 0.6ms ultra-fast snap
	decayTau := 0.0032  // 3.2ms rapid decay

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)

		// High frequency sweep: 2400Hz -> 1200Hz
		freq := 1200.0 + 1200.0*math.Exp(-t/0.003)
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		// Crisp tactile tone
		tone := 0.65*math.Sin(phase) + 0.35*math.Sin(1.5*phase)

		// Tactile snap noise burst during the initial 1.5ms
		snap := randFloat() * math.Exp(-t/0.0012) * 0.40

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		// Smooth tail fade
		remaining := numSamples - 1 - i
		if remaining < 20 {
			env *= float64(remaining) / 20.0
		}

		val := (tone + snap) * env * 22000.0
		if val > 32767.0 {
			val = 32767.0
		} else if val < -32768.0 {
			val = -32768.0
		}

		samples[i] = int16(val)
	}

	return samples
}

// EncodeWAV serializes 16-bit mono PCM samples into a standard 44-byte RIFF WAV byte array.
func EncodeWAV(samples []int16, sampleRate int) []byte {
	numSamples := len(samples)
	dataSize := numSamples * 2
	totalSize := 44 + dataSize
	buf := make([]byte, totalSize)

	// RIFF Chunk Descriptor
	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], uint32(36+dataSize))
	copy(buf[8:12], "WAVE")

	// "fmt " sub-chunk
	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16) // Subchunk1Size for PCM = 16
	binary.LittleEndian.PutUint16(buf[20:22], 1)  // AudioFormat 1 = PCM
	binary.LittleEndian.PutUint16(buf[22:24], 1)  // NumChannels 1 = Mono
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	byteRate := uint32(sampleRate * 1 * 2)
	binary.LittleEndian.PutUint32(buf[28:32], byteRate)
	binary.LittleEndian.PutUint16(buf[32:34], 2)  // BlockAlign = NumChannels * BitsPerSample / 8
	binary.LittleEndian.PutUint16(buf[34:36], 16) // BitsPerSample = 16

	// "data" sub-chunk
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))

	// PCM Sample payload
	offset := 44
	for _, s := range samples {
		binary.LittleEndian.PutUint16(buf[offset:offset+2], uint16(s))
		offset += 2
	}

	return buf
}
