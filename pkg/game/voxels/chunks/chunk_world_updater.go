package chunks

import (
	"time"

	"ant.com/ant/pkg/ant"
)

type ChunkWorldUpdater struct {
	Camera      *ant.Camera
	Scene       *ant.Scene
	initialized bool
	ChunkWorld  *ChunkWorld
}

func NewChunkWorldUpdater(camera *ant.Camera, scene *ant.Scene) *ChunkWorldUpdater {
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings)
	return &ChunkWorldUpdater{
		Camera:     camera,
		Scene:      scene,
		ChunkWorld: world,
	}
}

func (self *ChunkWorldUpdater) Update(dt *time.Duration) {
	if !self.initialized {
		for ci := -2; ci < 4; ci++ {
			for cj := -2; cj < 4; cj++ {
				chunk := self.ChunkWorld.ChunkBuilder.CreateChunk(self.ChunkWorld, ci, cj, -1)
				self.ChunkWorld.Region.SetChunkRegion(chunk)
			}
		}

		for _, chunk := range self.ChunkWorld.Region.Chunks {
			renderData := self.ChunkWorld.ChunkRenderDataBuilder.ChunkToRenderData(chunk)
			self.Scene.AddRenderData(renderData)
		}

		self.initialized = true
	}
}
