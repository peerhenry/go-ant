package ant

import (
	"github.com/go-gl/mathgl/mgl32"
)

func buildCubeWorld(windowWidth, windowHeight int) GameWorld {
	glslProgram := initGlslProgram("shaders/ads/vertex.glsl", "shaders/ads/fragment.glsl")
	uniformStore := setupUniforms(glslProgram.handle, windowWidth, windowHeight)
	objects := createGameObjects()

	return GameWorld{
		uniforms:    uniformStore,
		glslProgram: &glslProgram,
		objects:     objects,
	}
}

func setupUniforms(glslProgramHandle uint32, windowWidth, windowHeight int) *UniformStore {
	uniformStore := createUniformStore(glslProgramHandle, true)
	uniformStore.setMat4("ViewMatrix", mgl32.LookAtV(
		mgl32.Vec3{5, 3, 3}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	))
	uniformStore.setMat4("ProjectionMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0))
	loadImageFileToUniform("resources/atlas.png", "Tex", glslProgramHandle)
	return uniformStore
}

func createGameObjects() []*GameObject {
	cube1 := createCube(mgl32.Vec3{2, 0, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	cube2 := createCube(mgl32.Vec3{0, 2, 0}, mgl32.QuatRotate(1, mgl32.Vec3{1, 0, 0}))
	cube3 := createCube(mgl32.Vec3{0, 0, 2}, mgl32.QuatRotate(1, mgl32.Vec3{0, 1, 0}))
	cube4 := createCube(mgl32.Vec3{-2, 0, 0}, mgl32.QuatRotate(2, mgl32.Vec3{0, 1, 0}))
	return []*GameObject{cube1, cube2, cube3, cube4}
}
