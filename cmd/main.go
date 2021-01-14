package main

import (
	"log"
	"runtime"

	// "ant.com/ant/pkg/game/quad"
	// "ant.com/ant/pkg/game/text"
	// "ant.com/ant/pkg/game/cubes"
	"ant.com/ant/pkg/game/voxels"
)

const (
	windowWidth  = 1600
	windowHeight = 900
)

func main() {
	log.Println("Application starting...")
	runtime.LockOSThread()
	// game := quad.BuildGame(windowWidth, windowHeight, "resources/text-atlas-monospace-white-outlined-on-alpha.png")
	// game := quad.BuildGame(windowWidth, windowHeight, "resources/text-atlas-monospace-white-on-black.png")
	// game := text.BuildGame(windowWidth, windowHeight)
	// game := cubes.BuildGame(windowWidth, windowHeight)
	game := voxels.BuildGame(windowWidth, windowHeight)
	log.Println("Now running...")
	game.Run()
	log.Println("Application end")
}
