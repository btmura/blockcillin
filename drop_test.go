package main

import (
	"reflect"
	"testing"
)

func TestFindDrops(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		input *board
		want  []*drop
	}{
		{
			desc: "drop from above",
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{}},
						},
					},
					{
						cells: []*cell{
							{block: &block{invisible: true}},
						},
					},
				},
				ringCount: 2,
				cellCount: 1,
			},
			want: []*drop{
				{0, 0},
			},
		},
	} {
		got := findDrops(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("[%s] findDrops(%s) = %s, want %s", tt.desc, pp(tt.input), pp(got), pp(tt.want))
		}
	}
}
