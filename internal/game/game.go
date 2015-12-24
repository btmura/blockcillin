package game

import (
	"github.com/btmura/blockcillin/internal/audio"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const SecPerUpdate = 1.0 / 60.0

type Game struct {
	State gameState
	Menu  *Menu
	Board *Board
}

type gameState int32

const (
	GameNeverStarted gameState = iota
	GamePlaying
	GamePaused
)

func New() *Game {
	return &Game{
		Menu: newMenu(),
	}
}

func (g *Game) KeyCallback(win *glfw.Window, key glfw.Key, action glfw.Action) {
	if action != glfw.Press && action != glfw.Repeat {
		return
	}

	switch g.State {
	case GamePlaying:
		switch key {
		case glfw.KeyLeft:
			g.Board.moveLeft()

		case glfw.KeyRight:
			g.Board.moveRight()

		case glfw.KeyDown:
			g.Board.moveDown()

		case glfw.KeyUp:
			g.Board.moveUp()

		case glfw.KeySpace:
			g.Board.swap()

		case glfw.KeyEscape:
			g.State = GamePaused
			g.Menu.addContinueGame()
			audio.Play(audio.SoundSelect)
		}

	default:
		switch key {
		case glfw.KeyDown:
			g.Menu.moveDown()
			audio.Play(audio.SoundMove)

		case glfw.KeyUp:
			g.Menu.moveUp()
			audio.Play(audio.SoundMove)

		case glfw.KeyEnter, glfw.KeySpace:
			switch g.Menu.selectedItem() {
			case menuContinueGame:
				g.State = GamePlaying
				audio.Play(audio.SoundSelect)

			case menuNewGame:
				g.State = GamePlaying
				g.Board = newBoard(&boardConfig{
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
			switch g.State {
			case GamePaused:
				g.State = GamePlaying
				audio.Play(audio.SoundSelect)
			}
		}
	}
}

func (g *Game) Update() {
	switch g.State {
	case GamePlaying:
		g.Board.update()
		if g.Board.state == boardGameOver {
			g.State = GameNeverStarted
			g.Menu.removeContinueGame()
		}
	}
}
