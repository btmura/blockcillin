package main

import "math/rand"

// numRiseSteps is the steps in the rising animation for one ring's height.
const numRiseSteps float32 = 5.0 / secPerUpdate

type board struct {
	state boardState

	rings []*ring

	// chains of blocks that are scheduled to be cleared.
	chains []*chain

	y float32

	// ringCount is how many rings the board has.
	ringCount int

	// cellCount is how many cells are in each ring.
	cellCount int

	// riseStep is the current step in the rise animation.
	riseStep float32
}

type boardState int32

const (
	boardStatic boardState = iota
	boardRising
)

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
		b.rings = append(b.rings, newRing(b.cellCount))
	}

	return b
}

func newRing(cellCount int) *ring {
	r := &ring{}
	for i := 0; i < cellCount; i++ {
		c := &cell{
			block: &block{
				color: blockColor(rand.Intn(int(blockColorCount))),
			},
		}
		r.cells = append(r.cells, c)
	}
	return r
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
	b.updateBlocks()
	b.clearChains()
	b.dropBlocks()

	// Stop rising if chains are being cleared.
	if len(b.chains) > 0 {
		b.state = boardStatic
	} else {
		b.state = boardRising
	}

	switch b.state {
	case boardRising:
		if b.riseStep++; b.riseStep >= numRiseSteps {
			b.state = boardRising
			b.riseStep = 0

			// Add new ring once we've risen one ring higher.
			b.rings = append(b.rings, newRing(b.cellCount))
			b.ringCount++
			b.y++
		}
	}
}

func (b *board) updateBlocks() {
	for i := 0; i < b.ringCount; i++ {
		r := b.rings[i]
		for j := 0; j < b.cellCount; j++ {
			c := r.cells[j]
			c.block.update()
		}
	}
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

func (b *board) renderY(fudge float32) float32 {
	return linear(b.riseStep+fudge, b.y, 1, numRiseSteps)
}
