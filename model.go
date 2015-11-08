package main

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/kylelemons/godebug/pretty"
)

// Model is a 3D model consisting of an
type Model struct {

	// VBO is the shared Vertex Buffer Object.
	VBO *ModelBufferObject

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
	log.Printf("objs:\n%+v", pretty.Sprint(objs))

	m := &Model{
		VBO:     new(ModelBufferObject),
		TBO:     new(ModelBufferObject),
		IBOByID: map[string]*ModelBufferObject{},
	}

	var vertices []float32
	var texCoords []float32

	elementIndexMap := map[ObjFaceElement]uint16{}
	var nextIndex uint16

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

		var indices []uint16
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
	log.Printf("texCoords: %d", len(texCoordTable))

	gl.GenBuffers(1, &m.VBO.Name)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO.Name)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4 /* total bytes */, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.GenBuffers(1, &m.TBO.Name)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.TBO.Name)
	gl.BufferData(gl.ARRAY_BUFFER, len(texCoords)*4 /*total bytes */, gl.Ptr(texCoords), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return m
}
