package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
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

type ObjFace []int

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
	tokens := strings.Split(line, " ")
	if len(tokens) < 4 || tokens[0] != "f" {
		return nil, fmt.Errorf("invalid face spec: %q", line)
	}

	f := ObjFace{}
	for i, t := range tokens {
		if i == 0 {
			continue
		}
		vi, err := strconv.Atoi(t)
		if err != nil {
			return nil, err
		}
		f = append(f, vi)
	}
	return &f, nil
}
