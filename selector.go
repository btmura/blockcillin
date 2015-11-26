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
			if s.x++; s.x == s.cellCount {
				s.x = 0
			}
		}

	default:
		s.pulse++
	}
}

func (s *selector) renderX(fudge float32) float32 {
	move := func(delta float32) float32 {
		return linear(s.step+fudge, float32(s.x), delta, numMoveSteps)
	}

	switch s.state {
	case selectorMovingLeft:
		return move(-1)

	case selectorMovingRight:
		return move(1)
	}

	return float32(s.x)
}

func (s *selector) renderY(fudge float32) float32 {
	move := func(delta float32) float32 {
		return linear(s.step+fudge, float32(s.y), delta, numMoveSteps)
	}

	switch s.state {
	case selectorMovingUp:
		return move(-1)

	case selectorMovingDown:
		return move(1)
	}

	return float32(s.y)
}

func (s *selector) renderScale(fudge float32) float32 {
	return pulse(s.pulse+fudge, 1.0, 0.025, 0.1)
}
