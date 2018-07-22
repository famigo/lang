package symbol

import (
	"go/ast"

	"github.com/famigo/lang/data"
)

//Function represents a method or a function
type Function struct {
	Ident  *ast.Ident
	Name   string
	Rom    *data.ROM
	Rcv    *Variable
	Params []Variable
	Rets   []Variable
	Method bool
	Inline bool
}

//Scope returns the interface scope that this method implements or nil if it's not an implementation
func (f *Function) Scope() *Function {
	return nil
}
