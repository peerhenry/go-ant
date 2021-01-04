package ant

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL Version", version)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(100./256., 149./256., 237./256., 1.0)
}
