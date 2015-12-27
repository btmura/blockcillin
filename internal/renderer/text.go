package renderer

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type renderableText struct {
	text  string
	size  int
	color color.Color

	texture uint32
	width   float32
	height  float32
}

func createText(text string, size int, color color.Color, f *truetype.Font, textureUnit uint32) (*renderableText, error) {
	rgba, width, height, err := createTextImage(text, size, color, f)
	if err != nil {
		return nil, err
	}

	texture, err := createTexture(textureUnit, rgba)
	if err != nil {
		return nil, err
	}

	return &renderableText{
		text:  text,
		size:  size,
		color: color,

		texture: texture,
		width:   width,
		height:  height,
	}, nil
}

func (rt *renderableText) render(x, y float32) {
	m := newScaleMatrix(rt.width, rt.height, 1)
	m = m.mult(newTranslationMatrix(x, y, 0))
	gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
	gl.Uniform1i(textureUniform, int32(rt.texture)-1)
	textLineMesh.drawElements()
}

func createTextImage(text string, fontSize int, color color.Color, f *truetype.Font) (*image.RGBA, float32, float32, error) {
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
