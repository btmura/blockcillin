package main

import (
	"github.com/btmura/blockcillin/internal/audio"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type game struct {
	state gameState
	menu  *menu
	board *board
}

type gameState int32

const (
	gameNeverStarted gameState = iota
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
			audio.Play(audio.SoundSelect)
		}

	default:
		switch key {
		case glfw.KeyDown:
			g.menu.moveDown()
			audio.Play(audio.SoundMove)

		case glfw.KeyUp:
			g.menu.moveUp()
			audio.Play(audio.SoundMove)

		case glfw.KeyEnter, glfw.KeySpace:
			switch g.menu.selectedItem() {
			case menuContinueGame:
				g.state = gamePlaying
				audio.Play(audio.SoundSelect)

			case menuNewGame:
				g.state = gamePlaying
				g.board = newBoard(&boardConfig{
					ringCount:       10,
					cellCount:       15,
					filledRingCount: 2,
					spareRingCount:  2,
				})
				audio.Play(audio.SoundSelect)

			case menuExit:
				audio.Play(audio.SoundSelect)
				win.SetShouldClose(true)
			}

		case glfw.KeyEscape:
			switch g.state {
			case gamePaused:
				g.state = gamePlaying
				audio.Play(audio.SoundSelect)
			}
		}
	}
}

func (g *game) update() {
	switch g.state {
	case gamePlaying:
		g.board.update()
		if g.board.state == boardGameOver {
			g.state = gameNeverStarted
			g.menu.removeContinueGame()
		}
	}
}
