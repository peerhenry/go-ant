package ant

import "time"

type GameWorld struct {
	Uniforms    *UniformStore
	GlslProgram *GLSLProgram
	Objects     []*GameObject
}

func (self *GameWorld) Update(dt *time.Duration) {
	for _, object := range self.Objects {
		object.Update(dt)
	}
}

func (self *GameWorld) Render() {
	self.GlslProgram.Use()
	for _, object := range self.Objects {
		object.Draw(self.Uniforms)
	}
}
