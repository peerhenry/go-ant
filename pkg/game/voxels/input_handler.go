package voxels

import (
	"time"

	"ant.com/ant/pkg/game/voxels/chunks"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

type Cursor struct {
	xpos float64
	ypos float64
}

type Commands struct {
	forward      bool
	backward     bool
	left         bool
	right        bool
	up           bool
	down         bool
	fast         bool
	wireFrame    bool
	toggleNoclip bool
	click        bool
}

type InputHandler struct {
	commands      *Commands
	player        *chunks.Player
	drawWireFrame bool
}

func SetupInputHandling(window *glfw.Window, player *chunks.Player) *InputHandler {
	commands := new(Commands)
	cursor := new(Cursor)
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		dx := xpos - cursor.xpos
		dy := ypos - cursor.ypos

		dtheta := -dx * 0.005
		dphi := -dy * 0.005
		player.Camera.Rotate(dtheta, dphi)

		cursor.xpos = xpos
		cursor.ypos = ypos
	})

	// todo: mouse clicks
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				commands.click = true
			}
			if action == glfw.Release {
			}
		}
	})

	// todo: movement
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// ==== toggle wireframe ====
		if key == glfw.KeyF1 {
			if action == glfw.Press {
				commands.wireFrame = true
			}
			if action == glfw.Release {
				commands.wireFrame = false
			}
		}
		if key == glfw.KeyF2 {
			if action == glfw.Press {
				commands.toggleNoclip = true
			}
			if action == glfw.Release {
				commands.toggleNoclip = false
			}
		}
		// ==== Left shift go fast ====
		if key == glfw.KeyLeftShift {
			if action == glfw.Press {
				commands.fast = true
			}
			if action == glfw.Release {
				commands.fast = false
			}
		}
		// ==== Movement AZSD Space Ctrl ====
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

	return &InputHandler{
		commands:      commands,
		drawWireFrame: false,
		player:        player,
	}
}

func (self *InputHandler) Update(dt *time.Duration) {
	if self.commands.wireFrame {
		self.commands.wireFrame = false
		self.drawWireFrame = !self.drawWireFrame
		if self.drawWireFrame {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}
	}
	if self.commands.toggleNoclip {
		self.player.Noclip = !self.player.Noclip
		self.commands.toggleNoclip = false
	}
	if self.commands.click {
		chunks.RemoveBlock(self.player)
		self.commands.click = false
	}
	self.move(dt)
}

func (self *InputHandler) move(dt *time.Duration) {
	moveDir := mgl64.Vec3{0, 0, 0}
	isMoving := false
	speed := 5.0
	if self.player.Noclip {
		moveDir, isMoving = self.fly()
		speed = 30.0
		if self.commands.fast {
			speed = 70.0
		}
	} else {
		moveDir, isMoving = self.walk()
		if self.commands.fast {
			speed = 10.0
		}
	}

	dx := speed * dt.Seconds()
	if isMoving {
		moveDir = moveDir.Normalize()
		self.player.SuggestMovement(moveDir.Mul(dx))
	}
}

func (self *InputHandler) fly() (mgl64.Vec3, bool) {
	moveDir := mgl64.Vec3{0, 0, 0}
	isMoving := false

	if self.commands.forward {
		if !self.commands.backward {
			moveDir = moveDir.Add(self.player.Camera.Direction)
			isMoving = true
		}
	} else if self.commands.backward {
		moveDir = moveDir.Sub(self.player.Camera.Direction)
		isMoving = true
	}

	if self.commands.right {
		if !self.commands.left {
			moveDir = moveDir.Add(self.player.Camera.Right)
			isMoving = true
		}
	} else if self.commands.left {
		moveDir = moveDir.Sub(self.player.Camera.Right)
		isMoving = true
	}

	if self.commands.up {
		if !self.commands.down {
			moveDir = moveDir.Add(mgl64.Vec3{0, 0, 1})
			isMoving = true
		}
	} else if self.commands.down {
		moveDir = moveDir.Sub(mgl64.Vec3{0, 0, 1})
		isMoving = true
	}
	return moveDir, isMoving
}

func (self *InputHandler) walk() (mgl64.Vec3, bool) {
	moveDir := mgl64.Vec3{0, 0, 0}
	isMoving := false

	camDir := self.player.Camera.Direction

	if self.commands.forward {
		if !self.commands.backward {
			moveDir = moveDir.Add(mgl64.Vec3{camDir[0], camDir[1], 0})
			isMoving = true
		}
	} else if self.commands.backward {
		moveDir = moveDir.Sub(mgl64.Vec3{camDir[0], camDir[1], 0})
		isMoving = true
	}

	if self.commands.right {
		if !self.commands.left {
			moveDir = moveDir.Add(self.player.Camera.Right)
			isMoving = true
		}
	} else if self.commands.left {
		moveDir = moveDir.Sub(self.player.Camera.Right)
		isMoving = true
	}

	if self.commands.up {
		self.player.Jump()
	}

	return moveDir, isMoving
}

// todo
// func (self *InputHandler) swim() (mgl64.Vec3, bool) {

// }
