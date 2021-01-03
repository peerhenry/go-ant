package quad

import "ant.com/ant/pkg/ant"

func buildQuadWorld() ant.GameWorld {
	glslProgram := ant.InitGlslProgram("shaders/quad/vertex.glsl", "shaders/quad/fragment.glsl")
	quad := createQuad()
	objects := []*ant.GameObject{&quad}
	uniforms := ant.CreateUniformStore(glslProgram.Handle, true)
	ant.LoadImageFileToUniform("resources/atlas.png", "Tex", glslProgram.Handle)
	return ant.GameWorld{
		Uniforms:    uniforms,
		GlslProgram: &glslProgram,
		Objects:     objects,
	}
}
