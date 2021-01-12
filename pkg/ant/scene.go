package ant

import "time"

type Scene struct {
	Uniforms    *UniformStore
	GlslProgram *GLSLProgram
	Objects     []*GameObject
}

func (self *Scene) Update(dt *time.Duration) {
	for _, object := range self.Objects {
		object.Update(dt)
	}
}

func (self *Scene) Render() {
	self.GlslProgram.Use()
	for _, object := range self.Objects {
		object.Draw(self.Uniforms)
	}
}

func (self *Scene) Add(object *GameObject) {
	self.Objects = append(self.Objects, object)
}

func (self *Scene) Remove(object *GameObject) {
	panic("todo: cannot remove GameObject yet")
}
