package main

//go:generate go-bindata data

import (
	"encoding/binary"
	"log"
	"runtime"

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

	// Based upon: https://github.com/verdverm/go-wav/blob/master/wav.go
	type WAV struct {
		bChunkID  [4]byte
		chunkSize uint32
		bFormat   [4]byte

		bSubchunk1ID  [4]byte
		subchunk1Size uint32
		audioFormat   uint16
		numChannels   uint16
		sampleRate    uint32
		byteRate      uint32
		blockAlign    uint16
		bitsPerSample uint16

		bSubchunk2ID  [4]byte
		subchunk2Size uint32
		data          []byte
	}

	wav := WAV{}

	r := newAssetReader("data/sounds.wav")
	binary.Read(r, binary.BigEndian, &wav.bChunkID)
	binary.Read(r, binary.LittleEndian, &wav.chunkSize)
	binary.Read(r, binary.BigEndian, &wav.bFormat)

	binary.Read(r, binary.BigEndian, &wav.bSubchunk1ID)
	binary.Read(r, binary.LittleEndian, &wav.subchunk1Size)
	binary.Read(r, binary.LittleEndian, &wav.audioFormat)
	binary.Read(r, binary.LittleEndian, &wav.numChannels)
	binary.Read(r, binary.LittleEndian, &wav.sampleRate)
	binary.Read(r, binary.LittleEndian, &wav.byteRate)
	binary.Read(r, binary.LittleEndian, &wav.blockAlign)
	binary.Read(r, binary.LittleEndian, &wav.bitsPerSample)

	binary.Read(r, binary.BigEndian, &wav.bSubchunk2ID)
	binary.Read(r, binary.LittleEndian, &wav.subchunk2Size)

	log.Printf("WAV: %+v", wav)

	wav.data = make([]byte, wav.subchunk2Size)
	binary.Read(r, binary.LittleEndian, &wav.data)

	processAudio := func(out []int16) {
		// bChunkID:[82 73 70 70] chunkSize:705636 bFormat:[87 65 86 69]
		// bSubchunk1ID:[102 109 116 32] subchunk1Size:16 audioFormat:1 numChannels:2 sampleRate:44100 byteRate:176400 blockAlign:4 bitsPerSample:16
		// bSubchunk2ID:[100 97 116 97] subchunk2Size:705600 data:[]}
		var trim int
		for i := 0; i < len(out); i++ {
			var val int16
			if i*2 < len(wav.data) {
				val = int16(wav.data[i*2])
				val += int16(wav.data[i*2+1]) << 8
				trim += 2
			}
			out[i] = val
		}
		wav.data = wav.data[trim:]
	}

	stream, err := portaudio.OpenDefaultStream(0, 2, sampleRate, 0, processAudio)
	logFatalIfErr("portaudio.OpenDefaultStream", err)
	defer stream.Close()

	logFatalIfErr("s.Start", stream.Start())
	defer func() {
		logFatalIfErr("s.Stop", stream.Stop())
	}()

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
