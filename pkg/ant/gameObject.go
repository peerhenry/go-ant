package ant

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Drawable interface {
	draw()
}

type Updatable interface {
	update()
}

// type IGameObject interface {
// 	draw(renderState *GameRenderState)
// 	update()
// }

type GameObject struct {
	vao         uint32
	indexLength int32
	position    mgl32.Vec3
	orientation mgl32.Quat
}
