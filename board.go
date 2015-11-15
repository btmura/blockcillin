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
	return &board{
		rings: []*ring{
			{
				[]*cell{
					{red},
					{purple},
					{red},
					{blue},
					{red},
					{cyan},
					{red},
					{green},
					{red},
					{yellow},
				},
			},
			{
				[]*cell{
					{red},
					{purple},
					{red},
					{blue},
					{red},
					{cyan},
					{red},
					{green},
					{red},
					{yellow},
				},
			},
		},
		ringCount: 2,
		cellCount: 10,
	}
}
