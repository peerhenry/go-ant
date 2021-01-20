package chunks

import (
	"math"
	"time"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl64"
)

const fallAcceleration = 10.0 // m/s/s
const maxFallSpeed = -30.0    // m/s
const playerHeight = 1.8

type Player struct {
	Camera    *ant.Camera
	isFalling bool
	Velocity  mgl64.Vec3
	world     *ChunkWorld
}

func NewPlayer(camera *ant.Camera, world *ChunkWorld) *Player {
	return &Player{
		Camera:    camera,
		isFalling: true,
		Velocity:  mgl64.Vec3{0, 0, 0},
	}
}

func (self *Player) Update(dt *time.Duration) {
	// check if 10% of a voxel under players feet is a solid voxel
	feet := self.Camera.Position.Sub(mgl64.Vec3{0, 0, playerHeight + 0.1})
	if self.IntersectsWithSolidVoxel(feet) {
		self.isFalling = false
		// self.Camera.Position[2] =
	}
	// update velocity
	if self.isFalling {
		dv := dt.Seconds() * fallAcceleration
		newFallSpeed := math.Max(self.Velocity[2]-dv, maxFallSpeed)
		self.Velocity = mgl64.Vec3{self.Velocity[0], self.Velocity[1], newFallSpeed}
	}
	// update positions
	ds := self.Velocity.Mul(dt.Seconds())
	self.Camera.Translate(ds)
}

func (self *Player) IntersectsWithSolidVoxel(location mgl64.Vec3) bool {
	// var indexCoordinate []IndexCoordinate

	// location[0]

	// voxel := self.world.GetVoxelAt(indexCoordinate)
	// return voxel != AIR && voxel != WATER
	return false
}

// todo TDD this
func (self *Player) ToRegionCoord(location mgl64.Vec3) []IndexCoordinate {
	var output []IndexCoordinate
	return output
}
