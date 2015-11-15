package main

//go:generate go-bindata data

import (
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
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

	cameraPosition = vector3{0, 1, 3}

	blockColorByObjID = map[string]blockColor{
		"red_heart":       red,
		"purple_diamond":  purple,
		"blue_down_arrow": blue,
		"cyan_up_arrow":   cyan,
		"green_circle":    green,
		"yellow_star":     yellow,
	}
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

	mr, err := newAssetReader("data/models.obj")
	logFatalIfErr("newAssetReader", err)

	objs, err := readObjFile(mr)
	logFatalIfErr("readObjFile", err)

	meshes := createMeshes(objs)

	tr, err := newAssetReader("data/texture.png")
	logFatalIfErr("newAssetReader", err)

	texture, err := createTexture(tr)
	logFatalIfErr("createTexture", err)

	vs, err := getStringAsset("data/shader.vert")
	logFatalIfErr("getStringAsset", err)

	fs, err := getStringAsset("data/shader.frag")
	logFatalIfErr("getStringAsset", err)

	program, err := createProgram(vs, fs)
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

	gl.Uniform1i(textureUniform, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	meshByBlockColor := map[blockColor]*mesh{}
	for _, m := range meshes {
		meshByBlockColor[blockColorByObjID[m.id]] = m
	}

	b := newBoard()

	updateMatrix := func(row, col int) {
		localRotationY := float32(360.0 / b.cellCount * col)
		xq := newAxisAngleQuaternion(xAxis, toRadians(globalRotation[0]))
		yq := newAxisAngleQuaternion(yAxis, toRadians(globalRotation[1]+localRotationY))
		zq := newAxisAngleQuaternion(zAxis, toRadians(globalRotation[2]))
		qm := newQuaternionMatrix(yq.mult(xq).mult(zq).normalize())

		m := newScaleMatrix(0.25, 0.25, 0.25)
		m = m.mult(newTranslationMatrix(0, -0.5*float32(row), -1))
		m = m.mult(qm)
		gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])
	}

	gl.ClearColor(0, 0, 0, 0)
	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for row, r := range b.rings {
			for col, c := range r.cells {
				updateMatrix(row, col)
				meshByBlockColor[c.blockColor].drawElements()
			}
		}

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func makeProjectionMatrix(width, height int) matrix4 {
	aspect := float32(width) / float32(height)
	fovRadians := float32(math.Pi) / 2
	return newPerspectiveMatrix(fovRadians, aspect, 1, 2000)
}

func makeViewMatrix() matrix4 {
	targetPosition := vector3{}
	up := vector3{0, 1, 0}
	return newViewMatrix(cameraPosition, targetPosition, up)
}
