package cmds

import (
	"errors"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var ErrNeedQuestionId = errors.New("need question id")

var View = &cli.Command{
	Name:      "view",
	Usage:     "view a question by id",
	UsageText: "leetgo view 127",
	Action:    viewAction,
}

func viewAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return ErrNeedQuestionId
	}

	id := strings.Join(context.Args().Slice(), " ")
	return comp.NewViewer(id).Run()
}
