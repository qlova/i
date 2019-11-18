package compiler

import "io"

//CompileStatement compiles the next statement.
func (c *Compiler) CompileStatement() (err error) {
	var token = c.Scan()

	if len(token) > 2 && token[0] == '/' && token[1] == '/' {
		return c.compileComments(token)
	}

	switch token.String() {
	case "main":
		c.Main(func() {
			err = c.CompileBlock()
		})
	case "print":
		if !c.ScanIf('(') {
			return c.NewError("expecting (")
		}
		expression, err := c.ScanExpression()
		if err != nil {
			return err
		}
		c.Discard(c.Send(nil, expression))
		if !c.ScanIf(')') {
			return c.NewError("expecting )")
		}
		c.Discard(c.Send(nil, c.String("\n")))
	case "":
		return io.EOF
	case "\n":
		return nil
	default:
		err = c.Undefined(s("statement: " + token.String()))
	}

	return err
}
