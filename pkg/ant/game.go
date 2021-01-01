package ant

import (
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
	// glslProgram := initGlslProgram("shaders/quad/fragment.glsl", "shaders/quad/fragment.glsl")
	glslProgram := initGlslProgram("shaders/ads/vertex.glsl", "shaders/ads/fragment.glsl")
	cube1 := createCube(mgl32.Vec3{1, 0, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	cube2 := createCube(mgl32.Vec3{0, 1, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	cube3 := createCube(mgl32.Vec3{0, 0, 1}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	objects := []*GameObject{&cube1, &cube2, &cube3}
	// quad := createQuad()
	// drawables := []*Drawable{&quad}
	// updatables := []*Updatable{}
	renderState := buildGameRenderState(glslProgram.handle, windowWidth, windowHeight)
	world := GameWorld{
		renderState: &renderState,
		glslProgram: &glslProgram,
		objects:     objects,
		// drawables:   drawables,
		// updatables:  updatables,
	}
	return Game{
		window: window,
		world:  &world,
	}
}

func (game *Game) Run() {
	defer glfw.Terminate()
	for !game.window.ShouldClose() {
		game.update()
		game.draw()
	}
}

func (game *Game) update() {
	game.world.update()
}

func (game *Game) draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// todo: update on separate thread
	game.world.render()
	glfw.PollEvents()
	game.window.SwapBuffers()
}

// =================================================

func initGlslProgram(vertexShaderPath, fragmentShaderPath string) GLSLProgram {
	glslProgram := NewGLSLProgram()
	vertex := readFile(vertexShaderPath)
	fragment := readFile(fragmentShaderPath)
	glslProgram.CompileAndAttachShader(vertex, gl.VERTEX_SHADER)
	glslProgram.CompileAndAttachShader(fragment, gl.FRAGMENT_SHADER)
	glslProgram.Link()
	return glslProgram
}
