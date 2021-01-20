package chunks

import (
	"fmt"
	"strings"
)

func arrayToString(array interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array)), ","), "[]")
}

func mockChunkWorld(chunkSize int) *ChunkWorld {
	chunkSettings := NewChunkSettings(chunkSize, chunkSize, chunkSize)
	return NewChunkWorld(chunkSettings, HeightProviderConstant{-100})
}
