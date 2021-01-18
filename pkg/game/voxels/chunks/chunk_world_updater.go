package chunks

import (
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
	columnsQueue    map[ChunkColumnCoordinate]bool
	spawnedColumns  map[ChunkColumnCoordinate]bool
	renderChunks    map[IndexCoordinate]int
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
		columnsQueue:    make(map[ChunkColumnCoordinate]bool),
		spawnedColumns:  make(map[ChunkColumnCoordinate]bool),
		renderChunks:    make(map[IndexCoordinate]int),
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
		cam_ci := int(math.Floor(self.Camera.Position[0] / float64(chunkWidth)))
		cam_cj := int(math.Floor(self.Camera.Position[1] / float64(chunkDepth)))
		ci_min := cam_ci - 12
		ci_max := cam_ci + 12
		cj_min := cam_cj - 12
		cj_max := cam_cj + 12

		// todo: queue columns to spawn in circle rather than square
		for ci := ci_min; ci <= ci_max; ci++ {
			for cj := cj_min; cj <= cj_max; cj++ {
				_, alreadySpawned := self.spawnedColumns[[2]int{ci, cj}]
				_, alreadyQueued := self.columnsQueue[[2]int{ci, cj}]
				if !alreadySpawned && !alreadyQueued {
					self.columnsQueue[[2]int{ci, cj}] = true
				}
			}
		}

		// remove distant chunks
		for coord, renderIndex := range self.renderChunks {
			if coord.i < ci_min || coord.i > ci_max || coord.j < cj_min || coord.j > cj_max {
				self.Scene.RemoveRenderData(renderIndex)
				delete(self.renderChunks, coord)
				delete(self.spawnedColumns, [2]int{coord.i, coord.j})
				delete(self.columnsQueue, [2]int{coord.i, coord.j}) // in case they are queued
				self.ChunkWorld.DeleteChunk(coord)
			}
		}

		self.shouldSpawnMore = false
	}

	self.HandleChunkColumnQueue()
}

func (self *ChunkWorldUpdater) HandleChunkColumnQueue() {
	if len(self.columnsQueue) == 0 {
		return
	}
	newChunks := make(map[IndexCoordinate]*StandardChunk)

	// pop coordinate from queue
	columnCoord := self.PopColumsQueue()
	ci := columnCoord[0]
	cj := columnCoord[1]

	// create chunks in column
	newChunksInColumns := self.ChunkWorld.CreateChunksInColumn(ci, cj)
	for coord, chunk := range newChunksInColumns {
		newChunks[coord] = chunk
	}
	self.spawnedColumns[[2]int{ci, cj}] = true

	// create render data for chunks
	renderCount := 0
	for coord, chunk := range newChunks {
		if chunk.IsVisible() {
			existingIndex, alreadyExists := self.renderChunks[coord]
			renderData := self.ChunkWorld.ChunkRenderDataBuilder.ChunkToRenderData(chunk)
			if renderData != nil {
				renderCount += 1
				if !alreadyExists {
					index := self.Scene.AddRenderData(renderData)
					self.renderChunks[coord] = index
				} else {
					self.Scene.ReplaceRenderData(existingIndex, renderData)
				}
			}
		}
	}
}

func (self *ChunkWorldUpdater) PopColumsQueue() ChunkColumnCoordinate {
	for key := range self.columnsQueue {
		delete(self.columnsQueue, key)
		return key
	}
	return ChunkColumnCoordinate{} // dummy return to prevent compiler error
}
