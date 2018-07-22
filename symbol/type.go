package symbol

//Go baisc types
var (
	Int    = Type{kind: Basic, name: "int", size: 1}
	Int16  = Type{kind: Basic, name: "int16", size: 2}
	UInt16 = Type{kind: Basic, name: "uint16", size: 2}
	Int8   = Type{kind: Basic, name: "int8", size: 1}
	UInt8  = Type{kind: Basic, name: "uint8", size: 1}
	Byte   = Type{kind: Basic, name: "byte", size: 1, underlying: &UInt8}
	Bool   = Type{kind: Basic, name: "bool", size: 1}
)

//Kind is the kind of a type: basic, struct or interface
type Kind byte

//Kind of a type
const (
	Basic = Kind(iota)
	Struct
	Interface
)

//Type represents a declared type or a go type such int or byte
type Type struct {
	kind       Kind
	name       string
	size       int
	underlying *Type
	Scopes     []Function //method scopes of an interface type
	Fields     []Variable //fields of a struct type
}

//Name returns the type name
func (t *Type) Name() string {
	return t.name
}

//Size returns the type size in bytes
func (t *Type) Size() int {
	return t.size
}

//Underlying returns the underlying type or the type itself if it has no uderlying
func (t *Type) Underlying() *Type {
	if t.underlying != nil {
		return t.underlying.Underlying()
	}
	return t
}

//Implementations returns all methods that implement this interface
func (t *Type) Implementations() []Function {
	return nil
}
