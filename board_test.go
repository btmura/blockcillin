package main

import (
	"reflect"
	"testing"
)

func TestDropBlocks(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		board *board
		want  *board
	}{
		{
			desc: "drop from above",
			board: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{}},
						},
					},
					{
						cells: []*cell{
							{block: &block{state: blockCleared}},
						},
					},
				},
				ringCount: 2,
				cellCount: 1,
			},
			want: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{state: blockCleared}},
						},
					},
					{
						cells: []*cell{
							{block: &block{state: blockDroppingFromAbove}},
						},
					},
				},
				ringCount: 2,
				cellCount: 1,
			},
		},
	} {
		tt.board.dropBlocks()
		if !reflect.DeepEqual(tt.board, tt.want) {
			t.Errorf("[%s] board.dropBlocks() -> %s, want %s", tt.desc, pp(tt.board), pp(tt.want))
		}
	}
}
