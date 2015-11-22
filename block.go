package main

type block struct {
	state    blockState
	x        int
	moveStep float32
	color    blockColor
}

type blockState int32

const (
	blockStatic blockState = iota
	blockMovingFromLeft
	blockMovingFromRight
	blockCleared
)

func (b *block) moveFromLeft() {
	if b.state == blockStatic {
		b.state = blockMovingFromLeft
		b.x = -1
	}
}

func (b *block) moveFromRight() {
	if b.state == blockStatic {
		b.state = blockMovingFromRight
		b.x = 1
	}
}

func (b *block) clear() {
	b.state = blockCleared
}

func (b *block) update() {
	updateMove := func() bool {
		if b.moveStep++; b.moveStep >= numMoveSteps {
			b.state = blockStatic
			b.moveStep = 0
			return true
		}
		return false
	}

	switch b.state {
	case blockMovingFromLeft:
		if updateMove() {
			b.x++
		}

	case blockMovingFromRight:
		if updateMove() {
			b.x--
		}
	}
}

func (b *block) isDrawable() bool {
	return b.state != blockCleared
}

func (b *block) getX(fudge float32) float32 {
	bx := float32(b.x)
	move := func(delta float32) float32 {
		return linear(b.moveStep+fudge, bx, delta, numMoveSteps)
	}

	switch b.state {
	case blockMovingFromLeft:
		return move(1)

	case blockMovingFromRight:
		return move(-1)
	}

	return bx
}
