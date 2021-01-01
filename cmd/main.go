package main

import (
	"log"
	"runtime"

	"ant.com/ant/pkg/ant"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func main() {
	log.Println("Application starting...")
	runtime.LockOSThread()
	game := ant.BuildGame(windowWidth, windowHeight)
	log.Println("Now running...")
	game.Run()
	log.Println("Application end")
}
