package cmds

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
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

	langUsage     = "language, en or cn"
	codeLangUsage = "programing language"
	editorUsage   = "editor to use, vim or neovim"
)

var Config = &cli.Command{
	Name:      "config",
	Usage:     "show or set config of your leetgo project",
	UsageText: "show current config if no flags, use flags to set",
	Flags:     []cli.Flag{langFlag, codeLangFlag, editorFlag},
	Action:    configAction,
}

var (
	langFlag = &cli.StringFlag{
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

	codeLangFlag = &cli.StringFlag{
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

	editorFlag = &cli.StringFlag{
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
)

func configAction(context *cli.Context) error {
	show := func(cfg *config.Config) {
		buf := &strings.Builder{}
		buf.WriteString("|flag|current value|description|\n")
		buf.WriteString("| --- | --- | --- |\n")
		const rowTmp = "|-%s, --%s|%s|%v|\n"
		buf.WriteString(fmt.Sprintf(rowTmp, langShortKey, langKey, render.Info(config.DisplayLang(cfg.Language)), langUsage))
		buf.WriteString(fmt.Sprintf(rowTmp, codeLangShortKey, codeLangKey, render.Info(cfg.CodeLang), codeLangUsage))
		buf.WriteString(fmt.Sprintf(rowTmp, editorShortKey, editorKey, render.Info(cfg.Editor), editorUsage))
		fmt.Println(render.MarkDown(buf.String()))
	}
	cfg := &config.Config{}
	for _, v := range context.LocalFlagNames() {
		switch v {
		case langKey:
			cfg.Language = context.String(langKey)
		case codeLangKey:
			cfg.CodeLang = context.String(codeLangKey)
		case editorKey:
			cfg.Editor = context.String(editorKey)
		}
	}
	return comp.NewConfiger(cfg, len(context.LocalFlagNames()) > 0, show).Run()
}
