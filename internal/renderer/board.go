package renderer

import (
	"math"
	"strconv"

	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func renderBoard(g *game.Game, fudge float32) bool {
	if g.Board == nil {
		return false
	}

	b := g.Board

	metrics := newMetrics(g, fudge)

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

	gl.Uniform1i(textureUniform, int32(boardTexture)-1)

	for i := 0; i <= 2; i++ {
		gl.Uniform1f(grayscaleUniform, globalGrayscale)
		gl.Uniform1f(brightnessUniform, 0)
		gl.Uniform1f(alphaUniform, 1)

		if i == 0 {
			renderSelector(metrics)
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
						renderCellBlock(metrics, c, x, y)

					case game.BlockCracking, game.BlockCracked:
						renderCellFragments(metrics, c, x, y)
					}

				case 1: // draw transparent objects
					switch c.Block.State {
					case game.BlockExploding:
						renderCellFragments(metrics, c, x, y)
					}
					renderMarker(metrics, c.Marker, x, y)
				}
			}
		}
	}

	// Render the spare rings.

	// Set brightness to zero for all spare rings.
	gl.Uniform1f(brightnessUniform, 0)

	for y, r := range b.SpareRings {
		// Set grayscale value. First spare rings becomes colored. Rest are gray.
		grayscale := float32(1)
		if y == 0 {
			grayscale = easeInExpo(b.RiseProgress(fudge), 1, -1)
		}
		if grayscale < globalGrayscale {
			grayscale = globalGrayscale
		}
		gl.Uniform1f(grayscaleUniform, grayscale)

		// Set alpha value. Last spare ring fades in. Rest are opaque.
		alpha := float32(1)
		if y == len(b.SpareRings)-1 {
			alpha = easeInExpo(b.RiseProgress(fudge), 0, 1)
		}
		gl.Uniform1f(alphaUniform, alpha)

		// Render the spare rings below the normal rings.
		for x, c := range r.Cells {
			renderCellBlock(metrics, c, x, y+b.RingCount)
		}
	}

	return true
}

func renderSelector(metrics *metrics) {
	m := metrics.selectorMatrix()
	gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
	selectorMesh.drawElements()
}

func renderCellBlock(metrics *metrics, c *game.Cell, x, y int) {
	sx := float32(1)
	bv := float32(0)

	switch c.Block.State {
	case game.BlockDroppingFromAbove:
		sx = linear(c.Block.StateProgress(metrics.fudge), 1, -0.5)
	case game.BlockFlashing:
		bv = pulse(metrics.g.GlobalPulse+metrics.fudge, 0, 0.5, 1.5)
	}
	gl.Uniform1f(brightnessUniform, bv)

	m := newScaleMatrix(sx, 1, 1)
	m = m.mult(metrics.blockMatrix(c.Block, x, y))
	gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
	blockMeshes[c.Block.Color].drawElements()
}

func renderCellFragments(metrics *metrics, c *game.Cell, x, y int) {
	const (
		nw = iota
		ne
		se
		sw
	)

	render := func(sc, rx, ry, rz float32, dir int) {
		m := newScaleMatrix(sc, sc, sc)
		m = m.mult(newTranslationMatrix(rx, ry, rz))
		m = m.mult(metrics.blockMatrix(c.Block, x, y))
		gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
		fragmentMeshes[c.Block.Color][dir].drawElements()
	}

	ease := func(start, change float32) float32 {
		return easeOutCubic(c.Block.StateProgress(metrics.fudge), start, change)
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
		j = pulse(c.Block.StateProgress(metrics.fudge), 0, 0.5, 1.5)
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

func renderMarker(metrics *metrics, m *game.Marker, x, y int) {
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
		ty := metrics.globalTranslationY + cellTranslationY*-float32(y) + easeOutCubic(m.StateProgress(metrics.fudge), 0, 0.5)
		tz := metrics.globalTranslationZ + cellTranslationZ/2 + 0.1

		ry := metrics.globalRotationY + metrics.cellRotationY*(-float32(x)+metrics.selectorRelativeX())
		yq := newAxisAngleQuaternion(yAxis, ry)
		qm := newQuaternionMatrix(yq.normalize())

		gl.Uniform1f(brightnessUniform, 0)
		gl.Uniform1f(alphaUniform, easeOutCubic(m.StateProgress(metrics.fudge), 1, -1))

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
