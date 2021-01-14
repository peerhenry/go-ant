package ant

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Game struct {
	Window    *glfw.Window
	Scenes    []*Scene
	Init      func()
	PreDraw   func()
	PostDraw  func()
	PreUpdate func(dt *time.Duration)
	Update    func(dt *time.Duration)
	then      time.Time
}

func NewGame(window *glfw.Window) *Game {
	return &Game{
		Window:    window,
		Init:      func() {},
		PreDraw:   func() {},
		PostDraw:  func() {},
		PreUpdate: func(dt *time.Duration) {},
		Update:    func(dt *time.Duration) {},
		then:      time.Now(),
	}
}

func (game *Game) AddScene(scene *Scene) int {
	index := len(game.Scenes)
	game.Scenes = append(game.Scenes, scene)
	return index
}

func (game *Game) Run() {
	defer glfw.Terminate()
	game.Init()
	for !game.Window.ShouldClose() {
		now := time.Now()
		dt := now.Sub(game.then)
		game.then = now
		// todo: update on separate thread
		game.PreUpdate(&dt)
		game.Update(&dt)
		game.Draw()
	}
}

func (self *Game) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	self.PreDraw()
	for _, scene := range self.Scenes {
		scene.Draw()
	}
	self.PostDraw()
	glfw.PollEvents()
	self.Window.SwapBuffers()
}
