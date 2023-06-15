package remote

import (
	"sync"

	"github.com/carlmjohnson/requests"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

type Clienter interface {
	GetAll() (*model.All, error)
	Search(keyWords string) (*model.QustionsResp, error)
	GetQuestion(meta *model.Meta) (*model.Question, error)
	GetToday() (res *model.Today, err error)
	Test(question *model.Question, typedCode, codeLang string) (string, error)
	Submit(question *model.Question, typedCode, codeLang string) (string, error)
	CheckResult(id string, question *model.Question, res model.RunResult) error
	GetSolutions(meta *model.Meta) (model.SolutionListResp, error)
	GetSolution(req *model.SolutionReq, meta *model.Meta) (*model.SolutionResp, error)
}

const (
	enDomain = "https://leetcode.com"
	cnDomain = "https://leetcode.cn"

	listLimit = 10
)

func GetAll() (*model.All, error) {
	return curCli().GetAll()
}

func Search(keyWords string) (*model.QustionsResp, error) {
	return curCli().Search(keyWords)
}

func GetQuestion(sp *model.Meta) (*model.Question, error) {
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

func GetSolutions(meta *model.Meta) (model.SolutionListResp, error) {
	return curCli().GetSolutions(meta)
}

func GetSolution(solution *model.SolutionReq, meta *model.Meta) (*model.SolutionResp, error) {
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
