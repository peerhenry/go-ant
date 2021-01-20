package chunks

type IHeightProvider interface {
	GetHeight(ai, aj int) int
}
