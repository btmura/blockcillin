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

	nextMenu  *Menu
	nextBoard *Board
	nextHUD   *HUD
	step      float32
}

//go:generate stringer -type=GameState
type GameState int32

const (
	GameInitial GameState = iota
	GamePlaying
	GamePaused
	GameExiting
)

var gameStateSteps = [...]float32{
	GameInitial: 0.5 / SecPerUpdate,
	GamePlaying: 0.5 / SecPerUpdate,
	GamePaused:  0.5 / SecPerUpdate,
	GameExiting: 0.5 / SecPerUpdate,
}

func New() *Game {
	return &Game{
		Menu: mainMenu,
	}
}

func (g *Game) KeyCallback(key glfw.Key, action glfw.Action) {
	if action != glfw.Press && action != glfw.Repeat {
		// Handle any release triggers.
		if key == glfw.KeyLeftAlt {
			g.Board.useManualRiseRate = false
		}
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

		case glfw.KeyLeftAlt:
			g.Board.useManualRiseRate = true

		case glfw.KeyEscape:
			g.setState(GamePaused)
			g.Menu = pausedMenu
			g.Menu.reset()
			audio.Play(audio.SoundSelect)
		}

	case GameInitial, GamePaused:
		switch key {
		case glfw.KeyLeft:
			g.Menu.moveLeft()

		case glfw.KeyRight:
			g.Menu.moveRight()

		case glfw.KeyDown:
			g.Menu.moveDown()

		case glfw.KeyUp:
			g.Menu.moveUp()

		case glfw.KeyEnter, glfw.KeySpace:
			switch g.Menu.focused() {
			case MenuNewGameItem:
				g.Menu.selectItem()
				g.Menu = newGameMenu
				g.Menu.reset()

			case MenuExit:
				g.Menu.selectItem()
				g.setState(GameExiting)

			case MenuOK:
				g.Menu.selectItem()
				g.setState(GamePlaying)

				var numBlockColors int
				switch difficultyItem.Selector.Value() {
				case MenuEasy:
					numBlockColors = maxBlockColors - 2

				case MenuMedium:
					numBlockColors = maxBlockColors - 1

				default:
					numBlockColors = maxBlockColors
				}

				speed := speedItem.Slider.Value

				b := newBoard(numBlockColors, speed)
				h := newHUD(speed)
				if g.Board == nil {
					g.Board = b
					g.HUD = h
				} else {
					g.nextBoard = b
					g.nextHUD = h
					g.Board.exit()
				}

			case MenuContinueGame:
				g.Menu.selectItem()
				g.setState(GamePlaying)

			case MenuQuit:
				g.Menu.selectItem()
				g.Menu = mainMenu
				g.Menu.reset()
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

		switch g.Board.State {
		case BoardLive:
			g.HUD.Speed = g.Board.speed
			g.HUD.Score += g.Board.numUpdateBlocksCleared * 10
			g.HUD.update()

		case BoardGameOver:
			if g.Board.StateDone() {
				g.Menu = gameOverMenu
				g.Menu.reset()
				g.setState(GameInitial)
			}

		case BoardExiting:
			if g.Board.StateDone() {
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
