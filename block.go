package main

const (
	// numSwapSteps is the number of steps to swap blocks.
	numSwapSteps = numMoveSteps

	// numClearSteps is the number of steps to clear a block.
	numClearSteps float32 = 0.5 / secPerUpdate

	// numDropSteps is the number of steps to drop blocks.
	numDropSteps = numMoveSteps / 2
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

	// swapStep is the current step in the swap animation.
	swapStep float32

	// clearStep is the current step in the clear animation.
	clearStep float32

	// dropStep is the current step in the drop animation.
	dropStep float32
}

type blockState int32

const (
	blockStatic blockState = iota
	blockSwappingFromLeft
	blockSwappingFromRight
	blockDroppingFromAbove
	blockClearing
)

func (b *block) clear() {
	b.state = blockClearing
}

func (b *block) clearImmediately() {
	b.state = blockStatic
	b.invisible = true
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

func (b *block) isClearable() bool {
	return b.state == blockStatic && !b.invisible
}

func (b *block) isSwappable() bool {
	return b.state == blockStatic
}

func (b *block) isDroppable() bool {
	return b.state == blockStatic && !b.invisible
}

func (b *block) isDropReady() bool {
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
	}
}

func (b *block) renderX(fudge float32) float32 {
	move := func(start, delta float32) float32 {
		return ease(b.swapStep+fudge, start, delta, numSwapSteps)
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
	switch b.state {
	case blockDroppingFromAbove:
		return ease(b.dropStep+fudge, 1, -1, numDropSteps)

	default:
		return 0
	}
}

func (b *block) renderAlpha(fudge float32) float32 {
	switch b.state {
	case blockClearing:
		return ease(b.clearStep+fudge, 1, -1, numClearSteps)

	default:
		if b.invisible {
			return 0
		}
		return 1
	}
}
