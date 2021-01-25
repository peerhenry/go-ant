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

func (self *VaoBuilder) AddVertexBuffer(location uint32, vertexSize int32, data *[]float32) {
	vbo := MakeFloatVbo(data)
	vboData := VboData{location, vertexSize, vbo}
	self.vbos[self.numberOfVbos] = vboData
	self.numberOfVbos += 1
}

func (self *VaoBuilder) AddIntegerBuffer(location uint32, vertexSize int32, data *[]int32) {
	vbo := makeIntegerVbo(data)
	vboData := VboData{location, vertexSize, vbo}
	self.vbos[self.numberOfVbos] = vboData
	self.numberOfVbos += 1
}

func (self *VaoBuilder) AddIndexBuffer(data *[]uint32) {
	self.ibo = makeIbo(data)
	self.useIbo = true
}

func (self *VaoBuilder) Build() uint32 {
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

func MakeFloatVbo(data *[]float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(*data), gl.Ptr(*data), gl.STATIC_DRAW)
	return vbo
}

func makeIntegerVbo(data *[]int32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(*data), gl.Ptr(*data), gl.STATIC_DRAW)
	return vbo
}

func makeIbo(data *[]uint32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(*data), gl.Ptr(*data), gl.STATIC_DRAW)
	return vbo
}
