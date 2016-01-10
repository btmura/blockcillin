package audio

import (
	"encoding/binary"
	"fmt"
	"io"
)

// wav is decoded WAV data.
//
// https://github.com/verdverm/go-wav/blob/master/wav.go
// http://soundfile.sapp.org/doc/WaveFormat/
type wav struct {
	chunkID   [4]byte
	chunkSize uint32
	format    [4]byte

	subchunk1ID   [4]byte
	subchunk1Size uint32
	audioFormat   uint16
	numChannels   uint16
	sampleRate    uint32
	byteRate      uint32
	blockAlign    uint16
	bitsPerSample uint16

	subchunk2ID   [4]byte
	subchunk2Size uint32
	data          []byte
}

func decodeWAV(r io.Reader) (*wav, error) {
	w := &wav{}
	var err error

	read := func(order binary.ByteOrder, dst interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(r, order, dst)
	}

	check := func(name string, value [4]byte, want string) {
		if err != nil {
			return
		}
		if str := string(value[:]); str != want {
			err = fmt.Errorf("%s should be %q, got %q", name, want, str)
		}
	}

	// Chunk 1 - RIFF

	read(binary.BigEndian, &w.chunkID)
	read(binary.LittleEndian, &w.chunkSize)
	read(binary.BigEndian, &w.format)
	check("chunkID", w.chunkID, "RIFF")
	check("format", w.format, "WAVE")

	// Chunk 2 - fmt

	read(binary.BigEndian, &w.subchunk1ID)
	read(binary.LittleEndian, &w.subchunk1Size)
	read(binary.LittleEndian, &w.audioFormat)
	read(binary.LittleEndian, &w.numChannels)
	read(binary.LittleEndian, &w.sampleRate)
	read(binary.LittleEndian, &w.byteRate)
	read(binary.LittleEndian, &w.blockAlign)
	read(binary.LittleEndian, &w.bitsPerSample)
	check("subchunk1ID", w.subchunk1ID, "fmt ")

	// Chunk 3 - data

	read(binary.BigEndian, &w.subchunk2ID)
	read(binary.LittleEndian, &w.subchunk2Size)
	w.data = make([]byte, w.subchunk2Size)
	read(binary.LittleEndian, &w.data)
	check("subchunk2ID", w.subchunk2ID, "data")

	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *wav) String() string {
	return fmt.Sprintf("numChannels: %d sampleRate: %d bitsPerSample: %d len(data): %d", w.numChannels, w.sampleRate, w.bitsPerSample, len(w.data))
}
