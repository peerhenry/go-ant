package ant

import (
	"testing"

	"fmt"

	"github.com/go-gl/mathgl/mgl64"
)

func TestLineCellIntersections_xyslope(t *testing.T) {
	// arrange
	a := mgl64.Vec3{0.25, 0.75, 4.5}
	b := mgl64.Vec3{1.25, 1.75, 4.5}
	expect := []([3]int){
		[3]int{0, 0, 4},
		[3]int{0, 1, 4},
		[3]int{1, 1, 4},
	}
	// act
	result := LineCellIntersections(a, b)
	// assert
	AssertNestedArrayEqual(t, expect, result)
}

func TestLineCellIntersections_xcol(t *testing.T) {
	// arrange
	a := mgl64.Vec3{-0.5, 0.5, 4.5}
	b := mgl64.Vec3{2.5, 0.5, 4.5}
	expect := []([3]int){
		[3]int{-1, 0, 4},
		[3]int{0, 0, 4},
		[3]int{1, 0, 4},
		[3]int{2, 0, 4},
	}
	// act
	result := LineCellIntersections(a, b)
	// assert
	AssertNestedArrayEqual(t, expect, result)
}

func TestLineCellIntersections_ycol(t *testing.T) {
	// arrange
	a := mgl64.Vec3{0.5, -0.5, 4.5}
	b := mgl64.Vec3{0.5, 2.5, 4.5}
	expect := []([3]int){
		[3]int{0, -1, 4},
		[3]int{0, 0, 4},
		[3]int{0, 1, 4},
		[3]int{0, 2, 4},
	}
	// act
	result := LineCellIntersections(a, b)
	// assert
	AssertNestedArrayEqual(t, expect, result)
}

func TestLineCellIntersections_zcol(t *testing.T) {
	// arrange
	a := mgl64.Vec3{0.5, 4.5, -0.5}
	b := mgl64.Vec3{0.5, 4.5, 2.5}
	expect := []([3]int){
		[3]int{0, 4, -1},
		[3]int{0, 4, 0},
		[3]int{0, 4, 1},
		[3]int{0, 4, 2},
	}
	// act
	result := LineCellIntersections(a, b)
	// assert
	AssertNestedArrayEqual(t, expect, result)
}

func TestLineCellIntersections_zcol_negative_direction(t *testing.T) {
	// arrange
	a := mgl64.Vec3{0.5, 4.5, 2.5}
	b := mgl64.Vec3{0.5, 4.5, -0.5}
	expect := []([3]int){
		[3]int{0, 4, -1},
		[3]int{0, 4, 0},
		[3]int{0, 4, 1},
		[3]int{0, 4, 2},
	}
	// act
	result := LineCellIntersections(a, b)
	// assert
	AssertNestedArrayEqual(t, expect, result)
}

// ======= helper functions =======

func AssertNestedArrayEqual(t *testing.T, expect, actual []([3]int)) {
	if !NestedArrayEqual(expect, actual) {
		t.Errorf("assertion failed\nexpected: %s\n actual: %s", NestedArrayToString(expect), NestedArrayToString(actual))
	}
}

func NestedArrayEqual(expect, actual []([3]int)) bool {
	for index, expectedEntry := range expect {
		actualEntry := actual[index]
		for index2, expectedValue := range expectedEntry {
			actualValue := actualEntry[index2]
			if expectedValue != actualValue {
				return false
			}
		}
	}
	return true
}

func ThreeIntArrayToString(a [3]int) string {
	return fmt.Sprintf("{%d, %d, %d}", a[0], a[1], a[2])
}

func NestedArrayToString(v []([3]int)) string {
	s := "{"
	for i, next := range v {
		s += ThreeIntArrayToString(next)
		if i < len(v)-1 {
			s += ","
		}
	}
	s += "}"
	return s
}
