package Catch

import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"
	

var Expression = compiler.Expression{
	Name: compiler.Translatable{
		compiler.English: "catch",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Expecting("(")
		c.Expecting(")")
		return c.Catch()
	},
}

var ErrorExpression = compiler.Expression{
	Name: compiler.Translatable{
		compiler.English: "error",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		return c.Error()
	},
}
