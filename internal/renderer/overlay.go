package renderer

import (
	"github.com/btmura/blockcillin/internal/game"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func renderOverlay(g *game.Game, fudge float32) {
	gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &orthoProjectionViewMatrix[0])
	gl.Uniform1f(grayscaleUniform, 0)
	gl.Uniform1f(brightnessUniform, 0)
	gl.Uniform1f(alphaUniform, 1)
	gl.Uniform1f(mixAmountUniform, 0)

	tx := (float32(winWidth) - speedText.width) / 2
	ty := float32(winHeight) - speedText.height*2

	m := newScaleMatrix(speedText.width, speedText.height, 1)
	m = m.mult(newTranslationMatrix(tx, ty, 0))
	gl.UniformMatrix4fv(modelMatrixUniform, 1, false, &m[0])
	gl.Uniform1i(textureUniform, int32(speedText.texture)-1)

	textLineMesh.drawElements()
}
