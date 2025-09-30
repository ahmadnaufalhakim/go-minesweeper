package main

import (
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
				minesweeper = WaitForNGBoard(screen, cfg)
			} else {
				minesweeper, _ = GenerateBoardWithStartCell(cfg)
			}
			RunGame(screen, minesweeper, gameOptions, ng)
		}
	}
}
