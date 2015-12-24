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
	yAxis = vector3{0, 1, 0}

	cameraPosition = vector3{0, 5, 25}
	targetPosition = vector3{}
	up             = yAxis

	ambientLightColor     = [3]float32{0.5, 0.5, 0.5}
	directionalLightColor = [3]float32{0.5, 0.5, 0.5}
	directionalVector     = [3]float32{0.5, 0.5, 0.5}

	titleFontSize    = 54
	menuItemFontSize = 36

	titleTextColor    = color.White
	menuItemTextColor = color.Gray{100}
)

type Renderer struct {
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

	// SizeCallback is the callback that GLFW should call when resizing the window.
	SizeCallback func(width, height int)

	// width is the current window's width reported by the SizeCallback.
	width int

	// height is the current window's height reported by the SizeCallback.
	height int

	// perspectiveProjectionViewMatrix is the perspective projection view matrix uniform value.
	perspectiveProjectionViewMatrix matrix4

	// orthoProjectionViewMatrix is the ortho projection view matrix uniform value.
	orthoProjectionViewMatrix matrix4

	selectorMesh   *mesh
	blockMeshes    map[game.BlockColor]*mesh
	fragmentMeshes map[game.BlockColor][4]*mesh
	textLineMesh   *mesh

	boardTexture uint32
	titleText    *rendererText
	menuItemText map[game.MenuItem]*rendererText
}

type rendererText struct {
	texture uint32
	width   float32
	height  float32
}

func NewRenderer() *Renderer {
	logFatalIfErr("gl.Init", gl.Init())
	log.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	rr := &Renderer{}
	var err error

	rr.program, err = createProgram(asset.MustString("data/shader.vert"), asset.MustString("data/shader.frag"))
	logFatalIfErr("createProgram", err)
	gl.UseProgram(rr.program)

	mustUniform := func(name string) int32 {
		l, err := getUniformLocation(rr.program, name)
		logFatalIfErr("getUniformLocation", err)
		return l
	}

	rr.projectionViewMatrixUniform = mustUniform("u_projectionViewMatrix")
	rr.modelMatrixUniform = mustUniform("u_modelMatrix")
	rr.normalMatrixUniform = mustUniform("u_normalMatrix")
	rr.ambientLightColorUniform = mustUniform("u_ambientLightColor")
	rr.directionalLightColorUniform = mustUniform("u_directionalLightColor")
	rr.directionalVectorUniform = mustUniform("u_directionalVector")
	rr.textureUniform = mustUniform("u_texture")
	rr.grayscaleUniform = mustUniform("u_grayscale")
	rr.brightnessUniform = mustUniform("u_brightness")
	rr.alphaUniform = mustUniform("u_alpha")

	vm := newViewMatrix(cameraPosition, targetPosition, up)
	nm := vm.inverse().transpose()
	gl.UniformMatrix4fv(rr.normalMatrixUniform, 1, false, &nm[0])

	gl.Uniform3fv(rr.ambientLightColorUniform, 1, &ambientLightColor[0])
	gl.Uniform3fv(rr.directionalLightColorUniform, 1, &directionalLightColor[0])
	gl.Uniform3fv(rr.directionalVectorUniform, 1, &directionalVector[0])

	rr.SizeCallback = func(width, height int) {
		if rr.width == width && rr.height == height {
			return
		}

		log.Printf("window size changed (%dx%d -> %dx%d)", int(rr.width), int(rr.height), width, height)
		gl.Viewport(0, 0, int32(width), int32(height))

		// Calculate new perspective projection view matrix.
		rr.width, rr.height = width, height
		fw, fh := float32(width), float32(height)
		aspect := fw / fh
		fovRadians := float32(math.Pi) / 3
		rr.perspectiveProjectionViewMatrix = vm.mult(newPerspectiveMatrix(fovRadians, aspect, 1, 2000))

		// Calculate new ortho projection view matrix.
		rr.orthoProjectionViewMatrix = newOrthoMatrix(fw, fh, fw /* use width as depth */)
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

	rr.selectorMesh = mm("selector")
	rr.blockMeshes = map[game.BlockColor]*mesh{}
	rr.fragmentMeshes = map[game.BlockColor][4]*mesh{}
	for c, id := range colorObjIDs {
		rr.blockMeshes[c] = mm(id)
		rr.fragmentMeshes[c] = [4]*mesh{
			mm(id + "_north_west"),
			mm(id + "_north_east"),
			mm(id + "_south_east"),
			mm(id + "_south_west"),
		}
	}
	rr.textLineMesh = mm("text_line")

	rr.boardTexture, err = createAssetTexture(gl.TEXTURE0, "data/texture.png")
	logFatalIfErr("createAssetTexture", err)

	font, err := freetype.ParseFont(asset.MustAsset("data/Orbitron Medium.ttf"))
	logFatalIfErr("freetype.ParseFont", err)

	rr.titleText, err = createText(gl.TEXTURE1, font, "b l o c k c i l l i n", titleFontSize, titleTextColor)
	logFatalIfErr("createText", err)

	rr.menuItemText = map[game.MenuItem]*rendererText{}
	var textureUnit uint32 = gl.TEXTURE2
	for item, text := range game.MenuItemText {
		rr.menuItemText[item], err = createText(textureUnit, font, text, menuItemFontSize, menuItemTextColor)
		logFatalIfErr("createText", err)
		textureUnit++
	}

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 0)

	return rr
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

func (rr *Renderer) Render(g *game.Game, fudge float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	rr.renderBoard(g, fudge)
	rr.renderMenu(g, fudge)
}

func (rr *Renderer) renderBoard(g *game.Game, fudge float32) {
	if g.State != game.GamePlaying {
		return
	}

	gl.UniformMatrix4fv(rr.projectionViewMatrixUniform, 1, false, &rr.perspectiveProjectionViewMatrix[0])

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
		gl.UniformMatrix4fv(rr.modelMatrixUniform, 1, false, &m[0])

		rr.selectorMesh.drawElements()
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
		gl.Uniform1f(rr.brightnessUniform, bv)

		m := newScaleMatrix(sx, 1, 1)
		m = m.mult(blockMatrix(c.Block, x, y, fudge))
		gl.UniformMatrix4fv(rr.modelMatrixUniform, 1, false, &m[0])
		rr.blockMeshes[c.Block.Color].drawElements()
	}

	renderCellFragments := func(c *game.Cell, x, y int, fudge float32) {
		render := func(sc, rx, ry, rz float32, dir int) {
			m := newScaleMatrix(sc, sc, sc)
			m = m.mult(newTranslationMatrix(rx, ry, rz))
			m = m.mult(blockMatrix(c.Block, x, y, fudge))
			gl.UniformMatrix4fv(rr.modelMatrixUniform, 1, false, &m[0])
			rr.fragmentMeshes[c.Block.Color][dir].drawElements()
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
		gl.Uniform1f(rr.brightnessUniform, bv)
		gl.Uniform1f(rr.alphaUniform, av)

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
				gl.Uniform1f(rr.grayscaleUniform, easeInExpo(b.RiseStep+fudge, 1, -1, game.NumRiseSteps))
				gl.Uniform1f(rr.brightnessUniform, 0)
				gl.Uniform1f(rr.alphaUniform, 1)
				for x, c := range r.Cells {
					renderCell(c, x, y+b.RingCount, fudge)
				}

			case i == 1 && y == 1: // draw transparent objects
				gl.Uniform1f(rr.grayscaleUniform, 1)
				gl.Uniform1f(rr.brightnessUniform, 0)
				gl.Uniform1f(rr.alphaUniform, easeInExpo(b.RiseStep+fudge, 0, 1, game.NumRiseSteps))
				for x, c := range r.Cells {
					renderCell(c, x, y+b.RingCount, fudge)
				}
			}
		}
	}
}

func (rr *Renderer) renderMenu(g *game.Game, fudge float32) {
	alpha := float32(1)
	switch g.State {
	case game.GameInitial, game.GamePaused:
		alpha = g.StateProgress(fudge)

	case game.GamePlaying, game.GameExiting:
		alpha = 1.0 - g.StateProgress(fudge)
	}

	// Don't render the menu if it is invisible.
	if alpha == 0 {
		return
	}

	gl.Enable(gl.BLEND)
	gl.UniformMatrix4fv(rr.projectionViewMatrixUniform, 1, false, &rr.orthoProjectionViewMatrix[0])
	gl.Uniform1f(rr.grayscaleUniform, 0)
	gl.Uniform1f(rr.alphaUniform, alpha)

	menu := g.Menu

	totalHeight := rr.titleText.height*2 + float32(menuItemFontSize*len(menu.Items)*2)
	ty := (float32(rr.height) + totalHeight) / 2

	renderMenuItem := func(text *rendererText, selected bool) {
		tx := (float32(rr.width) - text.width) / 2
		ty -= text.height

		m := newScaleMatrix(text.width, text.height, 1)
		m = m.mult(newTranslationMatrix(tx, ty, 0))
		gl.UniformMatrix4fv(rr.modelMatrixUniform, 1, false, &m[0])
		gl.Uniform1i(rr.textureUniform, int32(text.texture)-1)

		brightness := float32(0)
		if selected {
			brightness = 1
		}
		gl.Uniform1f(rr.brightnessUniform, brightness)
		rr.textLineMesh.drawElements()

		ty -= text.height
	}

	renderMenuItem(rr.titleText, false)
	for i, item := range menu.Items {
		renderMenuItem(rr.menuItemText[item], menu.SelectedIndex == i)
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
