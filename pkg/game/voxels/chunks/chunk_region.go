package chunks

type ChunkRegion struct {
	Chunks map[IndexCoordinate]*StandardChunk
}

func (self *ChunkRegion) SetChunkRegion(chunk *StandardChunk) {
	self.Chunks[chunk.Coordinate] = chunk
}

func (self *ChunkRegion) GetChunkRegion(index IndexCoordinate) (*StandardChunk, bool) {
	thing, ok := self.Chunks[index]
	return thing, ok
}
