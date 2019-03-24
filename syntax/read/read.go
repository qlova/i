package Read


import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"

var Name = compiler.Translatable{
	compiler.English: "read",
	compiler.Maori: "rīti",
}

var Expression = compiler.Expression {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Expecting("(")
		
		if c.ScanIf(")") {
			return c.Read(c.Rune('\n'))
		}
		
		var argument = c.ScanExpression()
		c.Expecting(")")
		
		switch argument.(type) {
			case Rune:
				
				return c.Read(argument)
			
			default:
				c.Unimplemented()
				
		}
		
		return nil
	},
}

/*var Statement = compiler.Statement {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) {
		c.DropType(Expression.OnScan(c))
	},
}*/
