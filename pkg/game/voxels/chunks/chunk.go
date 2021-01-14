package chunks

type IsVoxelTransparent func(i, j, k int) bool

type StandardChunk struct {
	Index         ChunkIndex
	voxels        *[]int
	visibleVoxels *[]int
	chunkSettings IChunkSettings
}

func (self *StandardChunk) IsTransparent(i, j, k int) bool {
	// todo: check adjacent chunk
	if self.chunkSettings.CoordinateIsOutOfBounds(i, j, k) {
		return true
	}
	index := self.chunkSettings.CoordinateToIndex(i, j, k)
	return (*self.voxels)[index] == 0
}
