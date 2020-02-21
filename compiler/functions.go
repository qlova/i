package compiler

import (
	"github.com/qlova/usm"
)

//Constraint is used to constrain arguments.
type Constraint struct {
	Type
}

//Argument is a concept or function argument.
type Argument struct {
	Constraint

	Name    string
	Owned   bool
	Default Expression
}

//Function is a 'i' function.
type Function struct {
	usm.Label

	Args   []Argument
	Throws bool
}

//ScanCall scans a function call.
func (c *Compiler) ScanCall() (args []Expression, err error) {
	if !c.ScanIf('(') {
		return nil, c.NewError("Expecting '('")
	}

	for {
		arg, err := c.ScanExpression()
		if err != nil {
			return nil, err
		}

		args = append(args, arg)

		if c.ScanIf(',') {
			continue
		}

		if !c.ScanIf(')') {
			return nil, c.NewError("Expecting ')'")
		}

		break
	}

	return
}
