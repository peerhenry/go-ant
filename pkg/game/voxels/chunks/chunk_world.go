package chunks

import (
	"math"

	"ant.com/ant/pkg/ant"
)

type HeightAtlasIndex struct {
	i int
	j int
}

type ChunkWorld struct {
	Region                 *ChunkRegion
	ChunkSettings          IChunkSettings
	ChunkBuilder           *ChunkBuilder
	ChunkRenderDataBuilder *ChunkRenderDataBuilder
	initialized            bool
	HeightAtlas            *map[HeightAtlasIndex]*[]int
	HeightAtlasMapSize     int
	Perlin                 *ant.Perlin
}

func NewChunkWorld(chunkSettings IChunkSettings) *ChunkWorld {
	chunkBuilder := NewChunkBuilder(chunkSettings)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	heightAtlas := make(map[HeightAtlasIndex]*[]int)
	return &ChunkWorld{
		Region:                 NewChunkRegion(),
		ChunkBuilder:           chunkBuilder,
		ChunkRenderDataBuilder: &ChunkRenderDataBuilder{chunkSettings, meshBuilder},
		ChunkSettings:          chunkSettings,
		HeightAtlas:            &heightAtlas,
		HeightAtlasMapSize:     64,
		Perlin:                 ant.NewPerlin(1, 1),
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
	chunkWidth := self.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkSettings.GetChunkDepth()
	for vi := 0; vi < chunkWidth; vi++ {
		for vj := 0; vj < chunkDepth; vj++ {
			ai := ci*chunkWidth + vi
			aj := cj*chunkDepth + vj
			height, pile := self.GetPileCount(ai, aj) // how high should this voxel column add voxels

			for dvk := 0; dvk <= pile; dvk++ {
				vk, chunkK := self.HeightToCoordinates(height - dvk)
				chunk := self.GetOrCreateChunkAt(IndexCoordinate{ci, cj, chunkK})
				voxel := STONE
				if dvk == 0 {
					if height < -4 {
						voxel = SAND
					} else if height > 8 {
						voxel = SNOWDIRT
					} else {
						voxel = GRASS
					}
				} else if dvk < 3 {
					voxel = DIRT
				}
				chunk.AddVisibleVoxel(vi, vj, vk, voxel)
			}
			// fill water
			if height < -6 {
				waterDepth := -6 - height
				for dk := 1; dk <= waterDepth; dk++ {
					vk, chunkK := self.HeightToCoordinates(height + dk)
					chunk := self.GetOrCreateChunkAt(IndexCoordinate{ci, cj, chunkK})
					chunk.AddVisibleVoxel(vi, vj, vk, WATER)
				}
			}
		}
	}
}

func (self *ChunkWorld) DropTree(ai, aj, height int) {
	surfaceHeight := self.GetHeight(ai, aj)
	ak := surfaceHeight + self.ChunkSettings.GetChunkHeight()/2 // height 0 corresponds to chunk middle
	if surfaceHeight < -4 {
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

func (self *ChunkWorld) GetPileCount(ai, aj int) (int, int) {
	pile := 0
	h := self.GetHeight(ai, aj)
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
	return h, pile
}

func (self *ChunkWorld) GetHeight(ai, aj int) int {
	heightMapSize := self.HeightAtlasMapSize
	vi, vj, hmi, hmj := self.HorizontalToHeightMapCoordinates(ai, aj)
	localHeightMap, ok := (*self.HeightAtlas)[HeightAtlasIndex{hmi, hmj}]
	if !ok {
		// generate height map for atlas
		newHeightMap := self.GetHeightMap(hmi, hmj)
		(*self.HeightAtlas)[HeightAtlasIndex{hmi, hmj}] = newHeightMap
		return (*newHeightMap)[vi*heightMapSize+vj]
	}
	return (*localHeightMap)[vi*heightMapSize+vj]
}

func (self *ChunkWorld) GetHeightMap(hmi, hmj int) *[]int {
	heightMapSize := self.HeightAtlasMapSize
	var heights []int
	for vi := 0; vi < heightMapSize; vi++ {
		for vj := 0; vj < heightMapSize; vj++ {
			// absolute voxel i & j in world
			ai := heightMapSize*hmi + vi
			aj := heightMapSize*hmj + vj
			h := 0

			// === one big slope ===
			// h := ai + aj
			// ====================

			// === random ===
			// rand.Seed(time.Now().UnixNano() + int64(vi*vj))
			// h := rand.Intn(3)
			// ====================

			// === perlin noise ===
			cellSize := 128.0
			amp := 80.0
			h = int(math.Round(amp * self.Perlin.Perlin(float64(ai)/cellSize, float64(aj)/cellSize)))
			// ====================

			heights = append(heights, h)
		}
	}
	return &heights
}

func (self *ChunkWorld) HorizontalToHeightMapCoordinates(i int, j int) (int, int, int, int) {
	size := self.HeightAtlasMapSize
	vi := i
	vj := j
	hmi := 0
	hmj := 0
	if vi >= size {
		vi = i % size
		hmi = i / size
	} else if vi < 0 {
		hmi = (vi+1)/size - 1
		vi = vi - hmi*size
	}
	if vj >= size {
		vj = j % size
		hmj = j / size
	} else if vj < 0 {
		hmj = (vj+1)/size - 1
		vj = vj - hmj*size
	}
	return vi, vj, hmi, hmj
}

// returns voxelcoordinate k in chunk, and chunkcoordinate k in region
func (self *ChunkWorld) HeightToCoordinates(h int) (int, int) {
	chunkHeight := self.ChunkSettings.GetChunkHeight()
	base := chunkHeight / 2
	rawK := base + h
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

	// alternative implementation
	// remainderk := rawK % chunkHeight
	// rankupk := int(math.Floor(float64(rawK) / float64(chunkHeight)))
	// if remainderk < 0 {
	// 	remainderk = remainderk + chunkHeight
	// }
	// return remainderk, rankupk
}
