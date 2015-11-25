package main

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

	// board is the board the selector is moving on.
	board *board

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

func newSelector(board *board) *selector {
	return &selector{board: board}
}

func (s *selector) moveUp() {
	if s.state == selectorStatic && s.y > 0 {
		s.state = selectorMovingUp
	}
}

func (s *selector) moveDown() {
	if s.state == selectorStatic && s.y < s.board.ringCount-1 {
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
				s.x = s.board.cellCount - 1
			}
		}

	case selectorMovingRight:
		if updateMove() {
			if s.x++; s.x == s.board.cellCount {
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
		return linear(s.step+fudge, sx, delta, numMoveSteps)
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
		return linear(s.step+fudge, sy, delta, numMoveSteps)
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
	return pulse(s.pulse+fudge, 1.0, 0.025, 0.1)
}
