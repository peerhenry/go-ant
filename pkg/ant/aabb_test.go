package ant

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestIntersection(t *testing.T) {
	one := AABB64{
		Min: mgl64.Vec3{0, 0, 0},
		Max: mgl64.Vec3{20, 30, 40},
	}
	two := AABB64{
		Min: mgl64.Vec3{10, 20, -10},
		Max: mgl64.Vec3{50, 60, 10},
	}
	expect := AABB64{
		Min: mgl64.Vec3{10, 20, 0},
		Max: mgl64.Vec3{20, 30, 10},
	}
	// act
	result := one.Intersection(two)
	// assert
	if !result.ApproxEqual(expect) {
		t.Errorf("expected %s , got %s", expect.ToString(), result.ToString())
		return
	}
}

func TestLineIntersectsFace(t *testing.T) {
	expectIntersectFace(t, mgl64.Vec3{0.5, 0.5, 2}, UP)
	expectIntersectFace(t, mgl64.Vec3{0.5, 0.5, -1}, DOWN)
	expectIntersectFace(t, mgl64.Vec3{2, 0.5, 0.5}, EAST)
	expectIntersectFace(t, mgl64.Vec3{-1, 0.5, 0.5}, WEST)
	expectIntersectFace(t, mgl64.Vec3{0.5, 2, 0.5}, NORTH)
	expectIntersectFace(t, mgl64.Vec3{0.5, -1, 0.5}, SOUTH)
}

// given a line from parameter to center of AABB { {0 0 0} {1 1 1} }
// which face should it intersect?
func expectIntersectFace(t *testing.T, from mgl64.Vec3, face Face) {
	aabb := AABB64{
		Min: mgl64.Vec3{0, 0, 0},
		Max: mgl64.Vec3{1, 1, 1},
	}
	allFaces := []Face{UP, DOWN, NORTH, SOUTH, EAST, WEST}
	for _, nextFace := range allFaces {
		intersects := aabb.LineIntersectsFace(from, mgl64.Vec3{0.5, 0.5, 0.5}, nextFace)
		if nextFace == face {
			if !intersects {
				t.Errorf("Expected to intersect %s face", nextFace.ToString())
			}
		} else {
			if intersects {
				t.Errorf("Did not expect to intersect %s face", nextFace.ToString())
			}
		}
	}
}
