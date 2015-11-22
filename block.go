package main

// numAlphaSteps is the number of steps a move animation takes.
const numAlphaSteps float32 = 0.5 / secPerUpdate

type block struct {
	state     blockState
	color     blockColor
	moveStep  float32
	alphaStep float32
}

type blockState int32

const (
	blockStatic blockState = iota
	blockSwapFromLeft
	blockSwapFromRight
	blockCleared
)

func (b *block) clear() {
	if b.isClearable() {
		b.state = blockCleared
	}
}

func (b *block) swapFromLeft() {
	if b.state == blockStatic {
		b.state = blockSwapFromLeft
	}
}

func (b *block) swapFromRight() {
	if b.state == blockStatic {
		b.state = blockSwapFromRight
	}
}

func (b *block) isClearable() bool {
	return b.state == blockStatic
}

func (b *block) isSwappable() bool {
	return b.state != blockSwapFromLeft && b.state != blockSwapFromRight
}

func (b *block) update() {
	switch b.state {
	case blockSwapFromLeft, blockSwapFromRight:
		if b.moveStep++; b.moveStep >= numMoveSteps {
			b.state = blockStatic
			b.moveStep = 0
		}

	case blockCleared:
		if b.alphaStep < numAlphaSteps {
			b.alphaStep++
		}
	}
}

func (b *block) getX(fudge float32) float32 {
	move := func(start, delta float32) float32 {
		return linear(b.moveStep+fudge, start, delta, numMoveSteps)
	}

	switch b.state {
	case blockSwapFromLeft:
		return move(-1, 1)

	case blockSwapFromRight:
		return move(1, -1)
	}

	return 0
}

func (b *block) getAlpha(fudge float32) float32 {
	if b.state == blockCleared {
		return linear(b.alphaStep, 1, -1, numAlphaSteps)
	}
	return 1.0
}
