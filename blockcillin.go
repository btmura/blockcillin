package main

//go:generate go-bindata data

import (
	"bytes"
	"io"
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	vertexShaderSource = `
		#version 330 core

		uniform mat4 u_projectionViewMatrix;
		uniform mat4 u_normalMatrix;
		uniform mat4 u_matrix;

		uniform vec3 u_ambientLight;
		uniform vec3 u_directionalLight;
		uniform vec3 u_directionalVector;

		layout (location = 0) in vec4 i_position;
		layout (location = 1) in vec4 i_normal;
		layout (location = 2) in vec2 i_texCoord;

		out vec2 texCoord;
		out vec3 lighting;

		void main(void) {
			gl_Position = u_projectionViewMatrix * u_matrix * i_position;

			texCoord = i_texCoord;

			vec4 transformedNormal = u_normalMatrix * vec4(i_normal.xyz, 1.0);
			float directional = max(dot(transformedNormal.xyz, u_directionalVector), 0.0);
			lighting = u_ambientLight + (u_directionalLight * directional);
		}
	`

	fragmentShaderSource = `
		#version 330 core

		uniform sampler2D u_texture;

		in vec2 texCoord;
		in vec3 lighting;

		out vec4 fragColor;

		void main(void) {
			vec4 texColor = texture2D(u_texture, texCoord);
			fragColor = vec4(texColor.rgb * lighting, texColor.a);
		}
	`
)

const (
	positionLocation = iota
	normalLocation
	texCoordLocation
)

var (
	xAxis = vector3{1, 0, 0}
	yAxis = vector3{0, 1, 0}
	zAxis = vector3{0, 0, 1}
)

var (
	ambientLight      = [3]float32{0.5, 0.5, 0.5}
	directionalLight  = [3]float32{0.5, 0.5, 0.5}
	directionalVector = [3]float32{0.5, 0.5, 0.5}
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

	mr, err := loadAsset("data/models.obj")
	logFatalIfErr("loadAsset", err)

	objs, err := readObjFile(mr)
	logFatalIfErr("readObjFile", err)

	meshes := createMeshes(objs)

	ta, err := loadAsset("data/texture.png")
	logFatalIfErr("loadAsset", err)

	texture, err := createTexture(ta)
	logFatalIfErr("createTexture", err)

	program, err := createProgram(vertexShaderSource, fragmentShaderSource)
	logFatalIfErr("createProgram", err)
	gl.UseProgram(program)

	projectionViewMatrixUniform, err := getUniformLocation(program, "u_projectionViewMatrix")
	logFatalIfErr("getUniformLocation", err)

	normalMatrixUniform, err := getUniformLocation(program, "u_normalMatrix")
	logFatalIfErr("getUniformLocation", err)

	matrixUniform, err := getUniformLocation(program, "u_matrix")
	logFatalIfErr("getUniformLocation", err)

	ambientLightUniform, err := getUniformLocation(program, "u_ambientLight")
	logFatalIfErr("getUniformLocation", err)

	directionalLightUniform, err := getUniformLocation(program, "u_directionalLight")
	logFatalIfErr("getUniformLocation", err)

	directionalVectorUniform, err := getUniformLocation(program, "u_directionalVector")
	logFatalIfErr("getUniformLocation", err)

	textureUniform, err := getUniformLocation(program, "u_texture")
	logFatalIfErr("getUniformLocation", err)

	gl.Uniform3fv(ambientLightUniform, 1, &ambientLight[0])
	gl.Uniform3fv(directionalLightUniform, 1, &directionalLight[0])
	gl.Uniform3fv(directionalVectorUniform, 1, &directionalVector[0])

	vm := makeViewMatrix()
	sizeCallback := func(w *glfw.Window, width, height int) {
		pvm := vm.mult(makeProjectionMatrix(width, height))
		gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &pvm[0])
		gl.Viewport(0, 0, int32(width), int32(height))
	}

	nm := vm.inverse().transpose()
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

		case glfw.KeyEscape:
			win.SetShouldClose(true)
		}
	})

	updateMatrix := func(i int) {
		localRotationY := float32(360.0 / len(meshes) * i)
		xq := newAxisAngleQuaternion(xAxis, toRadians(globalRotation[0]))
		yq := newAxisAngleQuaternion(yAxis, toRadians(globalRotation[1]+localRotationY))
		zq := newAxisAngleQuaternion(zAxis, toRadians(globalRotation[2]))
		qm := newQuaternionMatrix(yq.mult(xq).mult(zq).normalize())

		m := newScaleMatrix(0.25, 0.25, 0.25)
		m = m.mult(newTranslationMatrix(0, 0, -1))
		m = m.mult(qm)
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
			m.drawElements()
		}

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func loadAsset(name string) (io.Reader, error) {
	a, err := Asset(name)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(a), nil
}

func makeProjectionMatrix(width, height int) matrix4 {
	aspect := float32(width) / float32(height)
	fovRadians := float32(math.Pi) / 2
	return newPerspectiveMatrix(fovRadians, aspect, 1, 2000)
}

func makeViewMatrix() matrix4 {
	cameraPosition := vector3{0, 2, 3}
	targetPosition := vector3{}
	up := vector3{0, 1, 0}
	return newViewMatrix(cameraPosition, targetPosition, up)
}
