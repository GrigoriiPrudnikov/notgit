package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
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

func Decompress(b []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
