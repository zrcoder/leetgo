package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	_ "github.com/zellyn/kooky/browser/all"

	"github.com/zrcoder/leetgo/internal/log"

	home "github.com/mitchellh/go-homedir"
)

const (
	DefaultLanguage = "en"
	DefaultCodeLang = "go"
	cnLanguage      = "cn"
	DefaultEditor   = "neovim"
)

var (
	ErrConfigNotExist     = errors.New("no config found, try `leetgo config`")
	ErrInvalidLang        = errors.New("only cn or en language supported")
	ErrInvalidCodeLang    = errors.New("not supported code language")
	ErrUnSupporttedEditor = errors.New("only vim and neovim/nvim supported")

	CfgDir  = "."
	cfgFile = "leetgo.json"

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
		"emacs":  "emacs",
		"vscode": "code",
	}
	editorCmdOption = map[string][]string{
		"neovim": {"-p"},
		"nvim":   {"-p"},
		"vim":    {"-p"},
		"vscode": nil,
		"emacs":  nil,
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

func init() {
	homedir, _ := home.Dir()
	if homedir != "" {
		CfgDir = filepath.Join(homedir, ".leetgo")
		cfgFile = filepath.Join(CfgDir, cfgFile)
		_ = os.MkdirAll(CfgDir, os.ModePerm)
	}
}

func Write(cfg *Config) (*Config, error) {
	storedCfg, err := Get()
	if err != nil {
		if err == ErrConfigNotExist {
			storedCfg = defaultCfg
		} else {
			return nil, err
		}
	}
	cfg = merge(storedCfg, cfg)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	err = os.WriteFile(cfgFile, data, 0o640)
	log.Debug(err)
	return cfg, err
}

func merge(pre, cur *Config) *Config {
	// struct Config is very simple now, if it becomes complex, use json marshal unmarshal instead
	if cur.CodeLang != "" {
		pre.CodeLang = cur.CodeLang
	}
	if cur.Language != "" {
		pre.Language = cur.Language
	}
	if cur.Editor != "" {
		pre.Editor = cur.Editor
	}
	return pre
}

func read() ([]byte, error) {
	_, err := os.Stat(cfgFile)
	if err != nil {
		log.Debug(err)
		if os.IsNotExist(err) {
			err = ErrConfigNotExist
		}
		return nil, err
	}
	return os.ReadFile(cfgFile)
}

func Get() (*Config, error) {
	data, err := read()
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	res := &Config{}
	err = json.Unmarshal(data, &res)
	return res, err
}

func IsDefaultLang() bool {
	if os.Getenv("LCL") == DefaultLanguage {
		return true
	}

	cfg, err := Get()
	if err != nil {
		return true
	}
	return cfg.Language == DefaultLanguage
}

func GetCodeFileExt(codeLang string) string {
	return codeLangExtensionDic[codeLang]
}

func IsGolang(cfg *Config) bool {
	return cfg.CodeLang == "go" || cfg.CodeLang == "golang"
}

func LeetcodeLang(lang string) string {
	if lang == "go" {
		return "golang"
	}
	return lang
}

func DisplayLang(lang string) string {
	if lang == "golang" {
		return "go"
	}
	return lang
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

func GetEditorCmdOps(editor string) (string, []string) {
	return editorCmdDic[editor], editorCmdOption[editor]
}
