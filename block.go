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
// By default, a block is a visible red block that is not moving.
type block struct {
	// state is the block's state. Use only within this file.
	state blockState

	// color is the block's color. Red by default.
	color blockColor

	// invisible is whether the block is invisible. Visible by default.
	invisible bool

	// step is the current step in any animation.
	step float32

	// pulse is used to advance any pulsing animations.
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
	if (l.state == blockStatic || l.state == blockClearing) && (r.state == blockStatic || r.state == blockClearing) {
		*l, *r = *r, *l
		l.state, r.state = blockSwappingFromRight, blockSwappingFromLeft
	}
}

// drop drops the upper block into the lower block.
func (u *block) drop(d *block) {
	if u.state == blockStatic && !u.invisible && d.state == blockStatic && d.invisible {
		*u, *d = *d, *u
		d.state = blockDroppingFromAbove
	}
}

func (b *block) flash() {
	b.state = blockFlashing
	b.pulse = 0
}

func (b *block) hasCracked() bool {
	return b.state == blockCracked
}

func (b *block) explode() {
	b.state = blockExploding
}

func (b *block) hasExploded() bool {
	return b.state == blockExploded
}

func (b *block) clear() {
	b.state = blockClearing
}

func (b *block) isClearable() bool {
	return b.state == blockStatic && !b.invisible
}

func (b *block) isCleared() bool {
	return b.state == blockStatic && b.invisible
}

func (b *block) update() {
	switch b.state {
	case blockSwappingFromLeft, blockSwappingFromRight:
		if b.step++; b.step >= numSwapSteps {
			b.state = blockStatic
			b.step = 0
		}

	case blockDroppingFromAbove:
		if b.step++; b.step >= numDropSteps {
			b.state = blockStatic
			b.step = 0
		}

	case blockFlashing:
		if b.step++; b.step >= numFlashSteps {
			b.state = blockCracking
			b.step = 0
		} else {
			b.pulse++
		}

	case blockCracking:
		if b.step++; b.step >= numCrackSteps {
			b.state = blockCracked
			b.step = 0
		}

	case blockExploding:
		if b.step++; b.step >= numExplodeSteps {
			b.state = blockExploded
			b.invisible = true
			b.step = 0
		}

	case blockClearing:
		if b.step++; b.step >= numClearSteps {
			b.state = blockStatic
			b.invisible = true
			b.step = 0
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

	default:
		return 0
	}
}

func (b *block) relativeY(fudge float32) float32 {
	switch b.state {
	case blockDroppingFromAbove:
		return linear(b.step+fudge, 1, -1, numDropSteps)

	default:
		return 0
	}
}

func (b *block) brightness(fudge float32) float32 {
	switch b.state {
	case blockFlashing:
		return pulse(b.pulse+fudge, 0, 0.5, 1.5)

	default:
		return 0
	}
}

func (b *block) alpha(fudge float32) float32 {
	switch b.state {
	case blockExploding:
		return linear(b.step+fudge, 1, -1, numExplodeSteps)

	default:
		if b.invisible {
			return 0
		}
		return 1
	}
}
