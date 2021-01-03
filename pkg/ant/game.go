package ant

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Game struct {
	window *glfw.Window
	world  *GameWorld
}

func BuildGame(windowWidth, windowHeight int) Game {
	window := initGlfw(windowWidth, windowHeight)
	initOpenGL()

	world := buildCubeWorld(windowWidth, windowHeight)
	// world := buildQuadWorld()

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		pos := mgl32.Vec3{5, 3, 3}
		unit := mgl32.Vec3{-5, -3, -3}.Normalize()
		rot := mgl32.Rotate3DZ(float32(-xpos) / 200)
		target := pos.Add(rot.Mul3x1(unit))
		view := mgl32.LookAtV(
			mgl32.Vec3{5, 3, 3}, // eye
			target,              // center
			mgl32.Vec3{0, 0, 1}, // up
		)
		world.uniforms.setMat4("ViewMatrix", view)
	})

	return Game{
		window: window,
		world:  &world,
	}
}

func (game *Game) Run() {
	defer glfw.Terminate()
	for !game.window.ShouldClose() {
		// todo: update on separate thread
		game.update()
		game.draw()
	}
}

func (game *Game) update() {
	game.world.update()
}

func (game *Game) draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	game.world.render()
	glfw.PollEvents()
	game.window.SwapBuffers()
}

// =================================================

func initGlslProgram(vertexShaderPath, fragmentShaderPath string) GLSLProgram {
	glslProgram := NewGLSLProgram()
	log.Println("reading vertex shader")
	vertex := readFile(vertexShaderPath)
	log.Println("compiling vertex shader")
	glslProgram.CompileAndAttachShader(vertex, gl.VERTEX_SHADER)
	log.Println("reading fragment shader")
	fragment := readFile(fragmentShaderPath)
	log.Println("compiling fragment shader")
	glslProgram.CompileAndAttachShader(fragment, gl.FRAGMENT_SHADER)
	log.Println("linking shader program")
	glslProgram.Link()
	return glslProgram
}
