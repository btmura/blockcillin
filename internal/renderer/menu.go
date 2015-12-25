package renderer

import (
	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func renderMenu(g *game.Game, fudge float32) {
	alpha := float32(1)
	ease := func(start, change float32) float32 {
		return easeOutCubic2(g.StateProgress(fudge), start, change)
	}

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
	gl.Uniform1f(alphaUniform, alpha)
	gl.Uniform1f(mixAmountUniform, 0)

	menu := g.Menu

	totalHeight := titleText.height*2 + float32(menuItemFontSize*len(menu.Items)*2)
	ty := (float32(winHeight) + totalHeight) / 2

	renderMenuItem := func(text *rendererText, focused bool) {
		tx := (float32(winWidth) - text.width) / 2
		ty -= text.height

		m := newScaleMatrix(text.width, text.height, 1)
		m = m.mult(newTranslationMatrix(tx, ty, 0))
		gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
		gl.Uniform1i(textureUniform, int32(text.texture)-1)

		var brightness float32
		switch {
		case focused && menu.Selected:
			brightness = pulse(menu.Pulse+fudge, 1, 1, 1)

		case focused:
			brightness = pulse(menu.Pulse+fudge, 1, 0.3, 0.06)
		}
		gl.Uniform1f(brightnessUniform, brightness)
		textLineMesh.drawElements()

		ty -= text.height
	}

	renderMenuItem(titleText, false)
	for i, item := range menu.Items {
		renderMenuItem(menuItemText[item], menu.FocusedIndex == i)
	}
}
