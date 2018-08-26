package String

import "strconv"
import "github.com/qlova/script/compiler"
import "github.com/qlova/script"

var Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) *compiler.Type {
	if StringA, ok := a.Type.(script.String); ok && symbol == "+" {
		
		if StringB, ok := b.Type.(script.String); ok {
			return compiler.ScriptType(c.Script.AddStrings(StringA, StringB))
		}
		
	}
	return nil
}

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		if c.Token() == "string" {
			c.Expecting("(")
			var expression = c.ScanExpression()
			c.Expecting(")")
			
			n, err := c.Script.ToString(expression.Type)
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: err.Error(),
				})
			}
			
			return compiler.ScriptType(n)
		}
		
		if c.Token()[0] == '"' {
			
			_, err := strconv.Unquote(c.Token())
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Invalid String!",
				})
			}
			
			return compiler.ScriptType(c.Script.LiteralString(c.Token()))
		}
		
		return nil
	},
}
