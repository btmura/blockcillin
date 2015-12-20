package main

//go:generate go-bindata data

import (
	"log"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/gordonklaus/portaudio"
)

const secPerUpdate = 1.0 / 60.0

const sampleRate = 44100

func init() {
	// This is needed to arrange that main() runs on the main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	logFatalIfErr("glfw.Init", glfw.Init())
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(640, 480, "blockcillin", nil, nil)
	logFatalIfErr("glfw.CreateWindow", err)
	win.MakeContextCurrent()

	logFatalIfErr("gl.Init", gl.Init())
	log.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	logFatalIfErr("portaudio.Initialize", portaudio.Initialize())
	defer func() {
		logFatalIfErr("portaudio.Terminate", portaudio.Terminate())
	}()
	log.Printf("PortAudio version: %d %s", portaudio.Version(), portaudio.VersionText())

	// PortAudio sample from https://github.com/gordonklaus/portaudio/blob/master/examples/stereoSine.go
	s := newStereoSine(256, 320, sampleRate)
	defer s.Close()
	logFatalIfErr("s.Start", s.Start())
	time.Sleep(2 * time.Second)
	logFatalIfErr("s.Stop", s.Stop())

	rr := newRenderer()

	// Call the size callback to set the initial viewport.
	w, h := win.GetSize()
	rr.sizeCallback(w, h)
	win.SetSizeCallback(func(w *glfw.Window, width, height int) {
		rr.sizeCallback(width, height)
	})

	g := newGame()
	win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		g.keyCallback(w, key, action)
	})

	var lag float64
	prevTime := glfw.GetTime()
	for !win.ShouldClose() {
		currTime := glfw.GetTime()
		elapsed := currTime - prevTime
		prevTime = currTime
		lag += elapsed

		for lag >= secPerUpdate {
			g.update()
			lag -= secPerUpdate
		}
		fudge := float32(lag / secPerUpdate)

		rr.render(g, fudge)

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func logFatalIfErr(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", tag, err)
	}
}

// PortAudio sample from https://github.com/gordonklaus/portaudio/blob/master/examples/stereoSine.go

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func newStereoSine(freqL, freqR, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, s.processAudio)
	logFatalIfErr("portaudio.OpenDefaultStream", err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
	}
}
