package game

var (
	cubePositions = []float32{
		-1, -1, -1, // south
		-1, -1, 1,
		1, -1, -1,
		1, -1, 1,

		1, -1, -1, // east
		1, -1, 1,
		1, 1, -1,
		1, 1, 1,

		1, 1, -1, // north
		1, 1, 1,
		-1, 1, -1,
		-1, 1, 1,

		-1, 1, -1, // west
		-1, 1, 1,
		-1, -1, -1,
		-1, -1, 1,

		-1, -1, 1, // top
		-1, 1, 1,
		1, -1, 1,
		1, 1, 1,

		-1, 1, -1, // bottom
		-1, -1, -1,
		1, 1, -1,
		1, -1, -1,
	}
)

var (
	cubeNormals = []float32{
		0, -1, 0, // south normal -y
		0, -1, 0,
		0, -1, 0,
		0, -1, 0,
		1, 0, 0, // east normal +x
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		0, 1, 0, // north normal +y
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		-1, 0, 0, // west normal -x
		-1, 0, 0,
		-1, 0, 0,
		-1, 0, 0,
		0, 0, 1, // top normal +z
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, -1, // bottom normal -z
		0, 0, -1,
		0, 0, -1,
		0, 0, -1,
	}
)

// 36 ints
var (
	cubeIndices = []uint32{
		0, 1, 2, 2, 1, 3, // south
		4, 5, 6, 6, 5, 7, // east
		8, 9, 10, 10, 9, 11, // north
		12, 13, 14, 14, 13, 15, // west
		16, 17, 18, 18, 17, 19, // top
		20, 21, 22, 22, 21, 23, // bottom
	}
)

const pixelSize float32 = 1.0 / 512
const quadSize float32 = 1.0 / 16

func getCubeUvsAt(i, j byte) []float32 {
	left := float32(i)*quadSize + pixelSize
	right := float32(i+1)*quadSize - pixelSize
	top := float32(j)*quadSize + pixelSize
	bottom := float32(j+1)*quadSize - pixelSize
	return []float32{
		left, bottom,
		left, top,
		right, bottom,
		right, top,
	}
}

var dirt []float32 = getCubeUvsAt(2, 0)
var grassTop []float32 = getCubeUvsAt(0, 0)
var grassSide []float32 = getCubeUvsAt(3, 0)
var stone []float32 = getCubeUvsAt(1, 0)
var wood []float32 = getCubeUvsAt(4, 0)
var sand []float32 = getCubeUvsAt(2, 1)

var cubeUvsGrass []float32 = makeCubeFaces(
	grassSide,
	grassSide,
	grassSide,
	grassSide,
	grassTop,
	dirt)

var cubeUvsDirt []float32 = makeCubeFaces(
	dirt,
	dirt,
	dirt,
	dirt,
	dirt,
	dirt)

var cubeUvsStone []float32 = makeCubeFaces(
	stone,
	stone,
	stone,
	stone,
	stone,
	stone)

var cubeUvsSand []float32 = makeCubeFaces(
	sand,
	sand,
	sand,
	sand,
	sand,
	sand)

func makeCubeFaces(
	south []float32,
	east []float32,
	north []float32,
	west []float32,
	top []float32,
	bottom []float32,
) []float32 {
	thing := append(south, east...)
	thing = append(thing, north...)
	thing = append(thing, west...)
	thing = append(thing, top...)
	thing = append(thing, bottom...)
	return thing
}
