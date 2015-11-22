package main

import "math/rand"

type blockColor int32

const (
	red blockColor = iota
	purple
	blue
	cyan
	green
	yellow
	blockColorCount
)

type board struct {
	rings []*ring

	// ringCount is how many rings the board has.
	ringCount int

	// cellCount is how many cells are in each ring.
	cellCount int
}

type ring struct {
	cells []*cell
}

type cell struct {
	block *block
}

func newBoard() *board {
	b := &board{
		ringCount: 5,
		cellCount: 15,
	}

	for i := 0; i < b.ringCount; i++ {
		r := &ring{}
		for j := 0; j < b.cellCount; j++ {
			c := &cell{
				block: &block{
					color: blockColor(rand.Intn(int(blockColorCount))),
				},
			}
			r.cells = append(r.cells, c)
		}
		b.rings = append(b.rings, r)
	}

	return b
}

func (b *board) update() {
	for i := 0; i < b.ringCount; i++ {
		r := b.rings[i]
		for j := 0; j < b.cellCount; j++ {
			c := r.cells[j]
			c.block.update()
		}
	}

	for _, ch := range findChains(b) {
		for _, c := range ch.cells {
			r := b.rings[c.y]
			c := r.cells[c.x]
			c.block.clear()
		}
	}
}

func (b *board) swap(x, y int) {
	r := b.rings[y]
	li, ri := x, (x+1)%b.cellCount
	lc, rc := r.cells[li], r.cells[ri]

	// Swap cell contents.
	if lc.block.isSwappable() && rc.block.isSwappable() {
		lc.block, rc.block = rc.block, lc.block
		lc.block.swapFromRight()
		rc.block.swapFromLeft()
	}
}
