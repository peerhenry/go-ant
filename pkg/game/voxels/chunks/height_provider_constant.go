package chunks

type HeightProviderConstant struct {
	height int
}

func (self HeightProviderConstant) GetHeight(_, _ int) int {
	return self.height
}
