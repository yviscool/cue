//go:build windows

package acoustic

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"
)

var (
	winmmDLL             = syscall.NewLazyDLL("winmm.dll")
	kernel32DLL          = syscall.NewLazyDLL("kernel32.dll")
	procWaveOutOpen      = winmmDLL.NewProc("waveOutOpen")
	procWaveOutClose     = winmmDLL.NewProc("waveOutClose")
	procWaveOutPrepare   = winmmDLL.NewProc("waveOutPrepareHeader")
	procWaveOutUnprepare = winmmDLL.NewProc("waveOutUnprepareHeader")
	procWaveOutWrite     = winmmDLL.NewProc("waveOutWrite")
	procWaveOutReset     = winmmDLL.NewProc("waveOutReset")
	procPlaySound        = winmmDLL.NewProc("PlaySoundW")

	procCreateEventW        = kernel32DLL.NewProc("CreateEventW")
	procCloseHandle         = kernel32DLL.NewProc("CloseHandle")
	procWaitForSingleObject = kernel32DLL.NewProc("WaitForSingleObject")
)

const (
	waveFormatPCM = 1
	waveMapper    = ^uintptr(0) // (UINT)-1: system default audio endpoint
	callbackEvent = 0x00050000
	whdrDone      = 0x00000001
	whdrPrepared  = 0x00000002

	numRingBuffers = 4
	framesPerBuf   = 512              // 512 frames @ 44.1kHz = ~11.6ms
	bufSizeBytes   = framesPerBuf * 4 // 16-bit stereo = 4 bytes per frame
	maxVoices      = 32

	sndSync      = 0x0000
	sndAsync     = 0x0001
	sndNodefault = 0x0002
	sndMemory    = 0x0004
	sndNowait    = 0x2000
)

type waveFormatEx struct {
	wFormatTag      uint16
	nChannels       uint16
	nSamplesPerSec  uint32
	nAvgBytesPerSec uint32
	nBlockAlign     uint16
	wBitsPerSample  uint16
	cbSize          uint16
}

type waveHdr struct {
	lpData          uintptr
	dwBufferLength  uint32
	dwBytesRecorded uint32
	dwUser          uintptr
	dwFlags         uint32
	dwLoops         uint32
	lpNext          uintptr
	reserved        uintptr
}

type activeVoice struct {
	samples []int16
	pos     int
}

// windowsStreamEngine provides continuous zero-allocation polyphonic audio streaming
// directly to the Windows sound driver via waveOut with circular triple/quad buffering.
type windowsStreamEngine struct {
	mu         sync.Mutex
	hWaveOut   uintptr
	hEvent     uintptr
	rawBuffers [numRingBuffers][]byte
	headers    [numRingBuffers]waveHdr
	voices     []activeVoice

	stopCh    chan struct{}
	doneCh    chan struct{}
	closeOnce sync.Once
}

func newAudioEngine() audioEngine {
	engine, err := newWindowsStreamEngine()
	if err != nil {
		return &stubAudioEngine{}
	}
	return engine
}

func newWindowsStreamEngine() (*windowsStreamEngine, error) {
	if err := procWaveOutOpen.Find(); err != nil {
		return nil, err
	}
	if err := procCreateEventW.Find(); err != nil {
		return nil, err
	}

	sampleRate := uint32(SampleRate)
	channels := uint16(2)
	bits := uint16(16)
	align := channels * bits / 8
	wfx := waveFormatEx{
		wFormatTag:      waveFormatPCM,
		nChannels:       channels,
		nSamplesPerSec:  sampleRate,
		nAvgBytesPerSec: sampleRate * uint32(align),
		nBlockAlign:     align,
		wBitsPerSample:  bits,
		cbSize:          0,
	}

	hEvent, _, _ := procCreateEventW.Call(0, 0, 0, 0)
	if hEvent == 0 {
		return nil, errors.New("failed to create audio event handle")
	}

	var hWaveOut uintptr
	r, _, err := procWaveOutOpen.Call(
		uintptr(unsafe.Pointer(&hWaveOut)),
		waveMapper,
		uintptr(unsafe.Pointer(&wfx)),
		hEvent,
		0,
		callbackEvent,
	)
	if r != 0 {
		procCloseHandle.Call(hEvent)
		return nil, err
	}

	e := &windowsStreamEngine{
		hWaveOut: hWaveOut,
		hEvent:   hEvent,
		voices:   make([]activeVoice, 0, maxVoices),
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}

	for i := 0; i < numRingBuffers; i++ {
		e.rawBuffers[i] = make([]byte, bufSizeBytes)
		e.headers[i] = waveHdr{
			lpData:         uintptr(unsafe.Pointer(&e.rawBuffers[i][0])),
			dwBufferLength: uint32(bufSizeBytes),
		}
		procWaveOutPrepare.Call(hWaveOut, uintptr(unsafe.Pointer(&e.headers[i])), unsafe.Sizeof(e.headers[i]))
		procWaveOutWrite.Call(hWaveOut, uintptr(unsafe.Pointer(&e.headers[i])), unsafe.Sizeof(e.headers[i]))
	}

	go e.audioLoop()

	return e, nil
}

func (e *windowsStreamEngine) audioLoop() {
	defer close(e.doneCh)

	for {
		select {
		case <-e.stopCh:
			return
		default:
		}

		// Wait for next buffer ready or timeout (50ms)
		procWaitForSingleObject.Call(e.hEvent, 50)

		select {
		case <-e.stopCh:
			return
		default:
		}

		for i := 0; i < numRingBuffers; i++ {
			if e.headers[i].dwFlags&whdrDone != 0 {
				e.mixBuffer(e.rawBuffers[i])
				e.headers[i].dwFlags &^= whdrDone
				procWaveOutWrite.Call(e.hWaveOut, uintptr(unsafe.Pointer(&e.headers[i])), unsafe.Sizeof(e.headers[i]))
			}
		}
	}
}

func (e *windowsStreamEngine) mixBuffer(buf []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.voices) == 0 {
		for i := range buf {
			buf[i] = 0
		}
		return
	}

	for f := 0; f < framesPerBuf; f++ {
		var sumL, sumR float64
		for vi := 0; vi < len(e.voices); vi++ {
			v := &e.voices[vi]
			if v.pos < len(v.samples) {
				sumL += float64(v.samples[v.pos])
				sumR += float64(v.samples[v.pos+1])
				v.pos += 2
			}
		}

		// Soft limiter saturation to avoid harsh digital clipping
		left := softLimit(sumL)
		right := softLimit(sumR)

		offset := f * 4
		buf[offset] = byte(left)
		buf[offset+1] = byte(left >> 8)
		buf[offset+2] = byte(right)
		buf[offset+3] = byte(right >> 8)
	}

	// Prune completed voices
	n := 0
	for _, v := range e.voices {
		if v.pos < len(v.samples) {
			e.voices[n] = v
			n++
		}
	}
	e.voices = e.voices[:n]
}

func softLimit(x float64) int16 {
	if x > 32767.0 {
		return 32767
	}
	if x < -32768.0 {
		return -32768
	}
	return int16(x)
}

func (e *windowsStreamEngine) PlaySamples(samples []int16) {
	if len(samples) == 0 {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.voices) >= maxVoices {
		copy(e.voices, e.voices[1:])
		e.voices = e.voices[:len(e.voices)-1]
	}
	e.voices = append(e.voices, activeVoice{samples: samples, pos: 0})
}

func (e *windowsStreamEngine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.voices = e.voices[:0]
}

func (e *windowsStreamEngine) Close() {
	e.closeOnce.Do(func() {
		close(e.stopCh)
		<-e.doneCh

		procWaveOutReset.Call(e.hWaveOut)
		for i := 0; i < numRingBuffers; i++ {
			procWaveOutUnprepare.Call(e.hWaveOut, uintptr(unsafe.Pointer(&e.headers[i])), unsafe.Sizeof(e.headers[i]))
		}
		procWaveOutClose.Call(e.hWaveOut)
		procCloseHandle.Call(e.hEvent)
	})
}

// playWAVMemory plays in-memory WAV audio using Windows winmm.dll PlaySoundW with SND_NOWAIT.
func playWAVMemory(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	defer func() {
		_ = recover()
	}()

	if err := procPlaySound.Find(); err != nil {
		return err
	}

	flags := uintptr(sndMemory | sndAsync | sndNodefault | sndNowait)
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
