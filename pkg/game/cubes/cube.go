package cubes

import (
	"time"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Placement struct {
	position    mgl32.Vec3
	orientation mgl32.Quat
}

func createCube(position mgl32.Vec3, orientation mgl32.Quat, cubeType int) *ant.GameObject {
	vao := buildCubeVao(cubeType)
	placement := &Placement{position, orientation}
	indicesCount := int32(len(cubeIndices))
	rotationSpeed := 2.0
	return &ant.GameObject{
		Update: func(dt *time.Duration) {
			dAngle := float32(rotationSpeed * dt.Seconds())
			rotation := mgl32.QuatRotate(dAngle, mgl32.Vec3{0, 0, 1})
			placement.orientation = rotation.Mul(placement.orientation)
		},
		Draw: func(uniformStore *ant.UniformStore) {
			// get values for uniform calculations
			viewMatrix := uniformStore.GetMat4("ViewMatrix")
			projectionMatrix := uniformStore.GetMat4("ProjectionMatrix")
			// calculate uniforms
			var positionMatrix mgl32.Mat4 = mgl32.Translate3D(placement.position.X(), placement.position.Y(), placement.position.Z())
			var orientationMatrix mgl32.Mat4 = placement.orientation.Mat4()
			modelMatrix := positionMatrix.Mul4(orientationMatrix)
			modelView := viewMatrix.Mul4(modelMatrix)
			normalMatrix := modelView.Mat3()
			mvp := projectionMatrix.Mul4(modelView)
			// set uiniforms
			uniformStore.UniformMat4("ModelViewMatrix", modelView)
			uniformStore.UniformMat3("NormalMatrix", normalMatrix)
			uniformStore.UniformMat4("MVP", mvp)
			// draw
			gl.BindVertexArray(vao)
			gl.DrawElements(gl.TRIANGLES, indicesCount, gl.UNSIGNED_INT, nil)
		},
	}
}

func buildCubeVao(cubeType int) uint32 {
	var positions []float32
	for _, pos := range cubePositions {
		positions = append(positions, 0.4*pos)
	}
	vaoBuilder := new(ant.VaoBuilder)
	vaoBuilder.AddVertexBuffer(0, 3, &positions)
	vaoBuilder.AddVertexBuffer(1, 3, &cubeNormals)
	vaoBuilder.AddVertexBuffer(2, 2, getCubeUvs(cubeType))
	vaoBuilder.AddIndexBuffer(&cubeIndices)
	return vaoBuilder.Build()
}

func getCubeUvs(cubeType int) *[]float32 {
	switch cubeType {
	case GRASS:
		return &cubeUvsGrass
	case DIRT:
		return &cubeUvsDirt
	case STONE:
		return &cubeUvsStone
	case SAND:
		return &cubeUvsSand
	default:
		return &cubeUvsGrass
	}
}
