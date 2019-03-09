package For

import "github.com/qlova/script/compiler"
import . "github.com/qlova/script"

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "for",
	},
	
	OnScan: func(c *compiler.Compiler) {
		
		//For element in array. TODO other for loops.
		var variable = c.Scan()
		c.Expecting("in")
		var array = c.ScanExpression()
		
		if _, ok := array.(List); !ok {
			c.RaiseError(compiler.Translatable{
				compiler.English: "Can only iterate over lists!",
			})
		}
		
		c.ForEach("i", variable, array.(List), func(i Number, v Type, q *Script) {
			c.GainScope()
			c.SetVariable("i", i)
			c.SetVariable(variable, v)
			if c.ScanIf(":") {
				c.CompileBlock("", "\n")
			} else {
				c.CompileBlock("", "}")
			}
			c.LoseScope()
		})
	},
}

