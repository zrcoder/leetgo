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
	editorShortKey   = "e"
	langKey          = "lang"
	codeLangKey      = "code"
	editorKey        = "editor"

	langUsage     = "language for the app"
	codeLangUsage = "programing language to resolve the problems"
	editorUsage   = "editor to use"
)

var Config = &cli.Command{
	Name:      "config",
	UsageText: "show current config if no flags, use flags to set",
	Flags:     []cli.Flag{langFlag, codeLangFlag, editorFlag},
	Action:    configAction,
}

var langFlag = &cli.StringFlag{
	Name:    langKey,
	Aliases: []string{langShortKey},
	Value:   config.DefaultLanguage,
	Usage:   langUsage,
	Action: func(ctx *cli.Context, s string) error {
		if !config.SrpportedLang(s) {
			return config.ErrInvalidLan
		}
		return nil
	},
}

var codeLangFlag = &cli.StringFlag{
	Name:    codeLangKey,
	Aliases: []string{codeLangShortKey},
	Value:   config.DefaultCodeLang,
	Usage:   codeLangUsage,
	Action: func(ctx *cli.Context, s string) error {
		if !config.SupportedCodeLang(s) {
			return config.ErrInvalidCodeLan
		}
		return nil
	},
}

var editorFlag = &cli.StringFlag{
	Name:    editorKey,
	Aliases: []string{editorShortKey},
	Value:   config.DefaultEditor,
	Usage:   editorUsage,
	Action: func(ctx *cli.Context, s string) error {
		if !config.SupportedEditor(s) {
			return config.ErrUnSupporttedEditor
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
	const lineFmt = "|-%s, --%s|%s|%v|\n"
	buf.WriteString(fmt.Sprintf(lineFmt, langShortKey, langKey, render.Info(cfg.Language), langUsage))
	buf.WriteString(fmt.Sprintf(lineFmt, codeLangShortKey, codeLangKey, render.Info(cfg.CodeLang), codeLangUsage))
	buf.WriteString(fmt.Sprintf(lineFmt, editorShortKey, editorKey, render.Info(cfg.Editor), editorUsage))
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
		case langKey:
			curCfg.Language = context.String(langKey)
		case codeLangKey:
			curCfg.CodeLang = context.String(codeLangKey)
		case editorKey:
			curCfg.Editor = context.String(editorKey)
		}
	}
	return writeAndShow(curCfg)
}

func writeAndShow(cfg *config.Config) error {
	err := config.Write(cfg)
	if err != nil {
		return err
	}
	fmt.Println(render.Info("Succeed"))
	showConfig()
	return nil
}
