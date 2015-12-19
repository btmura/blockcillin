package main

import "log"

// mesh is an OBJ file object.
type mesh struct {
	// id is the object's ID in the OBJ file.
	id string

	// vbo is a buffer object name to a buffer with vertices.
	vbo uint32

	// nbo is a buffer object name to a buffer with normals.
	nbo uint32

	// tbo is a buffer object name to a buffer with texture coordinates.
	tbo uint32

	// ibo is a buffer object name to a buffer with indices.
	ibo uint32

	// count is how many elements to render.
	count int32
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

	for i, m := range meshes {
		m.vbo = vbo
		m.nbo = nbo
		m.tbo = tbo
		m.ibo = iboNames[i]
	}

	return meshes
}
