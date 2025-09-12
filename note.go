package main

import "math"

var (
	C4  = Note(-9)
	Cs4 = Note(-8)
	D4  = Note(-7)
	Ds4 = Note(-6)
	E4  = Note(-5)
	F4  = Note(-4)
	Fs4 = Note(-3)
	G4  = Note(-2)
	Gs4 = Note(-1)
	A4  = Note(0)
	As4 = Note(1)
	B4  = Note(2)

	C5  = Note(3)
	Cs5 = Note(4)
	D5  = Note(5)
	Ds5 = Note(6)
	E5  = Note(7)
	F5  = Note(8)
	Fs5 = Note(9)
	G5  = Note(10)
	Gs5 = Note(11)
	A5  = Note(12)
	As5 = Note(13)
	B5  = Note(14)

	C6  = Note(15)
	Cs6 = Note(16)
	D6  = Note(17)
	Ds6 = Note(18)
	E6  = Note(19)
	F6  = Note(20)
	Fs6 = Note(21)
	G6  = Note(22)
	Gs6 = Note(23)
	A6  = Note(24)
	As6 = Note(25)
	B6  = Note(26)
)

func Note(semitonesFromA4 int) float64 {
	return 440 * math.Pow(2, float64(semitonesFromA4)/12.0)
}
