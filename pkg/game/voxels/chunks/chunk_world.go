package chunks

import (
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
	heightGenerator        IHeightGenerator
}

func NewChunkWorld(chunkSettings IChunkSettings) *ChunkWorld {
	chunkBuilder := NewChunkBuilder(chunkSettings)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	heightAtlas := make(map[HeightAtlasIndex]*[]int)
	perlin := ant.NewPerlin(1, 6)
	return &ChunkWorld{
		Region:                 NewChunkRegion(),
		ChunkBuilder:           chunkBuilder,
		ChunkRenderDataBuilder: &ChunkRenderDataBuilder{chunkSettings, meshBuilder},
		ChunkSettings:          chunkSettings,
		HeightAtlas:            &heightAtlas,
		HeightAtlasMapSize:     64,
		heightGenerator:        NewPerlinHeightGenerator(perlin, 200.0, 512.0),
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
		surface_ak := self.GetSurfaceAk(ai, aj) // Ak stands for absolute k
		if ak < surface_ak {
			return UNDERGROUND
		}
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
					} else if height > 20 {
						voxel = SNOWDIRT
					} else {
						voxel = GRASS
					}
				} else if dvk < 3 {
					voxel = DIRT
				}
				if dvk == pile && vk != 0 {
					// fill the rest of the voxels in the column of the chunk
					for restk := 0; restk < vk; restk++ {
						chunk.AddInvisibleVoxel(vi, vj, restk, UNDERGROUND)
						// chunk.AddVisibleVoxel(vi, vj, restk, DIRT)
					}
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
	if surfaceHeight < -4 || surfaceHeight > 20 {
		return
	}
	ak := surfaceHeight + self.ChunkSettings.GetChunkHeight()/2 // height 0 corresponds to chunk middle
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

func (self *ChunkWorld) GetSurfaceAk(ai, aj int) int {
	h := self.GetHeight(ai, aj)
	return h + self.ChunkSettings.GetChunkHeight()/2
}

// todo use this
func (self *ChunkWorld) getHeightsForChunkColumn(ci, cj int) *[]int {
	var output []int
	chunkWidth := self.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkSettings.GetChunkDepth()
	for i := 0; i < chunkWidth; i++ {
		ai := ci*chunkWidth + i
		for j := 0; j < chunkDepth; j++ {
			aj := cj*chunkDepth + j
			h := self.GetHeight(ai, aj)
			output = append(output, h)
		}
	}
	return &output
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
			h := self.heightGenerator.GetHeight(ai, aj)
			heights = append(heights, h)
		}
	}
	return &heights
}

// returns
// - voxel i & j on heightmap
// - atlas index i, j for accessing height map
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
