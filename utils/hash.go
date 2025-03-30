package utils

import (
	"crypto/sha256"
	"fmt"
)

func Hash(b []byte) string {
	header := fmt.Sprintf("blob %d\x00\n", len(b))
	blob := append([]byte(header), b...)
	hash := sha256.Sum256(blob)
	hex := fmt.Sprintf("%x", hash)

	return hex
}
