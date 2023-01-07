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
	CacheDirKey = "directory"

	DefaultCacheDir = "leetgo"
	DefaultLanguage = "en"
	CodeLangGo      = "golang"
	CodeLangGoShort = "go"
	CodeLangJava    = "java"
	CodeLangPython  = "python"
	DefaultCodeLang = CodeLangGo
	cnLanguage      = "cn"

	TokenKey   = "csrftoken"
	SessionKey = "LEETCODE_SESSION"

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

type Config struct {
	Language string `json:"language"`
	CacheDir string `json:"cacheDir"`
	CodeLang string `json:"codeLang"`
	Token    string `json:"token"`
	Session  string `json:"session"`
}

var defaultCfg = &Config{
	Language: DefaultLanguage,
	CacheDir: DefaultCacheDir,
	CodeLang: DefaultCodeLang,
}

func Write(cfg *Config) error {
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

func Get() (*Config, error) {
	data, err := Read()
	if err != nil {
		return nil, err
	}
	res := &Config{}
	err = json.Unmarshal(data, &res)
	log.Dev(err)

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

	token := cfg.Token
	session := cfg.Token
	if token != "" && session != "" {
		return token, session, nil
	}

	domain := strings.TrimPrefix(enDomain, "https://")
	if cfg.Language == cnLanguage {
		domain = strings.TrimPrefix(cnDomain, "https://")
	}
	token, session, err = getCredentialsFromBrowser(domain)
	if err != nil {
		return "", "", err
	}

	cfg.Token = token
	cfg.Session = session
	return token, session, Write(cfg)
}

func UpdateCredentials() error {
	cfg, err := Get()
	if err != nil {
		return err
	}
	domain := strings.TrimPrefix(enDomain, "https://")
	if cfg.Language == cnLanguage {
		domain = strings.TrimPrefix(cnDomain, "https://")
	}
	token, session, err := getCredentialsFromBrowser(domain)
	if err != nil {
		return err
	}

	cfg.Token = token
	cfg.Session = session
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
