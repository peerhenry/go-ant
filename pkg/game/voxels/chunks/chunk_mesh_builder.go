package chunks

type ChunkMeshBuilder struct {
	chunkSettings IChunkSettings
}

func NewChunkMeshBuilder(chunkSettings IChunkSettings) *ChunkMeshBuilder {
	return &ChunkMeshBuilder{chunkSettings}
}

func (self *ChunkMeshBuilder) ChunkToMesh(chunk *StandardChunk) *ChunkMesh {
	// iterate over visible voxels
	var positions []float32
	var normals []float32
	var uvs []float32
	var indices []uint32
	var indexOffset uint32 = 0
	var indicesCount int32 = 0

	maybeAddFace := func(voxel, i, j, k, di, dj, dk, face int) {
		if chunk.IsTransparent(i+di, j+dj, k+dk) {
			nextPositions := self.GetQuadPositions(i, j, k, face)
			positions = append(positions, nextPositions[:]...)

			nextNormals := self.GetQuadNormals(i, j, k, face)
			normals = append(normals, nextNormals[:]...)

			nextUvs := self.GetQuadUvs(voxel, i, j, k, face)
			uvs = append(uvs, nextUvs...)

			nextIndices := []uint32{indexOffset, indexOffset + 1, indexOffset + 2, indexOffset + 2, indexOffset + 1, indexOffset + 3}
			indices = append(indices, nextIndices...)
			indexOffset = indexOffset + 4
			indicesCount += 6
		}
	}

	for _, index := range *chunk.visibleVoxels {
		i, j, k := self.chunkSettings.IndexToCoordinate(index)
		voxel := (*chunk.voxels)[index]
		maybeAddFace(voxel, i, j, k, 0, -1, 0, SOUTH)
		maybeAddFace(voxel, i, j, k, 1, 0, 0, EAST)
		maybeAddFace(voxel, i, j, k, 0, 1, 0, NORTH)
		maybeAddFace(voxel, i, j, k, -1, 0, 0, WEST)
		maybeAddFace(voxel, i, j, k, 0, 0, 1, TOP)
		maybeAddFace(voxel, i, j, k, 0, 0, -1, BOTTOM)
	}
	return &ChunkMesh{&positions, &normals, &uvs, &indices, indicesCount}
}

func (self *ChunkMeshBuilder) GetQuadPositions(i, j, k, face int) [12]float32 {
	size := self.chunkSettings.GetVoxelSize()
	ox := size * float32(i)
	oy := size * float32(j)
	oz := size * float32(k)
	xx := ox + size
	yy := oy + size
	zz := oz + size
	switch face {
	case SOUTH:
		return [12]float32{
			ox, oy, oz,
			ox, oy, zz,
			xx, oy, oz,
			xx, oy, zz,
		}
	case EAST:
		return [12]float32{
			xx, oy, oz,
			xx, oy, zz,
			xx, yy, oz,
			xx, yy, zz,
		}
	case NORTH:
		return [12]float32{
			xx, yy, oz,
			xx, yy, zz,
			ox, yy, oz,
			ox, yy, zz,
		}
	case WEST:
		return [12]float32{
			ox, yy, oz,
			ox, yy, zz,
			ox, oy, oz,
			ox, oy, zz,
		}
	case TOP:
		return [12]float32{
			ox, oy, zz,
			ox, yy, zz,
			xx, oy, zz,
			xx, yy, zz,
		}
	case BOTTOM:
		return [12]float32{
			ox, yy, oz,
			ox, oy, oz,
			xx, yy, oz,
			xx, oy, oz,
		}
	}
	panic("No direction given")
}

func (self *ChunkMeshBuilder) GetQuadNormals(i, j, k, face int) [12]float32 {
	switch face {
	case SOUTH:
		return [12]float32{
			0, -1, 0,
			0, -1, 0,
			0, -1, 0,
			0, -1, 0,
		}
	case EAST:
		return [12]float32{
			1, 0, 0,
			1, 0, 0,
			1, 0, 0,
			1, 0, 0,
		}
	case NORTH:
		return [12]float32{
			0, 1, 0,
			0, 1, 0,
			0, 1, 0,
			0, 1, 0,
		}
	case WEST:
		return [12]float32{
			-1, 0, 0,
			-1, 0, 0,
			-1, 0, 0,
			-1, 0, 0,
		}
	case TOP:
		return [12]float32{
			0, 0, 1,
			0, 0, 1,
			0, 0, 1,
			0, 0, 1,
		}
	case BOTTOM:
		return [12]float32{
			0, 0, -1,
			0, 0, -1,
			0, 0, -1,
			0, 0, -1,
		}
	}
	panic("No direction given")
}

func (self *ChunkMeshBuilder) GetQuadUvs(voxel, i, j, k, face int) []float32 {
	switch voxel {
	case GRASS:
		switch face {
		case TOP:
			return grassTop
		case BOTTOM:
			return dirt
		default:
			return grassSide
		}
	case DIRT:
		return dirt
	case STONE:
		return stone
	case SAND:
		return sand
	default:
		return dirt
	}
}
