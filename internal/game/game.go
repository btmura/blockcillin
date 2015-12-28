package game

import (
	"github.com/btmura/blockcillin/internal/audio"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	updatesPerSec = 60
	SecPerUpdate  = 1.0 / updatesPerSec
)

type Game struct {
	State GameState
	Menu  *Menu
	Board *Board
	HUD   *HUD

	// GlobalPulse is incremented each update so it can be used for any pulsing animation.
	GlobalPulse float32

	nextBoard *Board
	nextHUD   *HUD
	step      float32
}

type GameState int32

const (
	GameInitial GameState = iota
	GamePlaying
	GamePaused
	GameExiting
)

var gameStateSteps = map[GameState]float32{
	GameInitial: 0.5 / SecPerUpdate,
	GamePlaying: 0.5 / SecPerUpdate,
	GamePaused:  0.5 / SecPerUpdate,
	GameExiting: 0.5 / SecPerUpdate,
}

func New() *Game {
	return &Game{
		Menu: newMenu(),
	}
}

func (g *Game) KeyCallback(key glfw.Key, action glfw.Action) {
	if action != glfw.Press && action != glfw.Repeat {
		return
	}

	if g.StateProgress(0) < 1 {
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
			g.setState(GamePaused)
			g.Menu.pause()
			audio.Play(audio.SoundSelect)
		}

	case GameInitial, GamePaused:
		switch key {
		case glfw.KeyDown:
			g.Menu.moveDown()
			audio.Play(audio.SoundMove)

		case glfw.KeyUp:
			g.Menu.moveUp()
			audio.Play(audio.SoundMove)

		case glfw.KeyEnter, glfw.KeySpace:
			switch g.Menu.focused() {
			case MenuItemContinueGame:
				g.setState(GamePlaying)
				g.Menu.Selected = true
				audio.Play(audio.SoundSelect)

			case MenuItemNewGame:
				g.setState(GamePlaying)

				b := newBoard(10, 15, 2, 2)
				h := newHUD()
				if g.Board == nil {
					g.Board = b
					g.HUD = h
				} else {
					g.nextBoard = b
					g.nextHUD = h
					g.Board.exit()
				}

				g.Menu.Selected = true
				audio.Play(audio.SoundSelect)

			case MenuItemExit:
				g.setState(GameExiting)
				g.Menu.Selected = true
				audio.Play(audio.SoundSelect)
			}

		case glfw.KeyEscape:
			switch g.State {
			case GamePaused:
				g.setState(GamePlaying)
				audio.Play(audio.SoundSelect)
			}
		}
	}
}

func (g *Game) Update() {
	g.GlobalPulse++

	switch g.State {
	case GameInitial, GamePaused, GameExiting:
		g.step++

	case GamePlaying:
		g.step++
		g.Board.update()
		g.HUD.update()

		g.HUD.Speed = 1 + g.Board.totalBlocksCleared/10
		g.HUD.Score += g.Board.newBlocksCleared * 10

		switch g.Board.State {
		case BoardGameOver:
			g.Menu.gameOver()
			g.setState(GameInitial)

		default:
			if g.Board.done() {
				g.Board = g.nextBoard
				g.HUD = g.nextHUD
				g.nextBoard = nil
				g.nextHUD = nil
			}
		}
	}
}

func (g *Game) StateProgress(fudge float32) float32 {
	totalSteps := gameStateSteps[g.State]
	if totalSteps == 0 {
		return 1
	}

	if p := (g.step + fudge) / totalSteps; p < 1 {
		return p
	}
	return 1
}

func (g *Game) setState(state GameState) {
	g.State = state
	g.step = 0
}

func (g *Game) Done() bool {
	return g.State == GameExiting && g.StateProgress(0) >= 1
}
