package ant

type Scene struct {
	UniformStore *UniformStore
	GlslProgram  *GLSLProgram
	objects      []*RenderData
	PreRender    func()
	Render       func(*UniformStore, *RenderData)
}

func CreateScene(windowWidth, windowHeight int, vertexShaderPath, fragmentShaderPath string) *Scene {
	glslProgram := InitGlslProgram(vertexShaderPath, fragmentShaderPath)
	uniformStore := CreateUniformStore(glslProgram.Handle, true)

	return &Scene{
		UniformStore: uniformStore,
		GlslProgram:  &glslProgram,
		PreRender:    func() {},
	}
}

func (self *Scene) AddRenderData(renderData *RenderData) int {
	newIndex := len(self.objects)
	self.objects = append(self.objects, renderData)
	return newIndex
}

func (self *Scene) Draw() {
	self.PreRender()
	for _, data := range self.objects {
		self.GlslProgram.Use()
		self.Render(self.UniformStore, data)
	}
}
