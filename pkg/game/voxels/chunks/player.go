package chunks

import (
	"log"
	"math"
	"sort"
	"time"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl64"
)

const fallAcceleration = 20.0 // m/s/s
const maxFallSpeed = -30.0    // m/s
const playerCamHeight = 1.8
const playerBoxHeight = 2.0
const playerBoxSize = 0.4 // half the length of a horizontal side
const playerBoxRatio = (2 * playerBoxSize) / playerBoxHeight

type Player struct {
	Camera                  *ant.Camera
	worldUpdater            *ChunkWorldUpdater
	world                   *ChunkWorld
	isFalling               bool
	Velocity                mgl64.Vec3
	inputMovementSuggestion mgl64.Vec3
	Noclip                  bool
	JumpStrength            float64
	debugToggle             bool // debug
}

func NewPlayer(camera *ant.Camera, worldUpdater *ChunkWorldUpdater) *Player {
	return &Player{
		Camera:       camera,
		worldUpdater: worldUpdater,
		world:        worldUpdater.ChunkWorld,
		isFalling:    true,
		Velocity:     mgl64.Vec3{0, 0, 0},
		Noclip:       false,
		JumpStrength: 8.0,
		debugToggle:  false,
	}
}

func (self *Player) SuggestMovement(ds mgl64.Vec3) {
	self.inputMovementSuggestion = ds
}

func (self *Player) Update(dt *time.Duration) {
	if !self.Noclip {
		self.fall(dt)
	}
	ds := self.Velocity.Mul(dt.Seconds())
	translationSuggestion := self.inputMovementSuggestion.Add(ds)
	if translationSuggestion[0] != 0 || translationSuggestion[1] != 0 || translationSuggestion[2] != 0 {
		translation := self.clipFromVoxelCollisions(translationSuggestion)
		self.Camera.Translate(translation)
	}
	self.inputMovementSuggestion = mgl64.Vec3{0, 0, 0}
}

func (self *Player) Jump() {
	if !self.isFalling {
		// log.Println("now jumping") // debug
		self.Velocity = mgl64.Vec3{self.Velocity[0], self.Velocity[1], self.Velocity[2] + self.JumpStrength}
		self.isFalling = true
	}
}

func (self *Player) fall(dt *time.Duration) {
	if self.isFalling {
		dv := dt.Seconds() * fallAcceleration
		newFallSpeed := math.Max(self.Velocity[2]-dv, maxFallSpeed)
		self.Velocity = mgl64.Vec3{self.Velocity[0], self.Velocity[1], newFallSpeed}
	} else {
		moveDown := -0.1 * float64(self.world.ChunkSettings.GetVoxelSize())
		aabbDown := self.getStandingSquare(mgl64.Vec3{0, 0, moveDown})
		// check if the AABB intersects with voxels
		voxels := self.getIntersectingVoxelAABBs(aabbDown)
		standsOnVoxel := len(voxels) > 0
		if !standsOnVoxel {
			self.isFalling = true
		}
	}
}

func (self *Player) clipFromVoxelCollisions(translationSuggestion mgl64.Vec3) mgl64.Vec3 {
	futurePlayerAABB := self.getFutureAABB(translationSuggestion)
	voxels := self.getIntersectingVoxelAABBs(futurePlayerAABB)
	clipped := translationSuggestion
	// determine the translation line for the center of the intersection
	// determine which face of the voxel AABB that vector intersects
	// that face determines which direction needs to be cancelled
	for _, aabb := range voxels {
		intersection := aabb.Intersection(futurePlayerAABB)
		center := intersection.Center()
		lineOrigin := center.Sub(translationSuggestion)
		// determine which face the line intersects
		faces := GetFacingFaces(clipped)
		for _, face := range faces {
			if aabb.LineIntersectsFace(lineOrigin, center, face) {
				clipped = self.cancelComponent(clipped, face)
				if face == UP {
					self.isFalling = false
					self.Velocity = mgl64.Vec3{0, 0, 0}
				}
			}
		}
	}

	return clipped
}

func (self *Player) getIntersectingVoxelAABBs(futureAABB ant.AABB64) []ant.AABB64 {
	var intersections []ant.AABB64
	intersectingChunks := self.getIntersectingChunks(futureAABB)
	// log.Println("futureAABB", futureAABB)
	for _, chunk := range intersectingChunks {
		// we can optimize this; instead of iterating over all voxels in chunk we can calculate min max like we do with getting intersecting chunks
		for index, voxel := range *chunk.Voxels {
			if voxel != AIR {
				voxelAABB := chunk.GetVoxelAABB(index)
				if futureAABB.Intersects(voxelAABB) {
					// log.Println("adding voxel AABB:", coord.ToString(), voxel, voxelMin, voxelMax)
					intersections = append(intersections, voxelAABB)
				}
			}
		}
	}
	return intersections
}

func (self *Player) getFutureAABB(translation mgl64.Vec3) ant.AABB64 {
	destination := self.Camera.Position.Add(translation)
	playerMin := destination.Sub(mgl64.Vec3{playerBoxSize, playerBoxSize, playerCamHeight})
	playerMax := destination.Add(mgl64.Vec3{playerBoxSize, playerBoxSize, playerBoxSize})
	return ant.AABB64{
		Min: playerMin,
		Max: playerMax,
	}
}

func (self *Player) getStandingSquare(translation mgl64.Vec3) ant.AABB64 {
	destination := self.Camera.Position.Add(translation)
	playerMin := destination.Add(mgl64.Vec3{-playerBoxSize, -playerBoxSize, -playerCamHeight})
	playerMax := destination.Add(mgl64.Vec3{playerBoxSize, playerBoxSize, -playerCamHeight})
	return ant.AABB64{
		Min: playerMin,
		Max: playerMax,
	}
}

func (self *Player) getIntersectingChunks(aabb ant.AABB64) []*StandardChunk {
	cMin := self.world.ChunkSettings.GetChunkCoord(aabb.Min)
	cMax := self.world.ChunkSettings.GetChunkCoord(aabb.Max)
	var chunkCoords []IndexCoordinate
	for ci := cMin.i; ci <= cMax.i; ci++ {
		for cj := cMin.j; cj <= cMax.j; cj++ {
			for ck := cMin.k; ck <= cMax.k; ck++ {
				chunkCoords = append(chunkCoords, IndexCoordinate{ci, cj, ck})
			}
		}
	}
	var chunks []*StandardChunk
	for _, coord := range chunkCoords {
		// todo: refactor for higher rank regions
		// todo: wait until column is loaded
		chunk, ok := self.world.Region.GetChunk(coord)
		if ok {
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

// todo remove; just call the one from ChunkSettings, also move unit test
func (self *Player) ToRegionCoord(location mgl64.Vec3) []IndexCoordinate {
	return self.world.ChunkSettings.ToRegionCoord(location)
}

func (self *Player) cancelComponent(thing mgl64.Vec3, face Face) mgl64.Vec3 {
	delta := 0.0
	switch face {
	case EAST:
		return mgl64.Vec3{delta, thing[1], thing[2]}
	case WEST:
		return mgl64.Vec3{-delta, thing[1], thing[2]}
	case NORTH:
		return mgl64.Vec3{thing[0], delta, thing[2]}
	case SOUTH:
		return mgl64.Vec3{thing[0], -delta, thing[2]}
	case UP:
		return mgl64.Vec3{thing[0], thing[1], delta}
	case DOWN:
		return mgl64.Vec3{thing[0], thing[1], -delta}
	}
	return thing
}

// setup chunk sorting by distance
type dChunk struct {
	chunk    *StandardChunk
	distance float64
}
type ByDistance []dChunk

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Less(i, j int) bool { return a[i].distance < a[j].distance }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (self *Player) RemoveBlock() {
	// determine interaction line points
	p1 := self.Camera.Position
	p2 := self.Camera.Position.Add(self.Camera.Direction.Mul(20))
	scaleX := 1.0 / float64(self.world.ChunkSettings.GetChunkWidth())
	scaleY := 1.0 / float64(self.world.ChunkSettings.GetChunkDepth())
	scaleZ := 1.0 / float64(self.world.ChunkSettings.GetChunkHeight())
	// scale line with chunk dimensions
	unitSpace_p1 := mgl64.Vec3{p1[0] * scaleX, p1[1] * scaleY, p1[2] * scaleZ}
	unitSpace_p2 := mgl64.Vec3{p2[0] * scaleX, p2[1] * scaleY, p2[2] * scaleZ}
	cellIntersections := ant.LineCellIntersections(unitSpace_p1, unitSpace_p2)

	// get intersecting chunks
	var coords []IndexCoordinate
	for _, yo := range cellIntersections {
		coords = append(coords, IndexCoordinate{i: yo[0], j: yo[1], k: yo[2]})
	}
	chunks := self.world.Region.GetChunks(coords)

	if len(chunks) == 0 {
		log.Println("no chunks intersect") // debug
		return
	} else {
		log.Println("chunks intersect: ", len(chunks)) // debug
	}

	var dChunks []dChunk

	// calculate distances
	for _, chunk := range chunks {
		dChunks = append(dChunks, dChunk{chunk: chunk, distance: self.GetChunkDistance(chunk.Coordinate)})
	}
	// order chunks by distance
	sort.Sort(ByDistance(dChunks))

	tmin := math.MaxFloat64
	var targetChunk *StandardChunk = nil
	targetVoxelIndex := -1
	for _, dChunk := range dChunks {
		chunk := dChunk.chunk
		// loop over visible voxels in chunk
		for _, vIndex := range *chunk.VisibleVoxels {
			voxel := (*chunk.Voxels)[vIndex]
			if voxel != AIR {
				voxelAABB := chunk.GetVoxelAABB(vIndex)
				intersects, t := voxelAABB.LineIntersects(p1, p2)
				// todo: get interestion face for adding voxels
				if intersects && t < tmin {
					tmin = t
					targetChunk = chunk
					targetVoxelIndex = vIndex
				}
			}
		}
	}

	if tmin != math.MaxFloat64 {
		targetChunk.RemoveVoxel(targetVoxelIndex)
		self.worldUpdater.QueueForRebuild(targetChunk)
	} else {
		log.Println("no voxel intersect") // debug
	}
}

func (self *Player) GetChunkDistance(c IndexCoordinate) float64 {
	sizeX := float64(self.world.ChunkSettings.GetChunkWidth())
	sizeY := float64(self.world.ChunkSettings.GetChunkDepth())
	sizeZ := float64(self.world.ChunkSettings.GetChunkHeight())
	halfX := sizeX / 2
	halfY := sizeY / 2
	halfZ := sizeZ / 2
	chunkPos := mgl64.Vec3{
		float64(c.i)*sizeX + halfX,
		float64(c.j)*sizeY + halfY,
		float64(c.k)*sizeZ + halfZ,
	}
	d := self.Camera.Position.Sub(chunkPos)
	return d[0]*d[0] + d[1]*d[1] + d[2]*d[2]
}
