package chunks

import (
	"ant.com/ant/pkg/ant"
)

type ChunkWorld struct {
	Region                 *ChunkRegion
	ChunkSettings          IChunkSettings
	ChunkBuilder           *ChunkBuilder
	ChunkRenderDataBuilder *ChunkRenderDataBuilder
	initialized            bool
	HeightAtlas            *HeightAtlas
}

func NewChunkWorld(chunkSettings IChunkSettings) *ChunkWorld {
	chunkBuilder := NewChunkBuilder(chunkSettings)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	perlin := ant.NewPerlin(1, 6)
	return &ChunkWorld{
		Region:                 NewChunkRegion(),
		ChunkBuilder:           chunkBuilder,
		ChunkRenderDataBuilder: &ChunkRenderDataBuilder{chunkSettings, meshBuilder},
		ChunkSettings:          chunkSettings,
		HeightAtlas:            NewHeightAtlas(64, NewPerlinHeightGenerator(perlin, 200.0, 512.0)),
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
	ai := voxelCoordinate.i
	aj := voxelCoordinate.j
	ak := voxelCoordinate.k
	if ranks > 1 {
		// todo: work with arbitrary high ranks
		chunkCoordinate = regionCoordinate[1]
		ai = self.ChunkSettings.GetChunkWidth()*chunkCoordinate.i + voxelCoordinate.i
		aj = self.ChunkSettings.GetChunkDepth()*chunkCoordinate.j + voxelCoordinate.j
		ak = self.ChunkSettings.GetChunkHeight()*chunkCoordinate.k + voxelCoordinate.k
	} else {
		chunkCoordinate = IndexCoordinate{0, 0, 0}
	}

	chunk, ok := self.Region.Chunks[chunkCoordinate]
	if !ok {
		surface_ak := self.GetHeight(ai, aj) // Ak stands for absolute k
		if ak < surface_ak {
			return UNDERGROUND
		}
		return AIR
	}
	index := self.ChunkSettings.CoordinateToIndex(voxelCoordinate)
	return (*chunk.Voxels)[index]
}

func (self *ChunkWorld) DropTree(ai, aj, height int) {
	ak := self.GetHeight(ai, aj)
	if ak < -4 || ak > 20 {
		return
	}
	tree := GetStandardTree(height)
	for dCoord, voxel := range tree.Voxels {
		absoluteCoord := dCoord.Addijk(ai, aj, ak)
		normalized := self.ChunkSettings.NormalizeCoordinate([]IndexCoordinate{absoluteCoord})
		ranks := len(normalized)
		var chunk *StandardChunk
		if ranks > 1 {
			chunk = self.GetOrCreateChunkAt(normalized[1])
		} else {
			chunk = self.GetOrCreateChunkAt(IndexCoordinate{0, 0, 0})
		}
		voxelCoord := normalized[0]
		chunk.AddVisibleVoxel(voxelCoord.i, voxelCoord.j, voxelCoord.k, voxel)
	}
}

func (self *ChunkWorld) GetHeight(ai, aj int) int {
	h := self.HeightAtlas.GetHeight(ai, aj)
	return h + self.ChunkSettings.GetChunkHeight()/2
}

// todo use this
func (self *ChunkWorld) getHeightsForChunkColumn(ci, cj int) (*[]int, int, int) {
	var output []int
	chunkWidth := self.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkSettings.GetChunkDepth()
	min := MaxInt
	max := MinInt
	for i := 0; i < chunkWidth; i++ {
		ai := ci*chunkWidth + i
		for j := 0; j < chunkDepth; j++ {
			aj := cj*chunkDepth + j
			h := self.GetHeight(ai, aj)
			if h < min {
				min = h
			}
			if h > max {
				max = h
			}
			output = append(output, h)
		}
	}
	return &output, min, max
}

// returns voxelcoordinate k in chunk, and chunkcoordinate k in region
func (self *ChunkWorld) HeightToCoordinates(ak int) (int, int) {
	chunkHeight := self.ChunkSettings.GetChunkHeight()
	if ak >= chunkHeight {
		remainderk := ak % chunkHeight
		rankupk := ak / chunkHeight
		return remainderk, rankupk
	} else if ak < 0 {
		rankupk := (ak+1)/chunkHeight - 1
		remainderk := ak - rankupk*chunkHeight
		return remainderk, rankupk
	}
	return ak, 0

	// alternative implementation
	// remainderk := rawK % chunkHeight
	// rankupk := int(math.Floor(float64(rawK) / float64(chunkHeight)))
	// if remainderk < 0 {
	// 	remainderk = remainderk + chunkHeight
	// }
	// return remainderk, rankupk
}

func (self *ChunkWorld) CreateChunksInColumn(ci, cj int) map[IndexCoordinate]*StandardChunk {
	// get heights for chunks in column
	heights, min, _ := self.getHeightsForChunkColumn(ci, cj)
	_, minChunkK := self.HeightToCoordinates(min)
	// _, maxChunkK := self.HeightToCoordinates(max)
	chunkHeight := self.ChunkSettings.GetChunkHeight()
	// minChunkK := int(math.Floor(float64(min) / float64(chunkHeight)))

	output := make(map[IndexCoordinate]*StandardChunk)

	for index, h := range *heights {
		vi := index / self.ChunkSettings.GetChunkWidth()
		vj := index % self.ChunkSettings.GetChunkWidth()
		voxelK, topChunkK := self.HeightToCoordinates(h)
		for ck := minChunkK; ck <= topChunkK; ck++ {
			topK := chunkHeight - 1
			if ck == topChunkK {
				topK = voxelK
			}
			coord := IndexCoordinate{ci, cj, ck}
			chunk := self.GetOrCreateChunkAt(coord)
			output[coord] = chunk

			for vk := 0; vk <= topK; vk++ {
				chunk.AddVisibleVoxel(vi, vj, vk, DIRT)
				depth := topK - vk + chunkHeight*(topChunkK-ck)
				SetVoxelBasedOnHeight(chunk, vi, vj, vk, h, depth)
			}
		}

		// fill water
		if h < -6 {
			waterDepth := -6 - h
			for dk := 1; dk <= waterDepth; dk++ {
				vk, chunkK := self.HeightToCoordinates(h + dk)
				coord := IndexCoordinate{ci, cj, chunkK}
				chunk := self.GetOrCreateChunkAt(coord)
				output[coord] = chunk
				chunk.AddVisibleVoxel(vi, vj, vk, WATER)
			}
		}
	}

	return output
}

func SetVoxelBasedOnHeight(chunk *StandardChunk, vi, vj, vk, ak, depth int) {
	voxel := DIRT
	if depth == 0 {
		voxel = GRASS
	}
	if ak < -4 {
		voxel = SAND
	} else if ak > 20 {
		voxel = SNOWDIRT
	}
	chunk.AddVisibleVoxel(vi, vj, vk, voxel)
}

func (self *ChunkWorld) GetVoxelPileHeight(h, ai, aj int) int {
	pile := 0
	checkPile := func(xi, xj int) {
		d := h - self.GetHeight(xi, xj)
		if d > pile {
			pile = d
		}
	}
	checkPile(ai+1, aj)
	checkPile(ai-1, aj)
	checkPile(ai, aj-1)
	checkPile(ai, aj+1)
	return pile
}
