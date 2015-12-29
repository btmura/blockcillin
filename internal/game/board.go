package game

import (
	"log"
	"math/rand"

	"github.com/btmura/blockcillin/internal/audio"
)

type Board struct {
	// State is the board's state. Use only within this file.
	State BoardState

	// Rings are the rings with cells with blocks that the player can swap.
	Rings []*Ring

	// SpareRings are additional upcoming rings that the user cannot swap yet.
	SpareRings []*Ring

	// RingCount is how many rings the board has.
	RingCount int

	// CellCount is the fixed number of cells each ring can have.
	CellCount int

	// Selector is the selector the player uses to swap blocks.
	Selector *Selector

	// Y is vertical offset from 0 to 1 as the board rises one ring.
	Y float32

	// chains of blocks that are scheduled to be cleared.
	chains []*chain

	// chainHistory maps generation to slices of chains.
	// Generation X means X drops occurred to produce the chain.
	chainHistory [][]*chain

	// step is the current step in the rise animation that rises one ring.
	step float32

	// newBlocksCleared is the number of blocks cleared in the last update.
	newBlocksCleared int

	// totalBlocksCleared is the number of blocks cleared across all updates.
	totalBlocksCleared int

	// riseRate is how much to raise the board on each update.
	// It increases as the player scores more points.
	riseRate float32

	// useManualRiseRate is whether to use the manual rise rate on each update.
	useManualRiseRate bool
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
	BoardExiting
)

var boardStateSteps = map[BoardState]float32{
	BoardEntering: 2.0 / SecPerUpdate,
	BoardExiting:  2.0 / SecPerUpdate,
}

func newBoard(ringCount, cellCount, filledRingCount, spareRingCount int, riseRate float32) *Board {
	b := &Board{
		RingCount: ringCount,
		CellCount: cellCount,
		riseRate:  riseRate,
	}

	// Create the board's rings.
	//
	// 1. The board always has ringCount number of rows, but the top ones contain empty cells.
	// 2. As the board rises, empty top rings are pruned an replaced with spare ring rows.
	// 3. Spare ring rows are replenished as they are added to the board.

	for i := 0; i < b.RingCount; i++ {
		invisible := i < b.RingCount-filledRingCount
		b.Rings = append(b.Rings, newRing(b.CellCount, invisible))
	}

	for i := 0; i < spareRingCount; i++ {
		b.SpareRings = append(b.SpareRings, newRing(b.CellCount, false))
	}

	// Position the selector at the first filled ring.
	b.Selector = newSelector(b.RingCount, b.CellCount)
	b.Selector.Y = b.RingCount - filledRingCount

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

func (b *Board) exit() {
	b.setState(BoardExiting)
}

func (b *Board) update() {
	b.newBlocksCleared = 0

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
			return
		}

		// Don't rise if there are blocks with certain states.
		for _, r := range b.Rings {
			for _, c := range r.Cells {
				if !blockStateAllowRise[c.Block.State] {
					return
				}
			}
		}

		// Reset the chain history and dropped flags before rising.
		b.chainHistory = nil
		for _, r := range b.Rings {
			for _, c := range r.Cells {
				c.Block.dropped = false
			}
		}

		// Determine the rise rate.
		riseRate := b.riseRate
		if b.useManualRiseRate {
			riseRate = manualRiseRate
		}

		if b.Y += riseRate; b.Y > 1 {
			b.Y = 0

			// Check that topmost ring is empty, so that it can be removed.
			for _, c := range b.Rings[0].Cells {
				if c.Block.State != BlockCleared {
					b.State = BoardGameOver
					return
				}
			}

			// Trim off the topmost ring and add a new spare ring.
			b.Rings = append(b.Rings[1:], b.SpareRings[0])

			// Add a new spare ring, since one was taken away.
			b.SpareRings = append(b.SpareRings[1:], newRing(b.CellCount, false))

			// Adjust the selector down in case it was at the removed top ring.
			if b.Selector.Y--; b.Selector.Y < 0 {
				b.Selector.Y = 0
			}
		}

	case BoardExiting:
		b.step++
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
	// Find new chains of matching blocks.
	chains := findChains(b)

	var generations []int
	for _, ch := range chains {
		var hasDroppedBlock bool
		for _, cc := range ch.cells {
			// Start the clearing animation for each cell in the chain.
			block := b.cellAt(cc.x, cc.y).Block
			block.State = BlockFlashing

			hasDroppedBlock = hasDroppedBlock || block.dropped
			b.newBlocksCleared++
			b.totalBlocksCleared++
		}

		var foundGen int
		if hasDroppedBlock {
		genloop:
			for gen := len(b.chainHistory) - 1; gen >= 0; gen-- {
				for _, histChain := range b.chainHistory[gen] {
					for _, histCell := range histChain.cells {
						for _, cc := range ch.cells {
							if cc.x == histCell.x {
								foundGen = gen + 1
								break genloop
							}
						}
					}
				}
			}
		}

		generations = append(generations, foundGen)
	}

	var historyChanged bool
	for i, gen := range generations {
		for len(b.chainHistory) <= gen {
			b.chainHistory = append(b.chainHistory, nil)
		}
		b.chainHistory[gen] = append(b.chainHistory[gen], chains[i])
		historyChanged = true
	}

	if historyChanged {
		log.Print("history:")
		for gen, chains := range b.chainHistory {
			if len(chains) > 0 {
				log.Printf("\tgen %d: %d", gen, len(chains))
			}
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

func (b *Board) RiseProgress(fudge float32) float32 {
	if p := (b.Y + fudge/10.0); p < 1 {
		return p
	}
	return 1
}

func (b *Board) cellAt(x, y int) *Cell {
	return b.Rings[y].Cells[x]
}
