package chunks

type IChunkSettings interface {
	GetChunkWidth() int
	GetChunkDepth() int
	GetChunkHeight() int
	CoordinateToIndex(i, j, k int) int
	IndexToCoordinate(index int) (int, int, int)
	CoordinateIsOutOfBounds(i, j, k int) bool
	IndexIsOutOfBounds(index int) bool
	GetVoxelSize() float32
}

type StandardChunkSettings struct {
	chunkWidth            int
	chunkDepth            int
	chunkHeight           int
	chunkDepthTimesHeight int
	voxelSize             float32
}

func CreateStandardChunkSettings(chunkWidth, chunkDepth, chunkHeight int) *StandardChunkSettings {
	return &StandardChunkSettings{
		chunkWidth:            chunkWidth,
		chunkDepth:            chunkDepth,
		chunkHeight:           chunkHeight,
		chunkDepthTimesHeight: chunkDepth * chunkHeight,
		voxelSize:             1.0,
	}
}

func (self *StandardChunkSettings) CoordinateToIndex(i, j, k int) int {
	return self.chunkDepthTimesHeight*i + self.chunkHeight*j + k
}

func (self *StandardChunkSettings) IndexToCoordinate(index int) (int, int, int) {
	k := index % self.chunkHeight
	j := (index % self.chunkDepthTimesHeight) / self.chunkHeight
	i := index / self.chunkDepthTimesHeight
	return i, j, k
}

func (self *StandardChunkSettings) GetChunkWidth() int {
	return self.chunkWidth
}

func (self *StandardChunkSettings) GetChunkDepth() int {
	return self.chunkDepth
}

func (self *StandardChunkSettings) GetChunkHeight() int {
	return self.chunkHeight
}

func (self *StandardChunkSettings) CoordinateIsOutOfBounds(i, j, k int) bool {
	return i < 0 || i >= self.chunkWidth || j < 0 || j >= self.chunkDepth || k < 0 || k >= self.chunkHeight
}

func (self *StandardChunkSettings) IndexIsOutOfBounds(index int) bool {
	i, j, k := self.IndexToCoordinate(index)
	return self.CoordinateIsOutOfBounds(i, j, k)
}

func (self *StandardChunkSettings) GetVoxelSize() float32 {
	return self.voxelSize
}
