package chunks

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func arrayToString(array interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array)), ","), "[]")
}

func mockChunkWorld(chunkSize int) *ChunkWorld {
	chunkSettings := NewChunkSettings(chunkSize, chunkSize, chunkSize)
	return NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
}

func ExpectVec3Equals(t *testing.T, expect, actual mgl64.Vec3) {
	if actual[0] != expect[0] {
		t.Errorf("expected x to be %f, got %f", expect[0], actual[0])
	}
	if actual[1] != expect[1] {
		t.Errorf("expected y yo be %f, got %f", expect[1], actual[1])
	}
	if actual[2] != expect[2] {
		t.Errorf("expected z to be %f, got %f", expect[2], actual[2])
	}
}
