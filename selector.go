package main

import "github.com/btmura/blockcillin/internal/audio"

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

	// ringCount is how many rings the board has.
	ringCount int

	// cellCount is how many cells each ring has.
	cellCount int

	// step is the current step in any animations.
	step float32

	// pulse is used to advance any pulsing animations.
	pulse float32
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
		audio.Play(audio.SoundMove)
	}
}

func (s *selector) moveDown() {
	if s.state == selectorStatic && s.y < s.ringCount-1 {
		s.state = selectorMovingDown
		audio.Play(audio.SoundMove)
	}
}

func (s *selector) moveLeft() {
	if s.state == selectorStatic {
		s.state = selectorMovingLeft
		audio.Play(audio.SoundMove)
	}
}

func (s *selector) moveRight() {
	if s.state == selectorStatic {
		s.state = selectorMovingRight
		audio.Play(audio.SoundMove)
	}
}

func (s *selector) update() {
	updateMove := func() bool {
		if s.step++; s.step >= numMoveSteps {
			s.state = selectorStatic
			s.step = 0
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
			s.x = (s.x + 1) % s.cellCount
		}

	default:
		s.pulse++
	}
}

// nextPosition returns the next position of the selector for swapping.
// This can be different from the current position if the selector is moving.
func (s *selector) nextPosition() (int, int) {
	switch s.state {
	case selectorMovingUp:
		return s.x, s.y - 1

	case selectorMovingDown:
		return s.x, s.y + 1

	case selectorMovingLeft:
		x := s.x - 1
		if x < 0 {
			x = s.cellCount - 1
		}
		return x, s.y

	case selectorMovingRight:
		return (s.x + 1) % s.cellCount, s.y
	}
	return s.x, s.y
}
