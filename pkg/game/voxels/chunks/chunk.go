package chunks

import (
	"math"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl64"
)

type IsVoxelTransparent func(i, j, k int) bool

type StandardChunk struct {
	ChunkWorld    *ChunkWorld
	Region        *ChunkRegion
	Coordinate    IndexCoordinate
	Voxels        *[]Block
	VisibleVoxels map[int]void
}

func NewChunk(world *ChunkWorld, region *ChunkRegion, coord IndexCoordinate) *StandardChunk {
	vis := make(map[int]void)
	chunkVoxels := make([]Block, world.ChunkSettings.GetChunkVolume())
	for i := range chunkVoxels {
		chunkVoxels[i] = AIR
	}
	chunk := &StandardChunk{
		world,
		region,
		coord,
		&chunkVoxels,
		vis,
	}
	return chunk
}

func (self *StandardChunk) GetVoxel(i, j, k int) int {
	settings := self.ChunkWorld.ChunkSettings
	voxelIndexCoord := IndexCoordinate{i, j, k}
	if settings.CoordinateIsOutOfBounds(voxelIndexCoord) {
		rawRegionCoord := []IndexCoordinate{voxelIndexCoord, self.Coordinate}
		region2Coord, isOrigin := self.Region.GetRegionCoodinate()
		if !isOrigin {
			rawRegionCoord = append(rawRegionCoord, region2Coord...)
		}
		regionCoord := settings.NormalizeCoordinate(rawRegionCoord)
		voxel := self.ChunkWorld.GetVoxelAt(regionCoord)
		return voxel
	}
	index := settings.CoordinateToIndex(voxelIndexCoord)
	return (*self.Voxels)[index]
}

func (self *StandardChunk) IsTransparent(i, j, k int) bool {
	v := self.GetVoxel(i, j, k)
	return v == AIR || v == WATER
}

func (self *StandardChunk) AddVisibleVoxel(i, j, k, voxel Block) {
	wasTransparent := self.IsTransparent(i, j, k)
	isNowTransparent := voxel == AIR
	voxelIndexCoord := IndexCoordinate{i, j, k}
	settings := self.ChunkWorld.ChunkSettings
	index := settings.CoordinateToIndex(voxelIndexCoord)
	(*self.Voxels)[index] = voxel
	if wasTransparent && !isNowTransparent {
		self.SetVoxelVisibility(index, true)
	}
}

func (self *StandardChunk) SetVoxelVisibility(index int, visible bool) {
	if visible {
		self.VisibleVoxels[index] = VOID
	} else {
		delete(self.VisibleVoxels, index)
	}
}

func (self *StandardChunk) SetAllVoxels(voxel int) {
	for k := range *self.Voxels {
		(*self.Voxels)[k] = voxel
	}
}

func (self *StandardChunk) RemoveVoxel(index int) {
	(*self.Voxels)[index] = AIR
	self.SetVoxelVisibility(index, false)
	// todo: add surrounding voxels to visible ones
	chunkWidth := self.ChunkWorld.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkWorld.ChunkSettings.GetChunkDepth()
	chunkHeight := self.ChunkWorld.ChunkSettings.GetChunkHeight()
	coord := self.ChunkWorld.ChunkSettings.IndexToCoordinate(index)
	east := coord.Addijk(1, 0, 0)
	west := coord.Addijk(-1, 0, 0)
	north := coord.Addijk(0, 1, 0)
	south := coord.Addijk(0, -1, 0)
	up := coord.Addijk(0, 0, 1)
	down := coord.Addijk(0, 0, -1)
	if east.i < chunkWidth {
		eastIndex := self.ChunkWorld.ChunkSettings.CoordinateToIndex(east)
		self.SetVoxelVisibility(eastIndex, true)
	}
	if west.i >= 0 {
		westIndex := self.ChunkWorld.ChunkSettings.CoordinateToIndex(west)
		self.SetVoxelVisibility(westIndex, true)
	}
	if north.j < chunkDepth {
		northIndex := self.ChunkWorld.ChunkSettings.CoordinateToIndex(north)
		self.SetVoxelVisibility(northIndex, true)
	}
	if south.j >= 0 {
		southIndex := self.ChunkWorld.ChunkSettings.CoordinateToIndex(south)
		self.SetVoxelVisibility(southIndex, true)
	}
	if up.k < chunkHeight {
		upIndex := self.ChunkWorld.ChunkSettings.CoordinateToIndex(up)
		self.SetVoxelVisibility(upIndex, true)
	}
	if down.k >= 0 {
		downIndex := self.ChunkWorld.ChunkSettings.CoordinateToIndex(down)
		self.SetVoxelVisibility(downIndex, true)
	}
}

func (self *StandardChunk) AddInvisibleVoxel(i, j, k, voxel int) {
	voxelIndexCoord := IndexCoordinate{i, j, k}
	index := self.ChunkWorld.ChunkSettings.CoordinateToIndex(voxelIndexCoord)
	(*self.Voxels)[index] = voxel
}

func (self *StandardChunk) IsVisible() bool {
	return len(self.VisibleVoxels) > 0
}

func (self *StandardChunk) CalculateOrigin() mgl64.Vec3 {
	chunkWidth := self.ChunkWorld.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkWorld.ChunkSettings.GetChunkDepth()
	chunkHeight := self.ChunkWorld.ChunkSettings.GetChunkHeight()
	voxelSize := self.ChunkWorld.ChunkSettings.GetVoxelSize()
	return mgl64.Vec3{
		float64(float32(self.Coordinate.i*chunkWidth) * voxelSize),
		float64(float32(self.Coordinate.j*chunkDepth) * voxelSize),
		float64(float32(self.Coordinate.k*chunkHeight) * voxelSize),
	}
}

func (self *StandardChunk) GetVoxelAABB(index int) ant.AABB64 {
	settings := self.ChunkWorld.ChunkSettings
	coord := settings.IndexToCoordinate(index)
	voxelSize := float64(settings.GetVoxelSize())
	chunkOrigin := self.CalculateOrigin() // todo: put on state to improve performance
	positionInChunk := mgl64.Vec3{
		float64(coord.i) * voxelSize,
		float64(coord.j) * voxelSize,
		float64(coord.k) * voxelSize,
	}
	voxelMin := chunkOrigin.Add(positionInChunk)
	voxelMax := voxelMin.Add(mgl64.Vec3{voxelSize, voxelSize, voxelSize})
	return ant.AABB64{Min: voxelMin, Max: voxelMax}
}

func (self *StandardChunk) Intersect(p1, p2 mgl64.Vec3) (*IntersectionEvent, bool) {
	tmin := math.MaxFloat64
	targetVoxelIndex := -1
	for vIndex := range self.VisibleVoxels {
		voxel := (*self.Voxels)[vIndex]
		if voxel != AIR {
			voxelAABB := self.GetVoxelAABB(vIndex)
			intersects, t := voxelAABB.LineIntersects(p1, p2)
			// todo: get interestion face for adding voxels
			if intersects && t < tmin {
				tmin = t
				targetVoxelIndex = vIndex
			}
		}
	}
	if targetVoxelIndex == -1 {
		return nil, false
	}
	event := IntersectionEvent{
		Chunk:      self,
		VoxelIndex: targetVoxelIndex,
		Face:       -1, // todo; get interesecting face
	}
	return &event, true
}
