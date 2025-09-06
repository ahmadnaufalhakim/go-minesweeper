package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
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

func drawTitleItems(
	screen tcell.Screen,
	titleItems []string,
	offsetY int,
	opts *GameOptions,
) {
	for i, titleItem := range titleItems {
		DrawCentered(screen, offsetY+i, opts.Style, titleItem)
	}
}

func drawMenuItems(
	screen tcell.Screen,
	selected int, menuTitle string, menuItems []string,
	offsetY int,
	opts *GameOptions,
) {
	DrawCentered(screen, offsetY, opts.Style, menuTitle)
	for i, menuItem := range menuItems {
		menuStyle := DefaultStyle
		if selected >= 0 && i == selected {
			menuStyle = menuStyle.Background(tcell.ColorOrange)
		}
		DrawCentered(screen, offsetY+2+i*2, menuStyle, menuItem)
	}
}

func drawMainMenu(
	screen tcell.Screen, titleItems []string,
	selected int, difficulty string, difficultyNG string,
	opts *GameOptions,
) {
	w, h := screen.Size()

	titleHeight := len(titleItems)

	menuItems := []string{
		fmt.Sprintf("Play <%s>", strings.Repeat(" ", len(difficulty))),
		fmt.Sprintf("Play NG <%s>", strings.Repeat(" ", len(difficultyNG))),
		"Options",
		"Credits",
		"Quit",
	}
	menuHeight := (len(menuItems)+1)*2 - 1

	contentHeight := titleHeight + 4 + menuHeight
	titleOffsetY := (h-contentHeight)/2 - contentHeight%2

	drawTitleItems(screen, titleItems, titleOffsetY, opts)
	drawMenuItems(screen, selected, "⚑⚑⚑ Main Menu  ⚑⚑⚑", menuItems, titleOffsetY+titleHeight+4, opts)

	x := ((w-len(menuItems[0]))/2 - len(menuItems[0])%2) + (len(menuItems[0]) - (len(difficulty) + 1))
	y := titleOffsetY + titleHeight + 6
	DrawString(screen, x, y, DifficultyToStyle[difficulty], difficulty)
	xNG := ((w-len(menuItems[1]))/2 - len(menuItems[1])%2) + (len(menuItems[1]) - (len(difficultyNG) + 1))
	yNG := titleOffsetY + titleHeight + 8
	DrawString(screen, xNG, yNG, DifficultyToStyle[difficultyNG], difficultyNG)
}

func drawOptionsMenu(screen tcell.Screen, titleItems []string, selected int, opts *GameOptions) {
	_, h := screen.Size()

	titleHeight := len(titleItems)

	menuItems := []string{
		fmt.Sprintf("Show inner borders: <%v>", opts.ShowInnerBorders),
		fmt.Sprintf("Border style: <%v>", opts.BorderStyle),
		fmt.Sprintf("Background: <%v>", opts.Background),
		"Back",
	}
	menuHeight := (len(menuItems)+1)*2 - 1

	contentHeight := titleHeight + 4 + menuHeight
	titleOffsetY := (h-contentHeight)/2 - contentHeight%2

	drawTitleItems(screen, titleItems, titleOffsetY, opts)
	drawMenuItems(screen, selected, "⚑⚑⚑ Options  ⚑⚑⚑", menuItems, titleOffsetY+titleHeight+4, opts)
}

func drawCredits(screen tcell.Screen, titleItems []string, opts *GameOptions) {
	_, h := screen.Size()

	titleHeight := len(titleItems)

	menuItems := []string{
		"made with 💖",
		"by Ahmad Naufal Hakim 🤓",
	}
	menuHeight := len(menuItems)*2 - 1

	contentHeight := titleHeight + 4 + menuHeight
	titleOffsetY := (h-contentHeight)/2 - contentHeight%2

	drawTitleItems(screen, titleItems, titleOffsetY, opts)
	drawMenuItems(screen, -1, "⚑⚑⚑ Credits  ⚑⚑⚑", menuItems, titleOffsetY+titleHeight+2, opts)
}

func drawQuitConfirm(screen tcell.Screen, opts *GameOptions) {
	_, h := screen.Size()

	menuItems := []string{
		"[y] yeah, i'm washed chat 😭",
		"[n] nah, i'd win 😎",
	}

	contentHeight := 4
	contentOffsetY := (h-contentHeight)/2 - contentHeight%2

	DrawCentered(screen, contentOffsetY, opts.Style, "Are you sure you want to quit?")
	for i, menuItem := range menuItems {
		DrawCentered(screen, contentOffsetY+2+i, opts.Style, menuItem)
	}
}

func drawCustomInput(
	screen tcell.Screen, titleItems []string,
	selected int,
	cfg DifficultyConfig,
	buf string, errMsg string,
	opts *GameOptions,
) {
	_, h := screen.Size()

	titleHeight := len(titleItems)

	menuItems := []string{
		fmt.Sprintf("Rows: <%d>", cfg.Rows),
		fmt.Sprintf("Cols: <%d>", cfg.Cols),
		fmt.Sprintf("BombCount: %d", cfg.BombCount),
		"Start",
		"Back",
	}
	menuHeight := len(menuItems)*2 - 1

	contentHeight := titleHeight + 4 + menuHeight
	titleOffsetY := (h-contentHeight)/2 - contentHeight%2

	drawTitleItems(screen, titleItems, titleOffsetY, opts)
	drawMenuItems(screen, selected, "⚑⚑⚑ Custom Difficulty  ⚑⚑⚑", menuItems, titleOffsetY+titleHeight+2, opts)

	if selected == 2 {
		DrawCentered(screen, titleOffsetY+titleHeight+2+(len(menuItems)+1)*2, opts.Style, fmt.Sprintf("Typing: %s", buf))
	}
	if errMsg != "" {
		DrawCentered(screen, titleOffsetY+titleHeight+2+(len(menuItems)+2)*2, opts.Style, "Error: "+errMsg)
	}
}

func drawHelpHint(screen tcell.Screen, opts *GameOptions) {
	w, h := screen.Size()
	message := "W/S = up/down, A/D = change, Enter = select, Esc/Backspace = back"
	DrawString(screen, w-len(message)-1, h-1, opts.Style, message)
}

func RunMenu(screen tcell.Screen, opts *GameOptions) (GameState, *GameOptions, DifficultyConfig, bool) {
	page := PageMain
	titleItems := assets.RandomTitle()
	bgs := append([]string{"none"}, assets.ListBackgrounds()...)
	bgIndex := slices.Index(bgs, opts.Background)
	selected := 0
	difficulties := []string{"beginner", "intermediate", "advanced", "expert", "insane", "custom"}
	difficultiesNG := []string{"beginner", "intermediate", "advanced", "expert", "insane"}
	diffIndex := 0
	diffNGIndex := 0
	playingNG := false
	customCfg := DifficultyConfig{Rows: 9, Cols: 9, BombCount: 10}
	rowsOptions := make([]int, MAX_ROWS)
	for i := range MAX_ROWS {
		rowsOptions[i] = i + 1
	}
	colsOptions := make([]int, MAX_COLS)
	for i := range MAX_COLS {
		colsOptions[i] = i + 1
	}
	rowsIndex := 8
	colsIndex := 8
	inputBuffer := ""
	errorMsg := ""

	var menuCount int

	for {
		screen.Clear()
		DrawBackground(screen, bgs[bgIndex], false)
		switch page {
		case PageMain:
			drawMainMenu(screen, titleItems, selected, difficulties[diffIndex], difficultiesNG[diffNGIndex], opts)
			menuCount = 5
		case PageOptions:
			drawOptionsMenu(screen, titleItems, selected, opts)
			menuCount = 4
		case PageCredits:
			drawCredits(screen, titleItems, opts)
		case PageQuitConfirm:
			drawQuitConfirm(screen, opts)
		case PageCustomInput:
			drawCustomInput(screen, titleItems, selected, customCfg, inputBuffer, errorMsg, opts)
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
				case PageMain, PageOptions, PageCustomInput:
					selected = (selected - 1 + menuCount) % menuCount
				}
			case tcell.KeyDown:
				switch page {
				case PageMain, PageOptions, PageCustomInput:
					selected = (selected + 1) % menuCount
				}
			case tcell.KeyLeft:
				switch page {
				case PageMain:
					switch selected {
					case 0:
						diffIndex = (diffIndex - 1 + len(difficulties)) % len(difficulties)
					case 1:
						diffNGIndex = (diffNGIndex - 1 + len(difficultiesNG)) % len(difficultiesNG)
					}
				case PageOptions:
					switch selected {
					// Show inner borders
					case 0:
						opts.ShowInnerBorders = !opts.ShowInnerBorders
					// Border style
					case 1:
						opts.BorderStyle = (opts.BorderStyle - 1 + 2) % 2
					// Background
					case 2:
						bgIndex = (bgIndex - 1 + len(bgs)) % len(bgs)
						opts.Background = bgs[bgIndex]
					}
				case PageCustomInput:
					switch selected {
					case 0:
						rowsIndex = (rowsIndex - 1 + len(rowsOptions)) % len(rowsOptions)
						customCfg.Rows = rowsOptions[rowsIndex]
					case 1:
						colsIndex = (colsIndex - 1 + len(colsOptions)) % len(colsOptions)
						customCfg.Cols = colsOptions[colsIndex]
					}
				}
			case tcell.KeyRight:
				switch page {
				case PageMain:
					switch selected {
					case 0:
						diffIndex = (diffIndex + 1) % len(difficulties)
					case 1:
						diffNGIndex = (diffNGIndex + 1) % len(difficultiesNG)
					}
				case PageOptions:
					switch selected {
					// Show inner borders
					case 0:
						opts.ShowInnerBorders = !opts.ShowInnerBorders
					// Border style
					case 1:
						opts.BorderStyle = (opts.BorderStyle + 1) % 2
					// Background
					case 2:
						bgIndex = (bgIndex + 1) % len(bgs)
						opts.Background = bgs[bgIndex]
					}
				case PageCustomInput:
					switch selected {
					case 0:
						rowsIndex = (rowsIndex + 1) % len(rowsOptions)
						customCfg.Rows = rowsOptions[rowsIndex]
					case 1:
						colsIndex = (colsIndex + 1) % len(colsOptions)
						customCfg.Cols = colsOptions[colsIndex]
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
						if difficultiesNG[diffNGIndex] == "custom" {
							page = PageCustomInput
						} else {
							opts.Difficulty = DifficultyMap[difficultiesNG[diffNGIndex]]
							return StatePlaying, opts, DifficultyMap[difficultiesNG[diffNGIndex]], playingNG
						}
					// Options
					case 2:
						page = PageOptions
						selected = 0
					// Credits
					case 3:
						page = PageCredits
					// Quit
					case menuCount - 1:
						page = PageQuitConfirm
					}
				case PageOptions:
					if selected == menuCount-1 {
						page = PageMain
						selected = 0
					}
				case PageCustomInput:
					switch selected {
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
					case menuCount - 1:
						page = PageMain
						selected = 0
						inputBuffer = ""
						errorMsg = ""
					default:
						if inputBuffer != "" {
							val, err := strconv.Atoi(inputBuffer)
							if err == nil {
								if selected == 2 {
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
				if page != PageMain && page != PageCustomInput {
					page = PageMain
					selected = 0
					inputBuffer = ""
					errorMsg = ""
				}
				if len(inputBuffer) > 0 {
					inputBuffer = inputBuffer[:len(inputBuffer)-1]
				}
			case tcell.KeyRune:
				r := ev.Rune()
				if page == PageCustomInput && selected == 2 && unicode.IsDigit(r) {
					inputBuffer += string(r)
				} else {
					switch r {
					case 'w':
						if page == PageMain || page == PageOptions || page == PageCustomInput {
							selected = (selected - 1 + menuCount) % menuCount
						}
					case 's':
						if page == PageMain || page == PageOptions || page == PageCustomInput {
							selected = (selected + 1) % menuCount
						}
					case 'a':
						switch page {
						case PageMain:
							switch selected {
							case 0:
								diffIndex = (diffIndex - 1 + len(difficulties)) % len(difficulties)
							case 1:
								diffNGIndex = (diffNGIndex - 1 + len(difficultiesNG)) % len(difficultiesNG)
							}
						case PageOptions:
							switch selected {
							// Show inner borders
							case 0:
								opts.ShowInnerBorders = !opts.ShowInnerBorders
							// Border style
							case 1:
								opts.BorderStyle = (opts.BorderStyle - 1 + 2) % 2
							// Background
							case 2:
								bgIndex = (bgIndex - 1 + len(bgs)) % len(bgs)
								opts.Background = bgs[bgIndex]
							}
						case PageCustomInput:
							switch selected {
							case 0:
								rowsIndex = (rowsIndex - 1 + len(rowsOptions)) % len(rowsOptions)
								customCfg.Rows = rowsOptions[rowsIndex]
							case 1:
								colsIndex = (colsIndex - 1 + len(colsOptions)) % len(colsOptions)
								customCfg.Cols = colsOptions[colsIndex]
							}
						}
					case 'd':
						switch page {
						case PageMain:
							switch selected {
							case 0:
								diffIndex = (diffIndex + 1) % len(difficulties)
							case 1:
								diffNGIndex = (diffNGIndex + 1) % len(difficultiesNG)
							}
						case PageOptions:
							switch selected {
							// Show inner borders
							case 0:
								opts.ShowInnerBorders = !opts.ShowInnerBorders
							// Border style
							case 1:
								opts.BorderStyle = (opts.BorderStyle + 1) % 2
							case 2:
								bgIndex = (bgIndex + 1) % len(bgs)
								opts.Background = bgs[bgIndex]
							}
						case PageCustomInput:
							switch selected {
							case 0:
								rowsIndex = (rowsIndex + 1) % len(rowsOptions)
								customCfg.Rows = rowsOptions[rowsIndex]
							case 1:
								colsIndex = (colsIndex + 1) % len(colsOptions)
								customCfg.Cols = colsOptions[colsIndex]
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
