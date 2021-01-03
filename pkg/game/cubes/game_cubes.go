package cubes

import (
	"math"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func BuildCubeGame(windowWidth, windowHeight int) *ant.Game {
	window := ant.InitGlfw(windowWidth, windowHeight)
	ant.InitOpenGL()
	world := buildCubeWorld(windowWidth, windowHeight)
	setupInputHandling(window, &world)

	game := ant.NewGame(window, &world)
	return game
}

type Cursor struct {
	xpos float64
	ypos float64
}

type Camera struct {
	phi      float64
	theta    float64
	position mgl32.Vec3
}

const PHI_MAX = 0.5*math.Pi - 0.001

func (self *Camera) rotate(dtheta float64, dphi float64) {
	self.theta = self.theta + dtheta
	self.phi = self.phi + dphi
	for self.theta > 2*math.Pi {
		self.theta -= 2 * math.Pi
	}
	for self.theta < 0 {
		self.theta += 2 * math.Pi
	}
	for self.phi > PHI_MAX {
		self.phi = PHI_MAX
	}
	for self.phi < -PHI_MAX {
		self.phi = -PHI_MAX
	}
}

func setupInputHandling(window *glfw.Window, world *ant.GameWorld) {
	cursor := new(Cursor)
	cam := new(Camera)

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		dx := xpos - cursor.xpos
		dy := ypos - cursor.ypos

		dtheta := -dx * 0.005
		dphi := -dy * 0.005
		cam.rotate(dtheta, dphi)

		cursor.xpos = xpos
		cursor.ypos = ypos
		// pos := mgl32.Vec3{5, 3, 3}
		// unit := mgl32.Vec3{-5, -3, -3}.Normalize()
		// float32(-ypos) / 200
		// rotX := mgl32.Rotate3DX()
		// rotZ := mgl32.Rotate3DZ(float32(-xpos) / 200)
		// rot :=
		// dir :=
		// target := pos.Add(dir)
		// view := mgl32.LookAtV(
		// 	mgl32.Vec3{5, 3, 3}, // eye
		// 	target,              // center
		// 	mgl32.Vec3{0, 0, 1}, // up
		// )
		// world.Uniforms.SetMat4("ViewMatrix", view)
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
