package comp

import (
	"fmt"

	"github.com/zrcoder/rdbook"

	"github.com/zrcoder/leetgo/internal/render"
)

func NewServer(htmlSrc, port string) Component {
	return &Server{htmlSrc: htmlSrc, port: port}
}

type Server struct {
	htmlSrc string
	port    string
}

func (s Server) Run() error {
	fmt.Println(render.Infof("  Serving on http://localhost:%s", s.port))
	return rdbook.Serve(s.htmlSrc, "9999")
}
