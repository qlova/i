package List

import "github.com/qlova/script/compiler"
import . "github.com/qlova/script"
import "github.com/qlova/i/syntax/errors"

var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) compiler.Type {
	if list, ok := a.(List); ok && symbol == "[" {
		
		if n, ok := b.(Number); !ok {
			c.RaiseError(errors.IndexError())
		} else {
			c.Expecting("]")
			return c.Index(list, n)
		}
		
	}
	return nil
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {		
		if c.Token() == "[" {
			
			var elements []Type
			var expression = c.ScanExpression()
			elements = append(elements, expression)
			
			for {
				if c.ScanIf(",") {
					
					var extra = c.ScanType(expression)
					elements = append(elements, extra)
					
				} else {
					c.Expecting("]")
					break
				}
			}
			
			return c.List(elements...))
		}
		
		return nil
	},
}
