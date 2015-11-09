package main

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Model is a model with multiple IBOs sharing the same VBO and TBO.
type Model struct {
	// VBO is the shared Vertex Buffer Object.
	VBO *ModelBufferObject

	// NBO is the shared Normal Buffer Object.
	NBO *ModelBufferObject

	// TBO is the shared Texture Coord Buffer Object.
	TBO *ModelBufferObject

	// IBOByID is map from OBJ file ID to Index Buffer Object.
	IBOByID map[string]*ModelBufferObject
}

type ModelBufferObject struct {
	// Name is the OpenGL buffer name set by gl.GenBuffers.
	Name uint32

	// Count is the number of logical units in the buffer.
	Count int32
}

func CreateModel(objs []*Obj) *Model {
	m := &Model{
		VBO:     &ModelBufferObject{},
		NBO:     &ModelBufferObject{},
		TBO:     &ModelBufferObject{},
		IBOByID: map[string]*ModelBufferObject{},
	}

	var vertices []float32
	var normals []float32
	var texCoords []float32

	elementIndexMap := map[ObjFaceElement]uint16{}
	var nextIndex uint16

	// Collect the vertices and texture coords used by the objects.
	var vertexTable []*ObjVertex
	var normalTable []*ObjNormal
	var texCoordTable []*ObjTexCoord
	for _, o := range objs {
		for _, v := range o.Vertices {
			vertexTable = append(vertexTable, v)
		}
		for _, n := range o.Normals {
			normalTable = append(normalTable, n)
		}
		for _, tc := range o.TexCoords {
			texCoordTable = append(texCoordTable, tc)
		}

		var indices []uint16
		for _, f := range o.Faces {
			for _, e := range f {
				if _, exists := elementIndexMap[e]; !exists {
					elementIndexMap[e] = nextIndex
					nextIndex++

					v := vertexTable[e.VertexIndex-1]
					vertices = append(vertices, v.X, v.Y, v.Z)

					n := normalTable[e.NormalIndex-1]
					normals = append(normals, n.X, n.Y, n.Z)

					// Flip the y-axis to convert from OBJ to OpenGL.
					// OpenGL considers the origin to be lower left.
					// OBJ considers the origin to be upper left.
					tc := texCoordTable[e.TexCoordIndex-1]
					texCoords = append(texCoords, tc.S, 1.0-tc.T)
				}

				indices = append(indices, elementIndexMap[e])
			}
		}

		ibo := &ModelBufferObject{
			Count: int32(len(indices)),
		}
		m.IBOByID[o.ID] = ibo

		gl.GenBuffers(1, &ibo.Name)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo.Name)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*2 /* total bytes */, gl.Ptr(indices), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	}

	log.Printf("vertices: %d", len(vertexTable))
	log.Printf("normals: %d", len(normalTable))
	log.Printf("texCoords: %d", len(texCoordTable))

	loadBuffer := func(mbo *ModelBufferObject, data []float32) {
		gl.GenBuffers(1, &mbo.Name)
		gl.BindBuffer(gl.ARRAY_BUFFER, mbo.Name)
		gl.BufferData(gl.ARRAY_BUFFER, len(data)*4 /* total bytes */, gl.Ptr(data), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	loadBuffer(m.VBO, vertices)
	loadBuffer(m.NBO, normals)
	loadBuffer(m.TBO, texCoords)

	return m
}
