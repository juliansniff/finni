package frontend

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const vs = `
#version 410 core
layout (location = 0) in vec4 vertex;

out vec2 TexCoords;

void main()
{
	gl_Position = vec4(vertex.xy, 1.0, 1.0);
	TexCoords = vec2(vertex.zw);
}
` + "\x00"

const fs = `
#version 410 core
in vec2 TexCoords;
out vec4 color;

uniform sampler2D text;

void main()
{
	color = vec4(0.0, 0.0, 0.0, texture(text, TexCoords).r);
}
` + "\x00"

// Shader encapsulates the ID of a shader program and
// creates several utility functions to ease life.
type Shader struct {
	// ID is the program ID of the shader.
	ID uint32
}

func NewShader() (*Shader, error) {
	shader := &Shader{}
	var success int32
	var infoLog [512]uint8

	vSource, vFree := gl.Strs(vs + "\x00")
	fSource, fFree := gl.Strs(fs + "\x00")

	// compile vertex shader
	vertex := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertex, 1, vSource, nil)
	gl.CompileShader(vertex)
	gl.GetShaderiv(vertex, gl.COMPILE_STATUS, &success)
	if success == 0 {
		gl.GetShaderInfoLog(vertex, 512, nil, &infoLog[0])
		return shader, fmt.Errorf("error compiling vertex shader: %s", infoLog)
	}

	// compile fragment shader
	fragment := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragment, 1, fSource, nil)
	gl.CompileShader(fragment)
	gl.GetShaderiv(fragment, gl.COMPILE_STATUS, &success)
	if success == 0 {
		gl.GetShaderInfoLog(fragment, 512, nil, &infoLog[0])
		return shader, fmt.Errorf("error compiling fragment shader: %s", infoLog)
	}

	// create program
	shader.ID = gl.CreateProgram()
	gl.AttachShader(shader.ID, vertex)
	gl.AttachShader(shader.ID, fragment)
	gl.LinkProgram(shader.ID)
	gl.GetProgramiv(shader.ID, gl.LINK_STATUS, &success)
	if success == 0 {
		gl.GetProgramInfoLog(shader.ID, 512, nil, &infoLog[0])
		return shader, fmt.Errorf("error linking program: %s", infoLog)
	}

	// delete shaders
	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	// free memory
	vFree()
	fFree()

	return shader, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) SetBool(name string, value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), i)
}

func (s *Shader) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

func (s *Shader) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
}
