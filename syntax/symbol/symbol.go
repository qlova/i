package Symbol

import "strconv"
import "github.com/qlova/script/compiler"

var Expression = compiler.Expression{
	Detect: func(c *compiler.Compiler) compiler.Type {
		
		if c.Token() == "symbol" {
			c.Expecting("(")
			var expression = c.ScanExpression()
			c.Expecting(")")
			
			return expression.Value().Rune()
		}
		
		if c.Token()[0] == '\'' {
			
			v, _, _, err := strconv.UnquoteChar(c.Token()[1:], '\'')
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Invalid symbol: "+c.Token(),
				})
			}
			
			return c.Rune(v)
		}
		
		return nil
	},
}
