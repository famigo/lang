package parser

import (
	"go/parser"
	"log"
	"sync"

	"golang.org/x/tools/go/loader"
)

var parsed sync.Map

//ParsePackage builds the program AST
func ParsePackage(pkgpath string) *loader.Program {
	conf := loader.Config{ParserMode: parser.ParseComments | parser.AllErrors}
	conf.Import(pkgpath)
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}
	return prog
}
