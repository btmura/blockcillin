package renderer

import (
	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func renderHUD(g *game.Game, fudge float32) {
	gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &orthoProjectionViewMatrix[0])
	gl.Uniform1f(grayscaleUniform, 0)
	gl.Uniform1f(brightnessUniform, 0)
	gl.Uniform1f(alphaUniform, 1)
	gl.Uniform1f(mixAmountUniform, 0)

	tx := (float32(winWidth) - speedText.width) / 2
	ty := float32(winHeight) - speedText.height*2
	speedText.render(tx, ty)
}
