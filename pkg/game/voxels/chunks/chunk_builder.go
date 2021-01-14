package chunks

type IChunkBuilder interface {
	CreateChunkData() StandardChunk
}

type ChunkBuilder struct {
	chunkSettings IChunkSettings
}

func CreateStandardChunkBuilder(chunkSettings IChunkSettings) *ChunkBuilder {
	return &ChunkBuilder{
		chunkSettings,
	}
}

func (self *ChunkBuilder) CreateChunk(ci, cj, ck int) *StandardChunk {
	var chunkVoxels []int
	var visibleVoxels []int
	chunkWidth := self.chunkSettings.GetChunkWidth()
	chunkDepth := self.chunkSettings.GetChunkDepth()
	chunkHeight := self.chunkSettings.GetChunkHeight()
	for vi := 0; vi < chunkWidth; vi++ {
		for vj := 0; vj < self.chunkSettings.GetChunkDepth(); vj++ {
			for vk := 0; vk < self.chunkSettings.GetChunkHeight(); vk++ {
				chunkVoxels = append(chunkVoxels, self.getVoxel(vi, vj, vk))
				if vi == 0 || vi == chunkWidth-1 || vj == 0 || vj == chunkDepth-1 || vk == 0 || vk == chunkHeight-1 {
					index := self.chunkSettings.CoordinateToIndex(vi, vj, vk)
					visibleVoxels = append(visibleVoxels, index)
				}
			}
		}
	}
	return &StandardChunk{
		Index:         ChunkIndex{ci, cj, ck},
		voxels:        &chunkVoxels,
		visibleVoxels: &visibleVoxels,
		chunkSettings: self.chunkSettings,
	}
}

// todo: injecting this function
func (self *ChunkBuilder) getVoxel(i, j, k int) int {
	if k == (self.chunkSettings.GetChunkHeight() - 1) {
		return GRASS
	}
	if k > self.chunkSettings.GetChunkHeight()-5 {
		return DIRT
	}
	return STONE
}
