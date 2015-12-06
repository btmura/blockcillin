package main

import (
	"bufio"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	yAxis          = vector3{0, 1, 0}
	cameraPosition = vector3{0, 5, 25}
)

type renderer struct {
	menuMesh       *mesh
	selectorMesh   *mesh
	blockMeshes    map[blockColor]*mesh
	fragmentMeshes map[blockColor][4]*mesh

	boardTexture       uint32
	titleTextTexture   uint32
	newGameTextTexture uint32

	projectionViewMatrixUniform int32
	normalMatrixUniform         int32
	matrixUniform               int32
	textureUniform              int32
	grayscaleUniform            int32
	brightnessUniform           int32
	alphaUniform                int32

	sizeCallback func(width, height int)
}

func (rr *renderer) init() error {
	objs, err := readObjFile(newAssetReader("data/meshes.obj"))
	logFatalIfErr("readObjFile", err)

	meshes := createMeshes(objs)
	meshMap := map[string]*mesh{}
	for i, m := range meshes {
		log.Printf("mesh %2d: %s", i, m.id)
		meshMap[m.id] = m
	}
	mm := func(id string) *mesh {
		m, ok := meshMap[id]
		if !ok {
			log.Fatalf("mesh not found: %s", id)
		}
		return m
	}

	colorObjIDs := map[blockColor]string{
		red:    "red",
		purple: "purple",
		blue:   "blue",
		cyan:   "cyan",
		green:  "green",
		yellow: "yellow",
	}

	rr.menuMesh = mm("menu")
	rr.selectorMesh = mm("selector")
	rr.blockMeshes = map[blockColor]*mesh{}
	rr.fragmentMeshes = map[blockColor][4]*mesh{}

	for c, id := range colorObjIDs {
		rr.blockMeshes[c] = mm(id)
		rr.fragmentMeshes[c] = [4]*mesh{
			mm(id + "_north_west"),
			mm(id + "_north_east"),
			mm(id + "_south_east"),
			mm(id + "_south_west"),
		}
	}

	rr.boardTexture, err = createAssetTexture(gl.TEXTURE0, "data/texture.png")
	logFatalIfErr("createAssetTexture", err)

	font, err := freetype.ParseFont(MustAsset("data/Orbitron Medium.ttf"))
	logFatalIfErr("freetype.ParseFont", err)

	rr.titleTextTexture, err = createMenuTextTexture(gl.TEXTURE1, "b l o c k c i l l i n", font)
	logFatalIfErr("createMenuTextTexture", err)

	rr.newGameTextTexture, err = createMenuTextTexture(gl.TEXTURE2, "N E W   G A M E", font)
	logFatalIfErr("createMenuTextTexture", err)

	program, err := createProgram(assetString("data/shader.vert"), assetString("data/shader.frag"))
	logFatalIfErr("createProgram", err)
	gl.UseProgram(program)

	rr.projectionViewMatrixUniform, err = getUniformLocation(program, "u_projectionViewMatrix")
	logFatalIfErr("getUniformLocation", err)

	rr.normalMatrixUniform, err = getUniformLocation(program, "u_normalMatrix")
	logFatalIfErr("getUniformLocation", err)

	rr.matrixUniform, err = getUniformLocation(program, "u_matrix")
	logFatalIfErr("getUniformLocation", err)

	ambientLightUniform, err := getUniformLocation(program, "u_ambientLight")
	logFatalIfErr("getUniformLocation", err)

	directionalLightUniform, err := getUniformLocation(program, "u_directionalLight")
	logFatalIfErr("getUniformLocation", err)

	directionalVectorUniform, err := getUniformLocation(program, "u_directionalVector")
	logFatalIfErr("getUniformLocation", err)

	rr.textureUniform, err = getUniformLocation(program, "u_texture")
	logFatalIfErr("getUniformLocation", err)

	rr.grayscaleUniform, err = getUniformLocation(program, "u_grayscale")
	logFatalIfErr("getUniformLocation", err)

	rr.brightnessUniform, err = getUniformLocation(program, "u_brightness")
	logFatalIfErr("getUniformLocation", err)

	rr.alphaUniform, err = getUniformLocation(program, "u_alpha")
	logFatalIfErr("getUniformLocation", err)

	vm := makeViewMatrix()
	rr.sizeCallback = func(width, height int) {
		pvm := vm.mult(makeProjectionMatrix(width, height))
		gl.UniformMatrix4fv(rr.projectionViewMatrixUniform, 1, false, &pvm[0])
		gl.Viewport(0, 0, int32(width), int32(height))
	}

	nm := vm.inverse().transpose()
	gl.UniformMatrix4fv(rr.normalMatrixUniform, 1, false, &nm[0])

	ambientLight := [3]float32{0.5, 0.5, 0.5}
	directionalLight := [3]float32{0.5, 0.5, 0.5}
	directionalVector := [3]float32{0.5, 0.5, 0.5}

	gl.Uniform3fv(ambientLightUniform, 1, &ambientLight[0])
	gl.Uniform3fv(directionalLightUniform, 1, &directionalLight[0])
	gl.Uniform3fv(directionalVectorUniform, 1, &directionalVector[0])

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 0)

	return nil
}

func (rr *renderer) render(b *board, fudge float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	rr.renderBoard(b, fudge)
	rr.renderMenu()
}

func (rr *renderer) renderBoard(b *board, fudge float32) {
	const (
		nw = iota
		ne
		se
		sw
	)

	s := b.selector

	cellRotationY := float32(360.0 / b.cellCount)
	startRotationY := cellRotationY / 2
	cellTranslationY := float32(2.0)

	globalTranslationY := float32(0)
	globalTranslationZ := float32(4)

	selectorRelativeX := func(fudge float32) float32 {
		move := func(delta float32) float32 {
			return linear(s.step+fudge, float32(s.x), delta, numMoveSteps)
		}

		switch s.state {
		case selectorMovingLeft:
			return move(-1)

		case selectorMovingRight:
			return move(1)
		}

		return float32(s.x)
	}

	selectorRelativeY := func(fudge float32) float32 {
		move := func(delta float32) float32 {
			return linear(s.step+fudge, float32(s.y), delta, numMoveSteps)
		}

		switch s.state {
		case selectorMovingUp:
			return move(-1)
		case selectorMovingDown:
			return move(1)
		}
		return float32(s.y)
	}

	boardRelativeY := func(fudge float32) float32 {
		return linear(b.riseStep+fudge, float32(b.y), 1, numRiseSteps)
	}

	blockRelativeX := func(b *block, fudge float32) float32 {
		move := func(start, delta float32) float32 {
			return linear(b.step+fudge, start, delta, numSwapSteps)
		}

		switch b.state {
		case blockSwappingFromLeft:
			return move(-1, 1)

		case blockSwappingFromRight:
			return move(1, -1)
		}

		return 0
	}

	blockRelativeY := func(b *block, fudge float32) float32 {
		if b.state == blockDroppingFromAbove {
			return linear(b.step+fudge, 1, -1, numDropSteps)
		}
		return 0
	}

	blockMatrix := func(b *block, x, y int, fudge float32) matrix4 {
		ty := globalTranslationY + cellTranslationY*(-float32(y)+blockRelativeY(b, fudge))

		ry := startRotationY + cellRotationY*(-float32(x)-blockRelativeX(b, fudge)+selectorRelativeX(fudge))
		yq := newAxisAngleQuaternion(yAxis, toRadians(ry))
		qm := newQuaternionMatrix(yq.normalize())

		m := newTranslationMatrix(0, ty, globalTranslationZ)
		m = m.mult(qm)
		return m
	}

	renderSelector := func(fudge float32) {
		sc := pulse(s.pulse+fudge, 1.0, 0.025, 0.1)
		ty := globalTranslationY - cellTranslationY*selectorRelativeY(fudge)

		m := newScaleMatrix(sc, sc, sc)
		m = m.mult(newTranslationMatrix(0, ty, globalTranslationZ))
		gl.UniformMatrix4fv(rr.matrixUniform, 1, false, &m[0])

		rr.selectorMesh.drawElements()
	}

	renderCell := func(c *cell, x, y int, fudge float32) {
		sx := float32(1)
		bv := float32(0)

		switch c.block.state {
		case blockDroppingFromAbove:
			sx = linear(c.block.step+fudge, 1, -0.5, numDropSteps)
		case blockFlashing:
			bv = pulse(c.block.step+fudge, 0, 0.5, 1.5)
		}
		gl.Uniform1f(rr.brightnessUniform, bv)

		m := newScaleMatrix(sx, 1, 1)
		m = m.mult(blockMatrix(c.block, x, y, fudge))
		gl.UniformMatrix4fv(rr.matrixUniform, 1, false, &m[0])
		rr.blockMeshes[c.block.color].drawElements()
	}

	renderCellFragments := func(c *cell, x, y int, fudge float32) {
		render := func(sc, rx, ry, rz float32, dir int) {
			m := newScaleMatrix(sc, sc, sc)
			m = m.mult(newTranslationMatrix(rx, ry, rz))
			m = m.mult(blockMatrix(c.block, x, y, fudge))
			gl.UniformMatrix4fv(rr.matrixUniform, 1, false, &m[0])
			rr.fragmentMeshes[c.block.color][dir].drawElements()
		}

		ease := func(start, change float32) float32 {
			return easeOutCubic(c.block.step+fudge, start, change, numExplodeSteps)
		}

		var bv float32
		var av float32
		switch c.block.state {
		case blockCracking, blockCracked:
			av = 1
		case blockExploding:
			bv = ease(0, 1)
			av = ease(1, -1)
		}
		gl.Uniform1f(rr.brightnessUniform, bv)
		gl.Uniform1f(rr.alphaUniform, av)

		const (
			maxCrack  = 0.03
			maxExpand = 0.02
			maxJitter = 0.1
		)
		var rs float32
		var rt float32
		var j float32
		switch c.block.state {
		case blockCracking:
			rs = ease(1, 1+maxExpand)
			rt = ease(0, maxCrack)
			j = pulse(c.block.step+fudge, 0, 0.5, 1.5)
		case blockCracked:
			rs = 1
			rt = maxCrack
		case blockExploding:
			rs = ease(1, -1)
			rt = ease(maxCrack, math.Pi*0.75)
		}

		const szt = 0.5 // starting z translation since model is 0.5 in depth
		wx, ex := -rt, rt
		fz, bz := rt+szt, -rt-szt

		const amp = 1
		ny := rt + amp*float32(math.Sin(float64(rt)))
		sy := -rt + amp*(float32(math.Cos(float64(-rt)))-1)

		render(rs, wx+j, ny+j, fz, nw) // front north west
		render(rs, ex+j, ny+j, fz, ne) // front north east

		render(rs, wx+j, ny+j, bz, nw) // back north west
		render(rs, ex+j, ny+j, bz, ne) // back north east

		render(rs, wx+j, sy+j, fz, sw) // front south west
		render(rs, ex+j, sy+j, fz, se) // front south east

		render(rs, wx+j, sy+j, bz, sw) // back south west
		render(rs, ex+j, sy+j, bz, se) // back south east
	}

	globalTranslationY = cellTranslationY * (4 + boardRelativeY(fudge))

	gl.Uniform1i(rr.textureUniform, int32(rr.boardTexture)-1)

	for i := 0; i <= 2; i++ {
		gl.Uniform1f(rr.grayscaleUniform, 0)
		gl.Uniform1f(rr.brightnessUniform, 0)
		gl.Uniform1f(rr.alphaUniform, 1)

		switch i {
		case 0:
			gl.Disable(gl.BLEND)
			renderSelector(fudge)

		case 1:
			gl.Enable(gl.BLEND)
		}

		for y, r := range b.rings {
			for x, c := range r.cells {
				switch i {
				case 0: // draw opaque objects
					switch c.block.state {
					case blockStatic,
						blockSwappingFromLeft,
						blockSwappingFromRight,
						blockDroppingFromAbove,
						blockFlashing:
						renderCell(c, x, y, fudge)

					case blockCracking, blockCracked:
						renderCellFragments(c, x, y, fudge)
					}

				case 1: // draw transparent objects
					switch c.block.state {
					case blockExploding:
						renderCellFragments(c, x, y, fudge)
					}
				}
			}
		}

		for y, r := range b.spareRings {
			switch {
			case i == 0 && y == 0: // draw opaque objects
				gl.Uniform1f(rr.grayscaleUniform, easeInExpo(b.riseStep+fudge, 1, -1, numRiseSteps))
				gl.Uniform1f(rr.brightnessUniform, 0)
				gl.Uniform1f(rr.alphaUniform, 1)
				for x, c := range r.cells {
					renderCell(c, x, y+b.ringCount, fudge)
				}

			case i == 1 && y == 1: // draw transparent objects
				gl.Uniform1f(rr.grayscaleUniform, 1)
				gl.Uniform1f(rr.brightnessUniform, 0)
				gl.Uniform1f(rr.alphaUniform, easeInExpo(b.riseStep+fudge, 0, 1, numRiseSteps))
				for x, c := range r.cells {
					renderCell(c, x, y+b.ringCount, fudge)
				}
			}
		}
	}
}

func (rr *renderer) renderMenu() {
	gl.Enable(gl.BLEND)

	m := newScaleMatrix(5, 5, 5)
	m = m.mult(newTranslationMatrix(0, 0, 0))
	gl.UniformMatrix4fv(rr.matrixUniform, 1, false, &m[0])
	gl.Uniform1i(rr.textureUniform, int32(rr.titleTextTexture)-1)
	rr.menuMesh.drawElements()

	m = newScaleMatrix(5, 5, 5)
	m = m.mult(newTranslationMatrix(0, -2, 0))
	gl.UniformMatrix4fv(rr.matrixUniform, 1, false, &m[0])
	gl.Uniform1i(rr.textureUniform, int32(rr.newGameTextTexture)-1)
	rr.menuMesh.drawElements()
}

func makeProjectionMatrix(width, height int) matrix4 {
	aspect := float32(width) / float32(height)
	fovRadians := float32(math.Pi) / 3
	return newPerspectiveMatrix(fovRadians, aspect, 1, 2000)
}

func makeViewMatrix() matrix4 {
	targetPosition := vector3{}
	up := yAxis
	return newViewMatrix(cameraPosition, targetPosition, up)
}

func createAssetTexture(textureUnit uint32, name string) (uint32, error) {
	img, _, err := image.Decode(newAssetReader(name))
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return createTexture(textureUnit, rgba)
}

func createMenuTextTexture(textureUnit uint32, text string, f *truetype.Font) (uint32, error) {
	rgba, err := createMenuTextImage(f, text)
	if err != nil {
		return 0, err
	}

	texture, err := createTexture(textureUnit, rgba)
	if err != nil {
		return 0, err
	}

	return texture, nil
}

func createMenuTextImage(f *truetype.Font, text string) (*image.RGBA, error) {
	fg, bg := image.White, image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, 128, 32))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	const fontSize = 12.0

	c := freetype.NewContext()
	c.SetFont(f)
	c.SetDPI(72)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(fontSize)>>6))
	if _, err := c.DrawString(text, pt); err != nil {
		return nil, err
	}

	return rgba, nil
}

func writeDebugPNG(rgba *image.RGBA) {
	outFile, err := ioutil.TempFile("", "debug")
	logFatalIfErr("ioutil.TempFile", err)
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	logFatalIfErr("png.Encode", png.Encode(b, rgba))
	logFatalIfErr("bufio.Flush", b.Flush())
	log.Printf("wrote %s", outFile.Name())
}
