package ant

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type AABB struct {
	Min Vec3
	Max Vec3
}

func (self *AABB) Center() Vec3 {
	return self.Min.Add(self.Max).Mul(0.5)
}

func (self *AABB) Intersects(other AABB) bool {
	if self.Max[0] < other.Min[0] || self.Min[0] > other.Max[0] ||
		self.Max[1] < other.Min[1] || self.Min[1] > other.Max[1] ||
		self.Max[2] < other.Min[2] || self.Min[2] > other.Max[2] {
		return false
	}
	return true
}

// === 64 bit variant ===

type AABB64 struct {
	Min mgl64.Vec3
	Max mgl64.Vec3
}

func (self AABB64) Center() mgl64.Vec3 {
	return self.Min.Add(self.Max).Mul(0.5)
}

func (self AABB64) Validate() bool {
	return self.Min[0] < self.Max[0] && self.Min[1] < self.Max[1] && self.Min[2] < self.Max[2]
}

func (self AABB64) ToString() string {
	return "AABB64 { Min: [" +
		fmt.Sprintf("%v, %v, %v", self.Min[0], self.Min[1], self.Min[2]) +
		"] Max: [" +
		fmt.Sprintf("%v, %v, %v", self.Max[0], self.Max[1], self.Max[2]) +
		"] }"
}

func (self AABB64) ApproxEqual(other AABB64) bool {
	return self.Min.ApproxEqual(other.Min) && self.Max.ApproxEqual(other.Max)
}

func (self AABB64) Intersects(other AABB64) bool {
	if self.Max[0] < other.Min[0] || self.Min[0] > other.Max[0] ||
		self.Max[1] < other.Min[1] || self.Min[1] > other.Max[1] ||
		self.Max[2] < other.Min[2] || self.Min[2] > other.Max[2] {
		return false
	}
	return true
}

// todo: write test
func (self AABB64) Intersection(other AABB64) AABB64 {
	xmin := math.Max(self.Min[0], other.Min[0])
	xmax := math.Min(self.Max[0], other.Max[0])
	ymin := math.Max(self.Min[1], other.Min[1])
	ymax := math.Min(self.Max[1], other.Max[1])
	zmin := math.Max(self.Min[2], other.Min[2])
	zmax := math.Min(self.Max[2], other.Max[2])
	return AABB64{
		Min: mgl64.Vec3{xmin, ymin, zmin},
		Max: mgl64.Vec3{xmax, ymax, zmax},
	}
}

func (self AABB64) LineIntersectsFace(p, q mgl64.Vec3, face Face) bool {
	ds := q.Sub(p)
	switch face {
	case UP:
		if ds[2] >= 0 {
			return false
		}
		// quad must run CCW
		a := mgl64.Vec3{self.Min[0], self.Min[1], self.Max[2]}
		b := mgl64.Vec3{self.Max[0], self.Min[1], self.Max[2]}
		c := mgl64.Vec3{self.Max[0], self.Max[1], self.Max[2]}
		d := mgl64.Vec3{self.Min[0], self.Max[1], self.Max[2]}
		return IntersectLineQuad64(p, q, a, b, c, d)
	case DOWN:
		if ds[2] <= 0 {
			return false
		}
		a := mgl64.Vec3{self.Min[0], self.Min[1], self.Min[2]}
		b := mgl64.Vec3{self.Min[0], self.Max[1], self.Min[2]}
		c := mgl64.Vec3{self.Max[0], self.Max[1], self.Min[2]}
		d := mgl64.Vec3{self.Max[0], self.Min[1], self.Min[2]}
		return IntersectLineQuad64(p, q, a, b, c, d)
	case EAST:
		if ds[0] >= 0 {
			return false
		}
		a := mgl64.Vec3{self.Max[0], self.Min[1], self.Min[2]}
		b := mgl64.Vec3{self.Max[0], self.Max[1], self.Min[2]}
		c := mgl64.Vec3{self.Max[0], self.Max[1], self.Max[2]}
		d := mgl64.Vec3{self.Max[0], self.Min[1], self.Max[2]}
		return IntersectLineQuad64(p, q, a, b, c, d)
	case WEST:
		if ds[0] <= 0 {
			return false
		}
		a := mgl64.Vec3{self.Min[0], self.Min[1], self.Min[2]}
		b := mgl64.Vec3{self.Min[0], self.Min[1], self.Max[2]}
		c := mgl64.Vec3{self.Min[0], self.Max[1], self.Max[2]}
		d := mgl64.Vec3{self.Min[0], self.Max[1], self.Min[2]}
		return IntersectLineQuad64(p, q, a, b, c, d)
	case NORTH:
		if ds[1] >= 0 {
			return false
		}
		a := mgl64.Vec3{self.Min[0], self.Max[1], self.Min[2]}
		b := mgl64.Vec3{self.Min[0], self.Max[1], self.Max[2]}
		c := mgl64.Vec3{self.Max[0], self.Max[1], self.Max[2]}
		d := mgl64.Vec3{self.Max[0], self.Max[1], self.Min[2]}
		return IntersectLineQuad64(p, q, a, b, c, d)
	case SOUTH:
		if ds[1] <= 0 {
			return false
		}
		a := mgl64.Vec3{self.Min[0], self.Min[1], self.Min[2]}
		b := mgl64.Vec3{self.Max[0], self.Min[1], self.Min[2]}
		c := mgl64.Vec3{self.Max[0], self.Min[1], self.Max[2]}
		d := mgl64.Vec3{self.Min[0], self.Min[1], self.Max[2]}
		return IntersectLineQuad64(p, q, a, b, c, d)
	}
	// determine the four points for face
	return false
}
