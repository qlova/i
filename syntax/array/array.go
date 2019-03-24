package Array

import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"
import "github.com/qlova/i/syntax/errors"


var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) compiler.Type {
	
	if a.Value().IsArray() && symbol == "[" {

		if !b.Value().Is(c.Int()) {
			c.RaiseError(errors.IndexError())
		} else {
			c.Expecting("]")

			return a.Value().Array().Index(b.Value().Int())
		}

	}
	
	return nil
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if c.GetVariable(c.Token()).Defined {
			
			var name = c.Token()
			var variable = c.Variable(name)

			if variable.Value().IsArray() && c.Peek() == "[" {
				c.Scan()
				
				var index = c.ScanType(c.Int()).(Int)
				c.Expecting("]")
				
				c.Expecting("=")
			
				variable.Value().Array().Modify(index, c.ScanType(variable.Value().Array().Subtype()))
				return true
			}
			
		}
		return false
	},
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {
		
		if c.Token() == "#" {
			var collection = c.Shunt(c.Expression(), 5) //Custom operator precidence.
			
			if collection.Value().IsArray() {
				return c.Len(collection)
			} else {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Cannot take the length of type "+collection.LanguageType().Name(),
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
				
				return c.Array(elements...)
			}
			
			var length = c.ScanType(c.Int()).(Int)
			c.Expecting(")")
			
			return c.Int().Array(int(length.Literal().Int64()))
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
