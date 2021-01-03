package ant

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type GLSLProgram struct {
	Handle    uint32
	Linked    bool
	logString string
}

func NewGLSLProgram() GLSLProgram {
	glProgram := gl.CreateProgram()
	log.Println("glProgram created with handle", glProgram)
	return GLSLProgram{glProgram, false, ""}
}

func (p GLSLProgram) Link() bool {
	gl.LinkProgram(p.GetHandle())
	p.Linked = true
	return true
}

func (p GLSLProgram) Use() bool {
	gl.UseProgram(p.Handle)
	return true
}

func (p GLSLProgram) Log() bool {
	return false
}

func (p GLSLProgram) GetHandle() uint32 {
	return p.Handle
}

func (p GLSLProgram) CompileAndAttachShader(source string, shaderType uint32) {
	shader, err := CompileShader(source, shaderType)
	check(err)
	gl.AttachShader(p.GetHandle(), shader)
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func InitGlslProgram(vertexShaderPath, fragmentShaderPath string) GLSLProgram {
	glslProgram := NewGLSLProgram()
	log.Println("reading vertex shader")
	vertex := ReadFile(vertexShaderPath)
	log.Println("compiling vertex shader")
	glslProgram.CompileAndAttachShader(vertex, gl.VERTEX_SHADER)
	log.Println("reading fragment shader")
	fragment := ReadFile(fragmentShaderPath)
	log.Println("compiling fragment shader")
	glslProgram.CompileAndAttachShader(fragment, gl.FRAGMENT_SHADER)
	log.Println("linking shader program")
	glslProgram.Link()
	return glslProgram
}
