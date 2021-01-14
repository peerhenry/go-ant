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
	NORTH  = 1
	EAST   = 2
	SOUTH  = 3
	WEST   = 4
	TOP    = 5
	BOTTOM = 6
)
