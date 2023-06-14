package remote

import (
	"sync"

	"github.com/carlmjohnson/requests"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

type Clienter interface {
	GetList() (*model.List, error)
	GetQuestion(meta *model.StatStatusPair) (*model.Question, error)
	GetToday() (res *model.Today, err error)
	Test(question *model.Question, typedCode, codeLang string) (string, error)
	Submit(question *model.Question, typedCode, codeLang string) (string, error)
	CheckResult(id string, question *model.Question, res model.RunResult) error
	GetSolutions(meta *model.StatStatusPair) (model.SolutionListResp, error)
	GetSolution(req *model.SolutionReq, meta *model.StatStatusPair) (*model.SolutionResp, error)
}

const (
	enDomain = "https://leetcode.com"
	cnDomain = "https://leetcode.cn"

	solutionsLimit = 10
)

func GetList() (*model.List, error) {
	return curCli().GetList()
}

func GetQuestion(sp *model.StatStatusPair) (*model.Question, error) {
	return curCli().GetQuestion(sp)
}

func GetToday() (res *model.Today, err error) {
	return curCli().GetToday()
}

func Test(question *model.Question, typedCode, codeLang string) (string, error) {
	return curCli().Test(question, typedCode, codeLang)
}

func Submit(question *model.Question, typedCode, codeLang string) (string, error) {
	return curCli().Submit(question, typedCode, codeLang)
}

func CheckResult(id string, question *model.Question, res model.RunResult) error {
	return curCli().CheckResult(id, question, res)
}

func GetSolutions(meta *model.StatStatusPair) (model.SolutionListResp, error) {
	return curCli().GetSolutions(meta)
}

func GetSolution(solution *model.SolutionReq, meta *model.StatStatusPair) (*model.SolutionResp, error) {
	return curCli().GetSolution(solution, meta)
}

func curCli() Clienter {
	var (
		cli     Clienter
		cliOnce = sync.Once{}
	)
	cliOnce.Do(func() {
		var token, session string
		var err error
		if config.IsDefaultLang() {
			token, session, err = getCredentials(enDomain)
			cli = newClient(enDomain, token, session)
		} else {
			token, session, err = getCredentials(cnDomain)
			cli = newClientCN(cnDomain, token, session)
		}
		if err != nil {
			log.Debug(err)
		}
	})
	return cli
}

func newClient(domain, token, session string) *client {
	return &client{
		domain: domain,
		rb: requests.New().BaseURL(domain).
			ContentType("application/json").
			Cookie("LEETCODE_SESSION", session).
			Cookie("csrftoken", token).
			Header("x-csrftoken", token),
	}
}

func newClientCN(domain, token, session string) *clientCN {
	return &clientCN{
		client: newClient(domain, token, session),
	}
}