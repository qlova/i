package compiler

import (
	"io"
)

//CompileStatement compiles the next statement.
func (c *Compiler) CompileStatement() (err error) {
	var token = c.Scan()

	if len(token) > 2 && token[0] == '/' && token[1] == '/' {
		return c.compileComments(token)
	}

	defer func() {
		//Error handling.
		if c.Peek().Is(";") {
			c.Scan()

			if err != nil {
				return
			}

			var directive = c.Scan()

			switch directive.String() {
			case "throw":
				var expr Expression
				expr, err = c.ScanExpression()

				c.Return(expr.Value)

			default:
				err = c.NewError("invalid error handling directive")
			}

		}
	}()

	switch token.String() {
	case "main":
		c.Main(func() {
			err = c.CompileBlock()
		})
	case "for":
		return c.ScanLoop()
	case "if":
		return c.ScanBranch()

	case "usm":
		return c.ScanInlineAssembly()

	case "return":
		var expr Expression
		expr, err = c.ScanExpression()

		c.Return(expr.Value)
		return err

	case "print":
		if !c.ScanIf('(') {
			return c.NewError("expecting (")
		}
		expression, err := c.ScanExpression()
		if err != nil {
			return err
		}
		c.Discard(c.Send(nil, expression.Value))
		if !c.ScanIf(')') {
			return c.NewError("expecting )")
		}
		c.Discard(c.Send(nil, c.String("\n")))
	case "":
		return io.EOF
	case "\n":
	default:

		//Implicit concept or type definition.
		if len(c.Scopes) == 0 {
			var name = token.String()
			var peek = c.Peek()
			if peek.Is("(") {
				concept, err := c.ScanConcept()
				if err != nil {
					return err
				}

				c.Concepts[name] = concept
				return nil
			}
		}

		//Variable declaration.
		if c.Peek().Is("$") {
			c.Scan()

			var name = token.String()
			if !c.ScanIf('=') {
				return c.NewError("Expecting =")
			}

			var value, err = c.ScanExpression()
			if err != nil {
				return err
			}

			c.DefineVariable(name, Variable{
				Type:     value.Type,
				Register: c.Var(value.Value),
			})

			return nil
		}

		err = c.Undefined(s("statement: " + token.String()))
	}

	return err
}
