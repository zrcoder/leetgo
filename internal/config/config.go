package config

import (
	"encoding/json"
	"errors"
	"os"

	_ "github.com/zellyn/kooky/browser/all"

	"github.com/zrcoder/leetgo/internal/log"
)

const (
	DefaultLanguage = "en"
	DefaultCodeLang = "go"
	cnLanguage      = "cn"
	DefaultEditor   = "neovim"

	enDomain = "https://leetcode.com"
	cnDomain = "https://leetcode.cn"

	configFile = "leetgo.json"
)

var (
	errConfigNotExist     = errors.New("no config found, try `leetgo init`")
	ErrInvalidLan         = errors.New("only cn or en language supported")
	ErrInvalidCodeLan     = errors.New("not supported code language")
	ErrUnSupporttedEditor = errors.New("only vim and neovim/nvim supported")

	codeLangExtensionDic = map[string]string{
		"go":      ".go",
		"golang":  ".go",
		"java":    ".java",
		"python":  ".py",
		"python2": ".py",
		"python3": ".py",
		"cpp":     ".cpp",
		"c++":     ".cpp",
		"c":       ".c",
		// TODO: add other language mappings
	}
	editorCmdDic = map[string]string{
		"neovim": "nvim",
		"nvim":   "nvim",
		"vim":    "vim",
	}
)

type Config struct {
	Language string `json:"language,omitempty"`
	CodeLang string `json:"codeLang,omitempty"`
	Editor   string `json:"editor,omitempty"`
}

var defaultCfg = &Config{
	Language: DefaultLanguage,
	CodeLang: DefaultCodeLang,
	Editor:   DefaultEditor,
}

func Write(cfg *Config) error {
	preCfg, err := Get()
	if err != nil {
		if err == errConfigNotExist {
			preCfg = defaultCfg
		} else {
			return err
		}
	}
	cfg = adapt(preCfg, cfg)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	err = os.WriteFile(configFile, data, 0640)
	log.Trace(err)
	return err
}

func adapt(pre, cur *Config) *Config {
	// struct Config is very simple now, if it becomes complex, use json marshal unmarshal instead
	if cur.CodeLang != "" {
		pre.CodeLang = cur.CodeLang
	}
	if cur.Language != "" {
		pre.Language = cur.Language
	}
	if pre.CodeLang == "golang" {
		pre.CodeLang = "go"
	}
	if cur.Editor != "" {
		pre.Editor = cur.Editor
	}
	return pre
}

func Read() ([]byte, error) {
	_, err := os.Stat(configFile)
	if err != nil {
		log.Trace(err)
		if os.IsNotExist(err) {
			err = errConfigNotExist
		}
		return nil, err
	}

	data, err := os.ReadFile(configFile)
	log.Trace(err)

	return data, err
}

func Get() (*Config, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	res := &Config{}
	err = json.Unmarshal(data, &res)
	log.Trace(err)
	return res, err
}

func Domain() string {
	cfg, err := Get()
	if err != nil {
		return enDomain
	}
	if cfg.Language == cnLanguage {
		return cnDomain
	}
	return enDomain
}

func IsCN(cfg *Config) bool {
	return cfg.Language == cnLanguage
}

func GetCodeFileExt(codeLang string) string {
	return codeLangExtensionDic[codeLang]
}

func IsGolang(cfg *Config) bool {
	return cfg.CodeLang == "go" || cfg.CodeLang == "golang"
}

func SrpportedLang(lang string) bool {
	return lang == "en" || lang == "cn"
}

func SupportedCodeLang(codeLang string) bool {
	_, ok := codeLangExtensionDic[codeLang]
	return ok
}

func SupportedEditor(editor string) bool {
	_, ok := editorCmdDic[editor]
	return ok
}
func GetEditorCmd(editor string) string {
	return editorCmdDic[editor]
}
