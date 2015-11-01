package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Obj struct {
	ID       string
	Vertices []*ObjVertex
	Faces    []*ObjFace
}

type ObjVertex struct {
	X float32
	Y float32
	Z float32
}

// ObjFace is a face described by vertex indices. Only triangles are supported.
type ObjFace [3]int

func ReadObjFile(r io.Reader) ([]*Obj, error) {
	var allObjs []*Obj
	var currentObj *Obj

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "o"):
			o, err := readObjObject(line)
			if err != nil {
				return nil, err
			}
			currentObj = o
			allObjs = append(allObjs, o)

		case strings.HasPrefix(line, "v"):
			v, err := readObjVertex(line)
			if err != nil {
				return nil, err
			}
			if currentObj == nil {
				return nil, errors.New("missing object ID")
			}
			currentObj.Vertices = append(currentObj.Vertices, v)

		case strings.HasPrefix(line, "f"):
			f, err := readObjFace(line)
			if err != nil {
				return nil, err
			}
			if currentObj == nil {
				return nil, errors.New("missing object ID")
			}
			currentObj.Faces = append(currentObj.Faces, f)
		}
	}

	return allObjs, nil
}

func readObjObject(line string) (*Obj, error) {
	o := &Obj{}
	if _, err := fmt.Sscanf(line, "o %s", &o.ID); err != nil {
		return nil, err
	}
	return o, nil
}

func readObjVertex(line string) (*ObjVertex, error) {
	v := &ObjVertex{}
	if _, err := fmt.Sscanf(line, "v %f %f %f", &v.X, &v.Y, &v.Z); err != nil {
		return nil, err
	}
	return v, nil
}

func readObjFace(line string) (*ObjFace, error) {
	f := &ObjFace{}
	if _, err := fmt.Sscanf(line, "f %d %d %d", &f[0], &f[1], &f[2]); err != nil {
		return nil, err
	}
	return f, nil
}
