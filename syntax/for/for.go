package For

import "github.com/qlova/script/compiler"

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "for",
	},
	
	OnScan: func(c *compiler.Compiler) {
		
		//For element in array. TODO other for loops.
		var variable = c.Scan()
		c.Expecting("in")
		var array = c.ScanExpression()
		
		if !array.Value().IsList() {
			c.RaiseError(compiler.Translatable{
				compiler.English: "Can only iterate over lists!",
			})
		} else {
			array.Value().List().ForEach(func() {
				c.GainScope()
				c.SetVariable("i", c.Index())
				c.SetVariable(variable, c.Value())
				if c.ScanIf(":") {
					c.CompileBlock("", "\n")
				} else {
					c.CompileBlock("", "}")
				}
				c.LoseScope()
			}, "i", variable)
		}
	},
}

