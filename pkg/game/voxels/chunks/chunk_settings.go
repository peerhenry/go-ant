package chunks

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type IChunkSettings interface {
	GetChunkWidth() int
	GetChunkDepth() int
	GetChunkHeight() int
	GetChunkVolume() int
	CoordinateToIndex(c IndexCoordinate) int
	CoordinateToIndexijk(i, j, k int) int
	IndexToCoordinate(index int) IndexCoordinate
	CoordinateIsOutOfBounds(c IndexCoordinate) bool
	IndexIsOutOfBounds(index int) bool
	GetVoxelSize() float32
	NormalizeCoordinate(c []IndexCoordinate) []IndexCoordinate
	ToRegionCoord(location mgl64.Vec3) []IndexCoordinate
	GetChunkCoord(location mgl64.Vec3) IndexCoordinate
}

type StandardChunkSettings struct {
	chunkWidth            int
	chunkDepth            int
	chunkHeight           int
	chunkDepthTimesHeight int
	voxelSize             float32
}

var _ IChunkSettings = (*StandardChunkSettings)(nil)

func NewChunkSettings(chunkWidth, chunkDepth, chunkHeight int) *StandardChunkSettings {
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
	i, j, k := self.IndexToCoordinateijk(index)
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

func (self *StandardChunkSettings) GetChunkVolume() int {
	return self.chunkWidth * self.chunkDepth * self.chunkHeight
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

func (self StandardChunkSettings) AddCoordinatei(coord []IndexCoordinate, i int) []IndexCoordinate {
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
	} else if newi < 0 {
		rankupdi := newi/self.chunkWidth - 1
		remainderi := rankupdi - rankupdi*self.chunkWidth
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{rankupdi, 0, 0})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i + rankupdi, rankup.j, rankup.k}
		}
		coord[0] = IndexCoordinate{remainderi, root.j, root.k}
	} else {
		coord[0] = IndexCoordinate{newi, root.j, root.k}
	}
	return coord
}

func (self StandardChunkSettings) AddCoordinateijk(coord []IndexCoordinate, i, j, k int) []IndexCoordinate {
	root := coord[0]
	coord[0] = IndexCoordinate{root.i + i, root.j + j, root.k + k}
	return self.NormalizeCoordinate(coord)
}

func (self StandardChunkSettings) NormalizeCoordinate(coord []IndexCoordinate) []IndexCoordinate {
	newi := coord[0].i
	if newi >= self.chunkWidth {
		remainderi := newi % self.chunkWidth
		rankupdi := newi / self.chunkWidth
		ranks := len(coord)
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{rankupdi, 0, 0})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i + rankupdi, rankup.j, rankup.k}
		}
		coord[0] = IndexCoordinate{remainderi, coord[0].j, coord[0].k}
	} else if newi < 0 {
		rankupdi := (newi+1)/self.chunkWidth - 1
		remainderi := newi - rankupdi*self.chunkWidth
		ranks := len(coord)
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{rankupdi, 0, 0})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i + rankupdi, rankup.j, rankup.k}
		}
		coord[0] = IndexCoordinate{remainderi, coord[0].j, coord[0].k}
	} else {
		coord[0] = IndexCoordinate{newi, coord[0].j, coord[0].k}
	}

	newj := coord[0].j
	if newj >= self.chunkDepth {
		remainderj := newj % self.chunkDepth
		rankupdj := newj / self.chunkDepth
		ranks := len(coord)
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{0, rankupdj, 0})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i, rankup.j + rankupdj, rankup.k}
		}
		coord[0] = IndexCoordinate{coord[0].i, remainderj, coord[0].k}
	} else if newj < 0 {
		rankupdj := (newj+1)/self.chunkDepth - 1
		remainderj := newj - rankupdj*self.chunkDepth
		ranks := len(coord)
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{0, rankupdj, 0})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i, rankup.j + rankupdj, rankup.k}
		}
		coord[0] = IndexCoordinate{coord[0].i, remainderj, coord[0].k}
	} else {
		coord[0] = IndexCoordinate{coord[0].i, newj, coord[0].k}
	}

	newk := coord[0].k
	if newk >= self.chunkHeight {
		remainderk := newk % self.chunkHeight
		rankupdk := newk / self.chunkHeight
		ranks := len(coord)
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{0, 0, rankupdk})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i, rankup.j, rankup.k + rankupdk}
		}
		coord[0] = IndexCoordinate{coord[0].i, coord[0].j, remainderk}
	} else if newk < 0 {
		rankupdk := (newk+1)/self.chunkHeight - 1
		remainderk := newk - rankupdk*self.chunkHeight
		ranks := len(coord)
		if ranks == 1 {
			coord = append(coord, IndexCoordinate{0, 0, rankupdk})
		} else {
			rankup := coord[1]
			coord[1] = IndexCoordinate{rankup.i, rankup.j, rankup.k + rankupdk}
		}
		coord[0] = IndexCoordinate{coord[0].i, coord[0].j, remainderk}
	} else {
		coord[0] = IndexCoordinate{coord[0].i, coord[0].j, newk}
	}

	return coord
}

func (self *StandardChunkSettings) ToRegionCoord(location mgl64.Vec3) []IndexCoordinate {
	i := int(math.Floor(location[0]))
	j := int(math.Floor(location[1]))
	k := int(math.Floor(location[2]))
	return self.AddCoordinateijk([]IndexCoordinate{IndexCoordinate{0, 0, 0}}, i, j, k)
}

func (self *StandardChunkSettings) GetChunkCoord(location mgl64.Vec3) IndexCoordinate {
	ci := int(math.Floor(location[0] / float64(self.chunkWidth)))
	cj := int(math.Floor(location[1] / float64(self.chunkDepth)))
	ck := int(math.Floor(location[2] / float64(self.chunkHeight)))
	return IndexCoordinate{ci, cj, ck}
}
