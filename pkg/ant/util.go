package ant

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Vec2 = mgl32.Vec2
type Vec3 = mgl32.Vec3
type Mat4 = mgl32.Mat4
type Mat3 = mgl32.Mat3

type Face int32

const (
	NORTH      Face = 0
	EAST       Face = 1
	SOUTH      Face = 2
	WEST       Face = 3
	UP         Face = 4
	DOWN       Face = 5
	NORTH_EAST Face = 6
	SOUTH_EAST Face = 7
	SOUTH_WEST Face = 8
	NORTH_WEST Face = 9
)

func (face Face) ToString() string {
	switch face {
	case NORTH:
		return "NORTH"
	case SOUTH:
		return "SOUTH"
	case EAST:
		return "EAST"
	case WEST:
		return "WEST"
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	}
	return "UNKNOWN"
}

// based on method from page 189 of "realtime collision detection"
// note that the order of abcd (clockwise/counterclockwise) matters
func IntersectLineQuad64(p, q, a, b, c, d mgl64.Vec3) bool {
	pq := q.Sub(p)
	pa := a.Sub(p)
	pb := b.Sub(p)
	pc := c.Sub(p)
	// Determine which triangle to test against by testing against diagonal first
	m := pc.Cross(pq)
	v := pa.Dot(m) // ScalarTriple(pq, pa, pc)

	if v >= 0 {
		// test intersection against triangle abc
		u := -pb.Dot(m) // ScalarTriple(pq, pc, pb)
		if u < 0 {
			return false
		}
		w := ScalarTriple64(pq, pb, pa)
		if w < 0 {
			return false
		}
	} else {
		pd := d.Sub(p)
		u := pd.Dot(m) // ScalarTriple(pq, pd, pc)
		if u < 0 {
			return false
		}
		w := ScalarTriple64(pq, pa, pd)
		if w < 0 {
			return false
		}
	}

	return true
}

func ScalarTriple64(a, b, c mgl64.Vec3) float64 {
	m := c.Cross(a)
	return b.Dot(m)
}
