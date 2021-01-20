package chunks

import (
	"testing"
)

func TestCreateChunkData(t *testing.T) {
	world := mockChunkWorld(7)
	chunkBuilder := world.ChunkBuilder
	chunk := chunkBuilder.CreateChunk(world, 0, 0, 0)
	result1 := chunk.IsTransparent(0, 0, 0)
	result2 := chunk.IsTransparent(4, 4, 4)
	if result1 {
		t.Errorf("Expected chunk.IsTransparent(0, 0, 0) to be false but was true")
	}
	if result2 {
		t.Errorf("Expected chunk.IsTransparent(6, 6, 6) to be false but was true")
	}
}
