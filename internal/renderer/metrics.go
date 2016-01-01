package renderer

import (
	"math"

	"github.com/btmura/blockcillin/internal/game"
)

type metrics struct {
	g *game.Game
	b *game.Board
	s *game.Selector

	fudge float32

	globalTranslationY float32
	globalTranslationZ float32
	globalRotationY    float32
	cellRotationY      float32
}

func newMetrics(g *game.Game, fudge float32) *metrics {
	b := g.Board
	s := b.Selector

	cellRotationY := float32(2*math.Pi) / float32(b.CellCount)
	globalRotationY := cellRotationY/2 + boardRotationY(b, fudge)
	globalTranslationY := cellTranslationY * (4 + boardTranslationY(b, fudge))
	globalTranslationZ := float32(4)

	return &metrics{
		g: g,
		b: b,
		s: s,

		fudge: fudge,

		cellRotationY:      cellRotationY,
		globalRotationY:    globalRotationY,
		globalTranslationY: globalTranslationY,
		globalTranslationZ: globalTranslationZ,
	}
}

func (ms *metrics) blockMatrix(b *game.Block, x, y int) matrix4 {
	ty := ms.globalTranslationY + cellTranslationY*(-float32(y)+blockRelativeY(b, ms.fudge))

	ry := ms.globalRotationY + ms.cellRotationY*(-float32(x)-blockRelativeX(b, ms.fudge)+selectorRelativeX(ms.s, ms.fudge))
	yq := newAxisAngleQuaternion(yAxis, ry)
	qm := newQuaternionMatrix(yq.normalize())

	m := newTranslationMatrix(0, ty, ms.globalTranslationZ)
	m = m.mult(qm)
	return m
}
