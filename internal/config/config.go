package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/zrcoder/leetgo/internal/trace"
)

const (
	LangKey    = "lang"
	ProjectKey = "project"

	DefaultProjectDir = "~/leetcode"
	DefaultLanguage   = "en"
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
	LangKey:    DefaultLanguage,
	ProjectKey: DefaultProjectDir,
}

func Write(cfg map[string]string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return trace.Wrap(err)
	}

	file, err := os.OpenFile(configFile, os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		return trace.Wrap(err)
	}
	defer file.Close()

	_, err = file.Write(data)

	return trace.Wrap(err)
}

func Read() ([]byte, error) {
	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			err = Write(defaultCfg)
		} else {
			return nil, trace.Wrap(err)
		}
	}
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return data, nil
}

func Get() (map[string]string, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	res := map[string]string{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return res, nil
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
