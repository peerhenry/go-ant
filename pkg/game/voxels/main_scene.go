package voxels

import (
	"ant.com/ant/pkg/ant"
	"ant.com/ant/pkg/game/voxels/chunks"

	"github.com/go-gl/mathgl/mgl32"
)

func buildScene(windowWidth, windowHeight int) ant.Scene {
	glslProgram := ant.InitGlslProgram("shaders/ads/vertex.glsl", "shaders/ads/fragment.glsl")
	uniformStore := setupUniforms(glslProgram.Handle, windowWidth, windowHeight)
	objects := createGameObjects()

	return ant.Scene{
		Uniforms:    uniformStore,
		GlslProgram: &glslProgram,
		Objects:     objects,
	}
}

func setupUniforms(glslProgramHandle uint32, windowWidth, windowHeight int) *ant.UniformStore {
	uniformStore := ant.CreateUniformStore(glslProgramHandle, true)
	uniformStore.SetMat4("ViewMatrix", mgl32.LookAtV(
		mgl32.Vec3{1, 0, 0}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	))
	uniformStore.SetMat4("ProjectionMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 100.0))
	ant.LoadImageFileToUniform("resources/atlas.png", "Tex", glslProgramHandle, 0)
	return uniformStore
}

func createGameObjects() []*ant.GameObject {
	chunk1 := chunks.BuildChunkGameObject(Vec3{0, 0, -32})
	chunk2 := chunks.BuildChunkGameObject(Vec3{-64, 0, -32})
	chunk3 := chunks.BuildChunkGameObject(Vec3{0, -64, -32})
	chunk4 := chunks.BuildChunkGameObject(Vec3{-64, -64, -32})
	return []*ant.GameObject{chunk1, chunk2, chunk3, chunk4}
}
