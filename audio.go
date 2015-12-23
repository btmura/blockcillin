package main

import (
	"log"
	"time"

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
	// soundBuffers is a map from sound to audio buffer.
	soundBuffers map[sound][]int16

	// soundQueue is a channel used to schedule the next sounds to play.
	soundQueue chan sound

	// done is a channel used to coordinate shutdown. Use the close method instead.
	done chan bool
}

func newAudioManager() *audioManager {
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

	soundBuffers := map[sound][]int16{
		soundMove:   makeBuffer("data/move.wav"),
		soundSelect: makeBuffer("data/select.wav"),
		soundSwap:   makeBuffer("data/swap.wav"),
		soundClear:  makeBuffer("data/clear.wav"),
	}

	return &audioManager{
		soundBuffers: soundBuffers,
		soundQueue:   make(chan sound, 100),
		done:         make(chan bool),
	}
}

func (a *audioManager) start() error {
	go func() {
		const (
			numInputChannels  = 0     /* no input - not recording */
			numOutputChannels = 2     /* stereo output */
			sampleRate        = 44100 /* samples per second */
			intervalMs        = 100
			framesPerBuffer   = sampleRate / 1000.0 * intervalMs /*  len(buf) == numChannels * framesPerBuffer */
			outputBufferSize  = numOutputChannels * framesPerBuffer
		)

		// Temporary buffers to read and write the next batch of data.
		tmpIn := make([]int16, outputBufferSize)
		tmpOut := make([]int16, outputBufferSize)

		// outputRingBuffer is a buffer that the PortAudio callback will read data from.
		// We periodically wake up to write data to it and PortAudio will wake up to read from it.
		outputRingBuffer := newRingBuffer(outputBufferSize * 10)

		// process is the callback that PortAudio will call when it needs audio data.
		process := func(out []int16) {
			outputRingBuffer.pop(tmpIn)
			for i := 0; i < len(out); i++ {
				out[i] = tmpIn[i]
			}
		}

		stream, err := portaudio.OpenDefaultStream(numInputChannels, numOutputChannels, sampleRate, framesPerBuffer, process)
		logFatalIfErr("portaudio.OpenDefaultStream", err)
		defer func() {
			logFatalIfErr("stream.Close", stream.Close())
		}()

		logFatalIfErr("stream.Start()", stream.Start())
		defer func() {
			// stream.Stop blocks until all samples have been played.
			logFatalIfErr("stream.Stop", stream.Stop())
			a.done <- true
		}()

		var active [][]int16
		quit := false

	loop:
		for {
			select {
			case <-time.After(intervalMs * time.Millisecond):
				// Play whatever sounds are in the queue at the time.
				n := len(a.soundQueue)
				for i := 0; i < n; i++ {
					active = append(active, a.soundBuffers[<-a.soundQueue])
				}

				// Fill temporary buffer with any active sounds buffers.
				for i := 0; i < len(tmpOut); i++ {
					// Combine active signals together.
					var v int16
					for j := 0; j < len(active); j++ {
						// Remove any buffers if they have no more samples.
						if len(active[j]) == 0 {
							active = append(active[:j], active[j+1:]...)
							j--
							continue
						}
						v += active[j][0]
						active[j] = active[j][1:]
					}
					tmpOut[i] = v
				}
				outputRingBuffer.push(tmpOut...)

				// Only quit waking until there are no more streams to play.
				if quit && len(active) == 0 {
					break loop
				}

			case <-a.done:
				close(a.soundQueue) // Prevent any new sounds from being scheduled.
				quit = true
			}
		}
	}()

	return nil
}

func (a *audioManager) play(s sound) {
	a.soundQueue <- s
}

func (a *audioManager) close() {
	a.done <- true // Signal audio manager to quit and play its last samples.
	<-a.done       // Wait for notification that all samples have been played.
}
