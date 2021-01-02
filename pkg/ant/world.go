package ant

type GameWorld struct {
	renderState *GameRenderState
	glslProgram *GLSLProgram
	objects     []*GameObject
	// drawables   []*Drawable
	// updatables  []*Updatable
}

func (self *GameWorld) update() {
	for _, object := range self.objects {
		object.update()
	}
}

func (self *GameWorld) render() {
	self.glslProgram.Use()
	for _, object := range self.objects {
		// object.draw(self.renderState)
		object.draw()
	}
}
