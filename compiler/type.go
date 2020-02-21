package compiler

import "strings"

//Type is an 'i' type.
type Type string

//With returns the type with a given subtype.
func (T Type) With(subtype Type) Type {
	return T + "." + subtype
}

//Is returns true if the type is a member of the specified collection type.
func (T Type) Is(collection Type) bool {
	if len(T) < len(collection) {
		return false
	}
	return T[:len(collection)] == collection
}

//Subtype returns the subtype of the type.
func (T Type) Subtype() Type {
	return Type(strings.SplitN(string(T), ".", 2)[1])
}

//Types is all of the types in i.
var Types = []Type{
	"integer", "string", "symbol", "metatype", "bit",
}

//Types in 'i'
const (
	Integer  Type = "integer"
	String   Type = "string"
	Symbol   Type = "symbol"
	Metatype Type = "metatype"
	Bit      Type = "bit"
	Byte     Type = "byte"

	Data Type = "data"

	//Special types.
	Newline Type = "newline"
	Switch  Type = "switch"
)

//NewMetatype returns a metatype with the specified type value as an expression.
func NewMetatype(t Type) Expression {
	return Expression{Metatype, nil, t}
}
