package chunks

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type void struct{}

var VOID void

type Vec2 = mgl32.Vec2
type Vec3 = mgl32.Vec3
type Mat4 = mgl32.Mat4
type Mat3 = mgl32.Mat3

type Block = int

const (
	UNDERGROUND Block = -1
	AIR         Block = 0
	GRASS       Block = 1
	DIRT        Block = 2
	STONE       Block = 3
	SAND        Block = 4
	TRUNK       Block = 5
	LEAVES      Block = 6
	SNOWDIRT    Block = 7
	WATER       Block = 8

	GRASS_1 Block = 101
	GRASS_2 Block = 102
	GRASS_3 Block = 103
	GRASS_4 Block = 104
	GRASS_5 Block = 105
	GRASS_6 Block = 106
	GRASS_7 Block = 107
	GRASS_8 Block = 108

	RED_FLOWER     Block = 109
	YELLOW_FLOWER  Block = 110
	RED_MUSHROOM   Block = 111
	BROWN_MUSHROOM Block = 112
)

func VoxelIsTransparent(voxel Block) bool {
	return voxel == AIR || voxel == WATER || VoxelIsXShaped(voxel)
}

func VoxelIsWalkThrough(voxel Block) bool {
	return voxel == GRASS_1 ||
		voxel == GRASS_2 ||
		voxel == GRASS_3 ||
		voxel == GRASS_4 ||
		voxel == GRASS_5 ||
		voxel == GRASS_6 ||
		voxel == GRASS_7 ||
		voxel == GRASS_8 ||
		voxel == RED_FLOWER ||
		voxel == YELLOW_FLOWER ||
		voxel == RED_MUSHROOM ||
		voxel == BROWN_MUSHROOM
}

type Face = ant.Face

const (
	NORTH      Face = ant.NORTH
	EAST       Face = ant.EAST
	SOUTH      Face = ant.SOUTH
	WEST       Face = ant.WEST
	TOP        Face = ant.UP   // deprecated
	BOTTOM     Face = ant.DOWN // deprecated
	UP         Face = ant.UP
	DOWN       Face = ant.DOWN
	NORTH_EAST Face = ant.NORTH_EAST
	SOUTH_EAST Face = ant.SOUTH_EAST
	SOUTH_WEST Face = ant.SOUTH_WEST
	NORTH_WEST Face = ant.NORTH_WEST
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
		faces = append(faces, WEST)
	} else if vec[0] < 0 {
		faces = append(faces, EAST)
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
