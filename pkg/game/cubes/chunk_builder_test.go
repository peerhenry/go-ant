package cubes

import (
	"fmt"
	"strings"
	"testing"
	// "ant.com/ant/pkg/game/quad"
	// "ant.com/ant/pkg/game/text"
)

func arrayToString(array interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array)), ","), "[]")
}

func TestCoordinateToIndex(t *testing.T) {
	chunkBuilder := NewChunkBuilder(7, 7, 7)
	expecti := 6
	expectj := 5
	expectk := 4
	expectedString := arrayToString([3]int{expecti, expectj, expectk})
	index := chunkBuilder.CoordinateToIndex(expecti, expectj, expectk)
	result := chunkBuilder.IndexToCoordinate(index)
	if result[0] != expecti || result[1] != expectj || result[2] != expectk {
		resultString := arrayToString(result)
		t.Errorf("Expected %s but got %s", expectedString, resultString)
	}
}
