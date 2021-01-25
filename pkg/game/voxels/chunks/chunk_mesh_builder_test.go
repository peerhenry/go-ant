package chunks

import (
	"testing"
	// "ant.com/ant/pkg/game/quad"
	// "ant.com/ant/pkg/game/text"
)

func TestChunkToMeshLengths(t *testing.T) {
	// Arrange
	chunkSettings := NewChunkSettings(2, 2, 2)
	world := NewChunkWorldBuilder().UseChunkSettings(chunkSettings).SetConstantHeight(-100).Build()
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	chunk := &StandardChunk{
		Voxels:        &[]int{1, 1, 1, 1, 1, 1, 1, 1},
		VisibleVoxels: &[]int{0, 1, 2, 3, 4, 5, 6, 7},
		ChunkWorld:    world,
	}
	// Act
	result := meshBuilder.ChunkToMesh(chunk)
	// Assert
	// 8*3*4*3 = 288 (8 voxels, 3 visible faces per voxel, 4 vertices per face, 3 floats per vertex)
	expected := 8 * 3 * 4 * 3
	posLength := len(*result.positions)
	if posLength != expected {
		t.Errorf("Expected len(result.positions) to be %d but got: %d", expected, posLength)
	}

	expectedNormals := 8 * 3 * 4
	normalsLength := len(*result.normalIndices)
	if normalsLength != expectedNormals {
		t.Errorf("Expected len(result.normals) to be %d but got: %d", expectedNormals, normalsLength)
	}
	// uvs got 2 floats per vertex so:
	// 8*3*4*2 = 192
	expectedUvs := 8 * 3 * 4 * 2
	uvsLength := len(*result.uvs)
	if uvsLength != expectedUvs {
		t.Errorf("Expected len(result.uvs) to be %d but got: %d", expectedUvs, uvsLength)
	}
	// 6 indices per face
	// 8 * 3 * 6 = 144
	expectedIndices := 8 * 3 * 6
	indicesLength := len(*result.indices)
	if indicesLength != expectedIndices {
		t.Errorf("Expected len(result.positions) to be %d but got: %d", expectedIndices, indicesLength)
	}
	if result.indicesCount != int32(expectedIndices) {
		t.Errorf("Expected indicesCount to be %d but got: %d", expectedIndices, indicesLength)
	}
}

func TestGetQuadPositionsForOriginSouth(t *testing.T) {
	chunkSettings := NewChunkSettings(7, 7, 7)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	pos := meshBuilder.GetQuadPositions(0, 0, 0, SOUTH)
	AssertOffset(t, 0, 0, 0, 0, &pos, chunkSettings)
	AssertOffset(t, 0, 0, 1, 1, &pos, chunkSettings)
	AssertOffset(t, 1, 0, 0, 2, &pos, chunkSettings)
	AssertOffset(t, 1, 0, 1, 3, &pos, chunkSettings)
}

func TestGetQuadPositionsForOriginEast(t *testing.T) {
	chunkSettings := NewChunkSettings(7, 7, 7)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	pos := meshBuilder.GetQuadPositions(0, 0, 0, EAST)
	AssertOffset(t, 1, 0, 0, 0, &pos, chunkSettings)
	AssertOffset(t, 1, 0, 1, 1, &pos, chunkSettings)
	AssertOffset(t, 1, 1, 0, 2, &pos, chunkSettings)
	AssertOffset(t, 1, 1, 1, 3, &pos, chunkSettings)
}

func TestGetQuadPositionsForOriginNorth(t *testing.T) {
	chunkSettings := NewChunkSettings(7, 7, 7)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	pos := meshBuilder.GetQuadPositions(0, 0, 0, NORTH)
	AssertOffset(t, 1, 1, 0, 0, &pos, chunkSettings)
	AssertOffset(t, 1, 1, 1, 1, &pos, chunkSettings)
	AssertOffset(t, 0, 1, 0, 2, &pos, chunkSettings)
	AssertOffset(t, 0, 1, 1, 3, &pos, chunkSettings)
}

func TestGetQuadPositionsForOriginWest(t *testing.T) {
	chunkSettings := NewChunkSettings(7, 7, 7)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	pos := meshBuilder.GetQuadPositions(0, 0, 0, WEST)
	AssertOffset(t, 0, 1, 0, 0, &pos, chunkSettings)
	AssertOffset(t, 0, 1, 1, 1, &pos, chunkSettings)
	AssertOffset(t, 0, 0, 0, 2, &pos, chunkSettings)
	AssertOffset(t, 0, 0, 1, 3, &pos, chunkSettings)
}

func TestGetQuadPositionsForOriginTop(t *testing.T) {
	chunkSettings := NewChunkSettings(7, 7, 7)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	pos := meshBuilder.GetQuadPositions(0, 0, 0, TOP)
	AssertOffset(t, 0, 0, 1, 0, &pos, chunkSettings)
	AssertOffset(t, 0, 1, 1, 1, &pos, chunkSettings)
	AssertOffset(t, 1, 0, 1, 2, &pos, chunkSettings)
	AssertOffset(t, 1, 1, 1, 3, &pos, chunkSettings)
}

func TestGetQuadPositionsForOriginBottom(t *testing.T) {
	chunkSettings := NewChunkSettings(7, 7, 7)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	pos := meshBuilder.GetQuadPositions(0, 0, 0, BOTTOM)
	AssertOffset(t, 0, 1, 0, 0, &pos, chunkSettings)
	AssertOffset(t, 0, 0, 0, 1, &pos, chunkSettings)
	AssertOffset(t, 1, 1, 0, 2, &pos, chunkSettings)
	AssertOffset(t, 1, 0, 0, 3, &pos, chunkSettings)
}

func TestGetQuadPositionsBottom(t *testing.T) {
	chunkSettings := NewChunkSettings(7, 7, 7)
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	pos := meshBuilder.GetQuadPositions(1, 2, 4, BOTTOM)
	AssertOffset(t, 1, 3, 4, 0, &pos, chunkSettings)
	AssertOffset(t, 1, 2, 4, 1, &pos, chunkSettings)
	AssertOffset(t, 2, 3, 4, 2, &pos, chunkSettings)
	AssertOffset(t, 2, 2, 4, 3, &pos, chunkSettings)
}

func AssertOffset(t *testing.T, expectX, expectY, expectZ float32, number int, pos *[12]float32, chunkSettings IChunkSettings) {
	size := chunkSettings.GetVoxelSize()
	offset := number * 3
	vertexExpected := [3]float32{size * expectX, size * expectY, size * expectZ}
	vertexIsOk := pos[offset] == vertexExpected[0] && pos[offset+1] == vertexExpected[1] && pos[offset+2] == vertexExpected[2]
	if !vertexIsOk {
		t.Errorf("Vertex failed: expected %s, got %s",
			arrayToString(vertexExpected),
			arrayToString([3]float32{pos[offset], pos[offset+1], pos[offset+2]}))
	}
}
