package main

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Mesh is a an object from an OBJ file with its VAO and element count.
type Mesh struct {
	// ID is the object's ID from the OBJ file.
	ID string

	// VAO is the Vertex Array Object name to use with gl.BindVertexArray.
	VAO uint32

	// Count is the number of elements to use with gl.DrawElements.
	Count int32
}

func CreateMeshes(objs []*Obj) []*Mesh {
	var vertexTable []*ObjVertex
	var normalTable []*ObjNormal
	var texCoordTable []*ObjTexCoord

	var vertices []float32
	var normals []float32
	var texCoords []float32

	elementIndexMap := map[ObjFaceElement]uint16{}
	var nextIndex uint16

	var meshes []*Mesh
	var iboNames []uint32

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

		meshes = append(meshes, &Mesh{
			ID:    o.ID,
			Count: int32(len(indices)),
		})
		iboNames = append(iboNames, CreateElementArrayBuffer(indices))
	}

	log.Printf("vertices: %d", len(vertexTable))
	log.Printf("normals: %d", len(normalTable))
	log.Printf("texCoords: %d", len(texCoordTable))

	vbo := CreateArrayBuffer(vertices)
	nbo := CreateArrayBuffer(normals)
	tbo := CreateArrayBuffer(texCoords)

	for i, m := range meshes {
		gl.GenVertexArrays(1, &m.VAO)
		gl.BindVertexArray(m.VAO)

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, nbo)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, tbo)
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, iboNames[i])
		gl.BindVertexArray(0)
	}

	return meshes
}
