package blob

type Blob struct {
	Type    string
	Hash    string
	Content []byte
	Size    int
}
