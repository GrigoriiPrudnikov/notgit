package commit

type Commit struct {
	Hash    string
	Time    int64
	Offset  string
	Author  string
	Message string
	Tree    string
	Parents []string
}
