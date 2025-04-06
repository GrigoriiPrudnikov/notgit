package blob

type Blob struct {
	Permission string
	Path       string
	Hash       string
	Content    []byte
}
