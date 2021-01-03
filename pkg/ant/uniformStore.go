package ant

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type IUniformStore interface {
	RegisterUniform(name string)
	GetLocation(name string) int32

	GetMat3(name string) mgl32.Mat3
	GetMat4(name string) mgl32.Mat4
	GetVec4(name string) mgl32.Vec4
	GetVec3(name string) mgl32.Vec3
	GetVec2(name string) mgl32.Vec2

	SetMat3(name string, value mgl32.Mat3)
	SetMat4(name string, value mgl32.Mat4)
	SetVec4(name string, value mgl32.Vec4)
	SetVec3(name string, value mgl32.Vec3)
	SetVec2(name string, value mgl32.Vec2)

	UniformMat3(name string, value mgl32.Mat3)
	UniformMat4(name string, value mgl32.Mat4)
	UniformVec4(name string, value mgl32.Vec4)
	UniformVec3(name string, value mgl32.Vec3)
	UniformVec2(name string, value mgl32.Vec2)
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

func CreateUniformStore(glslProgramHandle uint32, shouldList bool) *UniformStore {
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

func (self *UniformStore) RegisterUniform(name string) {
	self.locations[name] = gl.GetUniformLocation(self.glslProgramHandle, gl.Str(name+"\x00"))
}

func (uniforms *UniformStore) GetLocation(name string) int32 {
	location, ok := uniforms.locations[name]
	if !ok {
		panic("No uniform location is stored for name " + name)
	}
	return location
}

// value Getters

func (uniforms *UniformStore) GetMat4(name string) mgl32.Mat4 {
	value, ok := uniforms.mat4Map[name]
	if !ok {
		panic("No Mat4 value stored for name " + name)
	}
	return value
}

func (uniforms *UniformStore) GetMat3(name string) mgl32.Mat3 {
	value, ok := uniforms.mat3Map[name]
	if !ok {
		panic("No Mat3 value stored for name " + name)
	}
	return value
}

func (uniforms *UniformStore) GetVec2(name string) mgl32.Vec2 {
	value, ok := uniforms.vec2Map[name]
	if !ok {
		panic("No Vec2 value stored for name " + name)
	}
	return value
}

func (uniforms *UniformStore) GetVec3(name string) mgl32.Vec3 {
	value, ok := uniforms.vec3Map[name]
	if !ok {
		panic("No Vec3 value stored for name " + name)
	}
	return value
}

func (uniforms *UniformStore) GetVec4(name string) mgl32.Vec4 {
	value, ok := uniforms.vec4Map[name]
	if !ok {
		panic("No Vec4 value stored for name " + name)
	}
	return value
}

// value Setters

func (self *UniformStore) SetMat3(name string, value mgl32.Mat3) {
	self.mat3Map[name] = value
}

func (self *UniformStore) SetMat4(name string, value mgl32.Mat4) {
	self.mat4Map[name] = value
}

func (self *UniformStore) SetVec2(name string, value mgl32.Vec2) {
	self.vec2Map[name] = value
}

func (self *UniformStore) SetVec3(name string, value mgl32.Vec3) {
	self.vec3Map[name] = value
}

func (self *UniformStore) SetVec4(name string, value mgl32.Vec4) {
	self.vec4Map[name] = value
}

// uniform Setters

func (self *UniformStore) UniformMat3(name string, value mgl32.Mat3) {
	location := self.GetLocation(name)
	gl.UniformMatrix3fv(location, 1, false, &value[0])
}

func (self *UniformStore) UniformMat4(name string, value mgl32.Mat4) {
	location := self.GetLocation(name)
	gl.UniformMatrix4fv(location, 1, false, &value[0])
}

func (self *UniformStore) UniformVec2(name string, value mgl32.Vec2) {
	location := self.GetLocation(name)
	gl.Uniform2fv(location, 1, &value[0])
}

func (self *UniformStore) UniformVec3(name string, value mgl32.Vec3) {
	location := self.GetLocation(name)
	gl.Uniform3fv(location, 1, &value[0])
}

func (self *UniformStore) UniformVec4(name string, value mgl32.Vec4) {
	location := self.GetLocation(name)
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
		self.RegisterUniform(name)
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
