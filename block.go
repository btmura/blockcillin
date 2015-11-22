package main

// numClearSteps is the number of steps to clear a block.
const numClearSteps float32 = 0.5 / secPerUpdate

type block struct {
	state     blockState
	color     blockColor
	invisible bool
	moveStep  float32
	clearStep float32
}

type blockState int32

const (
	blockStatic blockState = iota
	blockSwappingFromLeft
	blockSwappingFromRight
	blockClearing
)

func (b *block) clear() {
	if b.isClearable() {
		b.state = blockClearing
	}
}

func (b *block) swapFromLeft() {
	if b.isSwappable() {
		b.state = blockSwappingFromLeft
	}
}

func (b *block) swapFromRight() {
	if b.isSwappable() {
		b.state = blockSwappingFromRight
	}
}

func (b *block) isClearable() bool {
	return b.state == blockStatic && !b.invisible
}

func (b *block) isSwappable() bool {
	return b.state == blockStatic
}

func (b *block) update() {
	switch b.state {
	case blockSwappingFromLeft, blockSwappingFromRight:
		if b.moveStep++; b.moveStep >= numMoveSteps {
			b.state = blockStatic
			b.moveStep = 0
		}

	case blockClearing:
		if b.clearStep++; b.clearStep >= numClearSteps {
			b.state = blockStatic
			b.invisible = true
			b.clearStep = 0
		}
	}
}

func (b *block) getX(fudge float32) float32 {
	move := func(start, delta float32) float32 {
		return linear(b.moveStep+fudge, start, delta, numMoveSteps)
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

func (b *block) getAlpha(fudge float32) float32 {
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
