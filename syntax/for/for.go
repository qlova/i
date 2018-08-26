package For

import "github.com/qlova/script/compiler"
import qlova "github.com/qlova/script"

var Statement = compiler.Statement{
	Name: compiler.Translatable{
		compiler.English: "for",
	},
	
	OnScan: func(c *compiler.Compiler) {
		
		//For element in array. TODO other for loops.
		var variable = c.Scan()
		c.Expecting("in")
		var array = c.ScanExpression()
		
		if _, ok := array.Type.(qlova.List); !ok {
			c.RaiseError(compiler.Translatable{
				compiler.English: "Can only iterate over lists!",
			})
		}
		
		c.Script.ForEachList("", variable, array.Type.(qlova.List), func(q *qlova.Script) {
			c.GainScope()
			c.SetVariable(variable, *compiler.ScriptType(array.Type.(qlova.List).Subtype()))
			if c.ScanIf(":") {
				c.CompileBlock("", "\n")
			} else {
				c.CompileBlock("", "}")
			}
			c.LoseScope()
		})
	},
}

