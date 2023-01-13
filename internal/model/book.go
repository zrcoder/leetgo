package model

import (
	"time"
)

const (
	SortByTime  = "time"
	SortByTitle = "Title"
)

// Doc is a chapter/directory or an article/file
type Doc struct {
	Title           string
	Time            time.Time
	MarkdownContent []byte
}
