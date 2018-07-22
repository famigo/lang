package compiler

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/types"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/famigo/lang/processor"

	"github.com/famigo/lang/constant"
	"github.com/famigo/lang/data"
	"github.com/famigo/lang/header"
	"golang.org/x/tools/go/loader"
)

type gendeclCmplr struct {
	grp     *sync.WaitGroup
	prog    *loader.Program
	pkginfo *loader.PackageInfo
	decl    *ast.GenDecl
	vspec   *ast.ValueSpec
	ident   *ast.Ident
	value   ast.Expr
	rom     *data.ROM
}

func newGenDeclCmplr(grp *sync.WaitGroup, prog *loader.Program, pkginfo *loader.PackageInfo) *gendeclCmplr {
	return &gendeclCmplr{
		grp:     grp,
		prog:    prog,
		pkginfo: pkginfo}
}

func (c *gendeclCmplr) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	switch node := node.(type) {
	case *ast.GenDecl:
		c.decl = node
	case *ast.CommentGroup:
		if rom := data.ParseRomPragma(node); rom != nil {
			c.rom = rom
		}
	case *ast.ValueSpec:
		c.vspec = node
	case *ast.Ident:
		c.ident = node

		obj := c.pkginfo.ObjectOf(node)
		if obj == nil || obj.Pkg() == nil {
			break
		}

		switch obj := obj.(type) {
		case *types.Const:
			c.compileConst(obj)
		case *types.Var:
			c.compileVar(obj)
		}
	}
	return c
}

func (c *gendeclCmplr) compile(decl ast.Decl) {
	c.grp.Add(1)
	go func() {
		ast.Walk(c, decl)
		c.grp.Done()
	}()
}

func (c *gendeclCmplr) compileConst(cons *types.Const) {
	if c.skipConst() {
		return
	}

	val, err := constant.ValueOf(cons)
	if err != nil {
		//TODO if rom != nil and cons.val is string then assign val as byte array else compile as an non-const variable
		return
	}

	nam := constant.NameOf(cons)

	if c.rom != nil {
		c.compileRom(nam, val)
		romchan <- c.rom
		return
	}

	if h := header.Of(cons.Type()); h != nil {
		headervalue, _ := strconv.ParseUint(val, 10, 16)
		h.Value = int16(headervalue)
		headerchan <- h
	}

	if nam != "" {
		constchan <- fmt.Sprintf("%s = %s", nam, val)
	}
}

func (c *gendeclCmplr) compileVar(v *types.Var) {
	if c.skipVar() {
		return
	}
	if c.rom != nil {
		if a, ok := v.Type().(*types.Array); ok {
			c.rom.Code = strconv.Itoa(int(a.Len()))
			romchan <- c.rom
		}
	}
}

func (c *gendeclCmplr) compileRom(nam string, val string) {
	if c.rom.Inc() != "" {
		inc := c.rom.Inc()
		if !path.IsAbs(inc) {
			inc = path.Join(build.Default.GOPATH, "src", c.pkginfo.Pkg.Path(), inc)
		}
		file, err := os.Open(inc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "#%s\n", c.pkginfo.Pkg.Path())
			fmt.Fprintf(os.Stderr, "%v: %v\n", c.prog.Fset.Position(c.rom.Doc().Pos()), err)
			os.Exit(1)
		}
		file.Close()
		ext := path.Ext(inc)
		if ext == "asm" || ext == "inc" || ext == "s" {
			c.rom.Code = fmt.Sprintf(`  .include "%s"`, inc)
		} else {
			c.rom.Code = fmt.Sprintf(`  .incbin "%s"`, inc)
		}
	} else {
		c.rom.Code = val
	}
	c.rom.Label = nam
}

func (c *gendeclCmplr) skipConst() bool {
	if _, load := compiled.LoadOrStore(c.ident, true); load {
		return true
	}
	return false
}

func (c *gendeclCmplr) skipVar() bool {
	if _, load := compiled.LoadOrStore(c.ident, true); load {
		return true
	}
	if !processor.IsProcessed(c.ident) {
		return c.rom == nil
	}
	return false
}
