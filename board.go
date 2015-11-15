package main

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
	blockColors []blockColor
}

func newBoard() *board {
	return &board{
		blockColors: []blockColor{red, purple, red, blue, red, cyan, red, green, red, yellow},
	}
}
