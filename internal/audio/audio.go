package audio

import (
	"log"
	"time"

	"github.com/btmura/blockcillin/internal/asset"
	"github.com/gordonklaus/portaudio"
)

// Sound is an enum that identifies a short sound in the game.
type Sound int

const (
	SoundMove Sound = iota
	SoundSelect
	SoundSwap
	SoundClear
)

// Play plays the given sound. It is overridden by Init.
var Play = func(s Sound) {}

// Terminate shuts down the audio system. It is overridden by Init.
var Terminate = func() {}

// Init starts the audio system, loads sound assets, and starts the sound loop.
func Init() error {
	if err := portaudio.Initialize(); err != nil {
		return err
	}

	log.Printf("PortAudio version: %d %s", portaudio.Version(), portaudio.VersionText())

	var err error
	makeBuffer := func(name string) []int16 {
		if err != nil {
			return nil
		}

		var w *wav
		if w, err = decodeWAV(asset.MustReader(name)); err != nil {
			return nil
		}

		log.Printf("%s: %+v", name, w)
		buf := make([]int16, len(w.data)/2)
		for i := 0; i < len(buf); i++ {
			buf[i] = int16(w.data[i*2])
			buf[i] += int16(w.data[i*2+1]) << 8
		}
		return buf
	}

	soundBuffers := map[Sound][]int16{
		SoundMove:   makeBuffer("data/move.wav"),
		SoundSelect: makeBuffer("data/select.wav"),
		SoundSwap:   makeBuffer("data/swap.wav"),
		SoundClear:  makeBuffer("data/clear.wav"),
	}
	if err != nil {
		return err
	}

	soundQueue := make(chan Sound, 10)
	Play = func(s Sound) {
		soundQueue <- s
	}

	done := make(chan bool)
	Terminate = func() {
		done <- true
		<-done
		close(done)
		logFatalIfErr("portaudio.Terminate", portaudio.Terminate())
	}

	go func() {
		const (
			numInputChannels  = 0     /* zero input - no recording */
			numOutputChannels = 2     /* stereo output */
			sampleRate        = 44100 /* samples per second */
			intervalMs        = 100
			framesPerBuffer   = sampleRate / 1000.0 * intervalMs /*  len(buf) == numChannels * framesPerBuffer */
			outputBufferSize  = numOutputChannels * framesPerBuffer
		)

		// Temporary buffers to read and write the next audio batch.
		tmpIn := make([]int16, outputBufferSize)
		tmpOut := make([]int16, outputBufferSize)

		// outputRingBuffer is a buffer that the PortAudio callback will read data from,
		// and that we will periodically wake up to write new data to.
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
			done <- true
		}()

		var active [][]int16
		quit := false

	loop:
		for {
			select {
			case <-time.After(intervalMs * time.Millisecond):
				// Play whatever sounds are in the queue at the time.
				n := len(soundQueue)
				for i := 0; i < n; i++ {
					active = append(active, soundBuffers[<-soundQueue])
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

			case <-done:
				close(soundQueue) // Prevent any new sounds from being scheduled.
				quit = true
			}
		}
	}()

	return nil
}

func logFatalIfErr(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", tag, err)
	}
}
