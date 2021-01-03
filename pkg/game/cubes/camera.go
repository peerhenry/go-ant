package cubes

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	phi      float64
	theta    float64
	position mgl32.Vec3
}

const PHI_MAX = 0.5*math.Pi - 0.001

func (self *Camera) Rotate(dtheta float64, dphi float64) {
	self.theta = self.theta + dtheta
	self.phi = self.phi + dphi
	for self.theta > 2*math.Pi {
		self.theta -= 2 * math.Pi
	}
	for self.theta < 0 {
		self.theta += 2 * math.Pi
	}
	for self.phi > PHI_MAX {
		self.phi = PHI_MAX
	}
	for self.phi < -PHI_MAX {
		self.phi = -PHI_MAX
	}
}

func (self *Camera) CalculateViewMatrix() mgl32.Mat4 {
	sint := math.Sin(self.theta)
	cost := math.Cos(self.theta)
	sinp := math.Sin(self.phi)
	cosp := math.Cos(self.phi)
	dir := mgl32.Vec3{float32(cost * cosp), float32(sint * cosp), float32(sinp)}
	target := self.position.Add(dir)
	return mgl32.LookAtV(
		self.position,       // eye
		target,              // center
		mgl32.Vec3{0, 0, 1}, // up
	)
}
