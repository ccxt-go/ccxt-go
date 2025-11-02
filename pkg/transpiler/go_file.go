package transpiler

import (
	"github.com/ccxt-go/ccxt-go/pkg/turing_parser"
	"github.com/emirpasic/gods/lists/arraylist"
)

type GoArgument struct {
	Name string
}

type GoFunction struct {
	Name      string
	Arguments arraylist.List

	Body *turing_parser.Expressions
}

type GoFile struct {
	Package string
	Imports []string

	ThisName     string
	OwnFunctions map[string]bool
}
