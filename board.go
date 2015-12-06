package main

import (
	"log"
	"math/rand"
)

// numRiseSteps is the steps in the rising animation for one ring's height.
const numRiseSteps float32 = 5.0 / secPerUpdate

type board struct {
	// state is the board's state. Use only within this file.
	state boardState

	// selector is the selector the player uses to swap blocks.
	selector *selector

	// rings are the rings with cells with blocks that the player can swap.
	rings []*ring

	// spareRings are additional upcoming rings that the user cannot swap yet.
	spareRings []*ring

	// chains of blocks that are scheduled to be cleared.
	chains []*chain

	// y is offset in unit rings to vertically center the board.
	y int

	// ringCount is how many rings the board has.
	ringCount int

	// cellCount is the fixed number of cells each ring can have.
	cellCount int

	// filledRingCount is how many rings at the bottom to initially fill.
	filledRingCount int

	// spareRingCount is how many spare rings at the bottom will be shown.
	spareRingCount int

	// riseStep is the current step in the rise animation that rises one ring.
	riseStep float32
}

type boardState int32

const (
	boardStatic boardState = iota
	boardRising
	boardGameOver
)

type ring struct {
	cells []*cell
}

type cell struct {
	block *block
}

type boardConfig struct {
	ringCount       int
	cellCount       int
	filledRingCount int
	spareRingCount  int
}

func newBoard(bc *boardConfig) *board {
	b := &board{
		ringCount:       bc.ringCount,
		cellCount:       bc.cellCount,
		filledRingCount: bc.filledRingCount,
		spareRingCount:  bc.spareRingCount,
	}

	b.selector = newSelector(b.ringCount, b.cellCount)

	// Position the selector at the first filled ring.
	b.selector.y = b.ringCount - bc.filledRingCount

	for i := 0; i < b.ringCount; i++ {
		invisible := i < b.ringCount-bc.filledRingCount
		b.rings = append(b.rings, newRing(b.cellCount, invisible))
	}

	for i := 0; i < bc.spareRingCount; i++ {
		b.spareRings = append(b.spareRings, newRing(b.cellCount, false))
	}

	return b
}

func newRing(cellCount int, invisible bool) *ring {
	r := &ring{}
	for i := 0; i < cellCount; i++ {
		state := blockStatic
		if invisible {
			state = blockCleared
		}
		c := &cell{
			block: &block{
				state: state,
				color: blockColor(rand.Intn(int(blockColorCount))),
			},
		}
		r.cells = append(r.cells, c)
	}
	return r
}

func (b *board) moveLeft() {
	b.selector.moveLeft()
}

func (b *board) moveRight() {
	b.selector.moveRight()
}

func (b *board) moveDown() {
	b.selector.moveDown()
}

func (b *board) moveUp() {
	b.selector.moveUp()
}

func (b *board) swap() {
	x, y := b.selector.nextPosition()

	// Check bounds since the selector can move above the rings.
	if y < 0 {
		return
	}

	li, ri := x, (x+1)%b.cellCount
	lc, rc := b.cellAt(li, y), b.cellAt(ri, y)
	lc.block.swap(rc.block)
}

func (b *board) update() {
	if b.state == boardGameOver {
		return
	}

	b.selector.update()

	for _, r := range b.rings {
		for _, c := range r.cells {
			c.block.update()
		}
	}

	// Drop blocks first to prevent mid-air chains.
	b.dropBlocks()
	b.clearChains()

	// Stop rising if chains are being cleared.
	if len(b.chains) > 0 {
		b.state = boardStatic
	} else {
		b.state = boardRising
	}

	// Continually raise the board one ring an a time.
	switch b.state {
	case boardRising:
		if b.riseStep++; b.riseStep >= numRiseSteps {
			for _, c := range b.rings[0].cells {
				if c.block.state != blockCleared {
					b.state = boardGameOver
					log.Print("game over")
					return
				}
			}

			b.state = boardRising
			b.riseStep = 0

			b.rings = append(b.rings[1:], b.spareRings[0])
			b.spareRings = append(b.spareRings[1:], newRing(b.cellCount, false))
			if b.selector.y--; b.selector.y < 0 {
				b.selector.y = 0
			}
		}
	}
}

func (b *board) dropBlocks() {
	// Start at the bottom and drop blocks as we move up.
	// This allows a vertical stack of blocks to simultaneously drop.
	for y := len(b.rings) - 1; y >= 1; y-- {
		for x, dc := range b.rings[y].cells {
			uc := b.cellAt(x, y-1)
			uc.block.drop(dc.block)
		}
	}
}

func (b *board) clearChains() {
	// Find new chains and mark the blocks to be cleared soon.
	chains := findChains(b)
	for _, ch := range chains {
		for _, cc := range ch.cells {
			b.cellAt(cc.x, cc.y).block.state = blockFlashing
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
			case c.block.state == blockCracked:
				c.block.state = blockExploding
				finished = false
				break loop

			case c.block.state != blockExploded:
				finished = false
				break loop
			}
		}

		// Clear the blocks and remove the chain once all animations are done.
		if finished {
			for _, cc := range ch.cells {
				b.cellAt(cc.x, cc.y).block.state = blockClearPausing
			}
			b.chains = append(b.chains[:i], b.chains[i+1:]...)
			i--
		}
	}
}

func (b *board) cellAt(x, y int) *cell {
	return b.rings[y].cells[x]
}
