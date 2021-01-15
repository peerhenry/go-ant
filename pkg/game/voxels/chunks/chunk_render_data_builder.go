package chunks

import (
	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl32"
)

type ChunkRenderDataBuilder struct {
	chunkSettings IChunkSettings
	meshBuilder   *ChunkMeshBuilder
}

func (self *ChunkRenderDataBuilder) ChunkToRenderData(chunk *StandardChunk) *ant.RenderData {
	position := Vec3{
		float32(self.chunkSettings.GetChunkWidth() * chunk.Coordinate.i),
		float32(self.chunkSettings.GetChunkDepth() * chunk.Coordinate.j),
		float32(self.chunkSettings.GetChunkHeight() * chunk.Coordinate.k),
	}
	mesh := self.meshBuilder.ChunkToMesh(chunk)
	vaoBuilder := new(ant.VaoBuilder)
	vaoBuilder.AddVertexBuffer(0, 3, mesh.positions)
	vaoBuilder.AddIntegerBuffer(1, 1, mesh.normalIndices)
	vaoBuilder.AddVertexBuffer(2, 2, mesh.uvs)
	vaoBuilder.AddIndexBuffer(mesh.indices)
	vao := vaoBuilder.Build()
	indicesCount := mesh.indicesCount
	var positionMatrix mgl32.Mat4 = mgl32.Translate3D(position.X(), position.Y(), position.Z())
	return &ant.RenderData{
		Transform:    positionMatrix,
		Vao:          vao,
		IndicesCount: indicesCount,
	}
}
