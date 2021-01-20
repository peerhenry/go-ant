package ant

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	phi       float64
	theta     float64
	Position  mgl64.Vec3
	Direction mgl64.Vec3
	Right     mgl64.Vec3
	relup     mgl64.Vec3 // for calculating view frustum
}

const PHI_MAX = 0.5*math.Pi - 0.001

func NewCamera() *Camera {
	return &Camera{
		phi:       0,
		theta:     0,
		Position:  mgl64.Vec3{0, 0, 0},
		Direction: mgl64.Vec3{0, 1, 0},
		Right:     mgl64.Vec3{1, 0, 0},
		relup:     mgl64.Vec3{0, 0, 0},
	}
}

func ToVec3(v mgl64.Vec3) Vec3 {
	return Vec3{float32(v[0]), float32(v[1]), float32(v[2])}
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
	self.Direction = mgl64.Vec3{cost * cosp, sint * cosp, sinp}
	self.Right = mgl64.Vec3{sint, -cost, 0}
	self.relup = mgl64.Vec3{-cost * sinp, -sint * sinp, cosp} // for view frustum
}

func (self *Camera) Translate(ds mgl64.Vec3) {
	self.Position = self.Position.Add(ds)
}

func (self *Camera) CalculateViewMatrix() mgl32.Mat4 {
	target := self.Position.Add(self.Direction)
	return mgl32.LookAtV(
		ToVec3(self.Position), // eye
		ToVec3(target),        // center
		Vec3{0, 0, 1},         // up
	)
}
