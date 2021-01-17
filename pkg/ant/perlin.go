package ant

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type Perlin struct {
	Seed          int64
	octaves       int
	octaveWeights []float64
	octaveScales  []float64
	gradientVecs  []mgl64.Vec2
}

func NewPerlin(seed int64, octaves int) *Perlin {
	var octaveWeights []float64
	var octaveScales []float64
	for octave := 0; octave < octaves; octave++ {
		halfPow := math.Pow(0.5, float64(octave))
		log.Println("adding halfPow", halfPow)
		octaveWeights = append(octaveWeights, halfPow)
		octaveScales = append(octaveScales, halfPow)
	}
	var gradientVecs []mgl64.Vec2
	gradientVecRadius := math.Sqrt(2)
	for n := 0; n < 32; n++ {
		angle := float64(n) * math.Pi / 16.0
		next := mgl64.Vec2{gradientVecRadius * math.Cos(angle), gradientVecRadius * math.Sin(angle)}
		gradientVecs = append(gradientVecs, next)
	}
	return &Perlin{
		Seed:          seed,
		octaves:       octaves,
		octaveWeights: octaveWeights,
		octaveScales:  octaveScales,
		gradientVecs:  gradientVecs,
	}
}

func (self *Perlin) ConfigureOctaves(octaveWeights, octaveScales []float64) {
	if len(octaveWeights) != len(octaveScales) {
		errorMsg := fmt.Sprintf("octaveWeights and octaveScales must be the samen length, but were respectively %d and %d", len(octaveWeights), len(octaveScales))
		panic(errorMsg)
	}
	self.octaves = len(octaveWeights)
}

func (self *Perlin) Perlin(x, y float64) float64 {
	total := 0.0
	weightTotal := 0.0
	for octave := 0; octave < self.octaves; octave++ {
		next := self.PerlinForOctave(x, y, octave)
		weight := self.octaveWeights[octave]
		weightTotal += weight
		total = next * weight
	}
	return total / weightTotal
}

// A - B
// | . |
// C - D

func (self *Perlin) PerlinForOctave(x, y float64, octave int) float64 {
	getDotProduct := func(celli, cellj float64) float64 {
		rand.Seed(self.Seed * int64(celli) * int64(cellj))
		hash := rand.Intn(len(self.gradientVecs))
		gradient := self.gradientVecs[hash]
		dx := x - celli
		dy := y - cellj
		return gradient[0]*dx + gradient[1]*dy
	}

	scale := self.octaveScales[octave]
	sx := scale * x
	sy := scale * y
	originx := math.Floor(sx)
	originy := math.Floor(sy)
	dotA := getDotProduct(originx, originy)
	dotB := getDotProduct(originx+1, originy)
	dotC := getDotProduct(originx, originy+1)
	dotD := getDotProduct(originx+1, originy+1)

	relX := x - originx
	mix1 := mix(dotA, dotB, SmoothStep(relX))
	mix2 := mix(dotC, dotD, SmoothStep(relX))
	relY := y - originy
	return mix(mix1, mix2, SmoothStep(relY))
}

func SmoothStep(f float64) float64 {
	return f * f * (3.0 - 2.0*f)
}

func mix(x, y, a float64) float64 {
	return x*(1-a) + y*a
}

// func AltSmooth(t float64) float64 {
// 	return t * t * t * (t*(t*6-15) + 10)
// }
