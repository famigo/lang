package compiler

import (
	"fmt"
	"sync"

	"github.com/famigo/lang/pkgs"

	"github.com/famigo/lang/processor"

	"github.com/famigo/lang/header"

	"github.com/famigo/lang/data"
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

var compiled = new(sync.Map)

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

//CompilePackage compiles a FamiGo package
func CompilePackage(pkgpath string) {
	prog := processor.ProcessPackage(pkgpath)

	pkgs.Name(prog)

	grp := new(sync.WaitGroup)
	for _, pkginfo := range prog.AllPackages {
		for _, file := range pkginfo.Files {
			for _, decl := range file.Decls {
				cmplr := newGenDeclCmplr(grp, prog, pkginfo)
				cmplr.compile(decl)
			}
		}
	}
	grp.Wait()
	Preview()
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
