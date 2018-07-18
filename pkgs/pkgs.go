package pkgs

import (
	"fmt"
	"go/types"
	"path"

	"golang.org/x/tools/go/loader"
)

var names map[string]string

//Name caches all the program packages names
func Name(prog *loader.Program) {
	size := len(prog.AllPackages)
	names = make(map[string]string, size)
	for _, pkginfo := range prog.AllPackages {
		NameOf(pkginfo.Pkg)
	}
}

//NameOf returns the name of the package
func NameOf(pkg *types.Package) string {
	if name, ok := names[pkg.Path()]; ok {
		return name
	}

	pkgname := Qualify(pkg)
	dir, file := path.Split(pkg.Path())
Check:
	for _, othername := range names {
		if othername == pkgname {
			dir, file = path.Split(dir)
			pkgname = fmt.Sprintf("%s.%s", file, pkgname)
			goto Check
		}
	}

	names[pkg.Path()] = pkgname

	return pkgname
}

//Qualify defines a package name
func Qualify(pkg *types.Package) string {
	return pkg.Name()
}
