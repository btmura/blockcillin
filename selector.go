package main

import "math"

type selectorState int32

const (
	static selectorState = iota
	movingUp
	movingDown
)

// numMoveFrames is the number of frames a move animation takes.
const numMoveFrames = 10

type selector struct {
	// state is the state of the selector.
	state selectorState

	y float32

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

func (s *selector) update() {
	updateState := func() {
		if s.moveFrame++; s.moveFrame == numMoveFrames {
			s.state = static
			s.moveFrame = 0
		}
	}

	switch s.state {
	case movingUp:
		s.y -= 0.1
		updateState()

	case movingDown:
		s.y += 0.1
		updateState()
	}

	s.scale = float32(1.0 + math.Sin(float64(s.pulse)*0.1)*0.025)
	s.pulse++
}
