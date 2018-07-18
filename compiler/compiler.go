package compiler

import (
	"log"
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"go/types"
	"os"
	"path"
	"strconv"

	"github.com/famigo/lang/constant"

	"github.com/famigo/lang/header"

	"github.com/famigo/lang/data"

	"golang.org/x/tools/go/loader"
)

var (
	consts  []string
	prgroms []*data.ROM
	chrroms []*data.ROM
	dmcroms []*data.ROM
	headers = make(map[string]int16, 11)
)

var (
	constchan  = make(chan string)
	romchan    = make(chan *data.ROM)
	headerchan = make(chan *header.Header)
)

var fset *token.FileSet

func init() {
	go func() {
		for {
			select {
			case cons := <-constchan:
				consts = append(consts, cons)
			case rom := <-romchan:
				switch rom.Kind() {
				case data.PRG:
					prgroms = append(prgroms, rom)
				case data.CHR:
					chrroms = append(chrroms, rom)
				case data.DMC:
					dmcroms = append(dmcroms, rom)
				}
			case h := <-headerchan:
				headers[h.Name()] = h.Value
			}
		}
	}()
}

//SetFset sets the file set
func SetFset(fs *token.FileSet) {
	fset = fs
}

//CompileGenDecl compiles all constants, registers, headers and roms
func CompileGenDecl(decl *ast.GenDecl, pkginfo *loader.PackageInfo) {
	for _, spec := range decl.Specs {
		if vspec, ok := spec.(*ast.ValueSpec); ok {
			for _, ident := range vspec.Names {
				o := pkginfo.ObjectOf(ident)
				if o == nil {
					continue
				}

				switch obj := o.(type) {
				case *types.Const:
					val, err := constant.ValueOf(obj)
					if err != nil {
						continue
					}

					if rom := data.VarRomOf(decl, vspec); rom != nil {
						if rom.Inc() != "" {
							inc := rom.Inc()
							if !path.IsAbs(rom.Inc()) {
								inc = path.Join(build.Default.GOPATH, "src", pkginfo.Pkg.Path(), inc)
							}
							file, err := os.Open(inc)
							if err != nil {
								doc := vspec.Doc
								if doc == nil {
									doc = decl.Doc
								}
								fmt.Fprintf(os.Stderr, "#%s\n", pkginfo.Pkg.Path())
								fmt.Fprintf(os.Stderr, "%v: %v\n", fset.Position(doc.Pos()), err)
								os.Exit(1)
							}
							file.Close()
							rom.Code = fmt.Sprintf(`.include "%s"`, inc)
						} else {
							rom.Code = val
						}
						rom.Label = constant.NameOf(obj)
						romchan <- rom
						continue
					}

					if h := header.Of(obj.Type()); h != nil {
						headervalue, _ := strconv.ParseUint(val, 10, 16)
						h.Value = int16(headervalue)
						headerchan <- h
					}

					if nam := constant.NameOf(obj); nam != "" {
						constchan <- fmt.Sprintf("%s = %s", nam, val)
					}
				case *types.Var:
					if rom := data.VarRomOf(decl, vspec); rom != nil {

						if a, ok := obj.Type().(*types.Array); ok {
							rom.Code = strconv.Itoa(int(a.Len()))
							romchan <- rom
						}
					}
				}
			}
		}
	}
}

func valueOf(decl *ast.GenDecl, pkginfo *loader.PackageInfo) (interface{}, error) {

	return nil, nil
}

//Preview compiled symbols
func Preview() {
	fmt.Println("---------")
	fmt.Println("constants")
	fmt.Println("---------")
	for _, cons := range consts {
		fmt.Println(cons)
	}

	fmt.Println("----")
	fmt.Println("roms")
	fmt.Println("----")
	for _, rom := range prgroms {
		fmt.Println(rom)
	}
	for _, rom := range chrroms {
		fmt.Println(rom)
	}
	for _, rom := range dmcroms {
		fmt.Println(rom)
	}

	fmt.Println("------")
	fmt.Println("header")
	fmt.Println("------")
	for hnam, hval := range headers {
		fmt.Printf(".%s %d\n", hnam, hval)
	}
}
