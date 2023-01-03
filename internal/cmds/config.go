package cmds

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/render"
)

const (
	lanShortKey      = "l"
	codeLangShortKey = "c"
	projectShortKey  = "p"
)

var Config = &cli.Command{
	Name:      "config",
	UsageText: "show current config if no flags, use flags to set",
	Flags:     []cli.Flag{langFlag, codeLangFlag, projectFlag},
	Action:    configAction,
}

var langFlag = &cli.StringFlag{
	Name:    config.LangKey,
	Aliases: []string{lanShortKey},
	Value:   config.DefaultLanguage,

	Usage: "language for the app",
}

var codeLangFlag = &cli.StringFlag{
	Name:    config.CodeLangKey,
	Aliases: []string{codeLangShortKey},
	Value:   config.DefaultCodeLang,
	Usage:   "programing language to resolve the problems",
}

var projectFlag = &cli.StringFlag{
	Name:    config.ProjectKey,
	Aliases: []string{projectShortKey},
	Value:   config.DefaultProjectDir,
	Usage:   "project directory",
}

func configAction(context *cli.Context) error {
	localFlags := context.LocalFlagNames()
	if len(localFlags) == 0 {
		return showConfig(context)
	}
	return setConfig(context, localFlags)
}

func showConfig(context *cli.Context) error {
	data, err := config.Read()
	if err != nil {
		return err
	}
	cfg := map[string]any{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	buf := &strings.Builder{}
	buf.WriteString("|key|value|\n| --- | --- |\n")
	for k, v := range cfg {
		buf.WriteString(fmt.Sprintf("|%s|%v|\n", k, v))
	}

	fmt.Println(render.Success("Current config:"))
	//md := fmt.Sprintf("Current config:\n\n---json\n\n```json\n%s\n```\n\n---json", data)
	fmt.Println(render.MarkDown(buf.String()))
	return nil
}

func setConfig(context *cli.Context, localFlags []string) error {
	info, err := config.Get()
	if err != nil {
		return err
	}

	for _, key := range localFlags {
		// ignore the short flags
		if _, exist := info[key]; !exist {
			continue
		}
		// take the long flags
		info[key] = context.String(key)
	}
	if info[config.CodeLangKey] == config.CodeLangGoShort {
		info[config.CodeLangKey] = config.CodeLangGo
	}
	return config.Write(info)
}
