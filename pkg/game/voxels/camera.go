package voxels

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	phi       float64
	theta     float64
	position  mgl32.Vec3
	direction mgl32.Vec3
	right     mgl32.Vec3
	relup     mgl32.Vec3 // for calculating view frustum
}

const PHI_MAX = 0.5*math.Pi - 0.001

func NewCamera() *Camera {
	return &Camera{
		phi:       0,
		theta:     0,
		position:  mgl32.Vec3{0, 0, 0},
		direction: mgl32.Vec3{0, 1, 0},
		right:     mgl32.Vec3{1, 0, 0},
		relup:     mgl32.Vec3{0, 0, 0},
	}
}

func ToVec3(x float64, y float64, z float64) Vec3 {
	return Vec3{float32(x), float32(y), float32(z)}
}

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
	sint := math.Sin(self.theta)
	cost := math.Cos(self.theta)
	sinp := math.Sin(self.phi)
	cosp := math.Cos(self.phi)
	self.direction = Vec3{float32(cost * cosp), float32(sint * cosp), float32(sinp)}
	self.right = ToVec3(sint, -cost, 0)
	self.relup = ToVec3(-cost*sinp, -sint*sinp, cosp) // for view frustum
}

func (self *Camera) Translate(dx, dy, dz float64) {
	self.position = self.position.Add(ToVec3(dx, dy, dz))
}

func (self *Camera) CalculateViewMatrix() mgl32.Mat4 {
	target := self.position.Add(self.direction)
	return mgl32.LookAtV(
		self.position,       // eye
		target,              // center
		mgl32.Vec3{0, 0, 1}, // up
	)
}
