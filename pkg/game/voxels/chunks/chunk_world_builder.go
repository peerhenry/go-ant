package chunks

type ChunkWorldBuilder struct {
	ChunkSettings  IChunkSettings
	HeightProvider IHeightProvider
	spawnTrees     bool
	WaterLevel     int
}

func NewChunkWorldBuilder() *ChunkWorldBuilder {
	return &ChunkWorldBuilder{
		ChunkSettings:  NewChunkSettings(32, 32, 8),
		HeightProvider: HeightProviderConstant{0},
		spawnTrees:     true,
		WaterLevel:     -6,
	}
}

func (self *ChunkWorldBuilder) SpawnTrees(t bool) *ChunkWorldBuilder {
	self.spawnTrees = t
	return self
}

func (self *ChunkWorldBuilder) SetWaterLevel(t int) *ChunkWorldBuilder {
	self.WaterLevel = t
	return self
}

func (self *ChunkWorldBuilder) UseChunkSettings(settings IChunkSettings) *ChunkWorldBuilder {
	self.ChunkSettings = settings
	return self
}

func (self *ChunkWorldBuilder) UseHeightProvider(provider IHeightProvider) *ChunkWorldBuilder {
	self.HeightProvider = provider
	return self
}

func (self *ChunkWorldBuilder) Build() *ChunkWorld {
	chunkBuilder := NewChunkBuilder(self.ChunkSettings)
	meshBuilder := NewChunkMeshBuilder(self.ChunkSettings)
	return &ChunkWorld{
		Region:                 NewChunkRegion(),
		ChunkBuilder:           chunkBuilder,
		ChunkRenderDataBuilder: &ChunkRenderDataBuilder{self.ChunkSettings, meshBuilder},
		ChunkSettings:          self.ChunkSettings,
		HeightAtlas:            self.HeightProvider,
		WaterLevel:             self.WaterLevel,
		SpawnTrees:             self.spawnTrees,
	}
}
