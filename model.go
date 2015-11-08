package main

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Model is a 3D model consisting of an ID and OpenGL handles.
type Model struct {
	// ID is the unique ID of the model.
	ID string

	// VBO is the Vertex Buffer Object handle.
	VBO uint32

	// TBO is the Texture Coord Buffer Object handle.
	TBO uint32

	// IBO is the Index Buffer Object handle.
	IBO uint32

	// NumIndices is the number of indices in the IBO.
	NumIndices int32
}

func CreateModels(objs []*Obj) []*Model {
	var models []*Model

	// Collect the vertices and texture coords used by the objects.
	var vertexTable []*ObjVertex
	var texCoordTable []*ObjTexCoord
	for _, o := range objs {
		for _, v := range o.Vertices {
			vertexTable = append(vertexTable, v)
		}
		for _, tc := range o.TexCoords {
			texCoordTable = append(texCoordTable, tc)
		}
	}

	// Parse each object's faces and create a corresponding model.
	for i, o := range objs {
		var vertices []float32
		var texCoords []float32
		var indices []uint16

		elementIndexMap := map[ObjFaceElement]uint16{}
		var nextIndex uint16

		for _, f := range o.Faces {
			for _, e := range f {
				if _, exists := elementIndexMap[e]; !exists {
					elementIndexMap[e] = nextIndex
					nextIndex++

					v := vertexTable[e.VertexIndex-1]
					vertices = append(vertices, v.X, v.Y, v.Z)

					// Flip the y-axis to convert from OBJ to OpenGL.
					// OpenGL considers the origin to be lower left.
					// OBJ considers the origin to be upper left.
					tc := texCoordTable[e.TexCoordIndex-1]
					texCoords = append(texCoords, tc.S, 1.0-tc.T)
				}

				indices = append(indices, elementIndexMap[e])
			}
		}

		m := &Model{
			ID:         o.ID,
			NumIndices: int32(len(indices)),
		}

		gl.GenBuffers(1, &m.VBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4 /* total bytes */, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.GenBuffers(1, &m.TBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, m.TBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(texCoords)*4 /*total bytes */, gl.Ptr(texCoords), gl.STATIC_DRAW)

		gl.GenBuffers(1, &m.IBO)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.IBO)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*2 /* total bytes */, gl.Ptr(indices), gl.STATIC_DRAW)

		models = append(models, m)

		log.Printf("model %d: %+v", i, m)
	}

	return models
}
