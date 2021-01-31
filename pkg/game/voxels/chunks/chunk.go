package chunks

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl64"
)

type IsVoxelTransparent func(i, j, k int) bool

type StandardChunk struct {
	ChunkWorld    *ChunkWorld
	Region        *ChunkRegion
	Coordinate    IndexCoordinate
	Voxels        *[]int
	VisibleVoxels *[]int
}

func NewChunk(world *ChunkWorld, region *ChunkRegion, coord IndexCoordinate) *StandardChunk {
	var vox []int
	var vis []int
	chunk := &StandardChunk{
		world,
		region,
		coord,
		&vox,
		&vis,
	}
	var chunkVoxels []int
	chunk.ForAll(func(i, j, k int) {
		chunkVoxels = append(chunkVoxels, AIR)
	})
	chunk.Voxels = &chunkVoxels
	return chunk
}

func (self *StandardChunk) ForAll(f func(i, j, k int)) {
	chunkWidth := self.ChunkWorld.ChunkSettings.GetChunkWidth()
	chunkDepth := self.ChunkWorld.ChunkSettings.GetChunkDepth()
	chunkHeight := self.ChunkWorld.ChunkSettings.GetChunkHeight()
	for vi := 0; vi < chunkWidth; vi++ {
		for vj := 0; vj < chunkDepth; vj++ {
			for vk := 0; vk < chunkHeight; vk++ {
				f(vi, vj, vk)
			}
		}
	}
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

func (self *StandardChunk) AddVisibleVoxel(i, j, k, voxel int) {
	wasTransparent := self.IsTransparent(i, j, k)
	isNowTransparent := voxel == AIR
	voxelIndexCoord := IndexCoordinate{i, j, k}
	index := self.ChunkWorld.ChunkSettings.CoordinateToIndex(voxelIndexCoord)
	(*self.Voxels)[index] = voxel
	if wasTransparent && !isNowTransparent {
		settings := self.ChunkWorld.ChunkSettings
		index := settings.CoordinateToIndexijk(i, j, k)
		cas := append(*self.VisibleVoxels, index)
		self.VisibleVoxels = &cas
	}
}

func (self *StandardChunk) RemoveVoxel(index int) {
	(*self.Voxels)[index] = AIR
	var newVisibleVoxels []int
	// todo: add surrounding voxels to visible ones
	for _, vi := range *self.VisibleVoxels {
		if vi != index {
			newVisibleVoxels = append(newVisibleVoxels, vi)
		}
	}
	self.VisibleVoxels = &newVisibleVoxels
}

func (self *StandardChunk) AddInvisibleVoxel(i, j, k, voxel int) {
	voxelIndexCoord := IndexCoordinate{i, j, k}
	index := self.ChunkWorld.ChunkSettings.CoordinateToIndex(voxelIndexCoord)
	(*self.Voxels)[index] = voxel
}

func (self *StandardChunk) IsVisible() bool {
	return len(*self.VisibleVoxels) > 0
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
