package cmds

import (
	"github.com/urfave/cli/v2"
)

const (
	markdownFlagName = "markdown"
)

var Book = &cli.Command{
	Name:        "book",
	Usage:       "generate a book for the picked questions",
	Subcommands: []*cli.Command{bookGen, bookServe},
}

var bookGen = &cli.Command{
	Name:   "gen",
	Usage:  "generate a book for the picked questions",
	Flags:  []cli.Flag{bookGenMarkdownFlag},
	Action: bookGenAction,
}

var bookServe = &cli.Command{
	Name:    "serve",
	Aliases: []string{"server"},
	Usage:   "generate a book for the picked questions and render it on the browser",
	Action:  bookServeAction,
}

var bookGenMarkdownFlag = &cli.BoolFlag{
	Name:    markdownFlagName,
	Aliases: []string{"m"},
}

func bookGenAction(context *cli.Context) error {
	return nil
}

func bookServeAction(context *cli.Context) error {
	return nil
}
