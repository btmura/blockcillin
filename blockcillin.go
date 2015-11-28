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

const secPerUpdate = 1.0 / 60.0

var (
	ambientLight      = [3]float32{0.5, 0.5, 0.5}
	directionalLight  = [3]float32{0.5, 0.5, 0.5}
	directionalVector = [3]float32{0.5, 0.5, 0.5}

	cameraPosition = vector3{0, 5, 25}

	blockColorByObjID = map[string]blockColor{
		"red":    red,
		"purple": purple,
		"blue":   blue,
		"cyan":   cyan,
		"green":  green,
		"yellow": yellow,
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

	mr, err := newAssetReader("data/meshes.obj")
	logFatalIfErr("newAssetReader", err)

	objs, err := readObjFile(mr)
	logFatalIfErr("readObjFile", err)

	meshes := createMeshes(objs)
	meshByBlockColor := map[blockColor]*mesh{}
	var selectorMesh *mesh
	for i, m := range meshes {
		log.Printf("mesh %d: %s", i, m.id)
		switch m.id {
		case "selector":
			selectorMesh = m
		default:
			if c, ok := blockColorByObjID[m.id]; ok {
				meshByBlockColor[c] = m
			}
		}
	}

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

	grayscaleUniform, err := getUniformLocation(program, "u_grayscale")
	logFatalIfErr("getUniformLocation", err)

	brightnessUniform, err := getUniformLocation(program, "u_brightness")
	logFatalIfErr("getUniformLocation", err)

	alphaUniform, err := getUniformLocation(program, "u_alpha")
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

	gl.Uniform1i(textureUniform, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	b := newBoard(&boardConfig{
		ringCount:       10,
		cellCount:       15,
		filledRingCount: 2,
		spareRingCount:  2,
	})
	s := b.selector

	win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press && action != glfw.Repeat {
			return
		}

		switch key {
		case glfw.KeyLeft:
			s.moveLeft()

		case glfw.KeyRight:
			s.moveRight()

		case glfw.KeyDown:
			s.moveDown()

		case glfw.KeyUp:
			s.moveUp()

		case glfw.KeySpace:
			b.swap()

		case glfw.KeyEscape:
			win.SetShouldClose(true)
		}
	})

	cellRotationY := float32(360.0 / b.cellCount)
	startRotationY := cellRotationY / 2
	cellTranslationY := float32(2.0)

	globalTranslationY := float32(0)
	globalTranslationZ := float32(4)

	renderSelector := func(fudge float32) {
		sc := s.scale(fudge)
		ty := globalTranslationY - cellTranslationY*s.relativeY(fudge)
		tz := globalTranslationZ

		m := newScaleMatrix(sc, sc, sc)
		m = m.mult(newTranslationMatrix(0, ty, tz))
		gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])

		selectorMesh.drawElements()
	}

	renderCell := func(c *cell, x, y int, fudge float32) {
		ry := startRotationY + cellRotationY*(-float32(x)-c.block.relativeX(fudge)+s.relativeX(fudge))
		ty := globalTranslationY + cellTranslationY*(-float32(y)+c.block.relativeY(fudge))
		tz := globalTranslationZ

		yq := newAxisAngleQuaternion(yAxis, toRadians(ry))
		qm := newQuaternionMatrix(yq.normalize())

		m := newTranslationMatrix(0, ty, tz)
		m = m.mult(qm)
		gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])

		meshByBlockColor[c.block.color].drawElements()
	}

	var lag float64
	prevTime := glfw.GetTime()

	gl.ClearColor(0, 0, 0, 0)
	for !win.ShouldClose() {
		currTime := glfw.GetTime()
		elapsed := currTime - prevTime
		prevTime = currTime
		lag += elapsed

		for lag >= secPerUpdate {
			s.update()
			b.update()
			lag -= secPerUpdate
		}
		fudge := float32(lag / secPerUpdate)

		globalTranslationY = cellTranslationY * (4 + b.relativeY(fudge))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Disable(gl.BLEND)
		gl.Uniform1f(grayscaleUniform, 0)
		gl.Uniform1f(brightnessUniform, 0)
		gl.Uniform1f(alphaUniform, 1)

		renderSelector(fudge)

		for i := 0; i <= 2; i++ {
			if i == 1 {
				gl.Enable(gl.BLEND)
			}

			gl.Uniform1f(grayscaleUniform, 0)
			for y, r := range b.rings {
				for x, c := range r.cells {
					alpha := c.block.alpha(fudge)

					switch {
					// First iteration: draw only opaque objects.
					case i == 0 && alpha >= 1.0:
						fallthrough

					// Second iteration: draw transparent objects.
					case i == 1 && alpha > 0 && alpha < 1:
						gl.Uniform1f(brightnessUniform, c.block.brightness(fudge))
						gl.Uniform1f(alphaUniform, alpha)
						renderCell(c, x, y, fudge)
					}
				}
			}

			gl.Uniform1f(brightnessUniform, 0)
			for y, r := range b.spareRings {
				alpha := b.spareRingAlpha(y, fudge)
				switch {
				// First iteration: draw only opaque objects.
				case i == 0 && alpha >= 1.0:
					fallthrough

				// Second iteration: draw transparent objects.
				case i == 1 && alpha > 0 && alpha < 1:
					gl.Uniform1f(grayscaleUniform, b.spareRingGrayscale(y, fudge))
					gl.Uniform1f(alphaUniform, alpha)
					for x, c := range r.cells {
						renderCell(c, x, y+b.ringCount, fudge)
					}
				}
			}

		}

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func makeProjectionMatrix(width, height int) matrix4 {
	aspect := float32(width) / float32(height)
	fovRadians := float32(math.Pi) / 3
	return newPerspectiveMatrix(fovRadians, aspect, 1, 2000)
}

func makeViewMatrix() matrix4 {
	targetPosition := vector3{}
	up := vector3{0, 1, 0}
	return newViewMatrix(cameraPosition, targetPosition, up)
}
