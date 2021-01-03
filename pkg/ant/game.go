package ant

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Game struct {
	Window *glfw.Window
	World  *GameWorld
}

func (game *Game) Run() {
	defer glfw.Terminate()
	for !game.Window.ShouldClose() {
		// todo: update on separate thread
		game.Update()
		game.Draw()
	}
}

func (game *Game) Update() {
	game.World.Update()
}

func (game *Game) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	game.World.Render()
	glfw.PollEvents()
	game.Window.SwapBuffers()
}
