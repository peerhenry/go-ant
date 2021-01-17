package chunks

type HeightAtlasIndex struct {
	i int
	j int
}

type HeightAtlas struct {
	Charts          map[HeightAtlasIndex]*[]int
	chartSize       int
	heightGenerator IHeightGenerator
}

func NewHeightAtlas(size int, generator IHeightGenerator) *HeightAtlas {
	return &HeightAtlas{
		Charts:          make(map[HeightAtlasIndex]*[]int),
		chartSize:       size,
		heightGenerator: generator,
	}
}

func (self *HeightAtlas) GetHeight(ai, aj int) int {
	vi, vj, hmi, hmj := self.GetChartCoordinates(ai, aj)
	chart, ok := self.Charts[HeightAtlasIndex{hmi, hmj}]
	if !ok {
		// generate height map for atlas
		newHeightMap := self.GetChart(hmi, hmj)
		self.Charts[HeightAtlasIndex{hmi, hmj}] = newHeightMap
		return (*newHeightMap)[vi*self.chartSize+vj]
	}
	return (*chart)[vi*self.chartSize+vj]
}

func (self *HeightAtlas) GetChart(hmi, hmj int) *[]int {
	var heights []int
	for vi := 0; vi < self.chartSize; vi++ {
		for vj := 0; vj < self.chartSize; vj++ {
			// absolute voxel i & j in world
			ai := self.chartSize*hmi + vi
			aj := self.chartSize*hmj + vj
			h := self.heightGenerator.GetHeight(ai, aj)
			heights = append(heights, h)
		}
	}
	return &heights
}

// returns
// - voxel i & j on chart
// - atlas index i, j for accessing height chart
func (self *HeightAtlas) GetChartCoordinates(i int, j int) (int, int, int, int) {
	size := self.chartSize
	vi := i
	vj := j
	hmi := 0
	hmj := 0
	if vi >= size {
		vi = i % size
		hmi = i / size
	} else if vi < 0 {
		hmi = (vi+1)/size - 1
		vi = vi - hmi*size
	}
	if vj >= size {
		vj = j % size
		hmj = j / size
	} else if vj < 0 {
		hmj = (vj+1)/size - 1
		vj = vj - hmj*size
	}
	return vi, vj, hmi, hmj
}
