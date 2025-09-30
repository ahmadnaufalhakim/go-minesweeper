package main

import (
	"slices"

	"github.com/ahmadnaufalhakim/go-minesweeper/assets"
	"github.com/gdamore/tcell/v2"
)

type Sprite struct {
	Char rune
	X, Y int
}

const (
	DEFAULT_MARGIN_X = 2
	DEFAULT_MARGIN_Y = 1
)

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

func DrawFrame(screen tcell.Screen, x, y int, style tcell.Style, strs []string, marginX, marginY int) {
	if marginX <= 0 {
		marginX = DEFAULT_MARGIN_X
	}
	if marginY <= 0 {
		marginY = DEFAULT_MARGIN_Y
	}
	strLens := make([]int, len(strs))
	for i := range strs {
		strLens[i] = len(strs[i])
	}
	maxStrLen := slices.Max(strLens)

	// Draw frame borders
	NewSprite('╔', x, y).Draw(screen, style)
	NewSprite('╗', x+2*marginX+maxStrLen+1, y).Draw(screen, style)
	NewSprite('╚', x, y+2*marginY+len(strs)+1).Draw(screen, style)
	NewSprite('╝', x+2*marginX+maxStrLen+1, y+2*marginY+len(strs)+1).Draw(screen, style)
	for j := x + 1; j < x+1+2*marginX+maxStrLen; j++ {
		NewSprite('═', j, y).Draw(screen, style)
		NewSprite('═', j, y+2*marginY+len(strs)+1).Draw(screen, style)
	}
	for i := y + 1; i < y+1+2*marginY+len(strs); i++ {
		NewSprite('║', x, i).Draw(screen, style)
		NewSprite('║', x+2*marginX+maxStrLen+1, i).Draw(screen, style)
	}
	// Draw frame fillings
	for i := y + 1; i < y+1+2*marginY+len(strs); i++ {
		for j := x + 1; j < x+1+2*marginX+maxStrLen; j++ {
			NewSprite(' ', j, i).Draw(screen, style)
		}
	}
}

func DrawOverlay(screen tcell.Screen, style tcell.Style, strs []string, marginX, marginY int) {
	w, h := screen.Size()
	strLens := make([]int, len(strs))
	for i := range strs {
		strLens[i] = len(strs[i])
	}
	maxStrLen := slices.Max(strLens)
	x := (w-maxStrLen)/2 - maxStrLen%2
	y := (h-len(strs))/2 - len(strs)%2

	DrawFrame(screen, x-1-marginX, y-1-marginY, style, strs, marginX, marginY)
	for i := range len(strs) {
		DrawCentered(screen, y+i, style, strs[i])
	}
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
