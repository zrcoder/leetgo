package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewBookGenerator(b *model.BookMeta) Component {
	return &Book{
		BookMeta: b,
		spinner:  newSpinner("Generating"),
	}
}

type Book struct {
	*model.BookMeta
	spinner *spinner.Spinner
}

func (b *Book) Run() error {
	b.spinner.Start()
	sitePath, markDowPath, err := local.Generate(b.BookMeta)
	b.spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Println(render.Infof("  Your book is generated in %s\n", sitePath))
	if markDowPath != "" {
		fmt.Println(render.Infof("  The markdown docs are in %s\n", markDowPath))
	}
	return nil
}
