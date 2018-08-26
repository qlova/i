package Array

import "github.com/qlova/script/compiler"
import "github.com/qlova/script"
import "github.com/qlova/i/syntax/errors"

var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) *compiler.Type {
	if array, ok := a.Type.(script.Array); ok && symbol == "[" {
		
		if n, ok := b.Type.(script.Number); !ok {
			c.RaiseError(errors.IndexError())
		} else {
			c.Expecting("]")
			return compiler.ScriptType(c.Script.IndexArray(array, n))
		}
		
	}
	return nil
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		if c.Token() == "#" {
			var array = c.Shunt(c.Expression(), 5) //Custom operator precidence.
			
			switch array.Type.(type) {
				case script.Array:
					return compiler.ScriptType(c.Script.LengthArray(array.Type.(script.Array)))
				
				default:
					c.RaiseError(compiler.Translatable{
						compiler.English: "Cannot take the length of type "+array.Type.Name(),
					})
			}
		}
		
		if c.Token() == "array" {
			c.Expecting("(")
			
			//Special case, we want to convert a list literal into an array literal!
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
				
				return compiler.ScriptType(c.Script.LiteralArray(elements[0], elements...))
			}
			
			var expression = c.ScanExpression()
			c.Expecting(")")
			
			n, err := c.Script.ToArray(script.Number(""), expression.Type)
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: err.Error(),
				})
			}
			
			return compiler.ScriptType(n)
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
