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
	gl.ClearColor(0.2, 0.4, 0.9, 1.0)
	// gl.ClearColor(0.9, 0.9, 0.9, 1.0)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	world := buildQuadWorld(windowWidth, windowHeight)

	return ant.NewGame(window, &world)
}
