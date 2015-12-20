package main

import "github.com/go-gl/glfw/v3.1/glfw"

type game struct {
	state       gameState
	menu        *menu
	board       *board
	keyCallback func(key glfw.Key, action glfw.Action) bool
}

type gameState int32

const (
	gameInitial gameState = iota
	gamePlaying
)

func newGame() *game {
	m := newMenu()
	b := newBoard(&boardConfig{
		ringCount:       10,
		cellCount:       15,
		filledRingCount: 2,
		spareRingCount:  2,
	})
	g := &game{
		menu:  m,
		board: b,
	}

	g.keyCallback = func(key glfw.Key, action glfw.Action) bool {
		if action != glfw.Press && action != glfw.Repeat {
			return false
		}

		if g.state == gamePlaying {
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

		switch key {
		case glfw.KeyDown:
			m.moveDown()
			return true

		case glfw.KeyUp:
			m.moveUp()
			return true

		case glfw.KeyEnter, glfw.KeySpace:
			switch m.items[m.selectedIndex] {
			case menuNewGame:
				g.state = gamePlaying
				break
			}
			return true

		default:
			return false
		}
	}
	return g
}

func (g *game) update() {
	switch g.state {
	case gamePlaying:
		g.board.update()
	}
}
