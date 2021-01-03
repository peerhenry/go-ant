package game

import (
	"ant.com/ant/pkg/ant"
)

func BuildQuadGame(windowWidth, windowHeight int) *ant.Game {
	window := ant.InitGlfw(windowWidth, windowHeight)
	ant.InitOpenGL()
	world := buildQuadWorld()

	return &ant.Game{
		Window: window,
		World:  &world,
	}
}
