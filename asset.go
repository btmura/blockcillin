package main

import (
	"bytes"
	"io"

	"github.com/btmura/blockcillin/internal/asset"
)

func newAssetReader(name string) io.Reader {
	return bytes.NewReader(asset.MustAsset(name))
}

func assetString(name string) string {
	return string(asset.MustAsset(name))
}
