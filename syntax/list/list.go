package List

import "github.com/qlova/script/compiler"
import "github.com/qlova/script"
import "github.com/qlova/i/syntax/errors"

var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) *compiler.Type {
	if list, ok := a.Type.(script.List); ok && symbol == "[" {
		
		if n, ok := b.Type.(script.Number); !ok {
			c.RaiseError(errors.IndexError())
		} else {
			c.Expecting("]")
			return compiler.ScriptType(c.Script.IndexList(list, n))
		}
		
	}
	return nil
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		if c.Token() == "#" {
			var list = c.Shunt(c.Expression(), 5) //Custom operator precidence.
			
			switch list.Type.(type) {
				case script.List:
					return compiler.ScriptType(c.Script.LengthList(list.Type.(script.List)))
				
				default:
					c.RaiseError(compiler.Translatable{
						compiler.English: "Cannot take the length of type "+list.Type.Name(),
					})
			}
		}
		
		if c.Token() == "list" {
			c.Expecting("(")
			
			//Special case, we want to convert a list literal into an list literal!
			if c.Peek() == "[" {
				c.Scan()
				
				var expression = c.ScanExpression()
				var elements []script.Type
				
				elements = append(elements, expression.Type)
				
				for {
					if c.Peek() == "," {
						
						c.Scan()
						expression = c.ScanExpression()
						
						if !expression.Type.Equals(elements[0]) {
							c.RaiseError(errors.Inconsistent(expression, elements[0]))
						}
						
						elements = append(elements, expression.Type)
						
					} else {
						c.Expecting("]")
						c.Expecting(")")
						break
					}
				}
				
				return compiler.ScriptType(c.Script.LiteralList(elements[0], elements...))
			}
			
			var expression = c.ScanExpression()
			c.Expecting(")")
			
			n, err := c.Script.ToList(script.Number(""), expression.Type)
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: err.Error(),
				})
			}
			
			return compiler.ScriptType(n)
		}
		
		if c.Token() == "[" {
			
			var elements []script.Type
			var expression = c.ScanExpression()
			elements = append(elements, expression.Type)
			
			for {
				if c.ScanIf(",") {
					
					var extra = c.ScanExpression()
					if !extra.Type.Equals(expression.Type) {
						c.RaiseError(errors.ExpectingType(expression, extra))
					}
					elements = append(elements, extra.Type)
					
				} else {
					c.Expecting("]")
					break
				}
			}
			
			return compiler.ScriptType(c.Script.LiteralList(expression.Type, elements...))
		}
		
		return nil
	},
}
