package chunks

import (
	"log"
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
		// log.Println("intersection with chunk:", intersectionEvent.Chunk.Coordinate.ToString())        // debug
		// log.Println("and voxel:", intersectionEvent.VoxelIndex)                                       // debug
		// log.Println("the voxel is:", (*intersectionEvent.Chunk.Voxels)[intersectionEvent.VoxelIndex]) // debug
		intersectionEvent.Chunk.RemoveVoxel(intersectionEvent.VoxelIndex)
		player.worldUpdater.QueueForRebuild(intersectionEvent.Chunk)
		// get adjacent chunks index coordinates with voxelindex
		adjacents := GetIntersectionAdjacentChunks(player.World.ChunkSettings, intersectionEvent)
		for _, yo := range adjacents {
			chunk, ok := player.World.Region.GetChunk(yo.ChunkCoordinate)
			if ok {
				index := player.GetSettings().CoordinateToIndex(yo.VoxelCoordinate)
				chunk.SetVoxelVisibility(index, true)
				player.worldUpdater.QueueForRebuild(chunk)
			} else {
				// check if chunk coordinate is underground and if yes, create a new one
				ci := yo.ChunkCoordinate.i * player.World.ChunkSettings.GetChunkWidth()
				cj := yo.ChunkCoordinate.j * player.World.ChunkSettings.GetChunkDepth()
				ck := yo.ChunkCoordinate.k * player.World.ChunkSettings.GetChunkHeight()
				surface_k := player.World.get_surface_k(ci, cj)
				if ck < surface_k {
					// spawn new full chunk
					newChunk := player.World.GenerateUndergroundChunk(yo.ChunkCoordinate)
					index := player.GetSettings().CoordinateToIndex(yo.VoxelCoordinate)
					newChunk.SetVoxelVisibility(index, true)
					player.worldUpdater.QueueForRebuild(newChunk)
				}
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
	for _, dChunk := range dChunks {
		chunk := dChunk.chunk
		event, ok := chunk.Intersect(p1, p2)
		if ok {
			return event, true
		}
	}
	return nil, false
}

func GetChunkDistance(player *Player, c IndexCoordinate) float64 {
	sizeX := float64(player.World.ChunkSettings.GetChunkWidth())
	sizeY := float64(player.World.ChunkSettings.GetChunkDepth())
	sizeZ := float64(player.World.ChunkSettings.GetChunkHeight())
	halfX := sizeX / 2.0
	halfY := sizeY / 2.0
	halfZ := sizeZ / 2.0
	chunkPos := mgl64.Vec3{
		float64(c.i)*sizeX + halfX,
		float64(c.j)*sizeY + halfY,
		float64(c.k)*sizeZ + halfZ,
	}
	d := player.Camera.Position.Sub(chunkPos)
	return d[0]*d[0] + d[1]*d[1] + d[2]*d[2]
}

// todo: is this a utility function that should be extracted from this file?
func GetIntersectionAdjacentChunks(settings IChunkSettings, event *IntersectionEvent) []ChunkAndVoxelCoordinate {
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
