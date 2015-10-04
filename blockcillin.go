package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	// This is needed to arrange that main() runs on the main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var vertexShaderSource = `
	void main() {
		gl_Position = gl_Vertex;
	}
` + "\x00"

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

	if err := gl.Init(); err != nil {
		log.Fatalf("gl.Init: %v", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s", version)

	vs, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatalf("compileShader: %v", err)
	}

	gl.ClearColor(0, 0, 0, 0)

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		win.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteShader(vs)
}

func compileShader(shaderSource string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(gl.VERTEX_SHADER)
	str := gl.Str(shaderSource)
	gl.ShaderSource(shader, 1, &str, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader: type: %d source: %q log: %q", shaderType, shaderSource, log)
	}

	return shader, nil
}
