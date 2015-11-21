package main

import "math"

type selector struct {
	y float32

	// scale is the scale of the selector to make it pulse.
	scale float32

	// pulse is an increasing counter used to calculate the pulsing amount.
	pulse int
}

func (s *selector) moveUp() {
	s.y--
}

func (s *selector) moveDown() {
	s.y++
}

func (s *selector) update() {
	s.scale = float32(1.0 + math.Sin(float64(s.pulse)*0.1)*0.025)
	s.pulse++
}
