package voxels

import (
	"fmt"
	"math"
	"time"

	"ant.com/ant/pkg/ant"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Hud struct {
	Draw   func()
	Update func(dt *time.Duration)
}

func BuildHud(windowWidth, windowHeight int) *Hud {
	hud := new(Hud)
	glslProgram := ant.InitGlslProgram("shaders/text/vertex.glsl", "shaders/text/fragment.glsl")
	uniformStore := ant.CreateUniformStore(glslProgram.Handle, true)
	ant.LoadImageFileToUniform("resources/text-atlas-monospace-white-on-alpha.png", "TextAtlas", glslProgram.Handle, 1)

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
	hud.Update = func(dt *time.Duration) {
		step += dt.Seconds()
		if step >= 1 {
			frameRate = math.Round((1/dt.Seconds())*100) / 100
			frameRateString := fmt.Sprintf("%f", frameRate)
			characters = stringToCharacterCodes("FPS: " + frameRateString)
			step -= 1
		}
	}
	hud.Draw = func() {
		glslProgram.Use()
		gl.Disable(gl.CULL_FACE)
		gl.Disable(gl.DEPTH_TEST)

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
	}
	return hud
}

// todo: unify with the one from cubedata

const dimension = 10 // characters per row/column
const imageWidth = 512
const fontAtlasPixelSize float32 = 1.0 / imageWidth
const fontAtlasQuadSize float32 = 1.0 / dimension

func getAtlasUvsAt(i, j byte) []float32 {
	left := float32(i)*fontAtlasQuadSize + fontAtlasPixelSize
	right := float32(i+1)*fontAtlasQuadSize - fontAtlasPixelSize
	top := float32(j)*fontAtlasQuadSize + fontAtlasPixelSize
	bottom := float32(j+1)*fontAtlasQuadSize - fontAtlasPixelSize
	return []float32{
		left, bottom,
		left, top,
		right, bottom,
		right, top,
	}
}

func stringToCharacterCodes(line string) []int32 {
	var chars []int32
	for _, c := range line {
		chars = append(chars, getAtlasIndex(c))
	}
	return chars
}

func getAtlasIndex(i rune) int32 {
	switch i {
	case ' ':
		return 0
	case '!':
		return 1
	case '"':
		return 2
	case '#':
		return 3
	case '$':
		return 4
	case '%':
		return 5
	case '&':
		return 6
	case '\'':
		return 7
	case '(':
		return 8
	case ')':
		return 9
	case '*':
		return 10
	case '+':
		return 11
	case ',':
		return 12
	case '-':
		return 13
	case '.':
		return 14
	case '/':
		return 15
	case '0':
		return 16
	case '1':
		return 17
	case '2':
		return 18
	case '3':
		return 19
	case '4':
		return 20
	case '5':
		return 21
	case '6':
		return 22
	case '7':
		return 23
	case '8':
		return 24
	case '9':
		return 25
	case ':':
		return 26
	case ';':
		return 27
	case '<':
		return 28
	case '=':
		return 29
	case '>':
		return 30
	case '?':
		return 31
	case '@':
		return 32
	case 'A':
		return 33
	case 'B':
		return 34
	case 'C':
		return 35
	case 'D':
		return 36
	case 'E':
		return 37
	case 'F':
		return 38
	case 'G':
		return 39
	case 'H':
		return 40
	case 'I':
		return 41
	case 'J':
		return 42
	case 'K':
		return 43
	case 'L':
		return 44
	case 'M':
		return 45
	case 'N':
		return 46
	case 'O':
		return 47
	case 'P':
		return 48
	case 'Q':
		return 49
	case 'R':
		return 50
	case 'S':
		return 51
	case 'T':
		return 52
	case 'U':
		return 53
	case 'V':
		return 54
	case 'W':
		return 55
	case 'X':
		return 56
	case 'Y':
		return 57
	case 'Z':
		return 58
	case '[':
		return 59
	case '\\':
		return 60
	case ']':
		return 61
	case '^':
		return 62
	case '_':
		return 63
	case '`':
		return 64
	case 'a':
		return 65
	case 'b':
		return 66
	case 'c':
		return 67
	case 'd':
		return 68
	case 'e':
		return 69
	case 'f':
		return 70
	case 'g':
		return 71
	case 'h':
		return 72
	case 'i':
		return 73
	case 'j':
		return 74
	case 'k':
		return 75
	case 'l':
		return 76
	case 'm':
		return 77
	case 'n':
		return 78
	case 'o':
		return 79
	case 'p':
		return 80
	case 'q':
		return 81
	case 'r':
		return 82
	case 's':
		return 83
	case 't':
		return 84
	case 'u':
		return 85
	case 'v':
		return 86
	case 'w':
		return 87
	case 'x':
		return 88
	case 'y':
		return 89
	case 'z':
		return 90
	case '{':
		return 91
	case '|':
		return 92
	case '}':
		return 93
	default:
		return 94
	}
}
