package main

import (
	"context"
	"log"

	"github.com/gdamore/tcell/v2"
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateQuit
	gameStateCount
)

const (
	TRIES              = 1500
	MAX_COMPONENT_SIZE = 18
)

var eventCh chan tcell.Event

func main() {
	// Initialize screen
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

	eventCh = make(chan tcell.Event, 128)
	go func() {
		for {
			ev := screen.PollEvent()
			eventCh <- ev
		}
	}()

	screen.SetStyle(DefaultStyle)

	gameOptions := NewGameOptions()
	InitSoundSystem(gameOptions)

	for {
		state, gameOptions, cfg, ng := RunMenu(screen, gameOptions)
		if state == StateQuit {
			break
		}
		if state == StatePlaying {
			var minesweeper *Minesweeper
			if ng {
				// Create a cancellable context for NG board generation.
				// cancel() can be called explicitly (when user presses
				// 'q') or at the end (via defer).
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel() // always call cancel eventually (avoid context leak)

				// Channel to receive the NG board generation result
				doneCh := make(chan *Minesweeper, 1)

				// Run NG board generation in a goroutine
				go func() {
					doneCh <- WaitForNGBoard(ctx, screen, cfg)
				}()

				generating := true
				for generating {
					select {
					case ev := <-eventCh:
						switch ev := ev.(type) {
						case *tcell.EventKey:
							// User cancels NG board generation with 'q' or Esc
							if (ev.Key() == tcell.KeyRune && ev.Rune() == 'q') || ev.Key() == tcell.KeyEsc {
								cancel() // triggers ctx.Done() inside WaitForNGBoard
							}
						case *tcell.EventResize:
							screen.Sync()
						}

					// NG board generation finishes (either success OR failed)
					case m := <-doneCh:
						minesweeper = m
						generating = false
					}
				}

				// If NG board generation is cancelled, go back to main menu
				if minesweeper == nil {
					continue
				}
			} else {
				minesweeper, err = GenerateBoardWithStartCell(cfg)
				if err != nil {
					log.Fatal(err)
				}
			}

			RunGame(screen, minesweeper, gameOptions, ng)
		}
	}
}
