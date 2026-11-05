package acoustic_test

import (
	"encoding/binary"
	"sync"
	"testing"
	"time"

	"github.com/yviscool/cue/internal/acoustic"
)

func TestWAVHeaderValidation(t *testing.T) {
	soundTypes := []acoustic.SoundType{
		acoustic.SoundThock,
		acoustic.SoundClick,
	}

	for _, st := range soundTypes {
		t.Run(st.String(), func(t *testing.T) {
			wav := acoustic.SynthesizeWAV(st)
			if len(wav) < 44 {
				t.Fatalf("WAV buffer too small: %d bytes (minimum 44)", len(wav))
			}

			// RIFF header
			if string(wav[0:4]) != "RIFF" {
				t.Errorf("expected 'RIFF' chunk ID, got %q", string(wav[0:4]))
			}

			chunkSize := binary.LittleEndian.Uint32(wav[4:8])
			if int(chunkSize)+8 != len(wav) {
				t.Errorf("chunk size mismatch: header says %d (+8 = %d), actual %d",
					chunkSize, chunkSize+8, len(wav))
			}

			if string(wav[8:12]) != "WAVE" {
				t.Errorf("expected 'WAVE' format, got %q", string(wav[8:12]))
			}

			// "fmt " subchunk
			if string(wav[12:16]) != "fmt " {
				t.Errorf("expected 'fmt ' subchunk ID, got %q", string(wav[12:16]))
			}

			subchunk1Size := binary.LittleEndian.Uint32(wav[16:20])
			if subchunk1Size != 16 {
				t.Errorf("expected subchunk1Size 16, got %d", subchunk1Size)
			}

			audioFormat := binary.LittleEndian.Uint16(wav[20:22])
			if audioFormat != 1 {
				t.Errorf("expected PCM audio format (1), got %d", audioFormat)
			}

			numChannels := binary.LittleEndian.Uint16(wav[22:24])
			if numChannels != 1 {
				t.Errorf("expected Mono (1 channel), got %d", numChannels)
			}

			sampleRate := binary.LittleEndian.Uint32(wav[24:28])
			if sampleRate != acoustic.SampleRate {
				t.Errorf("expected sample rate %d, got %d", acoustic.SampleRate, sampleRate)
			}

			byteRate := binary.LittleEndian.Uint32(wav[28:32])
			expectedByteRate := uint32(acoustic.SampleRate * 1 * 2)
			if byteRate != expectedByteRate {
				t.Errorf("expected byte rate %d, got %d", expectedByteRate, byteRate)
			}

			blockAlign := binary.LittleEndian.Uint16(wav[32:34])
			if blockAlign != 2 {
				t.Errorf("expected blockAlign 2, got %d", blockAlign)
			}

			bitsPerSample := binary.LittleEndian.Uint16(wav[34:36])
			if bitsPerSample != 16 {
				t.Errorf("expected bitsPerSample 16, got %d", bitsPerSample)
			}

			// "data" subchunk
			if string(wav[36:40]) != "data" {
				t.Errorf("expected 'data' subchunk ID, got %q", string(wav[36:40]))
			}

			dataSize := binary.LittleEndian.Uint32(wav[40:44])
			if int(dataSize)+44 != len(wav) {
				t.Errorf("data size mismatch: dataSize %d (+44 = %d), actual %d",
					dataSize, dataSize+44, len(wav))
			}
		})
	}
}

func countZeroCrossings(samples []int16) int {
	crossings := 0
	for i := 1; i < len(samples); i++ {
		if (samples[i-1] >= 0 && samples[i] < 0) || (samples[i-1] < 0 && samples[i] >= 0) {
			crossings++
		}
	}
	return crossings
}

func TestProceduralSynthesisCharacteristics(t *testing.T) {
	thockSamples := acoustic.SynthesizeSamples(acoustic.SoundThock)
	clickSamples := acoustic.SynthesizeSamples(acoustic.SoundClick)

	// Validate non-empty samples
	if len(thockSamples) == 0 {
		t.Fatal("thock samples should not be empty")
	}
	if len(clickSamples) == 0 {
		t.Fatal("click samples should not be empty")
	}

	// Duration checks:
	// Thock should be ~35ms (~1543 samples at 44100Hz)
	thockDurationMs := float64(len(thockSamples)) / float64(acoustic.SampleRate) * 1000.0
	if thockDurationMs < 30.0 || thockDurationMs > 45.0 {
		t.Errorf("expected Thock duration ~35ms, got %.2fms", thockDurationMs)
	}

	// Click should be ~15ms (~661 samples at 44100Hz)
	clickDurationMs := float64(len(clickSamples)) / float64(acoustic.SampleRate) * 1000.0
	if clickDurationMs < 10.0 || clickDurationMs > 20.0 {
		t.Errorf("expected Click duration ~15ms, got %.2fms", clickDurationMs)
	}

	// Peak amplitude check (waveform must not be silent)
	var maxThock int16
	for _, s := range thockSamples {
		if s > maxThock {
			maxThock = s
		}
	}
	if maxThock < 5000 {
		t.Errorf("Thock peak amplitude too low: %d (expected > 5000)", maxThock)
	}

	var maxClick int16
	for _, s := range clickSamples {
		if s > maxClick {
			maxClick = s
		}
	}
	if maxClick < 5000 {
		t.Errorf("Click peak amplitude too low: %d (expected > 5000)", maxClick)
	}

	// Acoustic frequency character check:
	// Click has high-frequency transient (1200-2400Hz) while Thock is low-frequency (180-220Hz).
	// Zero-crossing rate of Click must be significantly higher than Thock.
	thockCrossings := countZeroCrossings(thockSamples)
	clickCrossings := countZeroCrossings(clickSamples)

	thockCrossingRate := float64(thockCrossings) / (thockDurationMs / 1000.0)
	clickCrossingRate := float64(clickCrossings) / (clickDurationMs / 1000.0)

	if clickCrossingRate <= thockCrossingRate*2.0 {
		t.Errorf("Click crossing rate (%.1f/s) should be at least 2x higher than Thock crossing rate (%.1f/s)",
			clickCrossingRate, thockCrossingRate)
	}
}

func TestPlayerLifecycle(t *testing.T) {
	p := acoustic.NewPlayer()

	// Default state
	if !p.IsEnabled() {
		t.Errorf("expected Player to be enabled by default")
	}

	// Buffer validation
	if len(p.GetWAV(acoustic.SoundThock)) == 0 {
		t.Errorf("player thock WAV buffer is empty")
	}
	if len(p.GetWAV(acoustic.SoundClick)) == 0 {
		t.Errorf("player click WAV buffer is empty")
	}

	// Disabled playback safety
	p.SetEnabled(false)
	if p.IsEnabled() {
		t.Errorf("expected Player to be disabled after SetEnabled(false)")
	}
	p.PlayThock() // should silently no-op without panic
	p.PlayClick() // should silently no-op without panic

	// Enabled playback safety
	p.SetEnabled(true)
	if !p.IsEnabled() {
		t.Errorf("expected Player to be enabled after SetEnabled(true)")
	}

	// Playback execution (must not panic)
	p.PlayThock()
	time.Sleep(10 * time.Millisecond)
	p.PlayClick()
	time.Sleep(10 * time.Millisecond)
	p.PlaySound(acoustic.SoundThock)
	p.PlaySound(acoustic.SoundClick)
	p.Stop()
}

func TestConcurrentPlaybackSafety(t *testing.T) {
	p := acoustic.NewPlayer()
	var wg sync.WaitGroup

	// Concurrently play sounds from 10 goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 25; j++ {
				if id%2 == 0 {
					p.PlayThock()
				} else {
					p.PlayClick()
				}
				if j%10 == 0 {
					p.SetEnabled(p.IsEnabled())
				}
			}
		}(i)
	}

	wg.Wait()
	p.Stop()
}
