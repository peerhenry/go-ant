package ant

func buildQuadWorld() GameWorld {
	glslProgram := initGlslProgram("shaders/quad/vertex.glsl", "shaders/quad/fragment.glsl")
	quad := createQuad()
	objects := []*GameObject{&quad}
	uniforms := createUniformStore(glslProgram.handle, true)
	loadImageFileToUniform("resources/atlas.png", "Tex", glslProgram.handle)
	return GameWorld{
		uniforms:    uniforms,
		glslProgram: &glslProgram,
		objects:     objects,
	}
}
