package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/zellyn/kooky"
	_ "github.com/zellyn/kooky/browser/all"

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

	TokenKey   = "csrftoken"
	SessionKey = "LEETCODE_SESSION"

	cnTokenFlag   = "cnToken"
	cnSessionFlag = "cnSession"
	enTokenFlag   = "enToken"
	enSessionFlag = "enSession"

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
	data, _ := json.MarshalIndent(cfg, "", "  ")
	err := os.WriteFile(configFile, data, 0640)
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
		log.Dev(err)
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

func GetCredentials() (string, string, error) {
	log.Dev("get credentials from config")
	cfg, err := Get()
	if err != nil {
		return "", "", err
	}

	var domain, token, session string
	if cfg[LangKey] == cnLanguage {
		domain = strings.TrimPrefix(cnDomain, "https://")
		token = cfg[cnTokenFlag]
		session = cfg[cnSessionFlag]
	} else {
		strings.TrimPrefix(enDomain, "https://")
		token = cfg[enTokenFlag]
		session = cfg[enSessionFlag]
	}
	if token != "" && session != "" {
		return token, session, nil
	}

	token, session, err = getCredentialsFromBrowser(domain)
	if err != nil {
		return "", "", err
	}

	if cfg[LangKey] == cnLanguage {
		cfg[cnTokenFlag] = token
		cfg[cnSessionFlag] = session
	} else {
		cfg[enTokenFlag] = token
		cfg[enSessionFlag] = session
	}
	return token, session, Write(cfg)
}

func UpdateCredentials() error {
	cfg, err := Get()
	if err != nil {
		return err
	}
	domain := strings.TrimPrefix(enDomain, "https://")
	if cfg[LangKey] == cnLanguage {
		domain = strings.TrimPrefix(cnDomain, "https://")
	}
	token, session, err := getCredentialsFromBrowser(domain)
	if err != nil {
		return err
	}

	if cfg[LangKey] == DefaultLanguage {
		cfg[enTokenFlag] = token
		cfg[enSessionFlag] = session
	} else {
		cfg[cnTokenFlag] = token
		cfg[cnSessionFlag] = session
	}
	return Write(cfg)
}

func getCredentialsFromBrowser(domain string) (string, string, error) {
	log.Dev("get credentials from browser")
	tokenCookies := kooky.ReadCookies(
		kooky.Valid,
		kooky.DomainContains(domain),
		kooky.Name(TokenKey),
	)
	sessionCookies := kooky.ReadCookies(
		kooky.Valid,
		kooky.DomainContains(domain),
		kooky.Name(SessionKey),
	)
	if len(sessionCookies) == 0 || len(tokenCookies) == 0 {
		err := errors.New("failed to get credentials")
		log.Dev(err)
		return "", "", err
	}
	return tokenCookies[0].Value, sessionCookies[0].Value, nil
}
