package acoustic

import "math"

// The material voices below extend the mechanical archetypes in synth.go with a
// second family: sounds borrowed from wood, stone, paper, felt, water, crystal,
// silicon and a camera shutter. Each profile has a key voice (character keys)
// and a mod voice (modifiers), just like the original five. They share the same
// pre-rendered stereo grid, so the hot path stays zero-allocation.

// materialPRNG returns a deterministic low-discrepancy noise source. Determinism
// matters because every buffer is rendered once at startup.
func materialPRNG(seed uint32) func() float64 {
	state := seed
	return func() float64 {
		state = state*1664525 + 1013904223
		return (float64(state)/4294967296.0)*2.0 - 1.0
	}
}

// tailFade prevents a DC pop at the end of a material buffer.
func tailFade(remaining, fade int) float64 {
	if remaining < fade {
		return float64(remaining) / float64(fade)
	}
	return 1.0
}

// materialSample clamps a float sample into the 16-bit PCM range.
func materialSample(x float64) int16 {
	if x > 32767.0 {
		return 32767
	}
	if x < -32768.0 {
		return -32768
	}
	return int16(x)
}

// onePoleLP is a cheap one-pole low-pass used to shape noise into material
// textures (paper hiss, stone grit, felt air).
type onePoleLP struct {
	prev  float64
	alpha float64
}

func newOnePoleLP(cutoffHz float64) onePoleLP {
	dt := 1.0 / float64(SampleRate)
	rc := 1.0 / (2.0 * math.Pi * cutoffHz)
	return onePoleLP{alpha: dt / (rc + dt)}
}

func (f *onePoleLP) process(x float64) float64 {
	f.prev += f.alpha * (x - f.prev)
	return f.prev
}

// ---------------------------------------------------------------------------
// 柚木 Hinoki — a warm wooden bar struck by a soft mallet (~330 Hz, 45 ms).
// Fundamental + marimba-like 1:4 partial, granular 8 ms attack, no ring.
// ---------------------------------------------------------------------------

func synthesizeHinokiKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.045 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0x48A1C0DE)
	phase1, phase2, phase4 := 0.0, 0.0, 0.0
	attackSec := 0.0015
	decayTau := 0.011 / pitchFactor
	base := 330.0 * pitchFactor
	grainLP := newOnePoleLP(2800)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := base * (1.0 + 0.05*math.Exp(-t/0.010))
		phase1 += 2.0 * math.Pi * freq / float64(SampleRate)
		phase2 += 2.0 * math.Pi * freq * 2.0 / float64(SampleRate)
		phase4 += 2.0 * math.Pi * freq * 4.0 / float64(SampleRate)

		body := math.Sin(phase1) + 0.15*math.Sin(phase2) + 0.22*math.Sin(phase4)
		grain := grainLP.process(randFloat()) * math.Exp(-t/0.0035) * 0.35

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 40)

		samples[i] = materialSample((body + grain) * env * 21000.0)
	}
	return samples
}

func synthesizeHinokiModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.032 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0x48A1C0DF)
	phase1, phase2, phase4 := 0.0, 0.0, 0.0
	attackSec := 0.0010
	decayTau := 0.008 / pitchFactor
	base := 240.0 * pitchFactor
	grainLP := newOnePoleLP(2200)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := base * (1.0 + 0.04*math.Exp(-t/0.008))
		phase1 += 2.0 * math.Pi * freq / float64(SampleRate)
		phase2 += 2.0 * math.Pi * freq * 2.0 / float64(SampleRate)
		phase4 += 2.0 * math.Pi * freq * 4.0 / float64(SampleRate)

		body := math.Sin(phase1) + 0.12*math.Sin(phase2) + 0.16*math.Sin(phase4)
		grain := grainLP.process(randFloat()) * math.Exp(-t/0.0028) * 0.22

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 30)

		samples[i] = materialSample((body + grain) * env * 18000.0)
	}
	return samples
}

// ---------------------------------------------------------------------------
// 卵石 Pebble — two dense stones meeting (~120 Hz, 18 ms). No sustain, just a
// gritty, grounded thud that disappears before it can fatigue the ear.
// ---------------------------------------------------------------------------

func synthesizePebbleKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.018 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0x9EBB1E55)
	phase := 0.0
	attackSec := 0.0004
	decayTau := 0.0034 / pitchFactor
	base := 120.0 * pitchFactor
	gritLow := newOnePoleLP(700)
	gritHigh := newOnePoleLP(2600)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := base * (1.0 + 0.28*math.Exp(-t/0.004))
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		body := math.Sin(phase) + 0.18*math.Sin(1.73*phase)
		noise := randFloat()
		grit := gritHigh.process(noise-gritLow.process(noise)) * math.Exp(-t/0.0016) * 0.55

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 24)

		samples[i] = materialSample((body + grit) * env * 23000.0)
	}
	return samples
}

func synthesizePebbleModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.015 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0x9EBB1E56)
	phase := 0.0
	attackSec := 0.0003
	decayTau := 0.0028 / pitchFactor
	base := 92.0 * pitchFactor
	gritLow := newOnePoleLP(650)
	gritHigh := newOnePoleLP(2200)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := base * (1.0 + 0.25*math.Exp(-t/0.0035))
		phase += 2.0 * math.Pi * freq / float64(SampleRate)

		body := math.Sin(phase) + 0.16*math.Sin(1.73*phase)
		noise := randFloat()
		grit := gritHigh.process(noise-gritLow.process(noise)) * math.Exp(-t/0.0013) * 0.40

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 20)

		samples[i] = materialSample((body + grit) * env * 20000.0)
	}
	return samples
}

// ---------------------------------------------------------------------------
// 纸页 Paper — pure filtered noise, no pitch at all (~7 ms). A page turning,
// not a key being struck. The only tone-free voice in the family.
// ---------------------------------------------------------------------------

func synthesizePaperKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.010 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0xBA9E1F00)
	attackSec := 0.00030
	decayTau := 0.0020 / pitchFactor
	hpA := newOnePoleLP(1900)
	hpB := newOnePoleLP(1900)
	lp := newOnePoleLP(5200)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		noise := randFloat()
		high := noise - hpA.process(noise)
		high = high - hpB.process(high)
		band := lp.process(high)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 16)

		samples[i] = materialSample(band * env * 56000.0)
	}
	return samples
}

func synthesizePaperModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.014 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0xBA9E1F01)
	attackSec := 0.00045
	decayTau := 0.0030 / pitchFactor
	hpA := newOnePoleLP(1200)
	hpB := newOnePoleLP(1200)
	lp := newOnePoleLP(3800)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		noise := randFloat()
		high := noise - hpA.process(noise)
		high = high - hpB.process(high)
		band := lp.process(high)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 20)

		samples[i] = materialSample(band * env * 50000.0)
	}
	return samples
}

// ---------------------------------------------------------------------------
// 绒毡 Felt — a piano's soft pedal (~140 Hz, 50 ms). A mallet wrapped in wool:
// almost no attack, no brightness, just a warm hush.
// ---------------------------------------------------------------------------

func synthesizeFeltKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.050 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0xFE17C0DE)
	phase1, phase2 := 0.0, 0.0
	attackSec := 0.0020
	decayTau := 0.015 / pitchFactor
	base := 140.0 * pitchFactor
	airLP := newOnePoleLP(700)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := base * (1.0 + 0.12*math.Exp(-t/0.012))
		phase1 += 2.0 * math.Pi * freq / float64(SampleRate)
		phase2 += 2.0 * math.Pi * freq * 2.0 / float64(SampleRate)

		body := math.Sin(phase1) + 0.12*math.Sin(phase2)
		air := airLP.process(randFloat()) * math.Exp(-t/0.008) * 0.22

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 45)

		samples[i] = materialSample((body + air) * env * 21000.0)
	}
	return samples
}

func synthesizeFeltModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.040 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0xFE17C0DF)
	phase1, phase2 := 0.0, 0.0
	attackSec := 0.0016
	decayTau := 0.012 / pitchFactor
	base := 100.0 * pitchFactor
	airLP := newOnePoleLP(600)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := base * (1.0 + 0.10*math.Exp(-t/0.010))
		phase1 += 2.0 * math.Pi * freq / float64(SampleRate)
		phase2 += 2.0 * math.Pi * freq * 2.0 / float64(SampleRate)

		body := math.Sin(phase1) + 0.10*math.Sin(phase2)
		air := airLP.process(randFloat()) * math.Exp(-t/0.006) * 0.16

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 36)

		samples[i] = materialSample((body + air) * env * 17000.0)
	}
	return samples
}

// ---------------------------------------------------------------------------
// 雨滴 Droplet — a single drop landing in still water (~500→900 Hz over 8 ms,
// 70 ms tail). The upward glide is what makes the ear hear water, not a bell.
// ---------------------------------------------------------------------------

func synthesizeDropletKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.075 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0xD20F1E00)
	phase1, phase2 := 0.0, 0.0
	attackSec := 0.0008
	decayTau := 0.020 / pitchFactor
	base := 900.0 * pitchFactor
	airLP := newOnePoleLP(5000)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := base - 400.0*pitchFactor*math.Exp(-t/0.008)
		phase1 += 2.0 * math.Pi * freq / float64(SampleRate)
		phase2 += 2.0 * math.Pi * freq * 2.02 / float64(SampleRate)

		body := math.Sin(phase1) + 0.15*math.Sin(phase2)
		splash := airLP.process(randFloat()) * math.Exp(-t/0.0015) * 0.06

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 50)

		samples[i] = materialSample((body + splash) * env * 20000.0)
	}
	return samples
}

func synthesizeDropletModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.090 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0xD20F1E01)
	phase1, phase2 := 0.0, 0.0
	attackSec := 0.0010
	decayTau := 0.024 / pitchFactor
	base := 700.0 * pitchFactor
	airLP := newOnePoleLP(4200)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := base - 320.0*pitchFactor*math.Exp(-t/0.010)
		phase1 += 2.0 * math.Pi * freq / float64(SampleRate)
		phase2 += 2.0 * math.Pi * freq * 2.02 / float64(SampleRate)

		body := math.Sin(phase1) + 0.18*math.Sin(phase2)
		splash := airLP.process(randFloat()) * math.Exp(-t/0.0018) * 0.08

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 60)

		samples[i] = materialSample((body + splash) * env * 19000.0)
	}
	return samples
}

// ---------------------------------------------------------------------------
// 晶体 Crystal — an inharmonic shard (~3.2 kHz, 32 ms). Four detuned partials
// beat against each other; the result shimmers instead of ringing.
// ---------------------------------------------------------------------------

func synthesizeCrystalKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.032 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	phase1, phase1b, phase2, phase3 := 0.0, 0.0, 0.0, 0.0
	attackSec := 0.00030
	decayTau := 0.0075 / pitchFactor
	base := 3200.0 * pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		phase1 += 2.0 * math.Pi * base / float64(SampleRate)
		phase1b += 2.0 * math.Pi * base * 1.006 / float64(SampleRate)
		phase2 += 2.0 * math.Pi * base * 1.41 / float64(SampleRate)
		phase3 += 2.0 * math.Pi * base * 2.09 / float64(SampleRate)

		tone := 0.50*math.Sin(phase1) + 0.40*math.Sin(phase1b) +
			0.28*math.Sin(phase2) + 0.16*math.Sin(phase3)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 32)

		samples[i] = materialSample(tone * env * 14000.0)
	}
	return samples
}

func synthesizeCrystalModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.024 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	phase1, phase1b, phase2 := 0.0, 0.0, 0.0
	attackSec := 0.00025
	decayTau := 0.0055 / pitchFactor
	base := 1920.0 * pitchFactor

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		phase1 += 2.0 * math.Pi * base / float64(SampleRate)
		phase1b += 2.0 * math.Pi * base * 1.008 / float64(SampleRate)
		phase2 += 2.0 * math.Pi * base * 1.41 / float64(SampleRate)

		tone := 0.55*math.Sin(phase1) + 0.35*math.Sin(phase1b) + 0.20*math.Sin(phase2)

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 24)

		samples[i] = materialSample(tone * env * 12000.0)
	}
	return samples
}

// ---------------------------------------------------------------------------
// 街机 Chip — an 8-bit square blip (~660 Hz, 28 ms) with sample-rate reduction
// and coarse amplitude quantization. A deliberate, affectionate lo-fi voice.
// ---------------------------------------------------------------------------

func synthesizeChipKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.028 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	phase := 0.0
	attackSec := 0.0005
	decayTau := 0.009 / pitchFactor
	base := 660.0 * pitchFactor
	const holdSamples = 3
	const quantStep = 1800.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		if i%holdSamples == 0 {
			phase += 2.0 * math.Pi * base / float64(SampleRate)
		}

		square := 1.0
		if math.Sin(phase) < 0 {
			square = -1.0
		}

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 20)

		val := math.Round(square*env*14000.0/quantStep) * quantStep
		samples[i] = materialSample(val)
	}
	return samples
}

func synthesizeChipModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.020 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	phase := 0.0
	attackSec := 0.0004
	decayTau := 0.006 / pitchFactor
	base := 330.0 * pitchFactor
	const holdSamples = 4
	const quantStep = 2200.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		if i%holdSamples == 0 {
			phase += 2.0 * math.Pi * base / float64(SampleRate)
		}

		square := 1.0
		if math.Sin(phase) < 0 {
			square = -1.0
		}

		var env float64
		if t < attackSec {
			env = t / attackSec
		} else {
			env = math.Exp(-(t - attackSec) / decayTau)
		}
		env *= tailFade(numSamples-1-i, 16)

		val := math.Round(square*env*12000.0/quantStep) * quantStep
		samples[i] = materialSample(val)
	}
	return samples
}

// ---------------------------------------------------------------------------
// 快门 Shutter — a Leica-style two-stage release: a 1 ms blade click, then a
// 200 Hz mirror thump 12 ms later. The gap is the signature; without it the
// sound is just another click.
// ---------------------------------------------------------------------------

func synthesizeShutterKeyWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.042 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0x5A077E00)
	phasePing, phaseThump := 0.0, 0.0
	pingFreq := 3200.0 * pitchFactor
	thumpFreq := 205.0 * pitchFactor
	const secondStageSec = 0.012

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		phasePing += 2.0 * math.Pi * pingFreq / float64(SampleRate)
		phaseThump += 2.0 * math.Pi * thumpFreq / float64(SampleRate)

		blade := randFloat()*math.Exp(-t/0.0007)*0.75 +
			math.Sin(phasePing)*math.Exp(-t/0.0016)*0.30

		stageTwo := 0.0
		if t >= secondStageSec {
			dt := t - secondStageSec
			stageTwo = math.Sin(phaseThump)*math.Exp(-dt/0.006)*0.85 +
				randFloat()*math.Exp(-dt/0.0010)*0.30
		}

		val := (blade + stageTwo) * 19000.0
		val *= tailFade(numSamples-1-i, 36)
		samples[i] = materialSample(val)
	}
	return samples
}

func synthesizeShutterModWithPitch(pitchFactor float64) []int16 {
	if pitchFactor <= 0 {
		pitchFactor = 1.0
	}
	durationSec := 0.050 / math.Sqrt(pitchFactor)
	numSamples := int(durationSec * float64(SampleRate))
	samples := make([]int16, numSamples)

	randFloat := materialPRNG(0x5A077E01)
	phaseRatchet, phaseSettle := 0.0, 0.0
	ratchetFreq := 1300.0 * pitchFactor
	settleFreq := 320.0 * pitchFactor
	starts := [3]float64{0.000, 0.010, 0.020}
	amps := [3]float64{1.00, 0.75, 0.55}
	const settleSec = 0.032

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		phaseRatchet += 2.0 * math.Pi * ratchetFreq / float64(SampleRate)
		phaseSettle += 2.0 * math.Pi * settleFreq / float64(SampleRate)

		ticks := 0.0
		for k, start := range starts {
			if t >= start {
				dt := t - start
				ticks += amps[k] * (math.Sin(phaseRatchet)*math.Exp(-dt/0.0016) +
					randFloat()*math.Exp(-dt/0.0007)*0.55)
			}
		}

		settle := 0.0
		if t >= settleSec {
			dt := t - settleSec
			settle = math.Sin(phaseSettle) * math.Exp(-dt/0.005) * 0.45
		}

		val := (ticks + settle) * 18000.0
		val *= tailFade(numSamples-1-i, 40)
		samples[i] = materialSample(val)
	}
	return samples
}
