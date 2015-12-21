package main

import (
	"log"
	"time"

	"github.com/gordonklaus/portaudio"
)

func playSounds() {
	wav, err := decodeWAV(newAssetReader("data/sounds.wav"))
	logFatalIfErr("decodeWAV", err)
	log.Printf("WAV: %+v", wav)

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
	time.Sleep(3 * time.Second)
	defer func() {
		logFatalIfErr("s.Stop", stream.Stop())
	}()

}
