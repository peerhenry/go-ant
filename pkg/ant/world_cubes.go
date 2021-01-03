package ant

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func buildCubeWorld(windowWidth, windowHeight int) GameWorld {
	glslProgram := initGlslProgram("shaders/ads/vertex.glsl", "shaders/ads/fragment.glsl")
	uniformStore := setupUniforms(glslProgram.handle, windowWidth, windowHeight)
	objects := createGameObjects()

	return GameWorld{
		uniforms:    uniformStore,
		glslProgram: &glslProgram,
		objects:     objects,
	}
}

func setupUniforms(glslProgramHandle uint32, windowWidth, windowHeight int) *UniformStore {
	uniformStore := createUniformStore(glslProgramHandle)
	listActiveUniforms(glslProgramHandle)
	// register uniforms
	uniformStore.registerUniform("ModelViewMatrix")
	uniformStore.registerUniform("NormalMatrix")
	uniformStore.registerUniform("ProjectionMatrix")
	uniformStore.registerUniform("MVP")
	// set values
	uniformStore.setMat4("ViewMatrix", mgl32.LookAtV(
		mgl32.Vec3{5, 3, 3}, // eye
		mgl32.Vec3{0, 0, 0}, // center
		mgl32.Vec3{0, 0, 1}, // up
	))
	uniformStore.setMat4("ProjectionMatrix", mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 10.0))
	loadImageFileToUniform("resources/atlas.png", "Tex", glslProgramHandle)
	return uniformStore
}

func createGameObjects() []*GameObject {
	cube1 := createCube(mgl32.Vec3{1, 0, 0}, mgl32.QuatRotate(0, mgl32.Vec3{0, 0, 1}))
	cube2 := createCube(mgl32.Vec3{0, 1, 0}, mgl32.QuatRotate(1, mgl32.Vec3{1, 0, 0}))
	cube3 := createCube(mgl32.Vec3{0, 0, 1}, mgl32.QuatRotate(1, mgl32.Vec3{0, 1, 0}))
	return []*GameObject{cube1, cube2, cube3}
}

func listActiveUniforms(program uint32) {
	var count int32 = 0
	gl.GetProgramiv(program, gl.ACTIVE_UNIFORMS, &count)
	log.Println("Number of uniforms:", count)
	log.Println("=====")
	var i uint32 = 0
	for i < uint32(count) {
		getAndPrintUniform(program, i)
		i++
	}
}

func getAndPrintUniform(program uint32, i uint32) {
	var bufSize int32 = 64
	var length int32 = 0
	var size int32 = 0
	var xtype uint32 = 0
	var nameBuffer [64]byte
	gl.GetActiveUniform(
		program,
		i,
		bufSize,        // bufSize
		&length,        // length; the amount of characters written to buffer
		&size,          // size
		&xtype,         // xtype
		&nameBuffer[0], // name
	)
	name := numerBufferToString(nameBuffer[:])
	xtypeString := xtypeToString(xtype)
	logMessage := fmt.Sprintf("%s \t %s", xtypeString, name)
	log.Println(logMessage)
}

func numerBufferToString(nameBuffer []byte) string {
	var nameSlice []byte
	for _, element := range nameBuffer {
		if element == 0 {
			break
		} else {
			nameSlice = append(nameSlice, element)
		}
	}
	return string(nameSlice)
}

func xtypeToString(xtype uint32) string {
	switch xtype {
	case gl.FLOAT_VEC2:
		return "FLOAT_VEC2"
	case gl.FLOAT_VEC3:
		return "FLOAT_VEC3"
	case gl.FLOAT_VEC4:
		return "FLOAT_VEC4"
	case gl.FLOAT_MAT2:
		return "FLOAT_MAT2"
	case gl.FLOAT_MAT3:
		return "FLOAT_MAT3"
	case gl.FLOAT_MAT4:
		return "FLOAT_MAT4"
	case gl.SAMPLER_2D:
		return "SAMPLER_2D"
	default:
		return "unknown"
	}
}
