package main

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

type sound int

const (
	soundMove sound = iota
	soundSelect
	soundSwap
	soundClear
)

// playSound is a global function that must be overridden to use an audioManager.
var playSound = func(s sound) {}

type audioManager struct {
	soundQueue chan sound
	done       chan bool
}

func newAudioManager() *audioManager {
	return &audioManager{
		soundQueue: make(chan sound, 100),
		done:       make(chan bool),
	}
}

func (a *audioManager) start() {
	go func() {
		makeBuffer := func(name string) []int16 {
			wav, err := decodeWAV(newAssetReader(name))
			logFatalIfErr("decodeWAV", err)
			log.Printf("%s: %+v", name, wav)

			buf := make([]int16, len(wav.data)/2)
			for i := 0; i < len(buf); i++ {
				buf[i] = int16(wav.data[i*2])
				buf[i] += int16(wav.data[i*2+1]) << 8
			}
			return buf
		}

		buffers := map[sound][]int16{
			soundMove:   makeBuffer("data/move.wav"),
			soundSelect: makeBuffer("data/select.wav"),
			soundSwap:   makeBuffer("data/swap.wav"),
			soundClear:  makeBuffer("data/clear.wav"),
		}

		for s := range a.soundQueue {
			out := buffers[s]
			stream, err := portaudio.OpenDefaultStream(0 /*input channels */, 2, 44100, len(out), out)
			logFatalIfErr("portaudio.OpenDefaultStream", err)
			logFatalIfErr("stream.Start", stream.Start())
			stream.Write()
			logFatalIfErr("stream.Stop", stream.Stop())
			logFatalIfErr("stream.Close", stream.Close())
		}

		a.done <- true
	}()
}

func (a *audioManager) play(s sound) {
	a.soundQueue <- s
}

func (a *audioManager) stop() {
	close(a.soundQueue)
	<-a.done
	close(a.done)
}
