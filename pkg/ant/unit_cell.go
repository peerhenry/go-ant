package ant

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// choose x, y, or z to iterate over, based on length

// x
// - y
// - z

// y
// - x
// - z

// z
// - x
// - y

// Gets the cell indices of all cubic unit cells a line intersects
func LineCellIntersections(a, b mgl64.Vec3) []([3]int) {
	d := b.Sub(a)
	if d[0] > d[1] {
		if d[1] > d[2] {
			// dx > dy > dz
			return Master(a, b, [3]int{0, 1, 2})
		} else if d[0] > d[2] {
			// dx > dz > dy
			return Master(a, b, [3]int{0, 2, 1})
		} else {
			// dz > dx > dy
			return Master(a, b, [3]int{2, 0, 1})
		}
	} else {
		if d[0] > d[2] {
			// dy > dx > dz
			return Master(a, b, [3]int{1, 0, 2})
		} else if d[1] > d[2] {
			// dy > dz > dx
			return Master(a, b, [3]int{1, 2, 1})
		} else {
			// dz > dy > dx
			return Master(a, b, [3]int{2, 1, 0})
		}
	}
}

func Master(a, b mgl64.Vec3, order [3]int) []([3]int) {
	var output []([3]int)
	a_o1 := math.Floor(a[order[0]])
	b_o1 := math.Floor(b[order[0]])
	o1_min := math.Min(a_o1, b_o1)
	o1_max := math.Max(a_o1, b_o1)
	o1_cell_min := int(o1_min)
	o1_cell_max := int(o1_max)

	for o1 := o1_cell_min; o1 <= o1_cell_max; o1++ {
		o2_min_local, o2_max_local := GetLocalMinMax(a, b, o1, o1, o1_cell_min, o1_cell_max, order, false)
		for o2 := o2_min_local; o2 <= o2_max_local; o2++ {
			o3_min_local, o3_max_local := GetLocalMinMax(a, b, o1, o2, o1_cell_min, o1_cell_max, order, true)
			if o3_max_local-o3_min_local > 100 {
				return []([3]int){[3]int{1, 2, 3}}
			}
			for o3 := o3_min_local; o3 <= o3_max_local; o3++ {
				output = append(output, MakeOrder(order, o1, o2, o3))
			}
		}
	}
	return output
}

func GetLocalMinMax(a, b mgl64.Vec3, o1, o_n1, o1_cell_min, o1_cell_max int, order [3]int, forO3 bool) (int, int) {

	n1 := 0
	n2 := 1
	if forO3 {
		if b[order[1]] != a[order[1]] {
			// use third dimension line as function of second
			n1 = 1
			n2 = 2
		} else {
			// use third dimension line as function of first
			n1 = 0
			n2 = 2
		}
	}

	slope := (b[order[n2]] - a[order[n2]]) / (b[order[n1]] - a[order[n1]])
	ff := func(t1 float64) float64 {
		return a[order[n2]] + slope*(t1-a[order[n1]])
	}

	a_o_n2 := math.Floor(a[order[n2]])
	b_o_n2 := math.Floor(b[order[n2]])
	o_n2_min := math.Min(a_o_n2, b_o_n2)
	o_n2_max := math.Max(a_o_n2, b_o_n2)

	var p_o_n2 float64
	if o1 == o1_cell_min {
		p_o_n2 = o_n2_min
	} else {
		p_o_n2 = ff(float64(o_n1))
	}
	var q_o_n2 float64
	if o1 == o1_cell_max {
		q_o_n2 = o_n2_max
	} else {
		q_o_n2 = ff(float64(o_n1 + 1))
	}

	o_n2_min_local := math.Min(p_o_n2, q_o_n2)
	o_n2_max_local := math.Max(p_o_n2, q_o_n2)
	return int(o_n2_min_local), int(o_n2_max_local)
}

func MakeOrder(order [3]int, o1, o2, o3 int) [3]int {
	var output [3]int
	output[order[0]] = o1
	output[order[1]] = o2
	output[order[2]] = o3
	return output
}
