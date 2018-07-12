package processor

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	_ "strings"

	"golang.org/x/tools/go/loader"

	// for later use to parse dependencies
	_ "github.com/KyleBanks/depth"
)

//ProcessPackage process a famigo package
func ProcessPackage(pkgpath string) {
	mode := parser.ParseComments | parser.AllErrors
	conf := loader.Config{ParserMode: mode}
	conf.Import(pkgpath)

	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	dir := fmt.Sprintf("%s/src/%s", build.Default.GOPATH, pkgpath)
	pkgs, err := parser.ParseDir(prog.Fset, dir, nil, mode)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		file := ast.MergePackageFiles(pkg, ast.FilterImportDuplicates)
		for _, decl := range file.Decls {
			ast.Walk(skywalker{prog}, decl)
		}
	}

	ast.Print(prog.Fset, pkgs)
}

type skywalker struct {
	prg *loader.Program
}

func (w skywalker) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch t := n.(type) {
	case *ast.GenDecl:
		if t.Tok != token.VAR && t.Tok != token.CONST {
			return w
		}

		for _, spec := range t.Specs {
			if value, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range value.Names {
					// info, _, _ := w.prg.PathEnclosingInterval(name.Pos(), name.End())
					o := w.objectOf(name)
					if o != nil {
						fmt.Println(o)
					}
				}
			}
		}
	}

	return w
}

func (w *skywalker) objectOf(ident *ast.Ident) types.Object {
	for pname, pkg := range w.prg.AllPackages {
		log.Printf("looking for %s in %s ...", ident.Name, pname)
		if obj := pkg.ObjectOf(ident); obj != nil {
			log.Println("found!")
			return obj
		}
	}

	log.Println("nothing :(")
	return nil
}

//ProcessPkg process a famigo package
func ProcessPkg(pkgpath string) {
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	mode := parser.ParseComments | parser.AllErrors
	fset := token.NewFileSet()

	dir := fmt.Sprintf("%s/src/%s", build.Default.GOPATH, pkgpath)
	pkgs, err := parser.ParseDir(fset, dir, nil, mode)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		file := ast.MergePackageFiles(pkg, ast.FilterImportDuplicates)
		conf := types.Config{Importer: importer.Default()}
		pkg, err := conf.Check(pkgpath, fset, []*ast.File{file}, info)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(pkg.String())
		for ident := range info.Defs {
			fmt.Println(ident)
		}

		for _, decl := range file.Decls {
			w := walker{info}
			ast.Walk(w, decl)
		}
	}

	ast.Print(fset, pkgs)
}

type walker struct {
	info *types.Info
}

func (w walker) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch t := n.(type) {
	case *ast.GenDecl:
		if t.Tok != token.VAR && t.Tok != token.CONST {
			return w
		}

		for _, spec := range t.Specs {
			if value, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range value.Names {
					fmt.Println(w.info.Defs[name])
				}
			}
		}
	}

	return w
}

//ProcessFile process a famigo source code
func ProcessFile(filepath string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath, nil, parser.AllErrors)

	if err != nil {
		panic(err)
	}

	for _, decl := range file.Decls {
		var v visitor

		ast.Walk(v, decl)
	}

	ast.Print(fset, file)
	fmt.Println("done")
}

type visitor int

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch t := n.(type) {
	case *ast.GenDecl:
		if t.Tok != token.VAR && t.Tok != token.CONST {
			return v
		}

		for _, spec := range t.Specs {
			if value, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range value.Names {
					fmt.Println(name.Name)
				}
			}
		}
	}

	// fmt.Printf("%s%T\n", strings.Repeat("\t", int(v)), n)

	return v + 1
}
