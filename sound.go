package main

import (
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

const SAMPLERATE = beep.SampleRate(44100)

var (
	mixer  = &beep.Mixer{}
	format = beep.Format{
		SampleRate:  SAMPLERATE,
		NumChannels: 2,
		Precision:   2,
	}
	bombBuf      *beep.Buffer
	cellClearBuf *beep.Buffer
	winBuf       *beep.Buffer
)

// Setting up speaker and sound buffers
func InitSoundSystem() {
	// Set up the speaker
	speaker.Init(SAMPLERATE, SAMPLERATE.N(time.Second/10))
	speaker.Play(mixer)

	// Prepare sound buffers
	bombBuf = beep.NewBuffer(format)
	cellClearBuf = beep.NewBuffer(format)
	winBuf = beep.NewBuffer(format)

	// Define sound waves
	bombStreamer := NoiseWave(150 * time.Millisecond)
	cellClearStreamer := GlideSineWave(220, 880, 100*time.Millisecond)
	winStreamer := Phrase(
		ModSineWave(C4, 150*time.Millisecond, 2.5, 3, 0, 0),
		ModSineWave(E4, 150*time.Millisecond, 2.5, 3, 0, 0),
		ModSineWave(G4, 150*time.Millisecond, 2.5, 3, 0, 0),
		ModSineWave(C5, 350*time.Millisecond, 2.5, 3, 0, 0),

		// Distort(SineWave(C4, 150*time.Millisecond), HardClip, 2),
		// Distort(SineWave(E4, 150*time.Millisecond), HardClip, 2),
		// Distort(SineWave(G4, 150*time.Millisecond), HardClip, 2),
		// Distort(SineWave(C5, 350*time.Millisecond), HardClip, 2),

		// Distort(SineWave(C4, 150*time.Millisecond), SoftClip, 2),
		// Distort(SineWave(E4, 150*time.Millisecond), SoftClip, 2),
		// Distort(SineWave(G4, 150*time.Millisecond), SoftClip, 2),
		// Distort(SineWave(C5, 350*time.Millisecond), SoftClip, 2),

		// Distort(SineWave(C4, 150*time.Millisecond), BitCrush, 2),
		// Distort(SineWave(E4, 150*time.Millisecond), BitCrush, 2),
		// Distort(SineWave(G4, 150*time.Millisecond), BitCrush, 2),
		// Distort(SineWave(C5, 350*time.Millisecond), BitCrush, 2),

		Rest(300*time.Millisecond),
		ModChordWave(C4, []int{0, 7}, 50*time.Millisecond, 1, 1.5, 1.125, 1),
		Rest(50*time.Millisecond),
		ModChordWave(C4, []int{0, 7}, 400*time.Millisecond, 1, 1.5, 1.125, 1),
	)

	// Append sound waves to sound buffers
	bombBuf.Append(bombStreamer)
	cellClearBuf.Append(cellClearStreamer)
	winBuf.Append(winStreamer)
}

// Adding new sound to the mixer
func PlaySound(s beep.Streamer) {
	mixer.Add(s)
}

func BombSound() beep.Streamer {
	return bombBuf.Streamer(0, bombBuf.Len())
}

func CellClearSound() beep.Streamer {
	return cellClearBuf.Streamer(0, cellClearBuf.Len())
}

func WinSound() beep.Streamer {
	return winBuf.Streamer(0, winBuf.Len())
}
