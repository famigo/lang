package processor

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
	"sync"

	get "github.com/famigo/lang/getter"
	"golang.org/x/tools/go/loader"
)

//Process constants, headers, roms and vectors
func processTopLevelDecls(prog *loader.Program) {
	get.Init(prog)
	group := new(sync.WaitGroup)
	for _, pkginfo := range prog.AllPackages {
		get.NameOfPkg(pkginfo.Pkg)
		for _, file := range pkginfo.Files {
			group.Add(len(file.Decls))
			for _, decl := range file.Decls {
				go processTopLevelDecl(&decl, pkginfo, group)
			}
		}
	}
	group.Wait()
}

func processTopLevelDecl(decl *ast.Decl, pkginfo *loader.PackageInfo, group *sync.WaitGroup) {
	defer group.Done()

	switch d := (*decl).(type) {
	case *ast.GenDecl:
		processGenDecl(d, pkginfo)
	case *ast.FuncDecl:
		processVectorDecl(d, pkginfo)
	}
}

func processGenDecl(decl *ast.GenDecl, pkginfo *loader.PackageInfo) {
	for _, spec := range decl.Specs {
		if vspec, ok := spec.(*ast.ValueSpec); ok {
			for _, ident := range vspec.Names {
				o := pkginfo.ObjectOf(ident)
				if o == nil {
					continue
				}
				switch obj := o.(type) {
				case *types.Const:
					fmt.Printf("constant %s.%s = %s\n", pkginfo.Pkg.Name(), obj.Name(), obj.Val())
				case *types.Var:
					if decl.Doc != nil {
						doc := decl.Doc.List[len(decl.Doc.List)-1].Text
						if strings.Index(doc, "famigo:rom") != -1 {
							fmt.Printf("rom %s.%s\n", pkginfo.Pkg.Name(), obj.Name())
						} else {
							fmt.Printf("comment %s.%s = %s\n", pkginfo.Pkg.Name(), obj.Name(), doc)
						}
					}

				}
			}
		}
	}
}

func processVectorDecl(decl *ast.FuncDecl, pkginfo *loader.PackageInfo) {

}
