package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
)

func Compress(b []byte, variant string) []byte {
	header := fmt.Sprintf("%s %d\x00\n", variant, len(b))
	content := append([]byte(header), b...)

	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(content)
	w.Close()

	return compressed.Bytes()
}

func Decompress(b []byte) []byte {
	return b
}
