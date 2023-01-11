package cmds

import (
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
	Subcommands: []*cli.Command{bookGen, bookServe},
}

var bookGen = &cli.Command{
	Name:   "gen",
	Usage:  "generate a book for the picked questions",
	Flags:  []cli.Flag{sortFlag, genMarkdownFlag, genRepoFlag},
	Action: bookGenAction,
}

var bookServe = &cli.Command{
	Name:    "serve",
	Aliases: []string{"server"},
	Usage:   "serve the generated book html resources on localhost",
	Flags:   []cli.Flag{serveSourceFlag, servePortFlag},
	Action:  bookServeAction,
}

var genMarkdownFlag = &cli.BoolFlag{
	Name:    genMarkdownFlagName,
	Aliases: []string{genMarkDownFlagShortName},
	Usage:   "generate markdowns or not when generate or serve the book",
}

var sortFlag = &cli.StringFlag{
	Name:    sortDocFlagName,
	Aliases: []string{sortDocFlagShortName},
	Value:   "time",
	Usage:   "sort by time or title when generate or serve the book",
}

var genRepoFlag = &cli.StringFlag{
	Name:    genRepoFlagName,
	Aliases: []string{genRepoFlagShortName},
	Usage:   "if you want to publish your book as github/gitee pages, you should pass the repo",
}

var serveSourceFlag = &cli.StringFlag{
	Name:     serveSourceFlagName,
	Aliases:  []string{serveSourceFlagShortName},
	Usage:    "the source html resources' directory",
	Required: true,
}

var servePortFlag = &cli.StringFlag{
	Name:     servePortFlagName,
	Aliases:  []string{servePortShortFlagName},
	Required: true,
}

func bookGenAction(context *cli.Context) error {
	sortBy := context.String(sortDocFlagName)
	repo := context.String(genRepoFlagName)
	genMarkdown := context.Bool(genMarkdownFlagName)
	meta := &model.BookMeta{
		SortBy:       sortBy,
		Repo:         repo,
		GenMarkdowns: genMarkdown,
	}
	return comp.NewBookGenerator(meta).Run()
}

func bookServeAction(context *cli.Context) error {
	htmlSrc := context.String(serveSourceFlagName)
	port := context.String(servePortFlagName)
	return comp.NewServer(htmlSrc, port).Run()
}
