package main

import (
	"log"
	"runtime"

	"ant.com/ant/pkg/game/cubes"
)

const (
	windowWidth  = 1600
	windowHeight = 900
)

func main() {
	log.Println("Application starting...")
	runtime.LockOSThread()
	game := cubes.BuildCubeGame(windowWidth, windowHeight)
	log.Println("Now running...")
	game.Run()
	log.Println("Application end")
}
