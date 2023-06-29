package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/kalynv/proxx/game"
)

func main() {
	boardSize := 3
	blackHoles := 2
	theGame := game.NewGame(boardSize, blackHoles)
	adapter := newGameAdapter(theGame, os.Stdin, os.Stdout)

	adapter.Play()
}

func newGameAdapter(g *game.Game, in io.Reader, out io.Writer) *gameAdapter {
	return &gameAdapter{
		game: g,
		in:   in,
		out:  out,
	}
}

type gameAdapter struct {
	game *game.Game
	in   io.Reader
	out  io.Writer
}

func (ga *gameAdapter) Play() {
	fmt.Fprintln(ga.out, "Game started!")

	for !ga.game.Completed() {
		ga.displayBoard(presentCellAtGameTime)
		row, column := ga.readRevealCell()
		ga.game.RevealCell(row, column)
	}

	ga.displayBoard(presentCellRevealed)

	if ga.game.Won() {
		fmt.Fprintln(ga.out, "You won!")
	}
	if ga.game.Lost() {
		fmt.Fprintln(ga.out, "You lost.")
	}
	fmt.Fprintln(ga.out, "Game over")
}

func (ga *gameAdapter) readRevealCell() (row, column int) {
	fmt.Fprint(ga.out, "\nRevealing cell by row and column\n")
	for {
		var err error
		fmt.Fprint(ga.out, "Enter row: ")
		row, err = readInt(ga.in)
		if err != nil {
			fmt.Fprintf(ga.out, "%s\n", err.Error())

			continue
		}

		break
	}

	for {
		var err error
		fmt.Fprint(ga.out, "Enter column: ")
		column, err = readInt(ga.in)
		if err != nil {
			fmt.Fprintf(ga.out, "%s\n", err.Error())

			continue
		}

		break
	}

	return row, column
}

func readInt(r io.Reader) (int, error) {
	line, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(strings.Trim(line, " \n"))
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (ga *gameAdapter) displayBoard(presentCell func(c game.Cell) rune) {
	board := ga.game.GetState()

	presentedBoard := strings.Builder{}
	presentedBoard.WriteString("\nBoard:\n")

	for _, row := range board {
		presentedRow := strings.Builder{}
		for _, cell := range row {
			presentedRow.WriteRune(' ')
			presentedRow.WriteRune(presentCell(cell))
		}
		presentedBoard.WriteString(presentedRow.String())
		presentedBoard.WriteString("\n")
	}

	presentedBoard.WriteString("\n")

	fmt.Fprint(ga.out, presentedBoard.String())
}

func presentCellRevealed(c game.Cell) rune {
	return convertCellValue(c.Content)
}

func presentCellAtGameTime(c game.Cell) rune {
	switch c.State {
	case game.HiddenState:
		return 'H'
	case game.VisibleState:
		return convertCellValue(c.Content)
	default:
		return '_'
	}
}

func convertCellValue(v game.CellValue) rune {
	switch v {
	case game.BlackHoleCellValue:
		return '*'
	case game.ZeroCellValue:
		return '0'
	case game.OneCellValue:
		return '1'
	case game.TwoCellValue:
		return '2'
	case game.ThreeCellValue:
		return '3'
	case game.FourCellValue:
		return '4'
	case game.FiveCellValue:
		return '5'
	case game.SixCellValue:
		return '6'
	case game.SevenCellValue:
		return '7'
	case game.EightCellValue:
		return '8'
	default:
		return '_'
	}
}
