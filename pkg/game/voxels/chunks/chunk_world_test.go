package chunks

import "testing"

func TestHeightToCoordinates(t *testing.T) {
	ExpectHeightToGet(t, 2, 2, 0)
	ExpectHeightToGet(t, -2, 3, -1)
	ExpectHeightToGet(t, 15, 0, 3)
	ExpectHeightToGet(t, -15, 0, -3)
}

func ExpectHeightToGet(t *testing.T, height, expect_vk, expect_ck int) {
	world := mockChunkWorld(5)
	// Act
	vk, ck := world.HeightToCoordinates(height)
	// Assert
	if vk != expect_vk {
		t.Errorf("Expected vk: %d, got: %d", expect_vk, vk)
	}
	if ck != expect_ck {
		t.Errorf("Expected ck: %d, got: %d", expect_ck, ck)
	}
}

func TestGetOrCreateChunkAt(t *testing.T) {
	world := mockChunkWorld(5)
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
	world := mockChunkWorld(5)
	// Act
	world.CreateChunksInColumn(1, 2)
	// Assert no crash
}
