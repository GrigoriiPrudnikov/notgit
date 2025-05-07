package utils

import (
	"crypto/sha256"
	"fmt"
)

func Hash(kind string, content []byte) string {
	header := fmt.Sprintf("%s %d\x00\n", kind, len(content))
	data := append([]byte(header), content...)
	sum := sha256.Sum256(data)

	return fmt.Sprintf("%x", sum)
}
