package cmds

import (
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/tdoc"
)

const (
	SortByTime  = "time"
	SortByTitle = "title"
)

var Book = &cli.Command{
	Name:      "book",
	Usage:     "view all questiongs and codes in the local project as a book",
	UsageText: "leetgo book",
	Flags:     []cli.Flag{sortbyFlag, reverseFlag},
	Action:    bookAction,
}

var sortbyFlag = &cli.StringFlag{
	Name:    "sortby",
	Aliases: []string{"s"},
	Usage:   "sort questions by time/title, keep the original order by default",
}

var reverseFlag = &cli.BoolFlag{
	Name:    "reverse",
	Aliases: []string{"r"},
	Usage:   "reverse sort",
}

func bookAction(context *cli.Context) error {
	docs, err := local.GetMetaList()
	if err != nil {
		return err
	}
	sortby := context.String(sortbyFlag.Name)
	reverse := context.Bool(reverseFlag.Name)
	if sortby == SortByTime {
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].ModTime.Before(docs[j].ModTime)
		})
	} else if sortby == SortByTitle {
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].Title < docs[j].Title
		})
	}
	if reverse {
		i, j := 0, len(docs)-1
		for i < j {
			docs[i], docs[j] = docs[j], docs[i]
			i++
			j--
		}
	}
	return tdoc.Run(docs)
}
