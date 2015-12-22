package main

import "github.com/go-gl/glfw/v3.1/glfw"

type game struct {
	state gameState
	menu  *menu
	board *board
}

type gameState int32

const (
	gameInit gameState = iota
	gameNotStarted
	gamePlaying
	gamePaused
)

func newGame() *game {
	return &game{
		menu: newMenu(),
	}
}

func (g *game) keyCallback(win *glfw.Window, key glfw.Key, action glfw.Action) {
	if action != glfw.Press && action != glfw.Repeat {
		return
	}

	switch g.state {
	case gamePlaying:
		switch key {
		case glfw.KeyLeft:
			g.board.moveLeft()

		case glfw.KeyRight:
			g.board.moveRight()

		case glfw.KeyDown:
			g.board.moveDown()

		case glfw.KeyUp:
			g.board.moveUp()

		case glfw.KeySpace:
			g.board.swap()

		case glfw.KeyEscape:
			g.state = gamePaused
			g.menu.addContinueGame()
			playSound(soundSelect)
		}

	default:
		switch key {
		case glfw.KeyDown:
			g.menu.moveDown()
			playSound(soundMove)

		case glfw.KeyUp:
			g.menu.moveUp()
			playSound(soundMove)

		case glfw.KeyEnter, glfw.KeySpace:
			switch g.menu.selectedItem() {
			case menuContinueGame:
				g.state = gamePlaying
				playSound(soundSelect)

			case menuNewGame:
				g.state = gamePlaying
				g.board = newBoard(&boardConfig{
					ringCount:       10,
					cellCount:       15,
					filledRingCount: 2,
					spareRingCount:  2,
				})
				playSound(soundSelect)

			case menuExit:
				playSound(soundSelect)
				win.SetShouldClose(true)
			}

		case glfw.KeyEscape:
			switch g.state {
			case gamePaused:
				g.state = gamePlaying
				playSound(soundSelect)
			}
		}
	}
}

func (g *game) update() {
	switch g.state {
	case gameInit:
		playSound(soundSelect)
		g.state = gameNotStarted

	case gamePlaying:
		g.board.update()
		if g.board.state == boardGameOver {
			g.state = gameNotStarted
			g.menu.removeContinueGame()
		}
	}
}
