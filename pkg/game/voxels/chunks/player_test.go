package chunks

import (
	"testing"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl64"
)

func TestClipFromVoxelCollisions_ShouldClipY(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	cam.Position = mgl64.Vec3{0, 0, 5 + playerCamHeight + 0.01}
	world := NewChunkWorldBuilder().Build()
	player := NewPlayer(cam, world)
	player.isFalling = false
	world.CreateChunksInColumn(0, 0)
	world.CreateChunksInColumn(-1, 0)
	world.CreateChunksInColumn(0, -1)
	world.CreateChunksInColumn(-1, -1)
	// Act
	dx := 0.05
	dy := 0.05
	result := player.clipFromVoxelCollisions(mgl64.Vec3{dx, dy, 0.0})
	// Assert
	if result[2] != 0.0 {
		t.Errorf("expected clipped dz to be %f, got %f", 0.0, result[2])
		return
	}
}

func TestClipFromVoxelCollisions_ShouldStopFalling(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	cam.Position = mgl64.Vec3{0, 0, 5 + playerCamHeight + 0.05}
	world := NewChunkWorldBuilder().Build()
	world.CreateChunksInColumn(0, 0)
	world.CreateChunksInColumn(-1, 0)
	world.CreateChunksInColumn(0, -1)
	world.CreateChunksInColumn(-1, -1)
	player := NewPlayer(cam, world)
	player.isFalling = true
	// Act
	dx := 0.05
	dy := 0.05
	result := player.clipFromVoxelCollisions(mgl64.Vec3{dx, dy, -0.1})
	// Assert
	if result[0] != dx {
		t.Errorf("expected dx to be %f, got %f", dx, result[0])
		return
	}
	if result[1] != dy {
		t.Errorf("expected dy yo be %f, got %f", dy, result[1])
		return
	}
	if result[2] != 0.0 {
		t.Errorf("expected clipped dz to be %f, got %f", 0.0, result[2])
		return
	}
	if player.isFalling {
		t.Errorf("expected player.isFalling to be false, but was true")
		return
	}
}

func TestCancelComponent(t *testing.T) {
	// arrange
	cam := ant.NewCamera()
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{0})
	player := NewPlayer(cam, world)
	thing := mgl64.Vec3{1, 1, 1}
	expectX := mgl64.Vec3{0, 1, 1}
	// Act
	west := player.cancelComponent(thing, WEST)
	east := player.cancelComponent(thing, EAST)
	// Assert
	ExpectVec3Equals(t, expectX, west)
	ExpectVec3Equals(t, expectX, east)
}

func TestGetIntersectingVoxelAABBs_ExpectFour(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	cam.Position = mgl64.Vec3{0, 0, 5 + playerCamHeight + 0.05}
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{0})
	world.CreateChunksInColumn(0, 0)
	world.CreateChunksInColumn(-1, 0)
	world.CreateChunksInColumn(0, -1)
	world.CreateChunksInColumn(-1, -1)
	player := NewPlayer(cam, world)
	// Act
	playerBox := player.getFutureAABB(mgl64.Vec3{0, 0, -0.1})
	result := player.getIntersectingVoxelAABBs(playerBox)
	// Assert
	if len(result) != 4 {
		t.Errorf("expected %d intersecting voxels, got %d", 4, len(result))
		return
	}
}

func TestGetIntersectingVoxelAABBs_ExpectZero(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	cam.Position = mgl64.Vec3{0, 0, 1 + playerCamHeight + 0.05}
	chunkSettings := NewChunkSettings(2, 2, 2)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-1})
	world.CreateChunksInColumn(0, 0)
	world.CreateChunksInColumn(-1, 0)
	world.CreateChunksInColumn(0, -1)
	world.CreateChunksInColumn(-1, -1)
	if len(world.Region.Chunks) != 4 {
		t.Errorf("expected %d chunks, got %d", 4, len(world.Region.Chunks))
	}
	player := NewPlayer(cam, world)
	// Act
	playerBox := player.getFutureAABB(mgl64.Vec3{0, 0, 0})
	result := player.getIntersectingVoxelAABBs(playerBox)
	// Assert
	if len(result) != 0 {
		t.Errorf("expected %d intersecting voxels, got %d", 0, len(result))
		return
	}
}

func TestGetIntersectingChunks_ExpectFour(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	cam.Position = mgl64.Vec3{0, 0, 1 + playerCamHeight + 0.05}
	chunkSettings := NewChunkSettings(2, 2, 2)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{0})
	world.CreateChunksInColumn(0, 0)
	world.CreateChunksInColumn(-1, 0)
	world.CreateChunksInColumn(0, -1)
	world.CreateChunksInColumn(-1, -1)
	if len(world.Region.Chunks) != 4 {
		t.Errorf("expected %d chunks, got %d", 4, len(world.Region.Chunks))
	}
	player := NewPlayer(cam, world)
	// Act
	playerBox := player.getFutureAABB(mgl64.Vec3{0, 0, 0})
	result := player.getIntersectingChunks(playerBox)
	// Assert
	if len(result) != 4 {
		t.Errorf("expected %d intersecting chunks, got %d", 4, len(result))
		return
	}
}

func TestGetIntersectingChunksExpectZero(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
	player := NewPlayer(cam, world)
	playerMin := player.Camera.Position.Sub(mgl64.Vec3{0.2, 0.2, 1.8})
	playerMax := player.Camera.Position.Add(mgl64.Vec3{0.2, 0.2, 0.2})
	aabb := ant.AABB64{Min: playerMin, Max: playerMax}
	// Act
	chunks := player.getIntersectingChunks(aabb)
	// Assert
	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks, got, %d", len(chunks))
		return
	}
}

func TestGetIntersectingChunksExpectEight(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	cam.Position = mgl64.Vec3{32.1, 32.1, 8.5}
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
	world.GetOrCreateChunkAt(IndexCoordinate{0, 0, 0})
	world.GetOrCreateChunkAt(IndexCoordinate{1, 0, 0})
	world.GetOrCreateChunkAt(IndexCoordinate{0, 1, 0})
	world.GetOrCreateChunkAt(IndexCoordinate{1, 1, 0})
	world.GetOrCreateChunkAt(IndexCoordinate{0, 0, 1})
	world.GetOrCreateChunkAt(IndexCoordinate{1, 0, 1})
	world.GetOrCreateChunkAt(IndexCoordinate{0, 1, 1})
	world.GetOrCreateChunkAt(IndexCoordinate{1, 1, 1})
	player := NewPlayer(cam, world)
	playerMin := player.Camera.Position.Sub(mgl64.Vec3{0.2, 0.2, 1.8})
	playerMax := player.Camera.Position.Add(mgl64.Vec3{0.2, 0.2, 0.2})
	aabb := ant.AABB64{Min: playerMin, Max: playerMax}
	// Act
	chunks := player.getIntersectingChunks(aabb)
	// Assert
	if len(chunks) != 8 {
		t.Errorf("expected 8 chunks, got, %d", len(chunks))
		return
	}
}

func TestToRegionCoord(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
	player := NewPlayer(cam, world)
	// Act
	coords := player.ToRegionCoord(mgl64.Vec3{102.0, 101.0, 100.0})
	// Assert
	if len(coords) != 2 {
		t.Errorf("Expected region coordinate to have length 2, but was %d", len(coords))
		return
	}
	vCoord := coords[0]
	chunkCoord := coords[1]
	vCoordExpect := IndexCoordinate{6, 5, 4}
	if !vCoord.Equals(vCoordExpect) {
		t.Errorf("Expected vCoord to be %s but got %s", vCoordExpect.ToString(), vCoord.ToString())
	}
	chunkCoordExpected := IndexCoordinate{3, 3, 12}
	if !chunkCoord.Equals(chunkCoordExpected) {
		t.Errorf("Expected chunkCoord to be %s but got %s", chunkCoordExpected.ToString(), chunkCoord.ToString())
	}
}

// sanity
func TestToRegionCoordBigDistance(t *testing.T) {
	// Arrange
	cam := ant.NewCamera()
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
	player := NewPlayer(cam, world)
	// Act
	coords := player.ToRegionCoord(mgl64.Vec3{3200000.0, 3200000.0, 3200000.0}) // 3200 km away in each direction
	// Assert
	if len(coords) != 2 {
		t.Errorf("Expected region coordinate to have length 2, but was %d", len(coords))
		return
	}
	vCoord := coords[0]
	chunkCoord := coords[1]
	vCoordExpect := IndexCoordinate{0, 0, 0}
	if !vCoord.Equals(vCoordExpect) {
		t.Errorf("Expected vCoord to be %s but got %s", vCoordExpect.ToString(), vCoord.ToString())
	}
	chunkCoordExpected := IndexCoordinate{100000, 100000, 400000}
	if !chunkCoord.Equals(chunkCoordExpected) {
		t.Errorf("Expected chunkCoord to be %s but got %s", chunkCoordExpected.ToString(), chunkCoord.ToString())
	}
}

// func TestGetSextant_West(t *testing.T) {
// 	// Arrange
// 	cam := ant.NewCamera()
// 	cam.Position = mgl64.Vec3{playerBoxSize, playerBoxSize, playerCamHeight} // this puts player box min at origin 0,0,0
// 	chunkSettings := NewChunkSettings(32, 32, 8)
// 	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
// 	player := NewPlayer(cam, world)
// 	playerBoxRatio := (2 * playerBoxSize) / playerBoxHeight
// 	z_line := -1.0 + 0.05                    // z coordinate of aabb; raised slightly from directly below player box
// 	aabb_center_x := playerBoxRatio * z_line // given z, put x on the line from player box center through origin (in x-z plane)
// 	aabb_center_z := z_line + 0.05           // lift aabb z slightly more to put it in west sextant
// 	aabb_xmin := aabb_center_x - 0.5
// 	aabb_ymin := -0.5
// 	aabb_zmin := aabb_center_z - 0.5
// 	aabb := ant.AABB64{
// 		Min: mgl64.Vec3{aabb_xmin, aabb_ymin, aabb_zmin},
// 		Max: mgl64.Vec3{aabb_xmin + 1, aabb_ymin + 1, aabb_zmin + 1}}
// 	// Act
// 	result := player.GetSextant(mgl64.Vec3{0, 0, 0}, aabb)
// 	// Assert
// 	if result != WEST {
// 		t.Errorf("GetSextant expected %s, got %s", WEST.ToString(), result.ToString())
// 		return
// 	}
// }

// func TestGetSextant_Bottom(t *testing.T) {
// 	// Arrange
// 	cam := ant.NewCamera()
// 	chunkSettings := NewChunkSettings(32, 32, 8)
// 	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
// 	player := NewPlayer(cam, world)
// 	aabb := ant.AABB64{Min: mgl64.Vec3{-0.5, -0.5, -1.9}, Max: mgl64.Vec3{0.5, 0.5, -0.9}}
// 	// Act
// 	result := player.GetSextant(mgl64.Vec3{0, 0, 0}, aabb)
// 	// Assert
// 	if result != BOTTOM {
// 		t.Errorf("GetSextant expected %d, got, %d", BOTTOM, result)
// 		return
// 	}
// }

// func TestGetSextant(t *testing.T) {
// 	// Arrange
// 	cam := ant.NewCamera()
// 	chunkSettings := NewChunkSettings(32, 32, 8)
// 	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
// 	player := NewPlayer(cam, world)
// 	aabb := ant.AABB64{Min: mgl64.Vec3{-0.5, -0.5, -1.9}, Max: mgl64.Vec3{0.5, 0.5, -0.9}}
// 	// Act
// 	result := player.GetSextant(mgl64.Vec3{0, 0, 0}, aabb)
// 	// Assert
// 	if result != BOTTOM {
// 		t.Errorf("GetSextant expected %d, got, %d", BOTTOM, result)
// 		return
// 	}
// }
