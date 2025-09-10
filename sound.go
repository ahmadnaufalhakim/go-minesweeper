package main

import (
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

const SAMPLERATE = beep.SampleRate(44100)

var (
	mixer = &beep.Mixer{}
)

func InitSoundSystem() {
	speaker.Init(SAMPLERATE, SAMPLERATE.N(time.Second/10))
	speaker.Play(mixer)
}

func PlaySound(s beep.Streamer) {
	mixer.Add(s)
}

func BombSound() beep.Streamer {
	return NoiseWave(200 * time.Millisecond)
}

func CellClearSound() beep.Streamer {
	return GlideSineWave(220, 880, 200*time.Millisecond)
}

func WinSound() beep.Streamer {
	return GlideSineWave(880, 220, 200*time.Millisecond)
}
