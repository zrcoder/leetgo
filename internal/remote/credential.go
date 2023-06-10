package remote

import (
	"errors"
	"strings"

	"github.com/j178/kooky"
	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
)

const (
	TokenKey   = "csrftoken"
	SessionKey = "LEETCODE_SESSION"
)

func GetCredentials() (string, string, error) {
	domain := strings.TrimPrefix(config.Domain(), "https://")
	return getCredentialsFromBrowser(domain)
}

func getCredentialsFromBrowser(domain string) (string, string, error) {
	log.Debug("get credentials from browser")
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
		log.Debug(err)
		return "", "", err
	}
	return tokenCookies[0].Value, sessionCookies[0].Value, nil
}
