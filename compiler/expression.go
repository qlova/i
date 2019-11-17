package compiler

import (
	"strconv"

	"github.com/qlova/script/language"
)

//ScanExpression scans an expression.
func (c *Compiler) ScanExpression() (language.Value, error) {
	var token = c.Scan()

	//String
	if token[0] == '"' {
		var literal, err = strconv.Unquote(token.String())
		if err != nil {
			return nil, c.NewError("invalid string (", err, ")")
		}
		return c.Literal().String(literal), nil
	}

	return nil, c.NewError("invalid expression")
}
