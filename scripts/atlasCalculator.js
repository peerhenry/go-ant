const pixelSize = 1/512
const quadSize = 1/16

function getQuad(i) {
  return getQuadAt(i % 16, Math.floor(i/16))
}

function getQuadAt(i, j) {
  const left = i*quadSize + pixelSize
  const right = (i+1)*quadSize - pixelSize
  const top = j*quadSize + pixelSize
  const bottom = (j+1)*quadSize - pixelSize
  return [
    left, bottom,
    left, top,
    right, bottom,
    right, top,
  ]
}

console.log(getQuad(0))
console.log(getQuad(2))
console.log(getQuad(3))