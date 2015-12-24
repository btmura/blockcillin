package game

import "github.com/btmura/blockcillin/internal/audio"

const (
	// NumSwapSteps is how many steps to stay in the swapping states.
	NumSwapSteps = NumMoveSteps

	// NumDropSteps is how many steps to stay in the dropping state.
	NumDropSteps float32 = 0.05 / SecPerUpdate

	// numFlashSteps is how many steps to stay in the flashing state.
	numFlashSteps float32 = 0.5 / SecPerUpdate

	// numCrackSteps is how many steps to say in the cracking state.
	numCrackSteps float32 = 0.1 / SecPerUpdate

	// NumExplodeSteps is how many steps to stay in the exploding state.
	NumExplodeSteps float32 = 0.4 / SecPerUpdate

	// numClearPauseSteps is how many steps to stay in the clear pausing state.
	numClearPauseSteps float32 = 0.2 / SecPerUpdate
)

// Block is a block that can be put into a cell.
type Block struct {
	// State is the block's state.
	State BlockState

	// Color is the block's color. Red by default.
	Color BlockColor

	// Step is the current step in any animation.
	Step float32
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

		switch l.State {
		case BlockStatic:
			l.State = BlockSwappingFromRight
		case BlockClearPausing, BlockCleared:
			l.State = BlockCleared
		}

		switch r.State {
		case BlockStatic:
			r.State = BlockSwappingFromLeft
		case BlockClearPausing, BlockCleared:
			r.State = BlockCleared
		}

		l.reset()
		r.reset()
		audio.Play(audio.SoundSwap)
	}
}

// drop drops the upper block into the lower block.
func (u *Block) drop(d *Block) {
	if u.State == BlockStatic && d.State == BlockCleared {
		*u, *d = *d, *u
		u.State = BlockCleared
		d.State = BlockDroppingFromAbove
		u.reset()
		d.reset()
	}
}

// update advances the state machine by one update.
func (b *Block) update() {
	switch b.State {
	case BlockSwappingFromLeft, BlockSwappingFromRight:
		if b.Step++; b.Step >= NumSwapSteps {
			b.State = BlockStatic
			b.reset()
		}

	case BlockDroppingFromAbove:
		if b.Step++; b.Step >= NumDropSteps {
			b.State = BlockStatic
			b.reset()
		}

	case BlockFlashing:
		if b.Step++; b.Step >= numFlashSteps {
			b.State = BlockCracking
			b.reset()
		}

	case BlockCracking:
		if b.Step++; b.Step >= numCrackSteps {
			b.State = BlockCracked
			b.reset()
		}

	case BlockExploding:
		if b.Step++; b.Step >= NumExplodeSteps {
			b.State = BlockExploded
			b.reset()
		}

	case BlockClearPausing:
		if b.Step++; b.Step >= numClearPauseSteps {
			b.State = BlockCleared
			b.reset()
		}
	}
}

// reset resets the animation state.
func (b *Block) reset() {
	b.Step = 0
}
