package chunks

type IChunkSettings interface {
	GetChunkWidth() int
	GetChunkDepth() int
	GetChunkHeight() int
	CoordinateToIndex(c IndexCoordinate) int
	CoordinateToIndexijk(i, j, k int) int
	IndexToCoordinate(index int) IndexCoordinate
	CoordinateIsOutOfBounds(c IndexCoordinate) bool
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

func (self *StandardChunkSettings) CoordinateToIndex(c IndexCoordinate) int {
	return self.chunkDepthTimesHeight*c.i + self.chunkHeight*c.j + c.k
}

func (self *StandardChunkSettings) CoordinateToIndexijk(i, j, k int) int {
	return self.chunkDepthTimesHeight*i + self.chunkHeight*j + k
}

func (self *StandardChunkSettings) IndexToCoordinate(index int) IndexCoordinate {
	k := index % self.chunkHeight
	j := (index % self.chunkDepthTimesHeight) / self.chunkHeight
	i := index / self.chunkDepthTimesHeight
	return IndexCoordinate{i, j, k}
}

func (self *StandardChunkSettings) IndexToCoordinateijk(index int) (int, int, int) {
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

func (self *StandardChunkSettings) CoordinateIsOutOfBounds(c IndexCoordinate) bool {
	return c.i < 0 || c.i >= self.chunkWidth || c.j < 0 || c.j >= self.chunkDepth || c.k < 0 || c.k >= self.chunkHeight
}

func (self *StandardChunkSettings) IndexIsOutOfBounds(index int) bool {
	c := self.IndexToCoordinate(index)
	return self.CoordinateIsOutOfBounds(c)
}

func (self *StandardChunkSettings) GetVoxelSize() float32 {
	return self.voxelSize
}

func (self StandardChunkSettings) Addi(coord []IndexCoordinate, i int) []IndexCoordinate {
	ranks := len(coord)
	root := coord[0]
	newi := root.i + i
	if newi >= self.chunkWidth {
		remainderi := newi % self.chunkWidth
		rankupdi := newi / self.chunkWidth
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{rankupdi, 0, 0})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i + rankupdi, rankup.j, rankup.k}
		}
		coord[0] = IndexCoordinate{remainderi, root.j, root.k}
		return coord
	}
	if newi < 0 {
		rankupdi := newi/self.chunkWidth - 1
		remainderi := rankupdi - rankupdi*self.chunkWidth
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{rankupdi, 0, 0})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i + rankupdi, rankup.j, rankup.k}
		}
		coord[0] = IndexCoordinate{remainderi, root.j, root.k}
		return coord
	}
	return []IndexCoordinate{IndexCoordinate{newi, root.j, root.k}}
}
