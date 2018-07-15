package processor

import (
	"github.com/famigo/lang/parser"
)



//ProcessPackage process a FamiGo game package
func ProcessPackage(pkgpath string) {
	prog := parser.ParsePackage(pkgpath)

	processTopLevelDecls(prog)
}
