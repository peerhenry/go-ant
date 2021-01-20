package chunks

import "testing"

func TestIsTransparent(t *testing.T) {
	// create chunk
	world := mockChunkWorld(2)
	mychunk := &StandardChunk{
		ChunkWorld: world,
		Region:     world.Region,
		Coordinate: IndexCoordinate{0, 0, 0},
		Voxels:     &[]int{STONE, STONE, STONE, STONE, STONE, STONE, STONE, STONE},
	}
	// create six chunks around it
	world.Region.SetChunkRegion(CreateFullChunk(world, -1, 0, 0))
	world.Region.SetChunkRegion(CreateFullChunk(world, 1, 0, 0))
	world.Region.SetChunkRegion(CreateFullChunk(world, 0, -1, 0))
	world.Region.SetChunkRegion(CreateFullChunk(world, 0, 1, 0))
	world.Region.SetChunkRegion(CreateFullChunk(world, 0, 0, -1))
	world.Region.SetChunkRegion(CreateFullChunk(world, 0, 0, 1))
	// test TestIsTransparent for chunk coordinates that normalize to adjacent chunks
	expectOpaqueAt(t, mychunk, -1, 0, 0)
	expectOpaqueAt(t, mychunk, 2, 0, 0)
	expectOpaqueAt(t, mychunk, 0, -1, 0)
	expectOpaqueAt(t, mychunk, 0, 2, 0)
	expectOpaqueAt(t, mychunk, 0, 0, -1)
	expectOpaqueAt(t, mychunk, 0, 0, 2)
}

func expectOpaqueAt(t *testing.T, chunk *StandardChunk, i, j, k int) {
	if chunk.IsTransparent(i, j, k) != false {
		t.Errorf("Expected chunk NOT to be transparent at %d, %d, %d", i, j, k)
	}
}

func CreateFullChunk(world *ChunkWorld, ci, cj, ck int) *StandardChunk {
	var chunkVoxels []int
	var visibleVoxels []int
	chunkSettings := world.ChunkSettings
	chunkWidth := chunkSettings.GetChunkWidth()
	chunkDepth := chunkSettings.GetChunkDepth()
	chunkHeight := chunkSettings.GetChunkHeight()
	// set data of every voxel in the chunk
	for vi := 0; vi < chunkWidth; vi++ {
		for vj := 0; vj < chunkSettings.GetChunkDepth(); vj++ {
			for vk := 0; vk < chunkSettings.GetChunkHeight(); vk++ {
				voxel := STONE
				chunkVoxels = append(chunkVoxels, voxel)
				if vi == 0 || vi == chunkWidth-1 || vj == 0 || vj == chunkDepth-1 || vk == 0 || vk == chunkHeight-1 {
					if voxel != AIR {
						index := chunkSettings.CoordinateToIndexijk(vi, vj, vk)
						visibleVoxels = append(visibleVoxels, index)
					}
				}
			}
		}
	}
	return &StandardChunk{
		Coordinate:    IndexCoordinate{ci, cj, ck},
		Voxels:        &chunkVoxels,
		VisibleVoxels: &visibleVoxels,
		ChunkWorld:    world,
		Region:        world.Region,
	}
}
