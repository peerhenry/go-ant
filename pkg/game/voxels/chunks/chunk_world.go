package chunks

import (
	"math"
	"math/rand"
)

type ChunkWorld struct {
	Region                 *ChunkRegion
	ChunkSettings          IChunkSettings
	ChunkBuilder           *ChunkBuilder
	ChunkRenderDataBuilder *ChunkRenderDataBuilder
	initialized            bool
	HeightAtlas            IHeightProvider
	WaterLevel             int
	SpawnTrees             bool
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

func (self *ChunkWorld) GenerateUndergroundChunk(chunkCoordinate IndexCoordinate) *StandardChunk {
	newChunk := NewChunk(self, self.Region, chunkCoordinate)
	self.Region.Chunks[chunkCoordinate] = newChunk
	newChunk.SetAllVoxels(STONE) // todo: set voxels based on depth
	return newChunk
}

func (self *ChunkWorld) DeleteChunk(chunkCoordinate IndexCoordinate) {
	delete(self.Region.Chunks, chunkCoordinate)
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
		surface_ak := self.get_surface_k(ai, aj) // Ak stands for absolute k
		if ak < surface_ak {
			return UNDERGROUND
		}
		return AIR
	}
	index := self.ChunkSettings.CoordinateToIndex(voxelCoordinate)
	return (*chunk.Voxels)[index]
}

// drops it down onto the surface, return affected chunks
func (self *ChunkWorld) DropStructure(ai, aj int, tree *VoxelStructure) map[IndexCoordinate]*StandardChunk {
	chunks := make(map[IndexCoordinate]*StandardChunk)
	ak := self.get_surface_k(ai, aj)
	if ak < -4 || ak > 20 {
		return chunks
	}
	for dCoord, voxel := range tree.Voxels {
		absoluteCoord := dCoord.Addijk(ai, aj, ak)
		normalized := self.ChunkSettings.NormalizeCoordinate([]IndexCoordinate{absoluteCoord})
		ranks := len(normalized)
		var coord IndexCoordinate
		if ranks > 1 {
			coord = normalized[1]
		} else {
			coord = IndexCoordinate{0, 0, 0}
		}
		chunk := self.GetOrCreateChunkAt(coord)
		chunks[coord] = chunk
		voxelCoord := normalized[0]
		chunk.AddVisibleVoxel(voxelCoord.i, voxelCoord.j, voxelCoord.k, voxel)
	}
	return chunks
}

func (self *ChunkWorld) DropVoxel(ai, aj, voxel Block) *StandardChunk {
	ak := self.get_surface_k(ai, aj)
	absoluteCoord := IndexCoordinate{ai, aj, ak + 1}
	normalized := self.ChunkSettings.NormalizeCoordinate([]IndexCoordinate{absoluteCoord})
	ranks := len(normalized)
	var coord IndexCoordinate
	if ranks > 1 {
		coord = normalized[1]
	} else {
		coord = IndexCoordinate{0, 0, 0}
	}
	chunk := self.GetOrCreateChunkAt(coord)
	voxelCoord := normalized[0]
	chunk.AddVisibleVoxel(voxelCoord.i, voxelCoord.j, voxelCoord.k, voxel)
	return chunk
}

func (self *ChunkWorld) get_surface_k(ai, aj int) int {
	h := self.HeightAtlas.GetHeight(ai, aj)
	return h + (self.ChunkSettings.GetChunkHeight()-1)/2
}

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
			h := self.get_surface_k(ai, aj)
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

// takes an absolute k and returns voxelcoordinate k in chunk, and chunkcoordinate k in region
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
	chunkWidth := self.ChunkSettings.GetChunkWidth()
	chunkHeight := self.ChunkSettings.GetChunkHeight()
	// get heights for chunks in column
	heights, min, _ := self.getHeightsForChunkColumn(ci, cj)
	_, minChunkK := self.HeightToCoordinates(min)
	// _, maxChunkK := self.HeightToCoordinates(max)
	// minChunkK := int(math.Floor(float64(min) / float64(chunkHeight)))

	newChunks := make(map[IndexCoordinate]*StandardChunk)

	for index, h := range *heights {
		vi := index / chunkWidth
		vj := index % chunkWidth
		voxelK, topChunkK := self.HeightToCoordinates(h)
		for ck := minChunkK; ck <= topChunkK; ck++ {
			topK := chunkHeight - 1
			if ck == topChunkK {
				topK = voxelK
			}
			coord := IndexCoordinate{ci, cj, ck}
			chunk := self.GetOrCreateChunkAt(coord)
			newChunks[coord] = chunk

			for vk := 0; vk <= topK; vk++ {
				depth := (topK - vk + chunkHeight*(topChunkK-ck))
				SetVoxelBasedOnHeight(chunk, vi, vj, vk, h, uint(depth))
			}
		}

		// fill water
		if h < self.WaterLevel {
			waterDepth := self.WaterLevel - h
			for dk := 1; dk <= waterDepth; dk++ {
				vk, chunkK := self.HeightToCoordinates(h + dk)
				coord := IndexCoordinate{ci, cj, chunkK}
				chunk := self.GetOrCreateChunkAt(coord)
				newChunks[coord] = chunk
				if dk == waterDepth {
					chunk.AddVisibleVoxel(vi, vj, vk, WATER)
				} else {
					chunk.AddInvisibleVoxel(vi, vj, vk, WATER)
				}
			}
		}
	}

	if self.SpawnTrees {
		self.dropTrees(ci, cj, newChunks)
		self.dropFlowers(ci, cj, newChunks)
	}

	return newChunks
}

func (self *ChunkWorld) dropTrees(ci, cj int, newChunks map[IndexCoordinate]*StandardChunk) {
	chunkWidth := self.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkSettings.GetChunkDepth()
	cif := float64(ci)
	cjf := float64(cj)
	p := cif * cjf
	seed := int64(2.2*math.Cos(p+78.7) + 3.3*math.Sin(p+78.7))
	rand.Seed(seed)
	trees := rand.Intn(5) // max 5 trees
	for n := 0; n < trees; n++ {
		// pick a spot
		ai := rand.Intn(chunkWidth) + chunkWidth*ci
		aj := rand.Intn(chunkDepth) + chunkDepth*cj
		extraHeight := rand.Intn(7)
		tree := GetStandardTree(6 + extraHeight)
		chunks := self.DropStructure(ai, aj, tree)
		for coord, chunk := range chunks {
			newChunks[coord] = chunk
		}
	}
}

func (self *ChunkWorld) dropFlowers(ci, cj int, newChunks map[IndexCoordinate]*StandardChunk) {
	chunkWidth := self.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkSettings.GetChunkDepth()
	cif := float64(ci)
	cjf := float64(cj)
	p := cif * cjf
	seed := int64(12.2*math.Cos(p+178.7) + 13.3*math.Sin(p+178.7))
	rand.Seed(seed)
	flowers := rand.Intn(5) // max 5 flowers
	for n := 0; n < flowers; n++ {
		// pick a spot
		ai := rand.Intn(chunkWidth) + chunkWidth*ci
		aj := rand.Intn(chunkDepth) + chunkDepth*cj
		chunk := self.DropVoxel(ai, aj, RED_FLOWER)
		newChunks[chunk.Coordinate] = chunk
	}
}

func SetVoxelBasedOnHeight(chunk *StandardChunk, vi, vj, vk, ak int, depth uint) {
	voxel := DIRT
	if depth == 0 {
		if ak > 20 {
			voxel = SNOWDIRT
		} else {
			voxel = GRASS
		}
	}
	if ak < -4 && depth < 3 {
		voxel = SAND
	}
	if depth > 8 {
		voxel = STONE
	}
	chunk.AddVisibleVoxel(vi, vj, vk, voxel)
}

func (self *ChunkWorld) GetVoxelPileHeight(h, ai, aj int) int {
	pile := 0
	checkPile := func(xi, xj int) {
		d := h - self.get_surface_k(xi, xj)
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
