package main

const (
	// numSwapSteps is the number of steps to swap blocks.
	numSwapSteps = numMoveSteps

	// numClearSteps is the number of steps to clear a block.
	numClearSteps float32 = 0.5 / secPerUpdate

	// numDropSteps is the number of steps to drop blocks.
	numDropSteps = numMoveSteps
)

type block struct {
	// state is the block's state. Use only within this file.
	state blockState

	color     blockColor
	invisible bool
	swapStep  float32
	clearStep float32
	dropStep  float32
}

type blockState int32

const (
	blockStatic blockState = iota
	blockSwappingFromLeft
	blockSwappingFromRight
	blockDroppingFromAbove
	blockDroppingFromBelow
	blockClearing
)

func (b *block) clear() {
	b.state = blockClearing
}

func (b *block) swapFromLeft() {
	b.state = blockSwappingFromLeft
}

func (b *block) swapFromRight() {
	b.state = blockSwappingFromRight
}

func (b *block) dropFromAbove() {
	b.state = blockDroppingFromAbove
}

func (b *block) dropFromBelow() {
	b.state = blockDroppingFromBelow
}

func (b *block) isClearable() bool {
	return b.state == blockStatic && !b.invisible
}

func (b *block) isSwappable() bool {
	return b.state == blockStatic
}

func (b *block) isDroppable() bool {
	return b.state == blockStatic && !b.invisible
}

func (b *block) canReceiveDrop() bool {
	return b.state == blockStatic && b.invisible
}

func (b *block) update() {
	switch b.state {
	case blockSwappingFromLeft, blockSwappingFromRight:
		if b.swapStep++; b.swapStep >= numSwapSteps {
			b.state = blockStatic
			b.swapStep = 0
		}

	case blockClearing:
		if b.clearStep++; b.clearStep >= numClearSteps {
			b.state = blockStatic
			b.invisible = true
			b.clearStep = 0
		}

	case blockDroppingFromAbove:
		if b.dropStep++; b.dropStep >= numDropSteps {
			b.state = blockStatic
			b.dropStep = 0
		}

	case blockDroppingFromBelow:
		b.state = blockStatic
		b.invisible = true
	}
}

func (b *block) renderX(fudge float32) float32 {
	move := func(start, delta float32) float32 {
		return linear(b.swapStep+fudge, start, delta, numSwapSteps)
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

func (b *block) renderY(fudge float32) float32 {
	move := func(start, delta float32) float32 {
		return linear(b.dropStep+fudge, start, delta, numDropSteps)
	}

	switch b.state {
	case blockDroppingFromAbove:
		return move(1, -1)

	case blockDroppingFromBelow:
		return 0 // Nothing to animate.

	default:
		return 0
	}
}

func (b *block) renderAlpha(fudge float32) float32 {
	switch b.state {
	case blockClearing:
		return linear(b.clearStep+fudge, 1, -1, numClearSteps)

	default:
		if b.invisible {
			return 0
		}
		return 1
	}
}
