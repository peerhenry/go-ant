package chunks

type ChunkMeshBuilder struct {
	ChunkSettings IChunkSettings
}

func NewChunkMeshBuilder(chunkSettings IChunkSettings) *ChunkMeshBuilder {
	return &ChunkMeshBuilder{chunkSettings}
}

func (self *ChunkMeshBuilder) ChunkToMesh(chunk *StandardChunk) *ChunkMesh {
	// iterate over visible voxels
	var positions []float32
	var normalIndices []int32
	var uvs []float32
	var indices []uint32
	var indexOffset uint32 = 0
	var indicesCount int32 = 0

	maybeAddFace := func(voxel, vi, vj, vk, di, dj, dk int, face Face) {
		// addFace = this one is water and the other is air
		// or this one is not water and the other one is water or air
		other := chunk.GetVoxel(vi+di, vj+dj, vk+dk)
		shouldAddFace := (voxel != WATER && VoxelIsTransparent(other)) || (voxel == WATER && other == AIR)
		if shouldAddFace {
			nextPositions := self.GetQuadPositions(vi, vj, vk, face)
			positions = append(positions, nextPositions[:]...)

			nextNormals := self.GetQuadNormals(face)
			normalIndices = append(normalIndices, nextNormals[:]...)

			nextUvs := self.GetQuadUvs(voxel, vi, vj, vk, face)
			uvs = append(uvs, nextUvs...)

			nextIndices := []uint32{indexOffset, indexOffset + 1, indexOffset + 2, indexOffset + 2, indexOffset + 1, indexOffset + 3}
			indices = append(indices, nextIndices...)
			indexOffset = indexOffset + 4
			indicesCount += 6
		}
	}

	for _, index := range *chunk.VisibleVoxels {
		v := self.ChunkSettings.IndexToCoordinate(index)
		voxel := (*chunk.Voxels)[index]
		if voxel != UNDERGROUND && voxel != AIR {
			maybeAddFace(voxel, v.i, v.j, v.k, 0, -1, 0, SOUTH)
			maybeAddFace(voxel, v.i, v.j, v.k, 1, 0, 0, EAST)
			maybeAddFace(voxel, v.i, v.j, v.k, 0, 1, 0, NORTH)
			maybeAddFace(voxel, v.i, v.j, v.k, -1, 0, 0, WEST)
			maybeAddFace(voxel, v.i, v.j, v.k, 0, 0, 1, TOP)
			maybeAddFace(voxel, v.i, v.j, v.k, 0, 0, -1, BOTTOM)
		}
	}
	return &ChunkMesh{&positions, &normalIndices, &uvs, &indices, indicesCount}
}

func (self *ChunkMeshBuilder) GetQuadPositions(i, j, k int, face Face) [12]float32 {
	size := self.ChunkSettings.GetVoxelSize()
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

func (self *ChunkMeshBuilder) GetQuadNormals(face Face) [4]int32 {
	return [4]int32{
		int32(face),
		int32(face),
		int32(face),
		int32(face),
	}
}

func (self *ChunkMeshBuilder) GetQuadUvs(voxel, i, j, k int, face Face) []float32 {
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
	case TRUNK:
		switch face {
		case TOP:
			return trunkInner
		case BOTTOM:
			return trunkInner
		default:
			return trunk
		}
	case LEAVES:
		return leavesFull
	case SNOWDIRT:
		switch face {
		case TOP:
			return snow
		case BOTTOM:
			return dirt
		default:
			return snowDirt
		}
	case WATER:
		return water
	default:
		return dirt
	}
}

// =========== voxel UVS ============

const pixelSize float32 = 1.0 / 512
const quadSize float32 = 1.0 / 16

func getCubeUvsAt(i, j byte) []float32 {
	left := float32(i)*quadSize + pixelSize
	right := float32(i+1)*quadSize - pixelSize
	top := float32(j)*quadSize + pixelSize
	bottom := float32(j+1)*quadSize - pixelSize
	return []float32{
		left, bottom,
		left, top,
		right, bottom,
		right, top,
	}
}

var dirt []float32 = getCubeUvsAt(2, 0)
var grassTop []float32 = getCubeUvsAt(0, 0)
var grassSide []float32 = getCubeUvsAt(3, 0)
var stone []float32 = getCubeUvsAt(1, 0)
var planks []float32 = getCubeUvsAt(4, 0)
var sand []float32 = getCubeUvsAt(2, 1)
var trunk []float32 = getCubeUvsAt(4, 1)
var trunkInner []float32 = getCubeUvsAt(5, 1)
var leavesFull []float32 = getCubeUvsAt(5, 3)
var snowDirt []float32 = getCubeUvsAt(4, 4)
var snow []float32 = getCubeUvsAt(2, 4)
var water []float32 = getCubeUvsAt(13, 12)
