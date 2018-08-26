package ilang

import "github.com/qlova/script/compiler"

import (
	"github.com/qlova/i/syntax/main"
	"github.com/qlova/i/syntax/read"
	"github.com/qlova/i/syntax/print"
	"github.com/qlova/i/syntax/write"
	
	"github.com/qlova/i/syntax/for"
	
	"github.com/qlova/i/syntax/variables"
	
	"github.com/qlova/i/syntax/string"
	"github.com/qlova/i/syntax/array"
	"github.com/qlova/i/syntax/list"
	"github.com/qlova/i/syntax/number"
	"github.com/qlova/i/syntax/symbol"
	
	"github.com/qlova/i/syntax/concept"
)

var Syntax = compiler.NewSyntax("ilang")

func init() {
	Syntax.RegisterStatement(Main.Statement)
	Syntax.RegisterStatement(Concept.Statement)
	Syntax.RegisterStatement(Print.Statement)
	Syntax.RegisterStatement(Write.Statement)
	Syntax.RegisterStatement(Variables.Statement)
	Syntax.RegisterStatement(For.Statement)
	Syntax.RegisterStatement(Concept.Return)
	
	Syntax.RegisterExpression(Read.Expression)
	Syntax.RegisterExpression(Variables.Expression)
	Syntax.RegisterExpression(String.Expression)
	Syntax.RegisterExpression(Array.Expression)
	Syntax.RegisterExpression(List.Expression)
	Syntax.RegisterExpression(Number.Expression)
	Syntax.RegisterExpression(Symbol.Expression)
	Syntax.RegisterExpression(Concept.Expression)
	
	
	Syntax.RegisterShunt(Number.Shunt)
	Syntax.RegisterShunt(Array.Shunt)
	Syntax.RegisterShunt(String.Shunt)
	
	Syntax.RegisterOperator(")", -1)
	Syntax.RegisterOperator("]", -1)
	Syntax.RegisterOperator("or", 0)
	Syntax.RegisterOperator("and", 1)
	Syntax.RegisterOperator("=", 2)
	Syntax.RegisterOperator(">", 2)
	Syntax.RegisterOperator("<", 2)
	Syntax.RegisterOperator("+", 3)
	Syntax.RegisterOperator("-", 3)
	Syntax.RegisterOperator("*", 4)
	Syntax.RegisterOperator("^", 5)
	Syntax.RegisterOperator("/", 4)
	Syntax.RegisterOperator("%", 4)

	Syntax.RegisterOperator("(", 6)
	Syntax.RegisterOperator("[", 6)
	Syntax.RegisterOperator(".", 7)
}
