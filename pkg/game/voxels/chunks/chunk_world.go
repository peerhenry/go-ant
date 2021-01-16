package chunks

import (
	"time"

	"ant.com/ant/pkg/ant"
)

type ChunkWorld struct {
	Camera                 *ant.Camera
	Scene                  *ant.Scene
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

func (self *ChunkWorld) Update(dt *time.Duration) {
	if !self.initialized {
		for ci := -2; ci < 4; ci++ {
			for cj := -2; cj < 4; cj++ {
				chunk := self.ChunkBuilder.CreateChunk(self, ci, cj, -1)
				self.Region.SetChunkRegion(chunk)
				renderData := self.ChunkRenderDataBuilder.ChunkToRenderData(chunk)
				self.Scene.AddRenderData(renderData)
			}
		}
		self.initialized = true
	}
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
