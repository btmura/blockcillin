package main

import (
	"reflect"
	"testing"
)

func TestFindChains(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		input *board
		want  []*chain
	}{
		{
			desc: "cross",
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: green}},
							{block: &block{color: red}},
							{block: &block{color: green}},
						},
					},
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
					{
						cells: []*cell{
							{block: &block{color: green}},
							{block: &block{color: red}},
							{block: &block{color: green}},
						},
					},
				},
				ringCount: 3,
				cellCount: 3,
			},
			want: []*chain{
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
				},
				ringCount: 3,
				cellCount: 3,
			},
			want: []*chain{
				{
					color: red,
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
		input *board
		want  []*chain
	}{
		{
			desc: "first 3 horizontal match",
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: green}},
						},
					},
				},
				ringCount: 1,
				cellCount: 4,
			},
			want: []*chain{
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: green}},
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
				},
				ringCount: 1,
				cellCount: 5,
			},
			want: []*chain{
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: green}},
							{block: &block{color: red}},
						},
					},
				},
				ringCount: 1,
				cellCount: 4,
			},
			want: []*chain{
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: green}},
							{block: &block{color: blue}},
							{block: &block{color: blue}},
							{block: &block{color: blue}},
						},
					},
				},
				ringCount: 1,
				cellCount: 7,
			},
			want: []*chain{
				{
					color: blue,
					cells: []*chainCell{
						{4, 0},
						{5, 0},
						{6, 0},
					},
				},
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
				},
				ringCount: 1,
				cellCount: 4,
			},
			want: []*chain{
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
				},
				ringCount: 3,
				cellCount: 3,
			},
			want: []*chain{
				{
					color: red,
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
					},
				},
				{
					color: red,
					cells: []*chainCell{
						{0, 1},
						{1, 1},
						{2, 1},
					},
				},
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red, state: blockFlashing}},
							{block: &block{color: red}},
							{block: &block{color: green}},
						},
					},
				},
				ringCount: 1,
				cellCount: 5,
			},
		},
		{
			desc: "no match due to invisible blocks",
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red, invisible: true}},
							{block: &block{color: red}},
							{block: &block{color: green}},
						},
					},
				},
				ringCount: 1,
				cellCount: 5,
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
		input *board
		want  []*chain
	}{
		{
			desc: "first 3 vertical match",
			input: &board{
				rings: []*ring{
					{cells: []*cell{{block: &block{color: red}}}},
					{cells: []*cell{{block: &block{color: red}}}},
					{cells: []*cell{{block: &block{color: red}}}},
					{cells: []*cell{{block: &block{color: green}}}},
				},
				ringCount: 4,
				cellCount: 1,
			},
			want: []*chain{
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{cells: []*cell{{block: &block{color: blue}}}},
					{cells: []*cell{{block: &block{color: red}}}},
					{cells: []*cell{{block: &block{color: red}}}},
					{cells: []*cell{{block: &block{color: red}}}},
					{cells: []*cell{{block: &block{color: red}}}},
				},
				ringCount: 5,
				cellCount: 1,
			},
			want: []*chain{
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: blue}}}},
					{cells: []*cell{{block: &block{color: red}}}},
					{cells: []*cell{{block: &block{color: red}}}},
					{cells: []*cell{{block: &block{color: red}}}},
				},
				ringCount: 8,
				cellCount: 1,
			},
			want: []*chain{
				{
					color: green,
					cells: []*chainCell{
						{0, 0},
						{0, 1},
						{0, 2},
						{0, 3},
					},
				},
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green}}}},
				},
				ringCount: 4,
				cellCount: 1,
			},
			want: []*chain{
				{
					color: green,
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
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
				},
				ringCount: 3,
				cellCount: 3,
			},
			want: []*chain{
				{
					color: red,
					cells: []*chainCell{
						{0, 0},
						{0, 1},
						{0, 2},
					},
				},
				{
					color: red,
					cells: []*chainCell{
						{1, 0},
						{1, 1},
						{1, 2},
					},
				},
				{
					color: red,
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
			input: &board{
				rings: []*ring{
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green, state: blockFlashing}}}},
					{cells: []*cell{{block: &block{color: green}}}},
				},
				ringCount: 4,
				cellCount: 1,
			},
		},
		{
			desc: "no match due to clearing block",
			input: &board{
				rings: []*ring{
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green}}}},
					{cells: []*cell{{block: &block{color: green, invisible: true}}}},
					{cells: []*cell{{block: &block{color: green}}}},
				},
				ringCount: 4,
				cellCount: 1,
			},
		},
	} {
		got := findVerticalChains(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("[%s] findVerticalChains(%s) = %s, want %s", tt.desc, pp(tt.input), pp(got), pp(tt.want))
		}
	}
}
