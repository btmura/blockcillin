package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	vertexShaderSource = `
		void main(void) {
			gl_Position = gl_Vertex;
		}
	` + "\x00"

	fragmentShaderSource = `
		void main(void) {
			gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
		}
	` + "\x00"
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

	if err := gl.Init(); err != nil {
		log.Fatalf("gl.Init: %v", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s", version)

	program, err := createProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatalf("createProgram: %v", err)
	}
	gl.UseProgram(program)

	gl.ClearColor(0, 0, 0, 0)

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func createProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vs, err := createShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fs, err := createShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to create program: %q", log)
	}

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	return program, nil
}

func createShader(shaderSource string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
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
