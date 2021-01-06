package text

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func BuildGame(windowWidth, windowHeight int) *ant.Game {
	window := ant.InitGlfw(windowWidth, windowHeight)
	ant.InitOpenGL()
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	world := buildQuadWorld(windowWidth, windowHeight)

	return ant.NewGame(window, &world)
}
