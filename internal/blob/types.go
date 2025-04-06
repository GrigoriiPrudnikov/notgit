package blob

import "os"

type Blob struct {
	Mode    os.FileMode
	Path    string
	Hash    string
	Content []byte
}
