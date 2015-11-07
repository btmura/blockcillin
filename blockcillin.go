package main

import (
	"log"
	"math"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	vertexShaderSource = `
		uniform mat4 u_projection_view_matrix;
		uniform mat4 u_matrix;

		attribute vec4 a_position;
		attribute vec2 a_tex_coord;

		varying vec2 v_tex_coord;

		void main(void) {
			gl_Position = u_projection_view_matrix * u_matrix * a_position;
			v_tex_coord = a_tex_coord;
		}
	`

	fragmentShaderSource = `
		uniform sampler2D u_texture;

		varying vec2 v_tex_coord;

		void main(void) {
			gl_FragColor = texture2D(u_texture, v_tex_coord);
		}
	`
)

func init() {
	// This is needed to arrange that main() runs on the main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	logFatalIfErr := func(tag string, err error) {
		if err != nil {
			log.Fatalf("%s: %v", tag, err)
		}
	}

	mf, err := os.Open("models.obj")
	logFatalIfErr("os.Open", err)
	defer mf.Close()

	objs, err := ReadObjFile(mf)
	logFatalIfErr("ReadObjFile", err)

	tf, err := os.Open("texture.png")
	logFatalIfErr("os.Open", err)
	defer tf.Close()

	logFatalIfErr("glfw.Init", glfw.Init())
	defer glfw.Terminate()

	win, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	logFatalIfErr("glfw.CreateWindow", err)

	win.MakeContextCurrent()

	logFatalIfErr("gl.Init", gl.Init())

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s", version)

	program, err := createProgram(vertexShaderSource, fragmentShaderSource)
	logFatalIfErr("createProgram", err)
	gl.UseProgram(program)

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

	var vertices []float32
	var texCoords []float32
	var indices []uint16

	elementIndexMap := map[ObjFaceElement]uint16{}
	var nextIndex uint16
	for _, f := range objs[0].Faces {
		for _, e := range *f {
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

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4 /* total bytes */, gl.Ptr(vertices), gl.STATIC_DRAW)

	positionAttrib, err := getAttribLocation(program, "a_position")
	logFatalIfErr("getAttribLocation", err)

	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	var tbo uint32
	gl.GenBuffers(1, &tbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, tbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(texCoords)*4 /*total bytes */, gl.Ptr(texCoords), gl.STATIC_DRAW)

	texCoordAttrib, err := getAttribLocation(program, "a_tex_coord")
	logFatalIfErr("getAttribLocation", err)
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*2 /* total bytes */, gl.Ptr(indices), gl.STATIC_DRAW)

	texture, err := createTexture(tf)
	logFatalIfErr("createTexture", err)

	m := NewScaleMatrix(0.5, 0.5, 0.5)
	m = m.Mult(NewYRotationMatrix(toRadians(30.0)))
	m = m.Mult(NewXRotationMatrix(toRadians(30.0)))
	m = m.Mult(NewTranslationMatrix(0.0, 0.0, 0.0))

	matrixUniform, err := getUniformLocation(program, "u_matrix")
	logFatalIfErr("getUniformLocation", err)
	gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])

	textureUniform, err := getUniformLocation(program, "u_texture")
	logFatalIfErr("getUniformLocation", err)
	gl.Uniform1i(textureUniform, 0)

	projectionViewMatrixUniform, err := getUniformLocation(program, "u_projection_view_matrix")
	logFatalIfErr("getUniformLocation", err)
	vm := makeViewMatrix()
	sizeCallback := func(w *glfw.Window, width, height int) {
		pvm := vm.Mult(makeProjectionMatrix(width, height))
		gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &pvm[0])
		gl.Viewport(0, 0, int32(width), int32(height))
	}
	win.SetSizeCallback(sizeCallback)

	// Call the size callback to set the initial projection view matrix and viewport.
	w, h := win.GetSize()
	sizeCallback(win, w, h)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0, 0, 0, 0)
	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_SHORT, gl.Ptr(nil))

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func makeProjectionMatrix(width, height int) *Matrix4 {
	aspect := float32(width) / float32(height)
	fovRadians := float32(math.Pi) / 2
	return NewPerspectiveMatrix(fovRadians, aspect, 1, 2000)
}

func makeViewMatrix() *Matrix4 {
	cameraPosition := &Vector3{0, 0, 3}
	targetPosition := &Vector3{}
	up := &Vector3{0, 1, 0}
	return NewViewMatrix(cameraPosition, targetPosition, up)
}

func toRadians(degrees float32) float32 {
	return degrees * float32(math.Pi) / 180
}
