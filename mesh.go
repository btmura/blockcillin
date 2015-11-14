package main

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Mesh is a model with multiple IBOs sharing the same VBO and TBO.
type Mesh struct {
	VAOByID map[string]*MeshBufferObject

	// VBO is the shared Vertex Buffer Object.
	VBO *MeshBufferObject

	// NBO is the shared Normal Buffer Object.
	NBO *MeshBufferObject

	// TBO is the shared Texture Coord Buffer Object.
	TBO *MeshBufferObject

	// IBOByID is map from OBJ file ID to Index Buffer Object.
	IBOByID map[string]*MeshBufferObject
}

type MeshBufferObject struct {
	// Name is the OpenGL buffer name set by gl.GenBuffers.
	Name uint32

	// Count is the number of logical units in the buffer.
	Count int32
}

func CreateMesh(objs []*Obj) *Mesh {
	m := &Mesh{
		VBO:     &MeshBufferObject{},
		NBO:     &MeshBufferObject{},
		TBO:     &MeshBufferObject{},
		IBOByID: map[string]*MeshBufferObject{},
		VAOByID: map[string]*MeshBufferObject{},
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

		ibo := &MeshBufferObject{
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

	loadBuffer := func(mbo *MeshBufferObject, data []float32) {
		gl.GenBuffers(1, &mbo.Name)
		gl.BindBuffer(gl.ARRAY_BUFFER, mbo.Name)
		gl.BufferData(gl.ARRAY_BUFFER, len(data)*4 /* total bytes */, gl.Ptr(data), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	loadBuffer(m.VBO, vertices)
	loadBuffer(m.NBO, normals)
	loadBuffer(m.TBO, texCoords)

	for id, ibo := range m.IBOByID {
		var vaoName uint32
		gl.GenVertexArrays(1, &vaoName)
		gl.BindVertexArray(vaoName)

		gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO.Name)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, m.NBO.Name)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, m.TBO.Name)
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo.Name)
		gl.BindVertexArray(0)

		m.VAOByID[id] = &MeshBufferObject{
			Name:  vaoName,
			Count: ibo.Count,
		}
	}

	return m
}
