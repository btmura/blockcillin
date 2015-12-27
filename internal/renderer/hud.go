package renderer

import (
	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func renderHUD(g *game.Game, fudge float32) {
	gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &orthoProjectionViewMatrix[0])

	i := 1
	renderText := func(item game.HUDItem, val string) {
		text := hudItemText[item]
		x := float32(winWidth)/4*float32(i) - text.width/2
		y := float32(winHeight) - text.height*2
		text.render(x, y)

		var valWidth, valHeight float32
		for _, rune := range val {
			text := hudRuneText[rune]
			valWidth += text.width
			if text.height > valHeight {
				valHeight = text.height
			}
		}

		x = float32(winWidth)/4*float32(i) - valWidth/2
		y -= valHeight
		for _, rune := range val {
			text := hudRuneText[rune]
			text.render(x, y)
			x += text.width
		}

		i++
	}

	renderText(game.HUDItemSpeed, "1")
	renderText(game.HUDItemTime, "3:37")
	renderText(game.HUDItemScore, "1337")
}
