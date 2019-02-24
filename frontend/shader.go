package frontend

import (
	"fmt"
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Shader encapsulates the ID of a shader program and
// creates several utility functions to ease life.
type Shader struct {
	// ID is the program ID of the shader.
	ID uint32
}

func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	shader := &Shader{}
	var success int32
	var infoLog [512]uint8

	// read vertex shader
	b, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		return shader, fmt.Errorf("error reading vertex shader: %v", err)
	}
	vSource, vFree := gl.Strs(string(b) + "\x00")

	// read fragment shader
	b, err = ioutil.ReadFile(fragmentPath)
	if err != nil {
		return shader, fmt.Errorf("error reading fragment shader: %v", err)
	}
	fSource, fFree := gl.Strs(string(b) + "\x00")

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
