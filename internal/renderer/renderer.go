package renderer

import (
	"bufio"
	"fmt"
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

	menuTitleFontSize  = 54
	menuItemFontSize   = 36
	menuTitleTextColor = color.White
	menuItemTextColor  = color.Gray{100}
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
	boardTexture  uint32
	menuTitleText map[game.MenuTitle]rendererText
	menuItemText  map[game.MenuItem]rendererText
)

type rendererText struct {
	texture uint32
	width   float32
	height  float32
}

func Init() error {
	if err := gl.Init(); err != nil {
		return err
	}

	log.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	vs, err := asset.String("data/shader.vert")
	if err != nil {
		return err
	}

	fs, err := asset.String("data/shader.frag")
	if err != nil {
		return err
	}

	program, err := createProgram(vs, fs)
	if err != nil {
		return err
	}
	gl.UseProgram(program)

	var shaderErr error
	uniform := func(name string) int32 {
		var loc int32
		loc, shaderErr = getUniformLocation(program, name)
		return loc
	}

	projectionViewMatrixUniform = uniform("u_projectionViewMatrix")
	modelMatrixUniform = uniform("u_modelMatrix")
	normalMatrixUniform = uniform("u_normalMatrix")
	ambientLightColorUniform = uniform("u_ambientLightColor")
	directionalLightColorUniform = uniform("u_directionalLightColor")
	directionalVectorUniform = uniform("u_directionalVector")
	textureUniform = uniform("u_texture")
	grayscaleUniform = uniform("u_grayscale")
	brightnessUniform = uniform("u_brightness")
	alphaUniform = uniform("u_alpha")
	mixColorUniform = uniform("u_mixColor")
	mixAmountUniform = uniform("u_mixAmount")

	if shaderErr != nil {
		return shaderErr
	}

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

	r, err := asset.Reader("data/meshes.obj")
	if err != nil {
		return err
	}

	objs, err := decodeObjs(r)
	if err != nil {
		return err
	}

	meshes := createMeshes(objs)
	meshMap := map[string]*mesh{}
	for i, m := range meshes {
		log.Printf("mesh %d: %s", i, m.id)
		meshMap[m.id] = m
	}

	var meshErr error
	mm := func(id string) *mesh {
		m, ok := meshMap[id]
		if !ok {
			meshErr = fmt.Errorf("mesh not found: %s", id)
			return nil
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

	if meshErr != nil {
		return meshErr
	}

	var textureUnit uint32 = gl.TEXTURE0

	boardTexture, err = createAssetTexture(textureUnit, "data/texture.png")
	if err != nil {
		return err
	}
	textureUnit++

	font, err := freetype.ParseFont(asset.MustAsset("data/Orbitron Medium.ttf"))
	if err != nil {
		return err
	}

	menuTitleText = map[game.MenuTitle]rendererText{}
	for title, text := range game.MenuTitleText {
		menuTitleText[title], err = createText(textureUnit, font, text, menuTitleFontSize, menuTitleTextColor)
		if err != nil {
			return err
		}
		textureUnit++
	}

	menuItemText = map[game.MenuItem]rendererText{}
	for item, text := range game.MenuItemText {
		menuItemText[item], err = createText(textureUnit, font, text, menuItemFontSize, menuItemTextColor)
		if err != nil {
			return err
		}
		textureUnit++
	}

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 0)

	return nil
}

func createAssetTexture(textureUnit uint32, name string) (uint32, error) {
	r, err := asset.Reader(name)
	if err != nil {
		return 0, err
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return createTexture(textureUnit, rgba)
}

func createText(textureUnit uint32, f *truetype.Font, text string, fontSize int, color color.Color) (rendererText, error) {
	rgba, w, h, err := createTextImage(f, text, fontSize, color)
	if err != nil {
		return rendererText{}, err
	}

	t, err := createTexture(textureUnit, rgba)
	if err != nil {
		return rendererText{}, err
	}
	return rendererText{t, w, h}, nil
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

func Terminate() {}

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
