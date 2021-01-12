package quad

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func BuildGame(windowWidth, windowHeight int, filePath string) *ant.Game {
	window := ant.InitGlfw(windowWidth, windowHeight)
	ant.InitOpenGL()
	gl.ClearColor(0.2, 0.4, 0.9, 1.0)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	world := buildQuadScene(filePath)

	return ant.NewGame(window, &world)
}
