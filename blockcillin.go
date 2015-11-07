package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	vertexShaderSource = `
		uniform mat4 u_projection_view_matrix;
		uniform mat4 u_matrix;

		attribute vec4 a_position;
		attribute vec2 a_tex_coord;

		varying vec2 v_tex_coord;

		void main(void) {
			gl_Position = u_projection_view_matrix * u_matrix * a_position;
			v_tex_coord = a_tex_coord;
		}
	` + "\x00"

	fragmentShaderSource = `
		uniform sampler2D u_texture;

		varying vec2 v_tex_coord;

		void main(void) {
			gl_FragColor = texture2D(u_texture, v_tex_coord);
		}
	` + "\x00"
)

func init() {
	// This is needed to arrange that main() runs on the main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	mf, err := os.Open("models.obj")
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer mf.Close()

	objs, err := ReadObjFile(mf)
	if err != nil {
		log.Fatalf("ReadObjFile: %v", err)
	}

	tf, err := os.Open("texture.png")
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer tf.Close()

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

	var vertices []float32
	var texCoords []float32
	for _, o := range objs {
		for _, v := range o.Vertices {
			vertices = append(vertices, v.X, v.Y, v.Z)
		}
		for _, tc := range o.TexCoords {
			texCoords = append(texCoords, tc.S, tc.T)
		}
	}

	var indices []uint16
	for _, f := range objs[0].Faces {
		for _, e := range *f {
			indices = append(indices, uint16(e.VertexIndex-1) /* make 0-based */)
		}
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4 /* total bytes */, gl.Ptr(vertices), gl.STATIC_DRAW)

	positionAttrib := getAttribLocation(program, "a_position")
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	var tbo uint32
	gl.GenBuffers(1, &tbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, tbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(texCoords)*4 /*total bytes */, gl.Ptr(texCoords), gl.STATIC_DRAW)

	texCoordAttrib := getAttribLocation(program, "a_tex_coord")
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*2 /* total bytes */, gl.Ptr(indices), gl.STATIC_DRAW)

	t, err := createTexture(tf)
	if err != nil {
		log.Fatalf("createTexture: %v", err)
	}

	m := NewScaleMatrix(0.5, 0.5, 0.5)
	m = m.Mult(NewZRotationMatrix(toRadians(30.0)))
	m = m.Mult(NewTranslationMatrix(0.5, 0.5, 0.0))

	matrixUniform := getUniformLocation(program, "u_matrix")
	gl.UniformMatrix4fv(matrixUniform, 1, false, &m[0])

	textureUniform := getUniformLocation(program, "u_texture")
	gl.Uniform1i(textureUniform, 0)

	projectionViewMatrixUniform := getUniformLocation(program, "u_projection_view_matrix")
	vm := makeViewMatrix()
	sizeCallback := func(w *glfw.Window, width, height int) {
		pvm := vm.Mult(makeProjectionMatrix(width, height))
		gl.UniformMatrix4fv(projectionViewMatrixUniform, 1, false, &pvm[0])
		gl.Viewport(0, 0, int32(width), int32(height))
	}
	win.SetSizeCallback(sizeCallback)

	// Call the size callback to set the initial projection view matrix and viewport.
	w, h := win.GetSize()
	sizeCallback(win, w, h)

	gl.ClearColor(0, 0, 0, 0)
	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, t)

		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_SHORT, gl.Ptr(nil))

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

func getUniformLocation(program uint32, name string) int32 {
	u := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	if u == -1 {
		log.Fatalf("couldn't get uniform location: %q", name)
	}
	return u
}

func getAttribLocation(program uint32, name string) uint32 {
	a := gl.GetAttribLocation(program, gl.Str(name+"\x00"))
	if a == -1 {
		log.Fatalf("couldn't get attrib location: %q", name)
	}
	// Cast to uint32 for EnableVertexAttribArray and VertexAttribPointer better.
	return uint32(a)
}

func makeProjectionMatrix(width, height int) *Matrix4 {
	aspect := float32(width) / float32(height)
	fovRadians := float32(math.Pi) / 2
	return NewPerspectiveMatrix(fovRadians, aspect, 1, 2000)
}

func makeViewMatrix() *Matrix4 {
	cameraPosition := &Vector3{0, 0, 3}
	targetPosition := &Vector3{}
	up := &Vector3{0, 1, 0}
	return NewViewMatrix(cameraPosition, targetPosition, up)
}

func toRadians(degrees float32) float32 {
	return degrees * float32(math.Pi) / 180
}
