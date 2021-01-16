package chunks

type IsVoxelTransparent func(i, j, k int) bool

type StandardChunk struct {
	ChunkWorld    *ChunkWorld
	Region        *ChunkRegion
	Coordinate    IndexCoordinate
	Voxels        *[]int
	VisibleVoxels *[]int
}

func (self *StandardChunk) IsTransparent(i, j, k int) bool {
	settings := self.ChunkWorld.ChunkSettings
	voxelIndexCoord := IndexCoordinate{i, j, k}
	if settings.CoordinateIsOutOfBounds(voxelIndexCoord) {
		// region2Coord := self.Region.GetRegionCoodinate()
		// rawRegionCoord := append([]IndexCoordinate{voxelIndexCoord, self.Coordinate}, region2Coord...)
		// regionCoord := settings.NormalizeCoordinate(rawRegionCoord)
		// voxel := self.ChunkWorld.GetVoxelAt(regionCoord)
		// return voxel == AIR
		return true
	}
	index := settings.CoordinateToIndex(voxelIndexCoord)
	return (*self.Voxels)[index] == AIR
}
