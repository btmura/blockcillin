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
	ID        string
	Vertices  []*ObjVertex
	TexCoords []*ObjTexCoord
	Faces     []*ObjFace
}

type ObjVertex struct {
	X float32
	Y float32
	Z float32
}

type ObjTexCoord struct {
	S float32
	T float32
}

// numFaceElements is the number of required face elements. Only triangles are supported.
const numFaceElements = 3

// ObjFace is a face described by ObjFaceElements.
type ObjFace [numFaceElements]ObjFaceElement

// ObjFaceElement describes one point of a face.
type ObjFaceElement struct {
	// VertexIndex specifies a required vertex by global index starting from 1.
	VertexIndex int

	// TexCoordIndex specifies an optional texture coordinate by global index starting from 1.
	// It is 0 if no texture coordinate was specified.
	TexCoordIndex int
}

func ReadObjFile(r io.Reader) ([]*Obj, error) {
	var allObjs []*Obj
	var currentObj *Obj

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "o "):
			o, err := readObjObject(line)
			if err != nil {
				return nil, err
			}
			currentObj = o
			allObjs = append(allObjs, o)

		case strings.HasPrefix(line, "v "):
			v, err := readObjVertex(line)
			if err != nil {
				return nil, err
			}
			if currentObj == nil {
				return nil, errors.New("missing object ID")
			}
			currentObj.Vertices = append(currentObj.Vertices, v)

		case strings.HasPrefix(line, "vt "):
			tc, err := readObjTexCoord(line)
			if err != nil {
				return nil, err
			}
			if currentObj == nil {
				return nil, errors.New("missing object ID")
			}
			currentObj.TexCoords = append(currentObj.TexCoords, tc)

		case strings.HasPrefix(line, "f "):
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

func readObjTexCoord(line string) (*ObjTexCoord, error) {
	tc := &ObjTexCoord{}
	if _, err := fmt.Sscanf(line, "vt %f %f", &tc.S, &tc.T); err != nil {
		return nil, err
	}
	return tc, nil
}

func readObjFace(line string) (*ObjFace, error) {
	f := &ObjFace{}

	var specs [numFaceElements]string
	if _, err := fmt.Sscanf(line, "f %s %s %s", &specs[0], &specs[1], &specs[2]); err != nil {
		return nil, err
	}

	var err error
	makeElement := func(spec string) (ObjFaceElement, error) {
		tokens := strings.Split(spec, "/")
		if len(tokens) == 0 {
			return ObjFaceElement{}, errors.New("face has no elements")
		}

		e := ObjFaceElement{}

		e.VertexIndex, err = strconv.Atoi(tokens[0])
		if err != nil {
			return ObjFaceElement{}, err
		}

		if len(tokens) < 2 {
			return e, nil
		}

		e.TexCoordIndex, err = strconv.Atoi(tokens[1])
		if err != nil {
			return ObjFaceElement{}, err
		}

		return e, nil
	}

	for i, s := range specs {
		if f[i], err = makeElement(s); err != nil {
			return nil, err
		}
	}

	return f, nil
}
