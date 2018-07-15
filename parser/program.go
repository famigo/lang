package parser

import (
	"go/ast"
	"go/token"
	"go/types"
)

//Vectors maps each vector to its declaration node
var Vectors = make(map[string]*ast.Decl, 3)

//Files from each package
var Files []*ast.File

//Infos are the types iformation from each package
var Infos []*types.Info

//Fset is the AST file set
var Fset = token.NewFileSet()
