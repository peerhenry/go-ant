package ant

import "github.com/go-gl/mathgl/mgl32"

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
	uniformStore := createUniformStore(glslProgramHandle)
	// register uniforms
	uniformStore.registerUniform("ModelViewMatrix")
	// uniformStore.registerUniform("NormalMatrix")
	uniformStore.registerUniform("ProjectionMatrix")
	uniformStore.registerUniform("MVP")
	// set values
	uniformStore.setMat4("ViewMatrix", mgl32.LookAtV(
		mgl32.Vec3{5, 3, 3}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	))
	uniformStore.setMat4("ProjectionMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0))
	// todo: render textures cubes
	// loadImageFileToUniform("resources/atlas.png", "Tex", glslProgram.handle)
	return uniformStore
}

func createGameObjects() []*GameObject {
	cube1 := createCube(mgl32.Vec3{1, 0, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	cube2 := createCube(mgl32.Vec3{0, 1, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	cube3 := createCube(mgl32.Vec3{0, 0, 1}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	return []*GameObject{cube1, cube2, cube3}
}
