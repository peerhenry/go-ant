package chunks

type ChunkWorld struct {
	Region                 *ChunkRegion
	ChunkSettings          IChunkSettings
	ChunkBuilder           *ChunkBuilder
	ChunkRenderDataBuilder *ChunkRenderDataBuilder
	initialized            bool
}

func NewChunkWorld(chunkSettings IChunkSettings) *ChunkWorld {
	chunkBuilder := NewChunkBuilder(chunkSettings)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	return &ChunkWorld{
		Region:                 NewChunkRegion(),
		ChunkBuilder:           chunkBuilder,
		ChunkRenderDataBuilder: &ChunkRenderDataBuilder{chunkSettings, meshBuilder},
		ChunkSettings:          chunkSettings,
	}
}

func (self *ChunkWorld) GetOrCreateChunkAt(chunkCoordinate IndexCoordinate) *StandardChunk {
	chunk, ok := self.Region.Chunks[chunkCoordinate]
	if !ok {
		newChunk := NewChunk(self, self.Region, chunkCoordinate)
		self.Region.Chunks[chunkCoordinate] = newChunk
		return newChunk
	}
	return chunk
}

func (self *ChunkWorld) GetVoxelAt(regionCoordinate []IndexCoordinate) int {
	ranks := len(regionCoordinate)

	voxelCoordinate := regionCoordinate[0]
	var chunkCoordinate IndexCoordinate
	if ranks > 1 {
		// todo: work with arbitrary high ranks
		chunkCoordinate = regionCoordinate[1]
	} else {
		chunkCoordinate = IndexCoordinate{0, 0, 0}
	}

	chunk, ok := self.Region.Chunks[chunkCoordinate]
	if !ok {
		return AIR
	}
	index := self.ChunkSettings.CoordinateToIndex(voxelCoordinate)
	return (*chunk.Voxels)[index]
}

func (self *ChunkWorld) CreateChunksInColumn(ci, cj int) {
	heights, _, _ := self.GetChunkColumnHeights(ci, cj)
	chunkWidth := self.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkSettings.GetChunkDepth()
	for vi := 0; vi < chunkWidth; vi++ {
		for vj := 0; vj < chunkDepth; vj++ {
			height := (*heights)[vi*chunkWidth+vj]
			voxelK, chunkK := self.HeightToCoordinates(height)
			chunk := self.GetOrCreateChunkAt(IndexCoordinate{ci, cj, chunkK})
			chunk.AddVisibleVoxel(vi, vj, voxelK, STONE)
		}
	}
}

func (self *ChunkWorld) GetChunkColumnHeights(ci, cj int) (*[]int, int, int) {
	chunkWidth := self.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkSettings.GetChunkDepth()
	var heights []int
	minHeight := MaxInt
	maxHeight := MinInt
	for vi := 0; vi < chunkWidth; vi++ {
		for vj := 0; vj < chunkDepth; vj++ {
			ii := chunkWidth*ci + vi
			jj := chunkWidth*cj + vj
			h := 0
			if ii > 0 && jj > 0 {
				h = (ii + jj) / 10
			}
			heights = append(heights, 0)
			if h > maxHeight {
				maxHeight = h
			}
			if h < minHeight {
				minHeight = h
			}
		}
	}
	return &heights, minHeight, maxHeight
}

// returns voxelcoordinate k in chunk, and chunkcoordinate k in region
func (self *ChunkWorld) HeightToCoordinates(h int) (int, int) {
	chunkHeight := self.ChunkSettings.GetChunkHeight()
	base := chunkHeight / 2
	rawK := base + h

	// remainderk := rawK % chunkHeight
	// rankupk := int(math.Floor(float64(rawK) / float64(chunkHeight)))
	// if remainderk < 0 {
	// 	remainderk = remainderk + chunkHeight
	// }
	// return remainderk, rankupk

	if rawK >= chunkHeight {
		remainderk := rawK % chunkHeight
		rankupk := rawK / chunkHeight
		return remainderk, rankupk
	} else if rawK < 0 {
		rankupk := (rawK+1)/chunkHeight - 1
		remainderk := rawK - rankupk*chunkHeight
		return remainderk, rankupk
	}
	return rawK, 0
}
