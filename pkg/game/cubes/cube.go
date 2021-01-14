package cubes

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl32"
)

type Placement struct {
	position    mgl32.Vec3
	orientation mgl32.Quat
}

func createCube(position mgl32.Vec3, orientation mgl32.Quat, cubeType int) *ant.RenderData {
	vao := buildCubeVao(cubeType)
	placement := &Placement{position, orientation}
	indicesCount := int32(len(cubeIndices))
	// rotationSpeed := 2.0
	var positionMatrix mgl32.Mat4 = mgl32.Translate3D(placement.position.X(), placement.position.Y(), placement.position.Z())
	// var orientationMatrix mgl32.Mat4 = placement.orientation.Mat4()
	return &ant.RenderData{
		Transform:    positionMatrix,
		Vao:          vao,
		IndicesCount: indicesCount,
	}
}

// Update: func(dt *time.Duration) {
// 	dAngle := float32(rotationSpeed * dt.Seconds())
// 	rotation := mgl32.QuatRotate(dAngle, mgl32.Vec3{0, 0, 1})
// 	placement.orientation = rotation.Mul(placement.orientation)
// },

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
