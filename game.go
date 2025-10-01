package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

type GameOptions struct {
	Style            tcell.Style
	BorderStyle      BorderStyle
	ShowInnerBorders bool
	Background       string
	Volume           int
	Difficulty       DifficultyConfig

	bgIndex  int
	volIndex int
}

func NewGameOptions() *GameOptions {
	return &GameOptions{
		Style:            DefaultStyle,
		BorderStyle:      DefaultBorder,
		ShowInnerBorders: false,
		Background:       "none",
		Volume:           30,
		Difficulty:       DifficultyMap["beginner"],
		//TODO: debug for `ShowInnerBorders = true`

		bgIndex:  0,
		volIndex: 3,
	}
}

func (opts *GameOptions) ToggleInnerBorders() {
	opts.ShowInnerBorders = !opts.ShowInnerBorders
}

func (opts *GameOptions) NextBorderStyle(delta int) {
	opts.BorderStyle = BorderStyle((int(opts.BorderStyle) + delta + int(borderStyleCount)) % int(borderStyleCount))
}

func (opts *GameOptions) NextBackground(delta int, bgs []string) {
	opts.bgIndex = (opts.bgIndex + delta + len(bgs)) % len(bgs)
	opts.Background = bgs[opts.bgIndex]
}

func (opts *GameOptions) NextVolume(delta int, volPercentages []int) {
	opts.volIndex = (opts.volIndex + delta + len(volPercentages)) % len(volPercentages)
	opts.Volume = volPercentages[opts.volIndex]
	SetVolume(opts.Volume)
	PlaySound("cellClear")
}

func WaitForNGBoard(ctx context.Context, screen tcell.Screen, cfg DifficultyConfig) *Minesweeper {
	loadingMsg := "Generating NG board .."
	spinner := []rune{'|', '/', '-', '\\'}
	idx := 0
	attempt := 0

	minesweeperCh, progressCh := GenerateNGBoard(ctx, cfg, TRIES, MAX_COMPONENT_SIZE)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case minesweeper := <-minesweeperCh:
			if minesweeper == nil {
				screen.Clear()
				DrawOverlay(
					screen, FailedOverlayStyle,
					[]string{
						"Failed to generate NG board!😭",
						"Falling back to regular board ..",
					},
					DEFAULT_MARGIN_X, DEFAULT_MARGIN_Y,
				)
				screen.Show()
				time.Sleep(2000 * time.Millisecond)

				minesweeper, err := GenerateBoardWithStartCell(cfg)
				if err != nil {
					log.Fatal(err)
				}
				return minesweeper
			}

			screen.Clear()
			DrawOverlay(
				screen, SuccessOverlayStyle,
				[]string{
					"NG board successfully generated!😎",
					fmt.Sprintf("It only took %d attempts😤", attempt),
				},
				DEFAULT_MARGIN_X, DEFAULT_MARGIN_Y,
			)
			screen.Show()
			time.Sleep(2000 * time.Millisecond)
			return minesweeper

		case attempt = <-progressCh:

		case <-ticker.C:
			msgs := []string{
				fmt.Sprintf("%s %c", loadingMsg, spinner[idx%len(spinner)]),
				fmt.Sprintf("Attempt: %4d/%d", attempt, TRIES),
				"",
				"Press 'q' or 'Esc' key to cancel NG mode",
			}
			screen.Clear()
			DrawOverlay(screen, DefaultOverlayStyle, msgs, DEFAULT_MARGIN_X, DEFAULT_MARGIN_Y)
			screen.Show()
			idx++
		}
	}
}

func RunGame(screen tcell.Screen, m *Minesweeper, opts *GameOptions, ng bool) GameState {
	w, h := screen.Size()
	mScreenX := (w-(m.Cols+2))/2 - (m.Cols+2)%2
	mScreenY := (h-(m.Rows+2))/2 - (m.Rows+2)%2
	var err error

	screen.EnableMouse(tcell.MouseButtonEvents, tcell.MouseDragEvents)
	screen.EnablePaste()

	StopAllSounds()

	playing := true
	ox, oy := -1, -1
	var lastMouseButtons tcell.ButtonMask
	for playing {
		screen.Clear()
		DrawBackground(screen, opts.Background, m.IsGameOver && !m.IsWon)
		m.Draw(screen, opts.BorderStyle, opts.ShowInnerBorders, mScreenX, mScreenY)
		m.DrawSmiley(screen, mScreenY, opts.Style, lastMouseButtons)
		screen.Show()

		select {
		case ev := <-eventCh:
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
						StopAllSounds()
						if ng {
							ctx, cancel := context.WithCancel(context.Background())
							doneCh := make(chan *Minesweeper, 1)

							go func() {
								doneCh <- WaitForNGBoard(ctx, screen, opts.Difficulty)
							}()

							regenerating := true
							for regenerating {
								select {
								case newM := <-doneCh:
									cancel()
									m = newM
									regenerating = false
								case regEv := <-eventCh:
									switch regEv := regEv.(type) {
									case *tcell.EventKey:
										if regEv.Rune() == 'q' || regEv.Key() == tcell.KeyEsc {
											cancel()
											m = nil
											regenerating = false
										}
									}
								default:
								}
							}

							cancel()
							if m == nil {
								return StateMenu
							}
						} else {
							m, err = GenerateBoardWithStartCell(opts.Difficulty)
						}
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
											PlaySound("win")
										} else {
											PlaySound("bomb")
										}
									} else {
										PlaySound("cellClear")
									}
								}
							case tcell.Button2:
								m.Flag(row, col)
							}
						}
						ox, oy = -1, -1
						lastMouseButtons = tcell.ButtonNone
					}
				}
			}
		default:
		}
	}

	return StateMenu
}
