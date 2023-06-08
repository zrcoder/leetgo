package cmds

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/local"
)

var Edit = &cli.Command{
	Name:      "edit",
	Usage:     "edit codes to solve the question",
	UsageText: "leetgo edit 1",
	Action:    editAction,
}

func editAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return errors.New("please pass the question id")
	}

	cfg, err := config.Get()
	if err != nil {
		return err
	}

	id := context.Args().First()
	dir := local.GetDir(cfg, id)
	_, err = os.Stat(dir)
	if err != nil {
		return err
	}

	cmd := exec.Command("vim", fmt.Sprintf("+/%s", local.CodeStartFlag4Editor), local.GetCodeFile(cfg, id))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
