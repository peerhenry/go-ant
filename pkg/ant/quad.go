package ant

// type Quad struct {
// 	vao uint32
// }

// func createQuad() Quad {
// 	builder := new(VaoBuilder)
// 	builder.numberOfVbos = 0
// 	var quadPositions []float32 = []float32{
// 		-1.0, -1.0,
// 		-1.0, 1.0,
// 		1.0, -1.0,
// 		1.0, 1.0,
// 	}
// 	var uvs []float32 = []float32{
// 		0.0, 0.0,
// 		0.0, 1.0,
// 		1.0, 0.0,
// 		1.0, 1.0,
// 	}
// 	builder.addVertexBuffer(0, 2, &quadPositions)
// 	builder.addVertexBuffer(1, 2, &uvs)
// 	builder.addIndexBuffer(&[]uint32{
// 		0, 1, 2, 2, 1, 3,
// 	})
// 	vao := builder.build()
// 	return Quad{vao}
// }

// func (self *Quad) update() {}

// func (self *Quad) draw() {
// 	gl.BindVertexArray(self.vao)
// 	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
// }
