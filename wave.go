package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gopxl/beep"
)

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

func GlideSineWave(startFreq, endFreq float64, dur time.Duration) beep.Streamer {
	totalSamples := SAMPLERATE.N(dur)
	var sampleIndex int
	phase := .0

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

			// Compute phase increment per sample
			inc := 2 * math.Pi * freq / float64(SAMPLERATE)

			val := math.Sin(phase)
			samples[i][0], samples[i][1] = val, val

			phase += inc
			sampleIndex++
		}

		return len(samples), true
	})
}
