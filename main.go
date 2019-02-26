package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/juliansniff/finni/frontend"
)

var font *frontend.Font
var VAO, VBO uint32

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	font, _ = frontend.NewFont(120, 120)
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Panic("failed to initialize glfw: %v", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(1200, 800, "LearnOpengGL", nil, nil)
	if err != nil {
		log.Panic("failed to create window: %v", err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Panic("failed to initialize OpenGL: %v", err)
	}

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	shader, err := frontend.NewShader("shaders/shader.vs", "shaders/shader.fs")
	if err != nil {
		log.Panic(err)
	}
	shader.Use()

	projection := mgl32.Ortho2D(0, 800, 0, 600)
	gl.UniformMatrix4fv(shader.GetUniformLocation("projection"), 1, false, &projection[0])

	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.DYNAMIC_DRAW)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(1.0, 1.0, 0.9176, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		renderText(shader, 10, "func() {}")

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyL) == glfw.Press {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}
	if window.GetKey(glfw.KeyF) == glfw.Press {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

func renderText(s *frontend.Shader, scale float32, str string) {
	s.Use()
	gl.Uniform3f(s.GetUniformLocation("textColor"), 0.5, 0.5, 0.5)
	gl.BindVertexArray(VAO)

	x := 0

	for _, char := range str {
		r := rune(char)
		character, err := font.GetCharacter(r)
		if err != nil {
			log.Panic(err)
		}

		xpos := float32(x)
		var w, h float32
		w = float32(character.Image.Bounds().Size().X + 5)
		h = float32(character.Image.Bounds().Size().Y + 5)

		vertices := []float32{
			xpos, h, 0, 0,
			xpos, 0, 0, 1,
			xpos + w, 0, 1, 1,
			xpos, h, 0, 0,
			xpos + w, 0, 1, 1,
			xpos + w, h, 1, 0,
		}

		gl.BindTexture(gl.TEXTURE_2D, character.Texture)
		gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, 6*4*4, gl.Ptr(&vertices[0]))
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		x += 50
	}
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
