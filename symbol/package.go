package symbol

import (
	"go/ast"
)

//Package represents a package
type Package struct {
	Name   string
	Path   string
	Source *ast.Package
}
