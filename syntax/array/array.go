package Array

import "github.com/qlova/script/compiler"
import . "github.com/qlova/script"
import "github.com/qlova/i/syntax/errors"

var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) compiler.Type {
	if array, ok := a.(Array); ok && symbol == "[" {
		
		if n, ok := b.(Number); !ok {
			c.RaiseError(errors.IndexError())
		} else {
			c.Expecting("]")

			return c.Index(array, n)
		}
		
	}
	return nil
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if c.GetVariable(c.Token()).Defined {
			var name = c.Token()
			var variable = c.Variable(name)

			if array, ok := variable.(Array); ok && c.Peek() == "[" {
				c.Scan()
				
				var index = c.ScanType(c.Number()).(Number)
				c.Expecting("]")
				
				c.Expecting("=")
			
				c.Modify(array, index, c.ScanType(array.SubType()))
				return true
			}
		}
		return false
	},
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {
		
		if c.Token() == "#" {
			var array = c.Shunt(c.Expression(), 5) //Custom operator precidence.
			
			switch array.(type) {
				case Array:
					return c.Length(array)
					
				case List:
					return c.Length(array)
				
				default:
					c.RaiseError(compiler.Translatable{
						compiler.English: "Cannot take the length of type "+array.Name(),
					})
			}
		}
		
		if c.Token() == "array" {
			c.Expecting("(")
			
			//Special case, we want to convert a list literal into an array literal!
			if c.Peek() == "[" {
				c.Scan()
				
				var expression = c.ScanExpression()
				var elements []Type
				
				elements = append(elements, expression)
				
				for {
					if c.Peek() == "," {
						
						c.Scan()
						elements = append(elements, c.ScanType(elements[0]))
						
					} else {
						c.Expecting("]")
						c.Expecting(")")
						break
					}
				}
				
				return c.Array(elements)
			}
			
			var length = c.ScanType(c.Number()).(Number)
			c.Expecting(")")
			
			return c.Array(c.Number(), int(length.Literal.Int64()))
		}
		
		/*if c.Token() == "[" {
			
			var elements []script.Type
			var expression = c.ScanExpression()
			elements = append(elements, expression.Type)
			
			for {
				if c.Peek() == "," {
					
					var extra = c.ScanExpression()
					if !extra.Type.Equals(expression.Type) {
						c.RaiseError(errors.ExpectingType(expression, extra))
					}
					elements = append(elements, extra.Type)
					
				} else {
					c.Expecting("]")
				}
			}
			
			return compiler.ScriptType(c.Script.LiteralArray(expression.Type, elements...))
		}*/
		
		return nil
	},
}
