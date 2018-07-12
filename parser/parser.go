package parser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"sync"
)

//Package contains the parsed files and the symbols type from a package
type Package struct {
	File *ast.File
	Info *types.Info
}

//Program is a parsed package
type Program struct {
	Pkgs    map[string]*Package
	Vectors map[string]*ast.Decl
}

//Mode is the default parse mode with comments and reporting all errors
const Mode = parser.ParseComments | parser.AllErrors

const unknownPkgName = "?"

//Fset is the AST file set
var Fset = token.NewFileSet()

var mutex = sync.Mutex{}

//ParsePackage parses a package recursively
func ParsePackage(pkgpath string) *Program {
	prog := &Program{
		Pkgs:    make(map[string]*Package),
		Vectors: make(map[string]*ast.Decl, 3),
	}
	wait := make(chan bool, 1)
	parse(types.NewPackage(pkgpath, unknownPkgName), prog, wait)
	<-wait

	return prog
}

func parse(pkg *types.Package, prog *Program, done chan bool) {
	defer func() {
		done <- true
	}()

	if prog.reject(pkg) {
		return
	}

	pkgdir := fmt.Sprintf("%s/src/%s", build.Default.GOPATH, pkg.Path())
	pkgs, err := parser.ParseDir(Fset, pkgdir, nil, Mode)
	if err != nil {
		log.Fatal(err)
	}

	for _, parsed := range pkgs {
		file := ast.MergePackageFiles(parsed, ast.FilterImportDuplicates|ast.FilterUnassociatedComments)
		info := newInfo()
		conf := types.Config{Importer: importer.Default()}
		checked, err := conf.Check(pkg.Path(), Fset, []*ast.File{file}, info)
		if err != nil {
			log.Fatal(err)
		}

		pack, contains := prog.Pkgs[checked.Name()]
		if !contains {
			pack = new(Package)
			prog.Pkgs[checked.Name()] = pack
		}
		pack.File = file
		pack.Info = info

		wait := make(chan bool, len(checked.Imports()))
		for _, imported := range checked.Imports() {
			go parse(imported, prog, wait)
		}

		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				fname := fn.Name.Name
				if fname == "main" || fname == "nmi" || fname == "irq" {
					prog.Vectors[fname] = &decl
					break
				}
			}
		}

		for i := 0; i < len(checked.Imports()); i++ {
			<-wait
		}
	}
}

func newInfo() *types.Info {
	return &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
}

func (prog *Program) reject(pkg *types.Package) bool {
	if pkg.Name() == unknownPkgName {
		return false
	}

	mutex.Lock()
	_, contains := prog.Pkgs[pkg.Name()]
	if !contains {
		prog.Pkgs[pkg.Name()] = new(Package)
	}
	mutex.Unlock()

	return contains
}
