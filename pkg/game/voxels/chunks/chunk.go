package chunks

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

func (self *StandardChunk) AddInvisibleVoxel(i, j, k, voxel int) {
	voxelIndexCoord := IndexCoordinate{i, j, k}
	index := self.ChunkWorld.ChunkSettings.CoordinateToIndex(voxelIndexCoord)
	(*self.Voxels)[index] = voxel
}

func (self *StandardChunk) IsVisible() bool {
	return len(*self.VisibleVoxels) > 0
}
