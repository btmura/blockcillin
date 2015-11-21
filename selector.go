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

// numMoveFrames is the number of frames a move animation takes.
const numMoveFrames = 10

type selector struct {
	// state is the state of the selector.
	state selectorState

	// x is the column of the selector multiplied by 10 to avoid losing precision.
	// Example: x = 10 is 1 column down. x = 15 is 1.5 columns down.
	x int

	// y is the row of the selector multiplied by 10 to avoid losing precision.
	// Example: y = 10 is 1 row down. y = 15 is 1.5 rows down.
	y int

	// moveFrame is the current frame in the move animation from 0 to numMoveFrames.
	moveFrame int

	// scale is the scale of the selector to make it pulse.
	scale float32

	// pulse is an increasing counter used to calculate the pulsing amount.
	pulse int
}

func (s *selector) moveUp() {
	s.state = movingUp
}

func (s *selector) moveDown() {
	s.state = movingDown
}

func (s *selector) moveLeft() {
	s.state = movingLeft
}

func (s *selector) moveRight() {
	s.state = movingRight
}

func (s *selector) update() {
	updateState := func() {
		if s.moveFrame++; s.moveFrame == numMoveFrames {
			s.state = static
			s.moveFrame = 0
		}
	}

	switch s.state {
	case movingUp:
		s.y -= 1
		updateState()

	case movingDown:
		s.y += 1
		updateState()

	case movingLeft:
		s.x -= 1
		updateState()

	case movingRight:
		s.x += 1
		updateState()
	}

	s.scale = float32(1.0 + math.Sin(float64(s.pulse)*0.1)*0.025)
	s.pulse++
}
