package ant

import "time"

type UpdateCallback func(dt *time.Duration)
type DrawCallback func(uniforms *UniformStore)

type GameObject struct {
	Update UpdateCallback
	Draw   DrawCallback
}
