package ant

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type IUniformStore interface {
	registerUniform(name string)
	getLocation(name string) int32

	getMat3(name string) mgl32.Mat3
	getMat4(name string) mgl32.Mat4
	getVec4(name string) mgl32.Vec4
	getVec3(name string) mgl32.Vec3
	getVec2(name string) mgl32.Vec2

	setMat3(name string, value mgl32.Mat3)
	setMat4(name string, value mgl32.Mat4)
	setVec4(name string, value mgl32.Vec4)
	setVec3(name string, value mgl32.Vec3)
	setVec2(name string, value mgl32.Vec2)

	uniformMat3(name string, value mgl32.Mat3)
	uniformMat4(name string, value mgl32.Mat4)
	uniformVec4(name string, value mgl32.Vec4)
	uniformVec3(name string, value mgl32.Vec3)
	uniformVec2(name string, value mgl32.Vec2)
}

type UniformStore struct {
	glslProgramHandle uint32
	locations         map[string]int32
	mat3Map           map[string]mgl32.Mat3
	mat4Map           map[string]mgl32.Mat4
	vec2Map           map[string]mgl32.Vec2
	vec3Map           map[string]mgl32.Vec3
	vec4Map           map[string]mgl32.Vec4
}

func createUniformStore(glslProgramHandle uint32, shouldList bool) *UniformStore {
	store := new(UniformStore)
	store.glslProgramHandle = glslProgramHandle
	store.locations = make(map[string]int32)
	store.mat3Map = make(map[string]mgl32.Mat3)
	store.mat4Map = make(map[string]mgl32.Mat4)
	store.vec2Map = make(map[string]mgl32.Vec2)
	store.vec3Map = make(map[string]mgl32.Vec3)
	store.vec4Map = make(map[string]mgl32.Vec4)
	store.registerActiveUniforms(shouldList)
	return store
}

func (self *UniformStore) registerUniform(name string) {
	self.locations[name] = gl.GetUniformLocation(self.glslProgramHandle, gl.Str(name+"\x00"))
}

func (uniforms *UniformStore) getLocation(name string) int32 {
	location, ok := uniforms.locations[name]
	if !ok {
		panic("No uniform location is stored for name " + name)
	}
	return location
}

// value getters

func (uniforms *UniformStore) getMat4(name string) mgl32.Mat4 {
	value, ok := uniforms.mat4Map[name]
	if !ok {
		panic("No Mat4 value stored for name " + name)
	}
	return value
}

func (uniforms *UniformStore) getMat3(name string) mgl32.Mat3 {
	value, ok := uniforms.mat3Map[name]
	if !ok {
		panic("No Mat3 value stored for name " + name)
	}
	return value
}

func (uniforms *UniformStore) getVec2(name string) mgl32.Vec2 {
	value, ok := uniforms.vec2Map[name]
	if !ok {
		panic("No Vec2 value stored for name " + name)
	}
	return value
}

func (uniforms *UniformStore) getVec3(name string) mgl32.Vec3 {
	value, ok := uniforms.vec3Map[name]
	if !ok {
		panic("No Vec3 value stored for name " + name)
	}
	return value
}

func (uniforms *UniformStore) getVec4(name string) mgl32.Vec4 {
	value, ok := uniforms.vec4Map[name]
	if !ok {
		panic("No Vec4 value stored for name " + name)
	}
	return value
}

// value setters

func (self *UniformStore) setMat3(name string, value mgl32.Mat3) {
	self.mat3Map[name] = value
}

func (self *UniformStore) setMat4(name string, value mgl32.Mat4) {
	self.mat4Map[name] = value
}

func (self *UniformStore) setVec2(name string, value mgl32.Vec2) {
	self.vec2Map[name] = value
}

func (self *UniformStore) setVec3(name string, value mgl32.Vec3) {
	self.vec3Map[name] = value
}

func (self *UniformStore) setVec4(name string, value mgl32.Vec4) {
	self.vec4Map[name] = value
}

// uniform setters

func (self *UniformStore) uniformMat3(name string, value mgl32.Mat3) {
	location := self.getLocation(name)
	gl.UniformMatrix3fv(location, 1, false, &value[0])
}

func (self *UniformStore) uniformMat4(name string, value mgl32.Mat4) {
	location := self.getLocation(name)
	gl.UniformMatrix4fv(location, 1, false, &value[0])
}

func (self *UniformStore) uniformVec2(name string, value mgl32.Vec2) {
	location := self.getLocation(name)
	gl.Uniform2fv(location, 1, &value[0])
}

func (self *UniformStore) uniformVec3(name string, value mgl32.Vec3) {
	location := self.getLocation(name)
	gl.Uniform3fv(location, 1, &value[0])
}

func (self *UniformStore) uniformVec4(name string, value mgl32.Vec4) {
	location := self.getLocation(name)
	gl.Uniform4fv(location, 1, &value[0])
}

// registerActiveUniforms

func (self *UniformStore) registerActiveUniforms(shouldList bool) {
	var count int32 = 0
	gl.GetProgramiv(self.glslProgramHandle, gl.ACTIVE_UNIFORMS, &count)
	if shouldList {
		log.Println("Listing uniforms, count:", count)
	}
	var i uint32 = 0
	for i < uint32(count) {
		name := getUniformName(self.glslProgramHandle, i, shouldList)
		self.registerUniform(name)
		i++
	}
}

func getUniformName(program uint32, i uint32, shouldList bool) string {
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
	if shouldList {
		xtypeString := xtypeToString(xtype)
		logMessage := fmt.Sprintf("%s \t %s", xtypeString, name)
		log.Println(logMessage)
	}
	return name
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
