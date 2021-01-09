package cubes

import (
	"time"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// =============== chunk settings ===============

type IChunkSettings interface {
	GetChunkWidth() int
	GetChunkDepth() int
	GetChunkHeight() int
	CoordinateToIndex(i, j, k int) int
	IndexToCoordinate(index int) (int, int, int)
	CoordinateIsOutOfBounds(i, j, k int) bool
	IndexIsOutOfBounds(index int) bool
	GetVoxelSize() float32
}

type StandardChunkSettings struct {
	chunkWidth            int
	chunkDepth            int
	chunkHeight           int
	chunkDepthTimesHeight int
	voxelSize             float32
}

func CreateStandardChunkSettings(chunkWidth, chunkDepth, chunkHeight int) *StandardChunkSettings {
	return &StandardChunkSettings{
		chunkWidth:            chunkWidth,
		chunkDepth:            chunkDepth,
		chunkHeight:           chunkHeight,
		chunkDepthTimesHeight: chunkDepth * chunkHeight,
		voxelSize:             1.0,
	}
}

func (self *StandardChunkSettings) CoordinateToIndex(i, j, k int) int {
	return self.chunkDepthTimesHeight*i + self.chunkHeight*j + k
}

func (self *StandardChunkSettings) IndexToCoordinate(index int) (int, int, int) {
	k := index % self.chunkHeight
	j := (index % self.chunkDepthTimesHeight) / self.chunkHeight
	i := index / self.chunkDepthTimesHeight
	return i, j, k
}

func (self *StandardChunkSettings) GetChunkWidth() int {
	return self.chunkWidth
}

func (self *StandardChunkSettings) GetChunkDepth() int {
	return self.chunkDepth
}

func (self *StandardChunkSettings) GetChunkHeight() int {
	return self.chunkHeight
}

func (self *StandardChunkSettings) CoordinateIsOutOfBounds(i, j, k int) bool {
	return i < 0 || i >= self.chunkWidth || j < 0 || j >= self.chunkDepth || k < 0 || k >= self.chunkHeight
}

func (self *StandardChunkSettings) IndexIsOutOfBounds(index int) bool {
	i, j, k := self.IndexToCoordinate(index)
	return self.CoordinateIsOutOfBounds(i, j, k)
}

func (self *StandardChunkSettings) GetVoxelSize() float32 {
	return self.voxelSize
}

// =============== chunk builder ===============

type IChunkBuilder interface {
	CreateChunkData() StandardChunk
}

type ChunkBuilder struct {
	chunkSettings IChunkSettings
}

func CreateStandardChunkBuilder(chunkSettings IChunkSettings) *ChunkBuilder {
	return &ChunkBuilder{
		chunkSettings,
	}
}

func (self *ChunkBuilder) CreateChunk() *StandardChunk {
	var chunkVoxels []int
	var visibleVoxels []int
	chunkWidth := self.chunkSettings.GetChunkWidth()
	chunkDepth := self.chunkSettings.GetChunkDepth()
	chunkHeight := self.chunkSettings.GetChunkHeight()
	for i := 0; i < chunkWidth; i++ {
		for j := 0; j < self.chunkSettings.GetChunkDepth(); j++ {
			for k := 0; k < self.chunkSettings.GetChunkHeight(); k++ {
				chunkVoxels = append(chunkVoxels, self.getVoxel(i, j, k))
				if i == 0 || i == chunkWidth-1 || j == 0 || j == chunkDepth-1 || k == 0 || k == chunkHeight-1 {
					index := self.chunkSettings.CoordinateToIndex(i, j, k)
					visibleVoxels = append(visibleVoxels, index)
				}
			}
		}
	}
	return &StandardChunk{
		voxels:        &chunkVoxels,
		visibleVoxels: &visibleVoxels,
		chunkSettings: self.chunkSettings,
	}
}

// todo: injecting this function
func (self *ChunkBuilder) getVoxel(i, j, k int) int {
	if k == (self.chunkSettings.GetChunkHeight() - 1) {
		return GRASS
	}
	if k < 10 {
		return STONE
	}
	return DIRT
}

// =============== chunk ===============

type IsVoxelTransparent func(i, j, k int) bool

// type IChunk interface {
// 	IsTransparent(i, j, k int) bool
// }

type StandardChunk struct {
	voxels        *[]int
	visibleVoxels *[]int
	chunkSettings IChunkSettings
}

func (self *StandardChunk) IsTransparent(i, j, k int) bool {
	// todo: check adjacent chunk
	if self.chunkSettings.CoordinateIsOutOfBounds(i, j, k) {
		return true
	}
	index := self.chunkSettings.CoordinateToIndex(i, j, k)
	return (*self.voxels)[index] == AIR
}

// =============== chunk mesh builder ===============

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

// =============== chunk mesh ===============

type ChunkMesh struct {
	positions    *[]float32
	normals      *[]float32
	uvs          *[]float32
	indices      *[]uint32
	indicesCount int32
}

// =============== other shit ===============

func BuildChunkGameObject() *ant.GameObject {
	chunkSettings := CreateStandardChunkSettings(32, 32, 32)
	chunkBuilder := CreateStandardChunkBuilder(chunkSettings)
	chunk := chunkBuilder.CreateChunk()
	meshBuilder := NewChunkMeshBuilder(chunkSettings)
	mesh := meshBuilder.ChunkToMesh(chunk)
	vaoBuilder := new(ant.VaoBuilder)
	vaoBuilder.AddVertexBuffer(0, 3, mesh.positions)
	vaoBuilder.AddVertexBuffer(1, 3, mesh.normals)
	vaoBuilder.AddVertexBuffer(2, 2, mesh.uvs)
	vaoBuilder.AddIndexBuffer(mesh.indices)
	vao := vaoBuilder.Build()
	position := Vec3{0, 20, 0}
	indicesCount := mesh.indicesCount
	return &ant.GameObject{
		Update: func(dt *time.Duration) {},
		Draw: func(uniformStore *ant.UniformStore) {
			viewMatrix := uniformStore.GetMat4("ViewMatrix")
			projectionMatrix := uniformStore.GetMat4("ProjectionMatrix")
			// calculate uniforms
			var positionMatrix mgl32.Mat4 = mgl32.Translate3D(position.X(), position.Y(), position.Z())
			var orientationMatrix mgl32.Mat4 = mgl32.Ident4()
			modelMatrix := positionMatrix.Mul4(orientationMatrix)
			modelView := viewMatrix.Mul4(modelMatrix)
			normalMatrix := modelView.Mat3()
			mvp := projectionMatrix.Mul4(modelView)
			// set uiniforms
			uniformStore.UniformMat4("ModelViewMatrix", modelView)
			uniformStore.UniformMat3("NormalMatrix", normalMatrix)
			uniformStore.UniformMat4("MVP", mvp)
			// draw
			gl.BindVertexArray(vao)
			gl.DrawElements(gl.TRIANGLES, indicesCount, gl.UNSIGNED_INT, nil)
		},
	}
}

const (
	NORTH  = 1
	EAST   = 2
	SOUTH  = 3
	WEST   = 4
	TOP    = 5
	BOTTOM = 6
)
