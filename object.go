package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Object struct {
	vertices []*ObjectVertex
	faces    []*ObjectFace
}

type ObjectVertex struct {
	x float32
	y float32
	z float32
}

type ObjectFace [4]int

func ParseObjectSource(r io.Reader) ([]*Object, error) {
	obj := &Object{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "v"):
			v := &ObjectVertex{}
			if _, err := fmt.Sscanf(line, "v %f %f %f", &v.x, &v.y, &v.z); err != nil {
				return nil, err
			}
			obj.vertices = append(obj.vertices, v)

		case strings.HasPrefix(line, "f"):
			f := &ObjectFace{}
			if _, err := fmt.Sscanf(line, "f %d %d %d %d", &f[0], &f[1], &f[2], &f[3]); err != nil {
				return nil, err
			}
			obj.faces = append(obj.faces, f)
		}
	}
	return []*Object{obj}, nil
}
