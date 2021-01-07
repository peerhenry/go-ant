package text

import (
	"fmt"
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
	lineLength := 10
	charAspectRatio := float64(44) / float64(108)
	charHeightPixels := 44
	charWidthPixels := int(charAspectRatio * float64(charHeightPixels))
	lineHeightPixels := charHeightPixels
	lineWidthPixels := charWidthPixels * lineLength
	lineWidth := 2 * float32(lineWidthPixels) / float32(windowWidth)
	lineHeight := 2 * float32(lineHeightPixels) / float32(windowHeight)
	dimensions := mgl32.Vec2{lineWidth, lineHeight}
	marginTopPixels := 40
	marginLeftPixels := 40
	marginTop := 2 * float32(marginTopPixels) / float32(windowHeight)
	marginLeft := 2 * float32(marginLeftPixels) / float32(windowWidth)
	margin := mgl32.Vec2{marginLeft, marginTop}
	frameRate := 0.0
	var characters []int32 = stringToCharacterCodes("FPS: ?")
	var step float64 = 0
	return ant.GameObject{
		Update: func(dt *time.Duration) {
			step += dt.Seconds()
			if step >= 1 {
				frameRate = math.Round((1/dt.Seconds())*100) / 100
				frameRateString := fmt.Sprintf("%f", frameRate)
				characters = stringToCharacterCodes("FPS: " + frameRateString)
				step -= 1
			}
		},
		Draw: func(uniformStore *ant.UniformStore) {
			gl.BindVertexArray(vao)
			uniformStore.UniformVec2("DimensionsPixels", mgl32.Vec2{float32(lineWidthPixels), float32(lineHeightPixels)})
			uniformStore.UniformVec2("Dimensions", dimensions)
			uniformStore.UniformVec2("Margin", margin)

			uniformStore.UniformFloat("HalfPixel", 1.0/1024)
			uniformStore.UniformFloat("CharWidthPixels", float32(charWidthPixels))
			uniformStore.UniformFloat("CharHeightPixels", float32(charHeightPixels))
			uniformStore.UniformInts("Characters[0]", characters) // todo
			uniformStore.UniformInt("QuadsPerLine", 10)
			gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		},
	}
}
