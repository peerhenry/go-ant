package chunks

import (
	"log"
	"math"
	"time"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl64"
)

type ChunkColumnCoordinate [2]int

type ChunkWorldUpdater struct {
	Camera          *ant.Camera
	CameraAnchor    mgl64.Vec3
	Scene           *ant.Scene
	shouldSpawnMore bool
	ChunkWorld      *ChunkWorld
	spawnedColumns  map[ChunkColumnCoordinate]bool
}

func NewChunkWorldUpdater(camera *ant.Camera, scene *ant.Scene) *ChunkWorldUpdater {
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings)
	return &ChunkWorldUpdater{
		Camera:          camera,
		CameraAnchor:    camera.Position,
		Scene:           scene,
		ChunkWorld:      world,
		shouldSpawnMore: true,
		spawnedColumns:  make(map[ChunkColumnCoordinate]bool),
	}
}

func (self *ChunkWorldUpdater) Update(dt *time.Duration) {
	delta := self.Camera.Position.Sub(self.CameraAnchor)
	if delta.Len() > 50 {
		self.CameraAnchor = self.Camera.Position
		self.shouldSpawnMore = true
	}
	if self.shouldSpawnMore {
		chunkWidth := self.ChunkWorld.ChunkSettings.GetChunkWidth()
		chunkDepth := self.ChunkWorld.ChunkSettings.GetChunkDepth()
		startChunks := time.Now()
		cam_ci := int(math.Floor(self.Camera.Position[0] / float64(chunkWidth)))
		cam_cj := int(math.Floor(self.Camera.Position[1] / float64(chunkDepth)))
		ci_min := cam_ci - 4
		ci_max := cam_ci + 4
		cj_min := cam_cj - 4
		cj_max := cam_cj + 4

		newChunks := make(map[IndexCoordinate]*StandardChunk)

		for ci := ci_min; ci <= ci_max; ci++ {
			for cj := cj_min; cj < cj_max; cj++ {
				_, ok := self.spawnedColumns[[2]int{ci, cj}]
				if !ok {
					newChunksInColumns := self.ChunkWorld.CreateChunksInColumn(ci, cj)
					for coord, chunk := range newChunksInColumns {
						newChunks[coord] = chunk
					}
					self.spawnedColumns[[2]int{ci, cj}] = true
				}
			}
		}
		// drop some trees
		// for ci := -18; ci < 18; ci++ {
		// 	for cj := -18; cj < 18; cj++ {
		// 		rand.Seed(time.Now().UnixNano() + int64(ci*cj)) // todo world seed
		// 		extra := rand.Intn(7)
		// 		ddi := rand.Intn(7)
		// 		ddj := rand.Intn(8)
		// 		self.ChunkWorld.DropTree(ci*10+ddi, cj*10+ddj, 6+extra)
		// 	}
		// }
		elapsedChunks := time.Since(startChunks)
		log.Printf("Creating %d chunks took %s", len(newChunks), elapsedChunks)

		startRenderData := time.Now()
		renderCount := 0
		for _, chunk := range newChunks {
			if chunk.IsVisible() {
				renderData := self.ChunkWorld.ChunkRenderDataBuilder.ChunkToRenderData(chunk)
				if renderData != nil {
					renderCount += 1
					self.Scene.AddRenderData(renderData)
				}
			}
		}
		elapsed := time.Since(startRenderData)
		log.Printf("Converting %d chunks to render data took %s", renderCount, elapsed)
		self.shouldSpawnMore = false
	}
}
