package quad

import "ant.com/ant/pkg/ant"

func buildQuadScene(filePath string) ant.Scene {
	glslProgram := ant.InitGlslProgram("shaders/quad/vertex.glsl", "shaders/quad/fragment.glsl")
	quad := createQuad()
	objects := []*ant.GameObject{&quad}
	uniforms := ant.CreateUniformStore(glslProgram.Handle, true)
	ant.LoadImageFileToUniform(filePath, "Tex", glslProgram.Handle)
	return ant.Scene{
		Uniforms:    uniforms,
		GlslProgram: &glslProgram,
		Objects:     objects,
	}
}
