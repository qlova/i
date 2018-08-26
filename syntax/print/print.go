package Print

import "github.com/qlova/script/compiler"
	

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "print",
	},
	
	OnScan: func(c *compiler.Compiler) {
		c.Expecting("(")
		var expression = c.ScanExpression()
		
		s, err := c.Script.ToString(expression.Type)
		if err != nil {
			c.RaiseError(compiler.Translatable{
				compiler.English: err.Error(),
			})
		}

		c.Script.Print(s)
		
		for {
			if c.Peek() == "," {
				c.Scan()
				//TODO deal with other types.
				
				var expression = c.ScanExpression()
		
				s, err := c.Script.ToString(expression.Type)
				if err != nil {
					c.RaiseError(compiler.Translatable{
						compiler.English: err.Error(),
					})
				}

				c.Script.Print(s)
				
			} else {
				c.Expecting(")")
				break
			}
		}
		
		c.Script.Print(c.Script.LiteralString(`"\n"`))
		
		//Convert to String.
		
		

	},
}
