package main

import (
	"reflect"
	"testing"
)

func TestFindHorizontalChains(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		input *board
		want  []*chain
	}{
		{
			desc: "first 3 matches",
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
					cells: []*chainCell{
						{0, 0},
						{1, 0},
						{2, 0},
					},
				},
			},
		},
		{
			desc: "last 3 matches",
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
				ringCount: 1,
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
			desc: "whole ring matches",
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
			desc: "no match due to clearing block",
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
				ringCount: 1,
				cellCount: 5,
			},
		},
		{
			desc: "no match due to invisible block",
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
			desc: "first 3 matches",
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
					cells: []*chainCell{
						{0, 0},
						{0, 1},
						{0, 2},
					},
				},
			},
		},
	} {
		got := findVerticalChains(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("[%s] findVerticalChains(%s) = %s, want %s", tt.desc, pp(tt.input), pp(got), pp(tt.want))
		}
	}
}
