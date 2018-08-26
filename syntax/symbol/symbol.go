package Symbol

import "github.com/qlova/script/compiler"

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		if c.Token() == "symbol" {
			c.Expecting("(")
			var expression = c.ScanExpression()
			c.Expecting(")")
			
			n, err := c.Script.ToSymbol(expression.Type)
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: err.Error(),
				})
			}
			
			return compiler.ScriptType(n)
		}
		
		if c.Token()[0] == '\'' {
			return compiler.ScriptType(c.Script.LiteralSymbol(c.Token()))
		}
		
		return nil
	},
}
