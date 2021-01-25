package chunks

type ChunkWorldBuilder struct {
	chunkSettings  IChunkSettings
	HeightProvider IHeightProvider
	spawnTrees     bool
	WaterLevel     int
}

func NewChunkWorldBuilder() *ChunkWorldBuilder {
	return &ChunkWorldBuilder{
		chunkSettings:  NewChunkSettings(32, 32, 8),
		HeightProvider: HeightProviderConstant{0},
		spawnTrees:     false,
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
	self.chunkSettings = settings
	return self
}

func (self *ChunkWorldBuilder) UseHeightProvider(provider IHeightProvider) *ChunkWorldBuilder {
	self.HeightProvider = provider
	return self
}

func (self *ChunkWorldBuilder) SetConstantHeight(height int) *ChunkWorldBuilder {
	self.HeightProvider = NewHeightProviderConstant(height)
	return self
}

func (self *ChunkWorldBuilder) Build() *ChunkWorld {
	chunkBuilder := NewChunkBuilder(self.chunkSettings)
	meshBuilder := NewChunkMeshBuilder(self.chunkSettings)
	region := NewChunkRegion()
	return &ChunkWorld{
		Region:                 region,
		ChunkBuilder:           chunkBuilder,
		ChunkRenderDataBuilder: &ChunkRenderDataBuilder{self.chunkSettings, meshBuilder},
		ChunkSettings:          self.chunkSettings,
		HeightAtlas:            self.HeightProvider,
		WaterLevel:             self.WaterLevel,
		SpawnTrees:             self.spawnTrees,
	}
}
