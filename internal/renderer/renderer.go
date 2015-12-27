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
	menuTitleTextColor = color.White

	menuItemFontSize  = 36
	menuItemTextColor = color.Gray{100}

	hudFontSize  = 20
	hudTextColor = color.White
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
	blockMeshes    = map[game.BlockColor]*mesh{}
	fragmentMeshes = map[game.BlockColor][4]*mesh{}
	textLineMesh   *mesh
)

var (
	boardTexture  uint32
	menuTitleText = map[game.MenuTitle]*renderableText{}
	menuItemText  = map[game.MenuItem]*renderableText{}
	hudItemText   = map[game.HUDItem]*renderableText{}
	hudRuneText   = map[rune]*renderableText{}
)

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

	if err := initMeshes(); err != nil {
		return err
	}

	if err := initTextures(); err != nil {
		return err
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

func initMeshes() error {
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

	mm := func(id string) *mesh {
		if err != nil {
			return nil
		}
		m, ok := meshMap[id]
		if !ok {
			err = fmt.Errorf("mesh not found: %s", id)
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

	if err != nil {
		return err
	}
	return nil
}

func initTextures() error {
	var textureUnit uint32 = gl.TEXTURE0
	var err error

	boardTexture, err = createAssetTexture(textureUnit, "data/texture.png")
	if err != nil {
		return err
	}
	textureUnit++

	font, err := freetype.ParseFont(asset.MustAsset("data/Orbitron Medium.ttf"))
	if err != nil {
		return err
	}

	makeText := func(text string, size int, color color.Color) (rt *renderableText) {
		if err != nil {
			return nil
		}
		rt, err = createText(text, size, color, font, textureUnit)
		textureUnit++
		return
	}

	for title, text := range game.MenuTitleText {
		menuTitleText[title] = makeText(text, menuTitleFontSize, menuTitleTextColor)
	}
	for item, text := range game.MenuItemText {
		menuItemText[item] = makeText(text, menuItemFontSize, menuItemTextColor)
	}
	for item, text := range game.HUDItemText {
		hudItemText[item] = makeText(text, hudFontSize, hudTextColor)
	}

	runeStrs := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ":"}
	for _, str := range runeStrs {
		hudRuneText[[]rune(str)[0]] = makeText(str, hudFontSize, hudTextColor)
	}

	if err != nil {
		return err
	}
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

func Render(g *game.Game, fudge float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	if renderBoard(g, fudge) {
		renderHUD(g, fudge)
	}
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
