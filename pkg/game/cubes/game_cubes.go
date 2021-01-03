package cubes

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func BuildCubeGame(windowWidth, windowHeight int) *ant.Game {
	window := ant.InitGlfw(windowWidth, windowHeight)
	ant.InitOpenGL()

	// return BuildQuadGame()
	world := buildCubeWorld(windowWidth, windowHeight)

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		pos := mgl32.Vec3{5, 3, 3}
		unit := mgl32.Vec3{-5, -3, -3}.Normalize()
		rot := mgl32.Rotate3DZ(float32(-xpos) / 200)
		target := pos.Add(rot.Mul3x1(unit))
		view := mgl32.LookAtV(
			mgl32.Vec3{5, 3, 3}, // eye
			target,              // center
			mgl32.Vec3{0, 0, 1}, // up
		)
		world.Uniforms.SetMat4("ViewMatrix", view)
	})

	return &ant.Game{
		Window: window,
		World:  &world,
	}
}
