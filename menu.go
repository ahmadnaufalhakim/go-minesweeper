package main

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/ahmadnaufalhakim/go-minesweeper/assets"
	"github.com/gdamore/tcell/v2"
)

type MenuPage int

const (
	PageMain MenuPage = iota
	PageOptions
	PageCredits
	PageQuitConfirm
	PageCustomInput
)

func drawMainMenu(screen tcell.Screen, titleItems []string, selected int, difficulty string, opts *GameOptions) {
	_, h := screen.Size()

	titleHeight := len(titleItems)

	menuItems := []string{
		fmt.Sprintf("Play <%s>", difficulty),
		fmt.Sprintf("Play NG <%s>", difficulty),
		"Options",
		"Credits",
		"Quit",
	}
	menuHeight := (len(menuItems)+1)*2 - 1

	contentHeight := titleHeight + 4 + menuHeight
	cScreenY := (h-contentHeight)/2 - contentHeight%2

	for i, titleItem := range titleItems {
		DrawCentered(screen, cScreenY+i, opts.Style, titleItem)
	}

	DrawCentered(screen, cScreenY+titleHeight+4, opts.Style, "=== Main Menu ===")
	for i, menuItem := range menuItems {
		prefix := ""
		if i == selected {
			prefix = "> "
		}
		DrawCentered(screen, cScreenY+titleHeight+6+i*2, opts.Style, prefix+menuItem)
	}
}

func drawOptionsMenu(screen tcell.Screen, titleItems []string, selected int, opts *GameOptions) {
	_, h := screen.Size()

	titleHeight := len(titleItems)

	menuItems := []string{
		fmt.Sprintf("Show inner borders: %v", opts.ShowInnerBorders),
		fmt.Sprintf("Border style: %v", opts.BorderStyle),
		"Back",
	}
	menuHeight := (len(menuItems)+1)*2 - 1

	contentHeight := titleHeight + 4 + menuHeight
	cScreenY := (h-contentHeight)/2 - contentHeight%2

	for i, titleItem := range titleItems {
		DrawCentered(screen, cScreenY+i, opts.Style, titleItem)
	}

	DrawCentered(screen, cScreenY+titleHeight+4, opts.Style, "=== Options ===")
	for i, menuItem := range menuItems {
		prefix := ""
		if i == selected {
			prefix = "> "
		}
		DrawCentered(screen, cScreenY+titleHeight+6+i*2, opts.Style, prefix+menuItem)
	}
}

func drawCredits(screen tcell.Screen, titleItems []string, opts *GameOptions) {
	_, h := screen.Size()

	titleHeight := len(titleItems)

	menuItems := []string{
		"=== Credits ===",
		"made with ❤️by Ahmad Naufal Hakim 🤓",
		"[Esc] Back",
	}
	menuHeight := len(menuItems)*2 - 1

	contentHeight := titleHeight + 4 + menuHeight
	cScreenY := (h-contentHeight)/2 - contentHeight%2

	for i, titleItem := range titleItems {
		DrawCentered(screen, cScreenY+i, opts.Style, titleItem)
	}

	for i, menuItem := range menuItems {
		DrawCentered(screen, cScreenY+titleHeight+4+i*3, opts.Style, menuItem)
	}
}

func drawQuitConfirm(screen tcell.Screen, opts *GameOptions) {
	_, h := screen.Size()

	menuItems := []string{
		"[y] yeah, i'm washed chat 😭",
		"[n] nah, i'd win 😎",
	}

	contentHeight := 4
	cScreenY := (h-contentHeight)/2 - contentHeight%2

	DrawCentered(screen, cScreenY, opts.Style, "Are you sure you want to quit?")
	for i, menuItem := range menuItems {
		DrawCentered(screen, 2+cScreenY+i, opts.Style, menuItem)
	}
}

func drawCustomInput(
	screen tcell.Screen,
	selected int,
	cfg DifficultyConfig,
	buf string, errMsg string,
	opts *GameOptions,
) {
	menuItems := []string{
		fmt.Sprintf("Rows: %d", cfg.Rows),
		fmt.Sprintf("Cols: %d", cfg.Cols),
		fmt.Sprintf("BombCount: %d", cfg.BombCount),
		"Start",
		"Back",
	}
	DrawCentered(screen, 2, opts.Style, "=== Custom Difficulty ===")
	for i, menuItem := range menuItems {
		prefix := ""
		if i == selected {
			prefix = "> "
		}
		DrawCentered(screen, 4+i*2, opts.Style, prefix+menuItem)
	}

	if buf != "" {
		DrawCentered(screen, 16, opts.Style, fmt.Sprintf("Typing: %s", buf))
	}
	if errMsg != "" {
		DrawCentered(screen, 18, opts.Style, "Error: "+errMsg)
	}
}

func drawHelpHint(screen tcell.Screen, opts *GameOptions) {
	w, h := screen.Size()
	message := "W/S = up/down, A/D = change, Enter = select"
	DrawString(screen, w-len(message)-1, h-1, opts.Style, message)
}

func RunMenu(screen tcell.Screen, opts *GameOptions) (GameState, *GameOptions, DifficultyConfig, bool) {
	page := PageMain
	titleItems := assets.RandomTitle()
	selected := 0
	difficulties := []string{"beginner", "intermediate", "advanced", "expert", "insane", "custom"}
	diffIndex := 0
	playingNG := false
	customCfg := DifficultyConfig{Rows: 9, Cols: 9, BombCount: 10}
	fieldIndex := 0
	inputBuffer := ""
	errorMsg := ""

	var menuCount int

	for {
		screen.Clear()
		switch page {
		case PageMain:
			drawMainMenu(screen, titleItems, selected, difficulties[diffIndex], opts)
			menuCount = 5
		case PageOptions:
			drawOptionsMenu(screen, titleItems, selected, opts)
			menuCount = 3
		case PageCredits:
			drawCredits(screen, titleItems, opts)
		case PageQuitConfirm:
			drawQuitConfirm(screen, opts)
		case PageCustomInput:
			drawCustomInput(screen, fieldIndex, customCfg, inputBuffer, errorMsg, opts)
			menuCount = 5
		}
		drawHelpHint(screen, opts)
		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEsc:
				if page != PageMain {
					page = PageMain
					selected = 0
					inputBuffer = ""
					errorMsg = ""
				}
			case tcell.KeyUp:
				switch page {
				case PageMain, PageOptions:
					selected = (selected - 1 + menuCount) % menuCount
				case PageCustomInput:
					fieldIndex = (fieldIndex - 1 + menuCount) % menuCount
					inputBuffer = ""
				}
			case tcell.KeyDown:
				switch page {
				case PageMain, PageOptions:
					selected = (selected + 1) % menuCount
				case PageCustomInput:
					fieldIndex = (fieldIndex + 1) % menuCount
					inputBuffer = ""
				}
			case tcell.KeyLeft:
				if page == PageMain && (selected == 0 || selected == 1) {
					diffIndex = (diffIndex - 1 + len(difficulties)) % len(difficulties)
					opts.Difficulty = DifficultyMap[difficulties[diffIndex]]
				} else if page == PageOptions {
					switch selected {
					// Show inner borders
					case 0:
						opts.ShowInnerBorders = !opts.ShowInnerBorders
					// Border style
					case 1:
						opts.BorderStyle = (opts.BorderStyle - 1 + 2) % 2
					}
				}
			case tcell.KeyRight:
				if page == PageMain && (selected == 0 || selected == 1) {
					diffIndex = (diffIndex + 1) % len(difficulties)
					opts.Difficulty = DifficultyMap[difficulties[diffIndex]]
				} else if page == PageOptions {
					switch selected {
					// Show inner borders
					case 0:
						opts.ShowInnerBorders = !opts.ShowInnerBorders
					// Border style
					case 1:
						opts.BorderStyle = (opts.BorderStyle + 1) % 2
					}
				}
			case tcell.KeyEnter:
				switch page {
				case PageMain:
					switch selected {
					// Play
					case 0:
						playingNG = false
						if difficulties[diffIndex] == "custom" {
							page = PageCustomInput
						} else {
							opts.Difficulty = DifficultyMap[difficulties[diffIndex]]
							return StatePlaying, opts, DifficultyMap[difficulties[diffIndex]], playingNG
						}
					// Play NG
					case 1:
						playingNG = true
						if difficulties[diffIndex] == "custom" {
							page = PageCustomInput
						} else {
							opts.Difficulty = DifficultyMap[difficulties[diffIndex]]
							return StatePlaying, opts, DifficultyMap[difficulties[diffIndex]], playingNG
						}
					// Options
					case 2:
						page = PageOptions
						selected = 0
					// Credits
					case 3:
						page = PageCredits
					// Quit
					case 4:
						page = PageQuitConfirm
					}
				case PageOptions:
					if selected == 2 {
						page = PageMain
						selected = 0
					}
				case PageCustomInput:
					switch fieldIndex {
					// Start
					case 3:
						_, err := GenerateBoard(customCfg)
						if err != nil {
							errorMsg = err.Error()
						} else {
							opts.Difficulty = customCfg
							return StatePlaying, opts, customCfg, playingNG
						}
					// Back
					case 4:
						page = PageMain
						selected = 0
						inputBuffer = ""
						errorMsg = ""
					default:
						if inputBuffer != "" {
							val, err := strconv.Atoi(inputBuffer)
							if err == nil {
								switch fieldIndex {
								case 0:
									customCfg.Rows = val
								case 1:
									customCfg.Cols = val
								case 2:
									customCfg.BombCount = val
								}
							} else {
								errorMsg = err.Error()
							}
							inputBuffer = ""
						}
					}
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(inputBuffer) > 0 {
					inputBuffer = inputBuffer[:len(inputBuffer)-1]
				}
			case tcell.KeyRune:
				r := ev.Rune()
				if page == PageCustomInput && unicode.IsDigit(r) {
					inputBuffer += string(r)
				} else {
					switch r {
					case 'w':
						if page == PageMain || page == PageOptions {
							selected = (selected - 1 + menuCount) % menuCount
						}
					case 's':
						if page == PageMain || page == PageOptions {
							selected = (selected + 1) % menuCount
						}
					case 'a':
						if page == PageMain && (selected == 0 || selected == 1) {
							diffIndex = (diffIndex - 1 + len(difficulties)) % len(difficulties)
						} else if page == PageOptions {
							switch selected {
							// Show inner borders
							case 0:
								opts.ShowInnerBorders = !opts.ShowInnerBorders
							// Border style
							case 1:
								opts.BorderStyle = (opts.BorderStyle - 1 + 2) % 2
							}
						}
					case 'd':
						if page == PageMain && (selected == 0 || selected == 1) {
							diffIndex = (diffIndex + 1) % len(difficulties)
						} else if page == PageOptions {
							switch selected {
							// Show inner borders
							case 0:
								opts.ShowInnerBorders = !opts.ShowInnerBorders
							// Border style
							case 1:
								opts.BorderStyle = (opts.BorderStyle + 1) % 2
							}
						}
					case 'y':
						if page == PageQuitConfirm {
							return StateQuit, opts, DifficultyConfig{}, false
						}
					case 'n':
						if page == PageQuitConfirm {
							page = PageMain
							selected = 0
						}
					}
				}
			}
		}
	}
}
