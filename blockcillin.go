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
		#version 330 core

		// TODO(btmura): use uniforms to make these configurable
		const vec3 ambient_light = vec3(0.5, 0.5, 0.5);
		const vec3 directional_light_color = vec3(0.5, 0.5, 0.5);
		const vec3 directional_vector = vec3(0.5, 0.5, 0.5);

		uniform mat4 u_projection_view_matrix;
		uniform mat4 u_normal_matrix;
		uniform mat4 u_matrix;

		layout (location = 0) in vec4 a_position;
		layout (location = 1) in vec4 a_normal;
		layout (location = 2) in vec2 a_tex_coord;

		out vec2 v_tex_coord;
		out vec3 v_lighting;

		void main(void) {
			gl_Position = u_projection_view_matrix * u_matrix * a_position;
			v_tex_coord = a_tex_coord;

			vec4 transformedNormal = u_normal_matrix * vec4(a_normal.xyz, 1.0);
			float directional = max(dot(transformedNormal.xyz, directional_vector), 0.0);
			v_lighting = ambient_light + (directional_light_color * directional);
		}
	`

	fragmentShaderSource = `
		#version 330 core

		uniform sampler2D u_texture;

		in vec2 v_tex_coord;
		in vec3 v_lighting;

		out vec4 frag_color;

		void main(void) {
			vec4 tex_color = texture2D(u_texture, v_tex_coord);
			frag_color = vec4(tex_color.rgb * v_lighting, tex_color.a);
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

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

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

	meshes := CreateMeshes(objs)

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

	normalMatrixUniform, err := GetUniformLocation(program, "u_normal_matrix")
	logFatalIfErr("getUniformLocation", err)

	matrixUniform, err := GetUniformLocation(program, "u_matrix")
	logFatalIfErr("getUniformLocation", err)

	textureUniform, err := GetUniformLocation(program, "u_texture")
	logFatalIfErr("getUniformLocation", err)

	vm := makeViewMatrix()
	sizeCallback := func(w *glfw.Window, width, height int) {
		pvm := vm.Mult(makeProjectionMatrix(width, height))
		gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &pvm[0])
		gl.Viewport(0, 0, int32(width), int32(height))
	}

	nm := vm.Inverse().Transpose()
	gl.UniformMatrix4fv(normalMatrixUniform, 1, false, &nm[0])

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
		localRotationY := float32(360.0 / len(meshes) * i)
		xq := NewAxisAngleQuaternion(xAxis, toRadians(globalRotation[0]))
		yq := NewAxisAngleQuaternion(yAxis, toRadians(globalRotation[1]+localRotationY))
		zq := NewAxisAngleQuaternion(zAxis, toRadians(globalRotation[2]))
		qm := NewQuaternionMatrix(yq.Mult(xq).Mult(zq).Normalize())

		m := NewScaleMatrix(0.25, 0.25, 0.25)
		m = m.Mult(NewTranslationMatrix(0, 0, -1))
		m = m.Mult(qm)
		gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])
	}

	gl.Uniform1i(textureUniform, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0, 0, 0, 0)
	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for i, m := range meshes {
			updateMatrix(i)

			gl.BindVertexArray(m.VAO)
			gl.DrawElements(gl.TRIANGLES, m.Count, gl.UNSIGNED_SHORT, gl.Ptr(nil))
			gl.BindVertexArray(0)
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
