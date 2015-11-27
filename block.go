package main

const (
	// numSwapSteps is the steps in the swapping animation.
	numSwapSteps = numMoveSteps

	// numFlashSteps is the steps in the flashing animation.
	numFlashSteps float32 = 0.5 / secPerUpdate

	// numCrackSteps is the steps in the cracking animation.
	numCrackSteps float32 = 0.1 / secPerUpdate

	// numExplodeSteps is the steps in the exploding animation.
	numExplodeSteps float32 = 0.3 / secPerUpdate

	// numDropSteps is the steps in the dropping animation.
	numDropSteps float32 = 0.05 / secPerUpdate
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
	if l.state == blockStatic && r.state == blockStatic {
		l.state, r.state = blockSwappingFromRight, blockSwappingFromLeft
		l.color, r.color = r.color, l.color
		l.invisible, r.invisible = r.invisible, l.invisible
	}
}

// drop drops the upper block with the lower block.
func (u *block) drop(l *block) {
	if u.state == blockStatic && !u.invisible && l.state == blockStatic && l.invisible {
		l.state = blockDroppingFromAbove
		l.color = u.color
		u.invisible, l.invisible = true, false
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
	b.state = blockStatic
	b.invisible = true
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
			b.step = 0
			b.invisible = true
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
