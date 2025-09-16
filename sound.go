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
	FORMAT = beep.Format{
		SampleRate:  SAMPLERATE,
		NumChannels: 2,
		Precision:   2,
	}

	mixer      = &beep.Mixer{}
	volumeCtrl = &effects.Volume{
		Streamer: mixer,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	sounds = make(map[string]func() beep.Streamer)
)

func BackgroundLoop() beep.Streamer {
	// 1 bar duration
	barDur := 1500 * time.Millisecond

	// Chord
	chordCycleFactory := func() beep.Streamer {
		return Phrase(
			ModChordWave(C5, []int{0, 4}, barDur, .5, 4, .25, 2),
			ModChordWave(A4, []int{0, 3}, barDur, .5, 4, .25, 2),
			ModChordWave(F4, []int{0, 4}, barDur, .5, 4, .25, 2),
			ModChordWave(G4, []int{0, 4}, barDur, .5, 4, .25, 2),
		)
	}
	chordPart := beep.Seq(chordCycleFactory(), chordCycleFactory())

	// Bass
	baseCycleFactory := func() beep.Streamer {
		return Phrase(
			ModSineWave(E4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(E4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(E4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(E4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(E4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(D4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),

			ModSineWave(C4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(C4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(C4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(C4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(C4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(B3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),

			ModSineWave(A3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(A3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(A3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(A3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(A3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(A3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),

			ModSineWave(B3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(B3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(B3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(B3, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(C4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
			ModSineWave(D4, barDur/(6*2), 0, 0, .5, 4), Rest(barDur/(6*2)),
		)
	}
	bassPart := beep.Seq(baseCycleFactory(), baseCycleFactory())

	// Melody
	dtype := HardClip
	distortParam := 5.0
	melodyCycleOne := Phrase(
		Distort(SineWave(E5, barDur/3), dtype, distortParam),
		Distort(SineWave(C5, barDur/3), dtype, distortParam),
		Rest(barDur/3),

		Rest(barDur),

		Rest(2*barDur/3),
		Distort(SineWave(G5, barDur/6), dtype, distortParam),
		Rest(barDur/6),

		Distort(SineWave(A5, barDur/6), dtype, distortParam),
		Rest(barDur/6),
		Distort(SineWave(G5, barDur/6), dtype, distortParam),
		Distort(SineWave(F5, barDur/6), dtype, distortParam),
		Distort(SineWave(E5, barDur/6), dtype, distortParam),
		Rest(barDur/6),
	)
	melodyCycleTwo := Phrase(
		Distort(SineWave(D5, barDur/3), dtype, distortParam),
		Distort(SineWave(E5, barDur/6), dtype, distortParam),
		Distort(SineWave(C5, barDur/3), dtype, distortParam),
		Rest(barDur/6),

		Rest(barDur),

		Rest(barDur),

		Distort(SineWave(A4, barDur/6), dtype, distortParam),
		Distort(SineWave(G4, barDur/6), dtype, distortParam),
		Distort(SineWave(B4, barDur/6), dtype, distortParam),
		Distort(SineWave(G4, barDur/6), dtype, distortParam),
		Distort(SineWave(C5, barDur/6), dtype, distortParam),
		Distort(SineWave(D5, barDur/6), dtype, distortParam),
	)
	melodyPart := beep.Seq(melodyCycleOne, melodyCycleTwo)

	mixed := beep.Mix(chordPart, bassPart, melodyPart)
	volumeScale := 1.0 / math.Sqrt(3.0)
	scaled := &effects.Volume{
		Streamer: mixed,
		Base:     2,
		Volume:   math.Log2(volumeScale),
		Silent:   false,
	}

	buf := beep.NewBuffer(FORMAT)
	buf.Append(scaled)

	// Loop the phrase forever
	return beep.Loop(-1, buf.Streamer(0, buf.Len()))
}

func LoadSounds() {
	sounds["intro"] = func() beep.Streamer {
		return BackgroundLoop()
	}
	sounds["bomb"] = func() beep.Streamer {
		buf := beep.NewBuffer(FORMAT)
		buf.Append(NoiseWave(150 * time.Millisecond))
		return buf.Streamer(0, buf.Len())
	}
	sounds["cellClear"] = func() beep.Streamer {
		buf := beep.NewBuffer(FORMAT)
		buf.Append(GlideSineWave(220, 880, 100*time.Millisecond))
		return buf.Streamer(0, buf.Len())
	}
	sounds["win"] = func() beep.Streamer {
		buf := beep.NewBuffer(FORMAT)
		buf.Append(Phrase(
			Distort(SineWave(C4, 150*time.Millisecond), SoftClip, 2),
			Distort(SineWave(E4, 150*time.Millisecond), SoftClip, 2),
			Distort(SineWave(G4, 150*time.Millisecond), SoftClip, 2),
			Distort(SineWave(C5, 350*time.Millisecond), SoftClip, 2),

			Rest(300*time.Millisecond),
			ModChordWave(C4, []int{0, 7}, 50*time.Millisecond, 1, 1.5, 1.125, 1),
			Rest(50*time.Millisecond),
			ModChordWave(C4, []int{0, 7}, 400*time.Millisecond, 1, 1.5, 1.125, 1),
		))
		return buf.Streamer(0, buf.Len())
	}
}

// Setting up speaker and sound buffers
func InitSoundSystem(opts *GameOptions) {
	// Set up volume
	SetVolume(opts.Volume)

	// Set up the speaker
	speaker.Init(SAMPLERATE, SAMPLERATE.N(time.Second/10))
	speaker.Play(volumeCtrl)
	LoadSounds()
}

// Play sound by adding new sound to the mixer
func PlaySound(name string) {
	if factory, ok := sounds[name]; ok {
		mixer.Add(factory())
	}
}

func StopAllSounds() {
	mixer.Clear()
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
