package quad

import (
	"time"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func createQuad() ant.GameObject {
	var quadPositions []float32 = []float32{
		-1.0, -1.0,
		-1.0, 1.0,
		1.0, -1.0,
		1.0, 1.0,
	}
	var uvs []float32 = []float32{
		0.0, 1.0,
		0.0, 0.0,
		1.0, 1.0,
		1.0, 0.0,
	}
	builder := new(ant.VaoBuilder)
	builder.AddVertexBuffer(0, 2, &quadPositions)
	builder.AddVertexBuffer(1, 2, &uvs)
	builder.AddIndexBuffer(&[]uint32{
		0, 1, 2, 2, 1, 3,
	})
	vao := builder.Build()
	return ant.GameObject{
		Update: func(dt *time.Duration) {},
		Draw: func(uniformStore *ant.UniformStore) {
			gl.BindVertexArray(vao)
			gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		},
	}
}
