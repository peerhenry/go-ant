package chunks

type IsVoxelTransparent func(i, j, k int) bool

type StandardChunk struct {
	Coordinate    IndexCoordinate
	Voxels        *[]int
	VisibleVoxels *[]int
	ChunkSettings IChunkSettings
}

func (self *StandardChunk) IsTransparent(i, j, k int) bool {
	// todo: check adjacent chunk
	if self.ChunkSettings.CoordinateIsOutOfBounds(IndexCoordinate{i, j, k}) {
		return true
	}
	index := self.ChunkSettings.CoordinateToIndex(IndexCoordinate{i, j, k})
	return (*self.Voxels)[index] == 0
}
