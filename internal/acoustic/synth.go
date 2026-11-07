package acoustic

import (
	"encoding/binary"
	"math"
	"strings"
	"sync"
	"sync/atomic"
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

// SpatialZone discretizes the keyboard width into five acoustic stages.
type SpatialZone int

const (
	ZoneFarLeft     SpatialZone = iota // pan ≈ -0.65 (Esc, Tab, Caps, ShiftL, A, Q, Z)
	ZoneMidLeft                        // pan ≈ -0.30 (W, S, E, D, R, F, 2, 3, 4)
	ZoneCenter                         // pan ≈ 0.00 (T, G, B, Y, H, N, Space)
	ZoneMidRight                       // pan ≈ +0.30 (U, J, I, K, O, L, M)
	ZoneFarRight                       // pan ≈ +0.65 (P, Enter, Backspace, Del, ShiftR)
	numSpatialZones = 5
)

const numVariations = 3

var zonePanValues = [numSpatialZones]float64{
	-0.65, // ZoneFarLeft
	-0.30, // ZoneMidLeft
	0.00,  // ZoneCenter
	0.30,  // ZoneMidRight
	0.65,  // ZoneFarRight
}

var variationPitchFactors = [numVariations]float64{
	1.000, // Variation 0: nominal baseline tuning
	1.035, // Variation 1: snappy crisp transient (+3.5%)
	0.965, // Variation 2: warm deep body (-3.5%)
}

type audioEngine interface {
	PlaySamples(samples []int16)
	Stop()
	Close()
}

type stubAudioEngine struct{}

func (s *stubAudioEngine) PlaySamples([]int16) {}
func (s *stubAudioEngine) Stop()               {}
func (s *stubAudioEngine) Close()              {}

// SoundProfile identifies the acoustic character archetype.
type SoundProfile string

const (
	ProfileAuto       SoundProfile = "auto"
	ProfileThock      SoundProfile = "thock"
	ProfileClick      SoundProfile = "click"
	ProfileScissor    SoundProfile = "scissor"
	ProfileGlass      SoundProfile = "glass"
	ProfileTypewriter SoundProfile = "typewriter"
)

const numProfiles = 5

var soundProfiles = [numProfiles]SoundProfile{
	ProfileThock,
	ProfileClick,
	ProfileScissor,
	ProfileGlass,
	ProfileTypewriter,
}

func profileToIndex(p SoundProfile) int {
	switch p {
	case ProfileThock:
		return 0
	case ProfileClick:
		return 1
	case ProfileScissor:
		return 2
	case ProfileGlass:
		return 3
	case ProfileTypewriter:
		return 4
	default:
		return 0
	}
}

// Player manages audio playback and maintains pre-rendered sound buffers in memory.
type Player struct {
	mu           sync.RWMutex
	enabled      bool
	profile      SoundProfile
	themeProfile string
	thockWAV     []byte
	clickWAV     []byte
	spatialGrid  [numProfiles][2][numSpatialZones][numVariations][]byte
	spatialPCM   [numProfiles][2][numSpatialZones][numVariations][]int16
	engine       audioEngine
	varSeq       uint32
}

// NewPlayer initializes the acoustic player with synthesized mechanical sound buffers.
func NewPlayer() *Player {
	p := &Player{
		enabled: true,
		profile: ProfileAuto,
		engine:  newAudioEngine(),
	}

	// Pre-render the 5 (profiles) x 2 (types) x 5 (zones) x 3 (variations) stereo matrix
	// for zero-allocation, 0ms latency polyphonic audio playback at keystroke time.
	for pIdx, prof := range soundProfiles {
		for sIdx, st := range []SoundType{SoundThock, SoundClick} {
			for zIdx, pan := range zonePanValues {
				for vIdx, pitch := range variationPitchFactors {
					pcm := SynthesizeProfileStereoSamples(prof, st, pan, pitch)
					p.spatialPCM[pIdx][sIdx][zIdx][vIdx] = pcm
					p.spatialGrid[pIdx][sIdx][zIdx][vIdx] = EncodeInterleavedStereoWAV(pcm, SampleRate)
				}
			}
		}
	}

	// Default baseline WAV buffers for backward compatibility
	p.thockWAV = p.spatialGrid[0][0][2][0]
	p.clickWAV = p.spatialGrid[0][1][2][0]

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

// SetProfile configures the active acoustic archetype profile.
func (p *Player) SetProfile(prof SoundProfile) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.profile = prof
}

// GetProfile returns the configured acoustic profile.
func (p *Player) GetProfile() SoundProfile {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.profile
}

// SetThemeProfile sets the visual theme's signature acoustic archetype for ProfileAuto mode.
func (p *Player) SetThemeProfile(thmProf string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.themeProfile = thmProf
}

// PlaySpatial plays an equal-power panned acoustic waveform with organic micro-pitch variation
// based on the key's physical position on the keyboard.
func (p *Player) PlaySpatial(st SoundType, key string) {
	if !p.IsEnabled() {
		return
	}
	p.mu.RLock()
	prof := p.profile
	thmProf := p.themeProfile
	p.mu.RUnlock()

	targetProfile := prof
	if targetProfile == ProfileAuto || targetProfile == "" {
		if thmProf != "" {
			targetProfile = SoundProfile(thmProf)
		} else {
			targetProfile = ProfileThock
		}
	}

	pIdx := profileToIndex(targetProfile)
	pan := KeyToPan(key)
	zone := PanToZone(pan)
	soundIdx := 0
	if st == SoundClick {
		soundIdx = 1
	}
	v := int(atomic.AddUint32(&p.varSeq, 1) % numVariations)

	if p.engine != nil {
		pcm := p.spatialPCM[pIdx][soundIdx][zone][v]
		p.engine.PlaySamples(pcm)
	} else {
		wav := p.spatialGrid[pIdx][soundIdx][zone][v]
		_ = playWAVMemory(wav)
	}
}

// PlayThock plays a low-frequency, deep mechanical thock panned center.
func (p *Player) PlayThock() {
	p.PlaySpatial(SoundThock, "Space")
}

// PlayClick plays a crisp, bright tactile switch click panned center.
func (p *Player) PlayClick() {
	p.PlaySpatial(SoundClick, "Space")
}

// PlaySound plays the specified sound type.
func (p *Player) PlaySound(st SoundType) {
	p.PlaySpatial(st, "Space")
}

// GetWAV returns the raw mono WAV byte slice for a given sound type.
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
	if p.engine != nil {
		p.engine.Stop()
	}
	stopSoundMemory()
}

// Close gracefully releases any native streaming audio resources.
func (p *Player) Close() {
	if p.engine != nil {
		p.engine.Close()
	}
}

// KeyToPan returns the stereo panning coordinate [-0.65, +0.65] for a given key label.
func KeyToPan(key string) float64 {
	k := strings.ToLower(strings.TrimSpace(key))
	switch k {
	// Far Left (-0.65)
	case "esc", "escape", "`", "~", "tab", "⇥", "caps", "capslock",
		"shift", "lshift", "leftshift", "⇧", "left shift",
		"ctrl", "lctrl", "leftctrl", "control", "left control",
		"win", "lwin", "leftwin", "cmd", "super", "left cmd", "left command",
		"alt", "lalt", "leftalt", "left alt", "⌥",
		"q", "a", "z", "1", "!", "f1", "f2":
		return -0.65

	// Mid Left (-0.30)
	case "w", "s", "x", "e", "d", "c", "r", "f", "v",
		"2", "@", "3", "#", "4", "$", "5", "%",
		"f3", "f4", "f5":
		return -0.30

	// Center (0.00)
	case "t", "g", "b", "y", "h", "n",
		"6", "^", "7", "&",
		"f6", "f7", "space", "spacebar", " ",
		"middleclick", "mouse3", "wheel":
		return 0.0

	// Mid Right (+0.30)
	case "u", "j", "m", "i", "k", ",", "<",
		"8", "*", "o", "l", ".", ">", "9", "(",
		"f8", "f9":
		return 0.30

	// Far Right (+0.65)
	case "p", "[", "{", "]", "}", "\\", "|",
		";", ":", "'", "\"", "/", "?",
		"0", ")", "-", "_", "=", "+",
		"backspace", "⌫", "enter", "return", "↵",
		"rshift", "rightshift", "right shift",
		"rctrl", "rightctrl", "right control",
		"ralt", "rightalt", "right alt",
		"delete", "del", "insert", "home", "end", "pageup", "pagedown", "pgup", "pgdn",
		"up", "down", "left", "right", "↑", "↓", "←", "→",
		"f10", "f11", "f12",
		"rightclick", "mouse2":
		return 0.65

	// Mouse left click: subtle left bias
	case "leftclick", "mouse1", "click":
		return -0.20

	default:
		return 0.0
	}
}

// PanToZone maps a continuous pan coordinate to one of the 5 pre-rendered spatial zones.
func PanToZone(pan float64) SpatialZone {
	if pan <= -0.45 {
		return ZoneFarLeft
	}
	if pan <= -0.15 {
		return ZoneMidLeft
	}
	if pan <= 0.15 {
		return ZoneCenter
	}
	if pan <= 0.45 {
		return ZoneMidRight
	}
	return ZoneFarRight
}

// equalPowerStereo computes constant-power left and right gains for pan in [-1.0, 1.0].
// Guaranteed: leftGain^2 + rightGain^2 = 1.0, preserving uniform perceived acoustic loudness.
func equalPowerStereo(pan float64) (float64, float64) {
	if pan < -1.0 {
		pan = -1.0
	} else if pan > 1.0 {
		pan = 1.0
	}
	theta := (pan + 1.0) * (math.Pi / 4.0)
	return math.Cos(theta), math.Sin(theta)
}

// SynthesizeProfileStereoSamples procedurally renders an interleaved 16-bit PCM stereo sample slice
// for a specific acoustic profile archetype: [L0, R0, L1, R1, ...]
func SynthesizeProfileStereoSamples(prof SoundProfile, st SoundType, pan float64, pitchFactor float64) []int16 {
	var raw []int16
	switch prof {
	case ProfileThock:
		if st == SoundThock {
			raw = synthesizeThockWithPitch(pitchFactor)
		} else {
			raw = synthesizeThockKeyWithPitch(pitchFactor)
		}
	case ProfileClick:
		if st == SoundThock {
			raw = synthesizeClickModWithPitch(pitchFactor)
		} else {
			raw = synthesizeClickWithPitch(pitchFactor)
		}
	case ProfileScissor:
		if st == SoundThock {
			raw = synthesizeScissorModWithPitch(pitchFactor)
		} else {
			raw = synthesizeScissorKeyWithPitch(pitchFactor)
		}
	case ProfileGlass:
		if st == SoundThock {
			raw = synthesizeGlassModWithPitch(pitchFactor)
		} else {
			raw = synthesizeGlassKeyWithPitch(pitchFactor)
		}
	case ProfileTypewriter:
		if st == SoundThock {
			raw = synthesizeTypewriterModWithPitch(pitchFactor)
		} else {
			raw = synthesizeTypewriterKeyWithPitch(pitchFactor)
		}
	default:
		raw = synthesizeThockWithPitch(pitchFactor)
	}

	leftGain, rightGain := equalPowerStereo(pan)
	n := len(raw)
	out := make([]int16, n*2)
	for i := 0; i < n; i++ {
		out[i*2] = int16(float64(raw[i]) * leftGain)
		out[i*2+1] = int16(float64(raw[i]) * rightGain)
	}
	return out
}

// SynthesizeStereoSamples procedurally renders an interleaved 16-bit PCM stereo sample slice
// with equal-power panning and organic micro-pitch variation: [L0, R0, L1, R1, ...]
func SynthesizeStereoSamples(st SoundType, pan float64, pitchFactor float64) []int16 {
	var raw []int16
	switch st {
	case SoundThock:
		raw = synthesizeThockWithPitch(pitchFactor)
	case SoundClick:
		raw = synthesizeClickWithPitch(pitchFactor)
	default:
		raw = synthesizeThockWithPitch(pitchFactor)
	}

	leftGain, rightGain := equalPowerStereo(pan)
	n := len(raw)
	out := make([]int16, n*2)
	for i := 0; i < n; i++ {
		out[i*2] = int16(float64(raw[i]) * leftGain)
		out[i*2+1] = int16(float64(raw[i]) * rightGain)
	}
	return out
}

// SynthesizeStereoWAV procedurally renders a stereo 16-bit PCM WAV with equal-power panning
// and organic micro-pitch variation.
func SynthesizeStereoWAV(st SoundType, pan float64, pitchFactor float64) []byte {
	var raw []int16
	switch st {
	case SoundThock:
		raw = synthesizeThockWithPitch(pitchFactor)
	case SoundClick:
		raw = synthesizeClickWithPitch(pitchFactor)
	default:
		raw = synthesizeThockWithPitch(pitchFactor)
	}

	leftGain, rightGain := equalPowerStereo(pan)
	n := len(raw)
	left := make([]int16, n)
	right := make([]int16, n)

	for i := 0; i < n; i++ {
		lVal := float64(raw[i]) * leftGain
		rVal := float64(raw[i]) * rightGain
		left[i] = int16(lVal)
		right[i] = int16(rVal)
	}

	return EncodeStereoWAV(left, right, SampleRate)
}

// SynthesizeWAV procedurally synthesizes in-memory RIFF WAV bytes for the given SoundType.
func SynthesizeWAV(st SoundType) []byte {
	samples := SynthesizeSamples(st)
	return EncodeWAV(samples, SampleRate)
}

// SynthesizeSamples generates raw 16-bit mono PCM samples for a sound type at nominal pitch.
func SynthesizeSamples(st SoundType) []int16 {
	switch st {
	case SoundThock:
		return synthesizeThockWithPitch(1.0)
	case SoundClick:
		return synthesizeClickWithPitch(1.0)
	default:
		return synthesizeThockWithPitch(1.0)
	}
}

// synthesizeThockWithPitch generates a low-frequency mechanical switch bottom-out (~180Hz-220Hz, ~35ms)
// scaled by an organic pitch factor.
func synthesizeThockWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.035 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	// Deterministic pseudo-random number generator for plastic impact transient
	var prngState uint32 = 0x1337BEEF
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phase := 0.0
	attackSec := 0.0015
	decayTau := 0.010 / pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)

		// Frequency glide: starts at ~220Hz dropping towards ~180Hz
		freq := (180.0 + 40.0*math.Exp(-t/0.008)) * pitchFactor
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

// synthesizeClickWithPitch generates a crisp, bright tactile click transient (~1200Hz-2400Hz, ~15ms)
// scaled by an organic pitch factor.
func synthesizeClickWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.015 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	var prngState uint32 = 0x98765432
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phase := 0.0
	attackSec := 0.0006
	decayTau := 0.0032 / pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)

		// High frequency sweep: 2400Hz -> 1200Hz
		freq := (1200.0 + 1200.0*math.Exp(-t/0.003)) * pitchFactor
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
		samples[i] = softLimit(val)
	}

	return samples
}

// synthesizeThockKeyWithPitch generates a punchy, warm mechanical key bottom-out (~220Hz-260Hz, ~26ms).
func synthesizeThockKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.026 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	var prngState uint32 = 0x2468ACE0
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phase := 0.0
	attackSec := 0.0010
	decayTau := 0.0075 / pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := (200.0 + 60.0*math.Exp(-t/0.006)) * pitchFactor
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		body := math.Sin(phase) + 0.35*math.Sin(2.0*phase)
		transientNoise := randFloat() * math.Exp(-t/0.0018) * 0.28

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		remaining := numSamples - 1 - i
		if remaining < 20 {
			env *= float64(remaining) / 20.0
		}

		val := (body + transientNoise) * env * 24000.0
		samples[i] = softLimit(val)
	}
	return samples
}

// synthesizeClickModWithPitch generates a tactile modifier switch click (~1400Hz, ~20ms).
func synthesizeClickModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.020 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	var prngState uint32 = 0x13579BDF
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phase := 0.0
	attackSec := 0.0008
	decayTau := 0.0045 / pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := (800.0 + 800.0*math.Exp(-t/0.004)) * pitchFactor
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		tone := 0.70*math.Sin(phase) + 0.30*math.Sin(1.5*phase)
		snap := randFloat() * math.Exp(-t/0.0016) * 0.35

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		remaining := numSamples - 1 - i
		if remaining < 20 {
			env *= float64(remaining) / 20.0
		}

		val := (tone + snap) * env * 22000.0
		samples[i] = softLimit(val)
	}
	return samples
}

// synthesizeScissorKeyWithPitch models an Apple MacBook low-profile scissor switch (~1200Hz, ~9ms).
func synthesizeScissorKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.009 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	var prngState uint32 = 0xABCDEF01
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phase := 0.0
	attackSec := 0.0004
	decayTau := 0.0022 / pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := (850.0 + 550.0*math.Exp(-t/0.002)) * pitchFactor
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		tone := 0.70*math.Sin(phase) + 0.30*math.Sin(2.0*phase)
		snap := randFloat() * math.Exp(-t/0.0008) * 0.30

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		remaining := numSamples - 1 - i
		if remaining < 15 {
			env *= float64(remaining) / 15.0
		}

		val := (tone + snap) * env * 19000.0
		samples[i] = softLimit(val)
	}
	return samples
}

// synthesizeScissorModWithPitch models a damped aluminum scissor bottom-out (~700Hz, ~12ms).
func synthesizeScissorModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.012 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	phase := 0.0
	attackSec := 0.0006
	decayTau := 0.0035 / pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := (500.0 + 350.0*math.Exp(-t/0.003)) * pitchFactor
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		tone := 0.80*math.Sin(phase) + 0.20*math.Sin(2.0*phase)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		remaining := numSamples - 1 - i
		if remaining < 15 {
			env *= float64(remaining) / 15.0
		}

		val := tone * env * 18000.0
		samples[i] = softLimit(val)
	}
	return samples
}

// synthesizeGlassKeyWithPitch models a visionOS crystalline liquid glass tap (~1760Hz, ~18ms).
func synthesizeGlassKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.018 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	phase1 := 0.0
	phase2 := 0.0
	attackSec := 0.0003
	decayTau := 0.0050 / pitchFactor

	freq1 := 1760.0 * pitchFactor
	freq2 := 3520.0 * pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		phase1 += 2.0 * math.Pi * freq1 / float64(SampleRate)
		phase2 += 2.0 * math.Pi * freq2 / float64(SampleRate)

		tone := 0.80*math.Sin(phase1) + 0.20*math.Sin(phase2)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		remaining := numSamples - 1 - i
		if remaining < 20 {
			env *= float64(remaining) / 20.0
		}

		val := tone * env * 22000.0
		samples[i] = softLimit(val)
	}
	return samples
}

// synthesizeGlassModWithPitch models a deep crystal resonant glass modifier chime (~1175Hz, ~22ms).
func synthesizeGlassModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.022 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	phase := 0.0
	attackSec := 0.0004
	decayTau := 0.0065 / pitchFactor

	freq := 1175.0 * pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		tone := 0.85*math.Sin(phase) + 0.15*math.Sin(2.0*phase)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		remaining := numSamples - 1 - i
		if remaining < 20 {
			env *= float64(remaining) / 20.0
		}

		val := tone * env * 21000.0
		samples[i] = softLimit(val)
	}
	return samples
}

// synthesizeTypewriterKeyWithPitch models a 1984 IBM Model M buckling spring strike (~28ms).
func synthesizeTypewriterKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.028 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	var prngState uint32 = 0x54321678
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phaseImpact := 0.0
	phaseSpring := 0.0
	attackSec := 0.0005
	decayTau := 0.0065 / pitchFactor

	springFreq := 750.0 * pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freqImpact := (900.0 + 1300.0*math.Exp(-t/0.002)) * pitchFactor
		phaseImpact += 2.0 * math.Pi * freqImpact / float64(SampleRate)
		phaseSpring += 2.0 * math.Pi * springFreq / float64(SampleRate)

		impact := 0.50*math.Sin(phaseImpact) + randFloat()*math.Exp(-t/0.0015)*0.50
		springRing := 0.40 * math.Sin(phaseSpring) * math.Exp(-t/0.009)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		remaining := numSamples - 1 - i
		if remaining < 25 {
			env *= float64(remaining) / 25.0
		}

		val := (impact + springRing) * env * 24000.0
		samples[i] = softLimit(val)
	}
	return samples
}

// synthesizeTypewriterModWithPitch models an industrial typewriter carriage return spring clack (~32ms).
func synthesizeTypewriterModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.032 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	var prngState uint32 = 0x87654321
	randFloat := func() float64 {
		prngState = prngState*1664525 + 1013904223
		return (float64(prngState)/4294967296.0)*2.0 - 1.0
	}

	phaseSpring := 0.0
	attackSec := 0.0008
	decayTau := 0.0090 / pitchFactor

	springFreq := 450.0 * pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		phaseSpring += 2.0 * math.Pi * springFreq / float64(SampleRate)

		impact := randFloat() * math.Exp(-t/0.002) * 0.40
		spring := 0.70*math.Sin(phaseSpring) + 0.30*math.Sin(2.0*phaseSpring)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}

		remaining := numSamples - 1 - i
		if remaining < 25 {
			env *= float64(remaining) / 25.0
		}

		val := (impact + spring) * env * 23000.0
		samples[i] = softLimit(val)
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

// EncodeStereoWAV serializes 16-bit stereo PCM samples into a standard 44-byte RIFF WAV byte array.
func EncodeStereoWAV(left, right []int16, sampleRate int) []byte {
	numFrames := len(left)
	if len(right) < numFrames {
		numFrames = len(right)
	}
	dataSize := numFrames * 2 * 2 // 2 channels * 2 bytes per sample
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
	binary.LittleEndian.PutUint16(buf[22:24], 2)  // NumChannels 2 = Stereo
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	byteRate := uint32(sampleRate * 2 * 2) // sampleRate * NumChannels * BitsPerSample / 8 = sampleRate * 4
	binary.LittleEndian.PutUint32(buf[28:32], byteRate)
	binary.LittleEndian.PutUint16(buf[32:34], 4)  // BlockAlign = NumChannels * BitsPerSample / 8 = 4
	binary.LittleEndian.PutUint16(buf[34:36], 16) // BitsPerSample = 16

	// "data" sub-chunk
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))

	// Interleaved PCM samples payload: [L0, R0, L1, R1, ...]
	offset := 44
	for i := 0; i < numFrames; i++ {
		binary.LittleEndian.PutUint16(buf[offset:offset+2], uint16(left[i]))
		binary.LittleEndian.PutUint16(buf[offset+2:offset+4], uint16(right[i]))
		offset += 4
	}

	return buf
}

// EncodeInterleavedStereoWAV serializes interleaved 16-bit stereo PCM samples into a standard 44-byte RIFF WAV byte array.
func EncodeInterleavedStereoWAV(samples []int16, sampleRate int) []byte {
	numFrames := len(samples) / 2
	dataSize := numFrames * 4 // 2 channels * 2 bytes
	totalSize := 44 + dataSize
	buf := make([]byte, totalSize)

	// RIFF Chunk Descriptor
	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], uint32(36+dataSize))
	copy(buf[8:12], "WAVE")

	// "fmt " sub-chunk
	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16)
	binary.LittleEndian.PutUint16(buf[20:22], 1)
	binary.LittleEndian.PutUint16(buf[22:24], 2)
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(buf[28:32], uint32(sampleRate*4))
	binary.LittleEndian.PutUint16(buf[32:34], 4)
	binary.LittleEndian.PutUint16(buf[34:36], 16)

	// "data" sub-chunk
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))

	offset := 44
	for _, s := range samples {
		binary.LittleEndian.PutUint16(buf[offset:offset+2], uint16(s))
		offset += 2
	}
	return buf
}
