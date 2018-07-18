package compiler

import (
	"fmt"
	"go/ast"
	"go/types"
	"strconv"

	"github.com/famigo/lang/header"

	"github.com/famigo/lang/pkgs"

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
					var val string
					if basic, ok := obj.Type().Underlying().(*types.Basic); ok {
						if basic.Info() == types.IsString {
							continue
						}
						if basic.Info() == types.IsBoolean {
							istrue, _ := strconv.ParseBool(obj.Val().String())
							if istrue {
								val = "1"
							} else {
								val = "0"
							}
						} else {
							val = obj.Val().ExactString()
						}
					}

					if rom := data.VarRomOf(decl); rom != nil {
						rom.Label = nameOf(obj)
						rom.Code = val
						romchan <- rom
						continue
					}

					if h := header.Of(obj.Type()); h != nil {
						headervalue, _ := strconv.ParseUint(val, 10, 16)
						h.Value = int16(headervalue)
						headerchan <- h
					}

					if nam := nameOf(obj); nam != "" {
						constchan <- fmt.Sprintf("%s = %s", nam, val)
					}
				case *types.Var:
				}
			}
		}
	}
}

func nameOf(cons *types.Const) string {
	if cons.Name() == "_" {
		return ""
	}
	return fmt.Sprintf("%s.%s", pkgs.NameOf(cons.Pkg()), cons.Name())
}
