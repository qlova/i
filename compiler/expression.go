package compiler

import (
	"math/big"
	"strconv"

	"github.com/qlova/usm"
)

//Expression type.
type Expression struct {
	Type
	usm.Value

	Literal interface{}
}

//Equals returns a new expression of equality between expressions.
//Keep in mind that not all type can be tested for equality.
//They should be the same type.
func (a Expression) Equals(b Expression, c *Compiler) (Expression, error) {
	if a.Type != b.Type {
		return Expression{}, c.NewError("types not equal")
	}

	switch a.Type {
	case Integer:
		return Expression{a.Type, c.Equals(a.Value, b.Value), nil}, nil
	case Metatype:
		b := a.Literal == b.Literal
		return Expression{Bit, c.Bit(b), b}, nil
	default:
		return Expression{}, c.NewError(a.Type, " types cannot be compared.")
	}
}

//ScanExpression returns the next expression or an error.
func (c *Compiler) ScanExpression() (Expression, error) {
	var expression, err = c.scanExpression()
	if err != nil {
		return expression, err
	}

	return c.Shunt(expression, 0)
}

//scanExpression scans an expression.
func (c *Compiler) scanExpression() (Expression, error) {
	var token = c.Scan()

	if len(token) == 0 {
		return Expression{}, c.NewError("Empty Expression")
	}

	//String
	if token[0] == '"' {
		var literal, err = strconv.Unquote(token.String())
		if err != nil {
			return Expression{}, c.NewError("invalid string (", err, ")")
		}
		return Expression{String, c.String(literal), literal}, nil
	}

	//Symbol
	if token[0] == '\'' {
		var string = token.String()
		var literal, _, _, err = strconv.UnquoteChar(string[1:len(string)-1], '\'')
		if err != nil {
			return Expression{}, c.NewError("invalid string (", err, ")")
		}
		var big = big.NewInt(int64(literal))
		return Expression{Symbol, c.Number(big), big}, nil
	}

	//Integer expression.
	if _, err := strconv.Atoi(string(token)); err == nil {
		var big, success = big.NewInt(0).SetString(string(token), 10)
		if !success {
			return Expression{}, c.NewError("invalid integer")
		}
		return Expression{Integer, c.Number(big), big}, nil
	}

	if token.Is("\n") {
		return Expression{Newline, nil, nil}, nil
	}

	var s = token.String()

	//Type
	if s == "type" {
		if !c.ScanIf('(') {
			return Expression{}, c.NewError("expecting '('")
		}
		expression, err := c.ScanExpression()
		if err != nil {
			return Expression{}, err
		}
		if !c.ScanIf(')') {
			return Expression{}, c.NewError("expecting ')'")
		}

		if expression.Type == "" {
			return Expression{}, c.NewError("unknown type")
		}

		return NewMetatype(expression.Type), nil
	}

	//Variables
	if v, ok := c.GetVariable(s); ok {
		return Expression{v.Type, c.Get(v.Register), nil}, nil
	}

	//Concepts
	if concept, ok := c.Concepts[s]; ok {
		return c.CallConcept(concept)
	}

	//Metatypes.
	if !c.Peek().Is("(") {
		for _, T := range Types {
			if s == string(T) {
				return NewMetatype(T), nil
			}
		}
	}

	//Constructors / Casting.
	if s == string(Data) {
		var args, err = c.ScanCall()
		if err != nil {
			return Expression{}, err
		}

		if len(args) != 1 {
			return Expression{}, c.NewError("Expecting 1 argument to type.")
		}

		if args[0].Type != Integer {
			return Expression{}, c.NewError("Cannot cast to data from " + args[0].Type)
		}

		return Expression{Data, c.Create(args[0].Value), nil}, nil
	}

	if s == string(Byte) {
		var args, err = c.ScanCall()
		if err != nil {
			return Expression{}, err
		}

		if len(args) != 1 {
			return Expression{}, c.NewError("Expecting 1 argument to type.")
		}

		if args[0].Type != Symbol {
			return Expression{}, c.NewError("Cannot cast to byte from " + args[0].Type)
		}

		return Expression{Byte, args[0].Value, nil}, nil
	}

	return Expression{}, c.NewError("invalid expression: ", s)
}
