package chunks

import (
	"log"
	"math/rand"
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
		for ci := -6; ci < 7; ci++ {
			for cj := -6; cj < 7; cj++ {
				self.ChunkWorld.CreateChunksInColumn(ci, cj)
				// chunk := self.ChunkWorld.ChunkBuilder.CreateChunk(self.ChunkWorld, ci, cj, -1)
				// self.ChunkWorld.Region.SetChunkRegion(chunk)
			}
		}
		// drop some trees
		for ci := -18; ci < 18; ci++ {
			for cj := -18; cj < 18; cj++ {
				rand.Seed(time.Now().UnixNano() + int64(ci*cj)) // todo world seed
				extra := rand.Intn(7)
				ddi := rand.Intn(7)
				ddj := rand.Intn(8)
				self.ChunkWorld.DropTree(ci*10+ddi, cj*10+ddj, 6+extra)
			}
		}
		elapsedChunks := time.Since(startChunks)
		log.Printf("Creating chunks took %s", elapsedChunks)

		startRenderData := time.Now()
		for _, chunk := range self.ChunkWorld.Region.Chunks {
			if chunk.IsVisible() {
				renderData := self.ChunkWorld.ChunkRenderDataBuilder.ChunkToRenderData(chunk)
				if renderData != nil {
					self.Scene.AddRenderData(renderData)
				}
			}
		}
		elapsed := time.Since(startRenderData)
		log.Printf("Converting chunks to render data took %s", elapsed)
		self.initialized = true
	}
}
