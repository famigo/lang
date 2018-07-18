package constant

import (
	"fmt"
	"go/types"
	"strconv"

	"github.com/famigo/lang/pkgs"
)

//NonConstantError indicates an atempt to eval a non-constant variable
type NonConstantError struct {
}

func (e *NonConstantError) Error() string {
	return "non-constant value"
}

//NameOf returns the name of a constant
func NameOf(cons *types.Const) string {
	if cons.Name() == "_" {
		return ""
	}
	return fmt.Sprintf("%s.%s", pkgs.NameOf(cons.Pkg()), cons.Name())
}

//ValueOf returns the value of a constant
//
//Returns a NonConstantError if the value is a string, float or complex
func ValueOf(obj *types.Const) (string, error) {
	if basic, ok := obj.Type().Underlying().(*types.Basic); ok {
		if basic.Info()&types.IsBoolean == types.IsBoolean {
			istrue, _ := strconv.ParseBool(obj.Val().String())
			if istrue {
				return "1", nil
			}
			return "0", nil
		}
		if basic.Info()&types.IsInteger == types.IsInteger {
			return obj.Val().ExactString(), nil
		}
	}

	return "", new(NonConstantError)
}
