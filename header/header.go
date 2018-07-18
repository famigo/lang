package header

import (
	"go/types"
	"reflect"
	"strings"

	"github.com/famigo/header/ines"
	"github.com/famigo/header/nes2"
	"github.com/famigo/lang/pkgs"
)

var directives map[string]string

//Header of a ROM
type Header struct {
	name  string
	Value int16
}

//Name of the header
func (header *Header) Name() string {
	return header.name
}

func init() {
	var names = [...]string{
		reflect.TypeOf(ines.DefaultCHR).String(),
		reflect.TypeOf(ines.DefaultMAP).String(),
		reflect.TypeOf(ines.DefaultMIR).String(),
		reflect.TypeOf(ines.DefaultPRG).String(),
		reflect.TypeOf(nes2.DefaultBRAM).String(),
		reflect.TypeOf(nes2.DefaultCHRBRAM).String(),
		reflect.TypeOf(nes2.DefaultCHRRAM).String(),
		reflect.TypeOf(nes2.DefaultPRGRAM).String(),
		reflect.TypeOf(nes2.DefaultSUB).String(),
		reflect.TypeOf(nes2.DefaultTV).String(),
		reflect.TypeOf(nes2.DefaultVS).String()}

	directives = make(map[string]string, len(names))

	for _, name := range names {
		directive := strings.Replace(name, ".", "", 1)
		directive = strings.ToLower(directive)
		directives[name] = directive
	}
}

//Of return the proper header of an iNes or NES2.0 type
//
//Returns nil if the type is not a header type
func Of(typ types.Type) *Header {
	str := types.TypeString(typ, pkgs.Qualify)
	if directive, ok := directives[str]; ok {
		return &Header{name: directive}
	}
	return nil
}
