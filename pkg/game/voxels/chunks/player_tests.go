package chunks

import (
	"testing"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl64"
)

func TestToRegionCoord(t *testing.T) {
	cam := ant.NewCamera()
	chunkSettings := NewChunkSettings(32, 32, 8)
	world := NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
	player := NewPlayer(cam, world)
	// Act
	coords := player.ToRegionCoord(mgl64.Vec3{100.0, 100.0, 100.0})
	vCoord := coords[0]
	chunkCoord := coords[1]
	// Assert
	vCoordExpect := IndexCoordinate{0, 0, 0}
	if !vCoord.Equals(vCoordExpect) {
		t.Errorf("Expected vCoord to be %s but got %s", vCoordExpect.ToString(), vCoord.ToString())
	}
	chunkCoordExpected := IndexCoordinate{0, 0, 0}
	if !chunkCoord.Equals(chunkCoordExpected) {
		t.Errorf("Expected chunkCoord to be %s but got %s", chunkCoordExpected.ToString(), chunkCoord.ToString())
	}
}
