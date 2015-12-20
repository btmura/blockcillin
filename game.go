package main

import "github.com/go-gl/glfw/v3.1/glfw"

type game struct {
	state       gameState
	menu        *menu
	board       *board
	keyCallback func(win *glfw.Window, key glfw.Key, action glfw.Action)
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

	g.keyCallback = func(win *glfw.Window, key glfw.Key, action glfw.Action) {
		if action != glfw.Press && action != glfw.Repeat {
			return
		}

		switch g.state {
		case gamePlaying:
			switch key {
			case glfw.KeyLeft:
				b.moveLeft()

			case glfw.KeyRight:
				b.moveRight()

			case glfw.KeyDown:
				b.moveDown()

			case glfw.KeyUp:
				b.moveUp()

			case glfw.KeySpace:
				b.swap()
			}

		default:
			switch key {
			case glfw.KeyDown:
				m.moveDown()

			case glfw.KeyUp:
				m.moveUp()

			case glfw.KeyEnter, glfw.KeySpace:
				switch m.items[m.selectedIndex] {
				case menuNewGame:
					g.state = gamePlaying

				case menuExit:
					win.SetShouldClose(true)
				}
			}
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
