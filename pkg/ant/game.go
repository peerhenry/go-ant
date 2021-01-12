package ant

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Game struct {
	Window    *glfw.Window
	World     *Scene
	PreDraw   func()
	PostDraw  func()
	PreUpdate func(dt *time.Duration)
	then      time.Time
}

func NewGame(window *glfw.Window, world *Scene) *Game {
	return &Game{
		Window:    window,
		World:     world,
		PreDraw:   func() {},
		PostDraw:  func() {},
		PreUpdate: func(dt *time.Duration) {},
		then:      time.Now(),
	}
}

func (game *Game) Run() {
	defer glfw.Terminate()
	for !game.Window.ShouldClose() {
		now := time.Now()
		dt := now.Sub(game.then)
		game.then = now
		// todo: update on separate thread
		game.Update(&dt)
		game.Draw()
	}
}

func (game *Game) Update(dt *time.Duration) {
	game.PreUpdate(dt)
	game.World.Update(dt)
}

func (game *Game) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	game.PreDraw()
	game.World.Render()
	game.PostDraw()
	glfw.PollEvents()
	game.Window.SwapBuffers()
}
