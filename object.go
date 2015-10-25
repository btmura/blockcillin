package main

import "io"

type Object struct {
	vertices []*ObjectVertex
}

type ObjectVertex struct {
	x float32
	y float32
	z float32
}

func ParseObjectSource(r io.Reader) ([]*Object, error) {
	return nil, nil
}
