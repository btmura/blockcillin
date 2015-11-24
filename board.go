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
			c.block.clearSoon()
		}
	}

	// Append these new chains to the list.
	b.chains = append(b.chains, chains...)

clear:
	for len(b.chains) > 0 {
		// Clear each cell in each chain one by one.
		ch := b.chains[0]
		for _, cc := range ch.cells {
			r := b.rings[cc.y]
			c := r.cells[cc.x]
			switch {
			case c.block.isClearingSoon():
				c.block.clear()
				break clear

			case !c.block.isClearingDone():
				break clear
			}
		}

		// Now really clear when the entire chain has finished animating.
		for _, cc := range ch.cells {
			r := b.rings[cc.y]
			c := r.cells[cc.x]
			c.block.clearImmediately()
		}

		// Remove the chain after clearing all of its cells.
		b.chains = b.chains[1:]
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
				uc.block.clearImmediately()
				dc.block.dropFromAbove()
			}
		}
	}
}
