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
}
