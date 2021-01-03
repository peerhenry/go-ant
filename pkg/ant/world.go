package ant

type GameWorld struct {
	uniforms    *UniformStore
	glslProgram *GLSLProgram
	objects     []*GameObject
}

func (self *GameWorld) update() {
	for _, object := range self.objects {
		object.update()
	}
}

func (self *GameWorld) render() {
	self.glslProgram.Use()
	for _, object := range self.objects {
		object.draw(self.uniforms)
	}
}
