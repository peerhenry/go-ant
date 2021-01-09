
const chunkWidth = 16
const chunkDepth = 16
const chunkHeight = 16
const GRASS = 1
const DIRT = 2

const chunkVoxels = []

for(let i = 0; i<16; i++) {
  for(let j = 0; j<16; j++) {
    for(let k = 0; k<16; k++) {
      const index = k*chunkWidth*chunkDepth + j*chunkWidth + i
      chunkVoxels[index] = k == 15 ? GRASS : DIRT
    }
  }
}

console.log(chunkVoxels)