package chunks

import (
	"fmt"
)

type IndexCoordinate struct {
	i int
	j int
	k int
}

type RegionCoordinate = []IndexCoordinate

func (self IndexCoordinate) ToString() string {
	return "IndexCoordinate {" + fmt.Sprintf("%v, %v, %v", self.i, self.j, self.k) + "}"
}

func (self IndexCoordinate) Equals(other IndexCoordinate) bool {
	return self.i == other.i && self.j == other.j && self.k == other.k
}

func (self IndexCoordinate) Addijk(di, dj, dk int) IndexCoordinate {
	return IndexCoordinate{self.i + di, self.j + dj, self.k + dk}
}
