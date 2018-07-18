package processor

import (
	"fmt"
	"go/ast"
	"sync"

	"github.com/famigo/lang/compiler"

	 "github.com/famigo/lang/pkgs"
	"golang.org/x/tools/go/loader"
)

var vectors = make(map[string]*ast.FuncDecl, 3)
var mutex = new(sync.Mutex)

//Process all constants, registers, headers, roms and vectors
func processTopLevelDecls(prog *loader.Program) {
	pkgs.Name(prog)
	group := new(sync.WaitGroup)
	for _, pkginfo := range prog.AllPackages {
		for _, file := range pkginfo.Files {
			group.Add(len(file.Decls))
			for _, decl := range file.Decls {
				go processTopLevelDecl(decl, pkginfo, group)
			}
		}
	}
	group.Wait()
	fmt.Println(vectors)
}

func processTopLevelDecl(decl ast.Decl, pkginfo *loader.PackageInfo, group *sync.WaitGroup) {
	defer group.Done()

	switch d := decl.(type) {
	case *ast.GenDecl:
		compiler.CompileGenDecl(d, pkginfo)
	case *ast.FuncDecl:
		processVectorDecl(d, pkginfo)
	}
}

func processVectorDecl(decl *ast.FuncDecl, pkginfo *loader.PackageInfo) {
	name := decl.Name.Name
	if name == "main" || name == "nmi" || name == "irq" {
		mutex.Lock()
		vectors[name] = decl
		mutex.Unlock()
	}
}
