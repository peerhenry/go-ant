package ant

type GameWorld struct {
	Uniforms    *UniformStore
	GlslProgram *GLSLProgram
	Objects     []*GameObject
}

func (self *GameWorld) Update() {
	for _, object := range self.Objects {
		object.Update()
	}
}

func (self *GameWorld) Render() {
	self.GlslProgram.Use()
	for _, object := range self.Objects {
		object.Draw(self.Uniforms)
	}
}
