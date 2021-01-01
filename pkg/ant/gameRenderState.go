package ant

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// todo: make generic, maybe use a dictionary

type GameRenderState struct {
	viewMatrix           mgl32.Mat4
	projectionMatrix     mgl32.Mat4
	glslProgramHandle    uint32
	modelViewLocation    int32
	normalMatrixLocation int32
	projectionLocation   int32
	mvpLocation          int32
}

func buildGameRenderState(glslProgramHandle uint32, windowWidth, windowHeight int) GameRenderState {
	view := mgl32.LookAtV(
		mgl32.Vec3{5, 3, 3}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	)
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0)

	return GameRenderState{
		viewMatrix:           view,
		projectionMatrix:     projection,
		glslProgramHandle:    glslProgramHandle,
		modelViewLocation:    gl.GetUniformLocation(glslProgramHandle, gl.Str("ModelViewMatrix\x00")),
		normalMatrixLocation: gl.GetUniformLocation(glslProgramHandle, gl.Str("NormalMatrix\x00")),
		projectionLocation:   gl.GetUniformLocation(glslProgramHandle, gl.Str("ProjectionMatrix\x00")),
		mvpLocation:          gl.GetUniformLocation(glslProgramHandle, gl.Str("MVP\x00")),
	}
}

func (self *GameRenderState) SetModelMatrix(model mgl32.Mat4) {
	modelview := self.viewMatrix.Mul4(model)
	gl.UniformMatrix4fv(self.modelViewLocation, 1, false, &modelview[0])
	normalMatrix := modelview.Mat3()
	gl.UniformMatrix4fv(self.normalMatrixLocation, 1, false, &normalMatrix[0])
	mvp := self.projectionMatrix.Mul4(modelview)
	gl.UniformMatrix4fv(self.mvpLocation, 1, false, &mvp[0])
}
