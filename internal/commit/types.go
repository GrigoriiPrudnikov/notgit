package commit

type Commit struct {
	Time    int64
	Offset  string
	Author  string
	Message string
	Tree    string
	Parents []string
}
