package cmds

import (
	"errors"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var ErrNeedQuestionId = errors.New("need question id")

var Pick = &cli.Command{
	Name:      "pick",
	Usage:     "pick a question by id",
	UsageText: "leetgo pick 127",
	Action:    pickAction,
}

func pickAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return ErrNeedQuestionId
	}

	id := strings.Join(context.Args().Slice(), " ")
	return comp.NewPicker(id).Run()
}
