package cmds

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/utils/render"
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
	editorUsage   = "editor for coding: vim/neovim/emacs/vscode/cursor/zed"
)

var Config = &cli.Command{
	Name:   "config",
	Usage:  "init or config your leetgo project",
	Flags:  []cli.Flag{langFlag, codeLangFlag, editorFlag},
	Action: configAction,
}

var (
	langFlag = &cli.StringFlag{
		Name:    langKey,
		Aliases: []string{langShortKey},
		Value:   config.DefaultLanguage,
		Usage:   langUsage,
		Action: func(ctx *cli.Context, s string) error {
			if !config.SrpportedLang(s) {
				return config.ErrInvalidLang
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
				return config.ErrInvalidCodeLang
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
		buf.WriteString(fmt.Sprintf(rowTmp, langShortKey, langKey, config.DisplayLang(cfg.Language), langUsage))
		buf.WriteString(fmt.Sprintf(rowTmp, codeLangShortKey, codeLangKey, cfg.CodeLang, codeLangUsage))
		buf.WriteString(fmt.Sprintf(rowTmp, editorShortKey, editorKey, cfg.Editor, editorUsage))
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
