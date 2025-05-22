package cmds

import (
	"errors"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
	"github.com/zrcoder/leetgo/internal/log"
)

var View = &cli.Command{
	Name:      "view",
	Usage:     "view questions or solutions",
	UsageText: "leetgo view [-s] 127",
	Flags:     []cli.Flag{solutionFlag, sortbyFlag, reverseFlag},
	Action:    viewAction,
}

var solutionFlag = &cli.BoolFlag{
	Name:    "solution",
	Aliases: []string{"s"},
	Usage:   "view the most voted solution",
}

var sortbyFlag = &cli.StringFlag{
	Name:    "sortby",
	Aliases: []string{"b"},
	Usage:   "sort questions by time/title, keep the original order by default",
	Action: func(ctx *cli.Context, s string) error {
		if s != "" && s != "time" && s != "title" {
			return errors.New("only `time` and `title` supported to sort")
		}
		return nil
	},
}

var reverseFlag = &cli.BoolFlag{
	Name:    "reverse",
	Aliases: []string{"r"},
	Usage:   "reverse sort",
}

func viewAction(context *cli.Context) error {
	args := context.Args()
	if args.Len() == 0 {
		return errors.New("no question id provided")
	}

	id := strings.Join(args.Slice(), " ")
	solution := context.Bool(solutionFlag.Name)
	log.Debug("view<", id, "> solution?", solution)
	if solution {
		return comp.NewSolutionViewer(id).Run()
	}
	return comp.NewSingleViewer(id).Run()
}
