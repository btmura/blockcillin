package asset

import (
	"bytes"
	"io"
)

func Reader(name string) (io.Reader, error) {
	data, err := Asset(name)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func String(name string) (string, error) {
	data, err := Asset(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
