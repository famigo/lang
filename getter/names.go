package getter

import (
	"golang.org/x/tools/go/loader"
	"fmt"
	"go/types"
	"path"
)

var pkgnames map[string]string

//Init the caches
func Init(prog *loader.Program)  {
	pkgnames = make(map[string]string, len(prog.AllPackages))
}

//NameOfPkg returns the name of the package
func NameOfPkg(pkg *types.Package) string {
	if name, ok := pkgnames[pkg.Path()]; ok {
		return name
	}

	pkgname := pkg.Name()
	dir, file := path.Split(pkg.Path())
Check:
	for _, othername := range pkgnames {
		if othername == pkgname {
			dir, file = path.Split(dir)
			pkgname = fmt.Sprintf("%s.%s", file, pkgname)
			goto Check
		}
	}

	pkgnames[pkg.Path()] = pkgname

	return pkgname
}
