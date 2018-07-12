package processor

import (
	"go/ast"

	"github.com/famigo/lang/parser"
)

//ProcessPackage process a FamiGo game package
func ProcessPackage(pkgpath string) {
	prog := parser.ParsePackage(pkgpath)
	ast.Print(parser.Fset, prog.Vectors["main"])
}
