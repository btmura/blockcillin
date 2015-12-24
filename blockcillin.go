package main

//go:generate go-bindata -debug -pkg asset -o internal/asset/bindata.go data

import (
	"log"
	"runtime"

	"github.com/btmura/blockcillin/internal/audio"
	"github.com/btmura/blockcillin/internal/game"
	"github.com/btmura/blockcillin/internal/renderer"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	// This is needed to arrange that main() runs on the main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	log.Printf("GLFW version: %s", glfw.GetVersionString())
	logFatalIfErr("glfw.Init", glfw.Init())
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(640, 480, "blockcillin", nil, nil)
	logFatalIfErr("glfw.CreateWindow", err)
	win.MakeContextCurrent()

	audio.Init()
	defer audio.Terminate()

	rr := renderer.NewRenderer()

	// Call the size callback to set the initial viewport.
	w, h := win.GetSize()
	rr.SizeCallback(w, h)
	win.SetSizeCallback(func(w *glfw.Window, width, height int) {
		rr.SizeCallback(width, height)
	})

	g := game.New()
	win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		g.KeyCallback(w, key, action)
	})

	var lag float64
	prevTime := glfw.GetTime()
	for !win.ShouldClose() {
		currTime := glfw.GetTime()
		elapsed := currTime - prevTime
		prevTime = currTime
		lag += elapsed

		for lag >= game.SecPerUpdate {
			g.Update()
			lag -= game.SecPerUpdate
		}
		fudge := float32(lag / game.SecPerUpdate)

		rr.Render(g, fudge)

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func logFatalIfErr(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", tag, err)
	}
}
