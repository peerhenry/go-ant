package performance_tests

import (
	"testing"

	"ant.com/ant/pkg/ant"
)

func TestNoise(t *testing.T) {
	perlin := ant.NewPerlin(12345678, 3)
	min := 10000.0
	max := -10000.0
	for i := 0; i < 200; i++ {
		for j := 0; j < 200; j++ {
			h := perlin.Noise(float64(i)/100.0, float64(j)/100.0)
			if h < min {
				min = h
			}
			if h > max {
				max = h
			}
		}
	}
	t.Logf("min, max: %f, %f", min, max)
}

func TestPerlin(t *testing.T) {
	perlin := ant.NewPerlin(666, 1)
	h := perlin.Perlin(25.4, 74.3)
	if h == 0.0 {
		t.Errorf("Expected h not to be zero")
	}
}
