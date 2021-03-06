package cubes

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func buildCubeScene(windowWidth, windowHeight int) *ant.Scene {
	scene := ant.CreateScene(windowWidth, windowHeight, "shaders/ads/vertex.glsl", "shaders/ads/fragment.glsl")
	setupUniforms(windowWidth, windowHeight, scene)
	createGameObjects(scene)
	scene.Render = renderCube

	return scene
}

func setupUniforms(windowWidth, windowHeight int, scene *ant.Scene) {
	scene.UniformStore.SetMat4("ViewMatrix", mgl32.LookAtV(
		mgl32.Vec3{1, 0, 0}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	))
	scene.UniformStore.SetMat4("ProjectionMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 100.0))
	ant.LoadImageFileToUniform("resources/atlas.png", "Tex", scene.GlslProgram.Handle, 0)
}

func createGameObjects(scene *ant.Scene) {
	cube1 := createCube(mgl32.Vec3{2, 0, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}), GRASS)
	scene.AddRenderData(cube1)
	cube2 := createCube(mgl32.Vec3{0, 2, 0}, mgl32.QuatRotate(1, mgl32.Vec3{1, 0, 0}), STONE)
	scene.AddRenderData(cube2)
	cube3 := createCube(mgl32.Vec3{0, 0, 2}, mgl32.QuatRotate(1, mgl32.Vec3{0, 1, 0}), DIRT)
	scene.AddRenderData(cube3)
	cube4 := createCube(mgl32.Vec3{-2, 0, 0}, mgl32.QuatRotate(2, mgl32.Vec3{0, 1, 0}), SAND)
	scene.AddRenderData(cube4)
	// return []*ant.RenderData{cube1, cube2, cube3, cube4}
}

func renderCube(uniformStore *ant.UniformStore, data *ant.RenderData) {
	// get values for uniform calculations
	viewMatrix := uniformStore.GetMat4("ViewMatrix")
	projectionMatrix := uniformStore.GetMat4("ProjectionMatrix")
	// calculate uniforms
	modelMatrix := data.Transform
	modelView := viewMatrix.Mul4(modelMatrix)
	normalMatrix := modelView.Mat3()
	mvp := projectionMatrix.Mul4(modelView)
	// set uiniforms
	uniformStore.UniformMat4("ModelViewMatrix", modelView)
	uniformStore.UniformMat3("NormalMatrix", normalMatrix)
	uniformStore.UniformMat4("MVP", mvp)
	// draw
	gl.BindVertexArray(data.Vao)
	gl.DrawElements(gl.TRIANGLES, data.IndicesCount, gl.UNSIGNED_INT, nil)
}
