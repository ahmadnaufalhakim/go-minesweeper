package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type GameOptions struct {
	Style            tcell.Style
	BorderStyle      BorderStyle
	ShowInnerBorders bool
	Background       string
	Difficulty       DifficultyConfig
}

func NewGameOptions() *GameOptions {
	return &GameOptions{
		Style:            DefaultStyle,
		BorderStyle:      DefaultBorder,
		ShowInnerBorders: false,
		Background:       "none",
		Difficulty:       DifficultyMap["beginner"],
		//TODO: debug for `ShowInnerBorders = true`
	}
}

func RunGame(screen tcell.Screen, m *Minesweeper, opts *GameOptions, ng bool) GameState {
	w, h := screen.Size()
	mScreenX := (w-(m.Cols+2))/2 - (m.Cols+2)%2
	mScreenY := (h-(m.Rows+2))/2 - (m.Rows+2)%2
	var err error

	screen.EnableMouse(tcell.MouseButtonEvents, tcell.MouseDragEvents)
	screen.EnablePaste()

	playing := true
	ox, oy := -1, -1
	var lastMouseButtons tcell.ButtonMask
	for playing {
		screen.Clear()
		DrawBackground(screen, opts.Background, m.IsGameOver && !m.IsWon)
		m.Draw(screen, opts.BorderStyle, opts.ShowInnerBorders, mScreenX, mScreenY)
		if m.IsGameOver {
			var message string
			if m.IsWon {
				message = "You win!"
				DrawCentered(screen, mScreenY-3, opts.Style, "😎")
			} else {
				message = "You lose!"
				DrawCentered(screen, mScreenY-3, opts.Style, "😭")
			}
			DrawCentered(screen, mScreenY-2, opts.Style, message)
			DrawCentered(screen, mScreenY-1, opts.Style, "Press 'r' to create a new board, 'q' to quit to main menu.")
		} else {
			DrawCentered(screen, mScreenY-3, opts.Style, "🙂")
		}
		screen.Show()

		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEsc:
				return StateMenu
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'q':
					playing = false
				case 'r':
					m, err = GenerateBoard(opts.Difficulty)
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
					row, col, ok := m.ScreenToGrid(x, y, mScreenX, mScreenY, opts.ShowInnerBorders)
					if ok {
						switch lastMouseButtons {
						case tcell.Button1:
							if ok := m.Reveal(row, col, true); ok {
								if m.IsGameOver {
									if m.IsWon {
										PlaySound(WinSound())
									} else {
										PlaySound(BombSound())
									}
								} else {
									PlaySound(CellClearSound())
								}
							}
						case tcell.Button2:
							m.Flag(row, col)
						}
					}
					ox, oy = -1, -1
				}
			}
		}
	}

	return StateMenu
}
