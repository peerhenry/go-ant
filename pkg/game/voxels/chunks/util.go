package chunks

import "github.com/go-gl/mathgl/mgl32"

type Vec2 = mgl32.Vec2
type Vec3 = mgl32.Vec3
type Mat4 = mgl32.Mat4
type Mat3 = mgl32.Mat3

const (
	AIR   = 0
	GRASS = 1
	DIRT  = 2
	STONE = 3
	SAND  = 4
)

const (
	NORTH  int32 = 0
	EAST   int32 = 1
	SOUTH  int32 = 2
	WEST   int32 = 3
	TOP    int32 = 4
	BOTTOM int32 = 5
)

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

const (
	MaxInt  = 1<<(UintSize-1) - 1 // 1<<31 - 1 or 1<<63 - 1
	MinInt  = -MaxInt - 1         // -1 << 31 or -1 << 63
	MaxUint = 1<<UintSize - 1     // 1<<32 - 1 or 1<<64 - 1
)
