package renderer

import (
	"math"

	"github.com/btmura/blockcillin/internal/game"
)

const (
	cellTranslationY         = 2
	cellTranslationZ         = 2
	initialBoardTranslationY = 2 * cellTranslationY
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
}

func newMetrics(g *game.Game, fudge float32) *metrics {
	m := &metrics{
		g:                  g,
		b:                  g.Board,
		s:                  g.Board.Selector,
		fudge:              fudge,
		cellRotationY:      float32(2*math.Pi) / float32(g.Board.CellCount),
		globalTranslationZ: float32(4),
	}
	m.globalRotationY = m.cellRotationY/2 + m.boardRotationY()
	m.globalTranslationY = cellTranslationY * (4 + m.boardTranslationY())
	return m
}

func (m *metrics) selectorMatrix() matrix4 {
	sc := pulse(m.s.Pulse+m.fudge, 1.0, 0.025, 0.1)
	ty := m.globalTranslationY - cellTranslationY*m.selectorRelativeY()

	mtx := newScaleMatrix(sc, sc, sc)
	return mtx.mult(newTranslationMatrix(0, ty, m.globalTranslationZ))
}

func (m *metrics) blockMatrix(b *game.Block, x, y int) matrix4 {
	ty := m.globalTranslationY + cellTranslationY*(-float32(y)+m.blockRelativeY(b))

	ry := m.globalRotationY + m.cellRotationY*(-float32(x)-m.blockRelativeX(b)+m.selectorRelativeX())
	yq := newAxisAngleQuaternion(yAxis, ry)
	qm := newQuaternionMatrix(yq.normalize())

	mtx := newTranslationMatrix(0, ty, m.globalTranslationZ)
	mtx = mtx.mult(qm)
	return mtx
}

func (m *metrics) boardTranslationY() float32 {
	switch m.b.State {
	case game.BoardEntering:
		return easeOutCubic(m.b.StateProgress(m.fudge), -initialBoardTranslationY, initialBoardTranslationY)
	case game.BoardExiting:
		return easeOutCubic(m.b.StateProgress(m.fudge), m.b.Y, -initialBoardTranslationY)
	}
	return m.b.Y
}

func (m *metrics) boardRotationY() float32 {
	switch m.b.State {
	case game.BoardEntering:
		return easeOutCubic(m.b.StateProgress(m.fudge), math.Pi, -math.Pi)
	case game.BoardExiting:
		return easeOutCubic(m.b.StateProgress(m.fudge), 0, math.Pi)
	}
	return 0
}

func (m *metrics) selectorRelativeX() float32 {
	move := func(delta float32) float32 {
		return linear(m.s.StateProgress(m.fudge), float32(m.s.X), delta)
	}
	switch m.s.State {
	case game.SelectorMovingLeft:
		return move(-1)
	case game.SelectorMovingRight:
		return move(1)
	}
	return float32(m.s.X)
}

func (m *metrics) selectorRelativeY() float32 {
	move := func(delta float32) float32 {
		return linear(m.s.StateProgress(m.fudge), float32(m.s.Y), delta)
	}
	switch m.s.State {
	case game.SelectorMovingUp:
		return move(-1)
	case game.SelectorMovingDown:
		return move(1)
	}
	return float32(m.s.Y)
}

func (m *metrics) blockRelativeX(b *game.Block) float32 {
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

func (m *metrics) blockRelativeY(b *game.Block) float32 {
	if b.State == game.BlockDroppingFromAbove {
		return linear(b.StateProgress(m.fudge), 1, -1)
	}
	return 0
}
