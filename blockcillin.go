package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	vertexShaderSource = `
		uniform vec2 u_translation;
		uniform vec2 u_rotation;

		attribute vec2 a_position;

		void main(void) {
			vec2 rotatedPosition = vec2(
				a_position.x * u_rotation.y + a_position.y * u_rotation.x,
				a_position.y * u_rotation.y - a_position.x * u_rotation.y);

			vec2 translatedPosition = rotatedPosition + u_translation;

			gl_Position = vec4(translatedPosition, 0, 1);
		}
	` + "\x00"

	fragmentShaderSource = `
		void main(void) {
			gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
		}
	` + "\x00"

	vertices = []float32{
		-1.0, -1.0,
		0.0, 2.0,
		1.0, -1.0,
	}

	translation = []float32{
		0.25, 0.25,
	}
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

	translationUniform, err := getUniformLocation(program, "u_translation")
	if err != nil {
		log.Fatalf("getUniformLocation: %v", err)
	}
	gl.Uniform2fv(translationUniform, 1, &translation[0])

	rotationUniform, err := getUniformLocation(program, "u_rotation")
	if err != nil {
		log.Fatalf("getUniformLocation: %v", err)
	}

	degrees := 45.0
	radians := degrees * math.Pi / 180
	rotation := []float32{
		float32(math.Sin(radians)),
		float32(math.Cos(radians)),
	}
	gl.Uniform2fv(rotationUniform, 1, &rotation[0])

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	positionAttrib, err := getAttribLocation(program, "a_position")
	if err != nil {
		log.Fatalf("getAttribLocation: %v", err)
	}
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.ClearColor(0, 0, 0, 0)

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

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

func getUniformLocation(program uint32, name string) (int32, error) {
	u := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	if u == -1 {
		return 0, fmt.Errorf("couldn't get uniform location: %q", name)
	}
	return u, nil
}

func getAttribLocation(program uint32, name string) (uint32, error) {
	a := gl.GetAttribLocation(program, gl.Str(name+"\x00"))
	if a == -1 {
		return 0, fmt.Errorf("couldn't get attrib location: %q", name)
	}
	// Cast to uint32 for EnableVertexAttribArray and VertexAttribPointer better.
	return uint32(a), nil
}
