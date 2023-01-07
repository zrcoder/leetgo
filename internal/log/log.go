package log

import (
	"github.com/dgraph-io/badger/v3"
)

// release version do nothing

var Dev = func(x ...any) {}

var Logger badger.Logger = logger{}

type logger struct{}

func (l logger) Errorf(s string, i ...any) {}

func (l logger) Warningf(s string, i ...any) {}

func (l logger) Infof(s string, i ...any) {}

func (l logger) Debugf(s string, i ...any) {}
