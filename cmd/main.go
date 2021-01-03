package main

import (
	"log"
	"runtime"

	"ant.com/ant/pkg/game"
)

const (
	windowWidth  = 1600
	windowHeight = 900
)

func main() {
	log.Println("Application starting...")
	runtime.LockOSThread()
	game := game.BuildCubeGame(windowWidth, windowHeight)
	log.Println("Now running...")
	game.Run()
	log.Println("Application end")
}
