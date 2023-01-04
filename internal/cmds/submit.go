package cmds

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/remote"
)

var Submit = &cli.Command{
	Name:   "submit",
	Action: submitAction,
}

func submitAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return ErrNeedQuestionId
	}

	id := strings.Join(context.Args().Slice(), " ")
	submitReq, err := local.GetAnswer(id)
	if err != nil {
		if err == local.ErrNotFound {
			err = fmt.Errorf("%w, you should pick a question and type your answer first", err)
		}
		return err
	}

	info, err := remote.Submit(submitReq, 0)
	if err != nil {
		return err
	}

	fmt.Println(info)
	return nil
}
