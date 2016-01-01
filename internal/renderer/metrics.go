package renderer

import (
	"math"

	"github.com/btmura/blockcillin/internal/game"
)

const (
	cellTranslationY = 2
	cellTranslationZ = 2
)

type metrics struct {
	g     *game.Game
	b     *game.Board
	s     *game.Selector
	fudge float32

	globalTranslationY float32
	globalTranslationZ float32
	globalRotationY    float32
	cellRotationY      float32

	selectorMatrix matrix4
}

func newMetrics(g *game.Game, fudge float32) *metrics {
	b := g.Board
	s := b.Selector

	selectorRelativeX := func() float32 {
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

	selectorRelativeY := func() float32 {
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

	cellRotationY := float32(2*math.Pi) / float32(b.CellCount)

	boardRotationY := func() float32 {
		v := cellRotationY/2 + cellRotationY*selectorRelativeX()
		switch b.State {
		case game.BoardEntering:
			return v + easeOutCubic(b.StateProgress(fudge), math.Pi, -math.Pi)
		case game.BoardExiting:
			return v + easeOutCubic(b.StateProgress(fudge), 0, math.Pi)
		}
		return v
	}

	boardTranslationY := func() float32 {
		const initialBoardTranslationY = 2 * cellTranslationY
		switch b.State {
		case game.BoardEntering:
			return easeOutCubic(b.StateProgress(fudge), -initialBoardTranslationY, initialBoardTranslationY)
		case game.BoardExiting:
			return easeOutCubic(b.StateProgress(fudge), b.Y, -initialBoardTranslationY)
		}
		return b.Y
	}

	globalTranslationY := cellTranslationY * (4 + boardTranslationY())
	globalTranslationZ := float32(4)

	selectorMatrix := func() matrix4 {
		sc := pulse(s.Pulse+fudge, 1.0, 0.025, 0.1)
		ty := globalTranslationY - cellTranslationY*selectorRelativeY()
		mtx := newScaleMatrix(sc, sc, sc)
		return mtx.mult(newTranslationMatrix(0, ty, globalTranslationZ))
	}

	return &metrics{
		g:                  g,
		b:                  b,
		s:                  g.Board.Selector,
		fudge:              fudge,
		cellRotationY:      cellRotationY,
		globalTranslationY: globalTranslationY,
		globalTranslationZ: globalTranslationZ,
		globalRotationY:    boardRotationY(),
		selectorMatrix:     selectorMatrix(),
	}
}

func (m *metrics) blockMatrix(b *game.Block, x, y int) matrix4 {
	blockRelativeX := func() float32 {
		move := func(start, delta float32) float32 {
			return linear(b.StateProgress(m.fudge), start, delta)
		}
		switch b.State {
		case game.BlockSwappingFromLeft:
			return move(-1, 1)

		case game.BlockSwappingFromRight:
			return move(1, -1)
		}
		return 0
	}

	blockRelativeY := func() float32 {
		if b.State == game.BlockDroppingFromAbove {
			return linear(b.StateProgress(m.fudge), 1, -1)
		}
		return 0
	}

	ty := m.globalTranslationY - cellTranslationY*(float32(y)+blockRelativeY())
	ry := m.globalRotationY - m.cellRotationY*(float32(x)-blockRelativeX())

	mtx := newTranslationMatrix(0, ty, m.globalTranslationZ)
	mtx = mtx.mult(newQuaternionMatrix(newAxisAngleQuaternion(yAxis, ry).normalize()))
	return mtx
}
