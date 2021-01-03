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

	glslProgram := initGlslProgram("shaders/ads/vertex.glsl", "shaders/ads/fragment.glsl")
	cube1 := createCube(mgl32.Vec3{1, 0, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	cube2 := createCube(mgl32.Vec3{0, 1, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	cube3 := createCube(mgl32.Vec3{0, 0, 1}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	objects := []*GameObject{cube1, cube2, cube3}

	// glslProgram := initGlslProgram("shaders/quad/vertex.glsl", "shaders/quad/fragment.glsl")
	// quad := createQuad()
	// objects := []*GameObject{&quad}
	// // drawables := []*Drawable{&quad}
	// // updatables := []*Updatable{}
	// renderState := buildGameRenderState(glslProgram.handle, windowWidth, windowHeight)

	log.Println("Time to create uniform store with handle", glslProgram.handle)
	uniforms := createUniformStore(glslProgram.handle)
	// register uniforms
	uniforms.registerUniform("ModelViewMatrix")
	// uniforms.registerUniform("NormalMatrix")
	uniforms.registerUniform("ProjectionMatrix")
	uniforms.registerUniform("MVP")
	// set values
	uniforms.setMat4("ViewMatrix", mgl32.LookAtV(
		mgl32.Vec3{5, 3, 3}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	))
	uniforms.setMat4("ProjectionMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0))

	world := GameWorld{
		uniforms:    uniforms,
		glslProgram: &glslProgram,
		objects:     objects,
	}

	// loadImageFileToUniform("resources/atlas.png", "Tex", glslProgram.handle)

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
