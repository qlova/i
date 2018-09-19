package String

import "strconv"
import "github.com/qlova/script/compiler"
import . "github.com/qlova/script"



var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) compiler.Type {
	
	//Concatenate strings.
	if StringA, ok := a.(String); ok && symbol == "+" {
		
		if StringB, ok := b.(String); ok {
			return c.Join(StringA, StringB)
		}
		
	}
	
	
	return nil
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {
		
		//Cast to string.
		if c.Token() == "string" {
			c.Expecting("(")
			var expression = c.ScanExpression()
			c.Expecting(")")
			
			return c.ToString(expression)
		}
		
		
		//String Literal.
		if c.Token()[0] == '"' {
			
			str, err := strconv.Unquote(c.Token())
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Invalid String!",
				})
			}
			
			return c.Script.String(str)
		}
		
		return nil
	},
}
