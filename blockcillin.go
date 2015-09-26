package main

import (
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	// This is needed to arrange that main() runs on the main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalf("glfw.Init: %v", err)
	}
	defer glfw.Terminate()

	win, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		log.Fatalf("glfw.CreateWindow: %v", err)
	}

	win.MakeContextCurrent()

	for !win.ShouldClose() {
		// Do OpenGL stuff.
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
