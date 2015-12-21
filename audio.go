package main

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

var soundQueue = make(chan int, 100)

func playSound() {
	soundQueue <- 0
}

func processSounds() {
	wav, err := decodeWAV(newAssetReader("data/move.wav"))
	logFatalIfErr("decodeWAV", err)
	log.Printf("move.wav: %+v", wav)

	out := make([]int16, len(wav.data)/2)
	for i := 0; i < len(out); i++ {
		out[i] = int16(wav.data[i*2])
		out[i] += int16(wav.data[i*2+1]) << 8
	}

	const (
		numChannels = 2
		sampleRate  = 44100
	)
	stream, err := portaudio.OpenDefaultStream(0 /*input channels */, numChannels, sampleRate, len(out), out)
	logFatalIfErr("portaudio.OpenDefaultStream", err)
	defer stream.Close()

	for {
		select {
		case <-soundQueue:
			logFatalIfErr("stream.Start", stream.Start())
			stream.Write()
			logFatalIfErr("stream.Stop", stream.Stop())
		}
	}
}
