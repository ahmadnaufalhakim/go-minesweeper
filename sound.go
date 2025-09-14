package main

import (
	"math"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
)

const SAMPLERATE = beep.SampleRate(44100)

var (
	mixer      = &beep.Mixer{}
	volumeCtrl = &effects.Volume{
		Streamer: mixer,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	format = beep.Format{
		SampleRate:  SAMPLERATE,
		NumChannels: 2,
		Precision:   2,
	}
	sounds = make(map[string]*beep.Buffer)
)

func LoadSounds() {
	sounds["bomb"] = beep.NewBuffer(format)
	sounds["cellClear"] = beep.NewBuffer(format)
	sounds["win"] = beep.NewBuffer(format)

	sounds["bomb"].Append(NoiseWave(150 * time.Millisecond))
	sounds["cellClear"].Append(GlideSineWave(220, 880, 100*time.Millisecond))
	sounds["win"].Append(Phrase(
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
	))
}

// Setting up speaker and sound buffers
func InitSoundSystem() {
	// Set up the speaker
	speaker.Init(SAMPLERATE, SAMPLERATE.N(time.Second/10))
	speaker.Play(volumeCtrl)
	LoadSounds()
}

// Play sound by adding new sound to the mixer
func PlaySound(name string) {
	if buf, ok := sounds[name]; ok {
		mixer.Add(buf.Streamer(0, buf.Len()))
	}
}

func SetVolume(percent int) {
	// If percentage set to 0,
	// mute the volume controller
	if percent <= 0 {
		volumeCtrl.Silent = true
		return
	}
	volumeCtrl.Silent = false

	// Convert percentage to volume
	// 0% -> mute, 100% -> 0 dB
	vol := float64(percent) / 100.0
	volumeCtrl.Volume = 2 * math.Log2(vol)
}
