package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

type Cell struct {
	Value    int
	Revealed bool
	Flagged  bool
}

type Minesweeper struct {
	Rows              int
	Cols              int
	BombCount         int
	Grid              [][]Cell
	BombPositions     [][2]int
	PositionToValue   map[[2]int]int
	IsGameOver        bool
	IsWon             bool
	RevealedCount     int
	StartCell         *Cell
	StartCellPosition [2]int
}

type DifficultyConfig struct {
	Rows      int
	Cols      int
	BombCount int
}

const (
	CLEAR int = 0
	BOMB  int = -1

	MAX_ROWS int = 36
	MAX_COLS int = 160
)

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
	"advanced": {
		Rows:      16,
		Cols:      30,
		BombCount: 99,
	},
	"expert": {
		Rows:      24,
		Cols:      30,
		BombCount: 150,
	},
	"insane": {
		Rows:      30,
		Cols:      30,
		BombCount: 199,
	},
}

var directions = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func (m *Minesweeper) isOutOfBounds(row, col int) bool {
	return row < 0 || row >= m.Rows || col < 0 || col >= m.Cols
}

func (m *Minesweeper) getNeighborsOf(row, col int) [][2]int {
	neighbors := make([][2]int, 0)
	for _, direction := range directions {
		newRow := row + direction[0]
		newCol := col + direction[1]
		if !m.isOutOfBounds(newRow, newCol) {
			neighbors = append(neighbors, [2]int{newRow, newCol})
		}
	}

	return neighbors
}

func (m *Minesweeper) addBomb(row, col int) {
	m.Grid[row][col].Value = BOMB

	neighbors := m.getNeighborsOf(row, col)
	for _, neighbor := range neighbors {
		if m.Grid[neighbor[0]][neighbor[1]].Value != BOMB {
			m.Grid[neighbor[0]][neighbor[1]].Value++
			if _, ok := m.PositionToValue[neighbor]; !ok {
				m.PositionToValue[neighbor] = 0
			}
			m.PositionToValue[neighbor]++
		}
	}
}

func isAdjacent(row1, col1, row2, col2 int) bool {
	return (-1 <= row1-row2 && row1-row2 <= 1) && (-1 <= col1-col2 && col1-col2 <= 1)
}

func validateConfig(cfg DifficultyConfig) error {
	if cfg.Rows <= 0 || cfg.Cols <= 0 || cfg.BombCount <= 0 {
		return errors.New("rows, cols, and bombCount must be non-negative integer")
	} else if cfg.Rows > MAX_ROWS {
		return fmt.Errorf("maximum number of rows is capped at %d", MAX_ROWS)
	} else if cfg.Cols > MAX_COLS {
		return fmt.Errorf("maximum number of cols is capped at %d", MAX_COLS)
	}

	if cfg.BombCount >= cfg.Rows*cfg.Cols {
		return fmt.Errorf("too many bombCount, must be in the range of [1, %d]", cfg.Rows*cfg.Cols-1)
	}

	return nil
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
			style = FlagStyle
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

func (m *Minesweeper) Reveal(row, col int, userClick bool) bool {
	if m.IsGameOver {
		return false
	}
	cell := &m.Grid[row][col]

	// Don't reveal flagged cell when user clicked
	if userClick && cell.Flagged {
		return false
	}

	// Cell with bomb is clicked/revealed
	if cell.Value == BOMB {
		cell.Revealed = true
		m.IsGameOver = true
		return true
	}

	// Cell is already revealed
	if cell.Revealed {
		if userClick {
			return m.Chord(row, col)
		}
		return false
	}

	// Normal cell reveal
	cell.Revealed = true
	m.RevealedCount++
	if m.RevealedCount == m.Rows*m.Cols-m.BombCount {
		m.IsGameOver = true
		m.IsWon = true
		return true
	}

	// Flood-fill expansion
	if cell.Value == CLEAR {
		neighbors := m.getNeighborsOf(row, col)
		for _, neighbor := range neighbors {
			m.Reveal(neighbor[0], neighbor[1], false)
		}
	}

	return true
}

func (m *Minesweeper) Chord(row, col int) bool {
	cell := &m.Grid[row][col]
	ok := false

	if cell.Revealed && cell.Value > 0 {
		unflaggedCells := make([][2]int, 0, 8)
		flaggedCells := make([][2]int, 0, 8)

		neighbors := m.getNeighborsOf(row, col)
		for _, neighbor := range neighbors {
			if !m.Grid[neighbor[0]][neighbor[1]].Revealed {
				if m.Grid[neighbor[0]][neighbor[1]].Flagged {
					flaggedCells = append(flaggedCells, neighbor)
				} else {
					unflaggedCells = append(unflaggedCells, neighbor)
				}
			}
		}

		if len(flaggedCells) == cell.Value {
			for _, pos := range unflaggedCells {
				r, c := pos[0], pos[1]
				if m.Reveal(r, c, false) {
					ok = true
				}
			}
		}
	}

	return ok
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

func GenerateBoard(cfg DifficultyConfig) (*Minesweeper, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	m := &Minesweeper{
		Rows:              cfg.Rows,
		Cols:              cfg.Cols,
		BombCount:         cfg.BombCount,
		IsGameOver:        false,
		IsWon:             false,
		RevealedCount:     0,
		StartCell:         nil,
		StartCellPosition: [2]int{-1, -1},
	}

	m.Grid = make([][]Cell, m.Rows)
	for r := range m.Grid {
		m.Grid[r] = make([]Cell, m.Cols)
	}

	m.BombPositions = make([][2]int, 0)
	m.PositionToValue = make(map[[2]int]int)
	for len(m.BombPositions) < m.BombCount {
		pos := rand.Intn(m.Rows * m.Cols)
		r, c := pos/m.Cols, pos%m.Cols
		if m.Grid[r][c].Value != BOMB {
			m.BombPositions = append(m.BombPositions, [2]int{r, c})
			m.addBomb(r, c)
		}
	}

	return m, nil
}

func GenerateBoardWithStartCell(cfg DifficultyConfig) (*Minesweeper, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	m := &Minesweeper{
		Rows:          cfg.Rows,
		Cols:          cfg.Cols,
		BombCount:     cfg.BombCount,
		IsGameOver:    false,
		IsWon:         false,
		RevealedCount: 0,
	}

	m.Grid = make([][]Cell, m.Rows)
	for r := range m.Grid {
		m.Grid[r] = make([]Cell, m.Cols)
	}

	startCellPos := rand.Intn(m.Rows * m.Cols)
	startCellRow, startCellCol := startCellPos/m.Cols, startCellPos%m.Cols
	m.StartCellPosition = [2]int{startCellRow, startCellCol}
	m.StartCell = &m.Grid[startCellRow][startCellCol]

	m.BombPositions = make([][2]int, 0)
	m.PositionToValue = make(map[[2]int]int)
	for len(m.BombPositions) < m.BombCount {
		pos := rand.Intn(m.Rows * m.Cols)
		r, c := pos/m.Cols, pos%m.Cols
		if !isAdjacent(startCellRow, startCellCol, r, c) && m.Grid[r][c].Value != BOMB {
			m.BombPositions = append(m.BombPositions, [2]int{r, c})
			m.addBomb(r, c)
		}
	}

	return m, nil
}

func GenerateNGBoardWithStartCell(cfg DifficultyConfig, tries, maxComponentSize int) (*Minesweeper, error) {
	for attempt := 1; attempt <= tries; attempt++ {
		m, err := GenerateBoardWithStartCell(cfg)
		if err != nil {
			return nil, err
		}

		solvable, _, _ := m.DeterministicSolve(maxComponentSize)
		if solvable {
			return m, nil
		}
	}

	return nil, errors.New("no NG board satisfies the difficulty")
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
			cell := &m.Grid[i][j]
			if !cell.Revealed {
				if cell.Flagged {
					if m.IsGameOver && cell.Value != BOMB {
						char = '×'
					} else {
						char = '⚑'
					}
					style = FlagStyle
				} else {
					if m.StartCell == cell {
						char = '✓'
						style = StartCellStyle
					} else {
						char = ' '
						style = DefaultBorderStyle
					}
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

func (m *Minesweeper) DrawSmiley(
	screen tcell.Screen,
	screenY int, style tcell.Style,
	lastMouseButtons tcell.ButtonMask,
) {
	if m.IsGameOver {
		var message string
		if m.IsWon {
			message = "You win!"
			DrawCentered(screen, screenY-3, style, "😎")
		} else {
			message = "You lose!"
			DrawCentered(screen, screenY-3, style, "😭")
		}
		DrawCentered(screen, screenY-2, style, message)
		DrawCentered(screen, screenY-1, style, "Press 'r' to create a new board, 'q' to quit to main menu.")
	} else if lastMouseButtons == tcell.Button1 {
		DrawCentered(screen, screenY-3, style, "😮")
	} else {
		DrawCentered(screen, screenY-3, style, "🙂")
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
	if m.isOutOfBounds(row, col) {
		return -1, -1, false
	}

	return row, col, true
}
