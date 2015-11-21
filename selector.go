package main

import "math"

type selectorState int32

const (
	static selectorState = iota
	movingUp
	movingDown
	movingLeft
	movingRight
)

// numMoveSteps is the number of frames a move animation takes.
const numMoveSteps = 10

type selector struct {
	// state is the state of the selector.
	state selectorState

	// x is the selector's current column.
	// It changes only after the move animation is complete.
	x int

	// y is the selector's current row.
	// It changes only after the move animation is complete.
	y int

	// moveStep is the current frame in the move animation from 0 to numMoveSteps.
	moveStep int

	// scale is the scale of the selector to make it pulse.
	scale float32

	// pulse is an increasing counter used to calculate the pulsing amount.
	pulse int
}

func (s *selector) moveUp() {
	if s.state == static {
		s.state = movingUp
	}
}

func (s *selector) moveDown() {
	if s.state == static {
		s.state = movingDown
	}
}

func (s *selector) moveLeft() {
	if s.state == static {
		s.state = movingLeft
	}
}

func (s *selector) moveRight() {
	if s.state == static {
		s.state = movingRight
	}
}

func (s *selector) update() {
	updateMove := func() bool {
		if s.moveStep++; s.moveStep == numMoveSteps {
			s.state = static
			s.moveStep = 0
			return true
		}
		return false
	}

	switch s.state {
	case movingUp:
		if updateMove() {
			s.y--
		}

	case movingDown:
		if updateMove() {
			s.y++
		}

	case movingLeft:
		if updateMove() {
			s.x--
		}

	case movingRight:
		if updateMove() {
			s.x++
		}

	default:
		s.scale = float32(1.0 + math.Sin(float64(s.pulse)*0.1)*0.025)
		s.pulse++
	}
}

func (s *selector) getX(fudge float32) float32 {
	sx := float32(s.x)
	dx := (float32(s.moveStep) + fudge) / numMoveSteps

	switch s.state {
	case movingLeft:
		return sx - dx

	case movingRight:
		return sx + dx

	default:
		return sx
	}
}

func (s *selector) getY(fudge float32) float32 {
	sy := float32(s.y)
	dy := (float32(s.moveStep) + fudge) / numMoveSteps

	switch s.state {
	case movingUp:
		return sy - dy

	case movingDown:
		return sy + dy

	default:
		return sy
	}
}
