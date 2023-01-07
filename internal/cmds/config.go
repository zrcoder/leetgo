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
	langShortKey     = "l"
	codeLangShortKey = "c"
	cacheDirShortKey = "d"

	langUsage     = "language for the app"
	codeLangUsage = "programing language to resolve the problems"
	cacheDirUsage = "cache directory"
)

var cfg = &config.Config{}

var Config = &cli.Command{
	Name:      "config",
	UsageText: "show current config if no flags, use flags to set",
	Flags:     []cli.Flag{langFlag, codeLangFlag, cacheDirFlag},
	Action:    configAction,
}

var langFlag = &cli.StringFlag{
	Name:        config.LangKey,
	Aliases:     []string{langShortKey},
	Value:       config.DefaultLanguage,
	Usage:       langUsage,
	Destination: &cfg.Language,
}

var codeLangFlag = &cli.StringFlag{
	Name:        config.CodeLangKey,
	Aliases:     []string{codeLangShortKey},
	Value:       config.DefaultCodeLang,
	Usage:       codeLangUsage,
	Destination: &cfg.CodeLang,
}

var cacheDirFlag = &cli.StringFlag{
	Name:        config.CacheDirKey,
	Aliases:     []string{cacheDirShortKey},
	Value:       config.DefaultCacheDir,
	Usage:       cacheDirUsage,
	Destination: &cfg.CacheDir,
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
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	buf := &strings.Builder{}
	buf.WriteString("|item|description|\n| --- | --- |\n")
	buf.WriteString(fmt.Sprintf("|%s|%v|\n", cfg.Language, langUsage))
	buf.WriteString(fmt.Sprintf("|%s|%v|\n", cfg.CodeLang, codeLangUsage))
	buf.WriteString(fmt.Sprintf("|%s|%v|\n", cfg.CacheDir, cacheDirUsage))
	fmt.Println(render.Success("CURRENT:"))
	fmt.Println(render.MarkDown(buf.String()))
	return nil
}

func setConfig(context *cli.Context, localFlags []string) error {
	curCfg, err := config.Get()
	if err != nil {
		return err
	}
	adapt(curCfg, cfg)
	if curCfg.CodeLang == config.CodeLangGoShort {
		curCfg.CodeLang = config.CodeLangGo
	}
	return config.Write(curCfg)
}

func adapt(dest, src *config.Config) {
	if src.Language != "" {
		dest.Language = src.Language
	}
	if src.CacheDir != "" {
		dest.CacheDir = src.CacheDir
	}
	if src.CodeLang != "" {
		dest.CodeLang = src.CodeLang
	}
}
