package chunks

type VoxelStructure struct {
	Voxels map[IndexCoordinate]int
}

func GetStandardTree(height int) *VoxelStructure {
	tree := make(map[IndexCoordinate]int)
	if height < 5 {
		return &VoxelStructure{
			Voxels: tree,
		}
	}
	for dk := 0; dk < height; dk++ {
		tree[IndexCoordinate{0, 0, dk}] = TRUNK
	}
	for di := -2; di <= 2; di++ {
		for dj := -2; dj <= 2; dj++ {
			if di != 0 || dj != 0 {
				tree[IndexCoordinate{di, dj, height - 3}] = LEAVES
				tree[IndexCoordinate{di, dj, height - 2}] = LEAVES
			}
		}
	}
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			if di != 0 || dj != 0 {
				tree[IndexCoordinate{di, dj, height - 1}] = LEAVES
			}
		}
	}
	tree[IndexCoordinate{0, 0, height}] = LEAVES
	tree[IndexCoordinate{1, 0, height}] = LEAVES
	tree[IndexCoordinate{-1, 0, height}] = LEAVES
	tree[IndexCoordinate{0, 1, height}] = LEAVES
	tree[IndexCoordinate{0, -1, height}] = LEAVES
	return &VoxelStructure{
		Voxels: tree,
	}
}
