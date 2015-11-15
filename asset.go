package main

import (
	"bytes"
	"io"
)

func newAssetReader(name string) (io.Reader, error) {
	a, err := Asset(name)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(a), nil
}

func getStringAsset(name string) (string, error) {
	a, err := Asset(name)
	if err != nil {
		return "", err
	}
	return string(a), nil
}
