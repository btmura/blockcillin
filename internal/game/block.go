package game

import "github.com/btmura/blockcillin/internal/audio"

// Block is a block that can be put into a cell.
type Block struct {
	// State is the block's state.
	State BlockState

	// Color is the block's color. Red by default.
	Color BlockColor

	// step is the current step in any animation.
	step float32
}

type BlockState int32

const (
	// BlockStatic is a visible and stationary block.
	BlockStatic BlockState = iota

	// BlockSwappingFromLeft is a visible block swapping from the left.
	BlockSwappingFromLeft

	// BlockSwappingFromRight is a visible block swapping from the right.
	BlockSwappingFromRight

	// BlockDroppingFromAbove is a visible block dropping from above.
	BlockDroppingFromAbove

	// BlockFlashing is a block within a chain that is flashing.
	// Automatically goes to the BlockCracking state.
	BlockFlashing

	// BlockCracking is a block within a chain that is cracking.
	// Automatically goes to the BlockCracked state.
	BlockCracking

	// BlockCracked is a block within a chain has finished cracking.
	// Manually change the state to BlockExploding when ready.
	BlockCracked

	// BlockExploding is a block within a chain that is exploding.
	// Automatically goes to the BlockExploded state.
	BlockExploding

	// BlockExploded is a block within a chain that has finished exploding.
	// Manually change the state to BlockClearPausing when ready.
	BlockExploded

	// BlockClearPausing is an invisible block but cannot be dropped into yet.
	// Automatically goes to the BlockCleared state.
	BlockClearPausing

	// BlockCleared is a an invisible block.
	BlockCleared
)

var blockStateSteps = map[BlockState]float32{
	BlockSwappingFromLeft:  0.1 / SecPerUpdate,
	BlockSwappingFromRight: 0.1 / SecPerUpdate,
	BlockDroppingFromAbove: 0.05 / SecPerUpdate,
	BlockFlashing:          0.5 / SecPerUpdate,
	BlockCracking:          0.1 / SecPerUpdate,
	BlockExploding:         0.4 / SecPerUpdate,
	BlockClearPausing:      0.2 / SecPerUpdate,
}

//go:generate stringer -type=BlockColor
type BlockColor int32

const (
	Red BlockColor = iota
	Purple
	Blue
	Cyan
	Green
	Yellow
	BlockColorCount
)

// swap swaps the left block with the right block.
func (l *Block) swap(r *Block) {
	if (l.State == BlockStatic || l.State == BlockClearPausing || l.State == BlockCleared) &&
		(r.State == BlockStatic || r.State == BlockClearPausing || r.State == BlockCleared) {
		*l, *r = *r, *l

		// swapped is whether an actual visible block was swapped.
		swapped := false

		switch l.State {
		case BlockStatic:
			l.setState(BlockSwappingFromRight)
			swapped = true
		case BlockClearPausing, BlockCleared:
			l.setState(BlockCleared)
		}

		switch r.State {
		case BlockStatic:
			r.setState(BlockSwappingFromLeft)
			swapped = true
		case BlockClearPausing, BlockCleared:
			r.setState(BlockCleared)
		}

		if swapped {
			audio.Play(audio.SoundSwap)
		}
	}
}

// drop drops the upper block into the lower block.
func (u *Block) drop(d *Block) {
	if u.State == BlockStatic && d.State == BlockCleared {
		*u, *d = *d, *u
		u.setState(BlockCleared)
		d.setState(BlockDroppingFromAbove)
	}
}

// update advances the state machine by one update.
func (b *Block) update() {
	advance := func(nextState BlockState) {
		if b.step++; b.step >= blockStateSteps[b.State] {
			b.setState(nextState)
		}
	}

	switch b.State {
	case BlockSwappingFromLeft, BlockSwappingFromRight, BlockDroppingFromAbove:
		advance(BlockStatic)

	case BlockFlashing:
		advance(BlockCracking)

	case BlockCracking:
		advance(BlockCracked)

	case BlockExploding:
		advance(BlockExploded)

	case BlockClearPausing:
		advance(BlockCleared)
	}
}

func (b *Block) StateProgress(fudge float32) float32 {
	totalSteps := blockStateSteps[b.State]
	if totalSteps == 0 {
		return 1
	}

	if p := (b.step + fudge) / totalSteps; p < 1 {
		return p
	}
	return 1
}

func (b *Block) setState(state BlockState) {
	b.State = state
	b.step = 0
}
