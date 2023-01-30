package comp

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/mdoc"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewBook(sortBy, port string) Component {
	return &Book{
		sortBy:  sortBy,
		port:    port,
		spinner: newSpinner("Generating"),
	}
}

type Book struct {
	sortBy string
	port   string

	spinner *spinner.Spinner
}

func (b *Book) Run() error {
	b.spinner.Start()
	docPath, err := local.Generate(b.sortBy)
	b.spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Println(render.Infof("Your book is generated in %s\n", docPath))

	cfg := mdoc.GetConfig()
	cfg.HttpPort = b.port
	cfg.DocsDirectory = docPath

	fmt.Printf("Serving on http://localhost:%s\n", b.port)
	return mdoc.Serve(cfg)
}
