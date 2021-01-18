package voxels

import (
	"time"

	"ant.com/ant/pkg/ant"
	"ant.com/ant/pkg/game/voxels/chunks"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
)

func BuildGame(windowWidth, windowHeight int) *ant.Game {
	window := ant.InitGlfw(windowWidth, windowHeight)
	ant.InitOpenGL()
	gl.ClearColor(100./256., 149./256., 237./256., 1.0) // todo: put this in a better place
	cam := ant.NewCamera()
	cam.Position = mgl64.Vec3{0, 0, 30}
	inputHandler := SetupInputHandling(window, cam)
	game := ant.NewGame(window)
	scene := BuildChunkScene(windowWidth, windowHeight)
	game.AddScene(scene)
	hud := BuildHud(windowWidth, windowHeight)
	chunkWorldUpdater := chunks.NewChunkWorldUpdater(cam, scene)
	game.Update = func(dt *time.Duration) {
		inputHandler.Update(dt)
		hud.Update(dt)
		chunkWorldUpdater.Update(dt)
	}
	game.PreDraw = func() {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.Enable(gl.BLEND)
		view := cam.CalculateViewMatrix()
		scene.UniformStore.SetMat4("ViewMatrix", view)
	}
	game.PostDraw = func() {
		hud.Draw()
	}
	return game
}
