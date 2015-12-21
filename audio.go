package main

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

const (
	numChannels = 2
	sampleRate  = 44100
)

func playSounds() {
	wav, err := decodeWAV(newAssetReader("data/move.wav"))
	logFatalIfErr("decodeWAV", err)
	log.Printf("WAV: %+v", wav)

	out := make([]int16, 8192)
	stream, err := portaudio.OpenDefaultStream(0 /*input channels */, numChannels, sampleRate, len(out), out)
	logFatalIfErr("portaudio.OpenDefaultStream", err)
	defer stream.Close()

	logFatalIfErr("stream.Start", stream.Start())
	defer func() {
		logFatalIfErr("stream.Stop", stream.Stop())
	}()

	j := 0
	for remaining := len(wav.data) / 2; remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		for i := 0; i < len(out); i++ {
			out[i] = int16(wav.data[j])
			out[i] += int16(wav.data[j+1]) << 8
			j += 2
		}
		stream.Write()
	}
}
