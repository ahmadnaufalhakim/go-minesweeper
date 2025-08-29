package main

import "github.com/gdamore/tcell/v2"

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
	start := (w-len(str))/2 - len(str)%2
	DrawString(screen, start, y, DefaultStyle, str)
}

func (s *Sprite) Draw(screen tcell.Screen, style tcell.Style) {
	screen.SetContent(s.X, s.Y, s.Char, nil, style)
}
