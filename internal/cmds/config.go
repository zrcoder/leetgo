package cmds

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/render"
)

const (
	langShortKey     = "l"
	codeLangShortKey = "c"

	langUsage     = "language for the app"
	codeLangUsage = "programing language to resolve the problems"
)

var Config = &cli.Command{
	Name:      "config",
	UsageText: "show current config if no flags, use flags to set",
	Flags:     []cli.Flag{langFlag, codeLangFlag},
	Action:    configAction,
}

var langFlag = &cli.StringFlag{
	Name:    config.LangKey,
	Aliases: []string{langShortKey},
	Value:   config.DefaultLanguage,
	Usage:   langUsage,
	Action: func(ctx *cli.Context, s string) error {
		if !config.AllowedLang[s] {
			return config.ErrInvalidLan
		}
		return nil
	},
}

var codeLangFlag = &cli.StringFlag{
	Name:    config.CodeLangKey,
	Aliases: []string{codeLangShortKey},
	Value:   config.DefaultCodeLang,
	Usage:   codeLangUsage,
	Action: func(ctx *cli.Context, s string) error {
		if !config.AllowedCodeLang[s] {
			return config.ErrInvalidCodeLan
		}
		return nil
	},
}

func configAction(context *cli.Context) error {
	localFlags := context.LocalFlagNames()
	if len(localFlags) == 0 {
		return showConfig()
	}
	err := setConfig(context, localFlags)
	if err != nil {
		return err
	}
	fmt.Println(render.Info("Success"))
	return nil
}

func showConfig() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	buf := &strings.Builder{}
	buf.WriteString("|flag|current value|description|\n| --- | --- | --- |\n")
	buf.WriteString(fmt.Sprintf("|-%s, --%s|%s|%v|\n",
		langShortKey, config.LangKey, render.Info(cfg.Language), langUsage))
	buf.WriteString(fmt.Sprintf("|-%s, --%s|%s|%v|\n",
		codeLangShortKey, config.CodeLangKey, render.Info(cfg.CodeLang), codeLangUsage))
	fmt.Println(render.MarkDown(buf.String()))
	return nil
}

func setConfig(context *cli.Context, localFlags []string) error {
	curCfg, err := config.Get()
	if err != nil {
		return err
	}
	for _, v := range localFlags {
		switch v {
		case config.LangKey:
			curCfg.Language = context.String(config.LangKey)
		case config.CodeLangKey:
			curCfg.CodeLang = context.String(config.CodeLangKey)
		}
	}
	return config.Write(curCfg)
}
