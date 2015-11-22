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
	}{} {
		got := findDrops(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("[%s] findDrops(%s) = %s, want %s", tt.desc, pp(tt.input), pp(got), pp(tt.want))
		}
	}
}
