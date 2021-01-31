package chunks

type ChunkRegion struct {
	Chunks map[IndexCoordinate]*StandardChunk
}

func NewChunkRegion() *ChunkRegion {
	return &ChunkRegion{
		Chunks: make(map[IndexCoordinate]*StandardChunk),
	}
}

func (self *ChunkRegion) SetChunk(chunk *StandardChunk) {
	self.Chunks[chunk.Coordinate] = chunk
}

func (self *ChunkRegion) GetChunk(index IndexCoordinate) (*StandardChunk, bool) {
	thing, ok := self.Chunks[index]
	return thing, ok
}

func (self *ChunkRegion) GetChunks(coords []IndexCoordinate) []*StandardChunk {
	var chunks []*StandardChunk
	for _, ix := range coords {
		chunk, ok := self.Chunks[ix]
		if ok {
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

func (self *ChunkRegion) GetRegionCoodinate() ([]IndexCoordinate, bool) {
	return []IndexCoordinate{IndexCoordinate{0, 0, 0}}, true // boolean returns if coordinate is origin
}
