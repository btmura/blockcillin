package main

import (
	"log"
	"math/rand"
)

const (
	// numRiseSteps is the steps in the rising animation for one ring's height.
	numRiseSteps float32 = 5.0 / secPerUpdate

	// numSpareRings is how many spare rings to create.
	numSpareRings = 2
)

type board struct {
	// state is the board's state. Use only within this file.
	state boardState

	// selector is the selector that swaps blocks.
	selector *selector

	// rings containing cells which in turn contain blocks.
	rings []*ring

	spareRings []*ring

	// chains of blocks that are scheduled to be cleared.
	chains []*chain

	// y is offset in whole rings to render the board.
	y float32

	// ringCount is how many rings the board has.
	ringCount int

	// cellCount is how many cells each ring has.
	// Stays fixed once the game starts.
	cellCount int

	// riseStep is the current step in the rise animation that rises one ring.
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
		ringCount: 3,
		cellCount: 15,
	}

	b.selector = newSelector(b.ringCount, b.cellCount)

	for i := 0; i < b.ringCount; i++ {
		b.rings = append(b.rings, newRing(b.cellCount))
	}

	for i := 0; i < numSpareRings; i++ {
		b.spareRings = append(b.spareRings, newRing(b.cellCount))
	}

	b.y = float32(-b.ringCount - numSpareRings)

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

func (b *board) cellAt(x, y int) *cell {
	return b.rings[y].cells[x]
}

func (b *board) swap(x, y int) {
	li, ri := x, (x+1)%b.cellCount
	lc, rc := b.cellAt(li, y), b.cellAt(ri, y)

	// Swap cell contents and start animations.
	if lc.block.isSwappable() && rc.block.isSwappable() {
		lc.block, rc.block = rc.block, lc.block
		lc.block.swapFromRight()
		rc.block.swapFromLeft()
	}
}

func (b *board) update() {
	for _, r := range b.rings {
		for _, c := range r.cells {
			c.block.update()
		}
	}

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
		// Prune empty rings at the top.
		// Do this only when the board is rising without pending chains,
		// because the findChains algorithm stores relative coordinates.
	loop:
		for len(b.rings) > 0 {
			for _, c := range b.rings[0].cells {
				if !c.block.isCleared() {
					break loop
				}
			}

			b.rings = b.rings[1:]
			b.ringCount--
			b.selector.ringCount--
			b.y--

			// Adjust the selector up.
			b.selector.y--

			log.Printf("removed ring: %d", b.ringCount)
		}

		// Continually raise the board one ring an a time.
		if b.riseStep++; b.riseStep >= numRiseSteps {
			b.state = boardRising
			b.riseStep = 0

			// Transfer new spare ring and add a new spare.
			b.rings = append(b.rings, b.spareRings[0])
			b.spareRings = append(b.spareRings[1:], newRing(b.cellCount))

			b.ringCount++
			b.selector.ringCount++
			b.y++

			log.Printf("added ring: %d", b.ringCount)
		}
	}
}

func (b *board) clearChains() {
	// Find new chains and mark the blocks to be cleared soon.
	chains := findChains(b)
	for _, ch := range chains {
		for _, cc := range ch.cells {
			b.cellAt(cc.x, cc.y).block.flash()
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
			c := b.cellAt(cc.x, cc.y)
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
				b.cellAt(cc.x, cc.y).block.clear()
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
			uc := b.cellAt(x, y-1)

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
