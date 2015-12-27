package game

import (
	"log"
	"math/rand"

	"github.com/btmura/blockcillin/internal/audio"
)

type Board struct {
	// State is the board's state. Use only within this file.
	State BoardState

	// Selector is the selector the player uses to swap blocks.
	Selector *Selector

	// Rings are the rings with cells with blocks that the player can swap.
	Rings []*Ring

	// SpareRings are additional upcoming rings that the user cannot swap yet.
	SpareRings []*Ring

	// Y is offset in unit rings to vertically center the board.
	Y int

	// RingCount is how many rings the board has.
	RingCount int

	// CellCount is the fixed number of cells each ring can have.
	CellCount int

	// chains of blocks that are scheduled to be cleared.
	chains []*chain

	// filledRingCount is how many rings at the bottom to initially fill.
	filledRingCount int

	// spareRingCount is how many spare rings at the bottom will be shown.
	spareRingCount int

	// step is the current step in the rise animation that rises one ring.
	step float32
}

type Ring struct {
	Cells []*Cell
}

type Cell struct {
	Block *Block
}

type BoardState int32

const (
	BoardEntering BoardState = iota
	BoardRising
	BoardGameOver
)

var boardStateSteps = map[BoardState]float32{
	BoardEntering: 2.0 / SecPerUpdate,
	BoardRising:   5.0 / SecPerUpdate,
}

type boardConfig struct {
	ringCount       int
	cellCount       int
	filledRingCount int
	spareRingCount  int
}

func newBoard(bc *boardConfig) *Board {
	b := &Board{
		RingCount:       bc.ringCount,
		CellCount:       bc.cellCount,
		filledRingCount: bc.filledRingCount,
		spareRingCount:  bc.spareRingCount,
	}

	b.Selector = newSelector(b.RingCount, b.CellCount)

	// Position the selector at the first filled ring.
	b.Selector.Y = b.RingCount - bc.filledRingCount

	for i := 0; i < b.RingCount; i++ {
		invisible := i < b.RingCount-bc.filledRingCount
		b.Rings = append(b.Rings, newRing(b.CellCount, invisible))
	}

	for i := 0; i < bc.spareRingCount; i++ {
		b.SpareRings = append(b.SpareRings, newRing(b.CellCount, false))
	}

	return b
}

func newRing(cellCount int, invisible bool) *Ring {
	r := &Ring{}
	for i := 0; i < cellCount; i++ {
		state := BlockStatic
		if invisible {
			state = BlockCleared
		}
		c := &Cell{
			Block: &Block{
				State: state,
				Color: BlockColor(rand.Intn(int(BlockColorCount))),
			},
		}
		r.Cells = append(r.Cells, c)
	}
	return r
}

func (b *Board) moveLeft() {
	if b.State != BoardRising {
		return
	}
	b.Selector.moveLeft()
}

func (b *Board) moveRight() {
	if b.State != BoardRising {
		return
	}
	b.Selector.moveRight()

}

func (b *Board) moveDown() {
	if b.State != BoardRising {
		return
	}
	b.Selector.moveDown()
}

func (b *Board) moveUp() {
	if b.State != BoardRising {
		return
	}
	b.Selector.moveUp()
}

func (b *Board) swap() {
	if b.State != BoardRising {
		return
	}

	x, y := b.Selector.nextPosition()

	// Check y since the selector can move above the rings.
	if y < 0 {
		return
	}

	li, ri := x, (x+1)%b.CellCount
	lc, rc := b.cellAt(li, y), b.cellAt(ri, y)
	lc.Block.swap(rc.Block)
}

func (b *Board) update() {
	switch b.State {
	case BoardEntering:
		if b.step++; b.step >= boardStateSteps[b.State] {
			b.setState(BoardRising)
		}

	case BoardRising:
		b.Selector.update()
		for _, r := range b.Rings {
			for _, c := range r.Cells {
				c.Block.update()
			}
		}

		// Drop blocks before clearing to prevent mid-air chains.
		b.dropBlocks()
		b.clearChains()

		// Don't rise if there are pending chains.
		if len(b.chains) > 0 {
			break
		}

		if b.step++; b.step >= boardStateSteps[b.State] {
			// Continually raise the board one ring an a time.
			b.setState(BoardRising)

			for _, c := range b.Rings[0].Cells {
				if c.Block.State != BlockCleared {
					b.State = BoardGameOver
					log.Print("game over")
					return
				}
			}

			b.Rings = append(b.Rings[1:], b.SpareRings[0])
			b.SpareRings = append(b.SpareRings[1:], newRing(b.CellCount, false))
			if b.Selector.Y--; b.Selector.Y < 0 {
				b.Selector.Y = 0
			}
		}
	}
}

func (b *Board) dropBlocks() {
	// Start at the bottom and drop blocks as we move up.
	// This allows a vertical stack of blocks to simultaneously drop.
	for y := len(b.Rings) - 1; y >= 1; y-- {
		for x, dc := range b.Rings[y].Cells {
			uc := b.cellAt(x, y-1)
			uc.Block.drop(dc.Block)
		}
	}
}

func (b *Board) clearChains() {
	// Find new chains and mark the blocks to be cleared soon.
	chains := findChains(b)
	for _, ch := range chains {
		for _, cc := range ch.cells {
			b.cellAt(cc.x, cc.y).Block.State = BlockFlashing
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
			case c.Block.State == BlockCracked:
				c.Block.State = BlockExploding
				audio.Play(audio.SoundClear)
				finished = false
				break loop

			case c.Block.State != BlockExploded:
				finished = false
				break loop
			}
		}

		// Clear the blocks and remove the chain once all animations are done.
		if finished {
			for _, cc := range ch.cells {
				b.cellAt(cc.x, cc.y).Block.State = BlockClearPausing
			}
			b.chains = append(b.chains[:i], b.chains[i+1:]...)
			i--
		}
	}
}

func (b *Board) StateProgress(fudge float32) float32 {
	totalSteps := boardStateSteps[b.State]
	if totalSteps == 0 {
		return 1
	}

	if p := (b.step + fudge) / totalSteps; p < 1 {
		return p
	}
	return 1
}

func (b *Board) setState(state BoardState) {
	b.State = state
	b.step = 0
}

func (b *Board) cellAt(x, y int) *Cell {
	return b.Rings[y].Cells[x]
}
