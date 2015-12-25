package renderer

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"

	"github.com/btmura/blockcillin/internal/asset"
	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	yAxis          = vector3{0, 1, 0}
	cameraPosition = vector3{0, 5, 25}
	targetPosition = vector3{}
	up             = yAxis

	ambientLightColor     = [3]float32{0.5, 0.5, 0.5}
	directionalLightColor = [3]float32{0.5, 0.5, 0.5}
	directionalVector     = [3]float32{0.5, 0.5, 0.5}
	blackColor            = [3]float32{}

	titleFontSize    = 54
	menuItemFontSize = 36

	titleTextColor    = color.White
	menuItemTextColor = color.Gray{100}
)

var (
	program                      uint32
	projectionViewMatrixUniform  int32
	modelMatrixUniform           int32
	normalMatrixUniform          int32
	ambientLightColorUniform     int32
	directionalLightColorUniform int32
	directionalVectorUniform     int32
	textureUniform               int32
	grayscaleUniform             int32
	brightnessUniform            int32
	alphaUniform                 int32
	mixColorUniform              int32
	mixAmountUniform             int32
)

var (
	// SizeCallback is the callback that GLFW should call when resizing the window.
	SizeCallback func(width, height int)

	// winWidth is the current window's width reported by the SizeCallback.
	winWidth int

	// winHeight is the current window's height reported by the SizeCallback.
	winHeight int

	// perspectiveProjectionViewMatrix is the perspective projection view matrix uniform value.
	perspectiveProjectionViewMatrix matrix4

	// orthoProjectionViewMatrix is the ortho projection view matrix uniform value.
	orthoProjectionViewMatrix matrix4
)

var (
	selectorMesh   *mesh
	blockMeshes    map[game.BlockColor]*mesh
	fragmentMeshes map[game.BlockColor][4]*mesh
	textLineMesh   *mesh
)

var (
	boardTexture uint32
	titleText    *rendererText
	menuItemText map[game.MenuItem]*rendererText
)

type rendererText struct {
	texture uint32
	width   float32
	height  float32
}

func Init() {
	logFatalIfErr("gl.Init", gl.Init())
	log.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	program, err := createProgram(asset.MustString("data/shader.vert"), asset.MustString("data/shader.frag"))
	logFatalIfErr("createProgram", err)
	gl.UseProgram(program)

	mustUniform := func(name string) int32 {
		l, err := getUniformLocation(program, name)
		logFatalIfErr("getUniformLocation", err)
		return l
	}

	projectionViewMatrixUniform = mustUniform("u_projectionViewMatrix")
	modelMatrixUniform = mustUniform("u_modelMatrix")
	normalMatrixUniform = mustUniform("u_normalMatrix")
	ambientLightColorUniform = mustUniform("u_ambientLightColor")
	directionalLightColorUniform = mustUniform("u_directionalLightColor")
	directionalVectorUniform = mustUniform("u_directionalVector")
	textureUniform = mustUniform("u_texture")
	grayscaleUniform = mustUniform("u_grayscale")
	brightnessUniform = mustUniform("u_brightness")
	alphaUniform = mustUniform("u_alpha")
	mixColorUniform = mustUniform("u_mixColor")
	mixAmountUniform = mustUniform("u_mixAmount")

	vm := newViewMatrix(cameraPosition, targetPosition, up)
	nm := vm.inverse().transpose()
	gl.UniformMatrix4fv(normalMatrixUniform, 1, false, &nm[0])

	gl.Uniform3fv(ambientLightColorUniform, 1, &ambientLightColor[0])
	gl.Uniform3fv(directionalLightColorUniform, 1, &directionalLightColor[0])
	gl.Uniform3fv(directionalVectorUniform, 1, &directionalVector[0])

	SizeCallback = func(width, height int) {
		if winWidth == width && winHeight == height {
			return
		}

		log.Printf("window size changed (%dx%d -> %dx%d)", int(winWidth), int(winHeight), width, height)
		gl.Viewport(0, 0, int32(width), int32(height))

		// Calculate new perspective projection view matrix.
		winWidth, winHeight = width, height
		fw, fh := float32(width), float32(height)
		aspect := fw / fh
		fovRadians := float32(math.Pi) / 3
		perspectiveProjectionViewMatrix = vm.mult(newPerspectiveMatrix(fovRadians, aspect, 1, 2000))

		// Calculate new ortho projection view matrix.
		orthoProjectionViewMatrix = newOrthoMatrix(fw, fh, fw /* use width as depth */)
	}

	objs, err := decodeObjs(asset.MustReader("data/meshes.obj"))
	logFatalIfErr("decodeObjs", err)

	meshes := createMeshes(objs)
	meshMap := map[string]*mesh{}
	for i, m := range meshes {
		log.Printf("mesh %d: %s", i, m.id)
		meshMap[m.id] = m
	}
	mm := func(id string) *mesh {
		m, ok := meshMap[id]
		if !ok {
			log.Fatalf("mesh not found: %s", id)
		}
		return m
	}

	colorObjIDs := map[game.BlockColor]string{
		game.Red:    "red",
		game.Purple: "purple",
		game.Blue:   "blue",
		game.Cyan:   "cyan",
		game.Green:  "green",
		game.Yellow: "yellow",
	}

	selectorMesh = mm("selector")
	blockMeshes = map[game.BlockColor]*mesh{}
	fragmentMeshes = map[game.BlockColor][4]*mesh{}
	for c, id := range colorObjIDs {
		blockMeshes[c] = mm(id)
		fragmentMeshes[c] = [4]*mesh{
			mm(id + "_north_west"),
			mm(id + "_north_east"),
			mm(id + "_south_east"),
			mm(id + "_south_west"),
		}
	}
	textLineMesh = mm("text_line")

	boardTexture, err = createAssetTexture(gl.TEXTURE0, "data/texture.png")
	logFatalIfErr("createAssetTexture", err)

	font, err := freetype.ParseFont(asset.MustAsset("data/Orbitron Medium.ttf"))
	logFatalIfErr("freetype.ParseFont", err)

	titleText, err = createText(gl.TEXTURE1, font, "b l o c k c i l l i n", titleFontSize, titleTextColor)
	logFatalIfErr("createText", err)

	menuItemText = map[game.MenuItem]*rendererText{}
	var textureUnit uint32 = gl.TEXTURE2
	for item, text := range game.MenuItemText {
		menuItemText[item], err = createText(textureUnit, font, text, menuItemFontSize, menuItemTextColor)
		logFatalIfErr("createText", err)
		textureUnit++
	}

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 0)
}

func createAssetTexture(textureUnit uint32, name string) (uint32, error) {
	img, _, err := image.Decode(asset.MustReader(name))
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return createTexture(textureUnit, rgba)
}

func createText(textureUnit uint32, f *truetype.Font, text string, fontSize int, color color.Color) (*rendererText, error) {
	rgba, w, h, err := createTextImage(f, text, fontSize, color)
	if err != nil {
		return nil, err
	}

	t, err := createTexture(textureUnit, rgba)
	if err != nil {
		return nil, err
	}
	return &rendererText{t, w, h}, nil
}

func createTextImage(f *truetype.Font, text string, fontSize int, color color.Color) (*image.RGBA, float32, float32, error) {
	// 1 pt = 1/72 in, 72 dpi = 1 in
	const dpi = 72

	fg, bg := image.NewUniform(color), image.Transparent

	c := freetype.NewContext()
	c.SetFont(f)
	c.SetDPI(dpi)
	c.SetFontSize(float64(fontSize))
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)

	// 1. Draw within small bounds to figure out bounds.
	// 2. Draw within final bounds.

	var rgba *image.RGBA
	w, h := 10, fontSize
	for i := 0; i < 2; i++ {
		rgba = image.NewRGBA(image.Rect(0, 0, w, h))
		draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)

		pt := freetype.Pt(0, int(c.PointToFixed(float64(fontSize))>>6))
		end, err := c.DrawString(text, pt)
		if err != nil {
			return nil, 0, 0, err
		}

		w = int(end.X >> 6)
	}

	return rgba, float32(w), float32(h), nil
}

func Render(g *game.Game, fudge float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	renderBoard(g, fudge)
	renderMenu(g, fudge)
}

func renderBoard(g *game.Game, fudge float32) {
	if g.Board == nil {
		return
	}

	var grayscale float32
	var darkness float32
	ease := func(start, change float32) float32 {
		return easeOutCubic2(g.StateProgress(fudge), start, change)
	}

	switch g.State {
	case game.GameInitial:
		grayscale = 1
		darkness = 0.8

	case game.GamePlaying:
		grayscale = ease(1, -1)
		darkness = ease(0.8, -0.8)

	case game.GamePaused:
		grayscale = ease(0, 1)
		darkness = ease(0, 0.8)

	case game.GameExiting:
		grayscale = 1
		darkness = ease(0.8, 1)
	}

	gl.Uniform3fv(mixColorUniform, 1, &blackColor[0])
	gl.Uniform1f(mixAmountUniform, darkness)

	gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &perspectiveProjectionViewMatrix[0])

	const (
		nw = iota
		ne
		se
		sw
	)

	b := g.Board
	s := b.Selector

	cellRotationY := float32(360.0 / b.CellCount)
	startRotationY := cellRotationY / 2
	cellTranslationY := float32(2.0)

	globalTranslationY := float32(0)
	globalTranslationZ := float32(4)

	selectorRelativeX := func(fudge float32) float32 {
		move := func(delta float32) float32 {
			return linear(s.Step+fudge, float32(s.X), delta, game.NumMoveSteps)
		}

		switch s.State {
		case game.SelectorMovingLeft:
			return move(-1)

		case game.SelectorMovingRight:
			return move(1)
		}

		return float32(s.X)
	}

	selectorRelativeY := func(fudge float32) float32 {
		move := func(delta float32) float32 {
			return linear(s.Step+fudge, float32(s.Y), delta, game.NumMoveSteps)
		}

		switch s.State {
		case game.SelectorMovingUp:
			return move(-1)
		case game.SelectorMovingDown:
			return move(1)
		}
		return float32(s.Y)
	}

	boardRelativeY := func(fudge float32) float32 {
		return linear(b.RiseStep+fudge, float32(b.Y), 1, game.NumRiseSteps)
	}

	blockRelativeX := func(b *game.Block, fudge float32) float32 {
		move := func(start, delta float32) float32 {
			return linear(b.Step+fudge, start, delta, game.NumSwapSteps)
		}

		switch b.State {
		case game.BlockSwappingFromLeft:
			return move(-1, 1)

		case game.BlockSwappingFromRight:
			return move(1, -1)
		}

		return 0
	}

	blockRelativeY := func(b *game.Block, fudge float32) float32 {
		if b.State == game.BlockDroppingFromAbove {
			return linear(b.Step+fudge, 1, -1, game.NumDropSteps)
		}
		return 0
	}

	blockMatrix := func(b *game.Block, x, y int, fudge float32) matrix4 {
		ty := globalTranslationY + cellTranslationY*(-float32(y)+blockRelativeY(b, fudge))

		ry := startRotationY + cellRotationY*(-float32(x)-blockRelativeX(b, fudge)+selectorRelativeX(fudge))
		yq := newAxisAngleQuaternion(yAxis, toRadians(ry))
		qm := newQuaternionMatrix(yq.normalize())

		m := newTranslationMatrix(0, ty, globalTranslationZ)
		m = m.mult(qm)
		return m
	}

	renderSelector := func(fudge float32) {
		sc := pulse(s.Pulse+fudge, 1.0, 0.025, 0.1)
		ty := globalTranslationY - cellTranslationY*selectorRelativeY(fudge)

		m := newScaleMatrix(sc, sc, sc)
		m = m.mult(newTranslationMatrix(0, ty, globalTranslationZ))
		gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])

		selectorMesh.drawElements()
	}

	renderCell := func(c *game.Cell, x, y int, fudge float32) {
		sx := float32(1)
		bv := float32(0)

		switch c.Block.State {
		case game.BlockDroppingFromAbove:
			sx = linear(c.Block.Step+fudge, 1, -0.5, game.NumDropSteps)
		case game.BlockFlashing:
			bv = pulse(c.Block.Step+fudge, 0, 0.5, 1.5)
		}
		gl.Uniform1f(brightnessUniform, bv)

		m := newScaleMatrix(sx, 1, 1)
		m = m.mult(blockMatrix(c.Block, x, y, fudge))
		gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
		blockMeshes[c.Block.Color].drawElements()
	}

	renderCellFragments := func(c *game.Cell, x, y int, fudge float32) {
		render := func(sc, rx, ry, rz float32, dir int) {
			m := newScaleMatrix(sc, sc, sc)
			m = m.mult(newTranslationMatrix(rx, ry, rz))
			m = m.mult(blockMatrix(c.Block, x, y, fudge))
			gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
			fragmentMeshes[c.Block.Color][dir].drawElements()
		}

		ease := func(start, change float32) float32 {
			return easeOutCubic(c.Block.Step+fudge, start, change, game.NumExplodeSteps)
		}

		var bv float32
		var av float32
		switch c.Block.State {
		case game.BlockCracking, game.BlockCracked:
			av = 1
		case game.BlockExploding:
			bv = ease(0, 1)
			av = ease(1, -1)
		}
		gl.Uniform1f(brightnessUniform, bv)
		gl.Uniform1f(alphaUniform, av)

		const (
			maxCrack  = 0.03
			maxExpand = 0.02
		)
		var rs float32
		var rt float32
		var j float32
		switch c.Block.State {
		case game.BlockCracking:
			rs = ease(1, 1+maxExpand)
			rt = ease(0, maxCrack)
			j = pulse(c.Block.Step+fudge, 0, 0.5, 1.5)
		case game.BlockCracked:
			rs = 1
			rt = maxCrack
		case game.BlockExploding:
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

	gl.Uniform1i(textureUniform, int32(boardTexture)-1)

	for i := 0; i <= 2; i++ {
		gl.Uniform1f(grayscaleUniform, grayscale)
		gl.Uniform1f(brightnessUniform, 0)
		gl.Uniform1f(alphaUniform, 1)

		if i == 0 {
			renderSelector(fudge)
		}

		for y, r := range b.Rings {
			for x, c := range r.Cells {
				switch i {
				case 0: // draw opaque objects
					switch c.Block.State {
					case game.BlockStatic,
						game.BlockSwappingFromLeft,
						game.BlockSwappingFromRight,
						game.BlockDroppingFromAbove,
						game.BlockFlashing:
						renderCell(c, x, y, fudge)

					case game.BlockCracking, game.BlockCracked:
						renderCellFragments(c, x, y, fudge)
					}

				case 1: // draw transparent objects
					switch c.Block.State {
					case game.BlockExploding:
						renderCellFragments(c, x, y, fudge)
					}
				}
			}
		}

		for y, r := range b.SpareRings {
			switch {
			case i == 0 && y == 0: // draw opaque objects
				finalGrayscale := easeInExpo(b.RiseStep+fudge, 1, -1, game.NumRiseSteps)
				if grayscale > finalGrayscale {
					finalGrayscale = grayscale
				}

				gl.Uniform1f(grayscaleUniform, finalGrayscale)
				gl.Uniform1f(brightnessUniform, 0)
				gl.Uniform1f(alphaUniform, 1)
				for x, c := range r.Cells {
					renderCell(c, x, y+b.RingCount, fudge)
				}

			case i == 1 && y == 1: // draw transparent objects
				gl.Uniform1f(grayscaleUniform, 1)
				gl.Uniform1f(brightnessUniform, 0)
				gl.Uniform1f(alphaUniform, easeInExpo(b.RiseStep+fudge, 0, 1, game.NumRiseSteps))
				for x, c := range r.Cells {
					renderCell(c, x, y+b.RingCount, fudge)
				}
			}
		}
	}
}

func renderMenu(g *game.Game, fudge float32) {
	alpha := float32(1)
	ease := func(start, change float32) float32 {
		return easeOutCubic2(g.StateProgress(fudge), start, change)
	}

	switch g.State {
	case game.GameInitial, game.GamePaused:
		alpha = ease(0, 1)

	case game.GamePlaying, game.GameExiting:
		alpha = ease(1, -1)
	}

	// Don't render the menu if it is invisible.
	if alpha == 0 {
		return
	}

	gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &orthoProjectionViewMatrix[0])
	gl.Uniform1f(grayscaleUniform, 0)
	gl.Uniform1f(alphaUniform, alpha)
	gl.Uniform1f(mixAmountUniform, 0)

	menu := g.Menu

	totalHeight := titleText.height*2 + float32(menuItemFontSize*len(menu.Items)*2)
	ty := (float32(winHeight) + totalHeight) / 2

	renderMenuItem := func(text *rendererText, focused bool) {
		tx := (float32(winWidth) - text.width) / 2
		ty -= text.height

		m := newScaleMatrix(text.width, text.height, 1)
		m = m.mult(newTranslationMatrix(tx, ty, 0))
		gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
		gl.Uniform1i(textureUniform, int32(text.texture)-1)

		var brightness float32
		switch {
		case focused && menu.Selected:
			brightness = pulse(menu.Pulse+fudge, 1, 1, 1)

		case focused:
			brightness = 1
		}
		gl.Uniform1f(brightnessUniform, brightness)
		textLineMesh.drawElements()

		ty -= text.height
	}

	renderMenuItem(titleText, false)
	for i, item := range menu.Items {
		renderMenuItem(menuItemText[item], menu.FocusedIndex == i)
	}
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

func logFatalIfErr(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", tag, err)
	}
}
