package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/juliansniff/finni/frontend"
)

var VAO, VBO uint32

const (
	width  = 800
	height = 600
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

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

	window, err := glfw.CreateWindow(width, height, "Finni", nil, nil)
	if err != nil {
		log.Panic("failed to create window: %v", err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Panic("failed to initialize OpenGL: %v", err)
	}

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.ClearColor(1.0, 1.0, 0.9176, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		w.SwapBuffers()
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	shader, err := frontend.NewShader()
	if err != nil {
		log.Panic(err)
	}
	shader.Use()

	// sq := frontend.NewSquare()
	f, err := frontend.NewFont()
	if err != nil {
		log.Panic(err)
	}
	buffer, err := frontend.NewBuffer("test.txt", f)
	if err != nil {
		log.Panic(err)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(1.0, 1.0, 0.9176, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		buffer.Draw()

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
