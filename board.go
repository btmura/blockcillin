package main

import "math/rand"

type board struct {
	rings []*ring

	// chains of blocks that are scheduled to be cleared.
	chains []*chain

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

func (b *board) swap(x, y int) {
	r := b.rings[y]
	li, ri := x, (x+1)%b.cellCount
	lc, rc := r.cells[li], r.cells[ri]

	// Swap cell contents and start animations.
	if lc.block.isSwappable() && rc.block.isSwappable() {
		lc.block, rc.block = rc.block, lc.block
		lc.block.swapFromRight()
		rc.block.swapFromLeft()
	}
}

func (b *board) update() {
	for i := 0; i < b.ringCount; i++ {
		r := b.rings[i]
		for j := 0; j < b.cellCount; j++ {
			c := r.cells[j]
			c.block.update()
		}
	}

	b.clearChains()
	b.dropBlocks()
}

func (b *board) clearChains() {
	// Find new chains and mark the blocks to be cleared soon.
	chains := findChains(b)
	for _, ch := range chains {
		for _, cc := range ch.cells {
			r := b.rings[cc.y]
			c := r.cells[cc.x]
			c.block.flash()
		}
	}

	// Append these new chains to the list.
	b.chains = append(b.chains, chains...)

	// Advance each chain - clearing one block at a time.
	for i := 0; i < len(b.chains); i++ {
		ch := b.chains[i]
		finished := true

	loop:
		// Animate each block one at a time. Break if it is still animating.
		for _, cc := range ch.cells {
			r := b.rings[cc.y]
			c := r.cells[cc.x]
			switch {
			case c.block.hasCracked():
				c.block.explode()
				finished = false
				break loop

			case !c.block.hasExploded():
				finished = false
				break loop
			}
		}

		// Clear the blocks and remove the chain once all animations are done.
		if finished {
			for _, cc := range ch.cells {
				r := b.rings[cc.y]
				c := r.cells[cc.x]
				c.block.clear()
			}
			b.chains = append(b.chains[:i], b.chains[i+1:]...)
			i--
		}
	}
}

func (b *board) dropBlocks() {
	// Start at the bottom and drop blocks as we move up.
	// This allows a vertical stack of blocks to simultaneously drop.
	for y := len(b.rings) - 1; y >= 1; y-- {
		for x, dc := range b.rings[y].cells {
			uc := b.rings[y-1].cells[x]

			if uc.block.isDroppable() && dc.block.isDropReady() {
				// Swap cell contents and start animations.
				uc.block, dc.block = dc.block, uc.block
				uc.block.clear()
				dc.block.dropFromAbove()
			}
		}
	}
}
