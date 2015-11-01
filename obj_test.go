package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestReadObjFile(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		input   string
		want    []*Obj
		wantErr error
	}{
		{
			desc: "missing object ID",
			input: `
				# Blender v2.76 (sub 0) OBJ File: ''
				# www.blender.org
				v -0.032120 -0.290752 -0.832947
				v 1.967880 -0.290752 -0.832947
				v -0.032120 -0.290752 -2.832947
				s off
				f 1 2 4
			`,
			wantErr: errors.New("missing object ID"),
		},
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
				f 2 3 4
				f 8 7 6
				f 5 6 2
				f 6 7 3
				f 3 7 8
				f 1 4 8
				f 1 2 4
				f 5 8 6
				f 1 5 2
				f 2 6 3
				f 4 3 8
				f 5 1 8
			`,
			want: []*Obj{
				{
					ID: "Cube",
					Vertices: []*ObjVertex{
						{1, -1, -1},
						{1, -1, 1},
						{-1, -1, 1},
						{-1, -1, -1},
						{1, 1, -0.999999},
						{0.999999, 1, 1.000001},
						{-1, 1, 1},
						{-1, 1, -1},
					},
					Faces: []*ObjFace{
						{2, 3, 4},
						{8, 7, 6},
						{5, 6, 2},
						{6, 7, 3},
						{3, 7, 8},
						{1, 4, 8},
						{1, 2, 4},
						{5, 8, 6},
						{1, 5, 2},
						{2, 6, 3},
						{4, 3, 8},
						{5, 1, 8},
					},
				},
			},
		},
		{
			desc: "multiple objects",
			input: `
				# Blender v2.76 (sub 0) OBJ File: ''
				# www.blender.org
				o Plane.001
				v 0.652447 0.140019 -0.450452
				v 2.652447 0.140019 -0.450452
				v 0.652447 0.140019 -2.450452
				v 2.652447 0.140019 -2.450452
				s off
				f 2 4 3
				f 1 2 3
				o Plane
				v -1.079860 0.672774 2.814899
				v 0.920140 0.672774 2.814899
				v -1.079860 0.672774 0.814900
				v 0.920140 0.672774 0.814900
				s off
				f 6 8 7
				f 5 6 7
			`,
			want: []*Obj{
				{
					ID: "Plane.001",
					Vertices: []*ObjVertex{
						{0.652447, 0.140019, -0.450452},
						{2.652447, 0.140019, -0.450452},
						{0.652447, 0.140019, -2.450452},
						{2.652447, 0.140019, -2.450452},
					},
					Faces: []*ObjFace{
						{2, 4, 3},
						{1, 2, 3},
					},
				},
				{
					ID: "Plane",
					Vertices: []*ObjVertex{
						{-1.079860, 0.672774, 2.814899},
						{0.920140, 0.672774, 2.814899},
						{-1.079860, 0.672774, 0.814900},
						{0.920140, 0.672774, 0.814900},
					},
					Faces: []*ObjFace{
						{6, 8, 7},
						{5, 6, 7},
					},
				},
			},
		},
	} {
		got, gotErr := ReadObjFile(strings.NewReader(tt.input))
		if !reflect.DeepEqual(got, tt.want) || !errorContains(gotErr, tt.wantErr) {
			t.Errorf("[%s] ReadObjFile(%q) = (%v, %v), want (%v, %v)", tt.desc, tt.input, pretty.Sprint(got), gotErr, pretty.Sprint(tt.want), tt.wantErr)
		}
	}
}
