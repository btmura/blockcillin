package audio

import (
	"encoding/binary"
	"fmt"
	"io"
)

// WAV is decoded WAV data.
//
// https://github.com/verdverm/go-wav/blob/master/wav.go
// http://soundfile.sapp.org/doc/WaveFormat/
type WAV struct {
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

func decodeWAV(r io.Reader) (*WAV, error) {
	wav := &WAV{}

	// Chunk 1 - RIFF

	if err := binary.Read(r, binary.BigEndian, &wav.chunkID); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.chunkSize); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &wav.format); err != nil {
		return nil, err
	}

	if chunkID := string(wav.chunkID[:]); chunkID != "RIFF" {
		return nil, fmt.Errorf(`chunkID should be "RIFF", got %q`, chunkID)
	}
	if format := string(wav.format[:]); format != "WAVE" {
		return nil, fmt.Errorf(`format should be "WAVE", got %q`, format)
	}

	// Chunk 2 - fmt

	if err := binary.Read(r, binary.BigEndian, &wav.subchunk1ID); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.subchunk1Size); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.audioFormat); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.numChannels); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.sampleRate); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.byteRate); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.blockAlign); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.bitsPerSample); err != nil {
		return nil, err
	}

	if subchunk1ID := string(wav.subchunk1ID[:]); subchunk1ID != "fmt " {
		return nil, fmt.Errorf(`subchunk1ID should be "fmt ", got %q`, subchunk1ID)
	}

	// Chunk 3 - data

	if err := binary.Read(r, binary.BigEndian, &wav.subchunk2ID); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.subchunk2Size); err != nil {
		return nil, err
	}
	wav.data = make([]byte, wav.subchunk2Size)
	if err := binary.Read(r, binary.LittleEndian, &wav.data); err != nil {
		return nil, err
	}

	if subchunk2ID := string(wav.subchunk2ID[:]); subchunk2ID != "data" {
		return nil, fmt.Errorf(`subchunk2ID should be "fmt ", got %q`, subchunk2ID)
	}

	return wav, nil
}

func (w *WAV) String() string {
	return fmt.Sprintf("chunkID: %s chunkSize: %d format: %s "+
		"subchunk1ID: %s subchunk1Size: %d audioFormat: %d numChannels: %d sampleRate: %d byteRate: %d blockAlign: %d bitsPerSample: %d "+
		"subchunk2ID: %s subchunk2Size: %d len(data): %d",
		w.chunkID, w.chunkSize, w.format,
		w.subchunk1ID, w.subchunk1Size, w.audioFormat, w.numChannels, w.sampleRate, w.byteRate, w.blockAlign, w.bitsPerSample,
		w.subchunk2ID, w.subchunk2Size, len(w.data))
}
