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
	vbo := ant.MakeFloatVbo(&data)
	ch.Draw = func() {
		glslProgram.Use()
		gl.Enable(gl.POINT_SIZE)
		gl.PointSize(5)

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.EnableVertexAttribArray(0)                         // use vertex attribute array
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil) // what to do with supplied array data
		gl.DrawArrays(gl.POINTS, 0, 1)
		gl.DisableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}
	return ch
}
