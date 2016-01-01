package game

import (
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

	// matches contains matches that are being cleared.
	matches []*match

	// chainLinks contains links to continue chains.
	chainLinks []*chainLink

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

	// swapIDCounter is the next non-zero swap ID to set on the next swapped blocks.
	swapIDCounter int
}

type Ring struct {
	Cells []*Cell
}

type Cell struct {
	Block  *Block
	Marker *Marker
}

//go:generate stringer -type=BoardState
type BoardState int32

const (
	BoardEntering BoardState = iota
	BoardRising
	BoardGameOver
	BoardExiting
)

var boardStateSteps = map[BoardState]float32{
	BoardEntering: 2.0 / SecPerUpdate,
	BoardGameOver: 2.0 / SecPerUpdate,
	BoardExiting:  2.0 / SecPerUpdate,
}

type chainLink struct {
	// matches contain matches that new matches must vertically drop on to.
	matches []*match

	// nextMatches is a temporary variable for the next matches.
	nextMatches []*match

	// level is the highest chain level assigned with the chain link.
	level int
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
			Marker: &Marker{},
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
	lc.Block.swap(rc.Block, b.nextSwapID())
}

func (b *Board) exit() {
	b.setState(BoardExiting)
}

func (b *Board) update() {
	b.newBlocksCleared = 0

	advance := func(nextState BoardState) {
		if b.step++; b.step >= boardStateSteps[b.State] {
			b.setState(nextState)
		}
	}

	switch b.State {
	case BoardEntering:
		advance(BoardRising)

	case BoardRising:
		// Update the state of the selector, blocks, and markers.
		b.Selector.update()
		for _, r := range b.Rings {
			for _, c := range r.Cells {
				c.Block.update()
				c.Marker.update()
			}
		}

		// Drop any blocks to prevent mid-air matches.
		b.dropBlocks()

		// Add new matches to the board and update any chains in progress.
		b.addNewMatches()

		// Update pending matches.
		b.updateMatches()

		// Reset swap IDs for stationary blocks after new matches have been found.
		for _, r := range b.Rings {
			for _, c := range r.Cells {
				if c.Block.State == BlockStatic && c.Block.swapID != 0 {
					c.Block.swapID = 0
				}
			}
		}

		// Don't rise if there are pending matches.
		if len(b.matches) > 0 {
			return
		}

		// Don't rise if there are blocks with certain states.
		for _, r := range b.Rings {
			for _, c := range r.Cells {
				if !blockStateRiseable[c.Block.State] {
					return
				}
			}
		}

		// Reset the chain links since we are rising again.
		b.chainLinks = nil

		// Reset drop IDs for stationary blocks since we are rising again.
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
			// Check that topmost ring is empty, so that it can be removed.
			for _, c := range b.Rings[0].Cells {
				if c.Block.State != BlockCleared {
					b.setState(BoardGameOver)
					return
				}
			}

			b.Y = 0

			// Trim off the topmost ring and add a new spare ring.
			b.Rings = append(b.Rings[1:], b.SpareRings[0])

			// Add a new spare ring, since one was taken away.
			b.SpareRings = append(b.SpareRings[1:], newRing(b.CellCount, false))

			// Adjust the selector down in case it was at the removed top ring.
			if b.Selector.Y--; b.Selector.Y < 0 {
				b.Selector.Y = 0
			}
		}

	case BoardGameOver:
		b.step++

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

func (b *Board) addNewMatches() {
	// Find new matches and append them to the overall list.
	matches := findGroupedMatches(b)
	b.matches = append(b.matches, matches...)

	var dirtyLinks []*chainLink

	for _, m := range matches {
		hasDroppedBlock := false
		for _, c := range m.cells {
			block := b.blockAt(c.x, c.y)
			block.State = BlockFlashing
			hasDroppedBlock = hasDroppedBlock || block.dropped
			b.newBlocksCleared++
			b.totalBlocksCleared++
		}

		var link *chainLink
		if hasDroppedBlock {
		linkLoop:
			for _, l := range b.chainLinks {
				for _, lm := range l.matches {
					for _, lc := range lm.cells {
						for _, mc := range m.cells {
							if lc.x == mc.x && lc.y <= mc.y {
								link = l
								break linkLoop
							}
						}
					}
				}
			}
		}

		if link == nil {
			link = &chainLink{}
			b.chainLinks = append(b.chainLinks, link)
		} else {
			link.level++
		}

		link.nextMatches = append(link.nextMatches, m)
		dirtyLinks = append(dirtyLinks, link)
		b.markerAt(m.cells[0].x, m.cells[0].y).show(len(m.cells), link.level)
	}

	for _, link := range dirtyLinks {
		link.matches = link.nextMatches
		link.nextMatches = nil
	}
}

func (b *Board) updateMatches() {
	// Update each match - clearing one block at a time.
	for i := 0; i < len(b.matches); i++ {
		m := b.matches[i]
		finished := true

	loop:
		// Animate each block one at a time. Break if it is still animating.
		for _, mc := range m.cells {
			block := b.blockAt(mc.x, mc.y)
			switch {
			case block.State == BlockCracked:
				block.State = BlockExploding
				audio.Play(audio.SoundClear)
				finished = false
				break loop

			case block.State != BlockExploded:
				finished = false
				break loop
			}
		}

		// Clear the blocks and remove the chain once all animations are done.
		if finished {
			for _, mc := range m.cells {
				b.blockAt(mc.x, mc.y).State = BlockClearPausing
			}
			b.matches = append(b.matches[:i], b.matches[i+1:]...)
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

func (b *Board) StateDone() bool {
	return b.StateProgress(0) >= 1
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

func (b *Board) nextSwapID() int {
	for {
		nextID := b.swapIDCounter
		b.swapIDCounter++
		if nextID != 0 {
			return nextID
		}
	}
}

func (b *Board) cellAt(x, y int) *Cell {
	return b.Rings[y].Cells[x]
}

func (b *Board) blockAt(x, y int) *Block {
	return b.cellAt(x, y).Block
}

func (b *Board) markerAt(x, y int) *Marker {
	return b.cellAt(x, y).Marker
}
