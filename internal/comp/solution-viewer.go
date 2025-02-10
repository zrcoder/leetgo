package comp

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/huh"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
)

func NewSolutionViewer(id string) Component {
	return &solutionViewer{id: id, spinner: newSpinner("Searching")}
}

type solutionViewer struct {
	spinner *spinner.Spinner
	id      string
}

func (s *solutionViewer) Run() error {
	meta, err := queryMeta(s.id)
	if err != nil {
		return err
	}

	s.spinner.Start()
	solutionsResp, err := remote.GetSolutions(meta)
	s.spinner.Stop()
	if err != nil {
		return err
	}

	reqs := solutionsResp.SolutionReqs()
	options := make([]huh.Option[model.SolutionReq], len(reqs))
	for i, req := range reqs {
		options[i] = huh.NewOption(fmt.Sprintf("%d. %s -- %s", i, req.Title, req.Author), req)
	}
	var req model.SolutionReq
	err = huh.NewForm(huh.NewGroup(huh.NewSelect[model.SolutionReq]().Title("Solutions").Options(options...).Value(&req))).Run()
	if err != nil {
		return err
	}

	rsp, err := remote.GetSolution(&req, meta)
	if err != nil {
		return err
	}

	content, err := rsp.RegularContent()
	if err != nil {
		return err
	}

	err = local.WriteTempMd(meta.FrontendID, []byte(content))
	if err != nil {
		return err
	}

	return openTempMD(meta.FrontendID)
}
