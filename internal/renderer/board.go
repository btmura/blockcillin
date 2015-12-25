package renderer

import (
	"math"

	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

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