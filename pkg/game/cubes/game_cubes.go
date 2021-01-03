package cubes

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func BuildCubeGame(windowWidth, windowHeight int) *ant.Game {
	window := ant.InitGlfw(windowWidth, windowHeight)
	ant.InitOpenGL()
	world := buildCubeWorld(windowWidth, windowHeight)

	cursor := new(Cursor)
	cam := new(Camera)
	setupInputHandling(window, &world, cursor, cam)

	game := ant.NewGame(window, &world)
	game.PreDraw = func() {
		view := cam.CalculateViewMatrix()
		world.Uniforms.SetMat4("ViewMatrix", view)
	}
	return game
}

type Cursor struct {
	xpos float64
	ypos float64
}

func setupInputHandling(window *glfw.Window, world *ant.GameWorld, cursor *Cursor, cam *Camera) {
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		dx := xpos - cursor.xpos
		dy := ypos - cursor.ypos

		dtheta := -dx * 0.005
		dphi := -dy * 0.005
		cam.Rotate(dtheta, dphi)

		cursor.xpos = xpos
		cursor.ypos = ypos
	})

	// todo: mouse clicks
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {})

	// todo: movement
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// world.Uniforms.SetMat4("Position", )
		if key == glfw.KeyA {
			// move forward
		}
		if key == glfw.KeyZ {
			// move backward
		}
		if key == glfw.KeyS {
			// move left
		}
		if key == glfw.KeyD {
			// move right
		}
		if key == glfw.KeySpace {
			// move up
		}
		if key == glfw.KeyLeftControl {
			// move down
		}
	})

	// todo: update projection matrix after resize
	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {})
}
