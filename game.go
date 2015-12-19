package main

import "github.com/go-gl/glfw/v3.1/glfw"

type game struct {
	state       gameState
	board       *board
	keyCallback func(key glfw.Key, action glfw.Action) bool
}

type gameState int32

const (
	gameInitial gameState = iota
)

func newGame() *game {
	b := newBoard(&boardConfig{
		ringCount:       10,
		cellCount:       15,
		filledRingCount: 2,
		spareRingCount:  2,
	})

	keyCallback := func(key glfw.Key, action glfw.Action) bool {
		if action != glfw.Press && action != glfw.Repeat {
			return false
		}

		switch key {
		case glfw.KeyLeft:
			b.moveLeft()
			return true

		case glfw.KeyRight:
			b.moveRight()
			return true

		case glfw.KeyDown:
			b.moveDown()
			return true

		case glfw.KeyUp:
			b.moveUp()
			return true

		case glfw.KeySpace:
			b.swap()
			return true

		default:
			return false
		}
	}

	return &game{
		board:       b,
		keyCallback: keyCallback,
	}
}

func (g *game) update() {
	g.board.update()
}
