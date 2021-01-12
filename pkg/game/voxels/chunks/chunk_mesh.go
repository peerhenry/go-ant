package chunks

type ChunkMesh struct {
	positions    *[]float32
	normals      *[]float32
	uvs          *[]float32
	indices      *[]uint32
	indicesCount int32
}
