package blob

import (
	"crypto/sha256"
	"fmt"
)

func Hash(b *Blob) {
	content := b.Content
	header := fmt.Sprintf("blob %d\x00\n", len(content))
	blob := append([]byte(header), content...)
	hash := sha256.Sum256(blob)
	hex := fmt.Sprintf("%x", hash)

	b.Hash = hex
}
