package main

import (
	"errors"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

type Cell struct {
	Value    int
	Revealed bool
	Flagged  bool
}

type Minesweeper struct {
	Rows          int
	Cols          int
	BombCount     int
	Grid          [][]Cell
	BombPositions [][2]int
	IsGameOver    bool
	IsWon         bool
	RevealedCount int
}

type DifficultyConfig struct {
	Rows      int
	Cols      int
	BombCount int
}

const CLEAR int = 0
const BOMB int = -1

var intToRune = map[int]rune{
	CLEAR: ' ',
	BOMB:  '¤',
	1:     '1',
	2:     '2',
	3:     '3',
	4:     '4',
	5:     '5',
	6:     '6',
	7:     '7',
	8:     '8',
	9:     '9',
}

var borderSets = map[BorderStyle]map[string]rune{
	BorderThin: {
		"topLeft":     '┌',
		"topRight":    '┐',
		"bottomLeft":  '└',
		"bottomRight": '┘',
		"horizontal":  '─',
		"vertical":    '│',
		"tUp":         '┴',
		"tDown":       '┬',
		"tRight":      '├',
		"tLeft":       '┤',
		"cross":       '┼',
	},
	BorderThick: {
		"topLeft":     '╔',
		"topRight":    '╗',
		"bottomLeft":  '╚',
		"bottomRight": '╝',
		"horizontal":  '═',
		"vertical":    '║',
		"tUp":         '╩',
		"tDown":       '╦',
		"tRight":      '╠',
		"tLeft":       '╣',
		"cross":       '╬',
	},
}

var DifficultyMap = map[string]DifficultyConfig{
	"beginner": {
		Rows:      9,
		Cols:      9,
		BombCount: 10,
	},
	"intermediate": {
		Rows:      16,
		Cols:      16,
		BombCount: 40,
	},
	"expert": {
		Rows:      16,
		Cols:      30,
		BombCount: 99,
	},
}

var directions = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func (m *Minesweeper) Reveal(row, col int, userClick bool) {
	if m.IsGameOver {
		return
	}
	cell := &m.Grid[row][col]
	if userClick {
		if cell.Flagged {
			return
		}
		if cell.Revealed {
			m.Chord(row, col)
			return
		}
		if cell.Value == BOMB {
			cell.Revealed = true
			m.IsGameOver = true
			return
		}
	} else {
		if cell.Revealed {
			return
		}
		if cell.Value == BOMB {
			cell.Revealed = true
			m.IsGameOver = true
			return
		}
	}

	cell.Revealed = true
	m.RevealedCount++
	if m.RevealedCount == m.Rows*m.Cols-m.BombCount {
		m.IsGameOver = true
		m.IsWon = true
	}

	if cell.Value == CLEAR {
		for _, direction := range directions {
			newRow := row + direction[0]
			newCol := col + direction[1]
			if newRow < 0 || newRow >= m.Rows || newCol < 0 || newCol >= m.Cols {
				continue
			}
			m.Reveal(newRow, newCol, false)
		}
	}
}

func (m *Minesweeper) Chord(row, col int) {
	cell := &m.Grid[row][col]
	if cell.Revealed && cell.Value > 0 {
		unflaggedCells := make([][2]int, 0, 8)
		flaggedCells := make([][2]int, 0, 8)
		for _, direction := range directions {
			newRow := row + direction[0]
			newCol := col + direction[1]
			if newRow < 0 || newRow >= m.Rows || newCol < 0 || newCol >= m.Cols {
				continue
			}
			if !m.Grid[newRow][newCol].Revealed {
				if m.Grid[newRow][newCol].Flagged {
					flaggedCells = append(flaggedCells, [2]int{newRow, newCol})
				} else {
					unflaggedCells = append(unflaggedCells, [2]int{newRow, newCol})
				}
			}
		}

		if len(flaggedCells) == cell.Value {
			for _, pos := range unflaggedCells {
				r, c := pos[0], pos[1]
				m.Reveal(r, c, false)
			}
		}
	}
}

func (m *Minesweeper) Flag(row, col int) {
	if m.IsGameOver {
		return
	}
	cell := &m.Grid[row][col]
	if !cell.Revealed {
		cell.Flagged = !cell.Flagged
	}
}

func addBomb(grid [][]Cell, rows, cols, row, col int) {
	grid[row][col].Value = BOMB
	for _, direction := range directions {
		newRow := row + direction[0]
		newCol := col + direction[1]
		if newRow < 0 || newRow >= rows || newCol < 0 || newCol >= cols {
			continue
		}
		if grid[newRow][newCol].Value != BOMB {
			grid[newRow][newCol].Value++
		}
	}
}

func GenerateBoard(cfg DifficultyConfig) (*Minesweeper, error) {
	if cfg.Rows <= 0 || cfg.Cols <= 0 || cfg.BombCount <= 0 {
		return nil, errors.New("rows, cols, and bombCount must be non-negative integer")
	}

	grid := make([][]Cell, cfg.Rows)
	for r := range grid {
		grid[r] = make([]Cell, cfg.Cols)
	}

	bombPositions := make([][2]int, 0)
	for len(bombPositions) < cfg.BombCount {
		pos := rand.Intn(cfg.Rows * cfg.Cols)
		r, c := pos/cfg.Cols, pos%cfg.Cols
		if grid[r][c].Value != BOMB {
			bombPositions = append(bombPositions, [2]int{r, c})
			addBomb(grid, cfg.Rows, cfg.Cols, r, c)
		}
	}

	m := &Minesweeper{
		Rows:          cfg.Rows,
		Cols:          cfg.Cols,
		BombCount:     cfg.BombCount,
		Grid:          grid,
		BombPositions: bombPositions,
		IsGameOver:    false,
		IsWon:         false,
		RevealedCount: 0,
	}

	return m, nil
}

func (m *Minesweeper) drawBombs(
	screen tcell.Screen,
	showInnerBorders bool,
	screenX, screenY int,
) {
	cellWidth, cellHeight := 1, 1
	if showInnerBorders {
		cellWidth, cellHeight = 2, 2
	}

	for _, pos := range m.BombPositions {
		var (
			char  rune
			style tcell.Style
		)
		r, c := pos[0], pos[1]
		cell := m.Grid[r][c]
		if cell.Flagged {
			continue
		}
		if m.IsWon {
			char = '⚑'
			style = tcell.StyleDefault.Background(tcell.ColorOrange).Foreground(tcell.ColorDarkRed)
		} else {
			char = intToRune[BOMB]
			style = ValueToCellStyle[cell.Value]
		}
		NewSprite(
			char,
			screenX+(cellWidth*c+1),
			screenY+(cellHeight*r+1),
		).Draw(screen, style)
	}
}

func (m *Minesweeper) Draw(
	screen tcell.Screen,
	border BorderStyle,
	showInnerBorders bool,
	screenX, screenY int,
) {
	runes := borderSets[border]
	cellWidth, cellHeight := 1, 1
	if showInnerBorders {
		cellWidth, cellHeight = 2, 2
	}

	NewSprite(
		runes["topLeft"],
		screenX, screenY,
	).Draw(screen, DefaultBorderStyle)
	for j := 0; j < m.Cols; j++ {
		NewSprite(
			runes["horizontal"],
			screenX+(cellWidth*j+1), screenY,
		).Draw(screen, DefaultBorderStyle)
		if j < m.Cols-1 && showInnerBorders {
			NewSprite(
				runes["tDown"],
				screenX+(cellWidth*j+2), screenY,
			).Draw(screen, DefaultBorderStyle)
		} else {
			NewSprite(
				runes["topRight"],
				screenX+(cellWidth*j+2), screenY,
			).Draw(screen, DefaultBorderStyle)
		}
	}

	for i := 0; i < m.Rows; i++ {
		NewSprite(
			runes["vertical"],
			screenX, screenY+(cellHeight*i+1),
		).Draw(screen, DefaultBorderStyle)
		for j := 0; j < m.Cols; j++ {
			var (
				char  rune
				style tcell.Style
			)
			cell := m.Grid[i][j]
			if !cell.Revealed {
				if cell.Flagged {
					if m.IsGameOver && cell.Value != BOMB {
						char = '×'
					} else {
						char = '⚑'
					}
					style = tcell.StyleDefault.Background(tcell.ColorOrange).Foreground(tcell.ColorDarkRed)
				} else {
					char = ' '
					style = DefaultBorderStyle
				}
			} else {
				char = intToRune[cell.Value]
				style = ValueToCellStyle[cell.Value]
			}
			NewSprite(
				char,
				screenX+(cellWidth*j+1), screenY+(cellHeight*i+1),
			).Draw(screen, style)
			NewSprite(
				runes["vertical"],
				screenX+(cellWidth*j+2), screenY+(cellHeight*i+1),
			).Draw(screen, DefaultBorderStyle)
		}

		if i < m.Rows-1 && showInnerBorders {
			NewSprite(
				runes["tRight"],
				screenX, screenY+(cellHeight*i+2),
			).Draw(screen, DefaultBorderStyle)
			for j := 0; j < m.Cols; j++ {
				NewSprite(
					runes["horizontal"],
					screenX+(cellWidth*j+1), screenY+(cellHeight*i+2),
				).Draw(screen, DefaultBorderStyle)
				if j < m.Cols-1 {
					NewSprite(
						runes["cross"],
						screenX+(cellWidth*j+2), screenY+(cellHeight*i+2),
					).Draw(screen, DefaultBorderStyle)
				} else {
					NewSprite(
						runes["tLeft"],
						screenX+(cellWidth*j+2), screenY+(cellHeight*i+2),
					).Draw(screen, DefaultBorderStyle)
				}
			}
		} else {
			NewSprite(
				runes["bottomLeft"],
				screenX, screenY+(cellHeight*i+2),
			).Draw(screen, DefaultBorderStyle)
			for j := 0; j < m.Cols; j++ {
				NewSprite(
					runes["horizontal"],
					screenX+(cellWidth*j+1), screenY+(cellHeight*i+2),
				).Draw(screen, DefaultBorderStyle)
				if j < m.Cols-1 {
					NewSprite(
						runes["tUp"],
						screenX+(cellWidth*j+2), screenY+(cellHeight*i+2),
					).Draw(screen, DefaultBorderStyle)
				} else {
					NewSprite(
						runes["bottomRight"],
						screenX+(cellWidth*j+2), screenY+(cellHeight*i+2),
					).Draw(screen, DefaultBorderStyle)
				}
			}
		}
	}

	if m.IsGameOver {
		m.drawBombs(screen, showInnerBorders, screenX, screenY)
	}
}

func (m *Minesweeper) ScreenToGrid(
	screenX, screenY,
	offsetX, offsetY int,
	showInnerBorders bool,
) (row, col int, ok bool) {
	relX := screenX - offsetX
	relY := screenY - offsetY
	cellWidth, cellHeight := 1, 1
	if showInnerBorders {
		cellWidth, cellHeight = 2, 2
	}

	if relX <= 0 || relY <= 0 {
		return -1, -1, false
	}

	if showInnerBorders && ((relX-1)%cellWidth == 0 || (relY-1)%cellHeight == 0) {
		return -1, -1, false
	}

	row = (relY - 1) / cellHeight
	col = (relX - 1) / cellWidth
	if row < 0 || row >= m.Rows || col < 0 || col >= m.Cols {
		return -1, -1, false
	}

	return row, col, true
}
