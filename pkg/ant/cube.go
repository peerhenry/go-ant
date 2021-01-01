package ant

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// todo: refactor with vaoBuilder

func createCube(position mgl32.Vec3, orientation mgl32.Quat) GameObject {

	var positions []float32
	for _, pos := range cubePositions {
		positions = append(positions, 0.4*pos)
	}

	// vao := makeVao(positions, cubeNormals, cubeIndices)
	vaoBuilder := new(VaoBuilder)
	vaoBuilder.addVertexBuffer(0, 3, positions)
	vaoBuilder.addVertexBuffer(1, 3, cubeNormals)
	vaoBuilder.addIndexBuffer(cubeIndices)
	vao := vaoBuilder.build()
	return GameObject{
		vao,
		int32(len(cubeIndices)),
		position,
		orientation,
	}
}

func (self *GameObject) update() {
	rotation := mgl32.QuatRotate(0.01, mgl32.Vec3{0, 0, 1})
	self.orientation = rotation.Mul(self.orientation)
}

func (self *GameObject) draw(renderState *GameRenderState) {
	var positionMatrix mgl32.Mat4 = mgl32.Translate3D(self.position.X(), self.position.Y(), self.position.Z())
	var orientationMatrix mgl32.Mat4 = self.orientation.Mat4()
	modelMatrix := positionMatrix.Mul4(orientationMatrix)
	renderState.SetModelMatrix(modelMatrix)
	gl.BindVertexArray(self.vao)
	gl.DrawElements(gl.TRIANGLES, self.indexLength, gl.UNSIGNED_INT, nil)
}
