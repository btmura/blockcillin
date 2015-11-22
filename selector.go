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

	// scale is the scale of the selector to make it pulse.
	scale float32

	// pulse is an increasing counter used to calculate the pulsing amount.
	pulse int

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

func newSelector(cellCount int) *selector {
	return &selector{cellCount: cellCount}
}

func (s *selector) moveUp() {
	if s.state == selectorStatic {
		s.state = selectorMovingUp
	}
}

func (s *selector) moveDown() {
	if s.state == selectorStatic {
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
		s.scale = float32(1.0 + math.Sin(float64(s.pulse)*0.1)*0.025)
		s.pulse++
	}
}

func (s *selector) getX(fudge float32) float32 {
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

func (s *selector) getY(fudge float32) float32 {
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
