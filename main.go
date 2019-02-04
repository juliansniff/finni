package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/juliansniff/finni/editing"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	fb := &editing.FileBuffer{
		File:   []byte("hello"),
		Cursor: 5,
	}

	fmt.Println(string(fb.File))
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
