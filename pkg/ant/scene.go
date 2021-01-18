package ant

type Scene struct {
	UniformStore    *UniformStore
	GlslProgram     *GLSLProgram
	renderIndexHead int
	objects         map[int]*RenderData
	PreRender       func()
	Render          func(*UniformStore, *RenderData)
}

func CreateScene(windowWidth, windowHeight int, vertexShaderPath, fragmentShaderPath string) *Scene {
	glslProgram := InitGlslProgram(vertexShaderPath, fragmentShaderPath)
	uniformStore := CreateUniformStore(glslProgram.Handle, true)

	return &Scene{
		UniformStore:    uniformStore,
		GlslProgram:     &glslProgram,
		PreRender:       func() {},
		objects:         make(map[int]*RenderData),
		renderIndexHead: 1,
	}
}

func (self *Scene) AddRenderData(renderData *RenderData) int {
	newIndex := self.renderIndexHead
	self.objects[newIndex] = renderData
	self.renderIndexHead++
	return newIndex
}

func (self *Scene) ReplaceRenderData(index int, renderData *RenderData) {
	self.objects[index] = renderData
}

func (self *Scene) RemoveRenderData(index int) {
	delete(self.objects, index)
}

func (self *Scene) Draw() {
	self.PreRender()
	for _, data := range self.objects {
		self.GlslProgram.Use()
		self.Render(self.UniformStore, data)
	}
}
