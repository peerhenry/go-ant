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
			uvs = append(uvs, nextUvs[:]...)

			nextIndices := []uint32{indexOffset, indexOffset + 1, indexOffset + 2, indexOffset + 2, indexOffset + 1, indexOffset + 3}
			indices = append(indices, nextIndices...)
			indexOffset = indexOffset + 4
			indicesCount += 6
		}
	}

	for index := range chunk.VisibleVoxels {
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

func (self *ChunkMeshBuilder) GetQuadUvs(voxel, i, j, k int, face Face) [8]float32 {
	switch voxel {
	case GRASS:
		switch face {
		case TOP:
			return uvs_grassTop
		case BOTTOM:
			return uvs_dirt
		default:
			return uvs_grassSide
		}
	case DIRT:
		return uvs_dirt
	case STONE:
		return uvs_stone
	case SAND:
		return uvs_sand
	case TRUNK:
		switch face {
		case TOP:
			return uvs_trunkInner
		case BOTTOM:
			return uvs_trunkInner
		default:
			return uvs_trunk
		}
	case LEAVES:
		return uvs_leavesFull
	case SNOWDIRT:
		switch face {
		case TOP:
			return uvs_snow
		case BOTTOM:
			return uvs_dirt
		default:
			return uvs_snowDirt
		}
	case WATER:
		return uvs_water

	case RED_FLOWER:
		return uvs_redFlower
	case YELLOW_FLOWER:
		return uvs_yellowFlower
	case RED_MUSHROOM:
		return uvs_redMushroom
	case BROWN_MUSHROOM:
		return uvs_brownMushroom
	case GRASS_1:
		return uvs_grass_1
	case GRASS_2:
		return uvs_grass_2
	case GRASS_3:
		return uvs_grass_3
	case GRASS_4:
		return uvs_grass_4
	case GRASS_5:
		return uvs_grass_5
	case GRASS_6:
		return uvs_grass_6
	case GRASS_7:
		return uvs_grass_7
	case GRASS_8:
		return uvs_grass_8

	default:
		return uvs_dirt
	}
}

// =========== voxel UVS ============

const pixelSize float32 = 1.0 / 512
const quadSize float32 = 1.0 / 16

func getCubeUvsAt(i, j byte) [8]float32 {
	left := float32(i)*quadSize + pixelSize
	right := float32(i+1)*quadSize - pixelSize
	top := float32(j)*quadSize + pixelSize
	bottom := float32(j+1)*quadSize - pixelSize
	return [8]float32{
		left, bottom,
		left, top,
		right, bottom,
		right, top,
	}
}

type QuadUvs = [8]float32

var uvs_dirt QuadUvs = getCubeUvsAt(2, 0)
var uvs_grassTop QuadUvs = getCubeUvsAt(0, 0)
var uvs_grassSide QuadUvs = getCubeUvsAt(3, 0)
var uvs_stone QuadUvs = getCubeUvsAt(1, 0)
var uvs_planks QuadUvs = getCubeUvsAt(4, 0)
var uvs_sand QuadUvs = getCubeUvsAt(2, 1)
var uvs_trunk QuadUvs = getCubeUvsAt(4, 1)
var uvs_trunkInner QuadUvs = getCubeUvsAt(5, 1)
var uvs_leavesFull QuadUvs = getCubeUvsAt(5, 3)
var uvs_snowDirt QuadUvs = getCubeUvsAt(4, 4)
var uvs_snow QuadUvs = getCubeUvsAt(2, 4)
var uvs_water QuadUvs = getCubeUvsAt(13, 12)
var uvs_redFlower QuadUvs = getCubeUvsAt(12, 0)
var uvs_yellowFlower QuadUvs = getCubeUvsAt(13, 0)
var uvs_redMushroom QuadUvs = getCubeUvsAt(12, 1)
var uvs_brownMushroom QuadUvs = getCubeUvsAt(13, 1)
var uvs_grass_1 QuadUvs = getCubeUvsAt(8, 5)
var uvs_grass_2 QuadUvs = getCubeUvsAt(9, 5)
var uvs_grass_3 QuadUvs = getCubeUvsAt(10, 5)
var uvs_grass_4 QuadUvs = getCubeUvsAt(11, 5)
var uvs_grass_5 QuadUvs = getCubeUvsAt(12, 5)
var uvs_grass_6 QuadUvs = getCubeUvsAt(13, 5)
var uvs_grass_7 QuadUvs = getCubeUvsAt(14, 5)
var uvs_grass_8 QuadUvs = getCubeUvsAt(15, 5)
