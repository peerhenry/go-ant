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
	cam := NewCamera()
	commands := new(Commands)
	setupInputHandling(window, &world, cursor, cam, commands)

	game := ant.NewGame(window, &world)
	game.PreUpdate = func() {
		move(commands, cam)
	}
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

type Commands struct {
	forward  bool
	backward bool
	left     bool
	right    bool
	up       bool
	down     bool
}

func setupInputHandling(window *glfw.Window, world *ant.GameWorld, cursor *Cursor, cam *Camera, commands *Commands) {
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
			if action == glfw.Press {
				commands.forward = true
			}
			if action == glfw.Release {
				commands.forward = false
			}
		}
		if key == glfw.KeyZ {
			if action == glfw.Press {
				commands.backward = true
			}
			if action == glfw.Release {
				commands.backward = false
			}
		}
		if key == glfw.KeyS {
			if action == glfw.Press {
				commands.left = true
			}
			if action == glfw.Release {
				commands.left = false
			}
		}
		if key == glfw.KeyD {
			if action == glfw.Press {
				commands.right = true
			}
			if action == glfw.Release {
				commands.right = false
			}
		}
		if key == glfw.KeySpace {
			if action == glfw.Press {
				commands.up = true
			}
			if action == glfw.Release {
				commands.up = false
			}
		}
		if key == glfw.KeyLeftControl {
			if action == glfw.Press {
				commands.down = true
			}
			if action == glfw.Release {
				commands.down = false
			}
		}
	})

	// todo: update projection matrix after resize
	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {})
}

func move(commands *Commands, cam *Camera) {
	moveDir := Vec3{0, 0, 0}
	isMoving := false

	if commands.forward {
		if !commands.backward {
			moveDir = moveDir.Add(cam.direction)
			isMoving = true
		}
	} else if commands.backward {
		moveDir = moveDir.Sub(cam.direction)
		isMoving = true
	}

	if commands.right {
		if !commands.left {
			moveDir = moveDir.Add(cam.right)
			isMoving = true
		}
	} else if commands.left {
		moveDir = moveDir.Sub(cam.right)
		isMoving = true
	}

	if commands.up {
		if !commands.down {
			moveDir = moveDir.Add(Vec3{0, 0, 1})
			isMoving = true
		}
	} else if commands.down {
		moveDir = moveDir.Sub(Vec3{0, 0, 1})
		isMoving = true
	}

	if isMoving {
		moveDir = moveDir.Normalize()
		cam.position = cam.position.Add(moveDir.Mul(0.05)) // need delta time
	}
}
