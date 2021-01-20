package chunks

import (
	"math"

	"ant.com/ant/pkg/ant"
)

type PerlinHeightGenerator struct {
	perlin    *ant.Perlin
	amplitude float64
	cellSize  float64
}

func NewPerlinHeightGenerator(perlin *ant.Perlin, amplitude float64, cellSize float64) *PerlinHeightGenerator {
	return &PerlinHeightGenerator{perlin: perlin, amplitude: amplitude, cellSize: cellSize}
}

func (self *PerlinHeightGenerator) GetHeight(ai, aj int) int {
	return int(math.Round(self.amplitude * self.perlin.Noise(float64(ai)/self.cellSize, float64(aj)/self.cellSize)))
}

type FlatHeightGenerator struct {
	height int
}

func (self *FlatHeightGenerator) GetHeight(ai, aj int) int {
	return self.height
}
