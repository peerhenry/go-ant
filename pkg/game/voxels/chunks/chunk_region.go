package chunks

type ChunkRegion struct {
	Chunks map[IndexCoordinate]*StandardChunk
}

func NewChunkRegion() *ChunkRegion {
	return &ChunkRegion{
		Chunks: make(map[IndexCoordinate]*StandardChunk),
	}
}

func (self *ChunkRegion) SetChunkRegion(chunk *StandardChunk) {
	self.Chunks[chunk.Coordinate] = chunk
}

func (self *ChunkRegion) GetChunkRegion(index IndexCoordinate) (*StandardChunk, bool) {
	thing, ok := self.Chunks[index]
	return thing, ok
}

func (self *ChunkRegion) GetRegionCoodinate() ([]IndexCoordinate, bool) {
	return []IndexCoordinate{IndexCoordinate{0, 0, 0}}, true // boolean returns if coordinate is origin
}
