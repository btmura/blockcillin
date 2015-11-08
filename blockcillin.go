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

var (
	xAxis = Vector3{1, 0, 0}
	yAxis = Vector3{0, 1, 0}
	zAxis = Vector3{0, 0, 1}
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

	logFatalIfErr("glfw.Init", glfw.Init())
	defer glfw.Terminate()

	win, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	logFatalIfErr("glfw.CreateWindow", err)

	win.MakeContextCurrent()

	logFatalIfErr("gl.Init", gl.Init())

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s", version)

	mf, err := os.Open("models.obj")
	logFatalIfErr("os.Open", err)
	defer mf.Close()

	objs, err := ReadObjFile(mf)
	logFatalIfErr("ReadObjFile", err)

	model := CreateModel(objs)

	tf, err := os.Open("texture.png")
	logFatalIfErr("os.Open", err)
	defer tf.Close()

	texture, err := CreateTexture(tf)
	logFatalIfErr("createTexture", err)

	program, err := CreateProgram(vertexShaderSource, fragmentShaderSource)
	logFatalIfErr("createProgram", err)
	gl.UseProgram(program)

	projectionViewMatrixUniform, err := GetUniformLocation(program, "u_projection_view_matrix")
	logFatalIfErr("getUniformLocation", err)

	matrixUniform, err := GetUniformLocation(program, "u_matrix")
	logFatalIfErr("getUniformLocation", err)

	positionAttrib, err := GetAttribLocation(program, "a_position")
	logFatalIfErr("getAttribLocation", err)

	texCoordAttrib, err := GetAttribLocation(program, "a_tex_coord")
	logFatalIfErr("getAttribLocation", err)

	textureUniform, err := GetUniformLocation(program, "u_texture")
	logFatalIfErr("getUniformLocation", err)

	vm := makeViewMatrix()
	sizeCallback := func(w *glfw.Window, width, height int) {
		pvm := vm.Mult(makeProjectionMatrix(width, height))
		gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &pvm[0])
		gl.Viewport(0, 0, int32(width), int32(height))
	}

	// Call the size callback to set the initial projection view matrix and viewport.
	w, h := win.GetSize()
	sizeCallback(win, w, h)
	win.SetSizeCallback(sizeCallback)

	globalRotation := []float32{0, 0, 0}
	win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press && action != glfw.Repeat {
			return
		}

		switch key {
		case glfw.KeyUp:
			globalRotation[0] += 5

		case glfw.KeyDown:
			globalRotation[0] -= 5

		case glfw.KeyLeft:
			globalRotation[1] -= 5

		case glfw.KeyRight:
			globalRotation[1] += 5

		case glfw.Key1:
			globalRotation[2] -= 5

		case glfw.Key0:
			globalRotation[2] += 5
		}
	})

	updateMatrix := func(i int) {
		localRotationY := float32(360.0 / len(model.IBOByID) * i)
		xq := NewAxisAngleQuaternion(xAxis, toRadians(globalRotation[0]))
		yq := NewAxisAngleQuaternion(yAxis, toRadians(globalRotation[1]+localRotationY))
		zq := NewAxisAngleQuaternion(zAxis, toRadians(globalRotation[2]))
		qm := NewQuaternionMatrix(yq.Mult(xq).Mult(zq).Normalize())

		m := NewScaleMatrix(0.25, 0.25, 0.25)
		m = m.Mult(NewTranslationMatrix(0, 0, -1))
		m = m.Mult(qm)
		gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, model.VBO.Name)
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, model.TBO.Name)
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.Uniform1i(textureUniform, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	var objIDs []string
	for id := range model.IBOByID {
		objIDs = append(objIDs, id)
	}

	gl.ClearColor(0, 0, 0, 0)
	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for i, id := range objIDs {
			updateMatrix(i)

			ibo := model.IBOByID[id]
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo.Name)
			gl.DrawElements(gl.TRIANGLES, ibo.Count, gl.UNSIGNED_SHORT, gl.Ptr(nil))
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		}

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func makeProjectionMatrix(width, height int) Matrix4 {
	aspect := float32(width) / float32(height)
	fovRadians := float32(math.Pi) / 2
	return NewPerspectiveMatrix(fovRadians, aspect, 1, 2000)
}

func makeViewMatrix() Matrix4 {
	cameraPosition := Vector3{0, 2, 3}
	targetPosition := Vector3{}
	up := Vector3{0, 1, 0}
	return NewViewMatrix(cameraPosition, targetPosition, up)
}
