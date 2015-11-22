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
			desc: "first 3 matches horizontally",
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
				cellCount: 4,
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
		{
			desc: "last 3 matches horizontally",
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: green}},
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red}},
						},
					},
				},
				cellCount: 4,
			},
			want: []*chain{
				{
					cells: []*chainCell{
						{1, 0},
						{2, 0},
						{3, 0},
					},
				},
			},
		},
		{
			desc: "wrap matches horizontally",
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
				cellCount: 4,
			},
			want: []*chain{
				{
					cells: []*chainCell{
						{3, 0},
						{0, 0},
						{1, 0},
					},
				},
			},
		},
		{
			desc: "multiple matches horizontally",
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
				cellCount: 7,
			},
			want: []*chain{
				{
					cells: []*chainCell{
						{4, 0},
						{5, 0},
						{6, 0},
					},
				},
				{
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
					},
				},
			},
		},
		{
			desc: "whole ring matches horizontally",
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
				cellCount: 4,
			},
			want: []*chain{
				{
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
			desc: "no match due to cleared block",
			input: &board{
				rings: []*ring{
					{
						cells: []*cell{
							{block: &block{color: red}},
							{block: &block{color: red}},
							{block: &block{color: red, state: blockClearing}},
							{block: &block{color: red}},
							{block: &block{color: green}},
						},
					},
				},
				cellCount: 5,
			},
		},
	} {
		got := findChains(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("[%s] findChains(%s) = %s, want %s", tt.desc, pp(tt.input), pp(got), pp(tt.want))
		}
	}
}
