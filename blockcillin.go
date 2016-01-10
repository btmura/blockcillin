package main

//go:generate go generate github.com/btmura/blockcillin/internal/asset
//go:generate go generate github.com/btmura/blockcillin/internal/audio
//go:generate go generate github.com/btmura/blockcillin/internal/game
//go:generate go generate github.com/btmura/blockcillin/internal/renderer

import (
	"flag"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/btmura/blockcillin/internal/audio"
	"github.com/btmura/blockcillin/internal/game"
	"github.com/btmura/blockcillin/internal/renderer"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var (
	fullScreen = flag.Bool("fs", true, "use fullscreen")
	seed       = flag.Int64("s", 0, "seed for the random number generator")
)

func init() {
	// This is needed to arrange that main() runs on the main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	flag.Parse()

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	log.Printf("seed: %d", *seed)
	rand.Seed(*seed)

	log.Printf("GLFW version: %s", glfw.GetVersionString())
	logFatalIfErr("glfw.Init", glfw.Init())
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()
	var fsMonitor *glfw.Monitor
	if *fullScreen {
		fsMonitor = monitor
	}
	win, err := glfw.CreateWindow(mode.Width, mode.Height, "blockcillin", fsMonitor, nil)
	logFatalIfErr("glfw.CreateWindow", err)
	win.MakeContextCurrent()

	logFatalIfErr("audio.Init", audio.Init())
	defer audio.Terminate()

	logFatalIfErr("renderer.Init", renderer.Init())
	defer renderer.Terminate()

	// Call the size callback to set the initial viewport.
	w, h := win.GetSize()
	renderer.SizeCallback(w, h)
	win.SetSizeCallback(func(w *glfw.Window, width, height int) {
		renderer.SizeCallback(width, height)
	})

	g := game.New()
	win.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		g.KeyCallback(key, action)
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

		renderer.Render(g, fudge)

		if g.Done() {
			win.SetShouldClose(true)
		}

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func logFatalIfErr(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", tag, err)
	}
}
