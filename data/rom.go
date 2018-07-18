package data

import (
	"fmt"
	"go/ast"
	"regexp"
)

//Kind can be PRG, CHR or DMC
type Kind string

const (
	//PRG kind of ROM
	PRG = "prg"
	//CHR kind of ROM
	CHR = "chr"
	//DMC kind of ROM
	DMC = "dmc"
)

const (
	//AnyBank delegates a ROM to be stored in any bank
	AnyBank = "?"
	//AllBanks delegates a ROM to be stored in all banks
	AllBanks = "*"
)

//ROM is a compiled variable or function commented with the directive famigo:rom
type ROM struct {
	kind  Kind
	bank  string
	inc   string
	Label string
	Code  string
}

//Kind of the ROM: PRG, CHR or DMC
func (rom *ROM) Kind() Kind {
	return rom.kind
}

//Bank to store the ROM
func (rom *ROM) Bank() string {
	return rom.bank
}

//Inc is the path - relative or absolute - of a file to being included
func (rom *ROM) Inc() string {
	return rom.inc
}

func (rom *ROM) String() string {
	return fmt.Sprintf("famigo:%v rom:%s inc:%s\n%s:\n\t%s", rom.kind, rom.bank, rom.inc, rom.Label, rom.Code)
}

var romregexp = regexp.MustCompile(`(?://|\\\*)\s*famigo:(prg|chr|dmc)(?:\s*rom:([0-9]|\?|\*))?(?:\s*inc:(.+))?`)

//VarRomOf parses the directive famigo:rom and return the proper ROM
//
//Returns nil if not commented or if not a variable nor a constant declaration
func VarRomOf(decl *ast.GenDecl, spec *ast.ValueSpec) *ROM {
	if rom := parsePragma(spec.Doc); rom != nil {
		return rom
	}
	if rom := parsePragma(decl.Doc); rom != nil {
		return rom
	}
	return nil
}

//FuncRomOf parses the directive famigo:rom and return the proper ROM
//
//Returns a default ROM with bank setted to AnyBank if not commented
func FuncRomOf(decl *ast.FuncDecl) *ROM {
	if rom := parsePragma(decl.Doc); rom != nil {
		return rom
	}
	return &ROM{
		kind: PRG,
		bank: AnyBank}
}

func parsePragma(comment *ast.CommentGroup) *ROM {
	if comment != nil {
		lastline := comment.List[len(comment.List)-1].Text
		match := romregexp.FindStringSubmatch(lastline)
		if match != nil {
			kindval := Kind(match[1])
			bankval := match[2]
			if bankval == "" {
				bankval = AnyBank
			}
			incval := match[3]
			return &ROM{
				kind: kindval,
				bank: bankval,
				inc:  incval}
		}
	}

	return nil
}
