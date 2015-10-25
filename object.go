package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Object struct {
	vertices []*ObjectVertex
}

type ObjectVertex struct {
	x float32
	y float32
	z float32
}

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
		}
	}
	return []*Object{obj}, nil
}
