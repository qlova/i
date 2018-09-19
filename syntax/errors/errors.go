package errors

import "github.com/qlova/script/compiler"
import "strings"

type Nameable interface {
	Name() string
}

func AssignmentMismatch(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Cannot assign a "+a.Name()+" value to a variable of type "+b.Name(),
	}
}

func ExpectingType(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Expecting a value of type "+a.Name()+" instead got a value of type "+b.Name(),
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
		compiler.English: a.Name()+" is not a numeric type!",
	}
}

func Single(a compiler.Type, symbol string, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "The relationship "+a.Name()+symbol+b.Name()+" is not defined!",
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

func Inconsistent(a, b Nameable) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "The usage here of the '"+a.Name()+"' type is inconsistent with the\n use of the '"+b.Name()+"' type before this!", 
	}
}

func NoSuchElement(a string, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "No such element '"+a+"' in type '"+b.Name()+"'!", 
	}
}
