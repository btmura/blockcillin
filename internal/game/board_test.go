package game

import (
	"reflect"
	"testing"
)

func TestDropBlocks(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		board *Board
		want  *Board
	}{
		{
			desc: "drop from above",
			board: &Board{
				Rings: []*ring{
					{
						Cells: []*Cell{
							{Block: &Block{}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{State: BlockCleared}},
						},
					},
				},
				RingCount: 2,
				CellCount: 1,
			},
			want: &Board{
				Rings: []*ring{
					{
						Cells: []*Cell{
							{Block: &Block{State: BlockCleared}},
						},
					},
					{
						Cells: []*Cell{
							{Block: &Block{State: BlockDroppingFromAbove}},
						},
					},
				},
				RingCount: 2,
				CellCount: 1,
			},
		},
	} {
		tt.board.dropBlocks()
		if !reflect.DeepEqual(tt.board, tt.want) {
			t.Errorf("[%s] board.dropBlocks() -> %s, want %s", tt.desc, pp(tt.board), pp(tt.want))
		}
	}
}
