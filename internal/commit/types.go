package commit

import "time"

type Commit struct {
	Date    time.Time
	Author  string
	Message string
	Parents []string
	Tree    string
}
