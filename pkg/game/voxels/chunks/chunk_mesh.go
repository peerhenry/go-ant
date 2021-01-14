package chunks

type ChunkMesh struct {
	positions     *[]float32
	normalIndices *[]int32
	uvs           *[]float32
	indices       *[]uint32
	indicesCount  int32
}
