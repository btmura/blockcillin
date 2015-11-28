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

	// numClearSteps is how many steps to stay in the clearing state.
	numClearSteps float32 = 0.2 / secPerUpdate
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
	blockStatic blockState = iota

	blockSwappingFromLeft
	blockSwappingFromRight
	blockDroppingFromAbove

	blockFlashing
	blockCracking
	blockCracked
	blockExploding
	blockExploded
	blockClearing
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
	if (l.state == blockStatic || l.state == blockClearing || l.state == blockCleared) &&
		(r.state == blockStatic || r.state == blockClearing || r.state == blockCleared) {
		*l, *r = *r, *l
		if l.state != blockCleared {
			l.state = blockSwappingFromRight
		}
		if r.state != blockCleared {
			r.state = blockSwappingFromLeft
		}
	}
}

// drop drops the upper block into the lower block.
func (u *block) drop(d *block) {
	if u.state == blockStatic && d.state == blockCleared {
		*u, *d = *d, *u
		d.state = blockDroppingFromAbove
	}
}

func (b *block) update() {
	reset := func() {
		b.step = 0
		b.pulse = 0
	}

	switch b.state {
	case blockSwappingFromLeft, blockSwappingFromRight:
		if b.step++; b.step >= numSwapSteps {
			b.state = blockStatic
			reset()
		}

	case blockDroppingFromAbove:
		if b.step++; b.step >= numDropSteps {
			b.state = blockStatic
			reset()
		}

	case blockFlashing:
		if b.step++; b.step >= numFlashSteps {
			b.state = blockCracking
			reset()
		} else {
			b.pulse++
		}

	case blockCracking:
		if b.step++; b.step >= numCrackSteps {
			b.state = blockCracked
			reset()
		}

	case blockExploding:
		if b.step++; b.step >= numExplodeSteps {
			b.state = blockExploded
			reset()
		}

	case blockClearing:
		if b.step++; b.step >= numClearSteps {
			b.state = blockCleared
			reset()
		}
	}
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

	case blockExploded, blockClearing, blockCleared:
		return 0
	}

	return 1
}
