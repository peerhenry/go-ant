package text

import (
	"math"
	"time"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func createQuad(windowWidth, windowHeight int) ant.GameObject {
	var quadIndices []int32 = []int32{1, 2, 3, 4}
	builder := new(ant.VaoBuilder)
	builder.AddIntegerBuffer(0, 1, &quadIndices)
	builder.AddIndexBuffer(&[]uint32{
		0, 1, 2, 2, 1, 3,
	})
	vao := builder.Build()
	lineLength := 3
	charSize := 128
	lineHeightPixels := charSize
	lineWidthPixels := charSize * lineLength
	dimensions := mgl32.Vec2{2 * float32(lineWidthPixels) / float32(windowWidth), 2 * float32(lineHeightPixels) / float32(windowHeight)}
	marginTopPixels := 20
	marginLeftPixels := 20
	margin := mgl32.Vec2{2 * float32(marginTopPixels) / float32(windowWidth), 2 * float32(marginLeftPixels) / float32(windowHeight)}
	frameRate := 0.0
	var characters []int32 = []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	// var characters []int32 = []int32{2, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	return ant.GameObject{
		Update: func(dt *time.Duration) {
			frameRate = math.Round((1/dt.Seconds())*100) / 100
		},
		Draw: func(uniformStore *ant.UniformStore) {
			gl.BindVertexArray(vao)
			uniformStore.UniformVec2("DimensionsPixels", mgl32.Vec2{512, 512})
			uniformStore.UniformVec2("Dimensions", dimensions)
			uniformStore.UniformVec2("Margin", margin)

			uniformStore.UniformInt("AtlasWidthPixels", 512)
			uniformStore.UniformFloat("CharWidthPixels", 51.2)
			uniformStore.UniformInts("Characters[0]", characters) // todo
			uniformStore.UniformInt("QuadsPerLine", 10)
			gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		},
	}
}
