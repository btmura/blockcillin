package game

import "github.com/btmura/blockcillin/internal/audio"

// NumMoveSteps is the number of steps a move animation takes.
const NumMoveSteps float32 = 0.1 / SecPerUpdate

type Selector struct {
	// State is the state of the selector.
	State selectorState

	// X is the selector's current column.
	// It changes only after the move animation is complete.
	X int

	// Y is the selector's current row.
	// It changes only after the move animation is complete.
	Y int

	// ringCount is how many rings the board has.
	ringCount int

	// cellCount is how many cells each ring has.
	cellCount int

	// Step is the current step in any animations.
	Step float32

	// pulse is used to advance any pulsing animations.
	Pulse float32
}

type selectorState int32

const (
	SelectorStatic selectorState = iota
	SelectorMovingUp
	SelectorMovingDown
	SelectorMovingLeft
	SelectorMovingRight
)

func newSelector(ringCount, cellCount int) *Selector {
	return &Selector{
		ringCount: ringCount,
		cellCount: cellCount,
	}
}

func (s *Selector) moveUp() {
	if s.State == SelectorStatic && s.Y > 0 {
		s.State = SelectorMovingUp
		audio.Play(audio.SoundMove)
	}
}

func (s *Selector) moveDown() {
	if s.State == SelectorStatic && s.Y < s.ringCount-1 {
		s.State = SelectorMovingDown
		audio.Play(audio.SoundMove)
	}
}

func (s *Selector) moveLeft() {
	if s.State == SelectorStatic {
		s.State = SelectorMovingLeft
		audio.Play(audio.SoundMove)
	}
}

func (s *Selector) moveRight() {
	if s.State == SelectorStatic {
		s.State = SelectorMovingRight
		audio.Play(audio.SoundMove)
	}
}

func (s *Selector) update() {
	updateMove := func() bool {
		if s.Step++; s.Step >= NumMoveSteps {
			s.State = SelectorStatic
			s.Step = 0
			return true
		}
		return false
	}

	switch s.State {
	case SelectorMovingUp:
		if updateMove() {
			s.Y--
		}

	case SelectorMovingDown:
		if updateMove() {
			s.Y++
		}

	case SelectorMovingLeft:
		if updateMove() {
			if s.X--; s.X < 0 {
				s.X = s.cellCount - 1
			}
		}

	case SelectorMovingRight:
		if updateMove() {
			s.X = (s.X + 1) % s.cellCount
		}

	default:
		s.Pulse++
	}
}

// nextPosition returns the next position of the selector for swapping.
// This can be different from the current position if the selector is moving.
func (s *Selector) nextPosition() (int, int) {
	switch s.State {
	case SelectorMovingUp:
		return s.X, s.Y - 1

	case SelectorMovingDown:
		return s.X, s.Y + 1

	case SelectorMovingLeft:
		X := s.X - 1
		if X < 0 {
			X = s.cellCount - 1
		}
		return X, s.Y

	case SelectorMovingRight:
		return (s.X + 1) % s.cellCount, s.Y
	}
	return s.X, s.Y
}
