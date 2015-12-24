package asset

import (
	"bytes"
	"io"
)

func MustReader(name string) io.Reader {
	return bytes.NewReader(MustAsset(name))
}

func MustString(name string) string {
	return string(MustAsset(name))
}
