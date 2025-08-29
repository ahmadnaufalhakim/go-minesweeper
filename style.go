package main

import "github.com/gdamore/tcell/v2"

type BorderStyle int

const (
	BorderThin BorderStyle = iota
	BorderThick
)

var DefaultBorder = BorderThin
var DefaultBorderStyle = tcell.StyleDefault.Background(COLOR_DARKGRAY).Foreground(tcell.ColorBlack)

var (
	COLOR_DARKGRAY  = tcell.NewRGBColor(96, 96, 96)
	COLOR_LIGHTGRAY = tcell.NewRGBColor(200, 200, 200)
	COLOR_BOMB      = tcell.ColorWhite
	COLOR_ONE       = tcell.ColorBlue
	COLOR_TWO       = tcell.NewRGBColor(0, 127, 0)
	COLOR_THREE     = tcell.ColorRed
	COLOR_FOUR      = tcell.NewRGBColor(0, 0, 127)
	COLOR_FIVE      = tcell.NewRGBColor(127, 0, 0)
	COLOR_SIX       = tcell.NewRGBColor(0, 127, 127)
	COLOR_SEVEN     = tcell.ColorBlack
	COLOR_EIGHT     = tcell.NewRGBColor(128, 128, 128)
)
var ValueToCellStyle = map[int]tcell.Style{
	CLEAR: tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(tcell.ColorReset).Bold(true),
	BOMB:  tcell.StyleDefault.Background(tcell.ColorRed).Foreground(COLOR_BOMB).Bold(true),
	1:     tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(COLOR_ONE).Bold(true),
	2:     tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(COLOR_TWO).Bold(true),
	3:     tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(COLOR_THREE).Bold(true),
	4:     tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(COLOR_FOUR).Bold(true),
	5:     tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(COLOR_FIVE).Bold(true),
	6:     tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(COLOR_SIX).Bold(true),
	7:     tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(COLOR_SEVEN).Bold(true),
	8:     tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(COLOR_EIGHT).Bold(true),
}

var DefaultStyle = tcell.StyleDefault.Background(COLOR_LIGHTGRAY).Foreground(tcell.ColorBlack)
var FlagStyle = tcell.StyleDefault.Background(tcell.ColorOrange).Foreground(tcell.ColorDarkRed)
