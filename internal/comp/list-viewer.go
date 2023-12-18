package comp

import (
	"sort"

	"github.com/zrcoder/tdoc"
	tmodel "github.com/zrcoder/tdoc/model"
)

type listViewer struct {
	sortby  string
	reverse bool
}

func (l *listViewer) Run() error {
	docs, err := getDocsFromLocal()
	if err != nil {
		return err
	}

	switch l.sortby {
	case "time":
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].ModTime.Before(docs[j].ModTime)
		})
	case "title":
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].Title < docs[j].Title
		})
	}
	if l.reverse {
		i, j := 0, len(docs)-1
		for i < j {
			docs[i], docs[j] = docs[j], docs[i]
			i++
			j--
		}
	}
	return tdoc.Run(docs, tmodel.Config{Title: "My Solutions"})
}
