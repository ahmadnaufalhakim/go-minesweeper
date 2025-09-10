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

func InitSoundSystem() {
	// Set up the speaker
	speaker.Init(SAMPLERATE, SAMPLERATE.N(time.Second/10))
	speaker.Play(mixer)

	// Pre-render each sound into a buffer
	bombBuf = beep.NewBuffer(format)
	bombBuf.Append(NoiseWave(200 * time.Millisecond))

	cellClearBuf = beep.NewBuffer(format)
	cellClearBuf.Append(GlideSineWave(220, 880, 150*time.Millisecond))

	winBuf = beep.NewBuffer(format)
	winBuf.Append(ChordWave(440, []int{3, 7, 10}, 150*time.Millisecond))
}

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
