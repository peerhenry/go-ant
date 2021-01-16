package chunks

type ChunkBuilder struct {
	chunkSettings IChunkSettings
}

func NewChunkBuilder(chunkSettings IChunkSettings) *ChunkBuilder {
	return &ChunkBuilder{
		chunkSettings,
	}
}

func (self *ChunkBuilder) CreateChunk(world *ChunkWorld, ci, cj, ck int) *StandardChunk {
	var chunkVoxels []int
	var visibleVoxels []int
	chunkWidth := self.chunkSettings.GetChunkWidth()
	chunkDepth := self.chunkSettings.GetChunkDepth()
	chunkHeight := self.chunkSettings.GetChunkHeight()
	// set data of every voxel in the chunk
	for vi := 0; vi < chunkWidth; vi++ {
		for vj := 0; vj < self.chunkSettings.GetChunkDepth(); vj++ {
			for vk := 0; vk < self.chunkSettings.GetChunkHeight(); vk++ {
				voxel := self.getVoxel(vi, vj, vk)
				chunkVoxels = append(chunkVoxels, voxel)
				if vi == 0 || vi == chunkWidth-1 || vj == 0 || vj == chunkDepth-1 || vk == 0 || vk == chunkHeight-1 {
					if voxel != AIR {
						index := self.chunkSettings.CoordinateToIndexijk(vi, vj, vk)
						visibleVoxels = append(visibleVoxels, index)
					} else {
						index := self.chunkSettings.CoordinateToIndexijk(vi, vj, vk-1)
						visibleVoxels = append(visibleVoxels, index)
					}
				}
			}
		}
	}
	return &StandardChunk{
		Coordinate:    IndexCoordinate{ci, cj, ck},
		Voxels:        &chunkVoxels,
		VisibleVoxels: &visibleVoxels,
		ChunkWorld:    world,
	}
}

// todo: injecting this function
func (self *ChunkBuilder) getVoxel(i, j, k int) int {
	if k == (self.chunkSettings.GetChunkHeight() - 1) {
		// if i%2 == 0 {
		// 	return AIR
		// }
		return GRASS
	}
	if k > self.chunkSettings.GetChunkHeight()-5 {
		return DIRT
	}
	return STONE
}
