package renderer

import (
	"math"
	"strconv"

	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	cellTranslationY         = 2
	cellTranslationZ         = 2
	initialBoardTranslationY = 2 * cellTranslationY
)

func renderBoard(g *game.Game, fudge float32) bool {
	if g.Board == nil {
		return false
	}

	const (
		nw = iota
		ne
		se
		sw
	)

	b := g.Board
	s := b.Selector

	gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &perspectiveProjectionViewMatrix[0])
	gl.Uniform3fv(mixColorUniform, 1, &blackColor[0])

	globalGrayscale := float32(1)
	globalDarkness := float32(0.8)
	var boardDarkness float32

	gameEase := func(start, change float32) float32 {
		return easeOutCubic(g.StateProgress(fudge), start, change)
	}
	boardEase := func(start, change float32) float32 {
		return easeOutCubic(b.StateProgress(fudge), start, change)
	}

	switch g.State {
	case game.GamePlaying:
		globalGrayscale = gameEase(1, -1)
		globalDarkness = gameEase(0.8, -0.8)

	case game.GamePaused:
		globalGrayscale = gameEase(0, 1)
		globalDarkness = gameEase(0, 0.8)

	case game.GameExiting:
		globalGrayscale = 1
		globalDarkness = gameEase(0.8, 1)
	}

	switch b.State {
	case game.BoardEntering:
		boardDarkness = boardEase(1, -1)

	case game.BoardExiting:
		boardDarkness = boardEase(0, 1)
	}

	finalDarkness := globalDarkness
	if finalDarkness < boardDarkness {
		finalDarkness = boardDarkness
	}
	gl.Uniform1f(mixAmountUniform, finalDarkness)

	cellRotationY := float32(2*math.Pi) / float32(b.CellCount)
	globalRotationY := cellRotationY/2 + boardRotationY(b, fudge)
	globalTranslationY := cellTranslationY * (4 + boardTranslationY(b, fudge))
	globalTranslationZ := float32(4)

	blockMatrix := func(b *game.Block, x, y int, fudge float32) matrix4 {
		ty := globalTranslationY + cellTranslationY*(-float32(y)+blockRelativeY(b, fudge))

		ry := globalRotationY + cellRotationY*(-float32(x)-blockRelativeX(b, fudge)+selectorRelativeX(s, fudge))
		yq := newAxisAngleQuaternion(yAxis, ry)
		qm := newQuaternionMatrix(yq.normalize())

		m := newTranslationMatrix(0, ty, globalTranslationZ)
		m = m.mult(qm)
		return m
	}

	renderSelector := func(fudge float32) {
		sc := pulse(s.Pulse+fudge, 1.0, 0.025, 0.1)
		ty := globalTranslationY - cellTranslationY*selectorRelativeY(s, fudge)

		m := newScaleMatrix(sc, sc, sc)
		m = m.mult(newTranslationMatrix(0, ty, globalTranslationZ))
		gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])

		selectorMesh.drawElements()
	}

	renderCellBlock := func(c *game.Cell, x, y int, fudge float32) {
		sx := float32(1)
		bv := float32(0)

		switch c.Block.State {
		case game.BlockDroppingFromAbove:
			sx = linear(c.Block.StateProgress(fudge), 1, -0.5)
		case game.BlockFlashing:
			bv = pulse(g.GlobalPulse+fudge, 0, 0.5, 1.5)
		}
		gl.Uniform1f(brightnessUniform, bv)

		m := newScaleMatrix(sx, 1, 1)
		m = m.mult(blockMatrix(c.Block, x, y, fudge))
		gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
		blockMeshes[c.Block.Color].drawElements()
	}

	renderMarker := func(m *game.Marker, x, y int, fudge float32) {
		switch m.State {
		case game.MarkerShowing:
			var val string
			switch {
			case m.ChainLevel > 0:
				val = "x" + strconv.Itoa(m.ChainLevel+1)

			case m.ComboLevel > 3:
				val = strconv.Itoa(m.ComboLevel)

			default:
				return
			}

			sc := float32(0.5)

			tx := -float32(len(val)-1) * sc
			ty := globalTranslationY + cellTranslationY*-float32(y) + easeOutCubic(m.StateProgress(fudge), 0, 0.5)
			tz := globalTranslationZ + cellTranslationZ/2 + 0.1

			ry := globalRotationY + cellRotationY*(-float32(x)+selectorRelativeX(s, fudge))
			yq := newAxisAngleQuaternion(yAxis, ry)
			qm := newQuaternionMatrix(yq.normalize())

			gl.Uniform1f(brightnessUniform, 0)
			gl.Uniform1f(alphaUniform, easeOutCubic(m.StateProgress(fudge), 1, -1))

			for _, rune := range val {
				text := markerRuneText[rune]

				m := newScaleMatrix(sc, sc, sc)
				m = m.mult(newTranslationMatrix(tx, ty, tz))
				m = m.mult(qm)
				gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])

				gl.Uniform1i(textureUniform, int32(text.texture)-1)
				squareMesh.drawElements()
				tx++
			}

			gl.Uniform1i(textureUniform, int32(boardTexture)-1)
		}
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
			return easeOutCubic(c.Block.StateProgress(fudge), start, change)
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
			j = pulse(c.Block.StateProgress(fudge), 0, 0.5, 1.5)
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

	gl.Uniform1i(textureUniform, int32(boardTexture)-1)

	for i := 0; i <= 2; i++ {
		gl.Uniform1f(grayscaleUniform, globalGrayscale)
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
						renderCellBlock(c, x, y, fudge)

					case game.BlockCracking, game.BlockCracked:
						renderCellFragments(c, x, y, fudge)
					}

				case 1: // draw transparent objects
					switch c.Block.State {
					case game.BlockExploding:
						renderCellFragments(c, x, y, fudge)
					}
					renderMarker(c.Marker, x, y, fudge)
				}
			}
		}

		for y, r := range b.SpareRings {
			switch {
			case i == 0 && y == 0: // draw opaque objects
				finalGrayscale := float32(1)
				if b.State == game.BoardRising {
					finalGrayscale = easeInExpo(b.RiseProgress(fudge), 1, -1)
					if globalGrayscale > finalGrayscale {
						finalGrayscale = globalGrayscale
					}
				}

				gl.Uniform1f(grayscaleUniform, finalGrayscale)
				gl.Uniform1f(brightnessUniform, 0)
				gl.Uniform1f(alphaUniform, 1)
				for x, c := range r.Cells {
					renderCellBlock(c, x, y+b.RingCount, fudge)
				}

			case i == 1 && y == 1: // draw transparent objects
				finalAlpha := float32(0)
				if b.State == game.BoardRising {
					finalAlpha = easeInExpo(b.RiseProgress(fudge), 0, 1)
				}

				gl.Uniform1f(grayscaleUniform, 1)
				gl.Uniform1f(brightnessUniform, 0)
				gl.Uniform1f(alphaUniform, finalAlpha)
				for x, c := range r.Cells {
					renderCellBlock(c, x, y+b.RingCount, fudge)
				}
			}
		}
	}

	return true
}

func boardTranslationY(b *game.Board, fudge float32) float32 {
	switch b.State {
	case game.BoardEntering:
		return easeOutCubic(b.StateProgress(fudge), -initialBoardTranslationY, initialBoardTranslationY)

	case game.BoardExiting:
		return easeOutCubic(b.StateProgress(fudge), b.Y, -initialBoardTranslationY)

	default:
		return b.Y
	}
}

func boardRotationY(b *game.Board, fudge float32) float32 {
	switch b.State {
	case game.BoardEntering:
		return easeOutCubic(b.StateProgress(fudge), math.Pi, -math.Pi)

	case game.BoardExiting:
		return easeOutCubic(b.StateProgress(fudge), 0, math.Pi)

	default:
		return 0
	}
}

func selectorRelativeX(s *game.Selector, fudge float32) float32 {
	move := func(delta float32) float32 {
		return linear(s.StateProgress(fudge), float32(s.X), delta)
	}

	switch s.State {
	case game.SelectorMovingLeft:
		return move(-1)

	case game.SelectorMovingRight:
		return move(1)
	}

	return float32(s.X)
}

func selectorRelativeY(s *game.Selector, fudge float32) float32 {
	move := func(delta float32) float32 {
		return linear(s.StateProgress(fudge), float32(s.Y), delta)
	}

	switch s.State {
	case game.SelectorMovingUp:
		return move(-1)
	case game.SelectorMovingDown:
		return move(1)
	}
	return float32(s.Y)
}

func blockRelativeX(b *game.Block, fudge float32) float32 {
	move := func(start, delta float32) float32 {
		return linear(b.StateProgress(fudge), start, delta)
	}

	switch b.State {
	case game.BlockSwappingFromLeft:
		return move(-1, 1)

	case game.BlockSwappingFromRight:
		return move(1, -1)
	}

	return 0
}

func blockRelativeY(b *game.Block, fudge float32) float32 {
	if b.State == game.BlockDroppingFromAbove {
		return linear(b.StateProgress(fudge), 1, -1)
	}
	return 0
}
