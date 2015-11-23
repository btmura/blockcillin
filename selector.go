package main

import "math"

// numMoveSteps is the number of steps a move animation takes.
const numMoveSteps float32 = 0.1 / secPerUpdate

type selector struct {
	// state is the state of the selector.
	state selectorState

	// x is the selector's current column.
	// It changes only after the move animation is complete.
	x int

	// y is the selector's current row.
	// It changes only after the move animation is complete.
	y int

	// moveStep is the current step in the move animation from 0 to numMoveSteps.
	moveStep float32

	// pulse is an increasing counter used to calculate the pulsing amount.
	pulse float32

	// ringCount is how many rings the board has.
	ringCount int

	// cellCount is how many cells are in a ring.
	cellCount int
}

type selectorState int32

const (
	selectorStatic selectorState = iota
	selectorMovingUp
	selectorMovingDown
	selectorMovingLeft
	selectorMovingRight
)

func newSelector(ringCount, cellCount int) *selector {
	return &selector{
		ringCount: ringCount,
		cellCount: cellCount,
	}
}

func (s *selector) moveUp() {
	if s.state == selectorStatic && s.y > 0 {
		s.state = selectorMovingUp
	}
}

func (s *selector) moveDown() {
	if s.state == selectorStatic && s.y < s.ringCount-1 {
		s.state = selectorMovingDown
	}
}

func (s *selector) moveLeft() {
	if s.state == selectorStatic {
		s.state = selectorMovingLeft
	}
}

func (s *selector) moveRight() {
	if s.state == selectorStatic {
		s.state = selectorMovingRight
	}
}

func (s *selector) canSwap() bool {
	return s.state == selectorStatic
}

func (s *selector) update() {
	updateMove := func() bool {
		if s.moveStep++; s.moveStep >= numMoveSteps {
			s.state = selectorStatic
			s.moveStep = 0
			return true
		}
		return false
	}

	switch s.state {
	case selectorMovingUp:
		if updateMove() {
			s.y--
		}

	case selectorMovingDown:
		if updateMove() {
			s.y++
		}

	case selectorMovingLeft:
		if updateMove() {
			if s.x--; s.x < 0 {
				s.x = s.cellCount - 1
			}
		}

	case selectorMovingRight:
		if updateMove() {
			if s.x++; s.x == s.cellCount {
				s.x = 0
			}
		}

	default:
		s.pulse++
	}
}

func (s *selector) renderX(fudge float32) float32 {
	sx := float32(s.x)
	move := func(delta float32) float32 {
		return linear(s.moveStep+fudge, sx, delta, numMoveSteps)
	}

	switch s.state {
	case selectorMovingLeft:
		return move(-1)

	case selectorMovingRight:
		return move(1)
	}

	return sx
}

func (s *selector) renderY(fudge float32) float32 {
	sy := float32(s.y)
	move := func(delta float32) float32 {
		return linear(s.moveStep+fudge, sy, delta, numMoveSteps)
	}

	switch s.state {
	case selectorMovingUp:
		return move(-1)

	case selectorMovingDown:
		return move(1)
	}

	return sy
}

func (s *selector) renderScale(fudge float32) float32 {
	return float32(1.0 + math.Sin(float64(s.pulse+fudge)*0.1)*0.025)
}
