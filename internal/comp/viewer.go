package comp

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/tdoc"
	tmodel "github.com/zrcoder/tdoc/model"

	"github.com/zrcoder/leetgo/internal/book"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/leetgo/utils/render"
)

func NewViewer(id string, solution bool) Component {
	return &viewer{id: id, solution: solution, spinner: newSpinner("Picking")}
}

type viewer struct {
	id       string
	solution bool

	spinner *spinner.Spinner
}

func (v *viewer) Run() error {
	isToday := v.id == "today"
	var today *model.Today
	var err error
	if isToday {
		today, err = remote.GetToday()
		if err != nil {
			return err
		}
		v.id = today.Question().FrontendID
	}

	if v.solution {
		return v.viewSolution(isToday, today)
	}

	return v.viewQuestion(isToday, today)
}

func (v *viewer) viewSolution(isToday bool, today *model.Today) error {
	v.spinner.Start()
	docs, err := v.getDocs(isToday, today)
	v.spinner.Stop()
	if err != nil {
		return err
	}

	return tdoc.Run(docs, tmodel.Config{Title: fmt.Sprintf("Most Voted Solutions for %s", v.id)})
}

func (v *viewer) getDocs(isToday bool, today *model.Today) ([]*tmodel.DocInfo, error) {
	var question model.Question
	var err error
	if isToday {
		question = today.Question()
	} else {
		q, err := query(v.id)
		if err != nil {
			return nil, err
		}
		question = *q
	}
	solutionsResp, err := remote.GetSolutions(&question)
	if err != nil {
		return nil, err
	}
	return book.GetMetaListFromSolutions(solutionsResp, &question)
}

func (v *viewer) viewQuestion(isToday bool, today *model.Today) error {
	printHint := func() {
		typeHint := fmt.Sprintf("Type `leetgo code %s` to solve it", v.id)
		fmt.Println(render.MarkDown(typeHint))
	}

	if local.Exist(v.id) {
		content, err := local.GetMarkdown(v.id)
		if err != nil {
			return err
		}
		fmt.Print(render.MarkDown(string(content)))
		printHint()
		return nil
	}

	var question *model.Question
	var err error
	v.spinner.Start()
	if isToday {
		question, err = remote.GetQuestion(today.Question().StatePair()) // faster than query(c.id)
	} else {
		question, err = query(v.id)
	}
	v.spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Print(render.MarkDown(question.MdContent))

	err = local.Write(question)
	if err != nil {
		return err
	}
	printHint()
	return nil
}

func query(frontendID string) (*model.Question, error) {
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
		return remote.GetQuestion(&sp)
	}
	return nil, fmt.Errorf("no questions found for `%s`", frontendID)
}
