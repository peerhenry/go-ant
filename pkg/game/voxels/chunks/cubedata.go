package chunks

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
