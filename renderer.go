package main

import (
	"image"
	"image/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type renderer struct {
	titleTextTexture   uint32
	newGameTextTexture uint32
}

func (r *renderer) init() error {
	return r.initMenu()
}

func (r *renderer) initMenu() error {
	font, err := freetype.ParseFont(MustAsset("data/Orbitron Medium.ttf"))
	if err != nil {
		return err
	}

	r.titleTextTexture, err = createTextTexture(font, "b l o c k c i l l i n", gl.TEXTURE1)
	if err != nil {
		return err
	}

	r.newGameTextTexture, err = createTextTexture(font, "N E W   G A M E", gl.TEXTURE2)
	if err != nil {
		return err
	}

	return nil
}

func createTextTexture(f *truetype.Font, text string, textureUnit uint32) (uint32, error) {
	rgba, err := createTextImage(f, text)
	if err != nil {
		return 0, err
	}

	texture, err := createTexture(textureUnit, rgba)
	if err != nil {
		return 0, err
	}

	return texture, nil
}

func createTextImage(f *truetype.Font, text string) (*image.RGBA, error) {
	fg, bg := image.White, image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, 128, 128))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetFont(f)
	c.SetDPI(72)
	c.SetFontSize(12)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(12)>>6))
	if _, err := c.DrawString(text, pt); err != nil {
		return nil, err
	}

	return rgba, nil
}

func (r *renderer) render() {
	r.renderMenu()
}

func (r *renderer) renderMenu() {
	m := newScaleMatrix(5, 5, 5)
	m = m.mult(newTranslationMatrix(0, 0, 0))
	gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])
	gl.Uniform1i(textureUniform, int32(r.titleTextTexture)-1)
	menuMesh.drawElements()

	m = newScaleMatrix(5, 5, 5)
	m = m.mult(newTranslationMatrix(0, -1, 0))
	gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])
	gl.Uniform1i(textureUniform, int32(r.newGameTextTexture)-1)
	menuMesh.drawElements()
}
