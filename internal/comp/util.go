package comp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	tmodel "github.com/zrcoder/tdoc/model"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
)

func queryMeta(frontendID string) (*model.StatStatusPair, error) {
	spinner := newSpinner("inquiring")
	spinner.Start()
	defer spinner.Stop()

	if frontendID == "today" {
		today, err := remote.GetToday()
		if err != nil {
			return nil, err
		}
		return today.Meta(), nil
	}

	list, err := remote.GetList()
	if err != nil {
		return nil, err
	}
	for _, sp := range list.StatStatusPairs {
		sp.Stat.FrontendID = sp.Stat.CalFrontendID()
		if sp.Stat.FrontendID != frontendID {
			continue
		}
		if sp.PaidOnly {
			err := fmt.Errorf("[%s. %s] is locked", sp.Stat.FrontendID, sp.Stat.QuestionTitle)
			log.Debug(err)
			return nil, err
		}
		return &sp, nil
	}
	return nil, fmt.Errorf("no questions found for `%s`", frontendID)
}

func queryMetas(key string) ([]model.StatStatusPair, error) {
	if key == "today" {
		today, err := remote.GetToday()
		if err != nil {
			return nil, err
		}
		return []model.StatStatusPair{*today.Meta()}, nil
	}

	return search(key)
}

func search(key string) ([]model.StatStatusPair, error) {
	list, err := remote.GetList()
	if err != nil {
		return nil, err
	}

	lower := strings.ToLower(key)
	var res []model.StatStatusPair
	for _, sp := range list.StatStatusPairs {
		sp.Stat.FrontendID = sp.Stat.CalFrontendID()
		oriLower := strings.ToLower(sp.Stat.QuestionTitle)
		for _, sep := range []string{" ", ". ", "."} {
			title := sp.Stat.FrontendID + sep + oriLower
			if strings.Contains(title, lower) {
				res = append(res, sp)
				break
			}
		}
	}
	if len(res) == 0 {
		log.Debug("no questions found")
		return nil, fmt.Errorf("no questions found for `%s`", key)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Stat.FrontendID < res[j].Stat.FrontendID
	})
	return res, nil
}

func regualarID(id string) string {
	if id != "today" {
		return id
	}

	spinner := newSpinner("Fetching today")
	spinner.Start()
	today, err := remote.GetToday()
	spinner.Stop()
	if err != nil {
		log.Debug(err)
		return id
	}
	return today.Meta().Stat.FrontendID
}

func getDocsFromLocal() ([]*tmodel.DocInfo, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	ids, err := local.GetPickedQuestionIds(cfg)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("you haven't pick any question yet")
	}
	log.Debug("local ids for book:", ids)

	docs := make([]*tmodel.DocInfo, len(ids))

	for i, id := range ids {
		id := id // for doc.Getter below, captrue current id
		doc := &tmodel.DocInfo{}
		question, err := local.GetQuestion(cfg, id)
		if err != nil {
			return nil, err
		}
		doc.Title = fmt.Sprintf("%s. %s", id, question.Title)
		mdFile := local.GetMarkdownFile(cfg, id)
		fi, err := os.Stat(mdFile)
		if err != nil {
			return nil, err
		}
		doc.ModTime = fi.ModTime()
		log.Debug(doc.Title, doc.ModTime)
		doc.Getter = func(filename string) ([]byte, error) {
			mdData, err := os.ReadFile(mdFile)
			if err != nil {
				return nil, err
			}
			codeData, err := local.GetTypedCode(cfg, id)
			if err != nil {
				return nil, err
			}
			noteDate, err := local.GetNotes(cfg, id)
			if err != nil {
				return nil, err
			}
			noteDate = bytes.TrimSpace(noteDate)
			mdData = bytes.TrimSpace(mdData)

			buf := bytes.NewBuffer(nil)
			buf.WriteString("\n\n## My Solution:\n\n")
			if len(noteDate) > 0 {
				buf.Write(noteDate)
				buf.WriteString("\n\n")
			}
			codeLang := config.DisplayLang(cfg.CodeLang)
			buf.WriteString(fmt.Sprintf("```%s\n", codeLang))
			buf.Write(codeData)
			buf.WriteString("\n```\n")

			mdData = append(mdData, buf.Bytes()...)
			return mdData, nil
		}
		docs[i] = doc
	}

	return docs, nil
}

func getDocsFromSolutions(solutionsResp model.SolutionListResp, meta *model.StatStatusPair) ([]*tmodel.DocInfo, error) {
	reqs := solutionsResp.SolutionReqs()
	docs := make([]*tmodel.DocInfo, len(reqs))
	for i, req := range reqs {
		req := req // for doc.Getter
		doc := &tmodel.DocInfo{}
		doc.Title = req.Title
		doc.Getter = func(s string) ([]byte, error) {
			rsp, err := remote.GetSolution(&req, meta)
			if err != nil {
				return nil, err
			}
			content, err := rsp.RegularContent()
			if err != nil {
				return nil, err
			}
			return append([]byte(content), '\n'), nil
		}
		docs[i] = doc
	}
	return docs, nil
}
