package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/zrcoder/leetgo/internal/log"
)

const (
	LangKey     = "lang"
	CodeLangKey = "code"
	ProjectKey  = "project"

	DefaultProjectDir = "leetgo"
	DefaultLanguage   = "en"
	CodeLangGo        = "golang"
	CodeLangGoShort   = "go"
	CodeLangJava      = "java"
	CodeLangPython    = "python"
	DefaultCodeLang   = CodeLangGo
	cnLanguage        = "cn"

	enDomain = "https://leetcode.com"
	cnDomain = "https://leetcode.cn"
)

var (
	configFile = ".leetgo"
)

func init() {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configFile = filepath.Join(configDir, configFile)
}

var defaultCfg = map[string]string{
	LangKey:     DefaultLanguage,
	ProjectKey:  DefaultProjectDir,
	CodeLangKey: DefaultCodeLang,
}

func Write(cfg map[string]string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Dev(err)
		return err
	}

	err = os.WriteFile(configFile, data, 0640)
	log.Dev(err)

	return err
}

func Read() ([]byte, error) {
	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			err = Write(defaultCfg)
		} else {
			log.Dev(err)
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configFile)
	log.Dev(err)

	return data, err
}

func Get() (map[string]string, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	res := map[string]string{}
	err = json.Unmarshal(data, &res)
	log.Dev(err)

	return res, err
}

func Domain() string {
	info, err := Get()
	if err != nil {
		return enDomain
	}
	if info[LangKey] == cnLanguage {
		return cnDomain
	}

	return enDomain
}

var CodeLangExtensionDic = map[string]string{
	CodeLangGo:     ".go",
	CodeLangJava:   ".java",
	CodeLangPython: ".py",
	// TODO: add other language mappings
}
