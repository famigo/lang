package parser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/types"
	"log"
	"os"
	"strings"
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

//ParsePackageOld parses a package recursively
func ParsePackageOld(pkgpath string) {
	files := make(chan *ast.File)
	infos := make(chan *types.Info)
	jobs := &sync.WaitGroup{}
	jobs.Add(1)

	go func() {
		for {
			select {
			case file, open := <-files:
				if open {
					Files = append(Files, file)

					for _, decl := range file.Decls {
						if fn, ok := isVector(decl); ok {
							Vectors[fn.Name.Name] = &decl
						}
					}
				} else {
					return
				}
			case info, open := <-infos:
				if open {
					Infos = append(Infos, info)
				} else {
					return
				}
			}
		}
	}()

	go parse(pkgpath, files, infos, jobs)

	jobs.Wait()
	close(files)
	close(infos)
}

func parse(pkgpath string, files chan<- *ast.File, infos chan<- *types.Info, jobs *sync.WaitGroup) {
	defer jobs.Done()

	if _, loaded := parsed.LoadOrStore(pkgpath, true); loaded {
		return
	}

	pkgs, err := parsePath(pkgpath)
	if err != nil {
		log.Fatal(err)
	}

	for _, parsed := range pkgs {
		file := ast.MergePackageFiles(parsed, ast.FilterImportDuplicates|ast.FilterUnassociatedComments)
		info := newInfo()
		conf := types.Config{Importer: importer.Default()}
		checked, err := conf.Check(pkgpath, Fset, []*ast.File{file}, info)
		if err != nil {
			log.Fatal(err)
		}

		files <- file
		infos <- info

		jobs.Add(len(checked.Imports()))
		for _, imported := range checked.Imports() {
			go parse(imported.Path(), files, infos, jobs)
		}
	}
}

func parsePath(path string) (map[string]*ast.Package, error) {
	filter := func(f os.FileInfo) bool {
		return !strings.HasSuffix(f.Name(), "_test.go")
	}
	mode := parser.ParseComments | parser.AllErrors
	pkgdir := fmt.Sprintf("%s/src/%s", build.Default.GOPATH, path)
	pkgs, err := parser.ParseDir(Fset, pkgdir, filter, mode)

	if err != nil {
		if _, notfound := err.(*os.PathError); notfound {
			pkgdir := fmt.Sprintf("%s/src/%s", build.Default.GOROOT, path)
			return parser.ParseDir(Fset, pkgdir, filter, mode)
		}
	}

	return pkgs, err
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

func isVector(decl ast.Decl) (*ast.FuncDecl, bool) {
	if fn, ok := decl.(*ast.FuncDecl); ok {
		fname := fn.Name.Name
		return fn, fname == "main" || fname == "nmi" || fname == "irq"
	}

	return nil, false
}
