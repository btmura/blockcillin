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
		uniform mat3 u_matrix;

		attribute vec2 a_position;

		void main(void) {
			vec3 position = vec3(a_position, 1) * u_matrix;
			gl_Position = vec4(position, 1);
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

	positionAttrib, err := getAttribLocation(program, "a_position")
	if err != nil {
		log.Fatalf("getAttribLocation: %v", err)
	}

	matrixUniform, err := getUniformLocation(program, "u_matrix")
	if err != nil {
		log.Fatalf("getUniformLocation: %v", err)
	}

	vertices := []float32{
		-1.0, -1.0,
		0.0, 1.0,
		1.0, -1.0,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	m := makeScaleMatrix(0.5, 0.5)
	m = multipleMatrices(m, makeRotationMatrix(toRadians(30.0)))
	m = multipleMatrices(m, makeTranslationMatrix(0.5, 0.5))
	gl.UniformMatrix3fv(matrixUniform, 1, true, &m[0])

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

func makeTranslationMatrix(x, y float32) []float32 {
	return []float32{
		1, 0, 0,
		0, 1, 0,
		x, y, 1,
	}
}

func makeRotationMatrix(radians float64) []float32 {
	c := float32(math.Cos(radians))
	s := float32(math.Sin(radians))
	return []float32{
		c, -s, 0,
		s, c, 0,
		0, 0, 1,
	}
}

func makeScaleMatrix(sx, sy float32) []float32 {
	return []float32{
		sx, 0, 0,
		0, sy, 0,
		0, 0, 1,
	}
}

func multipleMatrices(m, n []float32) []float32 {
	return []float32{
		m[0]*n[0] + m[1]*n[3] + m[2]*n[6], m[0]*n[1] + m[1]*n[4] + m[2]*n[7], m[0]*n[2] + m[1]*n[5] + m[2]*n[8],
		m[3]*n[0] + m[4]*n[3] + m[5]*n[6], m[3]*n[1] + m[4]*n[4] + m[5]*n[7], m[3]*n[2] + m[4]*n[5] + m[5]*n[8],
		m[6]*n[0] + m[7]*n[3] + m[8]*n[6], m[6]*n[1] + m[7]*n[4] + m[8]*n[7], m[6]*n[2] + m[7]*n[5] + m[8]*n[8],
	}
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
