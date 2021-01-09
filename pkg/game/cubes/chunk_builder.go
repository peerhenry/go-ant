package cubes

type ChunkData struct {
	voxels        *[]int
	visibleVoxels *[]int
}

type ChunkBuilder struct {
	chunkWidth            int
	chunkDepth            int
	chunkHeight           int
	chunkDepthTimesHeight int
}

func NewChunkBuilder(chunkWidth, chunkDepth, chunkHeight int) *ChunkBuilder {
	return &ChunkBuilder{
		chunkWidth:            chunkWidth,
		chunkDepth:            chunkDepth,
		chunkHeight:           chunkHeight,
		chunkDepthTimesHeight: chunkDepth * chunkHeight,
	}
}

func (self *ChunkBuilder) CreateChunkData() *ChunkData {
	var chunkVoxels []int
	var visibleVoxels []int
	for i := 0; i < self.chunkWidth; i++ {
		for j := 0; j < self.chunkDepth; j++ {
			for k := 0; k < self.chunkHeight; k++ {
				chunkVoxels = append(chunkVoxels, self.getVoxel(i, j, k))
				if i == 0 || i == self.chunkWidth-1 || j == 0 || j == self.chunkDepth-1 || k == 0 || k == self.chunkHeight-1 {
					index := self.CoordinateToIndex(i, j, k)
					visibleVoxels = append(visibleVoxels, index)
				}
			}
		}
	}
	return &ChunkData{
		&chunkVoxels,
		&visibleVoxels,
	}
}

func (self *ChunkBuilder) CoordinateToIndex(i, j, k int) int {
	return self.chunkDepthTimesHeight*i + self.chunkHeight*j + k
}

func (self *ChunkBuilder) IndexToCoordinate(index int) [3]int {
	k := index % self.chunkHeight
	j := (index % self.chunkDepthTimesHeight) / self.chunkHeight
	i := index / self.chunkDepthTimesHeight
	return [3]int{i, j, k}
}

func (self *ChunkBuilder) getVoxel(i, j, k int) int {
	if k == (self.chunkHeight - 1) {
		return GRASS
	}
	if k < 10 {
		return STONE
	}
	return DIRT
}

const (
	NORTH  = 1
	EAST   = 2
	SOUTH  = 3
	WEST   = 4
	TOP    = 5
	BOTTOM = 6
)

// func chunkDataToFaces(chunkData *ChunkData) {
// 	// iterate over visible voxels
// 	for _, coord := range *chunkData.visibleVoxels {

// 		// inspect adjacent chunkVoxels
// 		// if transparent, add face
// 		if isTransparent()
// 	}
// }

// func isTransparent(chunkVoxels *[]int, i, j, k int) bool {
// 	index := chunkDepth*chunkWidth*i + chunkDepth*j + k
// 	return (*chunkVoxels)[index] == 0
// }

// func buildChunk() *ant.GameObject {
// 	chunkVoxels := buildChunkData()
// }
