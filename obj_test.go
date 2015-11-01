package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadObjFile(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		input   string
		want    []*Obj
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
			want: []*Obj{
				{
					vertices: []*ObjVertex{
						{1, -1, -1},
						{1, -1, 1},
						{-1, -1, 1},
						{-1, -1, -1},
						{1, 1, -0.999999},
						{0.999999, 1, 1.000001},
						{-1, 1, 1},
						{-1, 1, -1},
					},
					faces: []*ObjFace{
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
		{
			desc: "valid cone",
			input: `
				# Blender v2.76 (sub 0) OBJ File: ''
				# www.blender.org
				v 0.586231 -0.022019 -0.000024
				v 0.586231 1.977981 0.999976
				v 0.781322 -0.022019 0.019190
				v 0.968915 -0.022019 0.076096
				v 1.141801 -0.022019 0.168506
				v 1.293338 -0.022019 0.292869
				v 1.417701 -0.022019 0.444405
				v 1.510111 -0.022019 0.617292
				v 1.567016 -0.022019 0.804885
				v 1.586231 -0.022019 0.999976
				v 1.567016 -0.022019 1.195066
				v 1.510111 -0.022019 1.382659
				v 1.417701 -0.022019 1.555546
				v 1.293338 -0.022019 1.707082
				v 1.141801 -0.022019 1.831445
				v 0.968914 -0.022019 1.923855
				v 0.781321 -0.022019 1.980761
				v 0.586231 -0.022019 1.999976
				v 0.391140 -0.022019 1.980761
				v 0.203547 -0.022019 1.923855
				v 0.030660 -0.022019 1.831445
				v -0.120876 -0.022019 1.707082
				v -0.245239 -0.022019 1.555545
				v -0.337649 -0.022019 1.382658
				v -0.394554 -0.022019 1.195065
				v -0.413769 -0.022019 0.999975
				v -0.394554 -0.022019 0.804884
				v -0.337648 -0.022019 0.617291
				v -0.245238 -0.022019 0.444404
				v -0.120875 -0.022019 0.292868
				v 0.030662 -0.022019 0.168505
				v 0.203549 -0.022019 0.076096
				v 0.391142 -0.022019 0.019190
				s off
				f 32 2 33
				f 1 2 3
				f 31 2 32
				f 30 2 31
				f 29 2 30
				f 28 2 29
				f 27 2 28
				f 26 2 27
				f 25 2 26
				f 24 2 25
				f 23 2 24
				f 22 2 23
				f 21 2 22
				f 20 2 21
				f 19 2 20
				f 18 2 19
				f 17 2 18
				f 16 2 17
				f 15 2 16
				f 14 2 15
				f 13 2 14
				f 12 2 13
				f 11 2 12
				f 10 2 11
				f 9 2 10
				f 8 2 9
				f 7 2 8
				f 6 2 7
				f 5 2 6
				f 4 2 5
				f 33 2 1
				f 3 2 4
				f 1 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33`,
		},
	} {
		got, gotErr := ReadObjFile(strings.NewReader(tt.input))
		if !reflect.DeepEqual(got, tt.want) || !errorContains(gotErr, tt.wantErr) {
			t.Errorf("[%s] ReadObjFile(%q) = (%v, %v), want (%v, %v)", tt.desc, tt.input, got, gotErr, tt.want, tt.wantErr)
		}
	}
}
