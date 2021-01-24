package chunks

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Vec2 = mgl32.Vec2
type Vec3 = mgl32.Vec3
type Mat4 = mgl32.Mat4
type Mat3 = mgl32.Mat3

const (
	UNDERGROUND = -1
	AIR         = 0
	GRASS       = 1
	DIRT        = 2
	STONE       = 3
	SAND        = 4
	TRUNK       = 5
	LEAVES      = 6
	SNOWDIRT    = 7
	WATER       = 8
)

func VoxelIsTransparent(voxel int) bool {
	return voxel == AIR || voxel == WATER
}

type Face = ant.Face

const (
	NORTH  Face = ant.NORTH
	EAST   Face = ant.EAST
	SOUTH  Face = ant.SOUTH
	WEST   Face = ant.WEST
	TOP    Face = ant.UP   // deprecated
	BOTTOM Face = ant.DOWN // deprecated
	UP     Face = ant.UP
	DOWN   Face = ant.DOWN
)

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

const (
	MaxInt  = 1<<(UintSize-1) - 1 // 1<<31 - 1 or 1<<63 - 1
	MinInt  = -MaxInt - 1         // -1 << 31 or -1 << 63
	MaxUint = 1<<UintSize - 1     // 1<<32 - 1 or 1<<64 - 1
)

func Vec3_32_to_64(vec Vec3) mgl64.Vec3 {
	return mgl64.Vec3{float64(vec[0]), float64(vec[1]), float64(vec[2])}
}

func Vec3_64_to_32(vec mgl64.Vec3) Vec3 {
	return Vec3{float32(vec[0]), float32(vec[1]), float32(vec[2])}
}

func GetFacingFaces(vec mgl64.Vec3) []Face {
	var faces []Face
	if vec[0] > 0 {
		faces = append(faces, EAST)
	} else if vec[0] < 0 {
		faces = append(faces, WEST)
	}
	if vec[1] > 0 {
		faces = append(faces, SOUTH)
	} else if vec[1] < 0 {
		faces = append(faces, NORTH)
	}
	if vec[2] > 0 {
		faces = append(faces, DOWN)
	} else if vec[2] < 0 {
		faces = append(faces, UP)
	}
	return faces
}
