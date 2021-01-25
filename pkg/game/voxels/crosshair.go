package voxels

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Crosshair struct {
	Draw func()
}

func BuildCrosshair() *Crosshair {
	ch := new(Crosshair)
	glslProgram := ant.InitGlslProgram("shaders/crosshair/vertex.glsl", "shaders/crosshair/fragment.glsl")

	data := []float32{0, 0, 0}
	vaoBuilder := new(ant.VaoBuilder)
	vaoBuilder.AddVertexBuffer(0, 3, &data)
	vao := vaoBuilder.Build()

	ch.Draw = func() {
		glslProgram.Use()
		gl.Disable(gl.CULL_FACE)
		gl.Disable(gl.DEPTH_TEST)
		gl.Enable(gl.POINT_SIZE)
		gl.PointSize(5)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.POINTS, 0, 1)
		gl.BindVertexArray(0)
	}
	return ch
}
