package chunks

import (
	"time"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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
	if k > self.chunkSettings.GetChunkHeight()-5 {
		return DIRT
	}
	return STONE
}

// =============== other shit ===============

func BuildChunkGameObject(position Vec3) *ant.GameObject {
	chunkSettings := CreateStandardChunkSettings(32, 32, 8)
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
