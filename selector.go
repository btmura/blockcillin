package main

type selector struct {
	y float32
}

func (s *selector) moveUp() {
	s.y--
}

func (s *selector) moveDown() {
	s.y++
}
