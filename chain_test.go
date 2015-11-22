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
			desc: "horizontal 3 in a row",
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
			},
			want: []*chain{
				{
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
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
