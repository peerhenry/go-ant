package ant

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Placement struct {
	position    mgl32.Vec3
	orientation mgl32.Quat
}

func createCube(position mgl32.Vec3, orientation mgl32.Quat) *GameObject {
	vao := buildCubeVao()
	placement := &Placement{position, orientation}
	indexLength := int32(len(cubeIndices))
	return &GameObject{
		update: func() {
			rotation := mgl32.QuatRotate(0.01, mgl32.Vec3{0, 0, 1})
			placement.orientation = rotation.Mul(placement.orientation)
		},
		draw: func(uniformStore *UniformStore) {
			// get values for uniform calculations
			viewMatrix := uniformStore.getMat4("ViewMatrix")
			projectionMatrix := uniformStore.getMat4("ProjectionMatrix")
			// calculate uniforms
			var positionMatrix mgl32.Mat4 = mgl32.Translate3D(placement.position.X(), placement.position.Y(), placement.position.Z())
			var orientationMatrix mgl32.Mat4 = placement.orientation.Mat4()
			modelMatrix := positionMatrix.Mul4(orientationMatrix)
			modelView := viewMatrix.Mul4(modelMatrix)
			normalMatrix := modelView.Mat3()
			mvp := projectionMatrix.Mul4(modelView)
			// set uiniforms
			uniformStore.uniformMat4("ModelViewMatrix", modelView)
			uniformStore.uniformMat3("NormalMatrix", normalMatrix)
			uniformStore.uniformMat4("MVP", mvp)
			// draw
			gl.BindVertexArray(vao)
			gl.DrawElements(gl.TRIANGLES, indexLength, gl.UNSIGNED_INT, nil)
		},
	}
}

func buildCubeVao() uint32 {
	var positions []float32
	for _, pos := range cubePositions {
		positions = append(positions, 0.4*pos)
	}
	vaoBuilder := new(VaoBuilder)
	vaoBuilder.addVertexBuffer(0, 3, positions)
	vaoBuilder.addVertexBuffer(1, 3, cubeNormals)
	vaoBuilder.addVertexBuffer(2, 2, cubeUvs)
	vaoBuilder.addIndexBuffer(cubeIndices)
	return vaoBuilder.build()
}
