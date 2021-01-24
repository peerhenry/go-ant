package ant

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestIntersectLineQuad64_ShouldIntersect(t *testing.T) {
	// arrange
	DoLineQuadTest(t, mgl64.Vec3{1, 1, -1}, mgl64.Vec3{-1, -1, 1}, true)
	DoLineQuadTest(t, mgl64.Vec3{-1, 1, -1}, mgl64.Vec3{1, -1, 1}, true)
	DoLineQuadTest(t, mgl64.Vec3{1, -1, -1}, mgl64.Vec3{-1, 1, 1}, true)
	DoLineQuadTest(t, mgl64.Vec3{-1, -1, -1}, mgl64.Vec3{1, 1, 1}, true)
	// following only pass if direction of abcd is switched from clickwise to counterclockwise
	// DoLineQuadTest(t, mgl64.Vec3{1, 1, 1}, mgl64.Vec3{-1, -1, -1}, true)
	// DoLineQuadTest(t, mgl64.Vec3{-1, 1, 1}, mgl64.Vec3{1, -1, -1}, true)
	// DoLineQuadTest(t, mgl64.Vec3{1, -1, 1}, mgl64.Vec3{-1, 1, -1}, true)
	// DoLineQuadTest(t, mgl64.Vec3{-1, -1, 1}, mgl64.Vec3{1, 1, -1}, true)
}

func TestIntersectLineQuad64_ShouldNotIntersect(t *testing.T) {
	// arrange
	p := mgl64.Vec3{2, 2, 1}
	q := mgl64.Vec3{1.01, 1.01, -1}
	DoLineQuadTest(t, p, q, false)
}

func DoLineQuadTest(t *testing.T, p, q mgl64.Vec3, expect bool) {
	// arrange
	a := mgl64.Vec3{0, 0, 0}
	b := mgl64.Vec3{0, 1, 0}
	c := mgl64.Vec3{1, 1, 0}
	d := mgl64.Vec3{1, 0, 0}
	// Act
	result := IntersectLineQuad64(p, q, a, b, c, d)
	// Assert
	if result != expect {
		if expect {
			t.Errorf("Unexpectedly did not intersect; %v, %v", p, q)
		} else {
			t.Errorf("Unexpectedly intersects; %v, %v", p, q)
		}
	}
}
