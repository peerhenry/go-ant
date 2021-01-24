package chunks

type HeightProviderConstant struct {
	height int
}

func NewHeightProviderConstant(height int) *HeightProviderConstant {
	return &HeightProviderConstant{height}
}

func (self HeightProviderConstant) GetHeight(_, _ int) int {
	return self.height
}
