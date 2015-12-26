package game

import (
	"reflect"
	"testing"
)

func TestFindChains(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		input *Board
		want  []*chain
	}{
		{
			desc: "cross",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Green}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Green}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{Color: Green}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Green}},
						},
					},
				},
				RingCount: 3,
				CellCount: 3,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{1, 0},
						{0, 1},
						{1, 1},
						{2, 1},
						{1, 2},
					},
				},
			},
		},
		{
			desc: "square",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
				},
				RingCount: 3,
				CellCount: 3,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
						{0, 1},
						{1, 1},
						{2, 1},
						{0, 2},
						{1, 2},
						{2, 2},
					},
				},
			},
		},
	} {
		got := findChains(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("[%s] findChains(%s) = %s, want %s", tt.desc, pp(tt.input), pp(got), pp(tt.want))
		}
	}
}

func TestFindHorizontalChains(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		input *Board
		want  []*chain
	}{
		{
			desc: "first 3 horizontal match",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Green}},
						},
					},
				},
				RingCount: 1,
				CellCount: 4,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
					},
				},
			},
		},
		{
			desc: "last 4 horizontal match",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Green}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
				},
				RingCount: 1,
				CellCount: 5,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{1, 0},
						{2, 0},
						{3, 0},
						{4, 0},
					},
				},
			},
		},
		{
			desc: "wrap matches",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Green}},
							{Block: &Block{Color: Red}},
						},
					},
				},
				RingCount: 1,
				CellCount: 4,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{3, 0},
						{0, 0},
						{1, 0},
					},
				},
			},
		},
		{
			desc: "multiple matches",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Green}},
							{Block: &Block{Color: Blue}},
							{Block: &Block{Color: Blue}},
							{Block: &Block{Color: Blue}},
						},
					},
				},
				RingCount: 1,
				CellCount: 7,
			},
			want: []*chain{
				{
					color: Blue,
					cells: []*chainCell{
						{4, 0},
						{5, 0},
						{6, 0},
					},
				},
				{
					color: Red,
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
					},
				},
			},
		},
		{
			desc: "whole row matches",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
				},
				RingCount: 1,
				CellCount: 4,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
						{3, 0},
					},
				},
			},
		},
		{
			desc: "square",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
				},
				RingCount: 3,
				CellCount: 3,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
					},
				},
				{
					color: Red,
					cells: []*chainCell{
						{0, 1},
						{1, 1},
						{2, 1},
					},
				},
				{
					color: Red,
					cells: []*chainCell{
						{0, 2},
						{1, 2},
						{2, 2},
					},
				},
			},
		},
		{
			desc: "no match due to flashing blocks",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red, State: BlockFlashing}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Green}},
						},
					},
				},
				RingCount: 1,
				CellCount: 5,
			},
		},
		{
			desc: "no match due to invisible blocks",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red, State: BlockCleared}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Green}},
						},
					},
				},
				RingCount: 1,
				CellCount: 5,
			},
		},
	} {
		got := findHorizontalChains(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("[%s] findHorizontalChains(%s) = %s, want %s", tt.desc, pp(tt.input), pp(got), pp(tt.want))
		}
	}
}

func TestFindVerticalChains(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		input *Board
		want  []*chain
	}{
		{
			desc: "first 3 vertical match",
			input: &Board{
				Rings: []*Ring{
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
				},
				RingCount: 4,
				CellCount: 1,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{0, 0},
						{0, 1},
						{0, 2},
					},
				},
			},
		},
		{
			desc: "last 4 vertical match",
			input: &Board{
				Rings: []*Ring{
					{Cells: []*Cell{{Block: &Block{Color: Blue}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
				},
				RingCount: 5,
				CellCount: 1,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{0, 1},
						{0, 2},
						{0, 3},
						{0, 4},
					},
				},
			},
		},
		{
			desc: "multiple matches",
			input: &Board{
				Rings: []*Ring{
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Blue}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
					{Cells: []*Cell{{Block: &Block{Color: Red}}}},
				},
				RingCount: 8,
				CellCount: 1,
			},
			want: []*chain{
				{
					color: Green,
					cells: []*chainCell{
						{0, 0},
						{0, 1},
						{0, 2},
						{0, 3},
					},
				},
				{
					color: Red,
					cells: []*chainCell{
						{0, 5},
						{0, 6},
						{0, 7},
					},
				},
			},
		},
		{
			desc: "whole column matches",
			input: &Board{
				Rings: []*Ring{
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
				},
				RingCount: 4,
				CellCount: 1,
			},
			want: []*chain{
				{
					color: Green,
					cells: []*chainCell{
						{0, 0},
						{0, 1},
						{0, 2},
						{0, 3},
					},
				},
			},
		},

		{
			desc: "square",
			input: &Board{
				Rings: []*Ring{
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
							{Block: &Block{Color: Red}},
						},
					},
				},
				RingCount: 3,
				CellCount: 3,
			},
			want: []*chain{
				{
					color: Red,
					cells: []*chainCell{
						{0, 0},
						{0, 1},
						{0, 2},
					},
				},
				{
					color: Red,
					cells: []*chainCell{
						{1, 0},
						{1, 1},
						{1, 2},
					},
				},
				{
					color: Red,
					cells: []*chainCell{
						{2, 0},
						{2, 1},
						{2, 2},
					},
				},
			},
		},
		{
			desc: "no match due to flashing block",
			input: &Board{
				Rings: []*Ring{
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green, State: BlockFlashing}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
				},
				RingCount: 4,
				CellCount: 1,
			},
		},
		{
			desc: "no match due to clearing block",
			input: &Board{
				Rings: []*Ring{
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green, State: BlockCleared}}}},
					{Cells: []*Cell{{Block: &Block{Color: Green}}}},
				},
				RingCount: 4,
				CellCount: 1,
			},
		},
	} {
		got := findVerticalChains(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("[%s] findVerticalChains(%s) = %s, want %s", tt.desc, pp(tt.input), pp(got), pp(tt.want))
		}
	}
}
