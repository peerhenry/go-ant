package cubes

import (
	"ant.com/ant/pkg/ant"
)

type Hud struct {
	GlslProgram  *ant.GLSLProgram
	lines        []string
	shouldUpdate bool
}

func BuildHud() *Hud {
	hud := new(Hud)
	glslProgram := ant.InitGlslProgram("shaders/hud/vertex.glsl", "shaders/hud/fragment.glsl")
	hud.GlslProgram = &glslProgram
	// todo: load texture
	return hud
}

func (self *Hud) WriteLine(text string) {
	self.lines = append(self.lines, text)
	self.shouldUpdate = true
}

func (self *Hud) Clear(text string) {
	self.lines = nil
	self.shouldUpdate = true
}

func (self *Hud) Update() {
	if self.shouldUpdate {
		// build quads
		self.shouldUpdate = false
	}
}

func (self *Hud) Draw() {
	self.GlslProgram.Use()
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

func getAtlasUv(i byte) []float32 {
	return getAtlasUvsAt(i%10, i/10)
}

func getAtlasUvForChar(char rune) []float32 {
	i := getAtlasIndex(char)
	return getAtlasUv(i)
}

func getAtlasIndex(i rune) byte {
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
