package main

import "math/rand"

type blockColor int32

const (
	red blockColor = iota
	purple
	blue
	cyan
	green
	yellow
)

type board struct {
	rings []*ring

	// ringCount is how many rings the board has.
	ringCount int

	// cellCount is how many cells are in each ring.
	cellCount int
}

type ring struct {
	cells []*cell
}

type cell struct {
	blockColor blockColor
}

func newBoard() *board {
	b := &board{
		ringCount: 5,
		cellCount: 15,
	}

	for i := 0; i < b.ringCount; i++ {
		r := &ring{}
		for j := 0; j < b.cellCount; j++ {
			c := &cell{
				blockColor: blockColor(rand.Intn(6)),
			}
			r.cells = append(r.cells, c)
		}
		b.rings = append(b.rings, r)
	}

	return b
}
