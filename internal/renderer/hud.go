package renderer

import (
	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func renderHUD(g *game.Game, fudge float32) {
	gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &orthoProjectionViewMatrix[0])

	i := 1
	renderText := func(item game.HUDItem) {
		text := hudItemText[item]
		x := float32(winWidth)/4*float32(i) - text.width/2
		y := float32(winHeight) - text.height*2
		text.render(x, y)
		i++
	}

	renderText(game.HUDItemSpeed)
	renderText(game.HUDItemTime)
	renderText(game.HUDItemScore)
}
