package main

import (
	"github.com/ahmadnaufalhakim/go-minesweeper/assets"
	"github.com/gdamore/tcell/v2"
)

type Sprite struct {
	Char rune
	X, Y int
}

func NewSprite(char rune, x, y int) *Sprite {
	return &Sprite{
		Char: char,
		X:    x,
		Y:    y,
	}
}

func DrawString(screen tcell.Screen, x, y int, style tcell.Style, s string) {
	for i, c := range s {
		screen.SetContent(x+i, y, c, nil, style)
	}
}

func DrawCentered(screen tcell.Screen, y int, style tcell.Style, str string) {
	w, _ := screen.Size()
	x := (w-len(str))/2 - len(str)%2
	DrawString(screen, x, y, style, str)
}

func (s *Sprite) Draw(screen tcell.Screen, style tcell.Style) {
	screen.SetContent(s.X, s.Y, s.Char, nil, style)
}

func DrawBackground(screen tcell.Screen, name string, negative bool) {
	w, h := screen.Size()
	bg, ok := assets.LoadBackground(name)
	if !ok {
		return
	}
	for y := range h {
		for x := range w {
			xRef := int(float64(x * bg.Bounds().Max.X / w))
			yRef := int(float64(y * bg.Bounds().Max.Y / h))
			rRef, gRef, bRef, _ := bg.At(xRef, yRef).RGBA()

			r := int32(rRef >> 8)
			g := int32(gRef >> 8)
			b := int32(bRef >> 8)
			if negative {
				r = 255 - r
				g = 255 - g
				b = 255 - b
			}

			screen.SetContent(
				x, y,
				' ', nil,
				DefaultStyle.Background(tcell.NewRGBColor(r, g, b)))
		}
	}
}
