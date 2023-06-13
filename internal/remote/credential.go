package remote

import (
	"errors"
	"strings"

	"github.com/j178/kooky"
	_ "github.com/j178/kooky/browser/chrome"
	_ "github.com/j178/kooky/browser/edge"
	_ "github.com/j178/kooky/browser/firefox"
	_ "github.com/j178/kooky/browser/safari"
	"github.com/zrcoder/leetgo/internal/log"
)

var supporttedBrowsers = []string{
	"chrome",
	"edge",
	"firefox",
	"safari",
}

const (
	tokenKey   = "csrftoken"
	sessionKey = "LEETCODE_SESSION"
)

func getCredentials(domain string) (string, string, error) {
	domain = strings.TrimPrefix(domain, "https://")
	return getCredentialsFromBrowsers(domain)
}

func getCredentialsFromBrowsers(domain string) (string, string, error) {
	log.Debug("search credentials from browsers")
	filters := []kooky.Filter{
		kooky.DomainHasSuffix(domain),
		kooky.FilterFunc(
			func(cookie *kooky.Cookie) bool {
				return kooky.Name(sessionKey).Filter(cookie) ||
					kooky.Name(tokenKey).Filter(cookie)
			},
		),
	}
	for _, store := range kooky.FindCookieStores(supporttedBrowsers...) {
		log.Debug("search:", store.Browser(), "file", store.FilePath())
		cookies, err := store.ReadCookies(filters...)
		if err != nil {
			log.Debug("failed to read cookies:", err)
			continue
		}
		if len(cookies) < 2 {
			log.Debug("no cookie found for", store.Browser())
			continue
		}
		var session, csrfToken string
		for _, cookie := range cookies {
			if cookie.Name == sessionKey {
				session = cookie.Value
			}
			if cookie.Name == tokenKey {
				csrfToken = cookie.Value
			}
		}
		if session == "" || csrfToken == "" {
			log.Debug("no cookie found for", store.Browser(), ":", domain)
			continue
		}
		log.Debug("read credentials succeed for", store.Browser(), ":", domain)
		return csrfToken, session, nil
	}
	return "", "", errors.New("no credentials found from browsers")
}
