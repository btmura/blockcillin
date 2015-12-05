package main

import (
	"bytes"
	"io"
)

func newAssetReader(name string) io.Reader {
	return bytes.NewReader(MustAsset(name))
}

func assetString(name string) string {
	return string(MustAsset(name))
}
