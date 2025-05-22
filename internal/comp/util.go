package comp

import (
	"fmt"
	"sort"

	// tmodel "github.com/zrcoder/tdoc/model"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/leetgo/utils/exec"
	"github.com/zrcoder/leetgo/utils/huh"
)

const (
	noQuestionFoundForErrFmt = "no question found for `%s`"
)

func queryMeta(frontendID string) (*model.Meta, error) {
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

	return query(frontendID)
}

func queryMetas(key string) ([]model.Meta, error) {
	if key == "today" {
		today, err := remote.GetToday()
		if err != nil {
			return nil, err
		}
		return []model.Meta{*today.Meta()}, nil
	}

	return search(key)
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
	return today.Meta().FrontendID
}

func query(frontendID string) (*model.Meta, error) {
	allMap, err := local.GetAll()
	if err == nil {
		res, ok := allMap[frontendID]
		if ok {
			return &res, nil
		}
		err = local.ErrNotCached
	}
	if err != local.ErrNotCached {
		return nil, err
	}

	all, err := remote.GetAll()
	if err != nil {
		return nil, err
	}
	allMap = make(map[string]model.Meta, len(all.StatStatusPairs))
	for _, sp := range all.StatStatusPairs {
		meta := sp.Meta()
		meta.Transform()
		allMap[meta.FrontendID] = meta
	}
	err = local.WriteAll(allMap)
	if err != nil {
		return nil, err
	}
	if res, ok := allMap[frontendID]; ok {
		return &res, nil
	}
	return nil, fmt.Errorf(noQuestionFoundForErrFmt, frontendID)
}

func search(key string) ([]model.Meta, error) {
	list, err := remote.Search(key)
	if err != nil {
		return nil, err
	}

	res := list.Data.ProblemsetQuestionList.Questions

	if len(res) == 0 {
		log.Debug("no questions found")
		return nil, fmt.Errorf(noQuestionFoundForErrFmt, key)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].FrontendID < res[j].FrontendID
	})
	return res, nil
}

func askToCode(id string) error {
	edit := true
	err := huh.NewConfirm(
		"Solve the question now?",
		"Open the local code file with your favorite editor to edit the code.",
		&edit).Run()
	if err != nil {
		return err
	}

	if edit {
		return code(id)
	}
	return nil
}

func code(id string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	cmd, ops := config.GetEditorCmdOps(cfg.Editor)
	mdFile := local.GetMarkdownFile(id)
	codeFile := local.GetCodeFile(cfg, id)
	ops = append(ops, mdFile, codeFile)
	return exec.Run(local.GetDir(id), cmd, ops...)
}

func openTempMD(id string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	cmd, _ := config.GetEditorCmdOps(cfg.Editor)
	dir := local.GetDir(id)
	log.Debug("open temp markdown", dir, local.TempMDFile)
	return exec.Run(dir, cmd, local.TempMDFile)
}
