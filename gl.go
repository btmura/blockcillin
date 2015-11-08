package main

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"

	_ "image/png" // needed to decode PNGs
)

func CreateProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
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
	str := gl.Str(shaderSource + "\x00")
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

func CreateTexture(r io.Reader) (uint32, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	return texture, nil
}

func GetUniformLocation(program uint32, name string) (int32, error) {
	u := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	if u == -1 {
		return 0, fmt.Errorf("couldn't get uniform location: %q", name)
	}
	return u, nil
}

func GetAttribLocation(program uint32, name string) (uint32, error) {
	a := gl.GetAttribLocation(program, gl.Str(name+"\x00"))
	if a == -1 {
		return 0, fmt.Errorf("couldn't get attrib location: %q", name)
	}
	// Cast to uint32 for EnableVertexAttribArray and VertexAttribPointer better.
	return uint32(a), nil
}
