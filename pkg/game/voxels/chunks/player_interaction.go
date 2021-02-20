package chunks

import (
	"log"
	"math"
	"sort"

	"ant.com/ant/pkg/ant"
	"github.com/go-gl/mathgl/mgl64"
)

type dChunk struct {
	chunk    *StandardChunk
	distance float64
}
type ByDistance []dChunk

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Less(i, j int) bool { return a[i].distance < a[j].distance }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func RemoveBlock(player *Player) {
	intersectionEvent, ok := GetIntersectionEvent(player)
	if ok {
		intersectionEvent.Chunk.RemoveVoxel(intersectionEvent.VoxelIndex)
		player.worldUpdater.QueueForRebuild(intersectionEvent.Chunk)
		// get adjacent chunks index coordinates with voxelindex
		adjacents := GetAdjacentChunks(player.World.ChunkSettings, intersectionEvent)
		for _, yo := range adjacents {
			chunk, ok := player.World.Region.GetChunk(yo.ChunkCoordinate)
			if ok {
				player.worldUpdater.QueueForRebuild(chunk)
			} else {
				// check if chunk coordinate is underground and if yes, create a new one
			}
		}
	} else {
		log.Println("no voxel intersect") // debug
	}
}

type ChunkAndVoxelCoordinate struct {
	ChunkCoordinate IndexCoordinate
	VoxelCoordinate IndexCoordinate
}

type IntersectionEvent struct {
	Chunk      *StandardChunk
	VoxelIndex int
	Face       Face
}

func GetIntersectionEvent(player *Player) (*IntersectionEvent, bool) {
	// determine interaction line points
	p1 := player.Camera.Position
	p2 := player.Camera.Position.Add(player.Camera.Direction.Mul(20))
	scaleX := 1.0 / float64(player.World.ChunkSettings.GetChunkWidth())
	scaleY := 1.0 / float64(player.World.ChunkSettings.GetChunkDepth())
	scaleZ := 1.0 / float64(player.World.ChunkSettings.GetChunkHeight())
	// scale line with chunk dimensions
	unitSpace_p1 := mgl64.Vec3{p1[0] * scaleX, p1[1] * scaleY, p1[2] * scaleZ}
	unitSpace_p2 := mgl64.Vec3{p2[0] * scaleX, p2[1] * scaleY, p2[2] * scaleZ}
	cellIntersections := ant.LineCellIntersections(unitSpace_p1, unitSpace_p2)

	// get intersecting chunks
	var coords []IndexCoordinate
	for _, yo := range cellIntersections {
		coords = append(coords, IndexCoordinate{i: yo[0], j: yo[1], k: yo[2]})
	}
	chunks := player.World.Region.GetChunks(coords)

	if len(chunks) == 0 {
		log.Println("no chunks intersect") // debug
		return nil, false
	} else {
		log.Println("chunks intersect: ", len(chunks)) // debug
	}

	var dChunks []dChunk

	// calculate distances
	for _, chunk := range chunks {
		dChunks = append(dChunks, dChunk{chunk: chunk, distance: GetChunkDistance(player, chunk.Coordinate)})
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
		if targetVoxelIndex != -1 {
			break
		}
	}
	if targetVoxelIndex == -1 {
		return nil, false
	}
	event := IntersectionEvent{
		Chunk:      targetChunk,
		VoxelIndex: targetVoxelIndex,
		Face:       -1, // todo; get interesecting face
	}
	return &event, true
}

func GetChunkDistance(player *Player, c IndexCoordinate) float64 {
	sizeX := float64(player.World.ChunkSettings.GetChunkWidth())
	sizeY := float64(player.World.ChunkSettings.GetChunkDepth())
	sizeZ := float64(player.World.ChunkSettings.GetChunkHeight())
	halfX := sizeX / 2
	halfY := sizeY / 2
	halfZ := sizeZ / 2
	chunkPos := mgl64.Vec3{
		float64(c.i)*sizeX + halfX,
		float64(c.j)*sizeY + halfY,
		float64(c.k)*sizeZ + halfZ,
	}
	d := player.Camera.Position.Sub(chunkPos)
	return d[0]*d[0] + d[1]*d[1] + d[2]*d[2]
}

// todo: is this a utility function that should be extracted from this file?
func GetAdjacentChunks(settings IChunkSettings, event *IntersectionEvent) []ChunkAndVoxelCoordinate {
	maxi := settings.GetChunkWidth() - 1
	maxj := settings.GetChunkDepth() - 1
	maxk := settings.GetChunkHeight() - 1
	coord := settings.IndexToCoordinate(event.VoxelIndex)
	var adjacents []ChunkAndVoxelCoordinate
	if coord.i == 0 {
		adjacents = append(adjacents, ChunkAndVoxelCoordinate{
			ChunkCoordinate: event.Chunk.Coordinate.Addijk(-1, 0, 0),
			VoxelCoordinate: coord.SetI(maxi),
		})
	} else if coord.i == maxi {
		adjacents = append(adjacents, ChunkAndVoxelCoordinate{
			ChunkCoordinate: event.Chunk.Coordinate.Addijk(1, 0, 0),
			VoxelCoordinate: coord.SetI(0),
		})
	}
	if coord.j == 0 {
		adjacents = append(adjacents, ChunkAndVoxelCoordinate{
			ChunkCoordinate: event.Chunk.Coordinate.Addijk(0, -1, 0),
			VoxelCoordinate: coord.SetJ(maxj),
		})
	} else if coord.j == maxj {
		adjacents = append(adjacents, ChunkAndVoxelCoordinate{
			ChunkCoordinate: event.Chunk.Coordinate.Addijk(0, 1, 0),
			VoxelCoordinate: coord.SetJ(0),
		})
	}
	if coord.k == 0 {
		adjacents = append(adjacents, ChunkAndVoxelCoordinate{
			ChunkCoordinate: event.Chunk.Coordinate.Addijk(0, 0, -1),
			VoxelCoordinate: coord.SetK(maxk),
		})
	} else if coord.k == maxk {
		adjacents = append(adjacents, ChunkAndVoxelCoordinate{
			ChunkCoordinate: event.Chunk.Coordinate.Addijk(0, 0, 1),
			VoxelCoordinate: coord.SetK(0),
		})
	}
	return adjacents
}
