package ant

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type PerlinCacheKey struct {
	x float64
	y float64
}

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
		twoPow := math.Pow(2, float64(octave))
		octaveWeights = append(octaveWeights, halfPow)
		octaveScales = append(octaveScales, twoPow)
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

func (self *Perlin) Noise(x, y float64) float64 {
	total := 0.0
	weightTotal := 0.0
	for octave := 0; octave < self.octaves; octave++ {
		scale := self.octaveScales[octave]
		next := self.Perlin(scale*x, scale*y)
		weight := self.octaveWeights[octave]
		weightTotal += weight
		total = next * weight
	}
	return total / weightTotal
}

// A - B
// | . |
// C - D

func (self *Perlin) Perlin(x, y float64) float64 {
	getDotProduct := func(celli, cellj, dx, dy float64) float64 {
		randomAngle := 2920.0 * math.Sin(celli*21942.0+cellj*171324.0+8912.0) * math.Cos(celli*23157.0*cellj*217832.0+9758.0)
		return math.Cos(randomAngle)*dx + math.Sin(randomAngle)*dy
	}

	left := math.Floor(x)
	top := math.Floor(y)
	right := left + 1
	bottom := top + 1
	xFromLeft := x - left
	xFromRight := x - right
	yFromTop := y - top
	yFromBottom := y - bottom

	dotA := getDotProduct(left, top, xFromLeft, yFromTop)
	dotB := getDotProduct(right, top, xFromRight, yFromTop)
	dotC := getDotProduct(left, bottom, xFromLeft, yFromBottom)
	dotD := getDotProduct(right, bottom, xFromRight, yFromBottom)

	f := SmoothStep(xFromLeft)
	mix1 := mix(dotA, dotB, f)
	mix2 := mix(dotC, dotD, f)
	value := mix(mix1, mix2, SmoothStep(yFromTop))
	return value
}

func SmoothStep(f float64) float64 {
	return f * f * (3.0 - 2.0*f)
}

func mix(x, y, a float64) float64 {
	return x*(1-a) + y*a
}

func Fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}
