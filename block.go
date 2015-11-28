package main

const (
	// numSwapSteps is how many steps to stay in the swapping states.
	numSwapSteps = numMoveSteps

	// numDropSteps is how many steps to stay in the dropping state.
	numDropSteps float32 = 0.05 / secPerUpdate

	// numFlashSteps is how many steps to stay in the flashing state.
	numFlashSteps float32 = 0.5 / secPerUpdate

	// numCrackSteps is how many steps to say in the cracking state.
	numCrackSteps float32 = 0.1 / secPerUpdate

	// numExplodeSteps is how many steps to stay in the exploding state.
	numExplodeSteps float32 = 0.3 / secPerUpdate

	// numClearPauseSteps is how many steps to stay in the clear pausing state.
	numClearPauseSteps float32 = 0.2 / secPerUpdate
)

// block is a block that can be put into a cell.
type block struct {
	// state is the block's state.
	state blockState

	// color is the block's color. Red by default.
	color blockColor

	// step is the current step in any animation.
	step float32

	// pulse is the current step used in any pulsing animations.
	pulse float32
}

type blockState int32

const (
	// blockStatic is a visible and stationary block.
	blockStatic blockState = iota

	// blockSwappingFromLeft is a visible block swapping from the left.
	blockSwappingFromLeft

	// blockSwappingFromRight is a visible block swapping from the right.
	blockSwappingFromRight

	// blockDroppingFromAbove is a visible block dropping from above.
	blockDroppingFromAbove

	// blockFlashing is a block within a chain that is flashing.
	// Automatically goes to the blockCracking state.
	blockFlashing

	// blockCracking is a block within a chain that is cracking.
	// Automatically goes to the blockCracked state.
	blockCracking

	// blockCracked is a block within a chain has finished cracking.
	// Manually change the state to blockExploding when ready.
	blockCracked

	// blockExploding is a block within a chain that is exploding.
	// Automatically goes to the blockExploded state.
	blockExploding

	// blockExploded is a block within a chain that has finished exploding.
	// Manually change the state to blockClearPausing when ready.
	blockExploded

	// blockClearPausing is an invisible block but cannot be dropped into yet.
	// Automatically goes to the blockCleared state.
	blockClearPausing

	// blockCleared is a an invisible block.
	blockCleared
)

type blockColor int32

const (
	red blockColor = iota
	purple
	blue
	cyan
	green
	yellow
	blockColorCount
)

// swap swaps the left block with the right block.
func (l *block) swap(r *block) {
	if (l.state == blockStatic || l.state == blockClearPausing || l.state == blockCleared) &&
		(r.state == blockStatic || r.state == blockClearPausing || r.state == blockCleared) {
		*l, *r = *r, *l

		switch l.state {
		case blockStatic:
			l.state = blockSwappingFromRight
		case blockClearPausing, blockCleared:
			l.state = blockCleared
		}

		switch r.state {
		case blockStatic:
			r.state = blockSwappingFromLeft
		case blockClearPausing, blockCleared:
			r.state = blockCleared
		}

		l.reset()
		r.reset()
	}
}

// drop drops the upper block into the lower block.
func (u *block) drop(d *block) {
	if u.state == blockStatic && d.state == blockCleared {
		*u, *d = *d, *u
		u.state = blockCleared
		d.state = blockDroppingFromAbove
		u.reset()
		d.reset()
	}
}

// update advances the state machine one update.
func (b *block) update() {
	switch b.state {
	case blockSwappingFromLeft, blockSwappingFromRight:
		if b.step++; b.step >= numSwapSteps {
			b.state = blockStatic
			b.reset()
		}

	case blockDroppingFromAbove:
		if b.step++; b.step >= numDropSteps {
			b.state = blockStatic
			b.reset()
		}

	case blockFlashing:
		if b.step++; b.step >= numFlashSteps {
			b.state = blockCracking
			b.reset()
		} else {
			b.pulse++
		}

	case blockCracking:
		if b.step++; b.step >= numCrackSteps {
			b.state = blockCracked
			b.reset()
		}

	case blockExploding:
		if b.step++; b.step >= numExplodeSteps {
			b.state = blockExploded
			b.reset()
		}

	case blockClearPausing:
		if b.step++; b.step >= numClearPauseSteps {
			b.state = blockCleared
			b.reset()
		}
	}
}

// reset resets the animation state.
func (b *block) reset() {
	b.step = 0
	b.pulse = 0
}

func (b *block) relativeX(fudge float32) float32 {
	move := func(start, delta float32) float32 {
		return linear(b.step+fudge, start, delta, numSwapSteps)
	}

	switch b.state {
	case blockSwappingFromLeft:
		return move(-1, 1)

	case blockSwappingFromRight:
		return move(1, -1)
	}

	return 0
}

func (b *block) relativeY(fudge float32) float32 {
	if b.state == blockDroppingFromAbove {
		return linear(b.step+fudge, 1, -1, numDropSteps)
	}
	return 0
}

func (b *block) brightness(fudge float32) float32 {
	if b.state == blockFlashing {
		return pulse(b.pulse+fudge, 0, 0.5, 1.5)
	}
	return 0
}

func (b *block) alpha(fudge float32) float32 {
	switch b.state {
	case blockExploding:
		return linear(b.step+fudge, 1, -1, numExplodeSteps)

	case blockExploded, blockClearPausing, blockCleared:
		return 0
	}

	return 1
}
