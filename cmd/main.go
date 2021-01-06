package main

import (
	"log"
	"runtime"

	// "ant.com/ant/pkg/game/cubes"
	// "ant.com/ant/pkg/game/quad"
	"ant.com/ant/pkg/game/text"
)

const (
	windowWidth  = 1600
	windowHeight = 900
)

func main() {
	log.Println("Application starting...")
	runtime.LockOSThread()
	// game := cubes.BuildCubeGame(windowWidth, windowHeight)
	// game := quad.BuildGame(windowWidth, windowHeight)
	game := text.BuildGame(windowWidth, windowHeight)
	log.Println("Now running...")
	game.Run()
	log.Println("Application end")
}
