package symbol

import (
	"go/ast"

	"github.com/famigo/lang/data"
)

//Variable has info about package and local variables, also function parameters and returns
type Variable struct {
	Ident *ast.Ident
	Pkg   *Package
	Name  string
	Value interface{} //literal value of a constant, rom or header
	Rom   *data.ROM
	Typ   *Type
	Dim   []int16
	Ptr   bool
}
