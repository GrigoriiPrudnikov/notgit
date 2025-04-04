package blob

import "os"

type File struct {
	Mode    os.FileMode
	Name    string
	Hash    string
	Content []byte
}
