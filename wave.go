package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/generators"
)

type DistortionType int

const (
	HardClip DistortionType = iota
	SoftClip
	BitCrush
	distortionTypeCount
)

// Generate a random wave sample (bzz sound)
func NoiseWave(dur time.Duration) beep.Streamer {
	totalSamples := SAMPLERATE.N(dur)
	var sampleIndex int

	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			// Stop the stream when finished
			if sampleIndex >= totalSamples {
				return i, false
			}

			// Randomize samples [-1,1]
			samples[i][0] = rand.Float64()*2 - 1
			samples[i][1] = rand.Float64()*2 - 1

			sampleIndex++
		}

		return len(samples), true
	})
}

// Generates a pure sine wave with given frequency and duration
func SineWave(freq float64, dur time.Duration) beep.Streamer {
	totalSamples := SAMPLERATE.N(dur)
	var sampleIndex int

	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			// Stop the stream when finished
			if sampleIndex >= totalSamples {
				return i, false
			}

			t := float64(sampleIndex) / float64(SAMPLERATE)
			val := math.Sin(2 * math.Pi * freq * t)
			samples[i][0], samples[i][1] = val, val

			sampleIndex++
		}

		return len(samples), true
	})
}

// Generates a sine wave with optional vibrato (frequency modulation)
// and tremolo (amplitude modulation).
//
// Vibrato: adds `vibratoDepth * sin(2π * vibratoRate * t)` to
// the baseFreq
//
// Tremolo: multiplies `1 + tremoloDepth * sin(2π * tremoloRate * t)` to
// the baseFreq
func ModSineWave(
	baseFreq float64, dur time.Duration,
	vibratoDepth, vibratoRate float64,
	tremoloDepth, tremoloRate float64,
) beep.Streamer {
	totalSamples := SAMPLERATE.N(dur)
	var sampleIndex int

	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			// Stop the stream when finished
			if sampleIndex >= totalSamples {
				return i, false
			}

			t := float64(sampleIndex) / float64(SAMPLERATE)

			// Vibrato (frequency modulation)
			freqMod := vibratoDepth * math.Sin(2*math.Pi*vibratoRate*t)
			freq := baseFreq + freqMod

			// Tremolo (amplitude modulation)
			ampMod := 1 + tremoloDepth*math.Sin(2*math.Pi*tremoloRate*t)

			val := math.Sin(2*math.Pi*freq*t) * ampMod
			samples[i][0], samples[i][1] = val, val

			sampleIndex++
		}

		return len(samples), true
	})
}

// Generates a gliding sine wave from `startFreq` to `endFreq`
func GlideSineWave(startFreq, endFreq float64, dur time.Duration) beep.Streamer {
	totalSamples := SAMPLERATE.N(dur)
	var sampleIndex int

	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			// Stop the stream when finished
			if sampleIndex >= totalSamples {
				return i, false
			}

			// Progress 0..1 over the glide duration
			progress := float64(sampleIndex) / float64(totalSamples)

			// Interpolate frequency linearly
			freq := startFreq + (endFreq-startFreq)*progress

			// Compute absolute time in seconds
			t := float64(sampleIndex) / float64(SAMPLERATE)

			val := math.Sin(2 * math.Pi * freq * t)
			samples[i][0], samples[i][1] = val, val

			sampleIndex++
		}

		return len(samples), true
	})
}

func clip(val float64) float64 {
	if val > 1 {
		return 1
	} else if val < -1 {
		return -1
	}
	return val
}

func Distort(source beep.Streamer, dtype DistortionType, param float64) beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		n, ok = source.Stream(samples)
		for i := 0; i < n; i++ {
			switch dtype {
			case HardClip:
				// param -> gain
				valL := samples[i][0] * param
				valR := samples[i][1] * param
				samples[i][0] = clip(valL)
				samples[i][1] = clip(valR)
			case SoftClip:
				// param -> gain
				samples[i][0] = math.Tanh(samples[i][0] * param)
				samples[i][1] = math.Tanh(samples[i][1] * param)
			case BitCrush:
				// param -> bit depth
				levels := math.Pow(2, param)
				samples[i][0] = math.Round(samples[i][0]*levels) / levels
				samples[i][1] = math.Round(samples[i][1]*levels) / levels
			}
		}

		return n, ok
	})
}

// Generates silence. That's it.
func Rest(dur time.Duration) beep.Streamer {
	return generators.Silence(SAMPLERATE.N(dur))
}

// Generates a chord relative to the root frequency.
//
// intervals = semitone offsets (e.g., [0,4,7] for major triad; [0,3,7] for minor triad)
func ChordWave(root float64, intervals []int, dur time.Duration) beep.Streamer {
	streamers := make([]beep.Streamer, len(intervals))
	for i, semitone := range intervals {
		freq := root * math.Pow(2, float64(semitone)/12.0)
		streamers[i] = SineWave(freq, dur)
	}

	mixed := beep.Mix(streamers...)
	volumeScale := 1.0 / math.Sqrt(float64(len(intervals)))
	return &effects.Volume{
		Streamer: mixed,
		Base:     2,
		Volume:   math.Log2(volumeScale),
		Silent:   false,
	}
}

// Generates a modulated chord relative to the root frequency.
func ModChordWave(
	root float64, intervals []int, dur time.Duration,
	vibratoDepth, vibratoRate float64,
	tremoloDepth, tremoloRate float64,
) beep.Streamer {
	streamers := make([]beep.Streamer, len(intervals))
	for i, semitone := range intervals {
		freq := root * math.Pow(2, float64(semitone)/12.0)
		streamers[i] = ModSineWave(
			freq, dur,
			vibratoDepth, vibratoRate,
			tremoloDepth, tremoloRate,
		)
	}

	mixed := beep.Mix(streamers...)
	volumeScale := 1.0 / math.Sqrt(float64(len(intervals)))
	return &effects.Volume{
		Streamer: mixed,
		Base:     2,
		Volume:   math.Log2(volumeScale),
		Silent:   false,
	}
}

// Generate a musical phrase, i.e. arpeggios, chords, combinations
func Phrase(streamers ...beep.Streamer) beep.Streamer {
	return beep.Seq(streamers...)
}
