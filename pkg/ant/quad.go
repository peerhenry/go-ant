package ant

import "github.com/go-gl/gl/v4.1-core/gl"

func createQuad() GameObject {
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
	builder := new(VaoBuilder)
	builder.addVertexBuffer(0, 2, quadPositions)
	builder.addVertexBuffer(1, 2, uvs)
	builder.addIndexBuffer([]uint32{
		0, 1, 2, 2, 1, 3,
	})
	vao := builder.build()
	return GameObject{
		update: func() {},
		draw: func(uniformStore *UniformStore) {
			gl.BindVertexArray(vao)
			gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		},
	}
}

// func (self *GameObject) update() {}

// func (self *GameObject) draw() {
// 	gl.BindVertexArray(self.vao)
// 	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
// }
