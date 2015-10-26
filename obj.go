package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Obj struct {
	vertices []*ObjVertex
	faces    []*ObjFace
}

type ObjVertex struct {
	x float32
	y float32
	z float32
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
			if _, err := fmt.Sscanf(line, "v %f %f %f", &v.x, &v.y, &v.z); err != nil {
				return nil, err
			}
			obj.vertices = append(obj.vertices, v)

		case strings.HasPrefix(line, "f"):
			f := &ObjFace{}
			if _, err := fmt.Sscanf(line, "f %d %d %d %d", &f[0], &f[1], &f[2], &f[3]); err != nil {
				return nil, err
			}
			obj.faces = append(obj.faces, f)
		}
	}
	return []*Obj{obj}, nil
}
