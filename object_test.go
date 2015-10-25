package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseObjects(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		input   string
		want    []*Object
		wantErr error
	}{
		{
			desc: "valid cube",
			input: `
				# Blender v2.76 (sub 0) OBJ File: ''
				# www.blender.org
				o Cube
				v 1.000000 -1.000000 -1.000000
				v 1.000000 -1.000000 1.000000
				v -1.000000 -1.000000 1.000000
				v -1.000000 -1.000000 -1.000000
				v 1.000000 1.000000 -0.999999
				v 0.999999 1.000000 1.000001
				v -1.000000 1.000000 1.000000
				v -1.000000 1.000000 -1.000000
				s off
				f 1 2 3 4
				f 5 8 7 6
				f 1 5 6 2
				f 2 6 7 3
				f 3 7 8 4
				f 5 1 4 8`,
			want: []*Object{
				{
					vertices: []*ObjectVertex{
						{1, -1, -1},
						{1, -1, 1},
						{-1, -1, 1},
						{-1, -1, -1},
						{1, 1, -0.999999},
						{0.999999, 1, 1.000001},
						{-1, 1, 1},
						{-1, 1, -1},
					},
					faces: []*ObjectFace{
						{1, 2, 3, 4},
						{5, 8, 7, 6},
						{1, 5, 6, 2},
						{2, 6, 7, 3},
						{3, 7, 8, 4},
						{5, 1, 4, 8},
					},
				},
			},
		},
	} {
		got, gotErr := ParseObjects(strings.NewReader(tt.input))
		if !reflect.DeepEqual(got, tt.want) || !errorContains(gotErr, tt.wantErr) {
			t.Errorf("[%s] ParseObjects(%q) = (%v, %v), want (%v, %v)", tt.desc, tt.input, got, gotErr, tt.want, tt.wantErr)
		}
	}
}
