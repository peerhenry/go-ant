package ant

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VboData struct {
	location   uint32
	vertexSize int32
	handle     uint32
}

type VaoBuilder struct {
	vbos         [5]VboData // hardcoded max 5 vbos
	ibo          uint32
	useIbo       bool
	numberOfVbos int
}

func (self *VaoBuilder) addVertexBuffer(location uint32, vertexSize int32, data []float32) {
	vbo := makeFloatVbo(data)
	vboData := VboData{location, vertexSize, vbo}
	self.vbos[self.numberOfVbos] = vboData
	self.numberOfVbos += 1
}

func (self *VaoBuilder) addIndexBuffer(data []uint32) {
	self.ibo = makeIbo(data)
	self.useIbo = true
}

func (self *VaoBuilder) build() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	for i := 0; i < self.numberOfVbos; i++ {
		vbo := self.vbos[i]
		gl.EnableVertexAttribArray(vbo.location)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handle)
		gl.VertexAttribPointer(vbo.location, vbo.vertexSize, gl.FLOAT, false, 0, nil)
	}

	if self.useIbo {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, self.ibo)
	}
	gl.BindVertexArray(0)
	return vao
}
