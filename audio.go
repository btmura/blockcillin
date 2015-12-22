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
// It's used instead of passing an audioManager everywhere and making tests easier.
var playSound = func(s sound) {}

type audioManager struct {
	// soundBatch is the current batch of sounds that will be queued to play at once.
	soundBatch []sound

	// soundBatchQueue is the queue of sound batches that are ready to play.
	soundBatchQueue chan []sound

	// done is channel used to shutdown the audioManager and finish playing sounds.
	done chan bool
}

func newAudioManager() *audioManager {
	return &audioManager{
		soundBatchQueue: make(chan []sound, 100),
		done:            make(chan bool),
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

		play := func(buf []int16) {
			stream, err := portaudio.OpenDefaultStream(0 /*input channels */, 2, 44100, len(buf), buf)
			logFatalIfErr("portaudio.OpenDefaultStream", err)
			logFatalIfErr("stream.Start", stream.Start())
			stream.Write()
			logFatalIfErr("stream.Stop", stream.Stop())
			logFatalIfErr("stream.Close", stream.Close())
		}

		for sb := range a.soundBatchQueue {
			switch {
			// Play the specific buffer if the batch has only one sound.
			case len(sb) == 1:
				play(buffers[sb[0]])

			// Combine the buffers if the batch has multiple sounds.
			case len(sb) > 1:
				amp := 1.0 / float32(len(sb))
				log.Printf("mixing %d streams (%.2f)", len(sb), amp)

				bufSize := 0
				for _, s := range sb {
					if bs := len(buffers[s]); bs > bufSize {
						bufSize = bs
					}
				}

				buf := make([]int16, bufSize)
				for _, s := range sb {
					for i, v := range buffers[s] {
						buf[i] += int16(float32(v) * amp)
					}
				}
				play(buf)
			}
		}

		a.done <- true
	}()
}

func (a *audioManager) play(s sound) {
	a.soundBatch = append(a.soundBatch, s)
}

func (a *audioManager) flush() {
	if len(a.soundBatch) > 0 {
		a.soundBatchQueue <- a.soundBatch
		a.soundBatch = nil
	}
}

func (a *audioManager) stop() {
	a.flush()
	close(a.soundBatchQueue)
	<-a.done
	close(a.done)
}
