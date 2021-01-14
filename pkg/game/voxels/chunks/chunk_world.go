package chunks

import (
	"time"

	"ant.com/ant/pkg/ant"
)

type ChunkWorld struct {
	Camera                 *ant.Camera
	Region                 *ChunkRegion
	Scene                  *ant.Scene
	ChunkBuilder           *ChunkBuilder
	ChunkRenderDataBuilder *ChunkRenderDataBuilder
	initialized            bool
}

func NewChunkWorld(camera *ant.Camera, region *ChunkRegion, scene *ant.Scene) *ChunkWorld {
	chunkSettings := CreateStandardChunkSettings(32, 32, 8)
	chunkBuilder := CreateStandardChunkBuilder(chunkSettings)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	return &ChunkWorld{
		Camera:                 camera,
		Region:                 region,
		ChunkBuilder:           chunkBuilder,
		ChunkRenderDataBuilder: &ChunkRenderDataBuilder{chunkSettings, meshBuilder},
		Scene: scene,
	}
}

func (self *ChunkWorld) Update(dt *time.Duration) {
	if !self.initialized {
		for ci := -2; ci < 4; ci++ {
			for cj := -2; cj < 4; cj++ {
				chunk := self.ChunkBuilder.CreateChunk(ci, cj, -1)
				renderData := self.ChunkRenderDataBuilder.ChunkToRenderData(chunk)
				self.Region.SetChunkRegion(chunk)
				self.Scene.AddRenderData(renderData)
			}
		}
		self.initialized = true
	}
}
