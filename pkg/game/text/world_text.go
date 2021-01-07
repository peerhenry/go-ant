package text

import "ant.com/ant/pkg/ant"

func buildQuadWorld(windowWidth, windowHeight int) ant.GameWorld {
	glslProgram := ant.InitGlslProgram("shaders/text/vertex.glsl", "shaders/text/fragment.glsl")
	quad := createQuad(windowWidth, windowHeight)
	objects := []*ant.GameObject{&quad}
	uniforms := ant.CreateUniformStore(glslProgram.Handle, true)
	// ant.LoadImageFileToUniform("resources/text-atlas.png", "TextAtlas", glslProgram.Handle)
	// ant.LoadImageFileToUniform("resources/text-atlas-monospace-white-outlined-on-alpha-extra.png", "TextAtlas", glslProgram.Handle)
	ant.LoadImageFileToUniform("resources/text-atlas-monospace-white-on-alpha.png", "TextAtlas", glslProgram.Handle)
	return ant.GameWorld{
		Uniforms:    uniforms,
		GlslProgram: &glslProgram,
		Objects:     objects,
	}
}
