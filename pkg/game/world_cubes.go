package game

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl32"
)

func buildCubeWorld(windowWidth, windowHeight int) ant.GameWorld {
	glslProgram := ant.InitGlslProgram("shaders/ads/vertex.glsl", "shaders/ads/fragment.glsl")
	uniformStore := setupUniforms(glslProgram.Handle, windowWidth, windowHeight)
	objects := createGameObjects()

	return ant.GameWorld{
		Uniforms:    uniformStore,
		GlslProgram: &glslProgram,
		Objects:     objects,
	}
}

func setupUniforms(glslProgramHandle uint32, windowWidth, windowHeight int) *ant.UniformStore {
	uniformStore := ant.CreateUniformStore(glslProgramHandle, true)
	uniformStore.SetMat4("ViewMatrix", mgl32.LookAtV(
		mgl32.Vec3{5, 3, 3}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	))
	uniformStore.SetMat4("ProjectionMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0))
	ant.LoadImageFileToUniform("resources/atlas.png", "Tex", glslProgramHandle)
	return uniformStore
}

func createGameObjects() []*ant.GameObject {
	cube1 := createCube(mgl32.Vec3{2, 0, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}), GRASS)
	cube2 := createCube(mgl32.Vec3{0, 2, 0}, mgl32.QuatRotate(1, mgl32.Vec3{1, 0, 0}), STONE)
	cube3 := createCube(mgl32.Vec3{0, 0, 2}, mgl32.QuatRotate(1, mgl32.Vec3{0, 1, 0}), DIRT)
	cube4 := createCube(mgl32.Vec3{-2, 0, 0}, mgl32.QuatRotate(2, mgl32.Vec3{0, 1, 0}), SAND)
	return []*ant.GameObject{cube1, cube2, cube3, cube4}
}
