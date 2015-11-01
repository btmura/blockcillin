package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Obj struct {
	Vertices []*ObjVertex
	Faces    []*ObjFace
}

type ObjVertex struct {
	X float32
	Y float32
	Z float32
}

type ObjFace [4]int

func ReadObjFile(r io.Reader) ([]*Obj, error) {
	obj := &Obj{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "v"):
			v := &ObjVertex{}
			if _, err := fmt.Sscanf(line, "v %f %f %f", &v.X, &v.Y, &v.Z); err != nil {
				return nil, err
			}
			obj.Vertices = append(obj.Vertices, v)

		case strings.HasPrefix(line, "f"):
			f := &ObjFace{}
			if _, err := fmt.Sscanf(line, "f %d %d %d %d", &f[0], &f[1], &f[2], &f[3]); err != nil {
				return nil, err
			}
			obj.Faces = append(obj.Faces, f)
		}
	}
	return []*Obj{obj}, nil
}
