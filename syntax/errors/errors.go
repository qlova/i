package errors

import "github.com/qlova/script/compiler"
import "strings"

type Stringable interface {
	String() string
}

func AssignmentMismatch(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Cannot assign a "+a.String()+" value to a variable of type "+b.String(),
	}
}

func ExpectingType(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Expecting a value of type "+a.String()+" instead got a value of type "+b.String(),
	}
}

func IndexError() compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Arrays can only be indexed with numbers!",
	}
}

func UnknownType(a string) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Unknown Type: "+a,
	}
}

func MustBeNumeric(a compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: a.String()+" is not a numeric type!",
	}
}

func Single(a compiler.Type, symbol string, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "The relationship "+a.String()+symbol+b.String()+" is not defined!",
	}
}

func IsInvalidName(name string) bool {
	return strings.Contains(name, "_")
}


func InvalidName(name string) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Invalid name: "+name,
	}
}

func Inconsistent(a, b Stringable) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "The usage here of the '"+a.String()+"' type is inconsistent with the\n use of the '"+b.String()+"' type before this!", 
	}
}

func NoSuchElement(a string, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "No such element '"+a+"' in type '"+b.String()+"'!", 
	}
}
