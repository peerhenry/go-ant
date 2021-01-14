package chunks

type ChunkRegion struct {
	Chunks map[ChunkIndex]*StandardChunk
}

func (self *ChunkRegion) SetChunkRegion(chunk *StandardChunk) {
	self.Chunks[chunk.Index] = chunk
}

func (self *ChunkRegion) GetChunkRegion(index ChunkIndex) (*StandardChunk, bool) {
	thing, ok := self.Chunks[index]
	return thing, ok
}
