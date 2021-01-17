package chunks

import (
	"log"
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
		startChunks := time.Now()
		for ci := -3; ci < 4; ci++ {
			for cj := -3; cj < 4; cj++ {
				self.ChunkWorld.CreateChunksInColumn(ci, cj)
				// chunk := self.ChunkWorld.ChunkBuilder.CreateChunk(self.ChunkWorld, ci, cj, -1)
				// self.ChunkWorld.Region.SetChunkRegion(chunk)
			}
		}
		elapsedChunks := time.Since(startChunks)
		log.Printf("Creating chunks took %s", elapsedChunks)

		startRenderData := time.Now()
		for _, chunk := range self.ChunkWorld.Region.Chunks {
			renderData := self.ChunkWorld.ChunkRenderDataBuilder.ChunkToRenderData(chunk)
			self.Scene.AddRenderData(renderData)
		}
		elapsed := time.Since(startRenderData)
		log.Printf("Converting chunks to render data took %s", elapsed)
		self.initialized = true
	}
}
