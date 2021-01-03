package ant

type GameObject struct {
	update func()
	draw   func(uniforms *UniformStore)
}
