package cmds

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
	"github.com/zrcoder/leetgo/internal/model"
)

const (
	sortDocFlagName        = "sort"
	sortDocFlagShortName   = "s"
	servePortFlagName      = "port"
	servePortShortFlagName = "p"
)

var Book = &cli.Command{
	Name:   "book",
	Usage:  "generate and serve a web book for the picked questions",
	Flags:  []cli.Flag{sortFlag, servePortFlag},
	Action: bookAction,
}

var sortFlag = &cli.StringFlag{
	Name:    sortDocFlagName,
	Aliases: []string{sortDocFlagShortName},
	Value:   model.SortByTime,
	Usage:   "sort the docs by time or title",
}

var servePortFlag = &cli.StringFlag{
	Name:    servePortFlagName,
	Aliases: []string{servePortShortFlagName},
	Value:   "9999",
	Usage:   "the local port for serving",
}

func bookAction(context *cli.Context) error {
	sortBy := context.String(sortDocFlagName)
	if sortBy != model.SortByTime && sortBy != model.SortByTitle {
		return fmt.Errorf("we only suport sort the questions by %s/%s", model.SortByTime, model.SortByTitle)
	}
	port := context.String(servePortFlagName)
	return comp.NewBook(sortBy, port).Run()
}
