package quad

import "ant.com/ant/pkg/ant"

func buildQuadWorld(filePath string) ant.GameWorld {
	glslProgram := ant.InitGlslProgram("shaders/quad/vertex.glsl", "shaders/quad/fragment.glsl")
	quad := createQuad()
	objects := []*ant.GameObject{&quad}
	uniforms := ant.CreateUniformStore(glslProgram.Handle, true)
	ant.LoadImageFileToUniform(filePath, "Tex", glslProgram.Handle)
	return ant.GameWorld{
		Uniforms:    uniforms,
		GlslProgram: &glslProgram,
		Objects:     objects,
	}
}
