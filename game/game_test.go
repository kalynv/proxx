package game

import "testing"

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
