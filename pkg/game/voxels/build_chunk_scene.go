package voxels

import (
	"ant.com/ant/pkg/ant"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func BuildChunkScene(windowWidth, windowHeight int) *ant.Scene {
	scene := ant.CreateScene(windowWidth, windowHeight, "shaders/chunks/vertex.glsl", "shaders/chunks/fragment.glsl")
	scene.UniformStore.SetMat4("ViewMatrix", mgl32.LookAtV(
		mgl32.Vec3{1, 0, 0}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	))
	scene.UniformStore.SetMat4("ProjectionMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 100.0))
	ant.LoadImageFileToUniform("resources/atlas.png", "Tex", scene.GlslProgram.Handle, 0)

	scene.PreRender = prerender
	scene.Render = renderChunk

	return scene
}

func prerender() {
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)
	gl.Enable(gl.DEPTH_TEST)
}

func renderChunk(uniformStore *ant.UniformStore, data *ant.RenderData) {
	viewMatrix := uniformStore.GetMat4("ViewMatrix")
	projectionMatrix := uniformStore.GetMat4("ProjectionMatrix")
	modelMatrix := data.Transform
	modelView := viewMatrix.Mul4(modelMatrix)
	// normalMatrix := modelView.Mat3()
	mvp := projectionMatrix.Mul4(modelView)
	// set uiniforms
	// uniformStore.UniformMat4("ModelViewMatrix", modelView)
	// uniformStore.UniformMat3("NormalMatrix", normalMatrix)
	uniformStore.UniformMat4("MVP", mvp)
	// draw
	gl.BindVertexArray(data.Vao)
	gl.DrawElements(gl.TRIANGLES, data.IndicesCount, gl.UNSIGNED_INT, nil)
}
