package main

import (
	"image"
	"image/draw"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type renderer struct {
	boardTexture       uint32
	titleTextTexture   uint32
	newGameTextTexture uint32
}

func (r *renderer) init() error {
	var err error

	r.boardTexture, err = createAssetTexture(gl.TEXTURE0, "data/texture.png")
	if err != nil {
		return err
	}

	font, err := freetype.ParseFont(MustAsset("data/Orbitron Medium.ttf"))
	if err != nil {
		return err
	}

	r.titleTextTexture, err = createTextTexture(gl.TEXTURE1, "b l o c k c i l l i n", font)
	if err != nil {
		return err
	}

	r.newGameTextTexture, err = createTextTexture(gl.TEXTURE2, "N E W   G A M E", font)
	if err != nil {
		return err
	}

	return nil
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

func createTextTexture(textureUnit uint32, text string, f *truetype.Font) (uint32, error) {
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

func (r *renderer) render(b *board, fudge float32) {
	r.renderMenu()
	r.renderBoard(b, fudge)
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

func (r *renderer) renderBoard(b *board, fudge float32) {
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
		gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])

		selectorMesh.drawElements()
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
		gl.Uniform1f(brightnessUniform, bv)

		m := newScaleMatrix(sx, 1, 1)
		m = m.mult(blockMatrix(c.block, x, y, fudge))
		gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])
		blockMeshes[c.block.color].drawElements()
	}

	renderCellFragments := func(c *cell, x, y int, fudge float32) {
		render := func(sc, rx, ry, rz float32, dir int) {
			m := newScaleMatrix(sc, sc, sc)
			m = m.mult(newTranslationMatrix(rx, ry, rz))
			m = m.mult(blockMatrix(c.block, x, y, fudge))
			gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])
			fragmentMeshes[c.block.color][dir].drawElements()
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
		gl.Uniform1f(brightnessUniform, bv)
		gl.Uniform1f(alphaUniform, av)

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

	gl.Uniform1i(textureUniform, int32(r.boardTexture)-1)

	for i := 0; i <= 2; i++ {
		gl.Uniform1f(grayscaleUniform, 0)
		gl.Uniform1f(brightnessUniform, 0)
		gl.Uniform1f(alphaUniform, 1)

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
				gl.Uniform1f(grayscaleUniform, easeInExpo(b.riseStep+fudge, 1, -1, numRiseSteps))
				gl.Uniform1f(brightnessUniform, 0)
				gl.Uniform1f(alphaUniform, 1)
				for x, c := range r.cells {
					renderCell(c, x, y+b.ringCount, fudge)
				}

			case i == 1 && y == 1: // draw transparent objects
				gl.Uniform1f(grayscaleUniform, 1)
				gl.Uniform1f(brightnessUniform, 0)
				gl.Uniform1f(alphaUniform, easeInExpo(b.riseStep+fudge, 0, 1, numRiseSteps))
				for x, c := range r.cells {
					renderCell(c, x, y+b.ringCount, fudge)
				}
			}
		}
	}
}
