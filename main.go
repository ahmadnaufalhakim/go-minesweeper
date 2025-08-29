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

	for {
		state, gameOptions, cfg, ng := RunMenu(screen, gameOptions)
		if state == StateQuit {
			break
		}
		if state == StatePlaying {
			minesweeper, _ := GenerateBoard(cfg)
			RunGame(screen, minesweeper, gameOptions, ng)
		}
	}
}
