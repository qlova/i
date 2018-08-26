package Read

import "github.com/qlova/script"
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
			return *compiler.ScriptType(c.Script.ReadSymbol(c.Script.LiteralSymbol(`'\n'`)))
		}
		
		var argument = c.ScanExpression()
		c.Expecting(")")
		
		switch argument.Type.(type) {
			case script.Symbol:
				
				return *compiler.ScriptType(
					c.Script.ReadSymbol(argument.Type.(script.Symbol)),
				)
			
			default:
				c.Unimplemented()
				
		}
		
		return compiler.Type{}
	},
}

/*var Statement = compiler.Statement {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) {
		c.DropType(Expression.OnScan(c))
	},
}*/
