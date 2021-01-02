package ant

import (
	"image"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Game struct {
	window *glfw.Window
	world  *GameWorld
}

func BuildGame(windowWidth, windowHeight int) Game {
	window := initGlfw(windowWidth, windowHeight)
	initOpenGL()

	// glslProgram := initGlslProgram("shaders/ads/vertex.glsl", "shaders/ads/fragment.glsl")
	// cube1 := createCube(mgl32.Vec3{1, 0, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	// cube2 := createCube(mgl32.Vec3{0, 1, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	// cube3 := createCube(mgl32.Vec3{0, 0, 1}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	// objects := []*GameObject{&cube1, &cube2, &cube3}

	glslProgram := initGlslProgram("shaders/quad/vertex.glsl", "shaders/quad/fragment.glsl")
	quad := createQuad()
	objects := []*GameObject{&quad}
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

	log.Println("Reading texture atlas")
	i := readImage("resources/atlas.png")
	switch i.(type) {
	case *image.RGBA:
		log.Println("i is an *image.RGBA")
	case *image.NRGBA:
		log.Println("i is an *image.NRBGA")
		if nrgba, ok := i.(*image.NRGBA); ok {
			// img is now an *image.RGBA
			log.Println("image", nrgba.Bounds().Dy())
			loadTexture(nrgba)
			gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
			gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
			// set uniform
			texUniformLocation := gl.GetUniformLocation(glslProgram.handle, gl.Str("Tex\x00"))
			gl.Uniform1i(texUniformLocation, 0)
		} else {
			log.Println("Warning: could not extract NRGBA from image...")
		}
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
	log.Println("reading vertex shader")
	vertex := readFile(vertexShaderPath)
	log.Println("compiling vertex shader")
	glslProgram.CompileAndAttachShader(vertex, gl.VERTEX_SHADER)
	log.Println("reading fragment shader")
	fragment := readFile(fragmentShaderPath)
	log.Println("compiling fragment shader")
	glslProgram.CompileAndAttachShader(fragment, gl.FRAGMENT_SHADER)
	log.Println("linking shader program...")
	glslProgram.Link()
	return glslProgram
}
