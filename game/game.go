package game

import (
	"fmt"
	"math/rand"
)

type Cell struct {
	Content CellValue
	State   CellState
}

type CellValue int

const ZeroCellValue CellValue = 0
const OneCellValue CellValue = 1
const TwoCellValue CellValue = 2
const ThreeCellValue CellValue = 3
const FourCellValue CellValue = 4
const FiveCellValue CellValue = 5
const SixCellValue CellValue = 6
const SevenCellValue CellValue = 7
const EightCellValue CellValue = 8
const BlackHoleCellValue CellValue = -1

type CellState int

const (
	HiddenState CellState = iota
	VisibleState
)

type cellAddress struct {
	row    int
	column int
}

func NewGame(boardSize int, blackHolesNumber int) *Game {
	board := createGameState(boardSize, boardSize, Cell{Content: ZeroCellValue, State: HiddenState})
	blackHoleAddresses := generateBlackHoleAddresses(boardSize, boardSize, blackHolesNumber)
	replaceCells(board, blackHoleAddresses, Cell{Content: BlackHoleCellValue, State: HiddenState})
	updateNaboringBlackHolesCellValues(board)

	return &Game{failAt: nil, board: board}
}

// Game is a contrainer for a game state and implements methods to update state
// according game rules.
type Game struct {
	failAt *cellAddress
	board  [][]Cell
}

// GetState clones the current Game state
func (g *Game) GetState() [][]Cell {
	rows := make([][]Cell, len(g.board))
	for i, boardRow := range g.board {
		row := make([]Cell, len(boardRow))
		copy(row, boardRow)
		rows[i] = row
	}

	return rows
}

// Lost returns true if the game is lost. Otherwise false is returned.
func (g *Game) Lost() bool {
	return g.failAt != nil
}

// Won returns true if the game is won. Otherwise false is returned.
func (g *Game) Won() bool {
	if g.Lost() {
		return false
	}

	return g.won()
}

// won return true if all cells are revealed except of cells with
// BlackHoleCellValue in State field.
func (g *Game) won() bool {
	totalCells := 0
	visibleCells := 0
	blackHoleCells := 0

	for _, row := range g.board {
		for _, cell := range row {
			totalCells++
			if cell.State == VisibleState {
				visibleCells++
			}
			if cell.Content == BlackHoleCellValue {
				blackHoleCells++
			}
		}
	}

	return (totalCells - visibleCells) == blackHoleCells
}

// Completed returns true if the game is over.
func (g *Game) Completed() bool {
	return g.Lost() || g.Won()
}

// RevealCell update the cell State field to be `VisibleState`.
// If the game is completed, then RevealCell will panic.
// If supplied i row and j column can not address a cell in the game, then
// RevealCall panics.
func (g *Game) RevealCell(i, j int) {
	if g.Completed() {
		panic("game over")
	}

	address := cellAddress{
		row:    i,
		column: j,
	}

	cell := getCell(g.board, address)
	if cell == nil {
		panic("non-existing cell addressed")
	}

	if cell.Content == BlackHoleCellValue {
		g.failAt = &address

		return
	}

	if cell.State == VisibleState {
		panic("cell already visible")
	}

	cell.State = VisibleState

	if cell.Content == ZeroCellValue {
		g.revealSurrounding(address)
	}
}

func getCell(board [][]Cell, a cellAddress) *Cell {
	if a.row < 0 || a.row >= len(board) {
		return nil
	}
	if a.column < 0 || a.column >= len(board[a.row]) {
		return nil
	}

	return &board[a.row][a.column]
}

func (g *Game) revealSurrounding(ca cellAddress) {
	for _, currentAddress := range surroundingAddresses(ca) {
		cell := getCell(g.board, currentAddress)

		if cell == nil {
			continue
		}

		if cell.State == VisibleState {
			continue
		}

		if cell.Content == BlackHoleCellValue {
			panic("unexpected condition: revealing naboring black hole")
		}

		cell.State = VisibleState

		if cell.Content == ZeroCellValue {
			g.revealSurrounding(currentAddress)
		}
	}
}

func surroundingAddresses(a cellAddress) []cellAddress {
	if a.row < 0 {
		panic("row must not be negative")
	}
	if a.column < 0 {
		panic("column must not be negative")
	}

	addresses := make([]cellAddress, 0, 8)

	rowMin := a.row - 1
	if rowMin < 0 {
		rowMin = 0
	}
	rowMax := a.row + 1

	columnMin := a.column - 1
	if columnMin < 0 {
		columnMin = 0
	}
	columnMax := a.column + 1

	for i := rowMin; i <= rowMax; i++ {
		for j := columnMin; j <= columnMax; j++ {
			if i == a.row && j == a.column {
				continue
			}

			address := cellAddress{
				row:    i,
				column: j,
			}

			addresses = append(addresses, address)
		}
	}

	return addresses
}

func createGameState(rows int, columns int, defaultCell Cell) [][]Cell {
	stateRows := make([][]Cell, rows)
	for rowNumber := range stateRows {
		row := make([]Cell, columns)
		for i := range row {
			row[i] = defaultCell
		}
		stateRows[rowNumber] = row
	}

	return stateRows
}

func replaceCells(board [][]Cell, addresses []cellAddress, defaultCell Cell) {
	for _, a := range addresses {
		board[a.row][a.column] = defaultCell
	}
}

func generateBlackHoleAddresses(rows int, columns int, amount int) []cellAddress {
	if rows < 0 {
		panic("rows must not be negative")
	}
	if columns < 0 {
		panic("columns must not be negative")
	}
	if amount < 0 {
		panic("amount must not be negative")
	}
	if amount > rows*columns {
		panic("number of black holes does not fit in game board")
	}

	addresses := make([]cellAddress, 0, amount)

	for len(addresses) < amount {
		address := cellAddress{
			row:    rand.Intn(rows),
			column: rand.Intn(columns),
		}
		if addressInList(address, addresses) {
			continue
		}
		addresses = append(addresses, address)
	}

	return addresses
}

func addressInList(address cellAddress, list []cellAddress) bool {
	for _, listed := range list {
		if listed == address {
			return true
		}
	}

	return false
}

func updateNaboringBlackHolesCellValues(board [][]Cell) {
	for i, rows := range board {
		for j := range rows {
			cell := &board[i][j]

			if cell.Content == BlackHoleCellValue {
				continue
			}

			cell.Content = calculateNaboringBlackHoles(board, cellAddress{row: i, column: j})
		}
	}
}

func calculateNaboringBlackHoles(board [][]Cell, a cellAddress) CellValue {
	count := 0
	for _, currentAddress := range surroundingAddresses(a) {
		cell := getCell(board, currentAddress)

		if cell == nil {
			continue
		}

		if cell.Content == BlackHoleCellValue {
			count++
		}
	}

	switch count {
	case 0:
		return ZeroCellValue
	case 1:
		return OneCellValue
	case 2:
		return TwoCellValue
	case 3:
		return ThreeCellValue
	case 4:
		return FourCellValue
	case 5:
		return FiveCellValue
	case 6:
		return SixCellValue
	case 7:
		return SevenCellValue
	case 8:
		return EightCellValue
	default:
		panic("unexpected counted surrounding black holes: " + fmt.Sprint(count))
	}
}
