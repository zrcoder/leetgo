package cmds

import (
	"errors"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
	"github.com/zrcoder/leetgo/internal/log"
)

var ErrNeedQuestionId = errors.New("need question id")

var View = &cli.Command{
	Name:      "view",
	Usage:     "view a question or it's solution by id",
	UsageText: "leetgo view [-s] 127",
	Flags:     []cli.Flag{solutionFlag},
	Action:    viewAction,
}

var solutionFlag = &cli.BoolFlag{
	Name:    "solution",
	Aliases: []string{"s"},
	Usage:   "view the most voted solution",
}

func viewAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return ErrNeedQuestionId
	}

	id := strings.Join(context.Args().Slice(), " ")
	solution := context.Bool(solutionFlag.Name)
	log.Debug("view<", id, "> solution?", solution)
	return comp.NewViewer(id, solution).Run()
}
