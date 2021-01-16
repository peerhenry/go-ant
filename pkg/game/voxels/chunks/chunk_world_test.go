package chunks

import "testing"

func TestHeightToCoordinates(t *testing.T) {
	ExpectHeightToGet(t, 15, 2, 3)
	ExpectHeightToGet(t, -15, 2, -3) // test negative
	ExpectHeightToGet(t, -2, 0, 0)   // test boundaries
	ExpectHeightToGet(t, 3, 0, 1)
	ExpectHeightToGet(t, -7, 0, -1)
}

func ExpectHeightToGet(t *testing.T, h, evk, eck int) {
	settings := NewChunkSettings(5, 5, 5)
	world := NewChunkWorld(settings)
	// Act
	vk, ck := world.HeightToCoordinates(h)
	// Assert
	if vk != evk {
		t.Errorf("Expected vk: %d, got: %d", evk, vk)
	}
	if ck != eck {
		t.Errorf("Expected ck: %d, got: %d", eck, ck)
	}
}

func TestGetOrCreateChunkAt(t *testing.T) {
	settings := NewChunkSettings(5, 5, 5)
	world := NewChunkWorld(settings)
	// Act
	chunk := world.GetOrCreateChunkAt(IndexCoordinate{1, 2, 3})
	// Assert
	if chunk == nil {
		t.Errorf("chunk is nil!")
	}
	if len(*chunk.Voxels) != 125 {
		t.Errorf("length of voxels is not 125, but %d", len(*chunk.Voxels))
	}
}

func TestCreateChunksInColumn(t *testing.T) {
	settings := NewChunkSettings(5, 5, 5)
	world := NewChunkWorld(settings)
	// Act
	world.CreateChunksInColumn(1, 2)
	// Assert no crash
}
