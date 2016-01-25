package renderer

import (
	"strconv"

	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func renderMenu(g *game.Game, fudge float32) {
	ease := func(start, change float32) float32 {
		return easeOutCubic(g.StateProgress(fudge), start, change)
	}

	alpha := float32(1)
	switch g.State {
	case game.GameInitial, game.GamePaused:
		alpha = ease(0, 1)

	case game.GamePlaying, game.GameExiting:
		alpha = ease(1, -1)
	}

	// Don't render the menu if it is invisible.
	if alpha == 0 {
		return
	}

	gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &orthoProjectionViewMatrix[0])
	gl.Uniform1f(grayscaleUniform, 0)
	gl.Uniform1f(brightnessUniform, 0)
	gl.Uniform1f(alphaUniform, alpha)
	gl.Uniform1f(mixAmountUniform, 0)

	menu := g.Menu
	titleText := menuTitleText[menu.ID]
	totalHeight := titleText.height * 2
	for _, item := range menu.Items {
		totalHeight += float32(menuItemFontSize) * 2
		if len(item.Choices) > 0 {
			totalHeight += float32(menuItemFontSize) * 2
		}
	}

	currentY := (float32(winHeight) + totalHeight) / 2

	centerX := func(txt *renderableText) float32 {
		return (float32(winWidth) - txt.width) / 2
	}

	renderText := func(text *renderableText) {
		currentY -= text.height
		text.render(centerX(text), currentY)
		currentY -= text.height // add spacing for next item
	}

	// TODO(btmura): split these out into separate functions

	renderSlider := func(item *game.MenuItem) {
		val := strconv.Itoa(item.SliderValue)

		var valWidth, valHeight float32
		for _, rune := range val {
			text := menuRuneText[rune]
			valWidth += text.width
			if text.height > valHeight {
				valHeight = text.height
			}
		}

		currentY -= valHeight
		x := (float32(winWidth) - valWidth) / 2
		for _, rune := range val {
			text := menuRuneText[rune]
			text.render(x, currentY)
			x += text.width
		}
		currentY -= valHeight
	}

	renderMenuItem := func(index int, item *game.MenuItem) {
		var brightness float32
		if menu.FocusedIndex == index {
			switch {
			case menu.Selected:
				brightness = pulse(g.GlobalPulse+fudge, 1, 1, 1)

			case len(item.Choices) == 0:
				brightness = pulse(g.GlobalPulse+fudge, 1, 0.3, 0.06)

			default:
				brightness = 1
			}
		}
		gl.Uniform1f(brightnessUniform, brightness)
		renderText(menuItemText[item.ID])
		switch item.Type {
		case game.MenuChoice:
			if len(item.Choices) > 0 {
				selected := item.Choices[item.SelectedChoice]
				renderText(menuChoiceText[selected])
			}

		case game.MenuSlider:
			renderSlider(item)
		}
	}

	renderText(titleText)
	for i, item := range menu.Items {
		renderMenuItem(i, item)
	}
}
