package main

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// mesh is a an object from an OBJ file with its VAO and element count.
type mesh struct {
	// id is the object's ID from the OBJ file.
	id string

	// vao is the Vertex Array Object name to use with gl.BindVertexArray.
	vao uint32

	// count is the number of elements to use with gl.DrawElements.
	count int32
}

func (m *mesh) drawElements() {
	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, m.count, gl.UNSIGNED_SHORT, gl.Ptr(nil))
	gl.BindVertexArray(0)
}

func createMeshes(objs []*obj) []*mesh {
	var vertexTable []*objVertex
	var normalTable []*objNormal
	var texCoordTable []*objTexCoord

	var vertices []float32
	var normals []float32
	var texCoords []float32

	elementIndexMap := map[objFaceElement]uint16{}
	var nextIndex uint16

	var meshes []*mesh
	var iboNames []uint32

	for _, o := range objs {
		for _, v := range o.vertices {
			vertexTable = append(vertexTable, v)
		}
		for _, n := range o.normals {
			normalTable = append(normalTable, n)
		}
		for _, tc := range o.texCoords {
			texCoordTable = append(texCoordTable, tc)
		}

		var indices []uint16
		for _, f := range o.faces {
			for _, e := range f {
				if _, exists := elementIndexMap[e]; !exists {
					elementIndexMap[e] = nextIndex
					nextIndex++

					v := vertexTable[e.vertexIndex-1]
					vertices = append(vertices, v.x, v.y, v.z)

					n := normalTable[e.normalIndex-1]
					normals = append(normals, n.x, n.y, n.z)

					// Flip the y-axis to convert from OBJ to OpenGL.
					// OpenGL considers the origin to be lower left.
					// OBJ considers the origin to be upper left.
					tc := texCoordTable[e.texCoordIndex-1]
					texCoords = append(texCoords, tc.s, 1.0-tc.t)
				}

				indices = append(indices, elementIndexMap[e])
			}
		}

		meshes = append(meshes, &mesh{
			id:    o.id,
			count: int32(len(indices)),
		})
		iboNames = append(iboNames, createElementArrayBuffer(indices))
	}

	log.Printf("vertices: %d", len(vertexTable))
	log.Printf("normals: %d", len(normalTable))
	log.Printf("texCoords: %d", len(texCoordTable))

	vbo := createArrayBuffer(vertices)
	nbo := createArrayBuffer(normals)
	tbo := createArrayBuffer(texCoords)

	const (
		positionLocation = iota
		normalLocation
		texCoordLocation
	)

	for i, m := range meshes {
		gl.GenVertexArrays(1, &m.vao)
		gl.BindVertexArray(m.vao)

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(positionLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, nbo)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(normalLocation, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, tbo)
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(texCoordLocation, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, iboNames[i])
		gl.BindVertexArray(0)
	}

	return meshes
}
