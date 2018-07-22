package processor

import (
	"go/ast"
	"sync"

	"golang.org/x/tools/go/loader"

	"github.com/famigo/lang/parser"
)

var processed = new(sync.Map)

var vectors = make(map[string]*ast.FuncDecl, 3)

type processor struct {
	grp *sync.WaitGroup
}

func (proc *processor) process(node ast.Node) {
	proc.grp.Add(1)
	go proc.walk(node)
}

func (proc *processor) walk(node ast.Node)  {
	ast.Walk(proc, node)
	proc.grp.Done()
}

func (proc *processor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	if ident, ok := node.(*ast.Ident); ok {
		if _, loaded := processed.LoadOrStore(ident, true); loaded {
			return proc
		}
		if fun := funcDeclOf(ident); fun != nil && !isVector(fun) {
			proc.process(fun)
		}
	}
	return proc
}

func funcDeclOf(ident *ast.Ident) *ast.FuncDecl {
	if obj := ident.Obj; obj != nil {
		if decl := obj.Decl; decl != nil {
			if fun, ok := decl.(*ast.FuncDecl); ok {
				return fun
			}
		}
	}
	return nil
}

func isVector(fun *ast.FuncDecl) bool {
	name := fun.Name.Name
	return (name == "main" || name == "nmi" || name == "irq") &&
		fun.Recv == nil &&
		fun.Type.Params.NumFields() == 0
}

//ProcessPackage process a FamiGo game package
func ProcessPackage(pkgpath string) *loader.Program {
	prog := parser.ParsePackage(pkgpath)

	for _, file := range prog.Package(pkgpath).Files {
		for _, decl := range file.Decls {
			if fun, ok := decl.(*ast.FuncDecl); ok {
				if isVector(fun) {
					vectors[fun.Name.Name] = fun
					ast.Print(prog.Fset, file)
				}
			}
		}
	}

	grp := new(sync.WaitGroup)
	for _, vector := range vectors {
		proc := &processor{grp}
		proc.process(vector)
	}
	grp.Wait()

	return prog
}

//IsProcessed returns true if the identifier is processed
func IsProcessed(ident *ast.Ident) bool {
	_, proc := processed.Load(ident)
	return proc
}
