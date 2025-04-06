package blob

import "os"

type Blob struct {
	Mode    os.FileMode
	Name    string
	Hash    string
	Content []byte
}
