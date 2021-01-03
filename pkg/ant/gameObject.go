package ant

type UpdateCallback func()
type DrawCallback func(uniforms *UniformStore)

type GameObject struct {
	Update UpdateCallback
	Draw   DrawCallback
}
