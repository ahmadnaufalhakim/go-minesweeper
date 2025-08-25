package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

type GameOptions struct {
	Style            tcell.Style
	BorderStyle      BorderStyle
	Control          string
	ShowInnerBorders bool
}

func initGameOptions() GameOptions {
	return GameOptions{
		Style:            DefaultStyle,
		BorderStyle:      DefaultBorder,
		Control:          "mouse",
		ShowInnerBorders: false,
		//TODO: debug for `ShowInnerBorders = true`
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	err = screen.Init()
	if err != nil {
		log.Fatal(err)
	}

	quit := func() {
		r := recover()
		if r != nil {
			screen.Fini()
			log.Panic(r)
		} else {
			screen.Fini()
		}
	}
	defer quit()

	screen.SetStyle(DefaultStyle)
	screen.EnableMouse(tcell.MouseButtonEvents, tcell.MouseDragEvents)
	screen.EnablePaste()
	screen.Clear()

	gameOptions := initGameOptions()
	minesweeper, err := GenerateBoard(16, 30, 99)
	if err != nil {
		log.Fatal(err)
	}

	running := true
	ox, oy := -1, -1
	var lastMouseButtons tcell.ButtonMask
	for running {
		// Draw
		screen.Clear()
		minesweeper.Draw(screen, gameOptions.BorderStyle, gameOptions.ShowInnerBorders, 10, 10)
		DrawString(screen, 10, 9, gameOptions.Style, fmt.Sprintf("%+v", minesweeper.IsGameOver))
		screen.Show()

		// Wait for event
		ev := screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			key := ev.Key()
			if key == tcell.KeyRune {
				switch ev.Rune() {
				case 'q':
					running = false
				case 'r':
					minesweeper, err = GenerateBoard(16, 30, 99)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			btn := ev.Buttons()

			switch btn {
			case tcell.Button1, tcell.Button2:
				if ox < 0 && oy < 0 {
					ox, oy = x, y
					lastMouseButtons = btn
				}
			case tcell.ButtonNone:
				if ox >= 0 {
					row, col, ok := minesweeper.ScreenToGrid(x, y, 10, 10, gameOptions.ShowInnerBorders)
					if ok {
						switch lastMouseButtons {
						case tcell.Button1:
							minesweeper.Reveal(row, col, true)
						case tcell.Button2:
							minesweeper.Flag(row, col)
						}
					}

					ox, oy = -1, -1
				}
			}
		}
	}
}
