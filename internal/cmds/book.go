package cmds

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
	"github.com/zrcoder/leetgo/internal/model"
)

const (
	genMarkdownFlagName      = "markdown"
	genMarkDownFlagShortName = "m"
	sortDocFlagName          = "sortBy"
	sortDocFlagShortName     = "s"
	genRepoFlagName          = "repo"
	genRepoFlagShortName     = "r"
	servePortFlagName        = "port"
	servePortShortFlagName   = "p"
	serveSourceFlagName      = "src"
	serveSourceFlagShortName = "s"
)

var Book = &cli.Command{
	Name:        "book",
	Usage:       "generate or serve a book for the picked questions",
	Subcommands: []*cli.Command{bookServe},
}

var bookServe = &cli.Command{
	Name:    "serve",
	Aliases: []string{"server"},
	Usage:   "generate a book for the picked questions/solutions, and serve it on localhost",
	Flags:   []cli.Flag{sortFlag, servePortFlag},
	Action:  bookServeAction,
}

var sortFlag = &cli.StringFlag{
	Name:    sortDocFlagName,
	Aliases: []string{sortDocFlagShortName},
	Value:   model.SortByTime,
	Usage:   "sort by time or title when generate or serve the book",
}

var servePortFlag = &cli.StringFlag{
	Name:     servePortFlagName,
	Aliases:  []string{servePortShortFlagName},
	Required: true,
}

func bookServeAction(context *cli.Context) error {
	sortBy := context.String(sortDocFlagName)
	if sortBy != model.SortByTime && sortBy != model.SortByTitle {
		return fmt.Errorf("we only suport sort the questions by %s/%s", model.SortByTime, model.SortByTime)
	}
	port := context.String(servePortFlagName)
	return comp.NewBook(sortBy, port).Run()
}
