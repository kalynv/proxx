package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_GetState(t *testing.T) {
	game := &Game{
		failAt: nil,
		board: [][]Cell{
			{{Content: OneCellValue}},
			{{Content: TwoCellValue}, {Content: ThreeCellValue}},
			{{Content: FourCellValue}, {Content: FiveCellValue}, {Content: SixCellValue}},
		},
	}

	got := game.GetState()

	if len(got) != len(game.board) {
		t.Errorf("Want board length [%d], got [%d]", len(game.board), len(got))

		return
	}

	for i := range got {
		if len(got[i]) != len(game.board[i]) {
			t.Errorf("Want board row %d length [%d], got [%d]", i, len(game.board[i]), len(got[i]))

			continue
		}

		for j := range got[i] {
			if got[i][j] != game.board[i][j] {
				t.Errorf("Want cell(%d, %d) [%v], got [%v]", i, j, game.board[i][j], got[i][j])
			}

			// ensure values are cloned
			if &got[i][j] == &game.board[i][j] {
				t.Errorf("The same cell(%d, %d) instance", i, j)
			}
		}
	}
}

func TestGame_Lost(t *testing.T) {
	tests := []struct {
		name string
		game *Game
		want bool
	}{
		{
			name: "lost game",
			game: &Game{
				failAt: &cellAddress{row: 0, column: 0},
				board:  [][]Cell{{{Content: BlackHoleCellValue, State: HiddenState}}},
			},
			want: true,
		},
		{
			name: "game is not lost",
			game: &Game{
				failAt: &cellAddress{row: 0, column: 0},
				board:  [][]Cell{{{Content: BlackHoleCellValue, State: HiddenState}}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.game.Lost()

			if tt.want != got {
				t.Errorf("Want [%v], got [%v]", tt.want, got)
			}
		})
	}
}

func TestGame_Won(t *testing.T) {
	tests := []struct {
		name string
		game *Game
		want bool
	}{
		{
			name: "game is lost",
			game: &Game{
				failAt: &cellAddress{row: 0, column: 0},
				board: [][]Cell{
					{{Content: BlackHoleCellValue, State: HiddenState}, {Content: BlackHoleCellValue, State: HiddenState}},
					{{Content: TwoCellValue, State: HiddenState}, {Content: TwoCellValue, State: HiddenState}},
				},
			},
			want: false,
		},
		{
			name: "game is not completed",
			game: &Game{
				failAt: nil,
				board: [][]Cell{
					{{Content: BlackHoleCellValue, State: HiddenState}, {Content: BlackHoleCellValue, State: HiddenState}},
					{{Content: TwoCellValue, State: VisibleState}, {Content: TwoCellValue, State: HiddenState}},
				},
			},
			want: false,
		},
		{
			name: "all non black hole cells are visible",
			game: &Game{
				failAt: nil,
				board: [][]Cell{
					{{Content: BlackHoleCellValue, State: HiddenState}, {Content: BlackHoleCellValue, State: HiddenState}},
					{{Content: TwoCellValue, State: VisibleState}, {Content: TwoCellValue, State: VisibleState}},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.game.Won()

			if tt.want != got {
				t.Errorf("Want [%v], got [%v]", tt.want, got)
			}
		})
	}
}

func TestGame_RevealCell(t *testing.T) {
	tests := []struct {
		name             string
		game             *Game
		invokeRevealCell func(g *Game)
		wantBoard        [][]Cell
		wantLost         bool
		wantWon          bool
		wantCompleted    bool
	}{
		{
			name: "lost right away",
			game: func() *Game {
				board := createGameState(3, 3, Cell{Content: ZeroCellValue, State: HiddenState})
				blackHoleAddresses := []cellAddress{
					{row: 0, column: 0},
					{row: 1, column: 1},
					{row: 2, column: 0},
				}
				replaceCells(board, blackHoleAddresses, Cell{Content: BlackHoleCellValue, State: HiddenState})
				updateNaboringBlackHolesCellValues(board)

				return &Game{failAt: nil, board: board}
			}(),
			invokeRevealCell: func(g *Game) {
				g.RevealCell(1, 1)
			},
			wantBoard: [][]Cell{
				{
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: TwoCellValue, State: HiddenState},
					{Content: OneCellValue, State: HiddenState},
				},
				{
					{Content: ThreeCellValue, State: HiddenState},
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: OneCellValue, State: HiddenState},
				},
				{
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: TwoCellValue, State: HiddenState},
					{Content: OneCellValue, State: HiddenState},
				},
			},
			wantLost:      true,
			wantWon:       false,
			wantCompleted: true,
		},
		{
			name: "reveal cell",
			game: func() *Game {
				board := createGameState(3, 3, Cell{Content: ZeroCellValue, State: HiddenState})
				blackHoleAddresses := []cellAddress{
					{row: 0, column: 0},
					{row: 1, column: 1},
					{row: 2, column: 0},
				}
				replaceCells(board, blackHoleAddresses, Cell{Content: BlackHoleCellValue, State: HiddenState})
				updateNaboringBlackHolesCellValues(board)

				return &Game{failAt: nil, board: board}
			}(),
			invokeRevealCell: func(g *Game) {
				g.RevealCell(1, 0)
			},
			wantBoard: [][]Cell{
				{
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: TwoCellValue, State: HiddenState},
					{Content: OneCellValue, State: HiddenState},
				},
				{
					{Content: ThreeCellValue, State: VisibleState},
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: OneCellValue, State: HiddenState},
				},
				{
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: TwoCellValue, State: HiddenState},
					{Content: OneCellValue, State: HiddenState},
				},
			},
			wantLost:      false,
			wantWon:       false,
			wantCompleted: false,
		},
		{
			name: "won on revealing",
			game: func() *Game {
				board := createGameState(3, 3, Cell{Content: ZeroCellValue, State: HiddenState})
				blackHoleAddresses := []cellAddress{
					{row: 0, column: 0},
					{row: 1, column: 1},
					{row: 2, column: 0},
				}
				replaceCells(board, blackHoleAddresses, Cell{Content: BlackHoleCellValue, State: HiddenState})
				updateNaboringBlackHolesCellValues(board)

				return &Game{failAt: nil, board: board}
			}(),
			invokeRevealCell: func(g *Game) {
				g.RevealCell(0, 1)
				g.RevealCell(0, 2)
				g.RevealCell(1, 0)
				g.RevealCell(1, 2)
				g.RevealCell(2, 1)
				g.RevealCell(2, 2)
			},
			wantBoard: [][]Cell{
				{
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: TwoCellValue, State: VisibleState},
					{Content: OneCellValue, State: VisibleState},
				},
				{
					{Content: ThreeCellValue, State: VisibleState},
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: OneCellValue, State: VisibleState},
				},
				{
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: TwoCellValue, State: VisibleState},
					{Content: OneCellValue, State: VisibleState},
				},
			},
			wantLost:      false,
			wantWon:       true,
			wantCompleted: true,
		},
		{
			name: "open contiguos space of cells not bordering with black holes",
			game: func() *Game {
				board := createGameState(3, 3, Cell{Content: ZeroCellValue, State: HiddenState})
				blackHoleAddresses := []cellAddress{
					{row: 0, column: 0},
				}
				replaceCells(board, blackHoleAddresses, Cell{Content: BlackHoleCellValue, State: HiddenState})
				updateNaboringBlackHolesCellValues(board)

				return &Game{failAt: nil, board: board}
			}(),
			invokeRevealCell: func(g *Game) {
				g.RevealCell(0, 2)
			},
			wantBoard: [][]Cell{
				{
					{Content: BlackHoleCellValue, State: HiddenState},
					{Content: OneCellValue, State: VisibleState},
					{Content: ZeroCellValue, State: VisibleState},
				},
				{
					{Content: OneCellValue, State: VisibleState},
					{Content: OneCellValue, State: VisibleState},
					{Content: ZeroCellValue, State: VisibleState},
				},
				{
					{Content: ZeroCellValue, State: VisibleState},
					{Content: ZeroCellValue, State: VisibleState},
					{Content: ZeroCellValue, State: VisibleState},
				},
			},
			wantLost:      false,
			wantWon:       true,
			wantCompleted: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invokeRevealCell(tt.game)

			gotBoard := tt.game.GetState()
			gotLost := tt.game.Lost()
			gotWon := tt.game.Won()
			gotCompleted := tt.game.Completed()

			assert.Equal(t, tt.wantBoard, gotBoard)
			assert.Equal(t, tt.wantLost, gotLost)
			assert.Equal(t, tt.wantWon, gotWon)
			assert.Equal(t, tt.wantCompleted, gotCompleted)
		})
	}
}

func TestGame_Completed(t *testing.T) {
	tests := []struct {
		name string
		game *Game
		want bool
	}{
		{
			name: "game is lost",
			game: &Game{
				failAt: &cellAddress{row: 0, column: 0},
				board: [][]Cell{
					{{Content: BlackHoleCellValue, State: HiddenState}, {Content: BlackHoleCellValue, State: HiddenState}},
					{{Content: TwoCellValue, State: HiddenState}, {Content: TwoCellValue, State: HiddenState}},
				},
			},
			want: true,
		},
		{
			name: "game is not completed",
			game: &Game{
				failAt: nil,
				board: [][]Cell{
					{{Content: BlackHoleCellValue, State: HiddenState}, {Content: BlackHoleCellValue, State: HiddenState}},
					{{Content: TwoCellValue, State: VisibleState}, {Content: TwoCellValue, State: HiddenState}},
				},
			},
			want: false,
		},
		{
			name: "all non black hole cells are visible",
			game: &Game{
				failAt: nil,
				board: [][]Cell{
					{{Content: BlackHoleCellValue, State: HiddenState}, {Content: BlackHoleCellValue, State: HiddenState}},
					{{Content: TwoCellValue, State: VisibleState}, {Content: TwoCellValue, State: VisibleState}},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.game.Completed()

			if tt.want != got {
				t.Errorf("Want [%v], got [%v]", tt.want, got)
			}
		})
	}
}

func Test_calculateNaboringBlackHoles(t *testing.T) {
	tests := []struct {
		name        string
		board       [][]Cell
		cellAddress cellAddress
		want        CellValue
	}{
		{
			name: "cell with no surrounding black holes",
			board: [][]Cell{
				{{Content: ZeroCellValue}, {Content: ZeroCellValue}, {Content: ZeroCellValue}},
				{{Content: ZeroCellValue}, {Content: ZeroCellValue}, {Content: ZeroCellValue}},
				{{Content: ZeroCellValue}, {Content: ZeroCellValue}, {Content: ZeroCellValue}},
			},
			cellAddress: cellAddress{row: 0, column: 0},
			want:        ZeroCellValue,
		},
		{
			name: "cell with one surrounding black hole",
			board: [][]Cell{
				{{Content: BlackHoleCellValue}, {Content: ZeroCellValue}, {Content: ZeroCellValue}},
				{{Content: ZeroCellValue}, {Content: ZeroCellValue}, {Content: ZeroCellValue}},
				{{Content: ZeroCellValue}, {Content: ZeroCellValue}, {Content: ZeroCellValue}},
			},
			cellAddress: cellAddress{row: 0, column: 1},
			want:        OneCellValue,
		},
		{
			name: "cell with two surrounding black holes",
			board: [][]Cell{
				{{Content: ZeroCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
				{{Content: ZeroCellValue}, {Content: ZeroCellValue}, {Content: BlackHoleCellValue}},
				{{Content: ZeroCellValue}, {Content: ZeroCellValue}, {Content: ZeroCellValue}},
			},
			cellAddress: cellAddress{row: 0, column: 2},
			want:        TwoCellValue,
		},
		{
			name: "cell with three surrounding black holes",
			board: [][]Cell{
				{{Content: ZeroCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
				{{Content: ZeroCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
			},
			cellAddress: cellAddress{row: 1, column: 2},
			want:        ThreeCellValue,
		},
		{
			name: "cell with four surrounding black holes",
			board: [][]Cell{
				{{Content: ZeroCellValue}, {Content: OneCellValue}, {Content: ZeroCellValue}},
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}},
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
			},
			cellAddress: cellAddress{row: 2, column: 1},
			want:        FourCellValue,
		},
		{
			name: "cell with five surrounding black holes",
			board: [][]Cell{
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
				{{Content: ZeroCellValue}, {Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}},
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
			},
			cellAddress: cellAddress{row: 1, column: 0},
			want:        FiveCellValue,
		},
		{
			name: "cell with six surrounding black holes",
			board: [][]Cell{
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
				{{Content: BlackHoleCellValue}, {Content: ZeroCellValue}, {Content: BlackHoleCellValue}},
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: ZeroCellValue}},
			},
			cellAddress: cellAddress{row: 1, column: 1},
			want:        SixCellValue,
		},
		{
			name: "cell with seven surrounding black holes",
			board: [][]Cell{
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}},
				{{Content: BlackHoleCellValue}, {Content: ZeroCellValue}, {Content: BlackHoleCellValue}},
				{{Content: BlackHoleCellValue}, {Content: ZeroCellValue}, {Content: BlackHoleCellValue}},
			},
			cellAddress: cellAddress{row: 1, column: 1},
			want:        SevenCellValue,
		},
		{
			name: "cell with eight surrounding black holes",
			board: [][]Cell{
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}},
				{{Content: BlackHoleCellValue}, {Content: ZeroCellValue}, {Content: BlackHoleCellValue}},
				{{Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}, {Content: BlackHoleCellValue}},
			},
			cellAddress: cellAddress{row: 1, column: 1},
			want:        EightCellValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateNaboringBlackHoles(tt.board, tt.cellAddress)

			if tt.want != got {
				t.Errorf("Want: [%d], got: [%d]", tt.want, got)
			}
		})
	}
}

func Fuzz_generateBlackHoleAddresses(f *testing.F) {
	f.Add(5, 7, 9)
	f.Fuzz(func(t *testing.T, row, column, amount int) {
		if row < 0 || column < 0 || amount < 0 {
			return
		}
		if amount > row*column {
			return
		}
		got := generateBlackHoleAddresses(row, column, amount)
		if len(got) != amount {
			t.Errorf("Want generated black hole addresses: %d, got: %d", amount, len(got))
		}
		for _, address := range got {
			if address.row < 0 || address.row >= row || address.column < 0 || address.column >= column {
				t.Errorf(
					"Want address row in [%d,%d) and column in [%d,%d), got [%d, %d]",
					0,
					row,
					0,
					column,
					address.row,
					address.column,
				)
			}
		}

		countMap := make(map[cellAddress]int)
		for _, a := range got {
			countMap[a]++
		}
		for k, v := range countMap {
			if v != 1 {
				t.Errorf("Want one uniq address generated, got address [%d, %d] generated %d times", k.row, k.column, v)
			}
		}
	})
}
