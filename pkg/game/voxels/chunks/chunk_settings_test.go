package chunks

import "testing"

func TestCoordinateToIndex(t *testing.T) {
	chunkSettings := CreateStandardChunkSettings(7, 7, 7)
	expecti := 6
	expectj := 5
	expectk := 4
	expectedString := arrayToString([3]int{expecti, expectj, expectk})
	index := chunkSettings.CoordinateToIndexijk(expecti, expectj, expectk)
	resulti, resultj, resultk := chunkSettings.IndexToCoordinateijk(index)
	if resulti != expecti || resultj != expectj || resultk != expectk {
		resultString := arrayToString([3]int{resulti, resultj, resultk})
		t.Errorf("Expected %s but got %s", expectedString, resultString)
	}
}

func TestAddCoordinatei(t *testing.T) {
	// Arrange
	settings := CreateStandardChunkSettings(5, 5, 5)
	root := IndexCoordinate{0, 0, 0}
	coord := []IndexCoordinate{root}
	expectedRoot := IndexCoordinate{1, 0, 0}
	// Act
	result := settings.AddCoordinatei(coord, 1)
	// Assert
	resultRanks := len(result)
	resultRoot := result[0]
	if resultRanks != 1 {
		t.Errorf("Expected 1 rank but got %d", resultRanks)
	}
	if !resultRoot.Equals(expectedRoot) {
		t.Errorf("Expected %s but got %s", expectedRoot.ToString(), resultRoot.ToString())
	}
}

func TestAddCoordinateiOverBoundary(t *testing.T) {
	// Arrange
	settings := CreateStandardChunkSettings(5, 5, 5)
	root := IndexCoordinate{1, 3, 2}
	coord := []IndexCoordinate{root}
	expected0 := IndexCoordinate{3, 3, 2}
	expected1 := IndexCoordinate{1, 0, 0}
	// Act
	result := settings.AddCoordinatei(coord, 7)
	// Assert
	resultRanks := len(result)
	if resultRanks != 2 {
		t.Errorf("Expected 2 ranks but got %d", resultRanks)
	}
	coord0 := result[0]
	if !coord0.Equals(expected0) {
		t.Errorf("Expected coord0 to be %s but got %s", expected0.ToString(), coord0.ToString())
	}
	coord1 := result[1]
	if !coord1.Equals(expected1) {
		t.Errorf("Expected coord1 to be %s but got %s", expected1.ToString(), coord1.ToString())
	}
}

func TestAddCoordinateiOverMultipleBoundaries(t *testing.T) {
	// Arrange
	settings := CreateStandardChunkSettings(5, 5, 5)
	root := IndexCoordinate{1, 0, 0}
	coord := []IndexCoordinate{root}
	expected0 := IndexCoordinate{3, 0, 0}
	expected1 := IndexCoordinate{3, 0, 0}
	// Act
	result := settings.AddCoordinatei(coord, 17)
	// Assert
	resultRanks := len(result)
	if resultRanks != 2 {
		t.Errorf("Expected 2 ranks but got %d", resultRanks)
	}
	coord0 := result[0]
	if !coord0.Equals(expected0) {
		t.Errorf("Expected coord0 to be %s but got %s", expected0.ToString(), coord0.ToString())
	}
	coord1 := result[1]
	if !coord1.Equals(expected1) {
		t.Errorf("Expected coord1 to be %s but got %s", expected1.ToString(), coord1.ToString())
	}
}

func TestAddCoordinateiOnBoundary(t *testing.T) {
	// Arrange
	settings := CreateStandardChunkSettings(5, 5, 5)
	root := IndexCoordinate{3, 0, 0}
	coord := []IndexCoordinate{root}
	expected0 := IndexCoordinate{0, 0, 0}
	expected1 := IndexCoordinate{1, 0, 0}
	// Act
	result := settings.AddCoordinatei(coord, 2)
	// Assert
	resultRanks := len(result)
	if resultRanks != 2 {
		t.Errorf("Expected 2 ranks but got %d", resultRanks)
	}
	coord0 := result[0]
	if !coord0.Equals(expected0) {
		t.Errorf("Expected coord0 to be %s but got %s", expected0.ToString(), coord0.ToString())
	}
	coord1 := result[1]
	if !coord1.Equals(expected1) {
		t.Errorf("Expected coord1 to be %s but got %s", expected1.ToString(), coord1.ToString())
	}
}

func TestAddCoordinateiSubtraction(t *testing.T) {
	// Arrange
	settings := CreateStandardChunkSettings(5, 5, 5)
	root := IndexCoordinate{1, 0, 0}
	coord := []IndexCoordinate{root}
	expected0 := IndexCoordinate{4, 0, 0}
	expected1 := IndexCoordinate{-1, 0, 0}
	// Act
	result := settings.AddCoordinatei(coord, -2)
	// Assert
	resultRanks := len(result)
	if resultRanks != 2 {
		t.Errorf("Expected 2 ranks but got %d", resultRanks)
	}
	coord0 := result[0]
	if !coord0.Equals(expected0) {
		t.Errorf("Expected coord0 to be %s but got %s", expected0.ToString(), coord0.ToString())
	}
	coord1 := result[1]
	if !coord1.Equals(expected1) {
		t.Errorf("Expected coord1 to be %s but got %s", expected1.ToString(), coord1.ToString())
	}
}

func TestAddCoordinateiOverBoundaryWithExistingHigherRanks(t *testing.T) {
	// Arrange
	settings := CreateStandardChunkSettings(5, 5, 5)
	root := IndexCoordinate{1, 3, 2}
	coord := []IndexCoordinate{root, IndexCoordinate{2, 1, 3}}
	expected0 := IndexCoordinate{3, 3, 2}
	expected1 := IndexCoordinate{3, 1, 3}
	// Act
	result := settings.AddCoordinatei(coord, 7)
	// Assert
	resultRanks := len(result)
	if resultRanks != 2 {
		t.Errorf("Expected 2 ranks but got %d", resultRanks)
	}
	coord0 := result[0]
	if !coord0.Equals(expected0) {
		t.Errorf("Expected coord0 to be %s but got %s", expected0.ToString(), coord0.ToString())
	}
	coord1 := result[1]
	if !coord1.Equals(expected1) {
		t.Errorf("Expected coord1 to be %s but got %s", expected1.ToString(), coord1.ToString())
	}
}

func TestAddCoordinateiSubtractionWithExistingHigherRanks(t *testing.T) {
	// Arrange
	settings := CreateStandardChunkSettings(5, 5, 5)
	root := IndexCoordinate{1, 3, 2}
	coord := []IndexCoordinate{root, IndexCoordinate{2, 1, 3}}
	expected0 := IndexCoordinate{4, 3, 2}
	expected1 := IndexCoordinate{1, 1, 3}
	// Act
	result := settings.AddCoordinatei(coord, -2)
	// Assert
	resultRanks := len(result)
	if resultRanks != 2 {
		t.Errorf("Expected 2 ranks but got %d", resultRanks)
	}
	coord0 := result[0]
	if !coord0.Equals(expected0) {
		t.Errorf("Expected coord0 to be %s but got %s", expected0.ToString(), coord0.ToString())
	}
	coord1 := result[1]
	if !coord1.Equals(expected1) {
		t.Errorf("Expected coord1 to be %s but got %s", expected1.ToString(), coord1.ToString())
	}
}
