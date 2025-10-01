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
	TRIES              = 1000
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
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				doneCh := make(chan *Minesweeper, 1)
				go func() {
					doneCh <- WaitForNGBoard(ctx, screen, cfg)
				}()

				generating := true
				for generating {
					select {
					case ev := <-eventCh:
						switch ev := ev.(type) {
						case *tcell.EventKey:
							if (ev.Key() == tcell.KeyRune && ev.Rune() == 'q') || ev.Key() == tcell.KeyEsc {
								cancel()
							}
						case *tcell.EventResize:
							screen.Sync()
						}
					case m := <-doneCh:
						minesweeper = m
						generating = false
					}
				}

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
