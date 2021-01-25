package chunks

import (
	"math"
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
	world                   *ChunkWorld
	isFalling               bool
	Velocity                mgl64.Vec3
	inputMovementSuggestion mgl64.Vec3
	Noclip                  bool
	JumpStrength            float64
	debugToggle             bool // debug
}

func NewPlayer(camera *ant.Camera, world *ChunkWorld) *Player {
	return &Player{
		Camera:       camera,
		world:        world,
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
	self.fall(dt)
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
	settings := self.world.ChunkSettings
	voxelSize := float64(settings.GetVoxelSize())
	// log.Println("futureAABB", futureAABB)
	for _, chunk := range intersectingChunks {
		chunkOrigin := chunk.Origin()
		// we can optimize this; instead of iterating over all voxels in chunk we can calculate min max like we do with getting intersecting chunks
		for index, voxel := range *chunk.Voxels {
			if voxel != AIR {
				coord := settings.IndexToCoordinate(index)
				positionInChunk := mgl64.Vec3{
					float64(coord.i) * voxelSize,
					float64(coord.j) * voxelSize,
					float64(coord.k) * voxelSize,
				}
				voxelMin := chunkOrigin.Add(positionInChunk)
				voxelMax := voxelMin.Add(mgl64.Vec3{voxelSize, voxelSize, voxelSize})
				voxelAABB := ant.AABB64{Min: voxelMin, Max: voxelMax}
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

// backup because I am scared to delete the deprecated code in one go

// This does not work :/
// Because the sextant is determined based on centers of AABBs, not the point of collision and that causes errors
// suggestion:
// Calculate the intersection box, take the center point of that
// construct a vector of that point with the suggested translation
// use this vector to determine which face of the target AABB it is intersecting: that will be the cancel direction!
// func (self *Player) GetSextant(playerCenter mgl64.Vec3, aabb ant.AABB64) Face {
// 	aabbCenter := mgl64.Vec3{
// 		(aabb.Min[0] + aabb.Max[0]) / 2.0,
// 		(aabb.Min[1] + aabb.Max[1]) / 2.0,
// 		(aabb.Min[2] + aabb.Max[2]) / 2.0,
// 	}
// 	delta := aabbCenter.Sub(playerCenter)
// 	if math.Abs(delta[0]) > math.Abs(delta[1]) {
// 		if math.Abs(delta[0]) > playerBoxRatio*math.Abs(delta[2]) {
// 			// x direction
// 			if delta[0] > 0 {
// 				return EAST
// 			} else {
// 				return WEST
// 			}
// 		} else {
// 			// z direction
// 			if delta[2] > 0 {
// 				return TOP
// 			} else {
// 				return BOTTOM
// 			}
// 		}
// 	} else {
// 		if math.Abs(delta[1]) > playerBoxRatio*math.Abs(delta[2]) {
// 			// y direction
// 			if delta[1] > 0 {
// 				return NORTH
// 			} else {
// 				return SOUTH
// 			}
// 		} else {
// 			// z direction
// 			if delta[2] > 0 {
// 				return TOP
// 			} else {
// 				return BOTTOM
// 			}
// 		}
// 	}
// }

// func (self *Player) clipFromVoxelCollisions(translationSuggestion mgl64.Vec3) mgl64.Vec3 {
// 	// get player future AABB using nextTranslation
// 	// futureAABB := self.getFutureAABB(translationSuggestion)
// 	// voxels := self.getIntersectingVoxelAABBs(futureAABB)
// 	// camPos := self.Camera.Position
// 	// origin & destination are player center
// 	// correction := playerCamHeight - playerBoxHeight/2.0
// 	// origin := mgl64.Vec3{camPos[0], camPos[1], camPos[2] - correction}
// 	// destination := origin.Add(translationSuggestion)
// 	// For intersecting voxels, get the sextant relative to both origin and destination where the center of the aabb resides
// 	// log.Println("intersecting voxels count", len(voxels)) // debug
// 	// log.Println("futureAABB", futureAABB)                 // debug
// 	clipped := translationSuggestion

// 	// ===== deprecated, it doesnt work properly =====
// 	// for _, aabb := range voxels {
// 	// 	// log.Println("voxel AABB", aabb) // debug
// 	// 	originSextant := self.GetSextant(origin, aabb)
// 	// 	destinationSextant := self.GetSextant(destination, aabb)
// 	// 	cancelDirection := originSextant
// 	// 	if originSextant != destinationSextant {
// 	// 		// cancel move direction in one direction and AABB test again, if it fails, cancel the other Sextant
// 	// 		newSuggestion := self.cancelComponent(translationSuggestion, cancelDirection)
// 	// 		newFutureAABB := self.getFutureAABB(newSuggestion)
// 	// 		stillIntersects := newFutureAABB.Intersects(aabb)
// 	// 		if stillIntersects {
// 	// 			cancelDirection = destinationSextant
// 	// 		}
// 	// 	}
// 	// 	if !self.debugToggle {
// 	// 		// log.Println("collideWithVoxels cancelling", DirectionToString(cancelDirection))
// 	// 		// log.Println("originSextant", DirectionToString(originSextant))
// 	// 		// log.Println("destinationSextant", DirectionToString(destinationSextant))
// 	// 		// self.debugToggle = true
// 	// 	}
// 	// 	if cancelDirection == BOTTOM {
// 	// 		self.isFalling = false
// 	// 	}
// 	// 	clipped = self.cancelComponent(clipped, cancelDirection)
// 	// }
// 	// ============================================================

// 	// determine the translation vector for the center of the intersection
// 	// determine which face of the voxel AABB that vector intersects
// 	// that face determines which direction needs to be cancelled
// 	for _, aabb := range voxels {
// 		intersection := aabb.Intersection(futureAABB)
// 		center := intersection.Center()
// 		lineOrigin := center.Sub(clipped)
// 		// determine which face the ray intersects
// 		faces := GetFacingFaces(clipped)
// 		for _, face := range faces {
// 			aabb.IntersectsFace(lineOrigin, center, face)
// 			clipped = self.cancelComponent(clipped, cancelDirection)
// 		}
// 	}

// 	return clipped
// }
